# Probability

## Sin Comentarios

NUNCA comentarios en Go/TypeScript. Al modificar un archivo, eliminar TODOS los comentarios existentes en el.
Excepciones: `//go:generate`, `//go:build`, `//nolint`, `'use server'`, `'use client'`

## UTF-8 / highlight.js Bug

Archivos 500+ lineas: CERO non-ASCII (acentos, box-drawing, emojis). Archivos cortos: preferir ASCII.
- Strings con acentos: escapes `\u00XX` o `golang.org/x/text/transform`
- Limpieza rapida: `sed -i 's/a/a/g; s/e/e/g; s/i/i/g; s/o/o/g; s/u/u/g; s/n/n/g'`
- NUNCA sugerir actualizar Claude Code como solucion a este bug.

## Stack

**Back:** Go 1.23 + Gin + GORM + RabbitMQ + Redis + JWT | `/back/central` API :3050
**Front:** Next.js 16 + React 19 + TailwindCSS 4 | `/front/central` :3000 | `/front/website` Astro 5
**Infra:** PostgreSQL 15 :5433 | Redis :6379 | RabbitMQ :5672 | MinIO :9000 | Docker Compose

Monorepo multi-tenant: ordenes, productos, pagos, envios desde Shopify, Amazon, MercadoLibre, WhatsApp.

## Produccion - iptables CRITICO

Si el sitio deja de funcionar desde Internet (AWS/SGs siempre estan bien, el problema es siempre iptables):

```bash
sudo iptables -P FORWARD ACCEPT
sudo iptables -I FORWARD 1 -s 10.89.0.0/24 -j ACCEPT
sudo iptables -I FORWARD 2 -d 10.89.0.0/24 -j ACCEPT
```

Reiniciar contenedores UNA sola vez. NUNCA `iptables -F` ni `iptables -t nat -F`.
