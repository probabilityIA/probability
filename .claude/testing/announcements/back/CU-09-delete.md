# CU-09: Eliminar anuncios

## Endpoint
`DELETE /api/v1/announcements/:id`

## Precondiciones
- CU-02 ejecutado (anuncios de prueba existen)
- Token JWT de super admin

## Caso 9.1: Eliminar anuncio (soft delete)

**Request:**
```
DELETE /api/v1/announcements/{ANNOUNCEMENT_ID_3}
Authorization: Bearer {super_admin_token}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] message = "announcement deleted"

**Verificacion MCP:**
```sql
SELECT id, deleted_at FROM announcements WHERE id = {ANNOUNCEMENT_ID_3};
```
- [ ] deleted_at IS NOT NULL (soft delete)

## Caso 9.2: Verificar que el anuncio eliminado no aparece en listado

**Request:**
```
GET /api/v1/announcements
```

**Verificaciones:**
- [ ] data NO contiene un item con id = ANNOUNCEMENT_ID_3

## Caso 9.3: Verificar que el anuncio eliminado no se puede obtener

**Request:**
```
GET /api/v1/announcements/{ANNOUNCEMENT_ID_3}
```

**Verificaciones:**
- [ ] Status code = 404
- [ ] success = false

## Caso 9.4: Error - eliminar anuncio inexistente

**Request:**
```
DELETE /api/v1/announcements/999999
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false

## Caso 9.5: Eliminar anuncio con links y targets (cascade check)

**Request:**
```
DELETE /api/v1/announcements/{ANNOUNCEMENT_ID_4}
```

**Verificaciones:**
- [ ] Status code = 200

**Verificacion MCP:**
```sql
SELECT COUNT(*) FROM announcement_targets WHERE announcement_id = {ANNOUNCEMENT_ID_4};
```
- [ ] Los targets siguen en la BD pero el anuncio tiene deleted_at (GORM soft delete)

## Caso 9.6: Cleanup - eliminar anuncios de prueba restantes

Eliminar ANNOUNCEMENT_ID_1 y ANNOUNCEMENT_ID_2:
```
DELETE /api/v1/announcements/{ANNOUNCEMENT_ID_1}
DELETE /api/v1/announcements/{ANNOUNCEMENT_ID_2}
```

**Verificaciones:**
- [ ] Ambos retornan status 200

**Verificacion MCP final:**
```sql
SELECT COUNT(*) FROM announcements WHERE title LIKE 'E2E Test%' AND deleted_at IS NULL;
```
- [ ] count = 0 (todos los de prueba eliminados)
