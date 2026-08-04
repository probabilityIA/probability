# Bancolombia QR: confirmar especificacion del API Market

**Fecha:** 2026-08-03

## Contexto

Se creo la integracion de pagos Bancolombia QR (modulo
`back/central/services/integrations/pay/bancolombia`, gateway_code `bancolombia`,
integration_type `bancolombia_qr`, `in_development=true`). La autenticacion OAuth2
(client_credentials, token 20 min, media type `application/vnd.bancolombia.v4+json`)
esta confirmada por documentacion publica, pero la especificacion detallada del
producto "QR Code" requiere cuenta en el API Market
(https://developer.bancolombia.com/, no resolvia DNS al momento del desarrollo).
Los valores no confirmados son configurables via credenciales de plataforma sin
recompilar (`scope`, `token_path`, `qr_generate_path`).

## Urgente

- Ninguno (la integracion esta en `in_development=true` y no expone flujo real).

## Importante (antes de activar `in_development=false`)

1. Crear cuenta/app en el API Market sandbox y confirmar: path y body exactos del
   endpoint de generacion de QR (default actual `/v1/qr-management/codes`), scope
   real del producto (default `QR-Management:write:app`).
2. Confirmar el formato de la notificacion webhook (campos y esquema de firma) y
   ajustar `parseWebhookPayload`/`verifySignature` en
   `pay/bancolombia/internal/app/webhook.go` si aplica.
3. Registrar las credenciales sandbox como platform credentials del tipo
   `bancolombia_qr` y probar E2E: generar QR (cola `pay.requests` con
   `gateway_code=bancolombia`) y simular webhook a
   `/api/v1/webhooks/bancolombia/test`.
4. Correr la migracion (crea `payment_webhook_events`, copia los eventos de
   `bold_webhook_events` y crea el integration_type).
   Recordar: el deploy NO corre migraciones, hay que ejecutarla a mano contra RDS.
   IMPORTANTE: la migracion debe correr ANTES de desplegar el backend nuevo, que
   ya no sabe escribir en `bold_webhook_events`.

## Deseable

- Icono `integration-types/bancolombia.png` en S3.
- Endpoint de consulta de estado del QR (`qr_query_path` ya reservado en config).
- Politica de retencion para `payment_webhook_events` (hoy crece sin limite).
- Eliminar `bold_webhook_events` cuando el backend nuevo lleve dias estable en
  produccion. La migracion la deja intacta a proposito, como respaldo.

## Criterio de cierre

QR generado contra el sandbox real de Bancolombia y webhook de confirmacion
procesado end-to-end (payment_transaction pasa a completed); paths/scope
confirmados contra la especificacion oficial.
