# Modulo Shipments

Gestion de envios: cotizacion, generacion de guia con carrier, tracking, COD, y
renderizado de guias en formatos amigables para el vendedor.

## Responsabilidad

- CRUD de `Shipment` (entidad principal, una por orden o manual).
- Direcciones de origen del vendedor (`origin_addresses`).
- Cotizacion y generacion de guia via integracion de transporte
  (envioclick por defecto). Comunicacion asincrona RabbitMQ.
- Sincronizacion de estados (webhook + sync manual).
- Cobro COD (contra entrega).
- Render de guias por formato (ver "Guias por formato" abajo).
- Estadisticas por geozona.

Cross-module: el modulo **NO** importa repositorios de otros modulos.
Cuando necesita data de `orders`, `business`, `warehouses`, `geozones`,
replica los SELECT en `infra/secondary/repository/*_queries.go` (regla de
aislamiento del proyecto, ver `.claude/rules/backend-conventions.md`).

## Arquitectura

```
shipments/
+-- bundle.go                          # Ensamblaje (repo + uc + handlers + consumer)
+-- internal/
    +-- domain/                        # Entidades, ports, errores
    |   +-- shipment.go                # Entidad Shipment del dominio
    |   +-- entities.go                # DTOs request/response, CarrierInfo, etc
    |   +-- guide_format.go            # Entidad GuideFormat + constantes de estrategia
    |   +-- ports.go                   # IRepository, ICarrierResolver, IPDFUploader, ISSEPublisher, ITransportRequestPublisher
    |
    +-- app/                           # Casos de uso (logica)
    |   +-- usecases/constructor.go    # Aggregator (ShipmentCRUD + OriginAddress)
    |   +-- usecaseshipment/           # CRUD, COD, sync, render_guide
    |   +-- usecaseoriginaddress/      # CRUD de direcciones de origen
    |
    +-- infra/
        +-- primary/                   # Adaptadores de entrada
        |   +-- handlers/              # HTTP handlers Gin (uno por endpoint)
        |   +-- queue/consumer/        # Consumer de respuestas de transport (RabbitMQ)
        |
        +-- secondary/                 # Adaptadores de salida
            +-- repository/            # GORM, replica de queries de otros modulos
            +-- queue/                 # Publisher RabbitMQ (transport y SSE)
            +-- cache/                 # Lectura de shipping_margins en Redis
            +-- storage/               # Adapter de IS3Service para subir PDF
```

## Entidad principal: Shipment

Modelo GORM en `back/migration/shared/models/shipment.go`. Campos relevantes:

| Campo | Tipo | Uso |
|---|---|---|
| `OrderID` | uuid? | FK opcional a `orders` |
| `TrackingNumber`, `TrackingURL` | string | Identificacion del carrier |
| `Carrier`, `CarrierCode` | string | ENVIA, INTERRAPIDISIMO, etc |
| `GuideID`, `GuideURL` | string | ID y URL PDF que devuelve el carrier |
| `ProbabilityGuideURL` | string | (Reservado para cache S3 futuro, hoy no se escribe) |
| `Status` | string | pending / in_transit / delivered / failed / cancelled |
| `CarrierStatus`, `CarrierStatusDetail` | string | Estado homologado del carrier |
| `ShippingCost`, `InsuranceCost`, `TotalCost` | decimal | Costos cargados al wallet |
| `CarrierCost`, `AppliedMargin` | decimal | Lo que se lleva el carrier vs ganancia Probability |
| `CodCarrierFee`, `CodProbabilityMargin` | decimal | Comision COD partida (ver `project_cod_recaudo_report`) |
| `Weight`, `Height`, `Width`, `Length` | decimal | Dimensiones del paquete |
| `WarehouseID` | uint | Almacen origen |
| `IsTest` | bool | Shipment generado en modo testing |
| `Metadata` | jsonb | Eventos del carrier + envioclick_id_order |

## Endpoints HTTP

Prefijo: `/api/v1`. Todas las rutas bajo `/shipments` requieren JWT
excepto las publicas de tracking. Para super admin (`business_id = 0` en JWT),
pasar `?business_id=X` como query param.

### Publicas (sin JWT)

| Metodo | Path | Handler | Uso |
|---|---|---|---|
| GET | `/tracking/search?q=` | `PublicSearchTracking` | Busqueda publica por numero o nombre |
| GET | `/tracking/:tracking_number/history` | `PublicGetTrackingHistory` | Historial de un tracking |

### Shipments CRUD

| Metodo | Path | Uso |
|---|---|---|
| GET | `/shipments` | Listar paginado (filtros: status, carrier, business_id, date_from/to) |
| GET | `/shipments/:id` | Detalle |
| POST | `/shipments` | Crear manual (sin orden) |
| PUT | `/shipments/:id` | Actualizar |
| DELETE | `/shipments/:id` | Soft delete |
| GET | `/shipments/order/:order_id` | Shipments asociados a una orden |
| GET | `/shipments/tracking/:tracking_number` | Buscar por tracking |

### Operaciones con carrier

| Metodo | Path | Uso |
|---|---|---|
| POST | `/shipments/quote` | Cotizar envio (asincrono, devuelve correlation_id) |
| POST | `/shipments/generate` | Generar guia con carrier (asincrono) |
| POST | `/shipments/tracking/:tracking_number/track` | Forzar refresh de estado |
| POST | `/shipments/:id/cancel` | Cancelar guia |
| POST | `/shipments/cancel-batch` | Cancelar varias guias |
| POST | `/shipments/sync-status` | Sync masivo de estados desde carrier |

El flujo `quote`/`generate` publica request a RabbitMQ. La respuesta llega
al `infra/primary/queue/consumer` y se procesa actualizando el shipment +
emitiendo SSE.

### Guias por formato

Render on-demand del PDF de la guia adaptado a un formato del catalogo:

| Metodo | Path | Uso |
|---|---|---|
| GET | `/shipments/guide-formats[?carrier=ENVIA]` | Lista de formatos disponibles (para dropdown UI) |
| GET | `/shipments/:id/guide[?format=<code>][&download=1]` | Stream del PDF (`application/pdf`). Sin `format=`, usa el default del carrier. Con `download=1`, header `attachment`. |

Ver "Catalogo guide_formats" abajo.

### Direcciones de origen

| Metodo | Path | Uso |
|---|---|---|
| GET | `/shipments/origin-addresses` | Listar |
| POST | `/shipments/origin-addresses` | Crear |
| PUT | `/shipments/origin-addresses/:id` | Actualizar |
| DELETE | `/shipments/origin-addresses/:id` | Eliminar |

### COD (contra entrega)

| Metodo | Path | Uso |
|---|---|---|
| GET | `/shipments/cod` | Listar shipments COD pendientes/cobrados (filtros: status, is_paid, business_id) |
| POST | `/shipments/:id/collect-cod` | Marcar como cobrado (admin) |

### Estadisticas

| Metodo | Path | Uso |
|---|---|---|
| GET | `/shipments/stats/by-geozone` | Conteo + tasa de exito por geozona |

## Catalogo guide_formats

Tabla `guide_formats` en `back/migration/shared/models/guide_format.go`.
Sembrada en `back/migration/internal/infra/repository/migrate_guide_formats.go`
(idempotente, soft-delete de obsoletos como `interrapidisimo-compact`).

### Campos

```
carrier        string    -- ENVIA, INTERRAPIDISIMO, COORDINADORA, ENVIOCLICK, TCC, SERVIENTREGA, 99MINUTOS, DEPRISA
code           string UQ -- envia-compact, envia-10x15, etc
label          string    -- "Compacta (1 guia)", "Adhesiva 10x15 cm"
width_cm       decimal   -- ancho del PDF resultante
height_cm      decimal   -- alto del PDF resultante
adhesive       bool      -- si es para impresion adhesiva
strategy       string    -- passthrough | crop | resize
crop_*_frac    decimal   -- region a recortar del original, en fracciones (0..1) del page
source_page    int       -- pagina del PDF original a tomar (default 1)
is_default     bool      -- formato por defecto cuando no se pasa ?format=
sort_order     int       -- orden en la UI
is_active      bool      -- soft toggle
```

### Estrategias

| Estrategia | Que hace |
|---|---|
| `passthrough` | Descarga el PDF del carrier y lo devuelve sin cambios. Ej: `envia-original`. |
| `crop` | Aplica un CropBox/MediaBox con `pdfcpu` y trim. Ej: `envia-compact` extrae el tercio superior de las 3 copias de ENVIA. |
| `resize` | Crop + redimensiona la pagina a `width_cm x height_cm`. Auto-swap de dimensiones si source es horizontal y target vertical (o viceversa). Ej: `envia-10x15` adhesiva. |

### Catalogo inicial sembrado

```
ENVIA           original (passthrough 21.6x27.9) | compact* (crop 21.6x9.3) | 10x15 (resize)
COORDINADORA    original | compact* | 10x15
ENVIOCLICK      original | compact* | 10x15        (alias de Coordinadora)
INTERRAPIDISIMO original* (A4 21x29.7) | 10x15
TCC             original | compact* (top 30%) | 10x15
SERVIENTREGA    original (4 copias) | compact* (mitad sup) | 10x15
99MINUTOS       original* (ya 1 guia) | 10x15
DEPRISA         original* (ya 10x11) | 10x15
```
(`*` = `is_default=true`)

## Usecase render_guide

`internal/app/usecaseshipment/render_guide.go`. Flujo on-demand sin cache:

```
1. GetShipmentByID(id)         -- domain
2. Resolver format             -- ?format= o GetDefaultGuideFormat(carrier)
3. http.Get(shipment.GuideURL) -- descarga el PDF original del carrier S3
4. applyGuideFormat(pdf, fmt)  -- segun fmt.Strategy:
     passthrough -> noop
     crop        -> pdfcpu.AddBoxes(CropBox=region) + pdfcpu.Trim
     resize      -> crop + pdfcpu.Resize a dim:WxH (auto-swap orientation)
5. return PDF bytes (no se guarda en S3)
```

Estado actual: **on-demand puro, sin almacenamiento**. Cada request descarga
y reprocesa. Si en el futuro se requiere cache, agregar tabla
`guide_renders (shipment_id, format_code, s3_url)` y el campo
`Shipment.ProbabilityGuideURL` ya existe en el modelo (reservado).

## Comunicacion asincrona

El flujo "generar guia" no es sincrono porque el carrier puede tardar
segundos. Se usa RabbitMQ:

```
Handler GenerateGuide
  -> repo.GetActiveShippingCarrier(business_id)
  -> shipment placeholder en DB (status=pending)
  -> transportPub.Publish(TransportRequest{op:generate, carrier, payload}) -> RabbitMQ
  -> 202 Accepted con correlation_id

[envioclick worker en otro container]
  -> POST a la API de envioclick
  -> Publica response a queue de respuestas

infra/primary/queue/consumer/response_consumer.go
  -> Lee response
  -> repo.UpdateShipment(GuideURL, TrackingNumber, Status...)
  -> ssePublisher.PublishGuideGenerated(...) -> frontend
```

El frontend escucha por SSE (`/api/v1/notify/sse/...`) y actualiza la UI.

## Carrier resolution

El handler recibe `business_id` (de JWT o query param). Llama
`repo.GetActiveShippingCarrier(businessID)` que selecciona la integracion
de transporte activa para ese business desde tabla `integrations`. Si no
hay carrier activo, falla con error claro.

Esto evita hardcodear el carrier por business. Cada negocio puede tener su
propia integracion (envioclick, mensajeros locales, etc).

## Aislamiento de modulos

Este modulo NO importa nada de `services/modules/orders`, `business`, etc.
Las queries que cruzan tablas estan en `repository/*_queries.go`:

- `order_queries.go`: GetOrderBusinessID, GetOrderPublicTrackingByNumber, etc
- `carrier_queries.go`: GetActiveShippingCarrier, GetBusinessName
- `cod_queries.go`: agregaciones para listado COD
- `geozone_queries.go`: GetShipmentStatsByGeozone, ResolveShipmentGeozone
- `guide_format_queries.go`: catalogo de formatos
- `guide_pdf_context.go`: JOIN para contexto del PDF (preparado para overlay/branding futuro)
- `sync_queries.go`: bulk update de estados
- `wallet_queries.go`: cargo al wallet al generar guia

## Dependencias clave

- `github.com/jung-kurt/gofpdf` v1 — generacion de PDF (reservado para overlay/branding futuro)
- `github.com/pdfcpu/pdfcpu` v0.12 — manipulacion de PDF (crop, resize, trim) **requiere Go 1.25**
- `github.com/skip2/go-qrcode` — QR codes
- `github.com/boombuler/barcode` — Code128 (reservado para futuro)

## Reglas del proyecto que aplican aqui

Ver `.claude/rules/`:
- `architecture.md` — hexagonal estricta
- `backend-conventions.md` — aislamiento, migraciones, logging, super admin
- `testing.md` — tests E2E en `.claude/testing/shipments/`

## Cosas pendientes / backlog

- Cache S3 de PDFs renderizados (`probability_guide_url`)
- Masking de "ENVIOCLICK COLOMBIA S.A.S" en PDFs originales (ver task #24)
- Overlay con branding del business (logo, NIT, color)
- Frontend: preferencia de formato favorito por business
- Validacion de URL mock (`envioclick-mock.local`, `testing-guias`) para devolver mensaje amigable
