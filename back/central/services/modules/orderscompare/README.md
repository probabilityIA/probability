# Comparativo de ordenes (orderscompare)

Cruza las ordenes que existen en el canal de venta contra las que existen en
Probability, y permite crear en Probability las que faltan pasando por el flujo
normal de creacion de orden.

Es el equivalente, para ordenes, de lo que `shared/inventorycompare` es para
inventario. Vive en el hub "Mis integraciones" del front, pestana
**Comparar ordenes**.

## Que resuelve

Un canal puede tener ordenes que nunca llegaron a Probability: el webhook fallo,
la integracion se conecto despues de que la tienda ya vendia, o el consumidor
estuvo caido. Antes la unica salida era correr un `SyncOrders` completo, sin
saber que iba a entrar. Ahora se ve la lista, se elige y se importa.

## Endpoints

Ambos requieren JWT. `business_id` sale del token; solo el super admin
(`business_id = 0` en el token) lo manda por query param.

```
GET  /api/v1/orders-compare?integration_id=259&from=2026-08-01&to=2026-08-22&page=1&page_size=50&only_diff=true&q=1002
POST /api/v1/orders-compare/apply   { "integration_id": 259, "external_ids": ["1002","1003"] }
```

`GET` responde `{ success, data: { rows, totals, page, page_size, total, total_pages, checked_at, channel } }`.

Cada fila trae `action`:

| action | significado |
|---|---|
| `create` | esta en el canal y no en Probability: es la que se puede crear |
| `in_sync` | esta en las dos. Si ademas trae `status_mismatch` o `total_mismatch`, los datos difieren |
| `only_in_probability` | esta en Probability y el canal no la devolvio en el rango consultado |

Y dos campos de inventario: `moves_inventory` y `inventory_note`.

`POST /apply` responde `{ queued, skipped, failed, without_inventory, note }`.
`queued` son las publicadas a la cola: la orden se crea de forma asincrona, no
en la respuesta del endpoint. `skipped` son las que ya existian (el apply es
idempotente y ademas `orders` tiene indice unico `(integration_id, external_id)`).

## Como se crea la orden

**No** se escribe la tabla `orders` desde aca. El modulo le pide al canal que
publique la orden a la cola canonica `probability.orders.canonical`, exactamente
igual que un webhook o un `SyncOrders`. De ahi la toma el modulo `orders`
(`MapAndSaveOrder`), que valida cliente, productos, estados, geocodifica, guarda
y publica `order.created`. O sea: **el flujo es el normal, sin atajos.**

```
front -> POST /orders-compare/apply
      -> orderscompare (valida negocio y duplicados)
      -> canal.ImportChannelOrders (lee la orden del canal y publica el DTO canonico)
      -> cola probability.orders.canonical
      -> modulo orders: MapAndSaveOrder -> order.created
      -> modulo inventory: reserva stock
```

## Ordenes que NO mueven inventario

Una orden que en el canal ya esta entregada, despachada, cancelada o devuelta no
puede reservar stock hoy: esa mercancia ya salio (o nunca salio) y el stock
actual de Probability ya refleja la realidad. Reservarla descuadraria el
inventario.

Por eso el canal marca el DTO con `skip_inventory: true` antes de publicarlo. La
marca viaja asi:

```
canonical.ProbabilityOrderDTO.SkipInventory
  -> dtos.ProbabilityOrderDTO.SkipInventory   (modulo orders)
  -> event.Metadata["skip_inventory"]         (evento order.created)
  -> inventory/order_consumer.go              (ignora el evento y lo registra en el log)
```

La orden se crea completa (cliente, items, totales, estados); lo unico que no
ocurre es el movimiento de stock. En la UI esas filas salen marcadas con el
motivo, y el `apply` devuelve `without_inventory` + `note`.

La regla vive en un solo lugar: `shared/orderscompare/policy.go`
(`SkipsInventoryFor`). Si aparece un estado nuevo, se agrega ahi, no en el canal.

**El estado de la orden no siempre alcanza.** Hay canales donde la entrega no se
refleja en el estado de la orden: MercadoLibre deja `order.status = paid` aunque
el envio ya se entrego hace semanas, porque el dato vive en los tags y en el
shipment. Por eso `ChannelOrder` tiene `FulfillmentStatus`: el canal lo llena con
lo que sepa de la entrega y `SkipsInventoryFor` mira los dos. Un canal cuyo
estado de orden ya distingue la entrega (WooCommerce con `completed`) lo deja
vacio.

## Cuidado con los identificadores compuestos

El cruce se hace por `external_id`, asi que `ListChannelOrders` **tiene que
devolver el mismo identificador con el que la orden queda guardada**, no el que
al canal le resulte mas natural.

El caso que ya mordio: en MercadoLibre un carrito llega como una orden por
producto y `consolidatePack` las guarda como UNA orden de Probability con
`external_id = pack_id`. Devolver el id de la orden hacia que la misma compra
apareciera dos veces (una como "solo en Probability", otra como "falta aqui") y
que crearla generara un duplicado sin los items de las hermanas. Hoy
`ListChannelOrders` devuelve el `pack_id` cuando el pack trae mas de una orden, e
`ImportChannelOrders` acepta tanto un id de orden como uno de pack y consolida
antes de publicar.

## Como agregar un canal nuevo

El modulo no conoce ningun canal: habla con `IIntegrationContract` a traves del
registro de `integrations/core`. Para sumar un canal hay que implementar tres
metodos.

1. **Caso de uso del canal** (`internal/app/usecases/orders_compare.go`):

```go
func (uc *miCanalUseCase) ListChannelOrders(ctx context.Context, integrationID string, filters orderscompare.ChannelFilters) ([]orderscompare.ChannelOrder, error)
func (uc *miCanalUseCase) ImportChannelOrders(ctx context.Context, integrationID string, externalIDs []string) (orderscompare.ImportResult, error)
```

`ListChannelOrders` lee del API del canal y mapea al tipo comun. `ExternalID`
tiene que ser **el mismo valor** que el canal pone en `ExternalID` del DTO
canonico, o el cruce no encuentra nada y todo saldra como "por crear".

`ImportChannelOrders` lee cada orden completa, la mapea con el mapper que ya usa
`SyncOrders`, le pone `SkipInventory` segun `orderscompare.SkipsInventory(dto.Status)`
y la publica con el publisher de siempre.

2. **Interfaz del caso de uso**: agregar las dos firmas a `IMiCanalUseCase`.

3. **Adaptador core** (`internal/infra/secondary/core/orders_compare.go`):

```go
func (c *MiCanalCore) SupportsOrdersCompare() bool { return true }
func (c *MiCanalCore) ListChannelOrders(...)   { return c.useCase.ListChannelOrders(...) }
func (c *MiCanalCore) ImportChannelOrders(...) { return c.useCase.ImportChannelOrders(...) }
```

No hay que tocar este modulo ni el front: el canal aparece solo. Los canales que
no implementan nada heredan `BaseIntegration`, que devuelve `SupportsOrdersCompare() = false`
y el front no les ofrece la pestana.

Implementados hoy: Shopify (1), MercadoLibre (3), WooCommerce (4), VTEX (16),
Tiendanube (17), Jumpseller (33).

## Limites

- `limit` de lectura al canal: 200 por defecto, 1000 maximo.
- `page_size`: 50 por defecto, 200 maximo.
- `apply`: maximo 200 ordenes por lote.
- El comparativo consulta el canal **en vivo** en cada llamada; no hay snapshot
  en Redis como en inventario. Rangos amplios en canales lentos (VTEX pide el
  detalle orden por orden) tardan.

## Aislamiento multi-tenant

`resolveIntegration` valida que la integracion pertenezca al `business_id`
resuelto antes de leer nada. Un `integration_id` de otro negocio se rechaza,
tanto en el compare como en el apply.

## Pruebas

```bash
go test ./shared/orderscompare/... ./services/modules/orderscompare/...
```

`internal/app/usecase_test.go` usa dobles del repositorio y del registro de
canales: no necesita base de datos ni credenciales del canal.

### E2E en local contra el mock de Tiendanube

El simulador de `back/testing` sirve ordenes en `/v1/:store_id/orders` (puerto
9102). Trae 6 ordenes sembradas: 3 que mueven inventario (pending, paid, paid) y
3 que no (completed, cancelled, refunded).

```bash
./scripts/dev-db-switch.sh local
cd back/testing && go run cmd/main.go          # mock en :9102
# apuntar el canal al mock (solo en la base LOCAL):
#   update integration_types set base_url_test='http://localhost:9102/v1' where id=17;
#   update integrations set is_testing=true where id=<integracion tiendanube>;
cd back/central && go run cmd/main.go
```

Luego, en el hub "Mis integraciones" -> **Comparar ordenes**. Verificado el
2026-08-22: las 6 ordenes salen como "falta en Probability", las 3 terminales
marcadas "no mueve stock", y al crearlas el log del consumidor de inventario
escribe `Order marked as historical import, inventory movement skipped` solo para
esas 3. Volver a poner `is_testing=false` al terminar.
