# WhatsApp Business Cloud API Integration

Integración completa y bidireccional con WhatsApp Business Cloud API para gestión de conversaciones con clientes.

## Características

- ✅ **11 plantillas conversacionales** con flujo de confirmación de pedidos
- ✅ **Botones interactivos** (Quick Reply) para respuestas predefinidas
- ✅ **Webhooks bidireccionales** para recibir respuestas de usuarios en tiempo real
- ✅ **Persistencia de conversaciones** en PostgreSQL con ventana de 24h
- ✅ **Máquina de estados** para flujos conversacionales automatizados
- ✅ **Eventos de negocio** publicados en RabbitMQ para integración con otros servicios
- ✅ **Arquitectura hexagonal** con separación de capas (domain, app, infra)

## Arquitectura

### Estructura de Directorios

```
whatsApp/
├── bundle.go                           # Ensamblaje de componentes (DI)
├── internal/
│   ├── domain/                         # Núcleo - Reglas de negocio
│   │   ├── entities/                   # Entidades del dominio (100% puras)
│   │   │   ├── conversation.go         # Entidades de conversación y estados
│   │   │   ├── message.go              # Entidades de mensajes
│   │   │   └── template.go             # Catálogo de 11 plantillas
│   │   ├── dtos/                       # DTOs de dominio
│   │   │   └── send_message_request.go # DTOs de entrada
│   │   ├── ports/                      # Interfaces (contratos)
│   │   │   └── ports.go                # Repositorios, servicios externos
│   │   └── errors/                     # Errores del dominio
│   │       └── errors.go
│   ├── app/                            # Casos de uso
│   │   ├── usecasemessaging/
│   │   │   ├── constructor.go          # Interfaz y constructor
│   │   │   ├── send-template-message.go    # Envío de plantillas dinámicas
│   │   │   ├── handle-webhook.go           # Procesamiento de webhooks
│   │   │   ├── conversation-manager.go     # Máquina de estados conversacional
│   │   │   ├── send-message.go             # Caso de uso legacy
│   │   │   └── utils.go                    # Validación de teléfonos
│   │   └── usecasetestconnection/
│   │       └── test_connection.go      # Test de integración
│   └── infra/                          # Infraestructura
│       ├── primary/                    # Adaptadores de entrada
│       │   ├── handlers/
│       │   │   ├── constructor.go      # Interfaz IHandler
│       │   │   ├── routes.go           # Registro de rutas HTTP
│       │   │   ├── template_handler.go # Endpoint POST /send-template
│       │   │   ├── webhook_handler.go  # Endpoints GET/POST /webhook
│       │   │   ├── request/            # DTOs de entrada HTTP (con tags json)
│       │   │   │   ├── send_template.go
│       │   │   │   └── webhook_payload.go  # ✨ Webhooks de Meta (con tags)
│       │   │   └── response/           # DTOs de salida HTTP (con tags json)
│       │   │       └── send_template.go
│       │   ├── consumer/               # Consumidores Redis
│       │   │   └── consumerevent/
│       │   └── queue/                  # Consumidores RabbitMQ
│       │       └── consumerorder/
│       └── secondary/                  # Adaptadores de salida
│           ├── client/                 # Cliente HTTP WhatsApp API
│           ├── repository/             # Repositorios PostgreSQL
│           │   ├── constructor.go
│           │   ├── conversation_repository.go
│           │   ├── message_log_repository.go
│           │   └── mappers/            # Mappers domain ↔ models
│           │       ├── to_domain.go
│           │       └── to_model.go
│           ├── adapters/               # Adaptadores a otros módulos
│           └── queue/                  # Publisher RabbitMQ
│               └── webhook_publisher.go
```

### Flujo de Conversación

```
[INICIO] → confirmacion_pedido_contraentrega
            ├─ "Confirmar pedido" → pedido_confirmado [FIN]
            └─ "No confirmar" → menu_no_confirmacion
                                 ├─ "Presentar novedad" → tipo_novedad_pedido
                                 │                          ├─ "Cambio de dirección" → novedad_cambio_direccion [FIN]
                                 │                          ├─ "Cambio de productos" → novedad_cambio_productos [FIN]
                                 │                          └─ "Cambio medio de pago" → novedad_cambio_medio_pago [FIN]
                                 ├─ "Cancelar pedido" → confirmar_cancelacion_pedido
                                 │                       ├─ "Sí, cancelar" → motivo_cancelacion_pedido
                                 │                       │                    └─ [texto libre] → pedido_cancelado [FIN]
                                 │                       └─ "No, volver" → menu_no_confirmacion
                                 └─ "Asesor" → handoff_asesor [FIN - HUMANO]
```

## Setup Local con ngrok

### Prerrequisitos

1. **Cuenta de WhatsApp Business en Meta**
   - Registrarse en https://business.facebook.com
   - Crear app de WhatsApp Business
   - Obtener credenciales: Phone Number ID, Access Token

2. **ngrok instalado** (para exponer localhost)
   ```bash
   # Opción 1: ngrok (más común)
   curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
   echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
   sudo apt update && sudo apt install ngrok

   # Registrarse en ngrok.com y obtener authtoken
   ngrok config add-authtoken <TU_TOKEN>

   # Opción 2: cloudflared (URL permanente, recomendado)
   wget https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
   sudo dpkg -i cloudflared-linux-amd64.deb
   ```

### Configuración

#### 1. Variables de Entorno

Agregar al archivo `/back/central/.env`:

```env
# WhatsApp Cloud API - Credenciales
WHATSAPP_PHONE_NUMBER_ID=123456789012345  # ID del número de teléfono de WhatsApp Business
WHATSAPP_ACCESS_TOKEN=EAAxxxxxxxxxxxxxxx  # Access token de la app de Meta

# WhatsApp Webhooks - Desarrollo Local
WHATSAPP_VERIFY_TOKEN=mi_token_seguro_para_verificacion_dev_2026
WHATSAPP_WEBHOOK_SECRET=mi_secreto_hmac_para_firmas_sha256

# Webhook URL pública (se actualiza según el túnel)
# DESARROLLO LOCAL: Usar URL de ngrok/cloudflared
WHATSAPP_WEBHOOK_URL=https://1234-tu-ip.ngrok-free.app/api/integrations/whatsapp/webhook

# PRODUCCIÓN: Usar dominio real
# WHATSAPP_WEBHOOK_URL=https://api.probability.com/api/integrations/whatsapp/webhook
```

#### 2. Migraciones de Base de Datos

Ejecutar las migraciones SQL ubicadas en `/back/migration/shared/sql/`:

```bash
# Desde el directorio del proyecto
psql -h localhost -p 5433 -U postgres -d probability < back/migration/shared/sql/whatsapp_conversations.sql
psql -h localhost -p 5433 -U postgres -d probability < back/migration/shared/sql/whatsapp_message_logs.sql
```

O ejecutar el script de migración si existe en el proyecto.

#### 3. Iniciar Backend

```bash
cd back/central
go run cmd/main.go
```

El servidor debería iniciar en `http://localhost:8080`

#### 4. Iniciar Túnel ngrok

En una terminal separada:

```bash
# Opción A: ngrok (URL cambia cada reinicio)
ngrok http 8080
# Copiar la URL generada (ej: https://1234-5678-90ab.ngrok-free.app)

# Opción B: cloudflared (URL permanente)
cloudflared tunnel --url http://localhost:8080
# Copiar la URL generada (ej: https://random-words.trycloudflare.com)
```

⚠️ **IMPORTANTE**: Con ngrok gratuito, la URL cambia cada vez que reinicias. Debes actualizar:
1. La variable `WHATSAPP_WEBHOOK_URL` en `.env`
2. La URL del webhook en Meta Business Manager

#### 5. Configurar Webhook en Meta Business Manager

1. Ir a https://business.facebook.com/settings/whatsapp-business-accounts
2. Seleccionar tu cuenta WhatsApp Business
3. Ir a **Configuration** → **Webhooks**
4. Click **Edit** o **Configure**
5. Ingresar:
   - **Callback URL**: `https://TU-URL-NGROK.ngrok-free.app/api/integrations/whatsapp/webhook`
   - **Verify Token**: `mi_token_seguro_para_verificacion_dev_2026` (debe coincidir con tu `.env`)
6. Marcar campos de suscripción:
   - ☑️ `messages` (OBLIGATORIO)
   - ☑️ `message_template_status_update` (OPCIONAL)
7. Click **Verify and Save**
8. Meta enviará GET request → Tu endpoint debe retornar el challenge → ✅ Verified

### Verificación

#### 1. Verificar Webhook

```bash
curl "http://localhost:8080/api/integrations/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=mi_token_seguro_para_verificacion_dev_2026&hub.challenge=test123"
# Debe retornar: test123
```

#### 2. Enviar Plantilla de Prueba

```bash
curl -X POST http://localhost:8080/api/integrations/whatsapp/send-template \
  -H "Content-Type: application/json" \
  -d '{
    "template_name": "confirmacion_pedido_contraentrega",
    "phone_number": "+573001234567",
    "order_number": "ORD-12345",
    "business_id": 1,
    "variables": {
      "1": "Juan Pérez",
      "2": "TiendaDemo",
      "3": "ORD-12345",
      "4": "Calle 123 #45-67",
      "5": "2x Camiseta Azul, 1x Pantalón Negro"
    }
  }'
```

**Respuesta esperada:**
```json
{
  "message_id": "wamid.HBgNNTczMDA...",
  "status": "sent"
}
```

#### 3. Verificar en Base de Datos

```sql
-- Ver conversaciones activas
SELECT * FROM whatsapp_conversations WHERE expires_at > NOW();

-- Ver logs de mensajes
SELECT * FROM whatsapp_message_logs ORDER BY created_at DESC LIMIT 10;
```

#### 4. Prueba End-to-End

1. Enviar plantilla inicial → Cliente recibe WhatsApp con 2 botones
2. Cliente presiona "No confirmar" → Recibe menú con 3 opciones
3. Cliente presiona "Cancelar pedido" → Recibe confirmación
4. Cliente presiona "Sí, cancelar" → Recibe solicitud de motivo
5. Cliente envía texto libre "Ya no lo necesito" → Recibe confirmación de cancelación
6. Verificar:
   - Conversación en BD: `current_state = 'COMPLETED'`
   - Logs de mensajes: 5 registros (3 outbound, 2 inbound)
   - Evento publicado en RabbitMQ: `orders.whatsapp.cancelled`

## Endpoints HTTP

### POST `/api/integrations/whatsapp/send-template`

Envía una plantilla de WhatsApp con variables dinámicas.

**Request:**
```json
{
  "template_name": "confirmacion_pedido_contraentrega",
  "phone_number": "+573001234567",
  "order_number": "ORD-12345",
  "business_id": 1,
  "variables": {
    "1": "Juan Pérez",
    "2": "TiendaDemo",
    "3": "ORD-12345",
    "4": "Calle 123 #45-67",
    "5": "2x Camiseta Azul, 1x Pantalón Negro"
  }
}
```

**Response:**
```json
{
  "message_id": "wamid.HBgNNTczMDA...",
  "status": "sent"
}
```

### GET `/api/integrations/whatsapp/webhook`

Verificación del webhook (usado por Meta).

**Query Params:**
- `hub.mode=subscribe`
- `hub.verify_token=<token>`
- `hub.challenge=<challenge>`

**Response:** Retorna el challenge si el token es válido.

### POST `/api/integrations/whatsapp/webhook`

Recibe eventos de WhatsApp (mensajes, estados).

**Headers:**
- `X-Hub-Signature-256: sha256=<hmac_signature>`

**Body:** Webhook payload de Meta (ver estructura en `internal/infra/primary/handlers/request/webhook_payload.go`)

**Response:**
```json
{
  "status": "received"
}
```

## Eventos de RabbitMQ

### Colas

| Cola | Evento | Cuándo se publica |
|------|--------|-------------------|
| `orders.whatsapp.confirmed` | Pedido confirmado | Usuario presiona "Confirmar pedido" |
| `orders.whatsapp.cancelled` | Pedido cancelado | Usuario completa flujo de cancelación |
| `orders.whatsapp.novelty` | Novedad solicitada | Usuario selecciona tipo de novedad |
| `customer.whatsapp.handoff` | Handoff a humano | Usuario presiona "Asesor" |

### Formato de Evento

```json
{
  "event_type": "order.confirmed",
  "order_number": "ORD-12345",
  "phone_number": "+573001234567",
  "business_id": 1,
  "source": "whatsapp",
  "timestamp": 1706284800
}
```

## Plantillas Disponibles

| Nombre | Variables | Botones |
|--------|-----------|---------|
| `confirmacion_pedido_contraentrega` | {{1}}=Nombre, {{2}}=Tienda, {{3}}=#Orden, {{4}}=Dirección, {{5}}=Productos | "Confirmar pedido", "No confirmar" |
| `pedido_confirmado` | {{1}}=#Pedido | Ninguno |
| `menu_no_confirmacion` | {{1}}=#Pedido | "Presentar novedad", "Cancelar pedido", "Asesor" |
| `tipo_novedad_pedido` | Ninguna | "Cambio de dirección", "Cambio de productos", "Cambio medio de pago" |
| `confirmar_cancelacion_pedido` | {{1}}=#Pedido | "Sí, cancelar", "No, volver" |
| `motivo_cancelacion_pedido` | Ninguna | Ninguno (espera texto libre) |
| `pedido_cancelado` | {{1}}=#Pedido | Ninguno |
| `novedad_cambio_direccion` | Ninguna | Ninguno |
| `novedad_cambio_productos` | Ninguna | Ninguno |
| `novedad_cambio_medio_pago` | Ninguna | Ninguno |
| `handoff_asesor` | Ninguna | Ninguno |

## Troubleshooting

### Webhook no se verifica en Meta

- Verificar que ngrok esté corriendo y la URL sea correcta
- Verificar que `WHATSAPP_VERIFY_TOKEN` en `.env` coincida con el configurado en Meta
- Verificar logs del backend para ver el GET request

### Mensajes no llegan al usuario

- Verificar que `WHATSAPP_PHONE_NUMBER_ID` y `WHATSAPP_ACCESS_TOKEN` sean correctos
- Verificar que el número de teléfono esté en formato internacional (+573001234567)
- Verificar que el número esté registrado como destinatario de prueba en Meta (modo development)

### Webhooks no se reciben

- Verificar que ngrok esté corriendo
- Verificar que la URL en Meta coincida con la URL de ngrok
- Verificar logs de ngrok para ver si llegan requests
- Verificar que `WHATSAPP_WEBHOOK_SECRET` esté configurado correctamente

### Conversación no avanza

- Verificar que el texto del botón presionado coincida exactamente con el esperado
- Verificar en BD que `current_state` sea el esperado
- Verificar logs del backend para ver transiciones de estado

### Errores de compilación

Si hay errores al compilar, verificar:
```bash
cd back/central
go mod tidy
go build ./...
```

## Limitaciones Conocidas

- **ngrok gratuito**: URL cambia en cada reinicio
- **Ventana de 24h**: Las conversaciones expiran después de 24h sin actividad
- **Plantillas aprobadas**: Las plantillas deben estar aprobadas en Meta antes de enviarlas
- **Rate limits**: WhatsApp Cloud API tiene límites de mensajes por segundo/día

## Próximos Pasos

- [ ] Soporte para mensajes multimedia (imágenes, videos, documentos)
- [ ] Plantillas con headers dinámicos (imágenes, videos)
- [ ] List messages para menús más complejos
- [ ] Flow buttons para formularios in-app
- [ ] Analytics de conversiones por plantilla
- [ ] A/B testing de variantes de plantillas
- [ ] Integración con CRM para handoff automático

## Referencias

- [WhatsApp Cloud API Documentation](https://developers.facebook.com/docs/whatsapp/cloud-api)
- [Message Templates](https://developers.facebook.com/docs/whatsapp/cloud-api/guides/send-message-templates)
- [Webhooks Guide](https://developers.facebook.com/docs/whatsapp/cloud-api/webhooks)
- [ngrok Documentation](https://ngrok.com/docs)
- [cloudflared Documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/)
