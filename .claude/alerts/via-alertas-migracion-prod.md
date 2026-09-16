# Via: alertas en el chat - migracion pendiente en produccion

Fecha: 2026-09-16. Ticket: TKT-000096.

## Contexto

El asistente Via ahora tiene pestana de Alertas, aviso emergente y un canal de
notificacion nuevo (`notification_types.code = 'assistant'`, id 6) que se
configura en Notificaciones > Reglas igual que WhatsApp o correo.

El codigo del backend asume que existen las tablas `assistant_alerts` y
`assistant_alert_cursors` y que el canal `assistant` tiene **id 6**.

## Urgente

- [ ] Correr `migrateAssistantAlerts` en produccion con
  `./scripts/run-migration-prod.sh` (agregarla a `Migrate()` y dejarlo en cero
  despues) **antes** de que el deploy del backend quede activo. Sin la migracion:
  el consumidor de `ai.assistant.alerts` falla al insertar (5 reintentos y DLQ)
  y `GET /ai/assistant/alerts*` responde 500.
- [ ] Verificar `SELECT id, code FROM notification_types WHERE code='assistant'`
  devuelve 6. Si produccion ya tenia un id 6 con otro codigo, la migracion se
  detiene con error y hay que cambiar `NotificationTypeAssistant` en
  `services/events/internal/domain/dtos/event_types.go` y
  `assistantNotificationTypeID` en la migracion al id real.

## Importante

- [ ] Crear al menos una regla de prueba en el negocio Demo (26): canal Via,
  evento "Cambio de estado de la orden" filtrado a Cancelada, integracion
  Plataforma, y cancelar una orden para ver la burbuja.

## Criterio para cerrar

Migracion corrida en produccion, canal con id 6 verificado y una alerta real
vista en el chat de Via en produccion. Entonces mover TKT-000096 a `testing`.
