# CU-01: Campana programada con flujo, contra el mock de WhatsApp

Recorre un flujo completo de plantillas encadenadas **sin tocar Meta**. Todo el
trafico saliente lo recibe el mock de `back/testing/integrations/whatsappmock`, y
las respuestas del cliente las dispara el mismo mock como webhooks firmados.

## Precondiciones

1. Base local (`./scripts/dev-db-switch.sh status` debe decir LOCAL) con el
   negocio Demo (26).
2. Infra arriba: postgres 5434, redis, rabbitmq. Backend en 3050.
3. El flujo 1 del negocio 26 existe con estas transiciones:

   | origen | boton | destino |
   |---|---|---|
   | 26_saludo_inicial | Si | 26_dice_que_no |
   | 26_saludo_inicial | No | 26_me_jodo |
   | 26_dice_que_no | eres una rata | 26_hola |
   | 26_hola | Jodete | 26_me_jodo |

## Montaje del mock

```bash
cd back/testing
WHATSAPP_MOCK_PORT=9103 \
WEBHOOK_BASE_URL=http://localhost:3050 \
WHATSAPP_WEBHOOK_SECRET=<webhook_secret de integration:platform_creds:2> \
go run cmd/whatsappmock/main.go
```

Apuntar el backend al mock. Son dos cosas, y las dos son **solo locales**:

```bash
# 1. las credenciales de plataforma apuntan al mock en vez de a Meta
#    OJO: el backend repuebla esta clave al arrancar, asi que se escribe
#    DESPUES de levantar el backend
docker exec <redis> redis-cli -a <pass> GET integration:platform_creds:2
#    ... cambiar whatsapp_url a http://localhost:9103 y volver a hacer SET

# 2. la integracion de WhatsApp del negocio necesita numero propio
UPDATE integrations
SET config = config || '{"waba_id":"1302830408357767","phone_number_id":"1077369948787698"}'::jsonb
WHERE id = 61;
```

Sin (2) el lanzamiento falla con "las campanas de marketing solo salen desde el
numero propio del negocio", y los webhooks entrantes no resuelven a que negocio
pertenecen.

## Paso 1: aprobar las plantillas por el mock

Las plantillas nacen como borrador y una campana con flujo exige que **todas**
las plantillas destino esten aprobadas.

```bash
for id in 30 33 31 32; do
  curl -s -X POST localhost:3050/api/v1/whatsapp-templates/$id/submit \
    -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{}'
done
```

El mock las crea, responde un `meta_template_id` y a los 3 segundos manda el
webhook `message_template_status_update` con `APPROVED`.

**Esperado:** las 4 plantillas quedan en `approved` con `waba_id` poblado.

```sql
SELECT id, name, status FROM whatsapp_templates
WHERE business_id = 26 AND deleted_at IS NULL ORDER BY id;
```

## Paso 2: crear la campana programada con el flujo

```bash
curl -s -X POST localhost:3050/api/v1/whatsapp-campaigns \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{
    "name": "E2E flujo 1 mock",
    "flow_id": 1,
    "audience_type": "filtered_clients",
    "client_ids": [548378],
    "timezone": "America/Bogota",
    "send_window_start": "00:00",
    "send_window_end": "23:59",
    "scheduled_at": "<ahora + 45s en ISO UTC>",
    "daily_send_cap": 50,
    "batch_size": 10
  }'
```

**Ojo con la ventana de envio.** Si la hora de la prueba cae fuera de
`send_window_start`-`send_window_end`, el despachador salta la campana y no
manda nada. No es un bug: hay que abrir la ventana o probar en horario.

**Esperado:** `status=draft`, `FlowID=1` y `WhatsappTemplateID` resuelto solo a
la plantilla raiz del flujo (30).

## Paso 3: lanzar

```bash
curl -s -X POST localhost:3050/api/v1/whatsapp-campaigns/<id>/launch \
  -H "Authorization: Bearer $TOKEN"
```

**Esperado:** `status=scheduled` y `AudienceCount=1`. El despachador corre cada
2 minutos, asi que el primer envio sale en ese plazo.

## Paso 4: verificar la cadena

```bash
curl -s localhost:9103/_mock/sent
```

**Esperado con el guion por defecto** (rama izquierda del diagrama):

| # | plantilla enviada | respuesta del cliente |
|---|---|---|
| 1 | 26_saludo_inicial | Si |
| 2 | 26_dice_que_no | eres una rata |
| 3 | 26_hola | Jodete |
| 4 | 26_me_jodo | genial |

La cadena termina en el cuarto mensaje porque `26_me_jodo` no tiene transicion
para "genial": es el "sin respuesta" del diagrama.

En base de datos:

```sql
SELECT status, message_id FROM whatsapp_campaign_sends WHERE campaign_id = <id>;

SELECT l.direction, l.template_name, l.content
FROM whatsapp_message_logs l
JOIN whatsapp_conversations c ON c.id = l.conversation_id
WHERE c.business_id = 26 AND c.phone_number LIKE '%3164489436%'
ORDER BY l.created_at;
```

El envio de la campana queda `sent` con el `wamid` que devolvio el mock, y cada
plantilla encadenada deja su propia fila `outbound` en el log. Ese log es lo que
permite encadenar: el backend resuelve la plantilla de origen buscando
`whatsapp_message_logs.message_id = context.id` del webhook entrante.

## Paso 5: la otra rama

```bash
curl -s -X POST localhost:9103/_mock/reset
curl -s -X POST localhost:9103/_mock/script -H 'Content-Type: application/json' \
  -d '{"script":{"26_saludo_inicial":"No","26_me_jodo":""}}'
```

Repetir los pasos 2 y 3 con una campana nueva.

**Esperado:** solo dos envios, `26_saludo_inicial` y `26_me_jodo`, y la cadena
corta ahi porque la respuesta programada de `26_me_jodo` es vacia.

## Limpieza

El `whatsapp_url` que apunta al mock vive solo en Redis: **al reiniciar el
backend vuelve solo** al valor real, porque se repuebla desde
`integration_types.platform_credentials_encrypted`. No hay nada que revertir a
mano salvo el `config` de la integracion 61, que es dato de prueba del negocio
Demo en la base local.
