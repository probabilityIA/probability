# App mobile Flutter: 20 fases en una sesion

Fecha: 2026-08-20 (cierre 2026-08-21 madrugada)
Rama: `feat/mobile-app-flutter`
Worktree: `/home/cam/Desktop/probability-mobile`

## Que se pidio

Completar la app Flutter que ya existia en `mobile/mobile_central`, replicando
los modulos del front `front/central` con diseno minimalista base blanco y los
colores de marca. Reusar los endpoints existentes; si hacia falta un endpoint
agregado para la app, crearlo con `mobile` en la ruta. Mostrar los logos de las
integraciones igual que en la plataforma. Plan de 20 fases, una por vuelta.

## Estado inicial

223 archivos Dart con arquitectura hexagonal ya aplicada, `flutter analyze`
limpio, pero casi todo era listado sin detalle, con `colorSchemeSeed:
Colors.deepPurple` (Material por defecto, no la marca) y sin logo. No existia
plataforma web en el proyecto Flutter.

## Como se trabajo

- Worktree aparte para no tocar `main`.
- **Mock server propio** en `mobile/mock_server` (Node sin dependencias) para
  desarrollar sin levantar Go + Postgres + Redis + RabbitMQ. Devuelve 404 con
  el texto `mock sin ruta: METODO /path`, que sirvio de lista de trabajo.
- `mobile/dev.sh` levanta mock (:5199) y Flutter web (:5100).
- Cada pantalla se verifico en el navegador con Playwright antes de cerrar la
  fase. Flutter web no expone DOM, asi que se interactua despachando
  PointerEvent sinteticos sobre `flutter-view` y leyendo `flt-semantics`.

## Bugs reales encontrados en la app (no inventados por el mock)

1. **Billetera siempre en cero.** `wallet_repository` leia `response.data` sin
   desenvolver el campo `data` del sobre estandar, asi que el saldo llegaba en
   cero y el historial vacio.
2. **Notificaciones siempre vacias.** Mismo patron: el repositorio solo aceptaba
   una lista pelada y descartaba `{success, data}`. Se reviso el resto de repos:
   `inventory` ya lo manejaba bien, no habia mas casos.
3. **Foto de usuario nunca cargaba.** El backend guarda el avatar como llave de
   S3 (`avatars/xxx.jpg`) y `NetworkAvatar` la usaba como URL absoluta, asi que
   siempre caia a la inicial.
4. **Clientes sin total comprado.** El listado parseaba `CustomerInfo` en vez de
   `CustomerDetail`, descartando `order_count` y `total_spent` que el API si
   manda.
5. **Catalogo de integraciones incompleto.** Pedia la pagina por defecto y solo
   mostraba 10 de 29 integraciones.
6. **Paginacion que no acumulaba.** Los providers de ordenes, productos y
   clientes reemplazaban la lista en cada pagina en vez de concatenar, asi que
   el scroll infinito no servia.
7. **Items de orden sin tipar.** `orderItems` era `dynamic`; se creo la entidad
   `OrderLineItem`.
8. **JWT impreso en consola** en el login. Se elimino.
9. **Sesion incompleta.** Solo se persistia el usuario: al reiniciar se perdian
   los negocios y el flag de super admin.

## Contratos que se corrigieron contra la BD real

El mock inicial invento campos. Consultando el esquema se corrigio:

- `route` **no tiene columna `name`**; `route_stop` no guarda `order_number`,
  `item_count` ni `total_amount`. Las pantallas se rediseniaron con lo que si
  existe.
- Los conceptos de billetera salen de `transaction.type` + `transaction.concept`
  (USAGE/RECHARGE x GUIDE/RECHARGE/SUBSCRIPTION/REFUND/ADJUSTMENT/OTHER).
- Los logos de integracion son la llave S3 de `integration_types.image_url`,
  resuelta contra `https://probability-media-assets.s3.us-east-1.amazonaws.com`.
- `website_config`: los campos de contenido son objetos JSON, no cadenas.

## Reglas de negocio respetadas

- **Formula de tarifas**: `lib/shared/utils/rate_pricing.dart` es un port literal
  de `front/central/src/shared/utils/rate-pricing.ts`, con 7 pruebas que la
  fijan. Ninguna pantalla suma tarifas a mano. Cubre los dos errores que la
  regla marca como clasicos: omitir el seguro minimo y aplicar margen COD a una
  tarifa que no soporta contra entrega.
- **Costo de guia**: el detalle muestra `carrier_cost + applied_margin =
  total_cost` y deja la comision de contra entrega aparte, porque no va dentro
  de `cod_total`.
- **Cancelar guia**: el dialogo avisa que no se reembolsa el saldo debitado,
  que es el comportamiento real hoy (ver alerta guias-duplicadas-doble-cobro).
- **Recarga de billetera**: la hoja avisa que el cobro es real en la pasarela.
- **Bodega sin telefono o correo**: el detalle avisa que el carrier rechaza la
  guia, aunque la cotizacion funcione.

## Unico endpoint nuevo

`GET /api/v1/mobile/orders/:id/full` en `back/central/services/modules/mobile/`.
Devuelve orden + items + guia + factura en una sola respuesta; ahorra 3 llamadas
por apertura de detalle. Modulo nuevo con arquitectura hexagonal, repositorio
propio de solo lectura, `business_id` del token y validacion de que la orden
pertenezca al negocio antes de responder (404 si no, para no filtrar existencia).

Fuera de `mobile/` solo se toco eso mas 2 lineas en
`services/modules/bundle.go` (import y registro). Ningun comportamiento del
backend actual quedo modificado.

## Pendientes

Ver `mobile/PLAN-MOBILE.md`, seccion "Deuda conocida" y "Pendientes".
