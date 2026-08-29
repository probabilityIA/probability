# Migraciones

## Como funciona

`Migrate()` en `internal/infra/repository/constructor.go` esta **en cero**:
ejecutar `go run cmd/main.go` no migra nada.

Eso es a proposito. Antes se re-ejecutaban las 60+ migraciones historicas en cada
corrida: lento, ruidoso, y sin ninguna utilidad porque ya estaban aplicadas.

## Flujo para migrar algo

1. Escribir la migracion en su archivo (`XXX_descripcion_corta.go`).
2. Agregar la llamada dentro de `Migrate()`:

```go
func (r *Repository) Migrate(ctx context.Context) error {
    return r.migrateLoQueSea(ctx)
}
```

3. Correr `cd back/migration && go run cmd/main.go`.
4. Verificar el efecto en la base.
5. **Dejar `Migrate()` en cero otra vez** (`return nil`).
6. Registrar la corrida en la tabla de abajo.

Los DDL (AutoMigrate, `CREATE ... IF NOT EXISTS`) se conservan como archivo.
Los DML y seeds puntuales se pueden borrar una vez aplicados en produccion.

`migrateHistorico()` conserva el orden original de las migraciones ya aplicadas.
No se llama desde ningun lado; sirve como referencia y para reconstruir un
entorno desde cero si algun dia hace falta.

## Historico

| Fecha | Migracion | Que hizo | Entorno |
|-------|-----------|----------|---------|
| 2026-08-28 | `migrateChannelRawDataNullable` | Quita el `NOT NULL` de `order_channel_metadata.raw_data` para que la retencion pueda vaciar el JSON crudo del canal en las ordenes de mas de 90 dias sin borrar la fila (el mapeo orden-canal se conserva). Sin esto la purga falla con SQLSTATE 23502 | local |
| 2026-08-28 | `migrateLegalDocuments` + `migrateUserTourProgress` + `migrateLegalContent` | Crea `legal_documents`, `legal_acceptances` (nunca se habian corrido: la aceptacion de terminos llevaba desde el 2026-08-27 sin poder mostrarse) y `user_tour_progress`. Agrega `content_html` a `legal_documents` con el texto completo de los dos documentos en HTML, y reescribe `sha256` como el hash de ese HTML, no del PDF: asi la evidencia de aceptacion corresponde a lo que el usuario leyo en pantalla | local + produccion |
| 2026-08-25 | `migrateEmailLogGeneric` (aplicada como DDL a mano por el tunel, sin correr el binario) | Vuelve generica la tabla `email_logs` (antes solo notificaciones): agrega `module`, `reference_type`, `reference_id` (indice compuesto), `provider`, `provider_message_id`, `sent_by`, `sent_by_name`, y pone `DEFAULT 0` a `integration_id`/`config_id`. Primer uso: historico de correos de cortes de pago COD (`module=cod_report`, `reference_type=cod_payment_cut`). Cualquier modulo que mande correo debe registrar ahi | local + produccion |
| 2026-08-20 | `migrateFreeTrialSubscriptionPlans` | Agrega `payable`, `trial_duration_days` a `subscription_types` y `overage_accepted`, `overage_accepted_at`, `overage_amount_due`, `overage_amount_paid_at` a `business_subscriptions`. Siembra los planes `free` (50 envios incluidos, $1000 de excedente) y `trial` (15 dias, todos los modulos, no pagable). Nota: el primer intento de correrla quedo por error dentro de `migrateHistorico` (que nunca se ejecuta) y no aplico nada hasta corregirlo el mismo dia; ademas el seed inicial dejo `trial.payable = true` por un default de GORM (campo con `default:true` en el tag se omite del INSERT cuando el valor Go es `false`, el zero-value), corregido con `UPDATE` directo y quitando el `default:true` del tag. | produccion |
| 2026-08-18 | `migrateSubscriptionAuditLogs` + `migrateBusinessModuleOverrideExpiry` | Crea la tabla `subscription_audit_logs` (faltaba desde que se agrego el feature de auditoria de suscripciones, dejaba el endpoint de auditoria en 500 en silencio) y agrega `expires_at` a `business_module_overrides` | produccion |
| 2026-08-18 | `migrateSubscriptionTypeOverage` | Agrega a `subscription_types` los campos `included_shipments`, `shipment_overage_price`, `included_invoices`, `invoice_overage_price` para modelar planes personalizados con cuota incluida + costo por unidad adicional (envios/facturas) | produccion |
| 2026-08-19 | `migrateSubscriptionTypeOverage` (2da corrida, AutoMigrate agrego columnas nuevas) | Agrega a `subscription_types` `included_orders` y `order_overage_price`: el limite de "Sin intermediarios" es por ordenes creadas, no por facturas, asi que el excedente necesitaba su propio contador (paralelo a envios y facturas) | produccion |
| 2026-08-06 | `fixVigaCodCalibracionFallida` | Ajusto `cod_total` y `cod_carrier_fee` de VIG-0071, VIG-0072 y VIG-0069 al valor que liquida EnvioClick (3 ordenes) | produccion |
| 2026-08-12 | `migrateSyncRunItemParent` | Agrego `parent_ref`, `parent_label` y `variant_label` a `integration_sync_run_items` (+ indice en `parent_ref`) para agrupar variantes por publicacion en el comparativo | produccion |
| 2026-08-12 | `migrateSyncRunChannelNoSKU` | Agrego `channel_no_sku` a `integration_sync_runs` para contar los items del canal que no tienen SKU y no se pueden emparejar | produccion |
| 2026-08-12 | `migrateSyncRunSKUTypo` | Agrego `sku_typo` a `integration_sync_runs` para contar los posibles errores de digitacion detectados | produccion |
| 2026-08-12 | `migrateProductIntegrationLastPushedQty` | Agrego `last_pushed_qty` a `product_business_integrations` para no repetir el push cuando el stock no cambio | produccion |
| 2026-08-12 | `migrateSyncRunSKUChanged` | Agrego `sku_changed` a `integration_sync_runs` para contar los mapeos cuyo SKU cambio en el canal | produccion |
| 2026-08-12 | `migrateProductIntegrationLogisticType` | Agrego `external_logistic_type` a `product_business_integrations` para saber si la publicacion es de fulfillment (ML administra el stock y no se puede empujar). AutoMigrate arrastro tambien un DEFAULT '[]' en `subscription_types.features`, que ya coincidia con el modelo | produccion |
| 2026-08-12 | `migrateProductIntegrationVariantUnique` | Reemplazo el unico `idx_product_integration` de `(product_id, integration_id)` por `(product_id, integration_id, COALESCE(external_variant_id, ''))`, para que un producto pueda mapearse a varias variantes del mismo canal. Se quitaron los tags `uniqueIndex` del modelo: ahora el indice lo maneja este SQL | produccion |
| 2026-08-12 | `migrateSiigoReferrals` | Crea la tabla `siigo_referrals` (name, email, phone, order_range) para el nuevo modulo de referidos Siigo (formulario publico en front/website) | produccion |
| 2026-08-12 | `migrateSyncRunTypoEvidence` | Agrego `sku_spacing` a `integration_sync_runs` y a `integration_sync_run_items` las columnas `counterpart_sku`, `counterpart_name`, `channel_qty`, `own_qty`, `fix_side` y `pattern`: la evidencia que sostiene cada sugerencia de correccion de SKU (que dice cada sistema, cuanto stock tiene cada lado y de que lado conviene corregir) | produccion |
| 2026-08-12 | `migrateProductFieldProvenance` | Crea `product_field_origins` (estado actual: quien fue el ultimo en escribir cada campo de cada producto, canal o usuario) y `product_field_changes` (historial append-only con el valor anterior, tope 5 por producto+campo, retencion 12 meses). Base del motor de comparacion de datos entre canales: sin esto no se puede advertir "300 de estos productos vienen de WooCommerce" ni deshacer una aplicacion masiva | produccion |
| 2026-08-12 | `migrateProductFieldProvenance` (2da parte) | Agrego `channel_snapshot` (jsonb) y `snapshot_at` a `product_business_integrations`: la foto de como se ve el producto en cada canal, tomada en la comparacion. Sin esto, abrir el diff de un producto tendria que pegarle en vivo a la API de cada canal | produccion |

| 2026-08-13 | `migrateInventoryCompareSnapshot` | Crea `inventory_compare_snapshots`: la foto del ultimo comparativo de inventario por integracion (stock de Probability vs stock del canal, accion, motivo e imagen del producto). Sin esto, abrir "Sincronizar inventario" le pega en vivo a la API de todos los canales cada vez y quema rate limit | produccion |

| 2026-08-17 | `migrateShippingQuotes` | Agrego `error_message` a `shipping_quotes`: el motivo por el que fallo la generacion de la guia, ya sanitizado para el cliente (sin nombre del proveedor). Sin esto el modulo de cotizaciones no puede explicar por que no salio la guia | produccion |
| 2026-08-17 | `migrateTrazabilidadUsuario` | Agrego `created_by`/`created_by_name`/`updated_by`/`updated_by_name` a `shipments` y `shipping_quotes`, y `updated_by`/`updated_by_name` a `orders`. Sin esto no se puede saber que usuario genero una guia, cotizo o modifico una orden | produccion |
| 2026-08-17 | `fixVig0095CotizacionFallida` | Marco como `failed` las cotizaciones 6542 (VIG-0095, business 46) y 6441 (DEM-0040, business 26), que seguian en `guide_generated` con el shipment fallido, y les cargo el motivo real (correo invalido y telefono invalido). DML puntual, se puede borrar | produccion |
| 2026-08-18 | `migrateMeliStatusMappings` | Sembro los 8 mapeos de estado de MercadoLibre en `order_status_mappings` (integration_type_id 3), que no existian: por eso las ordenes de ML quedaban con `status_id` NULL y la UI mostraba el string crudo ("paid") en gris. Ademas hizo backfill del `status_id` de las ordenes de ML ya cargadas (10 filas). DML/seed, se puede borrar | produccion |
| 2026-08-18 | `migrateOrderChannelPack` | Agrego `orders.channel_pack_id` (+ indice): el id del carrito del canal cuando la orden consolida varias ordenes de MercadoLibre. Sin esto no hay como distinguir en la UI una orden normal de una que agrupa un pack, ni saber por que su numero es el del pack. Incluye backfill desde el JSON crudo | produccion |

Antes de esta fecha no habia registro: todas las migraciones listadas en
`migrateHistorico()` se aplicaron corriendo la cadena completa.
