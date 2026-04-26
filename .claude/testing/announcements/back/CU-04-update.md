# CU-04: Actualizar anuncios

## Endpoint
`PUT /api/v1/announcements/:id`

## Precondiciones
- CU-02 ejecutado (ANNOUNCEMENT_ID_1 existe)
- Token JWT de super admin

## Caso 4.1: Actualizar titulo y mensaje

**Request:**
```
PUT /api/v1/announcements/{ANNOUNCEMENT_ID_1}
Content-Type: application/json

{
  "category_id": 3,
  "title": "E2E Test - Bienvenida (Editado)",
  "message": "Mensaje actualizado en E2E",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 15,
  "is_global": true,
  "links": [],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] data.title = "E2E Test - Bienvenida (Editado)"
- [ ] data.message = "Mensaje actualizado en E2E"
- [ ] data.priority = 15

**Verificacion MCP:**
```sql
SELECT title, message, priority FROM announcements WHERE id = {ANNOUNCEMENT_ID_1} AND deleted_at IS NULL;
```
- [ ] title = 'E2E Test - Bienvenida (Editado)'

## Caso 4.2: Actualizar agregando links

**Request:**
```
PUT /api/v1/announcements/{ANNOUNCEMENT_ID_1}
Content-Type: application/json

{
  "category_id": 3,
  "title": "E2E Test - Bienvenida (Editado)",
  "message": "Ahora con links",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 15,
  "is_global": true,
  "links": [
    { "label": "Nuevo link", "url": "https://new.example.com", "sort_order": 0 }
  ],
  "target_ids": []
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] data.links tiene 1 elemento (los links se reemplazan, no se acumulan)

**Verificacion MCP:**
```sql
SELECT COUNT(*) FROM announcement_links WHERE announcement_id = {ANNOUNCEMENT_ID_1};
```
- [ ] count = 1

## Caso 4.3: Cambiar de global a segmentado

**Request:**
```
PUT /api/v1/announcements/{ANNOUNCEMENT_ID_1}
Content-Type: application/json

{
  "category_id": 3,
  "title": "E2E Test - Bienvenida (Editado)",
  "message": "Ahora segmentado",
  "display_type": "modal_text",
  "frequency_type": "once",
  "priority": 15,
  "is_global": false,
  "links": [],
  "target_ids": [1]
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] data.is_global = false
- [ ] data.targets tiene 1 elemento con business_id = 1

**Verificacion MCP:**
```sql
SELECT business_id FROM announcement_targets WHERE announcement_id = {ANNOUNCEMENT_ID_1};
```
- [ ] 1 registro con business_id = 1

## Caso 4.4: Error - actualizar anuncio inexistente

**Request:**
```
PUT /api/v1/announcements/999999
Content-Type: application/json

{
  "category_id": 1,
  "title": "No existe",
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

## Caso 4.5: Error - actualizar con display_type invalido

**Request:**
```
PUT /api/v1/announcements/{ANNOUNCEMENT_ID_1}

{ ... "display_type": "banner" ... }
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "invalid display type"
