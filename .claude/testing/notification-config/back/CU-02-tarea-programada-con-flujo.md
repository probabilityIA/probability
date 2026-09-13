# CU-02: Tarea programada que dispara un flujo

Comprueba el orden que sigue el cliente: **plantillas -> flujo -> tarea
programada**. La tarea ya no apunta a una plantilla suelta; apunta a un flujo y
la conversacion se encadena sola.

## Precondiciones

Las mismas de [CU-01](CU-01-flujo-campana-contra-mock.md): base local, mock de
WhatsApp en 9103, credenciales de plataforma apuntando al mock (se escriben
**despues** de arrancar el backend) y las 4 plantillas del flujo 1 aprobadas.

La columna `flow_id` de `scheduled_notification_rules` la crea la migracion
`migrateWhatsappFlows`.

## Paso 1: crear la tarea con el flujo

```bash
curl -s -X POST localhost:3050/api/v1/scheduled-notification-rules \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{
    "notification_type_id": 2,
    "flow_id": 1,
    "name": "E2E tarea programada con flujo",
    "segment_type": "customers_inactive",
    "days_without_purchase": 1,
    "send_window_start": "00:00",
    "send_window_end": "23:59",
    "cooldown_days": 1,
    "daily_send_cap": 5,
    "batch_size_cap": 1,
    "frequency_minutes": 1440,
    "requires_opt_in": true,
    "enabled": true
  }'
```

**Esperado:** la regla guarda `FlowID = 1` y ademas `WhatsappTemplateID = 30`,
que es la plantilla inicial del flujo resuelta por el backend. El cliente nunca
elige esa plantilla: sale del flujo.

**Sin `flow_id` el alta falla** con "la tarea programada necesita un flujo: crea
primero las plantillas, despues el flujo, y recien ahi la tarea".

## Paso 2: ejecutarla

```bash
curl -s -X POST localhost:3050/api/v1/scheduled-notification-rules/<id>/run \
  -H "Authorization: Bearer $TOKEN"
```

`run` ignora la ventana horaria a proposito: es la ejecucion manual. El
despachador automatico si la respeta.

**Esperado:** `MatchedCount` y `QueuedCount` en 1 (con un cliente que cumpla el
segmento).

## Paso 3: verificar la cadena

```bash
curl -s localhost:9103/_mock/sent
```

**Esperado**, con el guion por defecto:

| # | plantilla | respuesta simulada |
|---|---|---|
| 1 | 26_saludo_inicial | Si |
| 2 | 26_dice_que_no | eres una rata |
| 3 | 26_hola | Jodete |
| 4 | 26_me_jodo | genial (corta) |

O sea: la tarea programada solo dispara el primer mensaje; el resto lo encadenan
los botones, igual que en una campana.

```sql
SELECT id, name, whatsapp_template_id, flow_id FROM scheduled_notification_rules
WHERE business_id = 26;

SELECT rule_id, status, phone, message_id FROM scheduled_notification_sends
WHERE rule_id = <id>;
```

## Compatibilidad

Las reglas viejas que solo tienen `whatsapp_template_id` y `flow_id` nulo siguen
funcionando: `resolveTemplate` usa el flujo cuando esta y cae a la plantilla
cuando no. Lo que ya no se puede es **crear** una tarea sin flujo.
