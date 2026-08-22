# Plan - Modulo "Tus Integraciones" en la app (paridad con front)

Worktree: `/home/cam/Desktop/probability-mobile` (branch `feat/mobile-app-flutter`)
App: `mobile/mobile_central` | Ruta actual: `/core` (boton central de la nav)
Referencia: `front/central/src/services/modules/my-integrations/` (8.089 lineas, 23 componentes)

Base de este plan: `origin/main` a `7eadafdc`, ya mezclado en la rama. Eso trae
`comparativo de ordenes` (PR #83) y el fix de packs/entrega de MercadoLibre
(PR #84), que agregan un ambiente nuevo al modulo.

## Que existe hoy en la app

| Pieza | Estado |
|---|---|
| `/core` con pestanas Diagrama / Informe | hecho |
| Nucleo + tarjeta por canal con stats | hecho |
| Informe = OrdersReport (Vista general) | hecho |
| Acciones por canal: activar, probar, sincronizar ordenes | hecho |
| Motor de sincronizacion (SSE + correlation_id) | hecho (Fase 0) |
| Comparar productos / Actualizar productos / Sincronizar inventario | **falta** |
| Comparar ordenes | **falta** (nuevo desde main) |
| Crear / editar integracion | **falta** |
| Detalle de corrida | **falta** |
| Facturar | deshabilitado tambien en front |

## Endpoints (todos existen, ninguno hay que crear)

Cortesia del backend actual. El front pasa por route handlers `/internal/*` de
Next solo para no exponer el token; la app pega directo con su JWT.

| Uso | Endpoint |
|---|---|
| Ultimas corridas | `GET /integrations/sync-runs` |
| Detalle de una corrida | `GET /integrations/sync-runs/items` |
| Hallazgos | `GET /integrations/sync-runs/findings` |
| Items de un hallazgo | `GET /integrations/sync-runs/findings/items` |
| Matriz producto x canal | `GET /integrations/sync-runs/matrix` |
| Resumen de datos del canal | `GET /integrations/sync-runs/data-summary` |
| Vista previa de datos | `POST /integrations/sync-runs/data-preview` |
| Aplicar / deshacer datos | `POST /products/channel-data/apply` , `/undo` , `GET /batches` |
| Comparar ordenes | `GET /orders-compare` , `POST /orders-compare/apply` |
| Modulos internos del nucleo | `GET /businesses/{id}/configured-resources` |
| Eventos en vivo | `GET /notify/sse/order-notify?event_types=&business_id=` (sin JWT) |

Por proveedor (prefijo distinto segun el canal):

| Canal | typeId | Prefijo | Acciones |
|---|---|---|---|
| Shopify | 1 | `/integrations/shopify` | sync, reconcile, apply, associate |
| Mercado Libre | 3 | `/integrations/meli` | sync, compare, reconcile, apply, associate |
| WooCommerce | 4 | `/woocommerce` | sync, compare, reconcile, apply, associate |
| Siigo | 8 | `/siigo` | sync, compare, reconcile/start, apply |
| VTEX | 16 | `/vtex` | sync, reconcile, apply, associate |
| Tiendanube | 17 | `/tiendanube` | sync, compare, reconcile, apply, associate |
| Jumpseller | 33 | `/jumpseller` | sync, reconcile, apply, associate |

## El comparativo guardado (regla central del modulo)

El front **no obliga a correr nada para ver algo**. Al abrir cualquier ambiente
muestra el comparativo anterior, guardado, y correr de nuevo es una accion
aparte:

| Ambiente | De donde sale lo que se ve al abrir |
|---|---|
| Comparar productos | `GET /sync-runs/findings` (trae `compared_at` por canal) y `GET /sync-runs/matrix` |
| Sincronizar inventario | `POST /{prefijo}/inventory/compare` con `source: 'snapshot'`, que responde `from_cache: true` y `checked_at` |
| Actualizar productos | `GET /sync-runs/data-summary`, que trae `snapshot_at` |
| Comparar ordenes | no hay guardado: siempre se le pregunta al canal |

Ademas `GET /integrations/sync-runs` da la ultima corrida por canal y tipo, y
`GET /sync-runs/items` su detalle paginado.

Esto no es un detalle de UI: preguntarle al canal cuesta segundos y cuota de
API. **Toda pantalla de este modulo abre con lo guardado y dice cuando se
comparo**; el boton de comparar solo refresca.

## Reglas que aplican a todas las fases

- `.claude/rules/mobile-listados-memoria.md`: toda coleccion va con
  `PaginatedListView` o `ListView.builder` con ventana. **Ninguna** de estas
  tablas se pinta con `Column`. Las matrices son las de mayor riesgo: filas de
  productos x columnas de canales, con imagen por fila.
- Miniaturas de producto y logos de canal SIEMPRE con `cacheWidth`/`cacheHeight`.
- Sin comentarios en `.dart`. Archivos solo ASCII (`\u00XX` para acentos).
- Las tablas anchas del front no se portan como tabla. En 400 px cada fila del
  front se vuelve una tarjeta y cada columna de canal un chip. Ese es el
  trabajo de diseno real de cada fase, no el porteo de la logica.

---

## Fase 0 - Motor de sincronizacion (cimiento, sin UI nueva) - HECHO

Es el `sync-activity-context.tsx` del front (718 lineas), la pieza de la que
cuelgan las fases 2, 4, 5 y 6. Sin esto cada accion queda sin progreso ni
resultado.

Que se construye:

- `lib/core/network/sse_client.dart` - Dio con `ResponseType.stream`, parseo de
  `data:` linea a linea, reconexion con backoff, cierre al salir de la pantalla.
  El endpoint SSE no lleva JWT, asi que no hay que resolver headers.
- `lib/services/modules/my_integrations/domain/sync_entities.dart` -
  `SyncNodeState`, `SyncResult` (inventory / products / error), `SyncRunRecord`,
  `SyncRunDetail`, `DetailGroup`, `ProductActionKey`.
- `lib/services/modules/my_integrations/app/sync_providers.dart` - el registro
  de los 7 canales de la tabla de arriba: prefijo, prefijo de evento, que
  acciones soporta, nombre del campo `only_in_*`.
- `infra/repository/sync_runs_repository.dart` - ultimas corridas e items.
- `ui/providers/sync_activity_provider.dart` - el controlador: mapa
  `correlation_id -> integration_id`, buffer de eventos que llegan antes del
  id (el front topa en 100 por corrida y 20 corridas), timeout de 90 s, y el
  rescate por `sync-runs` cuando el timeout gana la carrera.

Verificacion: lanzar un sync de inventario de un canal y ver los estados
`queued -> active -> done` con el progreso real, sin UI nueva mas que la
tarjeta del canal.

Riesgo: el `EventSource` del navegador reconecta solo; en Dart hay que
escribirlo. Si el stream resulta fragil sobre datos moviles, el plan B es
poller de `sync-runs` cada 3 s mientras haya una corrida viva. Decidir con la
prueba en el telefono, no antes.

## Fase 1 - Barra de acciones y ambientes - HECHO

Reemplaza el `DefaultTabController` de dos pestanas por lo que el front tiene
en el encabezado del modal (`SyncActions`).

- Selector de vista Diagrama | Informe (segmentado, arriba).
- Fila de ambientes desplazable: Vista general, Comparar productos, Actualizar
  productos, Sincronizar inventario, Comparar ordenes, Facturar (deshabilitado).
- El ambiente decide que se pinta en Informe. Vista general = el OrdersReport
  que ya existe.
- Boton Reiniciar cuando una corrida termino.
- Paneles vacios con el texto de ayuda de cada ambiente para las fases que
  todavia no existen (el front tiene ese texto en `ReportView`).

Al terminar la fase el modulo ya tiene su forma final; lo que sigue es llenar
cada panel.

## Fase 2 - Sincronizar inventario

La accion mas usada a diario y la de dato mas simple: un numero por SKU y canal.

Front: `InventoryMatrixTable` (630) + `InventoryCompareTable` (525) +
`InventoryCompareModal`.

Diseno movil:

- **Modo canal** (equivale a `InventoryCompareTable`): selector de canal,
  buscar, "solo diferencias", y lista paginada de productos. Cada producto una
  tarjeta: miniatura, SKU y nombre, y a la derecha `Probability 12 / Canal 8`
  con el delta en color. Seleccion multiple y "Enviar stock" que llama
  `runInventoryOne(id, skus)`.
- **Modo matriz** (equivale a `InventoryMatrixTable`): misma lista, pero cada
  tarjeta trae un chip por canal con su cantidad; el chip dice si el producto
  ni siquiera esta publicado ahi. Es la unica forma de que una matriz de N
  canales entre en 400 px.
- Barrido por grupos con barra de avance (`GRUPO_INVENTARIO = 100`), igual que
  el front.
- Progreso en vivo por SSE (Fase 0).

Endpoints: `POST /{prefijo}/inventory/compare` y `/inventory/sync`.
Ojo: Shopify, VTEX y Jumpseller **no** tienen `compare`. La UI debe decir
"este canal solo permite enviar, no comparar" en vez de fallar.

## Fase 3 - Comparar ordenes - HECHO

Lo que llego de main. Es el panel mas facil de los pesados: no necesita SSE, es
un GET paginado y un POST.

Front: `OrdersCompareTable` (446).

- Filtros: canal, desde, hasta, buscar, "solo diferencias".
- 5 KPIs: en el canal, faltan aqui, en las dos, estado distinto, sin mover stock.
- Lista de ordenes como tarjetas: numero + external_id, cliente, fecha, estado
  del canal vs estado aqui, total, y la etiqueta de situacion
  (falta en Probability / solo en Probability / estado distinto / en las dos).
- Casilla por orden solo en las que faltan, "Crear en Probability", y el aviso
  de que las historicas se crean sin mover inventario.
- Solo aplica a los canales de `ORDERS_COMPARE_TYPE_IDS = [1,3,4,16,17,33]`.

Endpoints: `GET /orders-compare`, `POST /orders-compare/apply`.

## Fase 4 - Comparar productos

Front: `MatchMatrixTable` (365) + `FindingItemsTable` (316) +
`GlobalProductsMatrix` (203) + los hallazgos.

- Correr `reconcile` en todos los canales a la vez (el front lo hace en
  paralelo) con estado por canal.
- Pildoras: "Matriz de productos" + una por hallazgo con su conteo y su color
  (error / warn / info), desde `GET /sync-runs/findings`.
- Matriz: tarjeta por producto con un chip por canal; el chip distingue
  presente, ausente, y presente-con-SKU-distinto.
- Hallazgos: el front tiene tres formas de fila (comparativa lado a lado,
  cruzada con "presente en / falta en", y simple). En movil las tres se
  resuelven con una tarjeta de dos columnas Canal | Probability.
- Acciones por producto: asociar, crear en el canal, crear en Probability,
  actualizar en Probability, y crear en ambos lados.

Endpoints: `POST /{prefijo}/products/reconcile`, `/products/apply`,
`/products/associate`, y `GET /sync-runs/findings`, `/findings/items`, `/matrix`.

## Fase 5 - Actualizar productos (datos del canal)

Front: `ChannelDataTable` (184) + `ChannelDataApplyModal` (268).

- Resumen: una fila por campo (nombre, imagen, categoria...) y una celda por
  canal con "puede llenar N" / "puede reemplazar N". En movil: tarjeta por
  campo, chips por canal.
- Al tocar una celda se abre la vista previa: cuantos llenaria, cuantos
  reemplazaria, conflictos entre canales y una muestra de "antes / despues".
- Aplicar deja un `batch_id`; ofrecer "Deshacer" mientras el lote este a la
  vista. Es la unica accion del modulo que se puede revertir, y hay que
  aprovecharlo.

Endpoints: `GET /sync-runs/data-summary`, `POST /sync-runs/data-preview`,
`POST /products/channel-data/apply` y `/undo`.

## Fase 6 - Diagrama completo

Hoy el diagrama de la app es una version reducida. Falta:

- Nucleo con los modulos internos activos, desde
  `GET /businesses/{id}/configured-resources`, y la insignia con el numero de
  hallazgos.
- Grupos de servicios (mensajeria, facturacion) aparte de los canales.
- Por tarjeta: interruptor de `inventory_sync_enabled`, editar la integracion,
  y el estado vivo de la corrida (progreso y resultado de la Fase 0).
- Enlaces animados entre canales y nucleo. En el front son SVG con seis
  animaciones a la vez; en movil va una version sobria: linea con un punto que
  viaja solo mientras hay una corrida activa. En reposo, lineas estaticas.
  Un telefono de gama baja no puede pagar esas animaciones todo el tiempo.

## Fase 7 - Crear y editar integracion

La de mayor incertidumbre, por eso va al final.

- Editar: formulario por proveedor sobre `PUT /integrations/{id}`.
- Crear: solo categorias `ecommerce` e `invoicing`, como el front.
- Los canales por token (WooCommerce, VTEX, Siigo, Shopify) son un formulario
  normal.
- Los de OAuth (Mercado Libre, Tiendanube, Jumpseller) necesitan abrir el
  navegador externo y volver por deep link. Eso es infraestructura nueva en la
  app (`AndroidManifest` con el intent-filter, y una pantalla de retorno).
  Si esto se complica, la salida digna es que la app diga "conecta este canal
  desde la web" y no bloquee el resto del modulo.

## Fase 8 - Detalle de corrida

Front: `SyncDetailPanel` (650).

- Desde cada tarjeta de canal, ver el detalle de la ultima corrida:
  `GET /sync-runs/items` paginado, filtrado por grupo (actualizados, omitidos,
  fallidos, solo en el canal, sin SKU, SKU cambiado...) y con buscador.
- Cada item con su tono ok / warn / error.

## Fase 9 - Facturar

Deshabilitado tambien en el front. Se deja el boton apagado con el mismo texto
("Facturacion desde el hub: proximamente") y no se implementa.

---

## Orden sugerido y por que

1. **Fase 0** primero: sin el motor, las fases 2, 4 y 5 quedan sin progreso.
2. **Fase 1** enseguida: da la forma final del modulo y permite entregar los
   paneles de a uno sin volver a tocar la navegacion.
3. **Fase 3 (comparar ordenes)** antes que las de producto: es la mas nueva,
   la que no depende del motor, y la que se puede verificar en una sola sesion.
4. **Fase 2 (inventario)**: la de uso diario.
5. **Fases 4 y 5**: las matrices, el trabajo de diseno mas pesado.
6. **Fases 6, 8**: pulido del diagrama y detalle.
7. **Fase 7**: crear integracion, con el OAuth como riesgo aparte.
8. **Fase 9**: nada que hacer.

## Como se verifica cada fase

Con el negocio **"sin intermediarios"** (miles de ordenes) en el telefono
fisico por ADB, **solo lectura**: comparar, listar y filtrar. Ninguna prueba
dispara `apply`, `sync` ni `create` contra ese negocio, porque el ambiente
local apunta al RDS de produccion.

Ademas, por fase: `flutter analyze` limpio, `flutter test` sin regresiones
nuevas (hoy hay 29 fallas previas de `parseError`), y una medicion de
`dumpsys meminfo` en las fases 2, 4 y 5, que son las que traen imagenes de
producto a una lista larga.
