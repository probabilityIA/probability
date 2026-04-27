# Bold Integration - Resumen de implementacion

Estado: MVP funcional end-to-end (sandbox + webhook + auditoria).
Fecha: 2026-04-26
Branch: main (cambios sin commitear)

---

## 1. Que se construyo

### Flujo de recarga via Bold (wallet)

1. Usuario click "Proceder al Pago" -> escoge Bold.
2. Frontend (`wallet/page.tsx::handleSelectBold`) -> `GET /api/v1/pay/wallet/bold/signature?amount=X&currency=COP&business_id=Y`.
3. Backend (`bold_recharge.go::BoldGenerateSignature`):
   - Carga credenciales desde `integration_types.platform_credentials_encrypted` (AES-256-GCM, cache Redis).
   - Genera `order_id = "WLT" + uuid`.
   - Calcula `hash = sha256(order_id + amount + currency + secret_key)`.
   - **Crea `wallet_transaction` PENDING** con:
     - `reference = order_id`
     - `integration_type_id = 23` (Bold)
     - `integration_id = <id de business_integration>`
     - `gateway_request = JSON {order_id, amount, currency, hash, public_key, environment, is_sandbox, generated_at}`
   - Retorna `{order_id, hash, public_key, amount, currency, is_sandbox}`.
4. Frontend bifurca:
   - `is_sandbox=true` -> `POST /pay/wallet/bold/simulate` (mock interno, sin abrir Bold).
   - `is_sandbox=false` -> carga `checkout.bold.co/library/bold.js`, abre `BoldCheckout` con la firma. Retorno via webhook.

### Webhook de confirmacion

`POST /api/v1/webhooks/bold` (publico, sin JWT):

1. `webhook_handler.go::HandleWebhook` lee body raw + header `x-bold-signature`.
2. `webhook.go::HandleIncomingWebhook`:
   - Carga `BoldConfig` (secret de produccion via integrations core).
   - Verifica `HMAC-SHA256(base64(body), secret_key)` en hex.
   - Si invalida -> 401 + log warning.
   - Parsea CloudEvents envelope (`id`, `type`, `data.merchant_reference`, etc.).
   - Publica a queue `pay.bold.webhook.events`.
3. Consumer (`bold_webhook_consumer.go`) -> `bold_webhook_processor.go::ProcessBoldWebhookMessage`:
   - **Idempotencia**: insert en `bold_webhook_events` con unique constraint sobre `bold_event_id`.
   - **Routing por prefijo de reference**:
     - `WLT*` o `BOLD_SANDBOX_WLT*` -> `processWalletRechargeWebhook`:
       - Lookup `wallet_transaction` por reference.
       - Si pending -> marca COMPLETED, suma amount al wallet balance.
       - Guarda raw payload en `transaction.gateway_response`.
       - Linkea `bold_webhook_events.wallet_transaction_id`.
     - Otros -> path original (`findBoldPaymentTransaction` -> `payment_transaction` para ordenes).

### Modo sandbox (mock)

`POST /pay/wallet/bold/simulate` valida que `environment=sandbox` y aprueba el tx pending sin tocar Bold real. Para QA y demos sin credenciales reales.

### Auditoria (patron Shopify/Sofpyme)

Cada wallet_transaction tiene:
- `gateway_request` jsonb -> lo que generamos y enviamos al gateway.
- `gateway_response` jsonb -> el webhook completo recibido (CloudEvents envelope).
- `integration_type_id` + `integration_id` -> quien proceso el pago.
- Cross-link via `bold_webhook_events.wallet_transaction_id`.

### UI - Historial de pagos

Columna nueva "Metodo" muestra:
- Logo del integration_type (`it.image_url` con prefijo `URL_BASE_DOMAIN_S3`).
- Badge color por proveedor: Bold morado, Nequi rosa, Debito manual gris.

---

## 2. Archivos modificados / creados

### Backend

```
back/migration/shared/models/wallet.go                                 [+integration_type_id, +integration_id, +gateway_request, +gateway_response]
back/migration/shared/models/bold_webhook_event.go                     [+wallet_transaction_id FK]
back/migration/internal/infra/repository/migrate_wallet_tx_integration.go  NUEVO
back/migration/internal/infra/repository/constructor.go                [+migrateWalletTxIntegration]

back/central/services/modules/pay/internal/domain/entities/wallet.go   [+IntegrationTypeID, +IntegrationID, +GatewayRequest, +GatewayResponse, +IntegrationImageURL]
back/central/services/modules/pay/internal/domain/dtos/bold_dtos.go    [shape: hash/public_key/is_sandbox; +BoldSimulateDTO; +BoldBusinessIntegration]
back/central/services/modules/pay/internal/domain/ports/ports.go       [+GetWalletTransactionByReference, +SaveWalletTransactionGatewayResponse, +LinkBoldWebhookToWalletTransaction, +GetBoldIntegrationForBusiness, +BoldSimulatePayment, BoldGenerateSignature(businessID,...)]

back/central/services/modules/pay/internal/app/bold_recharge.go        [crea PENDING tx + gateway_request, sandbox detection, response shape nueva]
back/central/services/modules/pay/internal/app/bold_simulate.go        NUEVO  [aprueba tx pendiente; fallback crea si no existe]
back/central/services/modules/pay/internal/app/bold_webhook_processor.go  [routing wallet vs payment, save gateway_response, link wallet_transaction_id]

back/central/services/modules/pay/internal/infra/secondary/repository/wallet_repository.go  [+GetWalletTransactionByReference, +SaveWalletTransactionGatewayResponse, GetTransactionsByWalletID JOIN integration_types con prefijo S3]
back/central/services/modules/pay/internal/infra/secondary/repository/bold_credentials.go  [auto-detect sandbox via test_api_key, +GetBoldIntegrationForBusiness]
back/central/services/modules/pay/internal/infra/secondary/repository/bold_webhook.go     [fix uuid.Nil bug, +LinkBoldWebhookToWalletTransaction]

back/central/services/modules/pay/internal/infra/primary/handlers/bold_handlers.go   [GET signature con businessID, POST simulate, response uniforme {success, data}]
back/central/services/modules/pay/internal/infra/primary/handlers/wallet_routes.go   [GET /bold/signature, POST /bold/simulate]
back/central/services/modules/pay/internal/infra/primary/handlers/response/wallet_responses.go   [+integration_type_id, +integration_id, +integration_name, +integration_image_url, +gateway_request, +gateway_response]
back/central/services/modules/pay/internal/infra/primary/handlers/mappers/wallet_mapper.go       [integrationNameFromReference, jsonOrNil]
```

### Frontend

```
front/central/src/services/modules/pay/domain/ports.ts                  [+simulateBoldPayment]
front/central/src/services/modules/pay/infra/repository/api-repository.ts  [getBoldSignature parseo seguro, +simulateBoldPayment]
front/central/src/services/modules/pay/infra/actions/index.ts           [+simulateBoldPaymentAction]
front/central/src/app/(auth)/wallet/page.tsx                            [handleSelectBold detecta is_sandbox, script con onerror+timeout, columna Metodo con logo+badge]
```

---

## 3. Endpoints

| Endpoint | Auth | Descripcion |
|---|---|---|
| `GET /api/v1/pay/wallet/bold/signature?amount=&currency=&business_id=` | JWT | Genera firma + crea PENDING tx |
| `POST /api/v1/pay/wallet/bold/simulate` body `{order_id, amount}` | JWT | Mock: aprueba tx en sandbox |
| `GET /api/v1/pay/wallet/bold/status/:id` | JWT | Consulta status en Bold |
| `POST /api/v1/webhooks/bold` header `x-bold-signature` | publico (HMAC) | Recibe eventos de Bold |

---

## 4. Schema DB

### `transaction` (wallet recharges)

Columnas nuevas:
- `integration_type_id bigint` (FK integration_types)
- `integration_id bigint` (FK integrations)
- `gateway_request jsonb`
- `gateway_response jsonb`

### `bold_webhook_events`

Columna nueva:
- `wallet_transaction_id uuid` (FK transaction)

Mantiene la existente:
- `payment_transaction_id bigint` (FK payment_transaction para ordenes)

---

## 5. Validacion E2E

| Caso | Resultado |
|---|---|
| Demo recarga 22.000 via Bold sandbox | tx PENDING -> webhook HMAC valido -> COMPLETED, balance +22k |
| `gateway_request` poblado | OK (8 campos: order_id, hash, amount, currency, public_key, environment, is_sandbox, generated_at) |
| `gateway_response` poblado | OK (CloudEvents envelope completo del webhook) |
| `bold_webhook_events.wallet_transaction_id` linkado | OK |
| Idempotencia (mismo bold_event_id 2 veces) | Segundo evento ignorado, balance no se duplica |
| Frontend muestra logo Bold | OK (vía S3 prefijo URL_BASE_DOMAIN_S3) |
| Backfill historicas: BOLD_SANDBOX -> 23, MANUAL -> 22 | 6 + 115 rows actualizadas |
| Flujo Nequi (`POST /pay/wallet/recharge`) | Sin cambios, sigue funcionando |

---

## 6. Sandbox vs produccion - como cambiar

El integration_type Bold guarda 4 credenciales en `platform_credentials_encrypted`:
- `api_key` + `secret_key` -> produccion
- `test_api_key` + `test_secret_key` -> sandbox

`bold_credentials.go::GetBoldCredentials` detecta automatico:
- Si `test_api_key` y `test_secret_key` no estan vacias -> sandbox.
- Sino -> produccion.

Para ir a Bold real:
1. Super admin edita "Tipo de Integracion Bold" en `/integrations?tab=types`.
2. Vacia los campos test (o los reemplaza con sandbox real de Bold).
3. Configura webhook URL en panel Bold apuntando a `https://<dominio>/api/v1/webhooks/bold`.

---

## 7. Pendientes / mejoras opcionales

- **HMAC dual prod/sandbox**: actualmente el verificador usa solo el secret de produccion. Para que sandbox real de Bold funcione, hay que extender `verifySignature` para intentar tanto prod como sandbox secret y aceptar si alguno matchea.
- **MAN_DEB_** (debitos por generar guia): siguen sin logo. Opcionalmente crear un integration_type "Sistema" o renderizar un icono local.
- **Pestana Pagos para roles no super-admin**: el demo no ve la categoria "Pagos" en `/integrations` porque le falta el permiso `Integraciones-Pagos` (Read). Asignarlo desde IAM.
- **Cache Redis de creds**: ya existe via `GetCachedPlatformCredentials`; verificar TTL e invalidacion al editar.
- **Job de reconciliacion**: cron que cruza `bold_webhook_events` con Bold API para detectar webhooks perdidos despues de 5 reintentos.
- **Tests unitarios**: el caso de uso `bold_recharge`, el processor del webhook, y el `verifySignature` necesitan cobertura.

---

## 8. Documentos relacionados

- Plan original: `.claude/plan/bold-integration.md`
- Reglas arquitectura: `.claude/rules/architecture.md`
- Reglas backend: `.claude/rules/backend-conventions.md`
- README modulo: `back/central/services/integrations/pay/bold/README.md`
