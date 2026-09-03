# Infraestructura y Operaciones

## AWS CLI

Siempre `--profile probability --region us-east-1`.

## Todo el acceso es por AWS CLI. No hay puertos publicos ni .pem

Desde 2026-08-21 la unica forma de llegar a produccion es el AWS CLI sobre SSM.
Lo unico abierto a internet son **80 y 443** (nginx). Cerrados: 22, 3002, 3070 y
5432. El RDS quedo con `PubliclyAccessible=false` y su security group solo acepta
al SG del EC2 (`sg-03816f3607edc744b`).

**Prohibido reabrir un puerto a `0.0.0.0/0` para "probar rapido"**, y prohibido
volver a poner `PubliclyAccessible=true`. Si algo no conecta, el camino es el
tunel, no `authorize-security-group-ingress`. Cualquier excepcion se acuerda con
el usuario y se revierte el mismo dia.

```bash
./scripts/aws-tunnel.sh ensure   # tunel al RDS en 127.0.0.1:5433 (idempotente)
./scripts/aws-tunnel.sh status   # verifica que este arriba y prueba la conexion
./scripts/aws-tunnel.sh shell    # shell interactiva en el EC2
./scripts/aws-tunnel.sh stop
```

**El MCP `postgres-probability` apunta a `127.0.0.1:5433`, no al RDS.** Si una
consulta falla con "connection refused", el tunel esta abajo: correr
`./scripts/aws-tunnel.sh ensure` y reintentar. El tunel muere al cerrar la
terminal, asi que al iniciar sesion hay que levantarlo.

Requisito local: `session-manager-plugin` instalado (el script lo avisa si falta).

Usuarios IAM con acceso al tunel: `terraform-github-admin`, `santiago.camacho`
(politica `ProbabilityRdsTunnelSSM`: solo abrir el port-forwarding, sin permisos
sobre RDS ni EC2).

### Por que el contenedor necesita `dns: 169.254.169.253`

El resolver de Docker reenvia a 8.8.8.8, que devolvia la IP **publica** del RDS y
mandaba el trafico por el internet gateway. Los servicios que tocan la base
(`back-central-blue`/`back-central-green`, `back-testing`) llevan el DNS de la VPC en el compose para
resolver `172.31.96.15` e ir por red privada. **Cualquier servicio nuevo que
consulte la base necesita ese `dns:`**, o no conectara.

## Produccion (SSM, no SSH)

Desde 2026-08-21 el acceso a produccion va por AWS Systems Manager. No se usa
`probability.pem` ni el puerto 22. Requiere `session-manager-plugin` instalado.

Instancia: `i-0f3284d2a87127e57` (EC2 prod, us-east-1).

```bash
# Ejecutar comandos (equivalente a "ssh ... 'comando'")
aws ssm send-command --profile probability --region us-east-1 \
  --instance-ids i-0f3284d2a87127e57 --document-name AWS-RunShellScript \
  --parameters 'commands=["docker ps","docker logs --tail 50 central_reserve_prod_blue"]' \
  --query 'Command.CommandId' --output text
# ...luego leer el resultado:
aws ssm get-command-invocation --profile probability --region us-east-1 \
  --command-id <CMD_ID> --instance-id i-0f3284d2a87127e57 \
  --query '{status:Status,out:StandardOutputContent,err:StandardErrorContent}'

# Shell interactiva
aws ssm start-session --profile probability --region us-east-1 --target i-0f3284d2a87127e57
```

El comando de `send-command` corre como **root**. Para replicar el entorno de
siempre, envolver con `runuser -l ubuntu -c "..."`.

### Tunel a RDS (la unica forma de consultar la BD de produccion)

El puerto 5432 del RDS ya no esta abierto a internet. Para consultar:

```bash
aws ssm start-session --profile probability --region us-east-1 \
  --target i-0f3284d2a87127e57 \
  --document-name AWS-StartPortForwardingSessionToRemoteHost \
  --parameters '{"host":["database-1.capmmoe4cw2e.us-east-1.rds.amazonaws.com"],"portNumber":["5432"],"localPortNumber":["5433"]}'
```

Dejarlo corriendo en segundo plano y conectar contra `127.0.0.1:5433` con las
credenciales de `back/central/.env`. El MCP de postgres y cualquier cliente SQL
apuntan ahi.

Dir servidor: `/home/ubuntu/probability/infra/compose-prod/`
Solo docker/docker compose (podman desinstalado).
Si `docker compose up -d` falla por monitoring: `docker compose up -d rabbitmq redis back-central-blue back-testing front-central-blue front-website nginx front-testing`

**Blue-green:** `back-central` y `front-central` ya no existen como servicio; hay
un `-blue` y un `-green` y el activo lo dice
`infra/compose-prod/nginx-upstreams/active.conf`. Ver `.claude/rules/deploy.md`.

## Servicios de Desarrollo

Script: `./scripts/dev-services.sh`

```bash
./scripts/dev-services.sh status
./scripts/dev-services.sh start all          # infra + backend + frontend
./scripts/dev-services.sh restart backend    # detiene + limpia + inicia
./scripts/dev-services.sh logs backend 100
./scripts/dev-services.sh kill-zombies
./scripts/dev-services.sh ports
```

Puertos: infra 5433/6379/5672/9000 | backend :3050 | frontend :3000
NUNCA `go run cmd/main.go &` ni `nohup`. Siempre el script.

## Usuarios IAM

Cuenta AWS: `476702565908`. No usar la raiz para operar.

| Usuario | Para que |
|---|---|
| `terraform-github-admin` | infraestructura y CI |
| `santiago.camacho` | operacion: S3, tunel al RDS, shell del EC2 |
| `backend-s3-uploader` | credenciales de la app para subir a S3 |
| `docker-ecr-dev` | push de imagenes a ECR |
| `bedrik-IA`, `woo-store-power` | integraciones puntuales |

Politicas propias (no gestionadas por AWS):

- `ProbabilityRdsTunnelSSM` - solo abrir el port-forwarding al RDS. Sin permisos
  sobre RDS ni EC2. La tienen `terraform-github-admin` y `santiago.camacho`.
- `ProbabilityOpsSantiago` - S3 completo, lectura de RDS y EC2, `ssm:StartSession`
  y `ssm:SendCommand` **restringidos a `i-0f3284d2a87127e57`**, y gestion de sus
  propias access keys. Deliberadamente NO puede: crear o borrar EC2/RDS, tocar
  security groups, ni administrar IAM de terceros.

La guia de acceso que se le entrega a un operador nuevo (instalacion del CLI y
del session-manager-plugin, tunel, shell, compose, S3, rotacion de llaves) esta
en `~/Desktop/acceso-aws-santiago.md`. **No esta en el repo a proposito**: es un
entregable con instrucciones de credenciales.

Al dar de alta a alguien: crear el usuario, adjuntar las dos politicas de arriba,
entregar la llave por canal privado (nunca chat, correo ni git) y entregar el .md.

## ECR - retencion de imagenes

Cada repo tiene una **lifecycle policy** que expira las imagenes viejas. AWS la
evalua una vez cada ~24 h, asi que ver mas imagenes de las que dice la regla no
significa que este rota: son las acumuladas desde la ultima pasada.

| Repo | Mantiene |
|---|---|
| `probability-backend`, `probability-frontend`, `probability-website`, `probability-nginx` | 5 |
| `probability-testing-backend`, `probability-testing-frontend` | 3 |

**Produccion mantiene 5 y no 1 a proposito:** el rollback manual de
`.claude/rules/deploy.md` hace `docker tag <VERSION_ANTERIOR>` contra ECR. Con la
regla vieja de "mantener 1" la version anterior ya no estaba en el registro y el
rollback solo funcionaba si la imagen seguia cacheada en el disco del EC2 — que
el propio deploy limpia. Bajar ese numero a 1 otra vez deja el rollback sin red.

**Un repo nuevo NO hereda ninguna politica.** Hay que ponersela al crearlo o
crece sin limite: `probability-testing-backend` llego a 198 imagenes (2,69 GB,
61% del almacenamiento de ECR) simplemente porque nadie se la puso.

```bash
aws ecr put-lifecycle-policy --profile probability --region us-east-1 \
  --repository-name <repo> --lifecycle-policy-text \
  '{"rules":[{"rulePriority":1,"description":"Mantener las 5 imagenes mas recientes; el resto expira","selection":{"tagStatus":"any","countType":"imageCountMoreThan","countNumber":5},"action":{"type":"expire"}}]}'

# Que hay hoy, por repo
aws ecr describe-images --profile probability --region us-east-1 \
  --repository-name <repo> --query 'length(imageDetails)'
```

Antes de borrar imagenes a mano (`batch-delete-image`, irreversible), comprobar
que ninguna este en uso: sacar el digest de cada contenedor en produccion con
`docker inspect <id> --format '{{json .RepoDigests}}'` y verificar que no aparece
en la lista a borrar. Ojo con la diferencia entre **tags** e **imagenes**: cada
build empuja 4 tags de una misma imagen, asi que `list-images` cuenta hasta 4
veces mas que `describe-images`.

## Auditoria - CloudTrail

Trail `probability-trail`, activo desde 2026-08-21. Multi-region, con eventos
globales y validacion de integridad de archivos.

Bucket: `probability-cloudtrail-476702565908` (public access bloqueado, SSE-S3,
policy que solo deja escribir a CloudTrail y niega trafico sin TLS).
**Lifecycle: los logs se borran a los 90 dias.**

Registra **management events unicamente**: quien llamo que API de AWS, desde que
IP y cuando. Los *data events* (cada GET/PUT de objetos en S3) estan APAGADOS a
proposito porque se cobran a $0.10 por 100.000. No mandar el trail a CloudWatch
Logs: esa ingesta tambien se cobra.

Que SI queda grabado: `ssm:StartSession`, `ssm:SendCommand` (con el comando
enviado en los parametros), cambios de IAM, `ec2:AuthorizeSecurityGroupIngress`,
`rds:ModifyDBInstance`, logins a consola.

Que NO queda grabado: `docker logs`, salida de los contenedores, access logs de
nginx, queries de PostgreSQL, ni lo que alguien tipea dentro de una shell SSM
interactiva. CloudTrail audita la CUENTA, no la aplicacion.

```bash
aws cloudtrail lookup-events --profile probability --region us-east-1 \
  --lookup-attributes AttributeKey=Username,AttributeValue=santiago.camacho \
  --max-results 20 --query 'Events[].{T:EventTime,E:EventName,U:Username}' --output table
```

`lookup-events` da los ultimos 90 dias gratis incluso sin trail. Lo que aporta el
trail es la copia propia con validacion de integridad.

## Monitoreo y alertas

Todo lo activo hoy cabe en la capa gratuita.

| Que | Estado |
|---|---|
| 7 alarmas CloudWatch | activas (free tier: 10 alarmas) |
| Tema SNS `probability-alertas` | notifica a `secamc93@gmail.com` |
| RDS Performance Insights | 7 dias (el tramo gratuito) |
| RDS Enhanced Monitoring 60s | log group `RDSOSMetrics`, 30 dias |
| Cost Anomaly Detection | avisa a `probabilitysas@gmail.com` |
| EC2 detailed monitoring | **disabled** a proposito: el de 1 minuto se cobra |

Alarmas: `ec2-cpu-alta` (>85%), `ec2-creditos-cpu-bajos` (<30),
`ec2-status-check-falla`, `rds-almacenamiento-bajo` (<5 GB),
`rds-conexiones-altas` (>60), `rds-creditos-cpu-bajos` (<50),
`rds-memoria-libre-baja` (<100 MB).

**Al crear una alarma nueva, apuntarla al tema SNS.** Una alarma sin
`AlarmActions` no le avisa a nadie: es exactamente el error que estuvo vivo hasta
2026-08-21 (ver `.claude/bitacora/2026-08-21-monitoreo-sin-destinatario.md`).

Verificar que el destinatario sigue conectado:

```bash
aws sns list-subscriptions-by-topic --profile probability --region us-east-1 \
  --topic-arn arn:aws:sns:us-east-1:476702565908:probability-alertas --output text
```

Si aparece `PendingConfirmation` en vez de un ARN, **las alarmas no notifican**.
Cada correo de AWS trae un enlace de "unsubscribe" al pie: un clic accidental
deja el sistema mudo sin avisar.

Pendientes gratuitos, sin hacer: budget de costos (2 gratis), dashboard de
CloudWatch (3 gratis) y metricas de RAM/disco del EC2 via CloudWatch agent (hoy
nadie las vigila, y es lo que mas probablemente tumbe el `t4g.small`).

### Logs de contenedores en CloudWatch

Activo desde 2026-08-21 para el backend y el frontend, **solo lineas de
error**. Los 5 GB gratis de CloudWatch Logs son de INGESTA, no de almacenamiento:
borrar despues no devuelve nada, y no existe un tope duro que corte el envio. El
control es cuanto se manda.

Cadena: `docker logs -f` -> filtro grep -> archivo local -> CloudWatch agent.

- `/usr/local/bin/probability-logtail.sh <contenedor>` - resuelve el contenedor
  por prefijo de nombre (blue-green cambia el sufijo), hace `docker logs -f` y
  deja pasar solo `ERROR|WARN|FATAL|PANIC|panic:|Exception|error|warn|fatal`.
- `probability-logtail@<contenedor>.service` - unidad systemd template, con
  `Restart=always`. Reengancha sola cuando el deploy recrea el contenedor.
- Salida en `/var/log/probability/<contenedor>.log`, con logrotate (20 MB, 3
  copias) en `/etc/logrotate.d/probability`.
- CloudWatch agent (`amazon-cloudwatch-agent`) los sube a `/probability/back-central`
  y `/probability/front-central`, retencion **7 dias**.
- Permiso: `CloudWatchAgentServerPolicy` en el rol `probability-ec2-ecr-pull-role`.

**No se uso el log driver `awslogs`** a proposito: manda a AWS sin dejar copia
local y rompe `docker logs` en el servidor.

Volumen medido en produccion: el backend escribe ~10.800 lineas/hora y solo 26
pasan el filtro (0,24%). Sin filtro serian ~5,4 GB/mes; con filtro, decenas de MB.

La alarma `logs-ingesta-alta` vigila `IncomingBytes` de los dos grupos sumados
(metric math) y avisa si superan 10 MB en una hora. Existe porque el filtro NO
protege de un bucle: en el incidente de la cola de RabbitMQ el backend escupia
4,2 MB/min y eran todas lineas de ERROR, o sea justo las que el filtro deja pasar.

Para agregar otro contenedor:

```bash
sudo systemctl enable --now probability-logtail@<nombre_contenedor>.service
# y agregar el archivo al collect_list de
# /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json
```

Todos los servicios del compose de produccion tienen `logging` con
`max-size: 10m` / `max-file: 3` (el backend 50m/10). Un servicio nuevo sin ese
bloque escribe logs sin limite en el disco del EC2.

Los servicios `monitoring-api` y `monitoring-web` estan definidos en el compose
de produccion pero no corren, igual que las carpetas `grafana/`, `prometheus.yml`
y `alloy/`. En un `t4g.small` con 8 contenedores, Prometheus + Grafana se comen
la RAM que no sobra: se dejan apagados a proposito.

## GitHub

SIEMPRE `gh` CLI. NUNCA MCP de GitHub (problemas de autenticacion).

**OBLIGATORIO en cada terminal nueva antes de usar `gh`:**

```bash
eval "$(./scripts/gh-env.sh)"
```

Esto exporta `GH_TOKEN` con el PAT de `secamc93` y fija `GH_REPO=probabilityIA/probability`,
sin tocar el keyring global (donde `velocity` queda como default). Verificar con `gh auth status`
— debe decir `Active account: true` y `(GH_TOKEN)`, no `(keyring)`.

El token vive en `.mcp.json` o `.gh-token` (ambos en `.gitignore`). Si `gh auth status` muestra
`Bad credentials`, regenerar el PAT en https://github.com/settings/tokens y reemplazar en el archivo.

```bash
gh pr create --title "T" --body "B" --base main
gh pr merge <n> --squash
gh run list --limit 5
```

Feature branch sync: `git fetch origin && git merge main --no-edit && git push origin <branch>`
Si >50 conflictos: crear branch nuevo y rescatar codigo especifico.
