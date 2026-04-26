# SES - Gestion de acceso a produccion

## Estado actual (2026-04-24)

- SES sigue en **sandbox** (`ProductionAccessEnabled: false`)
- Solicitud DENIED bloqueada por Trust and Safety (case `177111462900461`)
- Motivo: key `AKIAW57NP4YKBNCM7PWO` del usuario `backend-s3-uploader` fue comprometida (detectada 2026-02-14)
- Key rotada el 2026-04-24. Nueva key activa: `AKIAW57NP4YKHS67JG6T`

## Paso 1 - PENDIENTE: Responder al case de seguridad

Re-abrir el case 177111462900461 con esta URL y pegar el texto de abajo:

```
https://console.aws.amazon.com/support/home?region=us-east-1#/case?displayId=177111462900461&language=en
```

---

Hello,

I apologize for the delayed response. I was not aware this case required
action until I noticed it while requesting SES production access.

I have now completed all required steps:

Step 1 - Access key rotated:
The exposed key AKIAW57NP4YKBNCM7PWO belonging to IAM user
backend-s3-uploader has been permanently deleted. A new access key has
been issued and deployed to all application environments.

Step 2 - CloudTrail review:
I reviewed CloudTrail logs for all activity using the compromised key.
I found reconnaissance attempts from multiple IPs (enumeration of S3
buckets, IAM roles, EC2, Lambda, Bedrock). All privilege escalation
attempts were denied due to the restrictive IAM policy on the
backend-s3-uploader user (scoped to S3 and SES only). No unauthorized
resource creation, data exfiltration, or destructive actions were found.

Step 3 - Account-wide review:
No unauthorized IAM users, roles, policies, EC2 instances, Lambda
functions, or other resources were found.

The key was exposed because it was inadvertently included in a
documentation file. It has been removed. Going forward, credentials will
be managed exclusively through environment variables and AWS Secrets
Manager.

I believe the account is now secure. Please let me know if any further
action is required.

Best regards,
<contact-name> - <contact-email>

---

## Paso 2 - Despues de que Trust and Safety cierre el case

Reenviar solicitud de aumento de limite SES con el texto de abajo.
CaseId del bloqueo SES: `177682962600052`

---

Hello,

Thank you for the review. Below are the additional details requested.

Production access request - resubmission with additional details.
Service: SES Sending Limits
Region: us-east-1
Mail Type: TRANSACTIONAL
Website: https://www.probabilityia.com.co
Estimated volume: under 10,000 emails per month

1. PLATFORM OVERVIEW
Probability is a B2B multi-tenant SaaS that centralizes e-commerce
operations (orders, payments, shipments) for businesses selling through
Shopify, Amazon, MercadoLibre and WhatsApp. Our customers are registered
businesses; the email recipients are the end customers of those
businesses.

2. HOW WE BUILD RECIPIENT LISTS
We do NOT import, buy or scrape email lists. Every recipient email is
obtained through explicit opt-in in one of two ways:
  a) The end customer placed an order on the business's Shopify/Amazon/
     MercadoLibre store (providing their email as part of checkout).
  b) The end customer registered directly on the business's Probability
     storefront, accepting terms of service and privacy policy.
Each business onboarded on Probability signs our terms agreeing that
they only upload consented data.

3. EMAIL TYPES SENT (all transactional, triggered by user action)
  - Order confirmation: sent immediately after an order is created.
  - Shipment status update: sent when a carrier webhook updates the
    tracking state (picked up, in transit, delivered, exception).
  - Payment notification: sent when a payment gateway confirms capture,
    refund or chargeback.
  - Password reset: sent only when the user clicks "forgot password".
  - Account verification: sent once when a new user registers.
  - Operational alerts: low stock, integration errors (sent to the
    business operator, not to end customers).
We do NOT send newsletters, promotions or marketing campaigns.

4. BOUNCE AND COMPLAINT HANDLING
  - We subscribe to SES bounce/complaint notifications via SNS.
  - Hard bounces: the recipient is added to our internal suppression
    list and the database record is flagged invalid. We never retry.
  - Complaints: the recipient is immediately suppressed and the
    originating business is notified so they can review consent.
  - We respect SES account-level suppression list.
  - Target bounce rate: below 2%. Target complaint rate: below 0.1%.

5. UNSUBSCRIBE
Although transactional emails do not legally require an unsubscribe
link, every email includes a footer with the originating business
contact details and a support mailbox to stop further communication.

6. COMPLIANCE
  - CAN-SPAM: accurate From, Reply-To, subject; no deceptive headers;
    physical address of the sending business in footer.
  - GDPR / Colombian Ley 1581: we only send on explicit consent, we
    honor erasure requests within 30 days.

7. INFRASTRUCTURE
  - Domain probabilityia.com.co is verified with DKIM (SUCCESS), SPF
    (include:amazonses.com) and DMARC (p=none, rua configured) already
    published.
  - From address: noreply@probabilityia.com.co
  - Emails are sent from our Go backend using the AWS SDK for Go v2
    (sesv2 SendEmail). No third-party relays.
  - Logging and monitoring of delivery metrics via CloudWatch.

Note: The previous security issue (case 177111462900461) has been fully
resolved. The compromised IAM key was rotated and all account resources
were verified clean.

Please let us know if any further detail is required.

Best regards,
<contact-email>
