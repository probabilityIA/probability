# Testing - Simuladores de APIs Externas

**Monolito Go** para simular APIs externas (Softpymes, Shopify, etc.) en entorno de testing.

## 🏗️ Arquitectura

```
testing/
+-- main.go                # ✅ Entry point monolito
+-- go.mod                 # ✅ Módulo único
+-- integrations/
|   +-- softpymes/         # ✅ Simulador Softpymes
|   |   +-- bundle.go      # Inicialización y servidor HTTP
|   |   +-- routes.go      # Handlers HTTP
|   |   +-- internal/      # Lógica interna
|   +-- shopify/           # Simulador Shopify (webhooks)
|   +-- whatsapp/          # (Futuro)
+-- shared/
    +-- log/               # Logger compartido
```

## 🚀 Uso

### Iniciar Testing Server

Al iniciar, se ejecutan **simultáneamente**:
1. **Servidor HTTP de Softpymes** (puerto 9090) - Para que el backend pueda consumir la API
2. **CLI Interactivo** - Menú para simular webhooks de Shopify/WhatsApp manualmente

```bash
cd /home/cam/Desktop/probability/back/testing

# Iniciar todo
go run cmd/main.go

# O usar Makefile
make run-testing
```

**Funcionalidades:**
- ✅ El backend puede llamar a `http://localhost:9090` para facturación (Softpymes)
- ✅ Desde la terminal puedes simular webhooks de Shopify (crear órdenes, marcar como pagado, etc.)
- ✅ Desde la terminal puedes simular respuestas de WhatsApp
- ✅ Desde la terminal puedes simular autenticación y facturas de Softpymes

### Configurar Puertos

```bash
# Variables de entorno
export SOFTPYMES_MOCK_PORT=9090

go run main.go
```

## 📋 Simuladores Disponibles

### 1. Softpymes ✅ (Puerto 9090)

**Endpoints:**
- `POST /oauth/integration/login/` - Autenticación
- `POST /app/integration/sales_invoice/` - Crear factura
- `POST /app/integration/search/documents/` - Buscar documentos
- `GET /health` - Health check

**Uso:**
```bash
export SOFTPYMES_API_URL=http://localhost:9090
cd ../central && go run cmd/main.go
```

### 2. Shopify (Webhooks) ✅

**Configuración del Business de Prueba:**

El simulador está configurado para usar datos reales del business de pruebas en la base de datos:

| Campo | Valor |
|-------|-------|
| **Business ID** | `7` |
| **Business Name** | `probability-dev` |
| **Integration ID** | `1` |
| **Integration Name** | `Shopify - pruebas` |
| **Shop Domain** | `tienda-test.myshopify.com` |
| **API Version** | `2024-01` |

**Cómo funciona:**
- Las órdenes simuladas incluyen `note_attributes` con metadatos (`_business_id`, `_integration_id`, etc.)
- El header `X-Shopify-Shop-Domain` se envía con `tienda-test.myshopify.com`
- El backend identifica la integración por el shop domain y procesa la orden correctamente

**Configuración:** Ver archivo `config.go` en `/integrations/shopify/internal/domain/`

## 🔧 Agregar Nuevo Simulador

1. Crear directorio: `integrations/nuevo-servicio/`
2. Crear `bundle.go` con función `New(logger, port)` y método `Start()`
3. Crear `routes.go` con `RegisterRoutes(router)`
4. Registrar en `main.go`:

```go
nuevoServicio := nuevoservicio.New(logger, port)
go nuevoServicio.Start()
```

---

**Módulo:** `github.com/secamc93/probability/back/testing`
**Última actualización:** 2026-02-09













