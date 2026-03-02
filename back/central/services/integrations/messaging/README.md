# Modulo de Mensajeria (messaging)

## Descripcion

Modulo de integracion de mensajeria del proyecto Probability. Actualmente implementa el proveedor **WhatsApp** (type_id=2) para comunicacion automatizada con clientes via la WhatsApp Cloud API de Meta. Incluye soporte para SMS como proveedor futuro (pendiente de implementacion).

El modulo se registra en el sistema de integraciones centralizado (`integrationCore`) y opera como un adaptador de salida para notificaciones y un flujo conversacional interactivo con los clientes.

## Proposito

- Enviar notificaciones automatizadas a clientes mediante plantillas de WhatsApp aprobadas por Meta.
- Gestionar un flujo conversacional con maquina de estados para la confirmacion de pedidos contra entrega.
- Recibir y procesar webhooks de Meta (mensajes entrantes, estados de entrega, respuestas de botones).
- Enviar alertas de monitoreo del servidor al administrador via WhatsApp.
- Publicar eventos de negocio (confirmacion, cancelacion, novedades, handoff) a RabbitMQ para que otros modulos los procesen.

## Entidades y Modelos de Dominio

### Conversation

Representa una conversacion activa de WhatsApp con ventana de 24 horas. Contiene el numero de telefono, numero de orden, business_id, estado actual de la maquina de estados y metadata asociada.

**Estados de la conversacion (ConversationState):**

```
START -> AWAITING_CONFIRMATION -> COMPLETED (pedido confirmado)
                               -> AWAITING_MENU_SELECTION -> AWAITING_NOVELTY_TYPE -> COMPLETED
                                                          -> AWAITING_CANCEL_CONFIRM -> AWAITING_CANCEL_REASON -> COMPLETED
                                                          -> HANDOFF_TO_HUMAN
```

### MessageLog

Registro de cada mensaje enviado o recibido dentro de una conversacion. Incluye direction (inbound/outbound), message_id de WhatsApp, nombre de plantilla, contenido y estado (sent, delivered, read, failed).

### TemplateMessage / TemplateDefinition

Representacion del mensaje de plantilla a enviar a la API de WhatsApp y catalogo de las 11 plantillas aprobadas en Meta:

| Plantilla | Variables | Botones | Descripcion |
|-----------|-----------|---------|-------------|
| confirmacion_pedido_contraentrega | nombre, tienda, numero_orden, direccion, productos | Confirmar pedido / No confirmar | Plantilla inicial de confirmacion |
| pedido_confirmado | numero_pedido | - | Confirmacion exitosa |
| menu_no_confirmacion | numero_pedido | Presentar novedad / Cancelar pedido / Asesor | Menu de opciones al no confirmar |
| tipo_novedad_pedido | - | Cambio de direccion / Cambio de productos / Cambio medio de pago | Seleccion de tipo de novedad |
| confirmar_cancelacion_pedido | numero_pedido | Si, cancelar / No, volver | Confirmacion antes de cancelar |
| motivo_cancelacion_pedido | - | - | Solicita motivo (texto libre) |
| pedido_cancelado | numero_pedido | - | Confirmacion de cancelacion |
| novedad_cambio_direccion | - | - | Recepcion de solicitud |
| novedad_cambio_productos | - | - | Recepcion de solicitud |
| novedad_cambio_medio_pago | - | - | Recepcion de solicitud |
| handoff_asesor | - | - | Espera para atencion humana |
| alerta_servidor | tipo_alerta, descripcion | - | Alerta de monitoreo del servidor |

## Endpoints HTTP

Registrados bajo el grupo `/integrations/whatsapp`:

| Metodo | Ruta | Handler | Descripcion |
|--------|------|---------|-------------|
| POST | `/whatsapp/send-template` | SendTemplate | Envia una plantilla de WhatsApp con variables dinamicas |
| GET | `/whatsapp/webhook` | VerifyWebhook | Verificacion del webhook por Meta (challenge/token) |
| POST | `/whatsapp/webhook` | ReceiveWebhook | Recibe eventos entrantes de mensajes y estados de entrega |

El endpoint POST `/webhook` valida la firma HMAC-SHA256 del payload, responde 200 inmediatamente (requisito de Meta: <5s) y procesa el webhook de forma asincrona en una goroutine.

## Consumidores RabbitMQ (entrada)

| Cola | Consumer | Descripcion |
|------|----------|-------------|
| `orders.confirmation.requested` | consumerorder | Recibe solicitudes de confirmacion de pedido y envia la plantilla de WhatsApp correspondiente al cliente |
| `monitoring.alerts` | consumeralert | Recibe alertas de monitoreo del servidor (RAM, CPU, disco) y las envia al telefono del administrador |

## Publicadores RabbitMQ (salida)

Los siguientes eventos se publican como resultado de las interacciones del usuario en el flujo conversacional:

| Cola | Evento | Descripcion |
|------|--------|-------------|
| `orders.whatsapp.confirmed` | order.confirmed | El cliente confirmo el pedido via WhatsApp |
| `orders.whatsapp.cancelled` | order.cancelled | El cliente cancelo el pedido (incluye motivo) |
| `orders.whatsapp.novelty` | order.novelty_requested | El cliente solicito una novedad (cambio de direccion, productos o medio de pago) |
| `customer.whatsapp.handoff` | customer.handoff_requested | El cliente solicito atencion de un agente humano |

## Integracion con Otros Modulos

- **integrations/core**: Se registra como integracion tipo WhatsApp (type_id=2) mediante `integrationCore.RegisterIntegration()`. Implementa `IIntegrationContract` incluyendo `RegisterRoutes`, `TestConnection` y `GetWebhookURL`.
- **modules (ModuleBundles)**: Recibe el bundle de modulos para acceso a funcionalidades compartidas.
- **services/events**: El modulo unificado de eventos consume eventos de ecommerce, consulta configuracion de notificaciones en Redis cache y encola mensajes a `orders.confirmation.requested`, que este modulo consume.
- **monitoring**: El modulo de monitoreo publica alertas a `monitoring.alerts`, que este modulo consume y envia al administrador.
- **orders**: Los eventos publicados por este modulo (`orders.whatsapp.confirmed`, `orders.whatsapp.cancelled`, `orders.whatsapp.novelty`) son consumidos por el modulo de ordenes para actualizar el estado de los pedidos.

## Arquitectura (Capas Hexagonales)

```
messaging/
  bundle.go                          # Bundle padre, registra proveedores en integrationCore
  whatsapp/
    bundle.go                        # Bundle de WhatsApp, wiring de dependencias
    internal/
      domain/                        # Nucleo de dominio (sin dependencias externas)
        entities/                    # Conversation, MessageLog, TemplateMessage, TemplateDefinition
        dtos/                        # WebhookPayloadDTO, StateTransitionDTO, SendMessageRequest, NotificationConfigData
        errors/                      # Errores tipados del dominio
        ports/                       # Interfaces: IWhatsApp, IConversationRepository, IMessageLogRepository,
                                     #   IIntegrationRepository, IEventPublisher, INotificationConfigRepository
      app/                           # Capa de aplicacion
        usecasemessaging/            # Caso de uso principal: envio de mensajes, manejo de webhooks,
                                     #   maquina de estados conversacional
        usecasetestconnection/       # Caso de uso de prueba de conexion
      infra/                         # Capa de infraestructura
        primary/                     # Adaptadores de entrada
          handlers/                  # Handlers HTTP (webhook, send-template)
            request/                 # Structs de request HTTP (con tags JSON)
            response/                # Structs de response HTTP
            mappers/                 # Mappers de infra a dominio
          queue/
            consumerorder/           # Consumer de confirmacion de ordenes
            consumeralert/           # Consumer de alertas de monitoreo
        secondary/                   # Adaptadores de salida
          client/                    # Cliente HTTP para WhatsApp Cloud API (Meta Graph API)
            request/                 # Request structs para la API
            response/                # Response structs de la API
            mappers/                 # Mappers de dominio a request de API
          repository/                # Repositorios: ConversationRepository, MessageLogRepository,
                                     #   IntegrationRepository (con descifrado AES-256-GCM)
            mappers/                 # Mappers de modelos GORM a dominio
          cache/                     # Cache Redis para configuraciones de notificacion
          queue/                     # WebhookPublisher (publica eventos de negocio a RabbitMQ)
      mocks/                         # Mocks para testing unitario
```

## Configuracion Requerida

Variables de entorno relevantes:

| Variable | Descripcion |
|----------|-------------|
| `WHATSAPP_URL` | URL base de la WhatsApp Cloud API (fallback si no esta en platform_credentials) |
| `WHATSAPP_VERIFY_TOKEN` | Token para verificacion del webhook por Meta |
| `WHATSAPP_WEBHOOK_SECRET` | Secret para validacion HMAC-SHA256 de webhooks |
| `WHATSAPP_PHONE_NUMBER_ID` | ID del numero de telefono (legacy, solo para SendMessage) |
| `WHATSAPP_TOKEN` | Access token (legacy, deprecado; ahora se obtiene de la BD encriptado) |
| `ENCRYPTION_KEY` | Clave AES-256 para descifrar credenciales almacenadas en la BD |

Las credenciales principales (`access_token`, `phone_number_id`) se almacenan encriptadas (AES-256-GCM) en la tabla `integrations` y se descifran en tiempo de ejecucion.
