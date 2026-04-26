# CU-05: Cambiar estado de anuncios

## Endpoint
`PATCH /api/v1/announcements/:id/status`

## Precondiciones
- CU-02 ejecutado (ANNOUNCEMENT_ID_1 existe con status "active")
- Token JWT de super admin

## Caso 5.1: Desactivar anuncio (active -> inactive)

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status
Content-Type: application/json

{ "status": "inactive" }
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] message = "status updated"

**Verificacion MCP:**
```sql
SELECT status FROM announcements WHERE id = {ANNOUNCEMENT_ID_1} AND deleted_at IS NULL;
```
- [ ] status = 'inactive'

## Caso 5.2: Reactivar anuncio (inactive -> active)

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status

{ "status": "active" }
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true

**Verificacion MCP:**
```sql
SELECT status FROM announcements WHERE id = {ANNOUNCEMENT_ID_1} AND deleted_at IS NULL;
```
- [ ] status = 'active'

## Caso 5.3: Cambiar a draft

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status

{ "status": "draft" }
```

**Verificaciones:**
- [ ] Status code = 200

## Caso 5.4: Cambiar a scheduled

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status

{ "status": "scheduled" }
```

**Verificaciones:**
- [ ] Status code = 200

## Caso 5.5: Error - status invalido

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status

{ "status": "archived" }
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "invalid status"

## Caso 5.6: Error - anuncio inexistente

**Request:**
```
PATCH /api/v1/announcements/999999/status

{ "status": "active" }
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false

## Caso 5.7: Error - sin campo status

**Request:**
```
PATCH /api/v1/announcements/{ANNOUNCEMENT_ID_1}/status

{}
```

**Verificaciones:**
- [ ] Status code = 400 (binding error, status es required)
