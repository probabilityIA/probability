# Blue-green: pasos manuales del primer despliegue

Fecha: 2026-09-02

El codigo de blue-green ya esta en `main`, pero el primer rollout tiene dos
cosas que el CI no puede hacer solo. Hasta que esten hechas, esta alerta sigue
abierta.

## Urgente 1: el orden del primer deploy

nginx tiene que quedar con la imagen nueva **antes** de que cambie el color, o
va a seguir buscando `back-central` / `front-central`, que ya no existen como
servicio, y el sitio se cae entero.

`bg_ensure_nginx` recrea nginx solo si le falta el montaje
`/etc/nginx/upstreams`, asi que el deploy de back o de front se autocorrige.
Aun asi, el orden seguro es:

1. Merge a `main` (dispara el workflow de nginx por `infra/nginx/**`).
2. Esperar a que **termine** el deploy de nginx.
3. Recien ahi dejar correr el de backend y el de frontend, **uno a la vez**
   (el `t4g.small` tiene ~770 MB libres y sin swap).
4. Verificar:

```bash
aws ssm send-command --profile probability --region us-east-1 \
  --instance-ids i-0f3284d2a87127e57 --document-name AWS-RunShellScript \
  --parameters 'commands=["docker ps --format \"{{.Names}}\"","cat /home/ubuntu/probability/infra/compose-prod/nginx-upstreams/active.conf","curl -sf -o /dev/null -w \"%{http_code}\\n\" https://www.probabilityia.com.co/api/v1/health"]' \
  --query 'Command.CommandId' --output text
```

Los contenedores `central_reserve_prod` y `frontend_prod` (sin sufijo de color)
los borra el propio deploy con `bg_drop_legacy`.

## Urgente 2: el envio de logs a CloudWatch se rompe

`probability-logtail@central_reserve_prod.service` y
`probability-logtail@frontend_prod.service` siguen a un contenedor por nombre
exacto. Con blue-green el nombre cambia de color en cada deploy, asi que esas
unidades quedan apuntando a un contenedor que ya no existe y **dejan de subir
lineas de error a CloudWatch, en silencio**.

Arreglo: que el script resuelva el contenedor por prefijo. En el EC2:

```bash
sudo tee /usr/local/bin/probability-logtail.sh >/dev/null <<'SH'
#!/bin/bash
# Sigue los logs de un contenedor por PREFIJO de nombre (blue-green cambia el
# sufijo _blue/_green en cada deploy) y deja pasar solo lineas de error.
PREFIX="$1"
OUT="/var/log/probability/${PREFIX}.log"
mkdir -p /var/log/probability
while true; do
  NAME=$(docker ps --filter "name=^${PREFIX}" --format '{{.Names}}' | head -1)
  if [ -z "$NAME" ]; then sleep 10; continue; fi
  docker logs -f --tail 0 "$NAME" 2>&1 \
    | grep --line-buffered -E 'ERROR|WARN|FATAL|PANIC|panic:|Exception|error|warn|fatal' \
    >> "$OUT"
  sleep 5
done
SH
sudo chmod +x /usr/local/bin/probability-logtail.sh

sudo systemctl disable --now probability-logtail@central_reserve_prod.service
sudo systemctl disable --now probability-logtail@frontend_prod.service
sudo systemctl enable --now probability-logtail@central_reserve_prod.service
sudo systemctl enable --now probability-logtail@frontend_prod.service
```

El prefijo sigue siendo `central_reserve_prod` / `frontend_prod`, asi que **no
hay que tocar** el `collect_list` del CloudWatch agent ni los log groups.

Verificar que llegan lineas:

```bash
ls -l /var/log/probability/
aws logs describe-log-streams --profile probability --region us-east-1 \
  --log-group-name /probability/back-central --order-by LastEventTime --descending --max-items 1
```

## Resuelto 2026-09-02 23:15: el switch quedo validado

Dos deploys de prueba del backend (commits `575dddab` y `2732077d`), con un
sondeo de 1 req/s a `/health` y `/` desde afuera durante toda la ventana:

| Deploy | Switch | Muestras | Fallos |
|---|---|---|---|
| 1 | -> green | 143 | 1 (`web=000`, con `/health` en 200 el mismo segundo) |
| 2 | green -> blue | 148 | **0** |

El segundo leyo bien el color activo (`green -> blue`), levanto el nuevo,
confirmo `/health` en 3 s, movio el upstream con `nginx -s reload`, drenó 15 s y
apago el viejo. **Cero caida medida.** Antes cada deploy eran 20-35 s duros.

Queda sin explicar del todo por que el primer deploy leyo `blue` estando green
activo: existia un contenedor `central_reserve_prod_blue` (lo retiro el propio
script) que no aparecia en `docker ps` cuando revise antes. No se reprodujo en el
segundo. Si vuelve a pasar, mirar `docker ps -a` completo antes de concluir.

## Pendiente: probar una reversion de verdad

Todavia no se comprobo que un deploy con la imagen rota deje el color viejo
sirviendo. Vale la pena forzarlo una vez (desplegando una imagen que no arranque)
y confirmar que el sitio nunca respondio 502.

## Deseable

- Blue-green para `front-website` y para el stack de testing.
- Que el deploy de nginx tampoco corte (hoy recrea el contenedor).

## Criterio para cerrar

Los dos urgentes hechos y verificados: primer deploy en verde con el sitio
arriba, y `/probability/back-central` recibiendo eventos despues de un cambio de
color.

## Lo que paso en el primer rollout (2026-09-02 22:04 UTC)

Los tres workflows (nginx, backend, frontend) salieron a la vez con el push y el
sitio se cayo ~10 minutos. Secuencia real:

1. nginx desplego bien y su entrypoint, **como root**, creo `active.conf`
   apuntando a `blue`.
2. El deploy del backend levanto `central_reserve_prod_green`, paso el health
   check, y al reescribir `active.conf` murio con
   `Permission denied` (el archivo era de root, el deploy es `ubuntu`).
3. Con `set -e` el script aborto ahi: green arriba y sin trafico, `active.conf`
   apuntando a `back-central-blue`, que **no existia** (el contenedor viejo se
   llamaba `central_reserve_prod`, sin color).
4. nginx no resuelve un upstream inexistente -> `exit 1` -> bucle de reinicio ->
   502 en todo el sitio.
5. El deploy del frontend fallo en cascada: su `startup.sh` espera al backend a
   traves de nginx, que estaba en el bucle.

Recuperacion manual: levantar `front-central-green`, escribir `active.conf` con
green/green, reiniciar nginx, borrar los contenedores legacy.

Correcciones ya en el codigo:

- `bg_init` hace `sudo chown -R ubuntu:ubuntu` del directorio de upstreams.
- El entrypoint de nginx hace `chown 1000:1000` del archivo que crea.
- `bg_active_color` ya no cree ciegamente al archivo: si el color que nombra no
  tiene contenedor corriendo, usa el que si lo tiene, y `bg_init` reescribe el
  archivo. Un `active.conf` apuntando a un color muerto ahora se corrige solo en
  el siguiente deploy en vez de tumbar el sitio.

Estado tras la recuperacion: backend y frontend sirviendo en **green**, legacy
borrado, RAM en 772 MB libres. El proximo deploy de cada uno ira a blue.
