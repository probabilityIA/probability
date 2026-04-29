# shipping_margins

Modulo que gestiona el margen comercial que cada negocio cobra al cliente
sobre el precio real de cada transportadora.

## Que hace

1. **Configuracion por business y carrier** — guarda en `shipping_margin`
   un valor `margin_amount` (suma al flete) y `insurance_margin` (suma al
   seguro) por cada transportadora soportada (servientrega, interrapidisimo,
   coordinadora, envia, tcc, deprisa, 99minutos, mipaquete, enviame).
2. **Aplicacion automatica al cotizar** — el consumer de respuestas de
   transporte (`shipments/.../response_consumer.go`) lee el margen del
   cache (Redis con fallback a DB) y lo suma a `flete` y `minimumInsurance`
   antes de devolver las opciones al frontend. El cliente solo ve el
   precio inflado.
3. **Persistencia por guia** — al crear cada `shipment`, el use case
   `applyCarrierCost` consulta el margen vigente y guarda en la fila:
   - `total_cost`: lo cobrado al cliente.
   - `carrier_cost`: el costo real del carrier (`total_cost - margen`).
   - `applied_margin`: snapshot del margen al momento de la guia.
4. **Auto-seed al integrar transporte** — cuando un business activa por
   primera vez una integracion de la categoria `shipping` (envioclick,
   enviame, mipaquete, tu), un observer dispara `EnsureDefaultsForBusiness`
   que inserta los 9 carriers con `margin_amount = 0`. El admin solo
   ajusta los valores que quiera.
5. **Reporte de ganancias** — endpoint
   `GET /api/v1/shipping-margins/profit-report?business_id&from&to&carrier`
   agrupa por carrier con totales de cobrado, costo real y ganancia.
   Renderizado en la pestana "Reporte de ganancias" del modulo.

## Acceso

Todos los endpoints validan `requireSuperAdmin`. El usuario business
nunca ve el modulo: el sidebar y la pagina `/shipping-margins` se
ocultan, y el backend rechaza con 403 cualquier llamada directa.

## Wiring

```
                    +--> shipping_margins (config + reporte)
                    |
PublicIntegration --+--> observer (auto-seed por business)
                    |
                    +--> shipments.response_consumer (suma margen al quote)
                    |
                    +--> shipments.usecaseshipment (persiste carrier_cost)
                    |
                    +--> wallet (debita total_cost por guia, sin cambios)
```

## Tabla

```sql
shipping_margin (id, business_id, carrier_code, carrier_name,
                 margin_amount, insurance_margin, is_active,
                 created_at, updated_at, deleted_at)
```

UNIQUE por `(business_id, carrier_code)` a nivel logico (validado por
`ExistsByCarrier`).

## Eliminacion

No existe. Las filas son configuracion de sistema; solo se editan o se
desactivan con `is_active = false`. El endpoint y boton DELETE fueron
removidos a proposito.
