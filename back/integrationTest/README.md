# Integration Test - Simulador de Webhooks

Proyecto de pruebas para simular webhooks de diferentes integraciones (Shopify, WhatsApp, MercadoLibre, etc.) y enviarlos al servidor.

## Estructura

```
integrationTest/
├── cmd/
│   └── main.go                    # Punto de entrada
├── integrations/
│   └── shopify/
│       ├── bundle.go              # Bundle del módulo Shopify
│       └── internal/
│           ├── app/
│           │   └── usecases/     # Casos de uso para simular órdenes
│           └── infra/
│               └── primary/
│                   └── client/   # Cliente HTTP para enviar webhooks
└── shared/
    ├── env/                       # Configuración de variables de entorno
    └── log/                       # Logger
```

## Uso

1. Copiar `.env.example` a `.env` y configurar:
   - `WEBHOOK_BASE_URL`: URL base del servidor (ej: http://localhost:8080)
   - `SHOPIFY_SHOP_DOMAIN`: Dominio de la tienda Shopify (debe coincidir con StoreID de una integración existente)
   - `SHOPIFY_API_VERSION`: Versión de la API de Shopify (opcional, default: 2024-10)

2. Ejecutar:
```bash
go run cmd/main.go
```

3. Seleccionar el topic del webhook a simular desde el menú interactivo.

## Módulos

### Shopify
Simula órdenes de Shopify y las envía como webhooks a:
- `POST {WEBHOOK_BASE_URL}/api/v1/integrations/shopify/webhook`

#### Topics soportados:
- `orders/create` - Orden creada
- `orders/paid` - Orden pagada
- `orders/updated` - Orden actualizada
- `orders/cancelled` - Orden cancelada
- `orders/fulfilled` - Orden cumplida
- `orders/partially_fulfilled` - Orden parcialmente cumplida
