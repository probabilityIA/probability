# CU-03: Listar y obtener anuncios

## Endpoints
- `GET /api/v1/announcements` (listado paginado)
- `GET /api/v1/announcements/:id` (detalle)

## Precondiciones
- CU-02 ejecutado (existen al menos 4 anuncios de prueba)
- Token JWT valido

## Caso 3.1: Listar anuncios sin filtros (paginacion default)

**Request:**
```
GET /api/v1/announcements
Authorization: Bearer {token}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] Respuesta contiene: data (array), total (int), page (int), page_size (int), total_pages (int)
- [ ] page = 1
- [ ] page_size = 20
- [ ] total >= 4 (los creados en CU-02)
- [ ] data.length <= page_size

## Caso 3.2: Listar con paginacion especifica

**Request:**
```
GET /api/v1/announcements?page=1&page_size=2
```

**Verificaciones:**
- [ ] data.length <= 2
- [ ] total_pages >= 2 (si total >= 4)
- [ ] page = 1

## Caso 3.3: Listar filtrando por status

**Request:**
```
GET /api/v1/announcements?status=scheduled
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] Todos los items en data tienen status = "scheduled"
- [ ] data contiene el anuncio "E2E Test - Programado futuro" (ANNOUNCEMENT_ID_3)

## Caso 3.4: Listar filtrando por busqueda

**Request:**
```
GET /api/v1/announcements?search=Bienvenida
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] data contiene al menos 1 item con title que contiene "Bienvenida"

## Caso 3.5: Listar con page_size > 100 (se normaliza a 20)

**Request:**
```
GET /api/v1/announcements?page_size=500
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] page_size = 20 (normalizado)

## Caso 3.6: Obtener anuncio por ID

**Request:**
```
GET /api/v1/announcements/{ANNOUNCEMENT_ID_2}
```

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] data.id = ANNOUNCEMENT_ID_2
- [ ] data.title = "E2E Test - Promo con links"
- [ ] data.links tiene 2 elementos
- [ ] data.category es un objeto con id, code, name, icon, color

## Caso 3.7: Obtener anuncio inexistente

**Request:**
```
GET /api/v1/announcements/999999
```

**Verificaciones:**
- [ ] Status code = 404
- [ ] success = false

## Caso 3.8: ID invalido (no numerico)

**Request:**
```
GET /api/v1/announcements/abc
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] success = false
- [ ] message contiene "invalid id"
