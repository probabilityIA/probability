# Plan Mobile - App Flutter Probability (20 fases)

Worktree: `/home/cam/Desktop/probability-mobile` (branch `feat/mobile-app-flutter`)
App: `mobile/mobile_central` (Flutter 3.38.6, Dart 3.10.7)

## Objetivo

Completar la app Flutter replicando los modulos del front `front/central`, con
diseno minimalista base blanco y los colores globales de la plataforma.
Reusar los endpoints existentes; cuando haga falta un endpoint agregado para
reducir llamadas desde la app, crearlo con `mobile` en la ruta
(`/api/v1/mobile/...`) para que quede claro que es de uso exclusivo del cliente
Flutter.

## Identidad visual

Extraida del logo y del front:

| Token | Valor | Uso |
|---|---|---|
| primary | `#5E17EB` | violeta Probability (logo) |
| primaryDark | `#4C1D95` | estados presionados, headers |
| accent | `#14F5A0` | verde menta, exito / destacados |
| secondary | `#F472B6` | rosa, badges secundarios |
| surface | `#FFFFFF` | base de toda la app |
| background | `#F9FAFB` | gray-50, fondo de listas |
| border | `#E5E7EB` | gray-200 |
| textPrimary | `#111827` | gray-900 |
| textMuted | `#6B7280` | gray-500 |
| success `#10B981` / warning `#F59E0B` / error `#EF4444` / info `#3B82F6` | | estados |

Radios `8/12/16`, sombras suaves, sin elevaciones fuertes, tipografia Inter.

## Estado base (auditoria inicial)

- 223 archivos Dart, arquitectura hexagonal ya aplicada
  (`domain/ | app/ | infra/repository/ | ui/{providers,screens}`).
- `flutter analyze`: sin issues.
- Faltaba plataforma web: agregada con `flutter create --platforms=web .`.
- Tema actual: `colorSchemeSeed: Colors.deepPurple` (Material por defecto,
  no es la marca). Se reemplaza en Fase 1.
- Sin logo aplicado.
- Muchos modulos son solo lista sin detalle ni acciones.

## Las 20 fases

| # | Fase | Contenido |
|---|---|---|
| 1 | Design System | tokens de marca, tema claro base blanco, tipografia, componentes base (AppScaffold, AppCard, AppButton, EmptyState, ErrorState, LoadingState, StatusChip, MoneyText), logo en assets |
| 2 | Shell y navegacion | bottom nav + drawer por modulo, selector de negocio, guardas de sesion, splash con logo |
| 3 | Auth | login, recuperar clave, verificacion, perfil, cambio de contrasena |
| 4 | Dashboard | KPIs, accesos rapidos, resumen de ventas |
| 5 | Ordenes - listado | filtros, busqueda, paginacion infinita, chips de estado |
| 6 | Ordenes - detalle y acciones | items, cliente, pagos, cambiar estado, cancelar, facturar |
| 7 | Envios / guias | lista, detalle, tracking, descargar guia, cancelar |
| 8 | Cotizador y generacion de guia | cotizar, comparar tarifas, generar guia |
| 9 | Contra entrega (COD) | reporte, cortes, conciliacion |
| 10 | Productos | lista, detalle, variantes, imagenes, precios |
| 11 | Inventario | stock por bodega, movimientos, ajustes, transferencias |
| 12 | Bodegas | bodegas, ubicaciones, ocupacion |
| 13 | Clientes | lista, detalle, direcciones, historico de compras |
| 14 | Facturacion | facturas, notas credito, emision DIAN, reintentos |
| 15 | Billetera | saldo, movimientos, recarga, alertas |
| 16 | Rutas / ultima milla | rutas, paradas, conductores, vehiculos |
| 17 | Integraciones | mis integraciones, catalogo, conectar, sincronizar, logs |
| 18 | Tienda online | storefront, catalogo publico, configuracion del sitio |
| 19 | Notificaciones | configuracion, canales, plantillas, historial |
| 20 | IAM + suscripcion + tickets | usuarios, roles, permisos, plan, soporte |

## Entorno local de trabajo

- Mock server HTTP en Node (`mobile/mock_server/`) que responde el contrato del
  backend real para poder desarrollar sin levantar Go + Postgres.
- App web: `flutter run -d web-server --web-port 5100`
  con `--dart-define=API_BASE_URL=http://localhost:5199/api/v1`.
- Verificacion visual con Playwright sobre `http://localhost:5100`.

## Cola de tareas pendientes (pedidas fuera de fase)

1. ~~Icono de la app con el logo de Probability~~ **LISTO**: mipmaps de Android
   (48 a 192 px), iconos de iOS, iconos y favicon de web, todos generados desde
   `logo.png` sobre fondo blanco. Nombre visible "Probability" en Android e iOS,
   manifest web con el violeta de marca.
2. **APK de produccion en el escritorio** al cerrar la fase 20 (ver abajo).
3. **Bitacora**: documentar todo lo hecho hoy en `.claude/bitacora/` al terminar.
4. **Pendientes anotados** en el plan y en alertas si aplica.
5. **Commit + push** de la rama al terminar (autorizado por el usuario).
6. **Apagar el equipo** al final.

## Alcance tocado fuera de `mobile/`

El usuario pidio que solo se afectara la app. Se cumple con una excepcion que el
mismo usuario pidio antes: el endpoint exclusivo de la app.

Fuera de `mobile/` solo hay:
- `back/central/services/modules/mobile/**` -> **modulo nuevo**, no toca nada
  existente.
- `back/central/services/modules/bundle.go` -> **2 lineas agregadas** (el import
  y `mobile.New(router, database)`).

Ningun comportamiento del backend actual fue modificado. Todo vive en la rama
`feat/mobile-app-flutter`, `main` esta intacta.

## Entrega final (al cerrar la fase 20)

Compilar el APK apuntando a **produccion** y dejarlo en el escritorio:

```bash
cd mobile/mobile_central
flutter build apk --release \
  --dart-define=APP_ENV=production \
  --dart-define=API_BASE_URL=https://www.probabilityia.com.co/api/v1
cp build/app/outputs/flutter-apk/app-release.apk ~/Desktop/probability-<fecha>.apk
```

`Environment.apiBaseUrl` ya resuelve a `https://www.probabilityia.com.co/api/v1`
cuando `APP_ENV=production`, asi que el `--dart-define` de la URL es redundante
pero se deja explicito para que el APK no dependa del default.

## Estado de avance

| Fase | Estado | Notas |
|---|---|---|
| 1 | LISTA | tokens de marca, tema, tipografia Inter empaquetada, logo, componentes base, login rediseniado y verificado en navegador |
| 2 | LISTA | bottom nav de 5 destinos, drawer agrupado con logo, AppScaffold comun, pantalla "Mas", ModuleTabsScaffold compartido por los 8 modulos con pestanias, gate de negocio para super admin |
| 3 | LISTA | recuperacion de clave en 4 pasos (canales/OTP/nueva clave), perfil, cambio de contrasena, sesion persistida completa |
| 4 | LISTA | saludo, 4 KPIs accionables, accesos rapidos, ordenes por canal con logo, envios por estado y transportadora con logo, top productos, clientes y ciudades |
| 5 | LISTA | listado con tarjeta de orden (logo de canal, estados, guia, COD), busqueda con debounce, filtros por estado, scroll infinito |
| 6 | LISTA | detalle con header de canal, items tipados, totales con bloque COD, cliente, envio, trazabilidad, copiar numero y cancelar orden |
| 7 | LISTA | lista de guias con logo de transportadora, filtros por estado, busqueda y scroll infinito; detalle con desglose de costo, contra entrega, destino, paquete y cancelacion |
| 8 | LISTA | primer endpoint exclusivo de la app en Go: `GET /api/v1/mobile/orders/:id/full`; el detalle de orden ya muestra guia y factura |
| 9 | LISTA | cotizador con formula portada de rate-pricing.ts, 7 pruebas que la fijan, desglose por tarifa y aviso cuando la tarifa no soporta contra entrega |
| 10 | LISTA | lista con miniatura, precio, stock y filtros; detalle con precios, margen calculado, inventario, dimensiones y canales donde esta publicado con su logo |
| 11 | LISTA | existencias por bodega con barra de ocupacion y alerta de reposicion, historial de movimientos con entrada/salida, y hoja de ajuste que escribe contra el API |
| 12 | LISTA | lista de bodegas con estado, predeterminada y despacho; detalle con direccion, codigo DANE, contacto con aviso de datos faltantes y ubicaciones internas |
| 13 | LISTA | directorio con busqueda y scroll infinito; detalle con KPIs de compras, contacto copiable y las ordenes reales del cliente |
| 14 | LISTA | facturas con logo del facturador (Siigo, Factus, Alegra), estado, motivo de rechazo de la DIAN, detalle con items, CUFE, reintento y cancelacion |
| 15 | LISTA | saldo con alerta de saldo bajo, KPIs de recargado y gastado en guias, movimientos tipificados por concepto real y hoja de recarga |
| 16 | LISTA | rutas con progreso y paradas fallidas, detalle con linea de tiempo de paradas y acciones iniciar/completar, conductores con licencia y estado, vehiculos con tipo y capacidad |
| 17 | LISTA | catalogo completo con los 29 logos reales agrupado por categoria, y acciones sobre cada integracion conectada (probar, sincronizar, activar/desactivar) |
| 18 | LISTA | catalogo publico en grilla con destacados y agotados, y configuracion del sitio con interruptor por seccion que guarda contra el API |
| 19 | LISTA | eventos con su canal y logo, contador de activos, filtro por canal e interruptor por evento que guarda contra el API |
| 20 | LISTA | usuarios con avatar, rol y marca de super admin; roles con nivel y alcance; permisos agrupados por recurso |

## APIs exclusivas de la app (`/api/v1/mobile/...`)

Regla acordada: cuando la app necesite varias llamadas para pintar una sola
pantalla, se crea un endpoint agregado con `mobile` en la ruta. Hasta la fase 6
**no hizo falta ninguno**: cada pantalla resuelve con un solo GET de los que ya
existen.

### `GET /api/v1/mobile/orders/:id/full` (fase 8)

Modulo `back/central/services/modules/mobile/`, arquitectura hexagonal, montado
en `services/modules/bundle.go`. Devuelve en una sola respuesta:

```
{ order, items[], shipment|null, invoice|null }
```

Ahorra 3 llamadas por apertura de detalle: antes hubiera necesitado
`/orders/:id`, `/shipments?order_id=` y `/invoices?order_id=`.

Cumple las reglas del repo:
- repositorio propio, solo SELECT, sin importar repos de otros modulos
- modelos desde `migration/shared/models`, sin `.Table()`
- aislamiento multi-tenant: `business_id` sale del token; el super admin lo pasa
  por query param, y ademas se valida que la orden pertenezca a ese negocio
  antes de devolver nada (404 si no, para no filtrar existencia)
- sin comentarios en el codigo Go

Cliente Dart en `lib/services/modules/mobile/`, con su provider
`OrderFullProvider`.

## La formula de tarifas

`lib/shared/utils/rate_pricing.dart` es un port literal de
`front/central/src/shared/utils/rate-pricing.ts`. **Ninguna pantalla suma
tarifas a mano**, igual que en el front. `test/shared/utils/rate_pricing_test.dart`
fija el comportamiento con 7 casos, incluidos los dos errores que la regla marca
como clasicos: omitir el seguro minimo y sumar el margen COD a una tarifa que no
soporta contra entrega.

## Reglas de negocio respetadas en guias

`.claude/rules/guias-contra-entrega.md` fija la formula del costo. La pantalla de
detalle la muestra explicita y cuadra la identidad
`carrier_cost + applied_margin == total_cost` con los datos reales de la tabla
`shipments`. La comision de contra entrega se muestra aparte porque **no** esta
incluida en `cod_total` (esa es la convencion, no un bug).

Al cancelar una guia el dialogo avisa que **no se reembolsa el saldo debitado**,
que es lo que ocurre hoy segun `.claude/alerts/guias-duplicadas-doble-cobro.md`.

## Estado de cada modulo: produccion, beta o desarrollo

`lib/shared/navigation/app_modules.dart` marca cada modulo con `ModuleStage`:

| Estado | Se ve en el menu | Insignia | Ruta alcanzable |
|---|---|---|---|
| `prod` | si | no | si |
| `beta` | si | "BETA" | si |
| `development` | **no** | - | **no** |

Un modulo en `development` desaparece del drawer y de la pantilla "Mas", su
grupo se oculta si se queda vacio, y el `redirect` del router manda a
`/dashboard` cualquier ruta suya: no basta con esconder el boton, un enlace
directo o una sesion restaurada tambien tienen que rebotar.

**Ultima milla (`/delivery`) esta en `development`** y por eso no aparece. El
codigo del modulo (rutas, conductores, vehiculos) queda intacto: para
reactivarlo se cambia una sola linea a `ModuleStage.beta` o `prod`.

Al ocultarlo salieron dos cosas de paso:

- La pestania inferior "Envios" apuntaba a `/delivery` (ultima milla), no a las
  guias. Se **elimina** de la barra: las guias ya viven como pestania dentro de
  Ordenes, y tenerlas dos veces confundia. La barra queda en cuatro: Inicio,
  Ordenes, Inventario y Mas. Ademas pasa a elegir la pestania **mas especifica**
  en vez de la primera que coincide, para que `/orders/shipments` marque
  Ordenes.
- `home_screen.dart` era codigo muerto (no lo referenciaba nadie) y contenia los
  accesos rapidos a Rutas, Conductores y Vehiculos. Se elimina.

Pendiente si se quiere llevar mas lejos: hoy el estado es una constante de la
app. Si se necesita prender modulos por negocio o sin publicar version, el
estado tendria que venir del backend como feature flag.

## Paginado dinamico y memoria (regla del proyecto)

`.claude/rules/mobile-listados-memoria.md` fija el patron obligatorio: todo
listado con scroll infinito sobre `ListView.builder`, nunca paginado numerado, y
nunca acumulando sin tope.

Infraestructura en `lib/shared/pagination/`:

- `PagedCollection` - lista rala con ventana deslizante y expulsion LRU. Al
  expulsar **no corre los indices**: deja el hueco en `null`, asi que el scroll
  no salta y el hueco se vuelve a pedir al entrar al viewport.
- `PagedListController` - refresh, loadMore, relleno de huecos y mantenimiento
  de la ventana. Tope por defecto: 8 paginas de 20 = 160 items vivos.
- `PaginatedListView` - el patron completo en un widget.
- `PageSizes` - `list` (20) y `catalog` (100).

Memoria de imagenes, que es el consumidor real en gama baja: `cacheWidth` y
`cacheHeight` obligatorios en todo `Image.network` (un logo de 40 px decodificaba
el PNG completo, ~1 MB por logo), mas el tope global del `ImageCache` en
`main.dart` (40 MB / 120 imagenes) via `ImageMemory.applyLowEndBudget()`.

Los 20 listados quedaron migrados. El mock server genera 640 ordenes, 420
clientes, 380 productos y 400 movimientos justamente para que la ventana
deslizante se pueda ver funcionando (con 64 registros nunca se expulsaba nada).

Cobertura: 27 pruebas nuevas entre `PagedCollection`, `PagedListController`,
`PaginatedListView` (widget tests que hacen scroll de verdad) y
`order_provider_paging_test.dart` sobre el provider real.

## Pendiente nuevo: filtro de stock de productos

Las pastillas "Con stock / Stock bajo / Agotados / Inactivos" del listado de
productos filtraban en el cliente sobre la pagina cargada, y `GetProductsParams`
**no tiene parametro de stock ni de estado**. Con paginacion real eso solo
filtraba los 20 registros cargados, asi que se quitaron de la app. Para
reponerlas hay que agregar el filtro al endpoint `GET /products` del backend
(`stock_status` o equivalente) y volver a pasarlas como filtro de servidor.

## Pendiente nuevo: los tests de providers prueban una copia

23 de los 29 archivos `*_provider_test.dart` definen una clase
`Testable<X>Provider` **dentro del propio test** y prueban esa copia, no el
provider real. Los unicos 6 que prueban el provider de verdad son login,
notification_config, invoicing, inventory, my_integrations y business. Es
cobertura aparente: se puede romper un provider real sin que falle una sola
prueba. Los cuatro providers de listado grandes (orders, shipments, products,
customers) ya aceptan `useCases` inyectado para poder probarlos de verdad;
`order_provider_paging_test.dart` es el ejemplo del patron.

## Pendientes al cerrar las 20 fases

1. **29 pruebas en rojo preexistentes** (ver deuda abajo). No las introdujo este
   trabajo; conviene una fase de limpieza.
2. **Pantallas que quedaron en modo lectura**: crear y editar productos,
   clientes, bodegas, conductores, vehiculos y rutas. La app hoy consulta y
   ejecuta acciones puntuales (ajustar stock, cancelar guia, reintentar factura,
   activar integracion), pero los formularios de alta completos siguen en la
   plataforma web.
3. **Generar guia desde el cotizador**: la pantalla cotiza y deja elegir tarifa,
   pero el boton final solo confirma la seleccion. Falta encadenar
   `POST /shipments/generate` con la bodega de origen y la orden.
4. **Conectar integracion desde la app**: el catalogo lista y filtra, pero el
   alta de credenciales sigue siendo web (cada proveedor tiene su formulario).
5. **Tickets y suscripcion**: quedaron fuera del alcance efectivo de la fase 20,
   que se concentro en IAM. Los endpoints ya estan mockeados.
6. **CORS del bucket S3**: en Flutter web los logos cargan por el fallback a
   `<img>`. Habilitar CORS en `probability-media-assets` limpiaria los errores
   de consola. No afecta a Android ni iOS.
7. **Push notifications**: no se toco. La app no tiene FCM configurado.

## Deuda conocida

- `flutter test`: **29 pruebas en rojo, todas preexistentes** (no las introdujo
  este trabajo, verificado por A/B con el cambio revertido):
  - 9 en `test/services/auth/business/ui/providers/business_provider_test.dart`
  - 20 en `test/services/modules/...`, la mayoria por `parseError`, que devuelve
    "Ocurrio un error inesperado" donde la prueba espera el mensaje del
    proveedor.
  El resto (2098) pasa. Conviene arreglarlas en una fase de limpieza aparte.

## Logos de marca

Los logos de integraciones salen de `integration_type.image_url`, igual que en el
front. Es una llave de S3 (`integration-types/xxx.png`) que
`BrandAssets.mediaUrl` convierte en URL absoluta contra
`https://probability-media-assets.s3.us-east-1.amazonaws.com`.
Los de transportadoras salen del mapa `BrandAssets.carrierLogos`, portado de
`front/central/src/shared/utils/carrier-logos.ts`.

Widget unico: `BrandLogo(name:, imageUrl:)`. Si no hay imagen o falla la carga
cae a una insignia con iniciales y color derivado del nombre.

**Nota web:** el bucket de S3 no manda `Access-Control-Allow-Origin`, asi que en
Flutter **web** la carga por XHR falla por CORS. `BrandLogo` usa
`webHtmlElementStrategy: WebHtmlElementStrategy.fallback`, que cae a un `<img>`
del DOM y muestra la imagen igual. En Android/iOS no existe CORS y carga
directo, asi que esto no afecta a la app real. Si se quiere limpiar los errores
de consola en web, hay que habilitar CORS en el bucket.

## Como levantar el entorno

```bash
mobile/dev.sh start        # mock :5199 + flutter web :5100
mobile/dev.sh restart-app  # hot restart
mobile/dev.sh logs 40
mobile/dev.sh stop
```

Login de prueba: cualquier correo y contrasena (el mock siempre autentica).

## Convenciones de codigo Dart

- Sin comentarios en el codigo.
- Sin caracteres non-ASCII en los archivos `.dart`: los textos con tilde usan
  escapes `\u00XX` (mismo criterio que Go y TypeScript en este repo).
- Componentes compartidos en `lib/shared/widgets/ui/` (barrel `ui.dart`).
- Colores solo desde `AppColors`; radios, sombras y espaciados desde
  `app_tokens.dart`. Prohibido `Colors.deepPurple` u otro color suelto.
