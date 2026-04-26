# Infraestructura y Operaciones

## AWS CLI

Siempre `--profile probability --region us-east-1`.

## Produccion (SSH)

```bash
# Conectar
ssh -i "/home/cam/Desktop/probability/probability.pem" ubuntu@ec2-3-224-189-33.compute-1.amazonaws.com

# Logs back
ssh -i ".../probability.pem" ubuntu@ec2-... "cd /home/ubuntu/probability/infra/compose-prod && docker compose logs --tail 50 back-central"
```

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

SIEMPRE `gh` CLI. NUNCA MCP de GitHub (problemas de autenticacion). Verificar cuenta con `gh auth status`.

```bash
gh pr create --title "T" --body "B" --base main
gh pr merge <n> --squash
gh run list --limit 5
```

Feature branch sync: `git fetch origin && git merge main --no-edit && git push origin <branch>`
Si >50 conflictos: crear branch nuevo y rescatar codigo especifico.
