# ALERTA: Pendientes integracion Siigo

Creada: 2026-06-10
Contexto: se implemento paridad de operaciones Siigo vs Softpymes via colas
(commit "feat(siigo): paridad de operaciones con softpymes via colas de facturacion").
Las 6 operaciones (cancel, check_status, compare, list_items, cash_receipt,
list_bank_accounts) compilan y pasan tests, pero quedan pendientes criticos.

## Urgente (riesgo o bloqueante de uso)

1. **Idempotencia en `retry` de Siigo**: el retry hace POST directo a /v1/invoices.
   Un mensaje RabbitMQ re-entregado puede facturar DOS VECES. Softpymes ya lo
   resuelve buscando la factura existente antes de crear (ListDocuments por fecha).
   Replicar ese patron en el consumer de Siigo usando ListInvoices.
2. **Config de Siigo no capturable desde la UI**: el form solo pide credenciales.
   Faltan campos para document_id, tax_id, seller_id, cash_receipt_document_id,
   cash_receipt_payment_id. Sin esto cash_receipt y la creacion bien configurada
   dependen de editar la config a mano en DB.

## Importante (antes de aprobar produccion)

3. **Validacion E2E contra cuenta Siigo real con escritura**: paths de
   /v1/invoices/{id}/annul y /v1/vouchers, valores reales de stamp.status, y el
   campo status del listado (del que depende detectar anuladas en compare) estan
   segun documentacion, NO verificados en vivo. La cuenta de pruebas actual es
   prestada y SOLO LECTURA (ver memoria: jamas POST/PUT/DELETE con ella).
4. **Seed inactivo**: invoicing_provider_types tiene Siigo con is_active = FALSE.
   Activar manualmente al salir a produccion.
5. **Partner-Id definitivo**: acordar con Siigo el Partner-Id camelCase; para
   volumen formal piden alianza de partner. Rate limit: 10 req/min en pruebas,
   100 req/min en produccion (compare paginado puede toparse el limite en pruebas).

## Deseable (roadmap)

6. Catalogos para configuracion: GET /v1/document-types, /v1/taxes, /v1/users,
   /v1/cost-centers (payment-types y products ya estan).
7. Moneda dual (*Presentment) y UnitPriceBase en DTOs de Siigo (ordenes USD).
8. SiigoIntegrationView en frontend (Softpymes ya tiene vista de config actual).
9. Notas de credito (POST /v1/credit-notes) - ningun proveedor las implementa.
10. Webhooks del proveedor y scheduler automatico de check_status / sync anuladas.

## Como cerrar esta alerta

Cuando los puntos 1-5 esten resueltos y validados E2E, eliminar este archivo.
