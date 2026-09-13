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
