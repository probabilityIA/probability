# WhatsApp: numero propio por negocio

Como un negocio envia y recibe WhatsApp desde SU numero, con Probability como
administrador de su cuenta, sin que el cliente haga tramites con Meta.

Contexto y fases: `.claude/project/whatsapp-numero-por-cliente.md`.

## Como se decide de que numero sale cada mensaje

Desde 2026-09-05 un negocio puede enviar y recibir desde **su propio numero**, con
Probability como administrador de su cuenta. Lo decide su fila de `integrations`
(tipo 2):

| `config` | Efecto |
|---|---|
| sin `phone_number_id`, o `use_platform_token: true` | sale por el numero de Probability (comportamiento historico) |
| `phone_number_id` + `waba_id` + credencial `access_token` | sale por el numero del negocio |

`credentials_cache.GetWhatsAppConfig` lee esa fila (config + credenciales
desencriptadas, cacheadas por `integrations/core` en `integration:meta:*` y
`integration:creds:*`) y solo cae a la plataforma cuando falta configuracion
propia. Si el negocio declara `phone_number_id` pero le falta el token, el envio
falla con un error **permanente**: es una configuracion rota, no un fallo
transitorio.

**Lo que NO migra al numero del negocio** (siempre sale del de Probability):

- alertas de monitoreo (`consumeralert`)
- OTP de login (`consumerauthotp`)
- respuestas del agente AI de la plataforma (`consumerai`)
- aviso de saldo bajo de billetera y de ventana de pago de suscripcion
  (`consumerwalletalert`, `consumersubscriptionalert`): son mensajes de
  Probability al negocio, por eso usan `SendPlatformTemplate`.

### Ruteo del webhook por numero

El webhook es uno solo para todos (la firma HMAC usa el `webhook_secret` de la
app, igual para todos). Cuando llega un mensaje sin conversacion activa ni sesion
humana, se resuelve el dueno con
`phone_number_id -> (integration_id, business_id)`: un SELECT sobre `integrations`
filtrando `config->>'phone_number_id'`, con indice en Redis
(`integration:idx:cfg:2:phone_number_id:<id>`) que se mantiene al cachear la
integracion y se invalida al cambiarle el `config`.

- Si el numero es de un negocio: se abre una conversacion `inbound` para ese
  negocio, se persiste el mensaje y se publica por SSE a su dashboard.
- Si no es de nadie (o es el numero de la plataforma): sigue el camino de
  siempre, el agente AI Sales.

### Onboarding manual (camino corto, sin Embedded Signup)

1. El cliente entra a su Business Manager -> Configuracion del negocio ->
   Cuentas de WhatsApp -> Socios -> agrega el Business ID de Probabilityapp con
   control total.
2. Se anota su `waba_id` y su `phone_number_id` y se suscribe la app:
   `POST /{waba_id}/subscribed_apps`.
3. En el front, en la integracion de WhatsApp del negocio, se activa "Usar mi
   propio numero" y se pegan `waba_id`, `phone_number_id` y el token.
4. Se pulsa "Crear las plantillas que faltan" y se espera la aprobacion de Meta.

Si el cliente revoca el acceso desde su Business Manager, Meta responde 190 / 401
y el clasificador lo trata como error permanente: el mensaje se descarta con
`Warn` en vez de reencolarse en bucle (ver
`.claude/rules/colas-errores-permanentes.md`).

## Plantillas por WABA

Las 13+ plantillas viven en el WABA de Probability. Un negocio con numero propio
necesita las suyas, aprobadas por Meta en SU cuenta.

- `POST /whatsapp/templates/provision` lee las plantillas del WABA de la
  plataforma (`GET /{waba_id}/message_templates`) y crea en el WABA del cliente
  las que le falten (`POST /{waba_id}/message_templates`). El catalogo origen es
  Meta, no una copia en Go: asi no se desincronizan.
- `GET /whatsapp/templates/status` devuelve el estado por plantilla. Se cachea en
  Redis (`whatsapp:templates:{integration_id}`, TTL 6 h) y `?refresh=true` fuerza
  la consulta a Meta.
- El webhook `message_template_status_update` actualiza ese cache usando el WABA
  del `entry.id` (`whatsapp:waba:{waba_id}` -> integration_id).

Mientras una plantilla no este `APPROVED`, los mensajes que la usan no se pueden
enviar: es el estado "cliente conectado pero sin plantillas aprobadas" que la UI
muestra explicitamente.

