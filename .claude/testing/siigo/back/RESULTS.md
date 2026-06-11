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
