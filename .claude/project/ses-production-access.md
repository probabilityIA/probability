# SES - Solicitud de acceso a produccion

Fecha primera solicitud: 2026-04-21 (PENDING -> DENIED)
Fecha segunda solicitud: 2026-04-21 (enviada desde UI de SES Console)
Proyecto: Probability
Cuenta AWS: 476702565908 (perfil `probability`)
Region: us-east-1
Plan de soporte: Basic (no hay API de Support)

## Timeline

1. **2026-04-21 22:38** — Identidad del dominio `probabilityia.com.co` recreada y verificada con DKIM exitosamente.
2. **2026-04-21 ~22:40** — Primera solicitud enviada via `aws sesv2 put-account-details`. Estado: PENDING.
3. **2026-04-21 22:47 GMT-0500** — AWS responde DENIED por falta de informacion ("Envia otra solicitud proporcionando mas informacion"). CaseId: `177682962600052`.
4. **2026-04-21** — Intento de reenvio por CLI: ConflictException (caso aun bloqueado aunque DENIED).
5. **2026-04-21** — Segunda solicitud enviada manualmente desde UI de SES Console con descripcion detallada (texto completo mas abajo).

## Datos del formulario (segunda solicitud)

| Campo | Valor |
|---|---|
| Mail type | Transactional |
| Website URL | https://www.probabilityia.com.co |
| Contact language | English |
| Additional contacts | <contact-email> |
| Production access enabled | true |

### Use case description enviado (segunda solicitud)

```
Production access request for SES Sending Limits in us-east-1.

Mail Type: TRANSACTIONAL
Website: https://www.probabilityia.com.co
Estimated volume: under 10,000 emails per month

1) PLATFORM OVERVIEW
Probability is a B2B multi-tenant SaaS that centralizes e-commerce
operations (orders, payments, shipments) for businesses selling through
Shopify, Amazon, MercadoLibre and WhatsApp. Our customers are registered
businesses; the email recipients are the end customers of those
businesses.

2) HOW WE BUILD RECIPIENT LISTS
We do NOT import, buy or scrape email lists. Every recipient email is
obtained through explicit opt-in in one of two ways:
a) The end customer placed an order on the business Shopify/Amazon/
   MercadoLibre store, providing their email at checkout.
b) The end customer registered directly on the business Probability
   storefront, accepting terms of service and privacy policy.
Each business onboarded on Probability signs our terms agreeing that
they only upload consented data.

3) EMAIL TYPES SENT (all transactional, triggered by user action)
- Order confirmation: sent immediately after an order is created.
- Shipment status update: sent when a carrier webhook updates the
  tracking state (picked up, in transit, delivered, exception).
- Payment notification: sent when a payment gateway confirms capture,
  refund or chargeback.
- Password reset: sent only when the user clicks "forgot password".
- Account verification: sent once when a new user registers.
- Operational alerts (low stock, integration errors): sent to the
  business operator, not to end customers.
We do NOT send newsletters, promotions or marketing campaigns.

4) BOUNCE AND COMPLAINT HANDLING
- We subscribe to SES bounce and complaint notifications via SNS.
- Hard bounces: the recipient is added to our internal suppression list
  and the database record is flagged invalid. We never retry.
- Complaints: the recipient is immediately suppressed and the
  originating business is notified so they can review consent.
- We respect the SES account-level suppression list.
- Target bounce rate: below 2%. Target complaint rate: below 0.1%.

5) UNSUBSCRIBE
Although transactional emails do not legally require an unsubscribe
link, every email includes a footer with the originating business
contact details and a support mailbox to stop further communication.

6) COMPLIANCE
- CAN-SPAM: accurate From, Reply-To, subject; no deceptive headers;
  physical address of the sending business in footer.
- GDPR / Colombian Ley 1581: we only send on explicit consent, and we
  honor erasure requests within 30 days.

7) INFRASTRUCTURE
- Domain probabilityia.com.co is verified. DKIM: SUCCESS.
  SPF: include:amazonses.com. DMARC: p=none with rua configured.
- From address: noreply@probabilityia.com.co
- Emails are sent from our Go backend using the AWS SDK for Go v2
  (sesv2 SendEmail). No third-party relays.
- Delivery metrics logged and monitored via CloudWatch.

Contact: <contact-email>
```

## Use case de la primera solicitud (denegada, para referencia)

```
Probability is a multi-tenant e-commerce management SaaS platform
(orders, payments, shipments) integrated with Shopify, Amazon,
MercadoLibre and WhatsApp. Estimated volume below 10,000 emails per
month. All emails are transactional: order confirmations, shipment
status updates, payment notifications, password recovery and operational
alerts to the end user. Recipients are end customers of businesses
registered on our platform who opted in when purchasing or registering.
We handle bounces and complaints via SES notifications. We do not send
marketing emails.
```

## Motivo del denial (textual AWS)

> Gracias por enviar tu solicitud para aumentar tus limites de envio. No
> podemos acceder a tu solicitud en este momento. Envia otra solicitud
> proporcionando mas informacion y nuestro equipo podra revisar la
> solicitud.

## Configuracion del dominio (verificada)

Dominio: `probabilityia.com.co`
- VerificationStatus: SUCCESS
- DkimStatus: SUCCESS
- DKIM tokens (publicados como CNAME en micom.co):
  - dl2n3c26j4n752zrbjqksytfi4kti3u5
  - tsv7kjextfdfvjdmd5ei323nqknjcnve
  - ojaxay3yckpbxboryxy3ejvsxbzfmk4n
- SPF: `v=spf1 include:amazonses.com ~all`
- DMARC: `v=DMARC1; p=none; rua=mailto:admin@probabilityia.com.co`

## Identidades creadas en SES

- probabilityia.com.co (DOMAIN, SUCCESS)
- <contact-email> (EMAIL_ADDRESS, SUCCESS) — creada para probar en sandbox

## DNS

- Nameservers del dominio cambiados de `ns1-3.wordpress.com` a
  `nameserver01-04.mi.com.co` el 2026-04-21.
- Panel DNS activo: micom.co

## Prueba realizada en sandbox

Envio exitoso desde `noreply@probabilityia.com.co` a `<contact-email>`.
MessageId: `0100019db3488f21-ff08691e-99aa-4ae8-8e8b-85a0dbbcce0b-000000`

## Verificacion mañana

```bash
aws sesv2 get-account --profile probability --region us-east-1 \
  --query '{Production:ProductionAccessEnabled,Review:Details.ReviewDetails,Quota:SendQuota}'
```

Esperado: `ProductionAccessEnabled: true` y `ReviewDetails.Status: GRANTED`.

Si sigue PENDING pasadas 48h o cambia a DENIED otra vez: revisar correo
en `<contact-email>` por pedido de mas info y ajustar use case.

Para cerrar el caso de soporte manualmente (si quedo bloqueado):
```
https://support.console.aws.amazon.com/support/home#/case/?displayId=177682962600052
```
Con plan Basic no hay API para cerrarlo, se hace solo por UI.

## Comandos utiles

```bash
aws sesv2 get-account --profile probability --region us-east-1
aws sesv2 list-email-identities --profile probability --region us-east-1
aws sesv2 get-email-identity --email-identity probabilityia.com.co --profile probability --region us-east-1
aws sesv2 send-email \
  --from-email-address "noreply@probabilityia.com.co" \
  --destination "ToAddresses=<email>" \
  --content '{"Simple":{"Subject":{"Data":"Test","Charset":"UTF-8"},"Body":{"Html":{"Data":"<p>test</p>","Charset":"UTF-8"}}}}' \
  --profile probability --region us-east-1
```

## Variables .env (ya configuradas)

```
SES_REGION=us-east-1
SES_ACCESS_KEY=***
SES_SECRET_KEY=***
FROM_EMAIL=noreply@probabilityia.com.co
```
