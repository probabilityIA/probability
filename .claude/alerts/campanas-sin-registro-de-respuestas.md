# Una campana nunca registra entregas ni respuestas

**Fecha:** 2026-09-12
**Modulo:** notificaciones (campanas de WhatsApp)

## Contexto

Salio de las pruebas E2E contra el mock local
(`.claude/testing/notification-config/back/RESULTS.md`). La campana se envio, el
cliente contesto cuatro veces y el contador de respuestas quedo en cero.

`whatsapp_campaign_sends.status` solo recibe dos valores: `sent` o `failed`, que
son los unicos que publica `consumercampaign`. Nadie escribe nunca `delivered`,
`read` ni `replied`, aunque la tabla los acepta y el contador existe.

El dato si esta en el sistema: los webhooks de estado actualizan
`whatsapp_message_logs` (ahi si se ven `delivered` y `read`, y las respuestas
entrantes quedan como filas `inbound`). Lo que falta es correlacionar ese log
con el envio de la campana por `message_id`.

## Items

### Importante

- `replied_count`, y cualquier metrica de entrega a nivel campana, hoy son
  siempre cero. El negocio ve "1 de 1 enviados, 0 respondieron" aunque el
  cliente haya contestado. Es una metrica que miente, no una que falta.
- Correlacionar `whatsapp_message_logs.message_id` con
  `whatsapp_campaign_sends.message_id` al procesar los webhooks de estado y las
  respuestas entrantes, y marcar el envio como `delivered` / `read` / `replied`.

### Deseable

- Mostrar en el detalle de la campana el embudo real (enviados, entregados,
  leidos, respondieron) en vez de solo enviados y respondieron.

## Criterio para cerrar

Una campana de prueba contra el mock termina con `replied_count` mayor a cero
cuando el cliente simulado responde, y los envios muestran `delivered` / `read`
en la tabla de resultados.
