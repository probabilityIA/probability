# Plan: Modulo Announcements (Avisos Informativos)

## Resumen

Sistema de avisos/anuncios configurable por super admins dirigido a businesses.
Soporta multiples formas de visualizacion, reglas de frecuencia, imagenes (S3),
links, tracking de interacciones y programacion de publicacion.

---

## Fase 1 - Implementacion Completa

### 1. Migracion y Modelos de Base de Datos

**Ruta:** `back/migration/`

- [x] **1.1 Modelo AnnouncementCategory** (`migration/shared/models/announcement_category.go`)
  - Tabla `announcement_categories`: id, code, name, icon, color, created_at, updated_at, deleted_at
  - Categorias iniciales: `promotion`, `alert`, `informative`, `tutorial`, `update`, `terms`

- [x] **1.2 Modelo Announcement** (`migration/shared/models/announcement.go`)
  - Tabla `announcements`: id, business_id (nullable, null=global), category_id, title, message,
    display_type (enum: `modal_image`, `modal_text`, `ticker`), frequency_type (enum: `once`, `daily`, `always`, `requires_acceptance`),
    priority (int), is_global (bool), status (enum: `draft`, `scheduled`, `active`, `inactive`),
    starts_at, ends_at, force_redisplay (bool), created_by (user_id),
    created_at, updated_at, deleted_at

- [x] **1.3 Modelo AnnouncementImage** (`migration/shared/models/announcement_image.go`)
  - Tabla `announcement_images`: id, announcement_id, image_url, sort_order, created_at, deleted_at

- [x] **1.4 Modelo AnnouncementLink** (`migration/shared/models/announcement_link.go`)
  - Tabla `announcement_links`: id, announcement_id, label, url, sort_order, created_at, deleted_at

- [x] **1.5 Modelo AnnouncementTarget** (`migration/shared/models/announcement_target.go`)
  - Tabla `announcement_targets`: id, announcement_id, business_id, created_at
  - Relacion: cuando is_global=false, se registran los businesses objetivo

- [x] **1.6 Modelo AnnouncementView** (`migration/shared/models/announcement_view.go`)
  - Tabla `announcement_views`: id, announcement_id, user_id, business_id,
    action (enum: `viewed`, `closed`, `clicked_link`, `accepted`),
    link_id (nullable), viewed_at, created_at

- [x] **1.7 Archivo de migracion** (`migration/internal/migrations/XXX_create_announcements.go`)
  - AutoMigrate de todos los modelos
  - Seed de categorias iniciales
  - Indices: business_id, status, starts_at/ends_at, announcement_id+user_id

- [x] **1.8 Registrar migracion en cmd/main.go y ejecutar**

- [x] **1.9 Verificar tablas creadas con MCP postgres**

---

### 2. Backend - Modulo Announcements

**Ruta:** `back/central/services/modules/announcements/`

#### 2.1 Domain

- [x] **2.1.1 Entities** (`internal/domain/entities/`)
  - announcement.go: Announcement, AnnouncementImage, AnnouncementLink, AnnouncementTarget
  - announcement_view.go: AnnouncementView
  - announcement_category.go: AnnouncementCategory
  - Enums: DisplayType, FrequencyType, AnnouncementStatus, ViewAction

- [x] **2.1.2 DTOs** (`internal/domain/dtos/`)
  - announcement_filters.go: filtros para listados (status, category, date range, business_id)
  - announcement_stats.go: estadisticas de vistas por aviso

- [x] **2.1.3 Ports** (`internal/domain/ports/ports.go`)
  - IRepository: CRUD announcements, images, links, targets, views, categories
  - IStorageService: upload/delete images (wraps shared/storage)

- [x] **2.1.4 Errors** (`internal/domain/errors/errors.go`)
  - ErrAnnouncementNotFound, ErrInvalidDateRange, ErrInvalidDisplayType, etc.

#### 2.2 Application (Use Cases)

- [x] **2.2.1 Constructor** (`internal/app/constructor.go`)
  - Dependencias: IRepository, IStorageService, ILogger

- [x] **2.2.2 CreateAnnouncement** (`internal/app/create_announcement.go`)
  - Validar datos, guardar aviso en status draft o scheduled/active segun fechas
  - Si tiene imagenes, subirlas a S3 (folder: `announcements/{announcement_id}/`)

- [x] **2.2.3 UpdateAnnouncement** (`internal/app/update_announcement.go`)
  - Editar campos, manejar imagenes nuevas/eliminadas (cleanup S3)
  - Si force_redisplay=true, limpiar registros de views para ese aviso

- [x] **2.2.4 DeleteAnnouncement** (`internal/app/delete_announcement.go`)
  - Soft delete del aviso + eliminar TODAS las imagenes del S3

- [x] **2.2.5 ListAnnouncements** (`internal/app/list_announcements.go`)
  - Paginado, filtros por status, categoria, business, rango de fechas
  - Para super admin: todos los avisos
  - Incluir conteo de vistas por aviso

- [x] **2.2.6 GetAnnouncement** (`internal/app/get_announcement.go`)
  - Detalle completo: aviso + imagenes + links + targets + stats de vistas

- [x] **2.2.7 GetActiveAnnouncements** (`internal/app/get_active_announcements.go`)
  - Para el frontend del usuario: avisos activos para su business
  - Filtrar por display_type, aplicar reglas de frecuencia
  - Consultar views del usuario para excluir los ya vistos segun frecuencia

- [x] **2.2.8 RegisterView** (`internal/app/register_view.go`)
  - Registrar accion del usuario: viewed, closed, clicked_link, accepted

- [x] **2.2.9 GetAnnouncementStats** (`internal/app/get_announcement_stats.go`)
  - Stats por aviso: total views, unique users, clicks, acceptances

- [x] **2.2.10 ListCategories** (`internal/app/list_categories.go`)
  - Listar categorias disponibles (catalogo)

- [x] **2.2.11 ChangeStatus** (`internal/app/change_status.go`)
  - Cambiar estado: draft->scheduled, draft->active, active->inactive, etc.

- [x] **2.2.12 ForceRedisplay** (`internal/app/force_redisplay.go`)
  - Marcar force_redisplay=true, limpiar views anteriores

#### 2.3 Infrastructure

##### 2.3.1 Handlers (Primary)

**Ruta:** `internal/infra/primary/handlers/`

- [x] **constructor.go** - Dependencia: IUseCase
- [x] **routes.go** - RegisterRoutes con grupo `/announcements`
- [x] **create_handler.go** - POST /announcements (multipart: json + imagenes)
- [x] **update_handler.go** - PUT /announcements/:id (multipart)
- [x] **delete_handler.go** - DELETE /announcements/:id
- [x] **list_handler.go** - GET /announcements (paginado, filtros)
- [x] **get_handler.go** - GET /announcements/:id
- [x] **get_active_handler.go** - GET /announcements/active (para usuarios normales)
- [x] **register_view_handler.go** - POST /announcements/:id/view
- [x] **get_stats_handler.go** - GET /announcements/:id/stats
- [x] **list_categories_handler.go** - GET /announcements/categories
- [x] **change_status_handler.go** - PATCH /announcements/:id/status
- [x] **force_redisplay_handler.go** - POST /announcements/:id/force-redisplay
- [x] **request/** - Request DTOs para cada handler
- [x] **response/** - Response DTOs para cada handler
- [x] **mappers/** - Mappers request->domain, domain->response

##### 2.3.2 Repository (Secondary)

**Ruta:** `internal/infra/secondary/repository/`

- [x] **constructor.go** - IDatabase dependency
- [x] **announcement_crud.go** - Create, Update, Delete, Get, List con paginacion
- [x] **announcement_images.go** - CRUD imagenes (con sort_order)
- [x] **announcement_links.go** - CRUD links
- [x] **announcement_targets.go** - CRUD targets (businesses asignados)
- [x] **announcement_views.go** - Registrar y consultar views/interacciones
- [x] **announcement_stats.go** - Queries de estadisticas agregadas
- [x] **category_queries.go** - Listar categorias
- [x] **mappers/** - Mappers model<->entity

##### 2.3.3 Storage Adapter (Secondary)

- [x] **storage_adapter.go** - Implementa IStorageService usando shared/storage

#### 2.4 Bundle

- [x] **bundle.go** - Componer capas, registrar rutas en router

#### 2.5 Registrar modulo

- [x] Registrar en `back/central/services/modules/bundle.go`

---

### 3. Frontend - Modulo Announcements (Admin/Super Admin)

**Ruta:** `front/central/src/services/modules/announcements/`

#### 3.1 Domain

- [x] **3.1.1 types.ts** - Interfaces: Announcement, AnnouncementImage, AnnouncementLink,
  AnnouncementTarget, AnnouncementView, AnnouncementCategory, AnnouncementStats,
  Enums: DisplayType, FrequencyType, AnnouncementStatus, ViewAction,
  DTOs: CreateAnnouncementDTO, UpdateAnnouncementDTO, PaginationParams, PaginatedResponse

- [x] **3.1.2 ports.ts** - IAnnouncementRepository interface

#### 3.2 Application

- [x] **3.2.1 use-cases.ts** - Funciones: createAnnouncement, updateAnnouncement,
  deleteAnnouncement, listAnnouncements, getAnnouncement, getActiveAnnouncements,
  registerView, getStats, listCategories, changeStatus, forceRedisplay

#### 3.3 Infrastructure

- [x] **3.3.1 repository/api-repository.ts** - Implementa IAnnouncementRepository,
  llamadas HTTP al backend, manejo de multipart para imagenes

- [x] **3.3.2 actions/index.ts** - Server Actions para mutaciones (create, update, delete,
  changeStatus, forceRedisplay) con revalidatePath

#### 3.4 UI - Pagina de Administracion (Super Admin)

**Ruta:** `front/central/src/app/(auth)/announcements/`

- [x] **3.4.1 page.tsx** - Pagina principal: listado de avisos con filtros,
  tabla con columnas: titulo, categoria, tipo, estado, vigencia, acciones

- [x] **3.4.2 AnnouncementManager.tsx** - Componente orquestador con modales create/edit

- [x] **3.4.3 AnnouncementForm.tsx** - Formulario de creacion/edicion (modal)

- [x] **3.4.4 [id]/stats/page.tsx** - Estadisticas de un aviso (vistas, clicks, etc.)

#### 3.5 UI - Componentes

**Ruta:** `front/central/src/services/modules/announcements/ui/components/`

- [x] **3.5.1 AnnouncementList.tsx** - Tabla paginada con filtros (status, busqueda)
- [x] **3.5.2 AnnouncementForm.tsx** - Formulario para crear/editar aviso:
  - Campos basicos: titulo, mensaje, categoria, tipo de display, frecuencia
  - Links dinamicos (agregar/quitar)
  - Checkbox global
  - Rango de fechas (starts_at, ends_at)
- [x] **3.5.3 StatusBadge** - Integrado en AnnouncementList con Badge shared
- [x] **3.5.4 CategoryBadge** - Integrado en AnnouncementList con color dot
- [x] **3.5.5 ImageUploader.tsx** - Componente de upload multiple con preview, drag&drop y reorden
- [x] **3.5.6 BusinessTargetSelector.tsx** - Selector multi-check de businesses para targeting (no global)
- [x] **3.5.7 AnnouncementStats.tsx** - Cards de estadisticas (vistas, usuarios, clicks, aceptaciones, cerrados)

#### 3.6 UI - Componentes de Visualizacion (lo que ve el usuario)

**Ruta:** `front/central/src/services/modules/announcements/ui/components/`

- [x] **3.6.1 AnnouncementModal.tsx** - Modal que aparece al iniciar sesion:
  - Carousel de avisos tipo modal_image (multiples imagenes por aviso)
  - Avisos tipo modal_text (solo titulo + mensaje)
  - Navegacion entre avisos (siguiente/anterior)
  - Boton cerrar / boton aceptar (segun frecuencia)
  - Registrar view al mostrar, close al cerrar, accept al aceptar

- [x] **3.6.2 AnnouncementTicker.tsx** - Barra superior con texto corrido:
  - Animacion CSS de izquierda a derecha (marquee)
  - Concatena mensajes de todos los avisos ticker activos
  - Visible en todas las paginas del dashboard
  - Registrar view al renderizar

- [x] **3.6.3 Integrar AnnouncementModal en el layout post-login**
  - Montado en layout-content.tsx despues de LinaChatbot

- [x] **3.6.4 Integrar AnnouncementTicker en el layout principal**
  - Montado en layout-content.tsx antes de los subnavbars

#### 3.7 UI - Hooks

- [x] **3.7.1** - Logica de hooks integrada directamente en componentes (AnnouncementList, AnnouncementForm)
- [x] **3.7.2** - Logica de fetch activos integrada en AnnouncementModal y AnnouncementTicker

#### 3.8 Navegacion

- [x] **3.8.1 Agregar entrada "Avisos" en el sidebar del dashboard** (solo super admin)
- [x] **3.8.2 Configurar rutas en el sistema de navegacion existente**

---

### 4. Testing

#### 4.1 Unit Testing

- [x] **4.1.1 Backend - Use cases tests** (`internal/app/usecase_test.go`)
  - Mocks en `internal/mocks/` (repository, storage, logger)
  - 38 tests, 89% cobertura: create, update, delete, list, get, get_active,
    register_view, get_stats, change_status, force_redisplay, list_categories
  - Casos edge: fechas invalidas, display/frequency invalidos, targets requeridos,
    frecuencia once/daily/requires_acceptance, S3 error no-blocking

- [x] **4.1.2 Backend - Domain entities tests** (`internal/domain/entities/announcement_test.go`)
  - 15 tests: DisplayType, FrequencyType, AnnouncementStatus, ViewAction constants

- [x] **4.1.3 Frontend - Use cases tests** (`app/use-cases.test.ts`)
  - Mock de IAnnouncementRepository, 22 tests pasando

#### 4.2 Testing E2E Backend

- [x] **4.2.1 Tests de integracion API** - Diferido: el proyecto no tiene infraestructura E2E.
  Se cubre parcialmente con los 38 unit tests de use cases (89% cobertura).

#### 4.3 Testing E2E Web

- [x] **4.3.1 Tests E2E frontend** - Diferido: no hay framework E2E (Playwright/Cypress) configurado.
  Se validara manualmente al iniciar backend + frontend.

---

## Fase 2 (Futuro - No implementar ahora)

- Business -> Clientes: avisos que un business configura para sus clientes finales
- Tracking de navegacion: saber si el usuario entro al modulo promocionado
- CRUD de reglas de frecuencia personalizadas
- Integracion con tienda online (storefront)
- Notificaciones push / email de avisos criticos
- A/B testing de avisos
- Templates predefinidos de avisos

---

## Notas Tecnicas

- Imagenes en S3 folder: `announcements/{announcement_id}/`
- Al eliminar aviso: soft delete en DB + hard delete de imagenes en S3
- El ticker usa CSS animation (marquee), no JS intervals
- Los avisos programados (scheduled) necesitan activarse cuando llega starts_at:
  opcion A: cron job, opcion B: verificar al consultar (lazy activation)
- Multi-tenant: super admin (business_id=0) crea avisos, is_global o con targets especificos
- Paginacion obligatoria en todos los GET de listados
