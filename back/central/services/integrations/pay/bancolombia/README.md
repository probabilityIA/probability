# Integracion Bancolombia QR

Pagos con codigo QR de Bancolombia + confirmacion de la transferencia via webhook.

## Flujo

```
pay.requests (gateway_code=bancolombia)
  -> pay router -> pay.bancolombia.requests
    -> consumer -> usecase.ProcessPayment
       -> OAuth2 client_credentials (token cacheado ~20 min)
       -> POST generacion de QR
       -> pay.responses (gateway_response: qr_value, qr_image_base64, raw)

Bancolombia POST /api/v1/webhooks/bancolombia (o /test)
  -> raw log en webhook_logs (retencion 15 dias)
  -> firma HMAC-SHA256 (si hay webhook_secret configurado)
  -> parse tolerante del payload
  -> pay.bancolombia.webhook.events
    -> modules/pay BancolombiaWebhookConsumer
       -> idempotencia por event_id (tabla bancolombia_webhook_events)
       -> match payment_transaction por reference
       -> APPROVED/SUCCESS -> completed | REJECTED/FAILED -> failed | CANCELLED/EXPIRED -> cancelled
       -> SSE payment completed/failed
```

## Credenciales (plataforma, integration_type `bancolombia_qr`)

| Clave | Uso |
|---|---|
| `client_id` / `client_secret` | OAuth2 produccion |
| `test_client_id` / `test_client_secret` | OAuth2 sandbox |
| `webhook_secret` / `test_webhook_secret` | Firma HMAC del webhook (opcional; sin secreto se acepta sin validar) |
| `merchant_id` / `merchant_name` | Datos del comercio para el QR |
| `scope` | Override del scope OAuth (default `QR-Management:write:app`) |
| `token_path` | Override del path del token (default `/security/oauth-provider/oauth2/token`) |
| `qr_generate_path` | Override del path de generacion de QR (default `/v1/qr-management/codes`) |

Base URLs en `integration_types`: prod `https://api.bancolombia.com`, sandbox `https://sb-api.bancolombia.com`.

## PENDIENTE: confirmar contra el API Market

La especificacion detallada del producto "QR Code" del API Market de Bancolombia
requiere cuenta en https://developer.bancolombia.com/ (el portal publico solo
documenta la autenticacion). Lo confirmado por documentacion publica:

- Token: POST `{base}/security/oauth-provider/oauth2/token`, header
  `Authorization: Basic base64(client_id:client_secret)`,
  `Content-Type: application/x-www-form-urlencoded`, body
  `grant_type=client_credentials&scope=...`, expira en 1200s.
- Media type: `Accept: application/vnd.bancolombia.v4+json`.

Lo que hay que confirmar al tener acceso al sandbox (ajustable via credenciales
de plataforma, sin recompilar):

1. Path exacto y body del endpoint de generacion de QR (`qr_generate_path`).
2. Scope exacto del producto QR (`scope`).
3. Formato y esquema de firma de la notificacion webhook (el parser es tolerante
   a variantes de nombres de campos y la firma soporta HMAC hex/base64 sobre el
   body crudo o en base64).

## Estado del pago saliente

El QR se publica en `pay.responses` como `status=success` con
`gateway_response.status=PENDING`: el QR generado no implica pago; el pago se
confirma solo cuando llega el webhook de transferencia.
