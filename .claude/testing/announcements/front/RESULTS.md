# Resultados E2E - Announcements Frontend

**Ultima ejecucion:** 2026-04-14
**Herramienta:** MCP Playwright (Google)
**Resumen:** 15 PASS / 0 FAIL de 15 casos (re-ejecucion completa)

---

## Bugs encontrados y corregidos durante testing

### Bug F1 - Server Actions usaban cookies() directo en vez de getAuthToken
- **Archivo:** `front/central/src/services/modules/announcements/infra/actions/index.ts`
- **Causa:** Importaba `cookies` de `next/headers` directamente en vez de usar `getAuthToken` de `@/shared/utils/server-auth`
- **Efecto:** Token no se obtenia correctamente, causando 401 en todas las llamadas
- **Fix:** Reemplazar `cookies()` por `getAuthToken()` (mismo patron que orders, customers, etc.)
- **Estado:** Corregido

### Bug F2 - Repositorio no extraia .data del response wrapper del backend
- **Archivo:** `front/central/src/services/modules/announcements/infra/repository/api-repository.ts`
- **Causa:** El backend siempre retorna `{success, data: ...}` pero el repo retornaba el wrapper completo. Metodos como `listCategories()` esperaban un array pero recibian `{success, data: [...]}`
- **Efecto:** `categories.map is not a function` al abrir el formulario
- **Fix:** Cada metodo del repo que espera un tipo especifico ahora extrae `.data` del response
- **Estado:** Corregido

### Bug F3 - get_active_handler.go bloqueaba super admin sin business_id
- **Archivo:** `back/central/services/modules/announcements/internal/infra/primary/handlers/get_active_handler.go`
- **Causa:** El handler requeria `business_id > 0`, pero super admin tiene `business_id = 0`
- **Efecto:** Error 500 al cargar la pagina de announcements como super admin (AnnouncementModal/Ticker crasheaban)
- **Fix:** Eliminar validacion obligatoria de business_id. En el repositorio, solo filtrar por targets cuando `business_id > 0`
- **Estado:** Corregido

### Bug F4 - Backend no tenia rutas de imagen registradas (404 en POST image)
- **Causa:** El backend corriendo era la version anterior al commit de rutas de imagen
- **Efecto:** `POST /api/v1/announcements/7/image` retornaba 404
- **Fix:** Reiniciar backend para cargar las rutas nuevas (`POST /:id/image`, `DELETE /:id/image/:imageId`)
- **Estado:** Corregido

### Bug F5 - AnnouncementImage usaba soft delete (acumulaba basura en DB)
- **Archivo:** `back/migration/shared/models/announcement_image.go`
- **Causa:** El modelo usaba `gorm.Model` que incluye `DeletedAt` (soft delete)
- **Efecto:** Las imagenes eliminadas permanecian en la tabla con `deleted_at` set, acumulando registros basura
- **Fix:** Reemplazar `gorm.Model` por campos explicitos sin `DeletedAt`. Migracion purga soft-deleted y elimina columna `deleted_at`
- **Estado:** Corregido

---

## TC-F01: Login y navegacion

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Login super admin | PASS | Credenciales ${AI_SUPER_ADMIN_EMAIL}, redirige a /home |
| Sidebar muestra "Anuncios" | PASS | Link a /announcements visible |
| Navegar a /announcements | PASS | Pagina carga correctamente |

## TC-F02: Pagina de listado

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Titulo y subtitulo | PASS | "Anuncios" + "Gestiona avisos informativos para los negocios" |
| Business Selector | PASS | Dropdown SUPER ADMIN con negocios |
| Boton "Nuevo anuncio" | PASS | Visible con icono |
| Buscador | PASS | Input + boton "Buscar" |
| Filtro de estados | PASS | Todos, Activo, Programado, Borrador, Inactivo |
| Tabla columnas | PASS | Titulo, Categoria, Tipo, Estado, Vigencia, Acciones |
| Estado vacio | PASS | "No hay anuncios registrados" |

## TC-F03: Crear anuncio modal_text global

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Abrir formulario | PASS | Modal "Nuevo anuncio" con todos los campos |
| Categorias cargadas | PASS | 6 categorias del seed |
| Crear con titulo/mensaje/categoria | PASS | POST 201, anuncio aparece en tabla |
| Badge "Global" | PASS | Visible en columna titulo |
| Estado "Activo" | PASS | Badge activo mostrado |

## TC-F04: Crear anuncio ticker

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Seleccionar tipo Ticker | PASS | Dropdown "Ticker (barra superior)" |
| Crear ticker global | PASS | POST 201, aparece en tabla |
| Tipo "Ticker" en tabla | PASS | Columna tipo muestra "Ticker" |
| 2 anuncios en lista | PASS | Paginacion "2 registros totales" |

## TC-F05: Filtros de busqueda y status

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Buscar por titulo "Ticker" | PASS | Solo 1 resultado (ticker) |
| Filtro estado "Borrador" | PASS | 0 resultados (ambos son Activo) |
| Filtro estado "Activo" | PASS | 2 resultados |
| Limpiar filtros | PASS | Boton "Limpiar" restaura vista |

## TC-F06: Editar anuncio

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Abrir modal edicion | PASS | "Editar anuncio" con campos pre-llenados |
| Modificar titulo | PASS | Titulo actualizado a "Ticker EDITADO" |
| Guardar cambios | PASS | PUT 200, tabla refleja cambio |

## TC-F07: Cambiar estado

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Desactivar anuncio | PASS | Badge cambia a "Inactivo" |
| Activar anuncio | PASS | Badge vuelve a "Activo" |

## TC-F08: Estadisticas

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Navegar a /announcements/:id/stats | PASS | Pagina carga con titulo y categoria |
| Cards de metricas | PASS | 5 cards: Vistas, Usuarios, Clicks, Aceptaciones, Cerrados |
| Valores iniciales | PASS | Todos en 0 (anuncio recien creado) |

## TC-F09: AnnouncementModal

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Modal aparece al navegar | PASS | Muestra anuncio modal_text activo |
| Titulo y mensaje correctos | PASS | Contenido del anuncio visible |
| Boton "Cerrar" | PASS | Cierra el modal |

## TC-F10: AnnouncementTicker

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Barra ticker visible | PASS | Barra morada superior con texto |
| Contenido correcto | PASS | Titulo + mensaje del anuncio ticker |

## TC-F11: Eliminar anuncios

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Confirm dialog | PASS | Muestra titulo del anuncio en mensaje |
| Eliminar primer anuncio | PASS | Desaparece de tabla, queda 1 |
| Eliminar segundo anuncio | PASS | Tabla vacia, "No hay anuncios registrados" |

## TC-F12: Subir imagen al editar anuncio

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Dropzone visible en formulario | PASS | "Arrastra imagenes o haz click para seleccionar", "0 / 10 imagenes" |
| Seleccionar archivo via file chooser | PASS | Preview aparece, contador sube a "1 / 10 imagenes" |
| Guardar con imagen nueva | PASS | PUT 200 + POST /image 201 (1565ms), imagen subida a S3 |
| Imagen guardada en DB | PASS | `announcement_images` con URL S3 correcta |
| Imagen visible en AnnouncementModal | PASS | `img` renderizado con URL de S3 |

## TC-F13: Multiples imagenes y carrusel

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Subir segunda imagen | PASS | Banner subido, "2 / 10 imagenes" en formulario |
| 2 registros en DB | PASS | sort_order 0 y 1, URLs distintas |
| Carrusel en AnnouncementModal | PASS | Botones prev/next, dots indicadores para 2 imagenes |
| Imagenes existentes visibles al editar | PASS | "2 / 10 imagenes" con previews de ambas |

## TC-F14: Eliminar imagen existente

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Boton eliminar en imagen | PASS | Icono trash rojo visible en hover |
| Eliminar primera imagen del formulario | PASS | Desaparece, queda "1 / 10 imagenes" |
| Guardar con imagen eliminada | PASS | DELETE /image/1 200, hard delete en DB y S3 |
| Solo queda banner en DB | PASS | 1 registro activo (id=2, banner) |
| Sin columna deleted_at | PASS | Migracion elimino columna, hard delete real |

## TC-F15: Cambiar negocios objetivo

| Caso | Resultado | Detalle |
|------|-----------|---------|
| Quitar Mystic Rose y Seb-dev-1 | PASS | Checkboxes desmarcados |
| Solo Demo seleccionado | PASS | "1 negocio(s) seleccionado(s)" |
| Targets actualizados en DB | PASS | Solo business_id=26 (Demo) |

---

## Archivos modificados durante testing (bugfixes)

| Archivo | Cambio |
|---------|--------|
| `front/.../announcements/infra/actions/index.ts` | `cookies()` -> `getAuthToken()`, agregar uploadImageAction/deleteImageAction |
| `front/.../announcements/infra/repository/api-repository.ts` | Extraer `.data` del response wrapper, agregar uploadImage/deleteImage |
| `front/.../announcements/domain/ports.ts` | Agregar uploadImage/deleteImage a IAnnouncementRepository |
| `front/.../announcements/domain/types.ts` | Agregar UploadImageResponse/DeleteImageResponse |
| `front/.../announcements/app/use-cases.ts` | Agregar uploadImage/deleteImage |
| `front/.../announcements/app/use-cases.test.ts` | Agregar mocks y tests para upload/delete image |
| `front/.../announcements/ui/components/AnnouncementForm.tsx` | Wiring completo de ImageUploader con estado, upload y delete |
| `back/.../announcements/.../get_active_handler.go` | Eliminar validacion obligatoria de business_id |
| `back/.../announcements/.../repository/announcement_active.go` | Solo filtrar por targets si business_id > 0 |
| `back/migration/shared/models/announcement_image.go` | Reemplazar gorm.Model por campos explicitos (hard delete) |
| `back/migration/.../migrate_announcements.go` | Purgar soft-deleted, eliminar columna deleted_at |
