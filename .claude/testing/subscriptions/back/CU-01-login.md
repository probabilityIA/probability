# CU-01: Login super admin y demo

## Objetivo
Obtener JWT de ambos roles para ejercitar el resto de los casos.

## Precondiciones
- Backend corriendo en `http://localhost:3050`.
- Variables cargadas: `set -a && source .env.ai && set +a`.

## Endpoint
```
POST /api/v1/auth/login
Content-Type: application/json
```

## Body (super admin)
```json
{
  "email": "{AI_SUPER_ADMIN_EMAIL}",
  "password": "{AI_SUPER_ADMIN_PASSWORD}"
}
```

## Body (demo / business normal)
```json
{
  "email": "{AI_DEMO_EMAIL}",
  "password": "{AI_DEMO_PASSWORD}"
}
```

## Verificaciones
- [ ] Status 200, `success=true`
- [ ] `data.token` presente -> guardar como SUPER_TOKEN / DEMO_TOKEN
- [ ] Token super admin: payload JWT trae `business_id: 0`
- [ ] Token demo: payload JWT trae `business_id: 26` (Demo)
