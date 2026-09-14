# CU-03: Conversaciones sin leer, orden, baja de marketing y etiquetas

Recorre la bandeja de Conversaciones contra el mock de WhatsApp, sin tocar Meta.
Cubre lo agregado el 2026-09-13: contador de mensajes sin leer por chat, total
de conversaciones sin leer, orden "sin leer primero", marca de leido compartida
por negocio y telefono, marca roja de baja ("Dejar de recibir") y etiquetas de
orden y campana.

## Precondiciones

1. Base local (`./scripts/dev-db-switch.sh status` = LOCAL), negocio Demo (26).
2. Migracion `migrateWhatsappConversationReads` aplicada en local.
3. Backend en 3050 y mock en 9103 (montaje en CU-01).
4. `whatsapp_url` de `integration:platform_creds:2` apuntando al mock. Se
   escribe DESPUES de arrancar el backend: al reiniciar vuelve a Meta.
5. Mock con respuestas automaticas apagadas, para que cada respuesta del
   cliente la dispare la prueba y no el guion:

```bash
curl -s -X POST localhost:9103/_mock/reset
curl -s -X POST localhost:9103/_mock/script -H 'Content-Type: application/json' \
  -d '{"auto_reply":false,"script":{}}'
```

API: `http://localhost:3050/api/v1/notification-configs/message-audit`.

## Paso 1: crear audiencia de tres telefonos

Campana con el flujo 1 a tres clientes (A, B, C) elegidos a mano, ventana
abierta, lanzada. Esperar el tick del despachador (hasta 2 min) hasta ver los
tres envios en `GET localhost:9103/_mock/sent`.

**Esperado:** tres envios `sent` y las tres conversaciones en el listado con
la etiqueta de la campana (`campaign_name`) y `unread_count = 0`.

## Paso 2: una respuesta deja el chat sin leer

`POST /_mock/reply {"phone": A, "button": "Si"}`.

**Esperado:** A con `unread_count = 1`, `unread_conversations = 1`.

## Paso 3: sin leer va primero aunque haya actividad mas reciente

B responde y se marca leido (`POST /conversations/{id}/read`). B queda con la
actividad mas reciente, pero leido.

**Esperado:** A (sin leer, mas vieja) sigue en la posicion 1 y B despues.

## Paso 4: marcar leido

`POST /conversations/{A}/read?business_id=26`.

**Esperado:** `{"success": true}`, A con `unread_count = 0`,
`unread_conversations = 0`, y sigue asi al volver a consultar.

## Paso 5: aislamiento por negocio

El mismo `POST .../read` con `business_id` de otro negocio.

**Esperado:** 404, y la marca de A no cambia.

## Paso 6: baja de marketing

`POST /_mock/reply {"phone": C, "button": "Dejar de recibir"}`.

**Esperado:** C con `opted_out = true` y `unread_count = 1`, en el listado y en
el detalle (`GET /conversations/{C}/messages`).

## Paso 7: varios mensajes suman

A responde dos veces mas.

**Esperado:** A con `unread_count = 2`.

## Paso 8: etiqueta de orden

Una conversacion con `order_number` de una orden existente del negocio.

**Esperado:** `order_id` poblado con el id de esa orden, en listado y detalle.

## Paso 9: la migracion es idempotente

Volver a correr `migrateWhatsappConversationReads`.

**Esperado:** no vuelve a rellenar: el conteo de filas y los `last_read_at` de
`whatsapp_conversation_reads` no cambian.

## No cubierto

- Cliente con `accepts_marketing = false` sin haber respondido la baja: la API
  de clientes no expone ese campo y las pruebas no escriben la base. Hoy nadie
  lo escribe, asi que ese camino solo aplica a datos cargados a mano.
