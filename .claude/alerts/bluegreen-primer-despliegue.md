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

## Importante: probar una reversion de verdad

Todavia no se comprobo en produccion que un deploy con la imagen rota deje el
color viejo sirviendo. Vale la pena forzarlo una vez (por ejemplo desplegando una
imagen que no arranque) y confirmar que el sitio nunca respondio 502.

## Deseable

- Blue-green para `front-website` y para el stack de testing.
- Que el deploy de nginx tampoco corte (hoy recrea el contenedor).

## Criterio para cerrar

Los dos urgentes hechos y verificados: primer deploy en verde con el sitio
arriba, y `/probability/back-central` recibiendo eventos despues de un cambio de
color.
