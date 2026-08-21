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
| Backend  | `back/central/**`, `back/migration/**` | 3050 |
| Frontend | `front/central/**` | 8080 |
| Website  | `front/website/**` | 8081 |
| Nginx    | `infra/nginx/**` | 80/443 |

Version tagging: `YYYY.DDD.N.XXXXXXX`. Script: `.github/scripts/generate-version.sh`

## Panic/Restart

Frontend y Nginx verifican dependencias al iniciar; si fallan hacen `exit 1` y docker reinicia (`restart: always`).
NO usar `depends_on` en compose. Scripts: `front/central/docker/startup.sh`, `infra/nginx/entrypoint.sh`

## Rollback Manual

```bash
aws ssm start-session --profile probability --region us-east-1 --target i-0f3284d2a87127e57
sudo su - ubuntu
cd ~/probability/infra/compose-prod
docker images | grep probability-backend
docker tag <ECR_URL>/probability-backend:<VERSION_ANTERIOR> <ECR_URL>/probability-backend:latest
docker compose up -d back-central
```

## Troubleshooting

- **Nginx 502:** `docker restart nginx_prod` (cachea IPs de upstreams)
- **Puerto ocupado:** `sudo fuser -k <PORT>/tcp`
- **Container stuck:** `docker rm -f <name>`
- **Site down:** containers corriendo? health checks? iptables FORWARD=ACCEPT (ver CLAUDE.md) ? DNS resuelve?
- **Frontend/Nginx en loop:** backend no disponible. `docker logs back-central` + `curl http://localhost:3050/health`
