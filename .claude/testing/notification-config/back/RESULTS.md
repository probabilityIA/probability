# Resultados E2E - notification-config (back)

Entorno: base local `127.0.0.1:5434` (negocio Demo 26), backend `:3050`,
mock de WhatsApp `:9103`. **Ningun mensaje salio a Meta.**

## 2026-09-12

| Caso | Resultado | Nota |
|---|---|---|
| CU-01 rama "Si" (cadena completa) | OK | 4 envios encadenados |
| CU-01 rama "No" (corta) | OK | 2 envios, termina donde debe |
| Aprobacion de plantillas por mock | OK | tras corregir un bug, ver abajo |
| CU-02 tarea programada con flujo | OK | regla 2, flujo 1, cadena de 4 |

### Rama "Si" - campana 2, cliente 548378 (+573164489436)

| # | plantilla | respuesta simulada |
|---|---|---|
| 1 | 26_saludo_inicial | Si |
| 2 | 26_dice_que_no | eres una rata |
| 3 | 26_hola | Jodete |
| 4 | 26_me_jodo | genial (sin transicion, corta) |

`whatsapp_message_logs` quedo con 8 filas alternando outbound/inbound. El envio
de campana quedo `sent` con el `wamid` del mock y la campana paso a `completed`
en el tick siguiente.

### Rama "No" - campana 3, cliente 548380 (+573187168251)

| # | plantilla | respuesta simulada |
|---|---|---|
| 1 | 26_saludo_inicial | No |
| 2 | 26_me_jodo | (guion vacio, corta) |

Confirma que el flujo bifurca por el texto del boton y que una rama sin
respuesta programada termina limpio.

## Bugs encontrados

### 1. Un webhook de estado de plantilla se perdia si la cache estaba fria (CORREGIDO)

`usecasetemplates.HandleStatusUpdate` hacia `return nil` dentro de la rama de
error de `UpdateStatusByWABA`, asi que cuando la clave `whatsapp:waba:<id>` no
existia en Redis **nunca llegaba a publicar el estado** y la plantilla se quedaba
en `pending` para siempre en nuestra base.

Sintoma: el log decia `Actualizacion de estado de plantilla recibida event=APPROVED`
seguido de `no se pudo actualizar el estado de la plantilla en cache
error="key not found: whatsapp:waba:1302830408357767"`, y la fila no cambiaba.

Corregido en `status.go`: la cache es un modelo de lectura, su fallo ya no impide
persistir. Con el arreglo las 4 plantillas pasaron a `approved`.

**Esto no es solo del entorno de pruebas.** En produccion la cache se llena
cuando alguien consulta el estado de plantillas; si esta fria (Redis reiniciado,
despliegue nuevo) toda aprobacion que mande Meta en esa ventana se perdia.

### 2. Una campana nunca registra las respuestas (SIN CORREGIR)

`whatsapp_campaign_sends.status` solo recibe `sent` o `failed`: son los unicos
valores que publica `consumercampaign`. Nada escribe `delivered`, `read` ni
`replied`, aunque la tabla y el contador existen.

Consecuencia: `replied_count` de la campana **siempre queda en 0**, incluso
cuando el cliente contesto (en esta prueba contesto 4 veces y el contador siguio
en 0). Lo mismo para cualquier metrica de entrega a nivel de campana.

El dato sin embargo si esta: los webhooks de estado actualizan
`whatsapp_message_logs`. Falta correlacionar ese log con el envio de campana por
`message_id`.

## Cosas que no son bugs pero muerden

- **Ventana de envio**: si la prueba corre fuera de `send_window_start`-`end`
  (por defecto 09:00-19:00) el despachador salta la campana en silencio.
- **Los contadores van un tick atras**: `UpdateCampaignCounters` solo corre al
  despachar un lote o al completar, cada 2 minutos.
- **Las tareas programadas ya aceptan flujos** (corregido el 2026-09-12, ver
  CU-02). Antes solo apuntaban a una plantilla suelta. Ahora el orden es el
  mismo en los dos caminos: plantillas -> flujo -> campana o tarea programada.
  Las reglas viejas con `flow_id` nulo siguen andando por compatibilidad, pero
  no se pueden crear nuevas sin flujo.
- **`integration:platform_creds:2` se repuebla al arrancar el backend** desde
  `integration_types.platform_credentials_encrypted`. Apuntar el mock escribiendo
  esa clave hay que hacerlo *despues* de levantar el backend, o se pierde.

## 2026-09-13 - Programacion por dias y modo de entrega (rama feat/campanas-programacion-dias)

Contra base local y el mock de WhatsApp, backend local.

| Caso | Resultado |
|---|---|
| Repetir sin cantidad de veces | OK: 400 "para repetir el envio indica cuantas veces se repite" |
| Fecha mal escrita (`15/09/2026`) | OK: rechazada |
| Intervalo en 0 | OK: rechazado |
| Lanzar con todas las fechas pasadas (campana 4) | OK: rechazado al lanzar |
| Solo manana (campana 5) | OK: lanza, envio queda `pending` hoy, no sale nada |
| Repetir hoy + en 7 dias (campana 6) | OK: vuelta 1 sale al mock (`round=1`, `sent`, con wamid) |

No probado en E2E (cubierto solo por tests unitarios): cambio de vuelta al
llegar la segunda fecha, pausa al acabarse las fechas con pendientes, y
completar tras la ultima fecha. Requieren esperar dias o escribir la base.

## 2026-09-13 - CU-03 bandeja: sin leer, orden, baja y etiquetas

Base local, backend local y mock de WhatsApp (respuestas automaticas apagadas).
Campana 7 "CU-03 bandeja sin leer" a tres clientes: A 573102020202, B 573105557788, C 573029998877.

| Paso | Resultado |
|---|---|
| 1. Tres conversaciones con etiqueta de campana | FAIL en la primera corrida, OK tras el fix (ver bug) |
| 2. A responde: `unread_count=1`, total 1 | OK |
| 3. B responde despues y se marca leido: A (sin leer, 19:49:50) queda antes que B (leido, 19:49:56) | OK |
| 4. Marcar A leido: 0 y persiste al reconsultar | OK |
| 5. `POST .../read` con otro negocio: 404, marca de A intacta, 0 filas en el negocio 1 | OK |
| 6. B responde "Dejar de recibir": `opted_out=true` (antes false), sin leer 1, tambien en el detalle | OK |
| 7. A responde dos veces: `unread_count=2` | OK |
| 8. Etiqueta de orden: `order_id` real de DEM-0045 en listado y detalle | OK |
| 9. Migracion re-ejecutada: 65 filas y misma `max(last_read_at)`, no re-rellena | OK |

**Bug encontrado y corregido:** los envios de campana guardan el telefono como
esta en el cliente (`3102020202`, sin indicativo) y la conversacion lo guarda con
`57`. El filtro `campaign_id` del listado y la etiqueta de campana comparaban
todos los digitos, asi que un cliente guardado sin +57 nunca aparecia en la
pestana Conversaciones de su campana. Ahora comparan los ultimos 10 digitos
(`message_audit_queries.go`). Anterior a este trabajo: tambien afectaba a la
pestana de conversaciones del detalle de campana.

**No cubierto:** cliente con `accepts_marketing=false` sin haber respondido la
baja (la API de clientes no expone el campo y la prueba no escribe la base).

**Observado, no es de este cambio:** cada respuesta del cliente crea una fila
nueva de `whatsapp_conversations` tipo `inbound` sin `campaign_id` para el mismo
telefono (A acumula tres filas en minutos). El listado las agrupa por telefono,
pero la tabla crece de mas.

## 2026-09-13 - Texto real y botones de las plantillas en Conversaciones

Base local, backend local, mock de WhatsApp.

| Caso | Resultado |
|---|---|
| Campana/flujo guardado como `nombre: parametros` muestra encabezado + cuerpo + pie | OK (`26_dice_que_no: Andres` -> "hola / dices que no ? Andres / que mal") |
| Plantilla del sistema guardada como "Plantilla: nombre" muestra el texto | OK tras llenar `body_text` desde Meta (`migrateSystemTemplateBodies`); 0 mensajes quedan sin texto en local |
| Vista previa de la lista usa el texto, no el nombre | OK |
| Botones de plantilla del negocio | OK (`26_me_jodo`: genial / que cagada) |
| Marketing agrega "Dejar de recibir" aunque no este guardado | OK (`26_saludo_inicial`: Si / No / Dejar de recibir) |
| Botones de plantillas del sistema (desde Meta) | OK (confirmaciones: Confirmar pedido / No confirmar) |
| Mensaje ya guardado con texto completo tambien muestra botones | OK |
| Mensajes entrantes no muestran botones | OK |

**Limite:** los envios hechos por `SendTemplateWithConversation` guardaban el
mensaje sin variables, asi que esos mensajes viejos muestran "-" donde iba el
dato (p. ej. `pedido_confirmado_v2`). Corregido hacia adelante: ahora se
guardan las variables.

## 2026-09-13 - CU-04 adjuntos en el chat y retencion de 1 ano

Base local, backend local, mock de WhatsApp. Archivos en el bucket real
`probability-chat-attachments` (privado, SSE-S3, solo TLS, expira a los 365 dias).

| Paso | Resultado |
|---|---|
| 1. Imagen con texto desde la bandeja | OK: 200, mock recibe `POST /media` + `type=image` con caption, detalle con `media.available=true`, URL firmada 200 `image/png`, objeto en `whatsapp/26/2026/09/` |
| 2. PDF sin texto | OK: `type=document` con `filename=factura.pdf`, `content` vacio |
| 3a. `.exe` | OK: 400 `invalid_media`, no llega a S3 ni al mock |
| 3b. Imagen de mas de 5 MB | OK: 400 |
| 3c. Archivo de mas de 10 MB | OK: 413 |
| 3d. Usuario normal con `business_id=1` en el form | OK: se ignora, el archivo queda en `whatsapp/26/`, nada en `whatsapp/1/` |
| 4. Cliente envia imagen con texto | OK: se descarga del mock, se guarda en S3, URL firmada 200, sin leer sube, vista previa "Asi llego" |
| 5. Cliente envia documento sin texto | OK: `filename=comprobante.pdf`, vista previa "[Documento]" |
| 6. Retencion | OK: `EXPLAIN` de los 3 DELETE planifica; hoy borraria 0 de 105 mensajes (con corte de 30 dias el filtro encuentra 51); el job corrio 2 min despues del arranque con `cutoff=2025-09-13` y borro 0/0/0; test unitario del corte a 365 dias |

**Encontrado de paso (anterior a este trabajo):** los archivos que mandaban los
clientes nunca se procesaban: el webhook no tenia campos de imagen, documento,
audio ni video, asi que el mensaje quedaba vacio. Corregido en este mismo cambio.

**Pendiente de limpieza:** 5 objetos de prueba en
`s3://probability-chat-attachments/whatsapp/26/2026/09/` (se dejaron para las
capturas; expiran solos a los 365 dias si no se borran).
