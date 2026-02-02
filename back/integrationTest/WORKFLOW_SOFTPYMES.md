# 🔄 Workflow de Testing de SoftPymes

## 📊 Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────┐
│            FLUJO COMPLETO DE TESTING - SOFTPYMES             │
└─────────────────────────────────────────────────────────────┘

   ┌──────────────────────┐                    ┌──────────────────────┐
   │   Central Backend    │                    │ SoftPymes Mock Server│
   │   (Sistema Real)     │  HTTP Requests     │  (Servidor HTTP)     │
   │   Puerto: 3050       │ ────────────────>  │  Puerto: 8082        │
   │                      │                    │                      │
   │  - Llama a SoftPymes │                    │  - Simula API        │
   │  - Crea facturas     │  <────────────────  │  - Retorna CUFEs     │
   │  - Guarda en BD      │    HTTP Responses   │  - Sin costos        │
   └──────────────────────┘                    └──────────────────────┘

   ┌──────────────────────┐
   │  Integration Test    │  (Opcional - solo para testing manual)
   │   (Simulador CLI)    │
   │                      │
   │  - Menú interactivo  │
   │  - Opciones 11-14    │
   │  - Testing directo   │
   └──────────────────────┘
```

## 🎯 Diferencia Clave con WhatsApp

| Aspecto | WhatsApp | SoftPymes |
|---------|----------|-----------|
| **Dirección** | Simulador → Central (webhooks) | Central → Mock Server (API calls) |
| **Puerto** | N/A (CLI) | 8082 (servidor HTTP) |
| **Uso** | Central recibe webhooks | Central llama al mock |

## 🔧 Variables de Entorno Configuradas

### `/back/central/.env` ✅

**Modo TESTING (Activo actualmente):**
```env
# 🧪 SoftPymes TESTING (Servidor mock local)
SOFTPYMES_API_URL=http://localhost:8082  # ⚠️ Mock server en puerto 8082
```

**Modo PRODUCCIÓN (Comentado):**
```env
# 📄 SoftPymes PRODUCCIÓN (API real)
# SOFTPYMES_API_URL=https://api-integracion.softpymes.com.co
```

## 🚀 Cómo Iniciar

### 1. Iniciar SoftPymes Mock Server

**Terminal 1 - Mock Server:**
```bash
cd /home/cam/Desktop/probability/back/integrationTest/integrations/softpymes
./start-server.sh
```

**Output esperado:**
```
🚀 SoftPymes Mock HTTP Server running on port 8082
📋 Endpoints available:
  POST /oauth/integration/login/
  POST /sales_invoice/
  POST /search/documents/notes/
  GET  /search/documents/
  GET  /health
```

**Verificar que está corriendo:**
```bash
curl http://localhost:8082/health
# Debe responder: {"status":"ok","service":"softpymes-mock","port":"8082"}
```

### 2. Iniciar Central Backend

**Terminal 2 - Central:**
```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go
# Debe decir: Server running on port 3050
```

### 3. (Opcional) Simulador CLI

**Terminal 3 - Simulador Interactivo:**
```bash
cd /home/cam/Desktop/probability/back/integrationTest
go run cmd/main.go
# Seleccionar opciones 11-14 para testing manual
```

## 🔄 Flujo de Autenticación y Facturación

### Desde el Central Backend

```
1. Central necesita crear factura
   ↓
2. Central llama: POST http://localhost:8082/oauth/integration/login/
   Body: {"apiKey": "...", "apiSecret": "..."}
   Header: Referer: https://tutienda.com
   ↓
3. Mock Server valida y retorna:
   {"accessToken": "spy_token_a1b2c3d4", "expiresInMin": 60}
   ↓
4. Central usa token para crear factura:
   POST http://localhost:8082/sales_invoice/
   Header: Authorization: Bearer spy_token_a1b2c3d4
   Body: {datos de factura}
   ↓
5. Mock Server crea factura simulada y retorna:
   {
     "success": true,
     "invoice_number": "SPY-1001",
     "cufe": "CUFE-...",
     "pdf_url": "https://softpymes-mock.local/..."
   }
   ↓
6. Central guarda factura en BD con CUFE y URLs
```

## 📋 Endpoints del Mock Server

### 1. POST `/oauth/integration/login/`

**Autenticación**

Request:
```json
{
  "apiKey": "test_key",
  "apiSecret": "test_secret"
}
```

Headers:
```
Referer: https://tutienda.com
Content-Type: application/json
```

Response (200 OK):
```json
{
  "success": true,
  "accessToken": "spy_token_a1b2c3d4",
  "expiresInMin": 60,
  "tokenType": "Bearer"
}
```

### 2. POST `/sales_invoice/`

**Crear Factura**

Headers:
```
Authorization: Bearer spy_token_a1b2c3d4
Content-Type: application/json
```

Request:
```json
{
  "order_id": "ORD-001",
  "customer": {
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "nit": "123456789"
  },
  "items": [
    {
      "description": "Producto Test",
      "quantity": 1,
      "unit_price": 100000,
      "tax": 19000,
      "total": 119000
    }
  ],
  "total": 100000
}
```

Response (200 OK):
```json
{
  "success": true,
  "message": "Invoice created successfully",
  "invoice_number": "SPY-1001",
  "external_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
  "invoice_url": "https://softpymes-mock.local/invoices/{uuid}",
  "pdf_url": "https://softpymes-mock.local/invoices/{uuid}.pdf",
  "xml_url": "https://softpymes-mock.local/invoices/{uuid}.xml",
  "cufe": "CUFE-a1b2c3d4e5f6g7h8",
  "issued_at": "2026-02-01T22:30:00Z"
}
```

### 3. POST `/search/documents/notes/`

**Crear Nota de Crédito**

Headers:
```
Authorization: Bearer spy_token_a1b2c3d4
Content-Type: application/json
```

Request:
```json
{
  "invoice_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
  "amount": 50000,
  "reason": "Devolución de producto",
  "note_type": "partial"
}
```

Response (200 OK):
```json
{
  "success": true,
  "message": "Credit note created successfully",
  "credit_note_number": "NC-2001",
  "external_id": "b2c3d4e5-f6g7-h8i9-j0k1-l2m3n4o5p6q7",
  "note_url": "https://softpymes-mock.local/credit-notes/{uuid}",
  "pdf_url": "https://softpymes-mock.local/credit-notes/{uuid}.pdf",
  "xml_url": "https://softpymes-mock.local/credit-notes/{uuid}.xml",
  "cufe": "CUFE-NC-i9j8k7l6m5n4",
  "issued_at": "2026-02-01T22:35:00Z"
}
```

### 4. GET `/health`

**Health Check**

Response (200 OK):
```json
{
  "status": "ok",
  "service": "softpymes-mock",
  "port": "8082"
}
```

## 🧪 Testing desde el Central

### Escenario 1: Facturar una Orden

1. **Central crea una orden** en el sistema
2. **Central llama a SoftPymes** (mock) para facturar
3. **Mock retorna CUFE y URLs**
4. **Central guarda** datos de factura en BD

**Verificar en BD:**
```sql
SELECT * FROM invoices WHERE order_id = 'ORD-001';
-- Debe tener cufe, invoice_number, pdf_url, xml_url
```

### Escenario 2: Crear Nota de Crédito

1. **Buscar factura** en BD que ya existe
2. **Central llama a SoftPymes** para crear nota de crédito
3. **Mock retorna** número de nota y CUFE
4. **Central vincula** nota con factura en BD

## 🐛 Debugging

### Ver logs del Mock Server

Los logs aparecen en la terminal donde corre el servidor:

```
[22:30:45] POST /oauth/integration/login/ - Status: 200 - Duration: 2ms
[22:30:46] POST /sales_invoice/ - Status: 200 - Duration: 5ms
```

### Ver logs del Central

```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go 2>&1 | grep softpymes
```

### Problemas comunes

| Problema | Causa | Solución |
|----------|-------|----------|
| "connection refused" | Mock server no está corriendo | Iniciar `./start-server.sh` |
| "port 8082 in use" | Puerto ocupado | Script mata proceso automáticamente |
| "invalid token" | Token no existe o expiró | Re-autenticar |
| "401 Unauthorized" | Falta header Authorization | Verificar que central envía Bearer token |

### Verificar flujo completo

```bash
# 1. Health check del mock
curl http://localhost:8082/health

# 2. Autenticar
curl -X POST http://localhost:8082/oauth/integration/login/ \
  -H "Content-Type: application/json" \
  -H "Referer: https://test.com" \
  -d '{"apiKey":"test","apiSecret":"secret"}'

# Guarda el token recibido
TOKEN="spy_token_xxxx"

# 3. Crear factura
curl -X POST http://localhost:8082/sales_invoice/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "order_id": "ORD-001",
    "customer": {"name": "Test", "email": "test@example.com", "nit": "123"},
    "items": [{"description": "Producto", "quantity": 1, "unit_price": 100000, "tax": 19000, "total": 119000}],
    "total": 100000
  }'
```

## 🔄 Cambiar a Producción

Cuando quieras usar **SoftPymes API real**:

1. **En `/back/central/.env`:**
   - Comentar: `SOFTPYMES_API_URL=http://localhost:8082`
   - Descomentar: `SOFTPYMES_API_URL=https://api-integracion.softpymes.com.co`

2. **Reiniciar Central Backend**

3. **Detener Mock Server** (ya no se necesita)

4. **El sistema ahora usa SoftPymes API real**

## ✅ Checklist de Verificación

Antes de empezar testing:

- [ ] Mock server corriendo en puerto 8082
- [ ] Health check responde: `curl http://localhost:8082/health`
- [ ] Central corriendo en puerto 3050
- [ ] Variable de entorno `SOFTPYMES_API_URL=http://localhost:8082` en central
- [ ] Base de datos corriendo
- [ ] Tablas de facturas creadas (migrations)

## 📊 Ventajas del Mock

1. **Sin Costos**: No consume créditos de SoftPymes
2. **Sin Límites**: No hay rate limiting
3. **Rápido**: Respuestas instantáneas
4. **Reproducible**: Mismos resultados siempre
5. **Offline**: No requiere internet
6. **CUFEs Válidos**: Genera CUFEs simulados sin DIAN
7. **Testing Seguro**: No afecta facturas reales

---

**Última actualización:** 2026-02-01
**Modo actual:** TESTING (Mock Server Local en puerto 8082)
