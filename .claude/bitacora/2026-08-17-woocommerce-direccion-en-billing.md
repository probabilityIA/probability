# WooCommerce: las ordenes entraban sin direccion de envio

Desde el 12/08/2026 las ordenes de Moto Mello (business 52) entraban con
`shipping_street`, `shipping_city` y geozona vacios. WooCommerce si mandaba la
direccion, pero en el bloque `billing`, y el mapper solo leia `shipping`.

## Sintoma

En el detalle de la orden no aparecia la direccion de destino ni el mapa.
Las ordenes recientes de la lista (148236, 148231, 148228, 148215...) todas
igual.

En base de datos:

```sql
SELECT order_number, shipping_street, shipping_city, shipping_lat
FROM orders WHERE business_id = 52 AND platform = 'woocommerce'
ORDER BY created_at DESC LIMIT 5;
-- todas con street = '', city = '', lat = NULL
```

Corte exacto por dia (business 52):

| Dia | Ordenes | Con direccion |
|-----|---------|---------------|
| 10 ago | 12 | 12 |
| 11 ago | 10 | 10 |
| **12 ago** | 16 | **8** |
| 13 ago en adelante | 44 | **0** |

Total afectado: **53 de 1.695** ordenes de WooCommerce del business 52.

## Diagnostico

El JSON crudo guardado en `order_channel_metadata` fue la prueba definitiva.
Orden 148236:

```
shipping.address_1 = ""              <- vacio
shipping.city      = ""
billing.address_1  = "Carrera 6 # 17 90"
billing.city       = "FUNZA"
billing.state      = "CO-CUN"
```

Y en una orden del 12/08 que si funciono (147763), `shipping.address_1` venia
lleno. O sea: el proveedor cambio de comportamiento a mitad del 12 de agosto,
casi seguro por un cambio en el checkout de la tienda.

### Hipotesis descartadas

- **"WooCommerce dejo de mandar la direccion"**: falsa. La manda completa, solo
  que en `billing`. Es el comportamiento normal de Woo cuando el comprador no
  pide una direccion de envio distinta a la de facturacion: el bloque `shipping`
  llega vacio y el envio va a la de billing.
- **"Las cotizaciones y guias salieron con destino malo"**: falsa. Las
  cotizaciones las genera el plugin de checkout con su propio payload, que si
  trae la direccion y el `daneCode`:

  ```json
  "destination": { "address": "Carrera 6 # 17 90 Urbanizacion la campina casa 4",
                   "suburb": "FUNZA", "daneCode": "25286000" }
  ```

  De las 53 afectadas, 47 tenian cotizacion con direccion completa y **ninguna
  tenia guia generada**, asi que no hubo envios con destino equivocado.

## Causa raiz

`woocommerce/internal/app/usecases/mapper/order_mapper.go` construia la
direccion de envio leyendo unicamente `order.Shipping`, sin caer a
`order.Billing` cuando el bloque venia vacio.

## Correccion

**Codigo** (commit `10e491ba`, desplegado): el mapper cae a `billing` cuando la
direccion de envio llega sin calle **ni** ciudad. La direccion de facturacion se
sigue guardando aparte. Test: `TestMapWooOrderToProbability_ShippingVacioUsaBilling`.

**Data** (produccion, 53 ordenes, en una transaccion):

1. Direccion desde `billing` del JSON crudo mas reciente (53 filas).
2. Geozona desde el `daneCode` de la cotizacion del checkout, cruzando
   `LEFT(daneCode,5) = geozones.code` con `type = 'city'` (47 filas).
3. Geozona por nombre de ciudad para las que no tenian cotizacion (5 filas).
4. Orden 148029 resuelta a mano: "RIONEGRO" existe en Antioquia (id 119) y en
   Santander (id 900); se uso `shipping_state = 'CO-ANT'` para desambiguar.

Respaldo de los valores previos en la tabla `backup_woo_addr_20260818`.

## Verificacion

```sql
SELECT count(*) AS total,
       count(*) FILTER (WHERE COALESCE(shipping_street,'') <> '') AS con_direccion,
       count(*) FILTER (WHERE destination_geozone_id IS NOT NULL) AS con_geozona
FROM orders WHERE id IN (SELECT id FROM backup_woo_addr_20260818);
-- 53 | 53 | 53
```

Ejemplo: 148236 quedo con "Carrera 6 # 17 90 Urbanizacion la campina casa 4",
FUNZA, geozona FUNZA.

## Pendientes

- Las 53 quedaron **sin `shipping_lat`/`shipping_lng`**: el geocodificador solo
  corre al crear la orden (`usecasecreateorder/geocode_order.go`), no al
  actualizar. El mapa muestra la zona de la ciudad, no el pin exacto. Si se
  quiere el pin hay que geocodificar esas direcciones aparte.
- Preguntar a Moto Mello que cambiaron en el checkout el 12/08. Con el fix ya no
  afecta, pero explica el corte exacto.
- Vale revisar si otros negocios con WooCommerce tienen el mismo patron: hoy
  business 26 tiene 152 de 214 ordenes con ciudad pero solo 52 con calle.
