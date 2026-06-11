# E2E Siigo - Resultados

Fecha: 2026-06-11

## Setup

- Mock Siigo: `back/testing` (servicio testing en tmux `prob:testing`), puerto 9095.
  Endpoints: /auth, /v1/auth, /v1/customers (POST/GET), /v1/invoices (POST/GET/list),
  /v1/invoices/:id, /v1/invoices/:id/annul, /v1/invoices/:id/stamp/errors,
  /v1/products, /v1/payment-types, /v1/vouchers, /v1/journals.
- Backend: `back/central` en tmux `prob:backend` (DB RDS prod, RabbitMQ/Redis local).
- Business demo id=26. Integracion Siigo creada via API: id=198, credencial
  `api_url=http://localhost:9095`, is_testing=false (el switch is_testing apunta a
  base_url_test=http://back-testing:9095, host docker no resoluble corriendo go run local;
  por eso se usa api_url como URL efectiva del mock).
- Operaciones disparadas publicando a la cola local `invoicing.siigo.requests`
  (identico a lo que hace el orquestador) y verificando mock + cola `invoicing.responses`.

## Resultados por operacion

| # | Operacion | Resultado | Evidencia |
|---|-----------|-----------|-----------|
| 1 | test_connection | OK | POST /integrations/test -> "Conexion probada exitosamente" (mock /auth) |
| 2 | create | OK | Factura FE-24446-1002, CUFE 64 chars, stamp Stamped |
| 3 | check_status | OK | GetInvoiceByID -> stamp_status=Stamped -> success |
| 4 | cash_receipt | OK | POST /v1/vouchers, balance 100000 -> 0 |
| 5 | list_items | OK | GET /v1/products -> count=3 |
| 6 | list_bank_accounts | OK | GET /v1/payment-types?document_type=RC -> accounts=3 |
| 7 | compare | OK | ListInvoices paginado -> CompareResponse en Redis + SSE |
| 8 | cancel | OK | POST /v1/invoices/{id}/annul -> status=annulled |

Cola siigo drenada (0 pendientes), respuestas consumidas por el orquestador (0 pendientes).

## Aislamiento multi-integracion (Softpymes intacto)

- Cero cambios de codigo en `softpymes/`.
- Cola `invoicing.softpymes.requests`: deliver_total=0 (nunca tocada).
- Mock softpymes coexiste en 9090.
- Business demo con AMBOS facturadores activos: Softpymes id=48 + Siigo id=198.

## Pendiente / notas

- Validacion contra cuenta Siigo real con escritura sigue pendiente (ver
  `.claude/alerts/siigo-pendientes.md`). Los paths annul/vouchers/stamp y los valores
  de stamp.status estan segun documentacion, verificados solo contra el mock.
- Integracion 198 (Siigo E2E Demo) quedo creada en la DB del business demo.

## E2E en PROD via WooCommerce (2026-06-11)

Flujo real: 20 webhooks WooCommerce (order.created) -> ordenes demo (business 26,
integracion woo 197) -> auto-factura a Siigo -> mock prod back-testing:9095.

- 20/20 webhooks 200, 20/20 ordenes creadas, 20/20 auto-facturadas (issued + CUFE).
- Operaciones Siigo validadas en prod sobre esas facturas:
  - cancel -> status cancelled.
  - cash_receipt -> voucher (sync log success). Requiere send_cash_receipt + cash_receipt_document_id/payment_id en config de la integracion.
  - credit_note -> NC-701 issued con CUFE. Requiere credit_note_document_id en config de la integracion.
- Validacion visual con Playwright en /invoicing/invoices (proveedor Siigo, Emitida/Cancelada).

### Bugs reales encontrados y corregidos (commits en main)
1. fix(woocommerce): mapper emitia payment status 'paid' pero ordenes espera 'completed' -> woo nunca quedaba is_paid.
2. fix(woocommerce): mapper no seteaba Invoiceable -> woo nunca facturable (regla COP como Shopify).
3. fix(invoicing): credit_notes.created_by_id NOT NULL no se asignaba -> 500 al crear NC.
4. feat(siigo): notas de credito end-to-end (no existian en ningun consumer).

### Config requerida en la integracion Siigo para todas las operaciones
document_id, payment_method_id, cash_receipt_document_id, cash_receipt_payment_id,
credit_note_document_id (todos en el config de la INTEGRACION, no solo en invoicing_config).
