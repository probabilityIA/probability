# WhatsApp numero propio por cliente - pendientes

Fecha: 2026-09-05
Plan: `.claude/project/whatsapp-numero-por-cliente.md`
Detalle tecnico: `back/central/services/integrations/messaging/whatsapp/NUMERO-PROPIO.md`
Rama: `claude/whatsapp-numero-por-cliente-s22bjf`

El codigo de las fases 1 a 5 esta escrito y compila; `go build ./...` y
`go test ./...` de `back/central` pasan, y el front typechequea. Lo que sigue NO
se ejecuto: esta sesion no tiene base de datos, ni Redis, ni token de Meta.

## Urgente (antes de desplegar)

1. **Correr la migracion.** `Migrate()` en
   `back/migration/internal/infra/repository/constructor.go` quedo apuntando a
   `migrateWhatsappInboundConversationType`, que extiende el CHECK de
   `whatsapp_conversations.conversation_type` para aceptar `inbound`.
   Sin ella, toda conversacion entrante creada por el ruteo por numero falla al
   persistirse con SQLSTATE 23514. Correrla en local, verificar, dejar
   `Migrate()` en cero y registrar la corrida en `back/migration/MIGRACIONES.md`.

2. **Probar la fase 1 contra la base local con el numero de test.** Un negocio
   sin `phone_number_id` propio debe seguir enviando por el numero de
   Probability (no romper a nadie), y uno con `phone_number_id` + `waba_id` +
   `access_token` debe enviar por el suyo. Es el caso que mas puede romper en
   produccion porque toca a TODOS los negocios de hoy.

3. **Poner `waba_id` en las credenciales de plataforma del integration_type 2.**
   Sin ese campo, `POST /whatsapp/templates/provision` no tiene de donde copiar
   las plantillas y devuelve error. Valor de produccion: `1302830408357767`.

## Importante

4. **Fase 3 sin verificar contra un WABA real.** El aprovisionamiento de
   plantillas y el consumo de `message_template_status_update` estan escritos
   pero nunca corrieron contra Meta. Necesitan un cliente conectado o un WABA de
   pruebas compartido con Probabilityapp.

5. **Quality rating por numero.** El riesgo que el propio plan senala sigue sin
   cubrir: hoy se vigila un `quality_rating` y con N clientes son N. La UI
   muestra el estado de plantillas, no la calidad del numero ni la ventana de
   24 h. Falta leer `GET /{phone_number_id}?fields=quality_rating` y exponerlo.

6. **Sin ticket.** La regla `.claude/rules/tickets.md` pide registrar el trabajo
   en un ticket y cerrarlo con el diagnostico. No se pudo: el modulo de tickets
   exige JWT de super admin contra un backend corriendo. Crear el ticket y pegar
   ahi el resumen de la rama.

## Deseable

7. **Embedded Signup.** El camino largo del plan sigue pendiente y depende de
   sacar la app de dev mode, publicar politica de privacidad / terminos /
   eliminacion de datos y pasar App Review con `business_management`. El codigo
   de las fases 1 a 5 sirve igual para ese camino; lo unico que cambia es la
   pantalla de conexion, que pasa a ser un boton.

## Cuando se cierra esta alerta

Cuando 1, 2 y 3 esten hechos y verificados, y 4 tenga al menos una corrida real
contra un WABA. Los puntos 5 a 7 se mueven a `ROADMAP.md` si siguen abiertos.
