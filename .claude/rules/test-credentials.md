# Credenciales de Testing (para IA)

Las credenciales reales viven en `.env.ai` (gitignored). Plantilla: `.env.ai.example`.

## Variables disponibles

### Super Admin
- `${AI_SUPER_ADMIN_EMAIL}`
- `${AI_SUPER_ADMIN_PASSWORD}`

### Usuario Business (no super admin)
- `${AI_DEMO_EMAIL}`
- `${AI_DEMO_PASSWORD}`

## Uso

Cargar el archivo antes de ejecutar tests E2E:
```bash
set -a && source .env.ai && set +a
```

Luego en HTTP requests usar `$AI_SUPER_ADMIN_EMAIL`, etc.
