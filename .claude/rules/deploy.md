# Deploy y CI/CD

## Workflows

Push a `main` dispara CI/CD en `.github/workflows/`. Build ARM64 -> ECR (4 tags) -> **SSM** -> deploy con 3 retries.
No entrar al servidor para verificar deploys: confiar en GitHub Actions.

### Deploy por SSM (desde 2026-08-21, ya no hay SSH)

El build sigue en GitHub Actions y empuja a ECR por HTTPS. Lo que cambio es el
transporte del deploy: `scp`/`ssh` con `probability.pem` fueron reemplazados por
S3 + `ssm send-command`.

- `.github/scripts/ssm-deploy.sh` - helper compartido: sube los artefactos a
  `s3://probability-deploy-artifacts/deploy/<run_id>/`, ejecuta el script remoto
  como usuario `ubuntu`, transmite la salida, propaga el exit code y limpia S3.
- `.github/scripts/deploy-<servicio>.sh` - la logica de cada deploy. Recibe
  `VERSION_FULL`, `VERSION_SHORT` y `DEPLOY_DIR` (donde quedaron los artefactos).

Para tocar un deploy se edita el `.sh`, no el YAML. El job del workflow solo
prepara `artifacts/` y llama al helper.

Gotchas que ya costaron tiempo:

- El comando remoto lo interpreta **dash**, no bash: `set -euo pipefail` falla
  con "Illegal option -o pipefail". Usar `set -e`.
- El script hay que mandarlo como **array de lineas** en `Parameters.commands`
  (`$cmd | split("\n")` con jq). Un solo string con `\n` no se interpreta.
- `send-command` corre como root; los scripts esperan `ubuntu`, por eso el helper
  usa `runuser -l ubuntu -c`.

Secrets ya sin uso: `EC2_SSH_KEY`, `EC2_HOST`, `EC2_USER`.

| Workflow | Paths | Puerto prod |
|----------|-------|-------------|
| Backend  | `back/central/**`, `back/migration/**` | 3050 blue / 3060 green |
| Frontend | `front/central/**` | 8080 blue / 8090 green |
| Website  | `front/website/**` | 8081 |
| Nginx    | `infra/nginx/**` | 80/443 |

## Blue-green (back-central y front-central)

Desde 2026-09-02 el backend y el frontend no se apagan para desplegar. Cada uno
tiene **dos colores** definidos en el compose (`back-central-blue` /
`back-central-green`, `front-central-blue` / `front-central-green`) y solo uno
recibe trafico a la vez.

**La fuente de verdad del color activo es un solo archivo**, en el host:
`infra/compose-prod/nginx-upstreams/active.conf`, montado en nginx como
`/etc/nginx/upstreams/active.conf`. No hay estado en ningun otro lado; para
saber que esta sirviendo, se lee ese archivo.

Secuencia de cada deploy (`bluegreen-lib.sh` + `deploy-<servicio>.sh`):

1. Pull de la imagen nueva. Si falla, no se toco nada y produccion sigue igual.
2. Se levanta el **color contrario** al activo, al lado del que esta sirviendo.
3. Se espera a que responda `/health` (hasta 150 s, sondeando cada 3 s).
   Si no responde: se borra el contenedor nuevo, el viejo sigue atendiendo y el
   deploy falla. **Esa es la reversion, y es automatica.**
4. Se reescribe `active.conf` con el color nuevo y se hace `nginx -s reload`
   (recarga en caliente, no corta conexiones). Si `nginx -t` falla, se vuelve a
   escribir el color viejo y se recarga.
5. Se drena el color viejo 15 s y recien ahi se apaga.

Consecuencias que hay que tener presentes:

- **El frontend ya no le habla al backend por `back-central:3050`.** Ese nombre
  no existe. Va por `http://nginx:8088/api/v1`, un `server` interno de nginx que
  no se publica al host y que apunta al upstream activo. Asi el frontend no
  necesita saber el color.
- Por eso **nginx no espera al frontend** en su entrypoint (seria un ciclo:
  front -> nginx -> back). Solo espera al backend activo.
- Los servicios `-green` estan en el perfil `green` de compose. Un
  `docker compose up -d` pelado levanta solo los azules; el script pide el verde
  por nombre con `--profile green`.
- Los nombres de contenedor cambiaron: `central_reserve_prod_blue|green` y
  `frontend_prod_blue|green`. Para logs, mirar cual esta activo primero.
- **Un deploy de nginx si corta** (se recrea el contenedor). Blue-green cubre
  back y front, no nginx ni website.
- En un `t4g.small` con ~770 MB libres y sin swap, durante el switch conviven
  dos copias del mismo servicio (front ~125 MB, back ~45 MB). Cabe, pero **no
  desplegar backend y frontend a la vez a proposito** si el servidor esta justo
  de RAM.

```bash
# Que color esta sirviendo
cat ~/probability/infra/compose-prod/nginx-upstreams/active.conf

# Switch manual (rollback rapido al color anterior, si sigue vivo)
cd ~/probability/infra/compose-prod
printf 'upstream probability_backend {\n    server back-central-blue:3050;\n}\n\nupstream probability_frontend {\n    server front-central-blue:3000;\n}\n' > nginx-upstreams/active.conf
docker exec nginx_prod nginx -t && docker exec nginx_prod nginx -s reload
```

### La pestana vieja y "Server Action not found"

Blue-green quita la caida del servidor, **no** arregla la pestana que un usuario
dejo abierta: cuando cambia el build de Next, los IDs de las Server Actions del
bundle viejo dejan de existir y el navegador muestra
`Server Action "<hash>" was not found on the server`.

Eso lo resuelve `VersionWatcher` (`front/central/src/shared/ui/version-watcher.tsx`):
compara `NEXT_PUBLIC_APP_VERSION` (horneado en el bundle) contra
`/api/app-version` (lo que responde el servidor que esta sirviendo). Si difieren
muestra un aviso con boton **Actualizar**, y si el usuario deja la pestana en
segundo plano la recarga sola a los 30 s, para no interrumpirlo a mitad de un
formulario.

La version la inyecta el workflow del frontend con
`--build-arg NEXT_PUBLIC_APP_VERSION=$VERSION_FULL`. En local vale `dev` y el
watcher se apaga solo.

Version tagging: `YYYY.DDD.N.XXXXXXX`. Script: `.github/scripts/generate-version.sh`

## Panic/Restart

Frontend y Nginx verifican su dependencia al iniciar; si falla hacen `exit 1` y docker reinicia (`restart: always`).
El frontend espera al backend (a traves de nginx) y nginx espera **solo** al backend del color activo:
esperar al frontend seria un ciclo. NO usar `depends_on` en compose.
Scripts: `front/central/docker/startup.sh`, `infra/nginx/entrypoint.sh`

## Rollback Manual

```bash
aws ssm start-session --profile probability --region us-east-1 --target i-0f3284d2a87127e57
sudo su - ubuntu
cd ~/probability/infra/compose-prod
docker images | grep probability-backend
docker tag <ECR_URL>/probability-backend:<VERSION_ANTERIOR> <ECR_URL>/probability-backend:latest
# Levantar el color que NO esta activo y recien ahi mover el upstream
docker compose --profile green up -d --no-deps back-central-green
docker exec nginx_prod nginx -t && docker exec nginx_prod nginx -s reload
```

## Troubleshooting

- **Nginx 502:** `docker restart nginx_prod` (cachea IPs de upstreams)
- **Puerto ocupado:** `sudo fuser -k <PORT>/tcp`
- **Container stuck:** `docker rm -f <name>`
- **Site down:** containers corriendo? health checks? iptables FORWARD=ACCEPT (ver CLAUDE.md) ? DNS resuelve?
- **Frontend/Nginx en loop:** backend no disponible. `docker logs central_reserve_prod_blue` (o `_green`, el que este activo) + `curl http://localhost:3050/health`
