# Probability

## Documentacion

- `ROADMAP.md` - orden de prioridades (P0/P1/P2). Solo ordena, no duplica
  contenido: el detalle vive en la alerta o documento que cada linea referencia.
  Consultar antes de preguntar "que sigue".
- `.claude/alerts/` - pendientes urgentes. Revisar al iniciar sesion.
- **Tickets** - todo bug, correccion o desarrollo se registra en un ticket y se
  cierra comentandolo con el diagnostico y la correccion. Buscar el ticket antes
  de empezar; si no existe, crearlo. Reglas: `.claude/rules/tickets.md`
- `.claude/bitacora/` - historico de soportes, incidentes y correcciones, un
  archivo por caso. Buscar aca antes de investigar un problema parecido.
  Reglas: `.claude/rules/bitacora.md`
- `AGENTS.md` - punto de entrada para agentes de IA que no sean Claude Code.
- `back/migration/MIGRACIONES.md` - flujo de migraciones e historico de corridas.

## Sin Comentarios

NUNCA comentarios en Go/TypeScript. Al modificar un archivo, eliminar TODOS los comentarios existentes en el.
Excepciones: `//go:generate`, `//go:build`, `//nolint`, `'use server'`, `'use client'`

## UTF-8 / highlight.js Bug

Archivos 500+ lineas: CERO non-ASCII (acentos, box-drawing, emojis). Archivos cortos: preferir ASCII.
- Strings con acentos: escapes `\u00XX` o `golang.org/x/text/transform`
- Limpieza rapida: `sed -i 's/a/a/g; s/e/e/g; s/i/i/g; s/o/o/g; s/u/u/g; s/n/n/g'`
- **Ese sed NUNCA se aplica a texto que ve el usuario.** Quitarle la tilde a un
  label o a un mensaje es una falta de ortografia, no una solucion al bug: ahi
  la tilde se escribe con `\u00XX` (`{'Facturaci\u00f3n'}` en JSX). Solo se usa
  sobre comentarios, nombres y texto interno.
- Ortografia del espanol en el front: skill `ortografia-front`
  (`.claude/skills/ortografia-front/`, incluye `revisar.py` para detectar y corregir).
- NUNCA sugerir actualizar Claude Code como solucion a este bug.

## Stack

**Back:** Go 1.23 + Gin + GORM + RabbitMQ + Redis + JWT | `/back/central` API :3050
**Front:** Next.js 16 + React 19 + TailwindCSS 4 | `/front/central` :3000 | `/front/website` Astro 5
**Infra:** PostgreSQL 17 (local :5434, tunel prod :5433) | Redis :6379 | RabbitMQ :5672 | Docker Compose (S3: AWS, sin MinIO local)

Monorepo multi-tenant: ordenes, productos, pagos, envios desde Shopify, Amazon, MercadoLibre, WhatsApp.

## WooCommerce de Pruebas (local)

WordPress + WooCommerce en Docker para probar la integracion (conexion, sync, webhooks).
Carpeta `/wordpress` (volumenes nombrados, no toca prod). Levantar: `cd wordpress && ./setup.sh`.
Tienda en `http://localhost:8088`, wp-admin admin/admin. Detalles: `wordpress/README.md`.

## Git - NUNCA hacer push sin permiso explicito

**PROHIBIDO hacer push automaticamente.** Siempre esperar autorizacion explicita del usuario.
- NO hacer push despues de commits
- NO hacer push como parte de flujos de trabajo
- Informar al usuario que hay cambios listos, pero esperar instruccion para hacer push

## Entorno local primero

Todo cambio se prueba contra la copia local de la base (business Demo) ANTES de
tocar produccion. PostgreSQL local en `127.0.0.1:5434` (el 5433 es el tunel a
produccion). Clonar: `./scripts/local-db-clone.sh 26`. Cambiar de base:
`./scripts/dev-db-switch.sh local|prod|status`.
Reglas y flujo: `.claude/rules/entorno-local.md`.

## Produccion - acceso solo por AWS CLI (SSM)

No hay SSH ni `.pem`, y la base no es alcanzable desde internet. Lo unico abierto
es 80/443. Para consultar la BD hay que levantar el tunel primero:
`./scripts/aws-tunnel.sh ensure` (el MCP de postgres apunta a `127.0.0.1:5433`).
Detalles y prohibiciones: `.claude/rules/infra-ops.md`.

## Produccion - iptables CRITICO

Si el sitio deja de funcionar desde Internet (AWS/SGs siempre estan bien, el problema es siempre iptables):

```bash
sudo iptables -P FORWARD ACCEPT
sudo iptables -I FORWARD 1 -s 10.89.0.0/24 -j ACCEPT
sudo iptables -I FORWARD 2 -d 10.89.0.0/24 -j ACCEPT
```

Reiniciar contenedores UNA sola vez. NUNCA `iptables -F` ni `iptables -t nat -F`.
