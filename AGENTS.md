# AGENTS.md

Instrucciones para cualquier agente de IA que trabaje en este repositorio
(Claude Code, Codex, Cursor, Copilot, Gemini u otro).

Este archivo es el punto de entrada. Las reglas completas viven en `CLAUDE.md`
y en `.claude/rules/`; aca esta lo minimo que hay que saber antes de tocar nada,
y donde buscar el resto.

## Antes de empezar

1. Leer `CLAUDE.md` (raiz): stack, reglas duras del proyecto.
2. Revisar `.claude/alerts/`: pendientes urgentes. Si hay una alerta del modulo
   que vas a tocar, leerla completa.
3. Buscar en `.claude/bitacora/`: historico de soportes y diagnosticos. Si el
   problema se parece a algo ya investigado, ahi esta el contexto y las
   hipotesis que ya se descartaron.

## Reglas duras

- **Sin comentarios en Go/TypeScript.** Al modificar un archivo, eliminar los
  que existan. Excepciones: `//go:generate`, `//go:build`, `//nolint`,
  `'use server'`, `'use client'`.
- **Nunca hacer push sin autorizacion explicita del usuario.** Tampoco como
  parte de un flujo automatico.
- **Nunca iniciar, reiniciar o detener el backend sin permiso.** Usar
  `./scripts/dev-services.sh`, nunca `go run cmd/main.go &` ni `nohup`.
- **El `.env` local apunta al RDS de produccion.** Cualquier migracion o
  backfill corrido en local escribe en produccion.
- **Migraciones**: `Migrate()` esta en cero a proposito. Se agrega la llamada
  que se necesita, se corre, se devuelve a cero y se anota en
  `back/migration/MIGRACIONES.md`. Nunca correr la cadena historica completa.
- **Aislamiento multi-tenant**: el `business_id` de un usuario normal sale
  SIEMPRE del token, jamas del body o query. Ver
  `.claude/rules/multi-tenant-security.md`.
- **UTF-8**: archivos de 500+ lineas sin caracteres no-ASCII (bug de
  highlight.js). En archivos cortos, preferir ASCII igual.

## Mapa de la documentacion

| Ruta | Que hay |
|------|---------|
| `CLAUDE.md` | Reglas del proyecto, stack, puertos, produccion |
| `.claude/rules/architecture.md` | Arquitectura hexagonal (Go y Next.js) |
| `.claude/rules/backend-conventions.md` | Convenciones Go, migraciones, logging, super admin |
| `.claude/rules/multi-tenant-security.md` | Aislamiento por negocio (obligatorio) |
| `.claude/rules/colas-errores-permanentes.md` | Consumidores RabbitMQ: error permanente vs transitorio |
| `.claude/rules/infra-ops.md` | AWS, SSH a produccion, servicios de desarrollo, `gh` |
| `.claude/rules/deploy.md` | CI/CD, rollback, troubleshooting |
| `.claude/rules/testing.md` | Casos de uso E2E |
| `.claude/rules/alerts.md` | Como usar y crear alertas |
| `.claude/rules/bitacora.md` | Como usar y escribir la bitacora |
| `.claude/alerts/` | Pendientes urgentes, se borran al resolverse |
| `.claude/bitacora/` | Historico de soportes y diagnosticos, un archivo por caso |
| `back/migration/MIGRACIONES.md` | Flujo de migraciones e historico de corridas |

## Estructura

- `back/central` - API Go (Gin + GORM), puerto 3050
- `back/migration` - migraciones, unica fuente de cambios de esquema
- `front/central` - Next.js, puerto 3000
- `front/website` - Astro
- `mobile/mobile_central` - Flutter
- `infra/` - Docker Compose y nginx

## Al terminar un trabajo relevante

- Si queda trabajo critico inconcluso: crear alerta en `.claude/alerts/`.
- Si se investigo un problema, se corrigio data en produccion, o se descubrio
  un comportamiento no documentado de un proveedor: crear entrada en
  `.claude/bitacora/`.
