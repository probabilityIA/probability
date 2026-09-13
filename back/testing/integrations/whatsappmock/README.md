# Mock de la Graph API de WhatsApp

Servidor HTTP que suplanta a `graph.facebook.com/v22.0` para probar el modulo de
notificaciones **sin tocar Meta**. Es el complemento del simulador interactivo de
`integrations/whatsapp`, que solo sabe mandar webhooks entrantes: este ademas
recibe lo que el backend envia.

## Que resuelve

| Camino | Antes | Ahora |
|---|---|---|
| Envio de plantilla | pegaba a Meta real | lo recibe el mock y devuelve un `wamid` |
| Alta de plantilla | quedaba `pending` esperando a Meta | el mock la aprueba por webhook |
| Respuesta del cliente | habia que mandarla a mano | el mock la dispara segun un guion |

## Como arranca

```bash
export WHATSAPP_WEBHOOK_SECRET=<el webhook_secret de integration:platform_creds:2>
export WEBHOOK_BASE_URL=http://localhost:3050
export WHATSAPP_MOCK_PORT=9103
cd back/testing && go run cmd/whatsappmock/main.go
```

Para que el backend le hable al mock en vez de a Meta, hay que apuntar
`whatsapp_url` de las credenciales de plataforma a `http://localhost:9103`.
Eso es **solo en local**: nunca se toca produccion.

## Endpoints que imita de Meta

| Metodo | Ruta | Para que |
|---|---|---|
| POST | `/:phone_number_id/messages` | enviar plantilla, devuelve `wamid` |
| POST | `/:waba_id/message_templates` | crear plantilla, luego la aprueba |
| GET | `/:waba_id/message_templates` | listar lo creado |
| POST | `/:app_id/uploads` | abrir sesion de subida de imagen |
| POST | `/:session_id` | cerrar la subida y devolver el handle |

## Endpoints de control del mock

| Metodo | Ruta | Para que |
|---|---|---|
| GET | `/_mock/sent` | todo lo que el backend envio, con `wamid` y parametros |
| GET | `/_mock/templates` | plantillas creadas y su estado |
| GET | `/_mock/script` | guion de respuestas vigente |
| POST | `/_mock/script` | cambiar el guion o apagar la respuesta automatica |
| POST | `/_mock/reply` | disparar una respuesta a mano |
| POST | `/_mock/reset` | vaciar el estado en memoria |

## El guion de respuestas

Cuando el backend manda una plantilla, el mock busca el nombre en el guion. Si
hay una respuesta, espera un momento y dispara tres webhooks firmados:
`delivered`, `read` y el mensaje de boton con `context.id` igual al `wamid` del
envio. Ese `context.id` es lo que le permite al backend saber a que plantilla
esta contestando el cliente y encadenar la siguiente.

Guion por defecto, el del flujo 1 del negocio Demo:

```json
{
  "26_saludo_inicial": "Si",
  "26_dice_que_no": "eres una rata",
  "26_hola": "Jodete",
  "26_me_jodo": "genial"
}
```

Para recorrer la otra rama:

```bash
curl -X POST localhost:9103/_mock/script -H 'Content-Type: application/json' \
  -d '{"script":{"26_saludo_inicial":"No"}}'
```

Una plantilla con respuesta `""` corta la rama, que es como se prueba el
"sin respuesta" del diagrama.

## Respuesta a mano

```bash
curl -X POST localhost:9103/_mock/reply -H 'Content-Type: application/json' \
  -d '{"phone":"+573164489436","button":"que cagada"}'
```

Contesta al ultimo mensaje que se le envio a ese numero.
