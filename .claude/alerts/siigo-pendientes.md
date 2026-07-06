# ALERTA: Pendientes integracion Siigo

Creada: 2026-06-10
Contexto: se implemento paridad de operaciones Siigo vs Softpymes via colas
(commit "feat(siigo): paridad de operaciones con softpymes via colas de facturacion").
Las 6 operaciones (cancel, check_status, compare, list_items, cash_receipt,
list_bank_accounts) compilan y pasan tests, pero quedan pendientes criticos.

## Urgente (riesgo o bloqueante de uso)

1. **Idempotencia en `retry` de Siigo**: RESUELTO EN CODIGO (2026-07-01, pendiente
   validar E2E con escritura real, ver punto 3). Se replico el patron de Softpymes:
   - Al crear se guarda el marcador "order:<orderID> | #<orderNumber>" en el campo
     "observations" de la factura Siigo (mapper).
   - Antes de todo POST /v1/invoices (create y retry), el cliente busca en Siigo una
     factura VIGENTE (no anulada) para esa orden via GET /v1/invoices filtrando por
     customer_identification + created_start/created_end y matcheando observations.
     Si existe, no re-factura (short-circuit devolviendo la existente + su CUFE).
   - Respeta la regla "cancelada libera la orden": las facturas anuladas se ignoran
     en el match, asi que una orden con su factura cancelada puede volver a facturarse
     (la regla "una sola vigente" ya la enforcea el modulo en InvoiceExistsForOrder,
     que excluye status failed/cancelled; es provider-agnostica).
   - Config opcional en la integracion: idempotency_check (bool, default true para
     desactivar el chequeo) e idempotency_lookback_days (int, default 30).
   Archivos: client/find_existing_invoice.go (nuevo), client/create_invoice.go,
   client/list_invoices.go, client/mappers/invoice.go, consumer/invoice_request_consumer.go,
   domain/dtos/invoice_types.go, domain/dtos/customer_types.go.
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
