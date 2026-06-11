# Sesion 2026-05-13: geozonas, probabilidad por carrier, facturacion masiva

## Resumen ejecutivo

Sesion larga enfocada en optimizar el modulo de geozonas (probabilidad de
entrega por transportadora) y la creacion masiva de facturas. Todo en main.

## Cambios en main (orden cronologico de commits)

| Commit  | Resumen |
|---------|---------|
| 33150aec | Tabla agregada geozone_carrier_stats + carrier_key normalizado + indices. Endpoint /probability/by-carrier baja de 18s a sub-ms. Indice compuesto en orders: listado baja de 600ms a 10ms. Backfill 1455 ordenes via match por nombre. |
| ecf32b9b | Baseline cascade (zona > global > baseline carrier) + marker P violeta sobre geozona + mapa en paso 2 del modal Generar Guia. |
| e1f64081 | Refactor BulkCreateInvoiceModal: paginado server-side, page_size hasta 200, filtros dinamicos (fecha, order_number, cliente), sort, seleccion persistente entre paginas, Set<string> para toggle O(1). |
| 2a568737 | Fix matching tolerante a sufijos D.C. en ResolveOrderGeozone + autocompletado obligatorio en OrderForm + origen mostrado en mapa del modal de Guia. Backfill 603 ordenes adicionales. |
| ce0a0bfb | Resolucion hasta nivel barrio: columna shipping_neighborhood + segundo paso SQL resolveOrderBarrio que matchea barrio dentro del ancestor chain de la city. 176 ordenes con barrio extraido de shipping_street. |
| 61f89df1 | Badge "estimado" en cards de efectividad (gris/desaturado) + tooltip explicativo. |
| 4f04a685 | Mapa dinamico por carrier en paso 2: click YA NO avanza, polygon se colorea rojo->verde segun rate, badge flotante con %, boton Continuar explicito. |
| 664da53d | Banner Destino prominente (violeta) + Origen como sub-linea tenue + cascade single-carrier + remueve Efectividad de recoleccion de las cards. |
| 4de29201 | Cascada real por niveles geograficos (barrio->...->country->__global__) sin baselines inventados. CarrierEffectivenessRates ahora muestra el origen del dato: "barrio Compartir - 3 envios", "tasa nacional - 1240 envios", etc. |
| 62d73266 | Fix UNNEST array placeholders no funcionan con GORM Raw. Reescritura del CTE con VALUES interpolado seguro. Esto rompia TODOS los endpoints de probabilidad (500 error -> sin datos en cards). |
| 3a8027f8 | Fix HTTP 404 en pagina /delivery/geozones. nginx prod tiene location /api/ -> backend, por lo que las BFFs de Next.js en /api/* eran inalcanzables. Reemplazado por server action getGeozonesForDisplayAction. |
| 8f2f1e84 | Remove DeliveryProbabilityByCarrier del modal de detalle de orden (no aplicaba en ese contexto). |

## Estado actual de geozonas

### Cascade de probabilidad (sin inventos)
Solo data real, bajando resolucion geografica:
```
barrio -> neighborhood -> admin_district -> locality -> city -> state -> country -> __global__ (tasa nacional)
```
Si total < 5 en el nivel matcheado, se marca is_estimated=true con
estimate_source="zone_low_sample". El frontend lo renderiza con barra
gris/desaturada.

Solo cuando el carrier no tiene ningun envio en BD, se cae al baseline
hardcoded (85% generico).

### Tabla agregada
`geozone_carrier_stats` con primary key (geozone_level, geozone_id,
carrier_key). Refresh completo via cron interno (cada 6h al arrancar el
backend), refresh incremental por evento RabbitMQ (orders.events fanout
-> queue geozones.probability.invalidate). Cache Redis fina con TTL 10min,
invalidacion dirigida por orden.

### Resolucion automatica de geozona en ordenes
- ResolveOrderGeozone se llama SIEMPRE al crear/actualizar (con o sin
  lat/lng) - antes solo corria si habia lat/lng.
- Normalizacion tolerante a sufijos D.C. en city y state.
- Fallback nivel barrio cuando shipping_neighborhood (o 3er segmento de
  shipping_street) matchea contra geozone tipo barrio cuya ancestor chain
  pasa por la city ya resuelta.

## Estado actual de facturacion masiva

### Backend
Endpoint `GET /api/v1/invoicing/invoices/invoiceable-orders` ahora soporta:
- `page`, `page_size` (max 200)
- `business_id` (super admin)
- `start_date`, `end_date` (created_at, ISO o YYYY-MM-DD)
- `order_number`, `customer_name`, `customer_email` (ILIKE)
- `payment_status_id`, `fulfillment_status_id`
- `sort_by` (whitelist: created_at, occurred_at, order_number, total_amount, customer_name)
- `sort_order` (asc | desc, default desc)

`InvoiceableOrdersFilter` como DTO tipado, sanitizado en domain antes de
llegar al repo.

### Frontend
`BulkCreateInvoiceModal` rehecho con paginacion server-side, selector
page_size, filtros colapsables, sort dropdown con 7 presets, busqueda como
filtro server-side con debounce 300ms, "Seleccionar pagina" vs "Seleccionar
todas las N coincidencias" (lleva un Set<string> persistente entre paginas
con toggle O(1)).

## Pendientes

### Alta prioridad
- [ ] **Persistir shipping_neighborhood al crear orden manual desde el
  OrderForm.** Hoy el form sigue concatenando barrio en shipping_street
  con `parts.join(' | ')`. El backfill llena las viejas pero las nuevas
  creadas por UI siguen con barrio "escondido" en street. Hay que:
  1. Modificar el DTO en back/central (orders/dtos): agregar
     `ShippingNeighborhood string`.
  2. En el handler create/update, pasar el valor al model y
     `populateOrderFields`.
  3. En el OrderForm.tsx, NO concatenar el barrio en shipping_street;
     enviarlo aparte.
  4. Eliminar el split_part del SQL resolveOrderBarrio (queda como
     fallback solo para ordenes viejas).

### Media prioridad
- [ ] Verificar que el cache de Redis no este invalidando cambios. Cuando
  el usuario probaba veia "sin datos" pero la BD estaba bien; podria haber
  un caso de stale cache para `geozones:probability:{business_id}:{order_id}`
  cuando se cambian datos en backend sin gatillar el event RabbitMQ.
- [ ] Asegurarse de que el badge "is_estimated" deje de aparecer
  automaticamente cuando la zona acumule >= 5 envios reales. Validar en
  ambiente real cuando se acumulen mas shipments.
- [ ] Implementar `is_cod` real con columna boolean indexada en orders
  (hoy el filtro hace ILIKE sobre payment_details JSONB, costoso). Bajaria
  los outliers de 900ms en GET /api/v1/orders.

### Baja prioridad / nice-to-have
- [ ] Indices GIN con pg_trgm sobre customer_name/email/phone para que las
  busquedas ILIKE no degeneren a seq scan.
- [ ] Dead letter exchange en RabbitMQ para mensajes que fallan N veces.
- [ ] Tabla `carrier_baselines` editable desde admin (en vez de hardcoded
  en baseline.go).
- [ ] Limpiar las 23 ordenes que aun tienen `geozone_state_id IS NULL`
  (probablemente direcciones internacionales o muy raras).
- [ ] Reescribir filtros de BulkCreateInvoiceModal usando el componente
  shared `DynamicFilters` (ya existe en shared/ui/dynamic-filters.tsx).

## Lecciones operacionales

1. **GORM no expande Go slices como arrays de Postgres.** Si pasas
   `[]string{"a","b"}` a una query con `?::text[]`, GORM emite
   `'a','b'::text[]` (lista de valores) en vez de `ARRAY['a','b']::text[]`.
   Para arrays: usar `pq.Array(...)` del driver `lib/pq`, O construir CTEs
   con VALUES interpolando con whitelist.

2. **nginx prod tiene `location /api/ -> backend`.** Eso significa que
   TODAS las BFFs de Next.js bajo `/app/api/...` son inalcanzables (caen
   al backend Go que devuelve 404). Si necesitas un BFF, ponelo bajo otra
   ruta (ej. `/app/bff/...`) o cambia la regla de nginx a
   `location /api/v1/ -> backend`.

3. **El bug de "sin datos" tenia 3 capas.** Primero pense que era cache de
   browser (no). Despues, que era backend cascade no aplicada al endpoint
   single-carrier (parcial). Finalmente, era el bug de UNNEST que rompia
   TODOS los queries de probabilidad con 500. Lograba lo correcto era
   pedir logs al server antes de seguir.

4. **El cache de prompt de 5min de Anthropic importa para sesiones largas.**
   Esta sesion lo aprovecho bien gracias a no tener gaps largos.

## Numeros finales del dia

- Endpoints optimizados: 2 (probability/by-carrier, orders listing)
- Speed up acumulado: probability/by-carrier de 18000ms a 0.3ms (~60000x),
  orders de 600ms a 10ms (60x)
- Ordenes con geozona resuelta retroactivamente: 1455 + 603 + 8 = 2066
- Tablas nuevas: 1 (geozone_carrier_stats)
- Columnas nuevas: 2 (carrier_key generated en shipments y geozone_monthly_stats, shipping_neighborhood en orders)
- Indices nuevos: 13
- Commits a main: 13
