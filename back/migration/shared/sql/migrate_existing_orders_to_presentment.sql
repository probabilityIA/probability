-- ============================================
-- Migraci?n de datos: Actualizar precios presentment desde raw_data
-- ============================================
-- Este script extrae valores de presentment_money del JSON raw_data
-- almacenado en order_channel_metadata y actualiza las nuevas columnas presentment
-- ============================================
-- IMPORTANTE: Este script solo actualiza ?rdenes de Shopify que tengan
-- presentment_money disponible en su raw_data JSON
-- ============================================

-- Actualizar ?rdenes con precios presentment desde raw_data
UPDATE orders o
SET 
    total_amount_presentment = COALESCE(
        NULLIF(
            (ocm.raw_data->'total_price_set'->'presentment_money'->>'amount')::NUMERIC,
            0
        ),
        total_amount_presentment
    ),
    subtotal_presentment = COALESCE(
        NULLIF(
            (ocm.raw_data->'subtotal_price_set'->'presentment_money'->>'amount')::NUMERIC,
            0
        ),
        subtotal_presentment
    ),
    tax_presentment = COALESCE(
        NULLIF(
            (ocm.raw_data->'total_tax_set'->'presentment_money'->>'amount')::NUMERIC,
            0
        ),
        tax_presentment
    ),
    discount_presentment = COALESCE(
        NULLIF(
            (ocm.raw_data->'total_discounts_set'->'presentment_money'->>'amount')::NUMERIC,
            0
        ),
        discount_presentment
    ),
    shipping_cost_presentment = COALESCE(
        NULLIF(
            (ocm.raw_data->'total_shipping_price_set'->'presentment_money'->>'amount')::NUMERIC,
            0
        ),
        shipping_cost_presentment
    ),
    currency_presentment = COALESCE(
        ocm.raw_data->>'presentment_currency',
        ocm.raw_data->'total_price_set'->'presentment_money'->>'currency_code',
        currency_presentment
    )
FROM order_channel_metadata ocm
WHERE o.id = ocm.order_id
  AND o.platform = 'shopify'
  AND ocm.channel_source = 'shopify'
  AND ocm.raw_data IS NOT NULL
  AND ocm.is_latest = true
  AND ocm.raw_data->'total_price_set'->'presentment_money'->>'amount' IS NOT NULL
  AND (o.total_amount_presentment = 0 OR o.total_amount_presentment IS NULL);

-- ============================================
-- Notas:
-- ============================================
-- 1. Este script es idempotente: puede ejecutarse m?ltiples veces sin duplicar datos
--    (usa COALESCE para no sobrescribir valores existentes)
--
-- 2. Solo procesa ?rdenes de Shopify con raw_data disponible y presentment_money
--
-- 3. Los valores presentment (moneda local) se extraen del JSON raw_data almacenado 
--    en order_channel_metadata
--
-- 4. Si un valor no existe en el JSON, se mantiene el valor existente o el default (0)
--
-- 5. Para items (order_items), la migraci?n es m?s compleja y requiere hacer match
--    entre items de BD y line_items del JSON. Los items se actualizar?n autom?ticamente
--    cuando se procesen nuevas ?rdenes con el c?digo actualizado.
--    Si necesita migrar items existentes, puede hacerlo ejecutando una sincronizaci?n
--    de las ?rdenes desde Shopify o creando un script adicional que haga match por SKU.
-- ============================================
