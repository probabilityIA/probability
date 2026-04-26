# CU-01: Setup - Activar modo test y obtener tokens

## Precondiciones
- Backend corriendo en localhost:3050
- Mock server corriendo en localhost:9091
- PostgreSQL accesible

## Caso 1.1: Obtener token super admin

**Request:**
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "${AI_SUPER_ADMIN_EMAIL}",
  "password": "${AI_SUPER_ADMIN_PASSWORD}"
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "data": {
    "token": "{SUPER_ADMIN_TOKEN}"
  }
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] token no vacio -> guardar como SUPER_ADMIN_TOKEN

## Caso 1.2: Obtener token usuario Demo

**Request:**
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "${AI_DEMO_EMAIL}",
  "password": "${AI_DEMO_PASSWORD}"
}
```

**Expected Response (200):**
```json
{
  "success": true,
  "data": {
    "token": "{DEMO_TOKEN}"
  }
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] token no vacio -> guardar como DEMO_TOKEN

## Caso 1.3: Activar modo test en integracion Demo

**Verificacion MCP (antes):**
```sql
SELECT id, name, is_testing, business_id FROM integrations WHERE id = 44;
```
- [ ] is_testing = false (estado actual)

**Accion:** Activar test mode via SQL
```sql
UPDATE integrations SET is_testing = true WHERE id = 44;
```

**Verificacion MCP (despues):**
```sql
SELECT id, name, is_testing, business_id FROM integrations WHERE id = 44;
```
- [ ] is_testing = true

## Caso 1.4: Verificar mock server operativo

**Request:**
```
GET http://localhost:9091/health
```

**Expected Response (200):**
```json
{
  "status": "ok",
  "service": "envioclick-mock"
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] status = "ok"
