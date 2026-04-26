# Deploy y CI/CD

## Workflows

Push a `main` dispara CI/CD en `.github/workflows/`. Build ARM64 -> ECR (4 tags) -> SSH EC2 -> deploy con 3 retries.
No hacer SSH para verificar deploys: confiar en GitHub Actions.

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
ssh -i ".../probability.pem" ubuntu@ec2-3-224-189-33.compute-1.amazonaws.com
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
