# Integracion Bold (Pasarela de Pagos)

Bold es una pasarela de pagos colombiana (datafonos, links de pago, boton web). Esta integracion permite generar links de pago, recibir notificaciones via webhook y consultar estado de transacciones.

## Que hace esta integracion

- Crea links de pago en `https://integrations.api.bold.co/online/link/v1` consumiendo solicitudes desde la cola `pay.bold.requests`.
- Publica respuestas estandarizadas en `pay.responses` para que `modules/pay` continue el flujo.
- Pendiente: recibir webhooks de Bold (`SALE_APPROVED`, `SALE_REJECTED`, `VOID_APPROVED`, `VOID_REJECTED`) y consultar estado de transacciones de forma robusta.

## Arquitectura (hexagonal)

```
internal/
├── domain/
│   ├── entities/        Resultado del link de pago
│   ├── ports/           IBoldClient, IIntegrationRepository, IResponsePublisher
│   └── errors/          Errores tipados
├── app/
│   ├── constructor.go
│   └── process_payment.go    Caso de uso: config -> link -> publish
└── infra/
    ├── primary/
    │   └── consumer/    Consume pay.bold.requests
    └── secondary/
        ├── client/      HTTP a Bold API
        ├── repository/  Lee credenciales encriptadas de integration_types
        └── queue/       Publica a pay.responses
```

## Credenciales

**NO se usan variables de entorno para credenciales Bold.** Patron del proyecto:

1. Credenciales viven en `integration_types.platform_credentials_encrypted` (bytea, AES-256-GCM)
2. Buscar por `code = 'bold_pay'`
3. La unica env var requerida es `ENCRYPTION_KEY` (clave maestra para descifrar)
4. Cache en Redis para no descifrar en cada request (TTL recomendado: 5-10 min, invalidar al actualizar)
5. Frontend (`my-integrations/`) puede solicitar la version descifrada via endpoint reveal con permisos de super admin

Estructura del JSON descifrado:
```json
{
  "api_key": "<bold-identity-key>",
  "secret_key": "<bold-secret-key-para-firma-y-webhook>",
  "environment": "sandbox" | "production"
}
```

## Endpoints externos (Bold)

| Operacion | Metodo | URL |
|---|---|---|
| Crear link de pago | POST | `https://integrations.api.bold.co/online/link/v1` |
| Consultar estado | GET | `https://integrations.api.bold.co/online/link/v1/{linkId}` |
| Webhook fallback (re-query) | GET | `https://integrations.api.bold.co/payments/webhook/notifications/{payment_id}` |

Headers de autenticacion:
```
Authorization: x-api-key <api_key>
Content-Type: application/json
```

## Webhook (a implementar)

- Endpoint propio sugerido: `POST /webhooks/bold` (en `modules/pay`, no aqui).
- Bold envia payload CloudEvents v1.0 con header `X-Bold-Signature`.
- Validar firma: `HMAC-SHA256(base64(body), secret_key)` en hex y comparar con header.
- Responder HTTP 200 en menos de 2s. Si no, Bold reintenta hasta 5 veces (15min, 1h, 4h, 8h, 24h).
- Eventos: `SALE_APPROVED`, `SALE_REJECTED`, `VOID_APPROVED`, `VOID_REJECTED`.
- Idempotencia obligatoria por `id` del evento.

## Flujo asincrono actual

```
modules/pay (request creada)
        v
pay.requests
        v
integrations/pay/router (enruta por gateway_code = 'bold')
        v
pay.bold.requests
        v
bold/consumer
        v
process_payment usecase -> bold_client.CreatePaymentLink
        v
pay.responses (resultado estandar)
        v
modules/pay (consume y actualiza estado)
```

## Sandbox / pruebas

Para forzar escenarios 3DS y motor antifraude, usar valores especificos en `total_amount`:
- `555001` -> 3DS aprobado
- `555002` -> 3DS rechazado
- (ver docs Bold para listado completo)

Limites por metodo:
- CREDIT_CARD, PSE, BOTON_BANCOLOMBIA, NEQUI: rangos entre $1.000 y $10.000.000 COP segun metodo.

## Documentacion oficial

- Portal: https://developers.bold.co
- API link de pagos: https://developers.bold.co/pagos-en-linea/api-link-de-pagos
- Webhooks: https://developers.bold.co/webhook
- Consulta de transacciones: https://developers.bold.co/pagos-en-linea/consulta-de-transacciones

## Estado actual

Ver `.claude/plan/bold-integration.md` para diagnostico completo y plan por fases.
