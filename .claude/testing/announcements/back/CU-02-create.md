# CU-02: Crear anuncios

## Endpoint
`POST /api/v1/announcements`

## Precondiciones
- Backend corriendo
- Token JWT de super admin (business_id = 0 en JWT)
- Categorias seeded (CU-01 verificado)

## Caso 2.1: Crear anuncio global tipo modal_text (happy path)

**Request:**
```
POST /api/v1/announcements
Authorization: Bearer {super_admin_token}
Content-Type: application/json

{
  "category_id": 3,
  "title": "E2E Test - Bienvenida",
  "message": "Este es un anuncio de prueba E2E",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 10,
  "is_global": true,
  "links": [],
  "target_ids": []
}
```

**Expected Response (201):**
```json
{
  "success": true,
  "data": {
    "id": "<any uint>",
    "business_id": null,
    "category_id": 3,
    "title": "E2E Test - Bienvenida",
    "message": "Este es un anuncio de prueba E2E",
    "display_type": "modal_text",
    "frequency_type": "once",
    "priority": 10,
    "is_global": true,
    "status": "active",
    "starts_at": null,
    "ends_at": null,
    "images": [],
    "links": [],
    "targets": []
  }
}
```

**Verificaciones:**
- [ ] Status code = 201
- [ ] success = true
- [ ] data.id es un entero > 0
- [ ] data.status = "active" (sin starts_at, se activa inmediatamente)
- [ ] data.is_global = true
- [ ] data.category_id = 3
- [ ] Guardar data.id como `ANNOUNCEMENT_ID_1` para tests siguientes

**Verificacion MCP:**
```sql
SELECT id, title, status, is_global, display_type, frequency_type FROM announcements WHERE title = 'E2E Test - Bienvenida' AND deleted_at IS NULL;
```
- [ ] 1 registro con status = 'active', is_global = true

## Caso 2.2: Crear anuncio con links

**Request:**
```json
{
  "category_id": 1,
  "title": "E2E Test - Promo con links",
  "message": "Visita nuestros enlaces",
  "display_type": "modal_image",
  "frequency_type": "daily",
  "priority": 5,
  "is_global": true,
  "links": [
    { "label": "Sitio web", "url": "https://example.com", "sort_order": 0 },
    { "label": "Documentacion", "url": "https://docs.example.com", "sort_order": 1 }
  ],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 201
- [ ] data.links tiene 2 elementos
- [ ] data.links[0].label = "Sitio web"
- [ ] data.links[1].label = "Documentacion"
- [ ] Guardar data.id como `ANNOUNCEMENT_ID_2`

**Verificacion MCP:**
```sql
SELECT id, label, url, sort_order FROM announcement_links WHERE announcement_id = {ANNOUNCEMENT_ID_2} ORDER BY sort_order;
```
- [ ] 2 links creados con sort_order 0 y 1

## Caso 2.3: Crear anuncio programado (scheduled)

**Request:**
```json
{
  "category_id": 5,
  "title": "E2E Test - Programado futuro",
  "message": "Este anuncio se activara en el futuro",
  "display_type": "ticker",
  "frequency_type": "always",
  "priority": 1,
  "is_global": true,
  "starts_at": "2027-12-01T00:00:00Z",
  "ends_at": "2027-12-31T23:59:59Z",
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 201
- [ ] data.status = "scheduled" (starts_at en el futuro)
- [ ] data.starts_at contiene "2027-12-01"
- [ ] data.ends_at contiene "2027-12-31"
- [ ] Guardar data.id como `ANNOUNCEMENT_ID_3`

## Caso 2.4: Crear anuncio NO global con targets

**Request:**
```json
{
  "category_id": 2,
  "title": "E2E Test - Segmentado",
  "message": "Solo para negocios especificos",
  "display_type": "modal_text",
  "frequency_type": "requires_acceptance",
  "priority": 20,
  "is_global": false,
  "links": [],
  "target_ids": [1, 2]
}
```

**Verificaciones:**
- [ ] Status code = 201
- [ ] data.is_global = false
- [ ] data.targets tiene 2 elementos
- [ ] data.targets contiene business_id 1 y 2
- [ ] Guardar data.id como `ANNOUNCEMENT_ID_4`

**Verificacion MCP:**
```sql
SELECT business_id FROM announcement_targets WHERE announcement_id = {ANNOUNCEMENT_ID_4};
```
- [ ] 2 registros con business_id 1 y 2

## Caso 2.5: Error - display_type invalido

**Request:**
```json
{
  "category_id": 1,
  "title": "Test invalido",
  "message": "Esto debe fallar",
  "display_type": "popup",
  "frequency_type": "once",
  "priority": 0,
  "is_global": true,
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "invalid display type"

## Caso 2.6: Error - no global sin targets

**Request:**
```json
{
  "category_id": 1,
  "title": "Test sin targets",
  "message": "No global pero sin targets",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 0,
  "is_global": false,
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "target"

## Caso 2.7: Error - fecha inicio despues de fecha fin

**Request:**
```json
{
  "category_id": 1,
  "title": "Test fechas invalidas",
  "message": "Rango invalido",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 0,
  "is_global": true,
  "starts_at": "2027-12-31T00:00:00Z",
  "ends_at": "2027-01-01T00:00:00Z",
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "starts_at must be before ends_at"

## Caso 2.8: Error - campos requeridos faltantes (sin title)

**Request:**
```json
{
  "category_id": 1,
  "message": "Sin titulo",
  "display_type": "modal_text",
  "frequency_type": "once",
  "is_global": true,
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
