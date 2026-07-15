# Inventario saliente filtrado por canal (product_business_integrations)

Fecha: 2026-07-15

## Contexto

Se agrego la asociacion producto<->canal en `product_business_integrations`:
- Al crear productos desde una integracion (upsert provider) se crea el mapping.
- El reconcile de Siigo detecta "existen en Probability pero no asociados a este
  canal" y permite asociarlos (todos o seleccionados).

Falta cerrar el circulo en la direccion SALIENTE.

## Alcance: SOLO canales de venta salientes

El filtro por asociacion aplica UNICAMENTE a canales de venta SALIENTES
(Probability -> WooCommerce, MercadoLibre, Shopify). NO aplica a Siigo.

Siigo es fuente de la verdad (ERP/facturacion), su sync es ENTRANTE
(Siigo -> Probability) y empareja por SKU a proposito. NO filtrar el entrante de
Siigo por `product_business_integrations` (ademas de que crearia un problema de
orden: crea+asocia y encola stock en la misma corrida, ambos async).

## Pendiente

### Importante
- La sincronizacion de inventario SALIENTE (Probability -> canales de venta:
  WooCommerce, MercadoLibre, Shopify) debe **notificar solo los productos
  asociados a ese canal** en `product_business_integrations`, no todo el catalogo.
  Objetivo: no empujar stock a un canal que no vende ese producto.
- Auditar cada integracion saliente para ver como selecciona hoy los productos a
  los que empuja stock (hoy varias cruzan por SKU / auto-mapean, no filtran por
  asociacion explicita).

### Deseable
- Verificar/crear indices en `product_business_integrations`: `integration_id` y
  `(product_id, integration_id)`, para que el filtro saliente sea rapido leyendo
  de DB (se decidio NO cachear en Redis: son lookups indexados baratos y el cache
  agregaria riesgo de stale justo en la exactitud que se busca).

## Criterio para cerrar

- Cada flujo de inventario saliente filtra por `product_business_integrations`.
- Indices confirmados en la tabla.
