# Resultados E2E - Announcements Backend

**Ultima ejecucion:** 2026-04-12
**Resumen:** 44 PASS / 3 FAIL de 47 casos

---

## CU-01: Listar categorias

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 1.1 Listar todas las categorias | PASS | 200, 6 categorias, DB confirma |
| 1.2 Sin token | PASS | 401 |

## CU-02: Crear anuncios

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 2.1 Crear global modal_text | PASS | 201, status=active, is_global=true |
| 2.2 Crear con links | **FAIL** | 201 pero links=[] en respuesta y DB. Bug: `Omit(clause.Associations)` en db.go |
| 2.3 Crear programado | PASS | 201, status=scheduled |
| 2.4 Crear no global con targets | **FAIL** | 201 pero targets=[] en DB. Mismo bug que 2.2 |
| 2.5 Error display_type invalido | PASS | 400 |
| 2.6 Error no global sin targets | PASS | 400 |
| 2.7 Error fechas invalidas | PASS | 400 |
| 2.8 Error campos requeridos | PASS | 400 |

## CU-03: Listar y obtener

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 3.1 Listar sin filtros | PASS | 200, total=4, page_size=20 |
| 3.2 Paginacion especifica | PASS | 2 items, total_pages=2 |
| 3.3 Filtrar por status | PASS | Solo scheduled |
| 3.4 Filtrar por busqueda | PASS | Coincide titulo |
| 3.5 page_size>100 normaliza | PASS | page_size=20 |
| 3.6 Obtener por ID | PASS | 200, category objeto completo |
| 3.7 Inexistente | PASS | 404 |
| 3.8 ID no numerico | PASS | 400 |

## CU-04: Actualizar

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 4.1 Actualizar titulo/mensaje | **FAIL parcial** | 200 pero priority no se actualiza. Falta en mapa Updates del repo |
| 4.2 Agregar links | PASS | 200, 1 link guardado (Update usa inserts separados) |
| 4.3 Cambiar a segmentado | PASS | 200, target guardado |
| 4.4 Error inexistente | PASS | 400 |
| 4.5 Error display_type invalido | PASS | 400 |

## CU-05: Cambiar estado

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 5.1 active -> inactive | PASS | 200, DB confirma |
| 5.2 inactive -> active | PASS | 200 |
| 5.3 -> draft | PASS | 200 |
| 5.4 -> scheduled | PASS | 200 |
| 5.5 Error status invalido | PASS | 400 |
| 5.6 Error inexistente | PASS | 400 |
| 5.7 Error sin campo status | PASS | 400 |

## CU-06: Vistas y estadisticas

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 6.1 Registrar viewed | PASS | 201, DB confirma viewed_at |
| 6.2 Registrar closed | PASS | 201 |
| 6.3 Registrar clicked_link | PASS | 201 |
| 6.4 Registrar accepted | PASS | 201 |
| 6.5 Obtener stats | PASS | 200, total_views=2, total_closed=1 |
| 6.6 Error action invalido | **FAIL** | 201 en vez de 400. No valida enum de action |
| 6.7 Error sin action | PASS | 400 |

## CU-07: Anuncios activos

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 7.1 Activos para business | PASS | 200, filtro de frecuencia funciona |
| 7.2 Segmentacion por business | PASS | Filtro de targets correcto |
| 7.3 Frecuencia once post-vista | PASS | No aparece despues de visto |
| 7.4 Super admin sin business | PASS | 400 "business_id is required" |

## CU-08: Force redisplay

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 8.1 Forzar redisplay | PASS | 200, force_redisplay=true, views=0 en DB |
| 8.2 Reaparece en /active | PASS | Anuncio visible de nuevo |
| 8.3 Error inexistente | PASS | 400 |

## CU-09: Eliminar

| Caso | Resultado | Detalle |
|------|-----------|---------|
| 9.1 Soft delete | PASS | 200, deleted_at IS NOT NULL |
| 9.2 No aparece en listado | PASS | Ausente |
| 9.3 No accesible por ID | PASS | 404 |
| 9.4 Error inexistente | PASS | 400 |
| 9.5 Delete con targets | PASS | 200 |
| 9.6 Cleanup | PASS | count=0 anuncios E2E |

---

## Bugs abiertos

### Bug 1 (CRITICO) - Links y targets no se persisten en Create
- **Archivo:** `back/central/shared/db/db.go` lineas 83-85
- **Causa:** `Omit(clause.Associations)` + `FullSaveAssociations: false` global
- **Casos:** CU-02.2, CU-02.4
- **Estado:** Abierto

### Bug 2 (MEDIO) - Campo priority no se actualiza en Update
- **Archivo:** `back/central/services/modules/announcements/internal/infra/secondary/repository/announcement_crud.go`
- **Causa:** `priority` ausente del mapa de `Updates`
- **Casos:** CU-04.1
- **Estado:** Abierto

### Bug 3 (BAJO) - RegisterView acepta actions invalidos
- **Archivo:** `register_view.go` o `register_view_handler.go`
- **Causa:** No valida enum ViewAction (viewed, closed, clicked_link, accepted)
- **Casos:** CU-06.6
- **Estado:** Abierto

---

## Observaciones

- Los test cases CU-02.4 y CU-04.3 usan target_ids [1,2] pero esos IDs no existen en esta BD. Se uso business_id=26 (Demo) durante la ejecucion.
- El Update de links/targets funciona correctamente porque usa inserts separados (ReplaceLinks, ReplaceTargets), no depende de GORM associations.
