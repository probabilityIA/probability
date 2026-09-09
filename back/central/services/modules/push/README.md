# Push (Firebase Cloud Messaging)

Notificaciones push a la app movil. **No es un sistema paralelo de
notificaciones**: es un canal mas del modulo `notification_config`, igual que
WhatsApp, SSE y Email. Cada negocio elige que eventos quiere recibir.

Ticket: TKT-000078.

## Como viaja un evento

```
publisher del modulo (shipments, orders, pay)
  -> exchange events.exchange
  -> services/events HandleEvent            lee las configs del negocio en Redis
  -> case dtos.NotificationTypePush         push_publisher.go
  -> cola notification.push.requests
  -> este modulo: consumer -> usecase -> FCM
```

El `notification_types` con `code = 'push'` y sus `notification_event_types`
los siembra `migrate_push_notifications.go`. Eventos sembrados:
`shipment.guide_generated`, `order.status_changed`, `order.delivered`,
`wallet.low_balance`.

Un evento sin mensaje definido en `buildMessage` se ignora en silencio (ACK),
no se reintenta.

## A quien le llega

A los **usuarios del negocio**, no al cliente final. Ese es el corte frente al
canal de WhatsApp, que le escribe al comprador.

Los destinatarios salen de `device_tokens`, que se llena cuando la app registra
su token tras el login. Un usuario puede tener varios dispositivos.

## Endpoints

| Metodo | Ruta | Que hace |
|---|---|---|
| POST | `/api/v1/push/devices` | registra o actualiza el token del dispositivo |
| GET | `/api/v1/push/devices` | lista los dispositivos del usuario |
| DELETE | `/api/v1/push/devices?token=...` | da de baja un dispositivo |

Todos exigen JWT. El `business_id` sale del token; solo el super admin lo manda
en el body.

## Tokens muertos

FCM responde `UNREGISTERED` o `INVALID_ARGUMENT` cuando el token ya no sirve
(app desinstalada, datos borrados). Ese caso **desactiva el token y hace ACK**:
reintentar no lo va a arreglar. Solo se reencola ante fallos transitorios (5xx,
timeout, red). Ver `.claude/rules/colas-errores-permanentes.md`.

## Configuracion

```env
FCM_PROJECT_ID=probability-app-6de7c
FCM_CREDENTIALS_FILE=/ruta/al/fcm-sender.json
```

Alternativa para produccion: `FCM_CREDENTIALS_JSON` con el contenido en linea.

**Sin esas variables el modulo arranca igual**, loguea un `Warn` y descarta los
eventos de push. No tumba el backend ni bloquea otros canales.

La cuenta de servicio es `probability-fcm-sender@probability-app-6de7c.iam.gserviceaccount.com`,
con rol `roles/firebasecloudmessaging.admin` y nada mas. La llave **no vive en el
repo**.

No se uso el SDK de Firebase Admin: se habla HTTP v1 directo firmando con
`golang.org/x/oauth2/google`, que ya estaba en el arbol de dependencias.

## Lado movil

`mobile/mobile_central/lib/services/modules/push/`. El registro del token pasa
tras el login (`main.dart`, `_syncPushRegistration`) y la baja en el logout.
La app pide permiso de notificaciones la primera vez; si el usuario lo niega, no
se registra token y no se rompe nada.

`data.route` del mensaje dice a que pantalla llevar al tocar la notificacion.

## Pendiente

- No existe evento propio de **novedad de envio**. Hoy se cubre por
  `order.status_changed` filtrando `status_code = delivery_novelty`.
- Notificaciones al **mensajero** (ruta asignada, parada nueva): fuera de esta
  entrega.
- En Android, con la app en primer plano el sistema no dibuja la notificacion.
  Para eso hace falta `flutter_local_notifications`, que no se agrego todavia.
