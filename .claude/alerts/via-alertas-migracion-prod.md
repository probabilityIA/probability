# Via: alertas en el chat - migracion pendiente en produccion

Fecha: 2026-09-16. Ticket: TKT-000096.

## Contexto

El asistente Via ahora tiene pestana de Alertas, aviso emergente y un canal de
notificacion nuevo (`notification_types.code = 'assistant'`, id 6) que se
configura en Notificaciones > Reglas igual que WhatsApp o correo.

El codigo del backend asume que existen las tablas `assistant_alerts` y
`assistant_alert_cursors` y que el canal `assistant` tiene **id 6**.

## Urgente

- [ ] Correr en produccion, en este orden y **antes** de que el deploy del
  backend quede activo, con `./scripts/run-migration-prod.sh` (agregarlas a
  `Migrate()` y dejarlo en cero despues):
  1. `migrateCancelRequestedStatus` (estado `cancel_requested` + plantilla)
  2. `migrateAssistantAlerts` (tablas de alertas + canal Via id 6)
  3. `migrateDefaultNotificationRules` (regla predeterminada en los negocios)
  El orden importa: la 2 necesita el estado de la 1 y la 3 necesita el evento de la 2. Sin la migracion:
  el consumidor de `ai.assistant.alerts` falla al insertar (5 reintentos y DLQ)
  y `GET /ai/assistant/alerts*` responde 500.
- [ ] Verificar `SELECT id, code FROM notification_types WHERE code='assistant'`
  devuelve 6. Si produccion ya tenia un id 6 con otro codigo, la migracion se
  detiene con error y hay que cambiar `NotificationTypeAssistant` en
  `services/events/internal/domain/dtos/event_types.go` y
  `assistantNotificationTypeID` en la migracion al id real.

## Importante

- [x] La plantilla `solicitud_cancelacion_recibida` (id 1782177306268808) quedo
  `APPROVED` en Meta el 2026-09-16.
- [ ] Informar a los negocios: los filtros de estado de las reglas de WhatsApp
  nunca se habian guardado (62 reglas en prod, 0 estados guardados). Desde este
  deploy si se guardan; quien los haya elegido tiene que volver a elegirlos.

- [ ] Crear al menos una regla de prueba en el negocio Demo (26): canal Via,
  evento "Cambio de estado de la orden" filtrado a Cancelada, integracion
  Plataforma, y cancelar una orden para ver la burbuja.

## Criterio para cerrar

Migracion corrida en produccion, canal con id 6 verificado y una alerta real
vista en el chat de Via en produccion. Entonces mover TKT-000096 a `testing`.
