# Alerta / Diseno: Sync de inventario Siigo -> Probability -> WooCommerce

Fecha: 2026-07-02
Estado: FASE 1 IMPLEMENTADA (compila + tests OK). FASES 2/3 PENDIENTES.

## ACTUALIZACION 2026-07-02: Siigo SI tiene webhook de stock (supera el enfoque pull)

Verificado por el usuario contra la API real (cuenta VELOCITYFULLCOMMERCESAS, Partner-Id: Velocity):
Siigo SI expone el topic `public.siigoapi.products.stock.update` que se dispara al cambiar el stock
de un producto. Esto es el evento que faltaba: el momento del REABASTECIMIENTO en Siigo.
=> El diseno pasa de "pull/polling" a EVENT-DRIVEN. La conclusion previa "no depender de webhook
de Siigo" queda SUPERADA.

Topics validos (verificados): products.stock.update (el que necesitamos), products.create, products.update.
Un topic por suscripcion. Auth: header Authorization=<token> (sin Bearer) + Partner-Id: Velocity.
Crear: POST /v1/webhooks {application_id, url, topic} -> 201 con id. Listar: GET. Borrar: DELETE /v1/webhooks/<id>.
NOTA cuenta prestada: en pruebas directas NO usar una URL que contenga "probability"; usar URL neutra de prueba.

Flujo revisado (recarga en Siigo -> Woo):
  recarga en Siigo -> webhook products.stock.update -> endpoint publico Probability
    -> relee stock absoluto del producto en Siigo (GET /v1/products) -> SyncProviderStock (existente)
      -> pushEcommerceStock -> cola -> consumer -> PUT WooCommerce (Fase 1 ya hecha).

Falta construir (Fase 2 replanteada):
- Endpoint receptor POST en modulo siigo (validado por company_key/firma del payload).
- Handler: resolver producto del payload -> resync de ese solo producto via SyncProviderStock.
- Reconciliacion completa nocturna SOLO como backstop (no mecanismo principal).
- PENDIENTE DATO: confirmar el shape del payload que envia Siigo (trae cantidad o solo id?).
  Plan: suscripcion temporal a webhook.site (URL neutra, sin "probability") + provocar cambio de stock.

## Progreso 2026-07-02 (Fase 1: push de stock Probability -> WooCommerce)

IMPLEMENTADO Y COMPILANDO (go build ./... OK, go test modulos OK):
- Cliente Woo UpdateProductStock: PUT /wp-json/wc/v3/products/{id} con {manage_stock:true, stock_quantity:N}.
  Archivo: services/integrations/ecommerce/woocommerce/internal/infra/secondary/client/update_product_stock.go
- IWooCommerceClient.UpdateProductStock (port) + WooClientMock actualizado.
- Usecase UpdateInventory: carga store_url del config + descifra consumer_key/secret + llama al cliente.
  Archivo: .../woocommerce/internal/app/usecases/update_inventory.go
- WooCommerceCore.UpdateInventory sobrescrito (antes ErrNotSupported).
- Consumer nuevo: escucha cola QueueWooInventoryStockPush y llama uc.UpdateInventory.
  Archivo: .../woocommerce/internal/infra/primary/queue/inventory_push_consumer.go
  Registrado en woocommerce/bundle.go (Start si rabbitMQ != nil).
- Publisher lado inventory: ISyncPublisher.PublishEcommerceStockPush + EcommerceStockPushMessage.
  Archivos: modules/inventory/.../domain/ports/ports.go, .../infra/secondary/queue/sync_publisher.go
- Enganche: useCase.pushEcommerceStock llamado dentro de updateProductTotalStock (adjust_stock.go).
  Se dispara tras CADA cambio de stock que pase por updateProductTotalStock (incluye el sync de Siigo
  via applyProviderStockItem y los ajustes manuales via AdjustStock).
- Cola nueva: shared/rabbitmq/queues.go -> QueueWooInventoryStockPush = "inventory.woocommerce.stock_push".

DECISIONES DE IMPLEMENTACION tomadas en Fase 1:
- Agregacion: pushEcommerceStock suma TODAS las bodegas del producto (total de updateProductTotalStock),
  NO por bodega. Coincide con "todas las bodegas activas del business -> un numero".
- Clamp: quantity = max(0, total) SOLO hacia Woo; products.stock_quantity guarda el total real.
- Filtro: solo integraciones con integration_type_code == "woocommerce" (id 4, code confirmado en DB)
  y con external_product_id no vacio. Los productos no fisicos ya se saltan aguas arriba
  (GetProductBySKU exige track_inventory).
- Matching a Woo: se usa external_product_id de product_business_integrations (ya existente),
  NO resolucion por SKU (GET products?sku=). Eso queda para Fase 2 si hace falta.
- Cola dedicada + Publish directo (patron Siigo), consumer secuencial: da serializacion de facto.
  Coalescing y particion por-SKU NO implementados aun (Fase 2).

VERIFICADO E2E 2026-07-02 (tramo cliente -> WooCommerce real, wordpress local :8088):
- client.UpdateProductStock ejercitado con el codigo real contra la tienda local.
  Producto 11: stock 40 -> 7 OK. Producto 22: manage_stock False->True y stock None->7 OK.
  Confirma que el payload {manage_stock:true, stock_quantity:N} activa y fija el stock.
- NO verificado aun end-to-end la cadena completa (publisher->cola->consumer) porque requiere
  backend corriendo y su DB apunta a RDS prod; el tramo interno queda cubierto por build + unit tests.

PENDIENTE INMEDIATO ANTES DE PRODUCCION:
- Verificar la cadena completa Siigo->Prob->cola->consumer->Woo con backend en un entorno seguro (no prod).
- Poblar product_business_integrations (mapeo producto<->external_product_id de Woo); sin eso no empuja.
- Confirmar SKU alineado (code Siigo == sku Probability == sku/producto Woo).

## (Diseno original abajo)

Estado historico: DISENO ACORDADO, IMPLEMENTACION PENDIENTE (no iniciada)
Modulos afectados: integrations/ecommerce/woocommerce, integrations/invoicing/siigo,
modules/inventory, modules/products, front (config de integraciones)

## Contexto

Se quiere cerrar el ciclo de inventario con Siigo como unica fuente de verdad:

  Siigo (inventario real) --> Probability (consume por pull) --> WooCommerce (refleja stock)

Hoy existe SOLO el tramo Siigo -> Probability (sync de stock por pull, commit 6efe6e65).
NO existe el tramo Probability -> WooCommerce: el cliente Woo solo hace pull de ordenes
(GetOrders, CreateWebhook). El slot generico UpdateInventory devuelve ErrNotSupported y
no tiene callers. Es trabajo nuevo.

WooCommerce core NO maneja multi-bodega: un solo stock_quantity por producto/variacion.
Por eso Probability debe AGREGAR (sumar) las bodegas activas del business a un solo numero
por SKU antes de empujar. No se necesita plugin multi-inventory.

Siigo NO sirve como emisor de eventos: GET /v1/webhooks de la cuenta real devolvio []
(sin suscripciones) y la doc oficial solo confirma el topic products.create (creacion),
sin evento de decremento de stock ni de factura. => No depender de webhook de Siigo.
El disparo del pull debe ser evento PROPIO de Probability (post-facturacion) + reconciliacion.

## Datos verificados de la cuenta Siigo (solo lectura, GET)

- Bodegas (GET /v1/warehouses): id 10 "PRUEBAS", id 90 "GENERAL" (ambas activas).
- Productos (GET /v1/products): 123 en total.
- Estructura de stock por producto: available_quantity (total) + warehouses[] {id,name,quantity}.
- OJO: la MAYORIA del stock vive en la bodega fantasma id -1 "Sin asignar", NO en 10/90.
- Aparecen cantidades NEGATIVAS y DECIMALES (-2315, 4574.01), muchas de items
  contables/servicio (COMISIONVENTAS, ALQUILERSOFTWARE, ALMACENAMIENTO), no productos fisicos.
- Matching de producto por code/SKU (ya usado por el sync actual).

## Decisiones tomadas (cerradas)

1. Bodega "Sin asignar" (id -1) se configura en el MISMO formulario de emparejamiento:
   - Modo single ("todo a una bodega"): se usa available_quantity total (ya incluye -1). Sin config extra.
   - Modo mapped: el dropdown de bodegas Siigo incluye -1 como opcion y es obligatorio
     mapearla a una bodega Probability (o marcarla "ignorar").
   - Nota: Woo recibe la SUMA de bodegas activas, asi que el mapeo de -1 afecta la exactitud
     interna por bodega, no el numero final que ve Woo (mientras se sumen todas).

2. Stock negativo se trata como 0 (max(0, qty)) para efectos de venta / push a Woo.
   Sub-decision pendiente: clampear SOLO al vender/empujar (guardando el real en Probability)
   vs clampear tambien en Probability. Recomendado: guardar real, clampear solo hacia Woo.

3. Filtrar productos fisicos vs no-fisicos por stock_control de Siigo, PERO igual importar
   todos a Probability para saber que existen:
   - stock_control=true  -> fisico -> track_inventory=true  -> cuenta y SE EMPUJA a Woo.
   - stock_control=false -> existe en Probability, track_inventory=false -> NO se empuja ni suma.
   - Consecuencia de scope: hoy el sync NO crea productos (match por SKU y si no existe, Skipped).
     "Que existan" exige agregar UPSERT de catalogo (crear/actualizar producto), no solo stock.

## Reglas de sanitizacion antes del push a Woo

- stock = max(0, round(available))  (sin negativos, sin decimales para unidades).
- Saltar productos sin stock_control (no fisicos) y sin SKU.
- Resolver el producto en Woo por SKU (GET wc/v3/products?sku=), incluyendo variaciones
  (products/{id}/variations/{vid}) para productos variables.
- Alinear SKU en los tres sistemas: code Siigo == sku Probability == sku WooCommerce.
  VALIDAR primero en la tienda real: si Woo no tiene SKU o difiere, el push no encuentra a quien actualizar.

## Arquitectura tecnica

Principios:
- Fuente de verdad = Siigo. Todo converge a valores ABSOLUTOS e IDEMPOTENTES. NUNCA decrementos relativos.
- El push a Woo es "pon el stock en X"; si Woo ya == X, no-op. Independiente del orden => a prueba de concurrencia.
- El talon de Aquiles NO es la concurrencia del decremento (ya es atomico via AdjustStockTx en transaccion),
  sino mezclar decrementos relativos con reconciliacion absoluta. Se evita usando SOLO absoluto.

Concurrencia (100 ordenes/min):
- NO serializar global (una-a-la-vez): una llamada lenta a Siigo/Woo congela todo. No escala.
- Serializar POR SKU (particion): dos ordenes del mismo SKU se ordenan; de SKUs distintos, en paralelo.
  Implementacion: RabbitMQ consistent-hash por SKU, o lock por SKU (Redis / SELECT FOR UPDATE en inventory_levels).
- Coalescing/debounce por SKU: si un SKU recibe N ordenes en poco tiempo, NO N pushes a Woo;
  se colapsa al ultimo valor absoluto y se empuja UNA vez.
- Reconciliacion periodica (job cada N min o nocturno): pull de Siigo + push absoluto a Woo,
  auto-sana cualquier drift. Es la red de seguridad; hace innecesario comparar en caliente tras cada orden.

Anti-patron a evitar:
- Comparar Woo vs Prob vs Siigo tras cada orden como GATE de correctitud. Los 3 se actualizan en
  momentos distintos (consistencia eventual) => siempre habra "diferencias" que son desfase, no error.
  Usar esa comparacion solo como AUDITORIA ASINCRONA de drift (alertar si un SKU lleva > X min descuadrado).

## Donde engancha en el codigo actual

- Sync stock Siigo->Prob: services/modules/inventory/internal/app/sync_provider_stock.go
  (aplica stock absoluto: delta = target - current; si 0 => Unchanged). Reusar patron absoluto.
- Canal Redis ya existe: shared/redis/channels.go -> "probability:inventory:state:events".
  Falta un CONSUMER WooCommerce que lo escuche.
- Cliente Woo (a extender): services/integrations/ecommerce/woocommerce/internal/infra/secondary/client/
  (hoy: get_orders.go, create_webhook.go). Agregar update de producto/variacion.
- Slot generico: integrations/core/internal/domain/provider_contract.go UpdateInventory (hoy ErrNotSupported).
  Sobrescribir en WooCommerceCore (.../woocommerce/internal/infra/secondary/core/core.go).
- Config de bodegas (front): SiigoInventorySection.tsx (agregar -1 al selector en modo mapped).
  Config del business para "que bodegas alimentan Woo".

## Trabajo a implementar (checklist)

Backend:
[ ] Cliente Woo: ResolveProductBySKU (GET wc/v3/products?sku=) + variaciones.
[ ] Cliente Woo: UpdateStock (PUT wc/v3/products/{id} o products/batch) con stock_quantity, manage_stock=true.
[ ] Sobrescribir UpdateInventory en WooCommerceCore (mapear SKU/business -> producto Woo).
[ ] Consumer Woo que escuche probability:inventory:state:events, con:
    - agregacion por SKU (suma de bodegas activas configuradas),
    - sanitizacion (max(0,round)), filtro fisicos (track_inventory),
    - coalescing por SKU, push absoluto idempotente (skip si igual).
[ ] Serializacion por SKU (consistent-hash o lock).
[ ] Job de reconciliacion periodica Siigo -> Woo.
[ ] Upsert de catalogo Siigo -> Probability (crear productos no existentes, set track_inventory por stock_control).
[ ] Config del business: modo bodegas + que bodegas alimentan Woo + clamp policy.

Frontend:
[ ] Selector de bodegas Siigo incluye "-1 Sin asignar" en modo mapped (obligatorio mapear/ignorar).
[ ] Config de "bodegas que alimentan WooCommerce".

## Decisiones pendientes (a resolver antes de codificar)

- Clamp de negativos: solo hacia Woo (guardando real en Prob) vs en ambos. (Recomendado: solo hacia Woo.)
- Upsert de catalogo: alcance real de "que existan en Probability" (crear todos vs solo marcar).
- Que bodegas alimentan Woo: todas las activas del business, o subconjunto configurable.
- Manejo de variaciones en Woo (productos variables) y su SKU por variacion.
- Timing del pull post-facturacion: Siigo descuenta al facturar; definir espera/reintento antes de
  jalar el stock (evitar leer Siigo antes de que descuente).

## Criterio para cerrar esta alerta

Cerrar cuando: (a) el push Probability->Woo este implementado con push absoluto idempotente,
serializacion por SKU y coalescing; (b) sanitizacion (negativos/decimales/filtro fisicos) aplicada;
(c) config de bodegas (incl. -1) en el form; (d) job de reconciliacion corriendo; y todo verificado
E2E contra la tienda WooCommerce de pruebas con la cuenta Siigo real. Actualizar items parciales
con fecha a medida que se resuelvan (no borrarlos).
