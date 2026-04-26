# CU-08: Forzar re-visualizacion

## Endpoint
`POST /api/v1/announcements/:id/force-redisplay`

## Precondiciones
- CU-06 ejecutado (ANNOUNCEMENT_ID_1 tiene views registradas)
- Token JWT de super admin

## Caso 8.1: Forzar redisplay (happy path)

**Pre-verificacion MCP:**
```sql
SELECT COUNT(*) FROM announcement_views WHERE announcement_id = {ANNOUNCEMENT_ID_1};
```
- [ ] count > 0 (hay views registradas de CU-06)

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_1}/force-redisplay
Authorization: Bearer {super_admin_token}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] message = "redisplay forced"

**Verificacion MCP:**
```sql
SELECT force_redisplay FROM announcements WHERE id = {ANNOUNCEMENT_ID_1} AND deleted_at IS NULL;
```
- [ ] force_redisplay = true

```sql
SELECT COUNT(*) FROM announcement_views WHERE announcement_id = {ANNOUNCEMENT_ID_1};
```
- [ ] count = 0 (las views fueron eliminadas)

## Caso 8.2: Verificar que el anuncio vuelve a aparecer en /active

**Despues de force-redisplay, ANNOUNCEMENT_ID_1 debe volver a aparecer para el usuario:**

**Request:**
```
GET /api/v1/announcements/active
Authorization: Bearer {business_user_token}
```

**Verificaciones:**
- [ ] ANNOUNCEMENT_ID_1 aparece en la lista (porque sus views fueron eliminadas)

## Caso 8.3: Error - anuncio inexistente

**Request:**
```
POST /api/v1/announcements/999999/force-redisplay
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
