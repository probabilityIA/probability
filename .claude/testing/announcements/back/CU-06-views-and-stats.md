# CU-06: Registrar vistas y consultar estadisticas

## Endpoints
- `POST /api/v1/announcements/:id/view`
- `GET /api/v1/announcements/:id/stats`

## Precondiciones
- CU-02 ejecutado (ANNOUNCEMENT_ID_1 existe)
- CU-05 ejecutado (ANNOUNCEMENT_ID_1 esta en status "active" o reactivado)
- Token JWT valido (user_id y business_id se extraen del JWT)

## Caso 6.1: Registrar view "viewed"

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_1}/view
Content-Type: application/json

{ "action": "viewed" }
```

**Verificaciones:**
- [ ] Status code = 201
- [ ] success = true
- [ ] message = "view registered"

**Verificacion MCP:**
```sql
SELECT id, announcement_id, user_id, action, viewed_at FROM announcement_views
WHERE announcement_id = {ANNOUNCEMENT_ID_1} AND action = 'viewed'
ORDER BY id DESC LIMIT 1;
```
- [ ] 1 registro con action = 'viewed'
- [ ] viewed_at no es null

## Caso 6.2: Registrar view "closed"

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_1}/view

{ "action": "closed" }
```

**Verificaciones:**
- [ ] Status code = 201

## Caso 6.3: Registrar view "clicked_link" con link_id

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_2}/view

{ "action": "clicked_link", "link_id": {LINK_ID_1} }
```

**Nota:** Usar el link_id obtenido de la respuesta de ANNOUNCEMENT_ID_2 en CU-02.

**Verificaciones:**
- [ ] Status code = 201

**Verificacion MCP:**
```sql
SELECT action, link_id FROM announcement_views
WHERE announcement_id = {ANNOUNCEMENT_ID_2} AND action = 'clicked_link';
```
- [ ] 1 registro con link_id = {LINK_ID_1}

## Caso 6.4: Registrar view "accepted"

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_4}/view

{ "action": "accepted" }
```

**Verificaciones:**
- [ ] Status code = 201

## Caso 6.5: Obtener estadisticas

**Request:**
```
GET /api/v1/announcements/{ANNOUNCEMENT_ID_1}/stats
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] data contiene: total_views, unique_users, total_clicks, total_acceptances, total_closed
- [ ] data.total_views >= 1 (por el viewed de caso 6.1)
- [ ] data.total_closed >= 1 (por el closed de caso 6.2)

**Verificacion MCP:**
```sql
SELECT
  COUNT(*) as total_views,
  COUNT(DISTINCT user_id) as unique_users
FROM announcement_views
WHERE announcement_id = {ANNOUNCEMENT_ID_1};
```
- [ ] total_views coincide con data.total_views

## Caso 6.6: Error - action invalido

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_1}/view

{ "action": "liked" }
```

**Verificaciones:**
- [ ] Status code = 400 o 201 (depende de si el backend valida el enum)

## Caso 6.7: Error - sin action (campo required)

**Request:**
```
POST /api/v1/announcements/{ANNOUNCEMENT_ID_1}/view

{}
```

**Verificaciones:**
- [ ] Status code = 400
