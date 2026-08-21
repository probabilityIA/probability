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

## Estado de avance

| Fase | Estado | Notas |
|---|---|---|
| 1 | LISTA | tokens de marca, tema, tipografia Inter empaquetada, logo, componentes base, login rediseniado y verificado en navegador |
| 2 | LISTA | bottom nav de 5 destinos, drawer agrupado con logo, AppScaffold comun, pantalla "Mas", ModuleTabsScaffold compartido por los 8 modulos con pestanias, gate de negocio para super admin |
| 3 | LISTA | recuperacion de clave en 4 pasos (canales/OTP/nueva clave), perfil, cambio de contrasena, sesion persistida completa |
| 4 | LISTA | saludo, 4 KPIs accionables, accesos rapidos, ordenes por canal con logo, envios por estado y transportadora con logo, top productos, clientes y ciudades |
| 5 | siguiente | ordenes: listado con filtros, busqueda y scroll infinito |

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
