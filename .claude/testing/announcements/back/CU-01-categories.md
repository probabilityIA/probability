# CU-01: Listar categorias de anuncios

## Endpoint
`GET /api/v1/announcements/categories`

## Precondiciones
- Backend corriendo
- Token JWT valido (cualquier usuario autenticado)
- Migracion de announcements ejecutada (seed de categorias)

## Caso 1.1: Listar todas las categorias

**Request:**
```
GET /api/v1/announcements/categories
Authorization: Bearer {token}
```

**Expected Response (200):**
```json
{
  "success": true,
  "data": [
    { "id": 1, "code": "promotion", "name": "Promocion", "icon": "tag", "color": "#10b981" },
    { "id": 2, "code": "alert", "name": "Alerta", "icon": "alert-triangle", "color": "#ef4444" },
    { "id": 3, "code": "informative", "name": "Informativo", "icon": "info", "color": "#3b82f6" },
    { "id": 4, "code": "tutorial", "name": "Tutorial", "icon": "book-open", "color": "#8b5cf6" },
    { "id": 5, "code": "update", "name": "Actualizacion", "icon": "refresh-cw", "color": "#f59e0b" },
    { "id": 6, "code": "terms", "name": "Terminos y Condiciones", "icon": "file-text", "color": "#6b7280" }
  ]
}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] data tiene 6 elementos
- [ ] Cada categoria tiene: id, code, name, icon, color

**Verificacion MCP:**
```sql
SELECT id, code, name, icon, color FROM announcement_categories WHERE deleted_at IS NULL ORDER BY id;
```
- [ ] 6 registros coinciden con la respuesta

## Caso 1.2: Sin token - Acceso denegado

**Request:**
```
GET /api/v1/announcements/categories
(sin Authorization header)
```

**Expected Response (401):**
- [ ] Status code = 401 o 403
