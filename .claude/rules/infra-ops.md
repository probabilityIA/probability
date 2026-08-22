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
(`back-central`, `back-testing`) llevan el DNS de la VPC en el compose para
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
  --parameters 'commands=["docker ps","docker logs --tail 50 central_reserve_prod"]' \
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
Si `docker compose up -d` falla por monitoring: `docker compose up -d rabbitmq redis back-central back-testing front-central front-website nginx front-testing`

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
