# CU-07: Obtener anuncios activos (endpoint de usuario)

## Endpoint
`GET /api/v1/announcements/active`

## Precondiciones
- CU-02 ejecutado (varios anuncios creados)
- CU-05 ejecutado (ANNOUNCEMENT_ID_1 en status "active")
- Token JWT de usuario con business_id > 0

## Caso 7.1: Obtener anuncios activos para un business

**Request:**
```
GET /api/v1/announcements/active
Authorization: Bearer {business_user_token}
```

**Nota:** El business_id y user_id se extraen del JWT, no van como query param.

**Verificaciones:**
- [ ] Status code = 200
- [ ] success = true
- [ ] data es un array
- [ ] Todos los items tienen status = "active" (implicito, el endpoint filtra)
- [ ] No aparecen anuncios con status "scheduled", "draft" o "inactive"
- [ ] Los anuncios globales (is_global=true) aparecen para cualquier business

## Caso 7.2: Anuncio segmentado solo aparece para business target

**Precondicion:** ANNOUNCEMENT_ID_4 (no global, targets: [1, 2])

Si el token tiene business_id = 1:
- [ ] ANNOUNCEMENT_ID_4 aparece en la lista

Si el token tiene business_id = 99 (no es target):
- [ ] ANNOUNCEMENT_ID_4 NO aparece en la lista

## Caso 7.3: Frecuencia "once" - no aparece despues de verlo

**Precondicion:**
- ANNOUNCEMENT_ID_1 tiene frequency_type = "once"
- CU-06 registro un "viewed" para este usuario

**Verificaciones:**
- [ ] ANNOUNCEMENT_ID_1 NO aparece en /active (ya fue visto con frecuencia "once")

## Caso 7.4: Error - business_id = 0 (super admin sin business seleccionado)

**Request:**
```
GET /api/v1/announcements/active
Authorization: Bearer {super_admin_token}
```

**Verificaciones:**
- [ ] Status code = 400
- [ ] message contiene "business_id is required"
