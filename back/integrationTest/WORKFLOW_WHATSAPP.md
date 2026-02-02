# 🔄 Workflow de Testing de WhatsApp

## 📊 Arquitectura del Sistema

```
┌─────────────────────────────────────────────────────────────┐
│                    FLUJO COMPLETO DE TESTING                 │
└─────────────────────────────────────────────────────────────┘

   ┌──────────────────────┐
   │  Integration Test    │  Puerto: N/A (no es servidor)
   │   (Simulador CLI)    │  Ubicación: /back/integrationTest
   │                      │
   │  - Menú interactivo  │
   │  - Simula webhooks   │
   │  - Envía a Central   │
   └──────┬───────────────┘
          │
          │ HTTP POST (webhooks)
          │ http://localhost:3050/api/integrations/whatsapp/webhook
          │
          ↓
   ┌──────────────────────┐
   │   Central Backend    │  Puerto: 3050
   │   (Sistema Real)     │  Ubicación: /back/central
   │                      │
   │  - Procesa webhooks  │
   │  - Actualiza BD      │
   │  - Lógica de negocio │
   └──────────────────────┘
```

## 🎯 Variables de Entorno Configuradas

### Integration Test (`.env`)

```env
# El simulador envía webhooks a esta URL
WEBHOOK_BASE_URL=http://localhost:3050

# Secret para firmar webhooks (debe coincidir con central)
WHATSAPP_WEBHOOK_SECRET=test_webhook_secret

# Delay antes de enviar respuestas automáticas
WHATSAPP_AUTO_REPLY_DELAY=2
```

### Central Backend (`.env`)

**Modo TESTING (Activo actualmente):**
```env
# Simulador local - NO usa Meta API real
WHATSAPP_URL="http://localhost:8081/v18.0/"
WHATSAPP_TOKEN="mock_token_for_testing"
WHATSAPP_PHONE_NUMBER_ID="123456789012345"
WHATSAPP_VERIFY_TOKEN="test_verify_token"
WHATSAPP_WEBHOOK_SECRET="test_webhook_secret"  # ⚠️ DEBE COINCIDIR con integrationTest
```

**Modo PRODUCCIÓN (Comentado):**
```env
# Meta API real - Descomentar para producción
# WHATSAPP_URL="https://graph.facebook.com/v22.0/"
# WHATSAPP_TOKEN="<token_real_de_meta>"
# WHATSAPP_PHONE_NUMBER_ID="479455621898882"
# WHATSAPP_VERIFY_TOKEN="<token_real>"
# WHATSAPP_WEBHOOK_SECRET="<secret_real>"
```

## 🚀 Cómo Usar

### 1. Iniciar Central Backend

```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go
```

**Verificar que esté corriendo:**
```bash
curl http://localhost:3050/health
# Debe responder: 200 OK
```

### 2. Iniciar Integration Test (Simulador)

```bash
cd /home/cam/Desktop/probability/back/integrationTest
go run cmd/main.go
```

**Menú que aparece:**
```
=== Simulador de Webhooks ===

📦 SHOPIFY:
1. orders/create (crear nueva orden aleatoria)
2. orders/paid (marcar orden como pagada)
...

💬 WHATSAPP:
8. Simular respuesta de usuario (manual)
9. Simular respuesta automática (por template)
10. Listar conversaciones almacenadas

0. Salir
```

### 3. Simular Respuesta Manual (Opción 8)

**Ejemplo de uso:**

```
Opción: 8
Número de teléfono (ej: +573001234567): +573001234567
Respuesta del usuario (ej: Confirmar pedido): Confirmar pedido
```

**Lo que sucede internamente:**

1. **Espera 2 segundos** (WHATSAPP_AUTO_REPLY_DELAY)
2. **Envía webhook de estado "delivered"** al central
   ```
   POST http://localhost:3050/api/integrations/whatsapp/webhook
   Header: X-Hub-Signature-256: sha256=<hmac>
   Body: { status: "delivered", ... }
   ```
3. **Espera 500ms**
4. **Envía webhook de estado "read"**
5. **Espera 500ms**
6. **Envía webhook con respuesta del usuario**
   ```json
   {
     "object": "whatsapp_business_account",
     "entry": [{
       "changes": [{
         "value": {
           "messages": [{
             "from": "+573001234567",
             "type": "button",
             "button": {
               "payload": "Confirmar pedido",
               "text": "Confirmar pedido"
             }
           }]
         }
       }]
     }]
   }
   ```

### 4. Simular Respuesta Automática (Opción 9)

**Ejemplo de uso:**

```
Opción: 9
Número de teléfono (ej: +573001234567): +573001234567
Nombre del template (ej: confirmacion_pedido_contraentrega): confirmacion_pedido_contraentrega
```

**Lo que sucede:**

El simulador automáticamente responde según el template:
- `confirmacion_pedido_contraentrega` → "Confirmar pedido"
- `menu_no_confirmacion` → "Presentar novedad"
- `solicitud_cancelacion` → "Sí, cancelar"
- etc.

### 5. Listar Conversaciones (Opción 10)

Muestra todas las conversaciones simuladas en memoria:

```
💬 Conversaciones almacenadas (3):
  1. +573001234567 - Estado: COMPLETED - Orden: ORD-001 - Mensajes: 2
  2. +573002345678 - Estado: WAITING_CONFIRMATION - Orden: ORD-002 - Mensajes: 1
  3. +573003456789 - Estado: HANDOFF_TO_HUMAN - Orden: ORD-003 - Mensajes: 4
```

## 🔐 Validación de Firma HMAC-SHA256

El simulador firma todos los webhooks con HMAC-SHA256:

```go
// Cálculo de firma
secret := "test_webhook_secret"
h := hmac.New(sha256.New, []byte(secret))
h.Write(payloadBytes)
signature := hex.EncodeToString(h.Sum(nil))

// Header enviado
X-Hub-Signature-256: sha256=<signature>
```

El **Central Backend** DEBE validar esta firma:

```go
// En el webhook handler del central
expectedSignature := calculateHMAC(webhookSecret, requestBody)
if receivedSignature != expectedSignature {
    return errors.New("invalid signature")
}
```

## 📋 Mapeo Template → Respuesta

| Template | Respuesta Automática | Estado Final Esperado |
|----------|---------------------|----------------------|
| `confirmacion_pedido_contraentrega` | "Confirmar pedido" | COMPLETED |
| `menu_no_confirmacion` | "Presentar novedad" | WAITING_INPUT |
| `solicitud_novedad` | "Otro" | WAITING_INPUT |
| `novedad_otro` | "El producto llegó dañado" | HANDOFF_TO_HUMAN |
| `solicitud_cancelacion` | "Sí, cancelar" | WAITING_INPUT |
| `motivo_cancelacion` | "Ya no lo necesito" | COMPLETED |
| `pedido_cancelado` | (sin respuesta) | - |
| `pedido_confirmado` | (sin respuesta) | - |

## 🧪 Escenarios de Testing

### Escenario 1: Happy Path (Confirmación)

1. **Central envía template** `confirmacion_pedido_contraentrega` al usuario real
2. **En el simulador:** Opción 9 → template `confirmacion_pedido_contraentrega`
3. **Simulador responde:** "Confirmar pedido"
4. **Central procesa:** Orden pasa a CONFIRMED
5. **Verificar en BD:**
   ```sql
   SELECT * FROM whatsapp_conversations WHERE phone_number = '+573001234567';
   -- current_state debe ser 'COMPLETED'

   SELECT * FROM orders WHERE id = <order_id>;
   -- order_status debe actualizarse
   ```

### Escenario 2: Cancelación Completa

1. **Central envía:** `confirmacion_pedido_contraentrega`
2. **Simulador opción 8:** "No confirmar"
3. **Central envía:** `menu_no_confirmacion`
4. **Simulador opción 8:** "Cancelar pedido"
5. **Central envía:** `solicitud_cancelacion`
6. **Simulador opción 8:** "Sí, cancelar"
7. **Central envía:** `motivo_cancelacion`
8. **Simulador opción 8:** "Ya no lo necesito"
9. **Central procesa:** Orden cancelada

### Escenario 3: Reportar Novedad

1. **Central envía:** `confirmacion_pedido_contraentrega`
2. **Simulador opción 8:** "No confirmar"
3. **Central envía:** `menu_no_confirmacion`
4. **Simulador opción 8:** "Presentar novedad"
5. **Central envía:** `solicitud_novedad`
6. **Simulador opción 8:** "Otro"
7. **Central envía:** `novedad_otro`
8. **Simulador opción 8:** "El paquete llegó dañado"
9. **Central procesa:** Handoff to human

## 🐛 Debugging

### Ver logs del Central

```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go 2>&1 | grep whatsapp
```

### Ver logs del Simulador

Los logs aparecen directamente en la consola donde corre el simulador:

```
22:41:52 INF Simulando respuesta de usuario phone_number=+573001234567 response="Confirmar pedido"
22:41:54 INF Webhook de estado enviado message_id=wamid.HBg... status=delivered
22:41:55 INF Webhook de estado enviado message_id=wamid.HBg... status=read
22:41:55 INF Webhook de mensaje enviado button_text="Confirmar pedido" phone_number=+573001234567
```

### Verificar que el webhook llega al Central

```bash
# En el central, agregar logs en el handler:
logger.Info().Msg("Webhook recibido de WhatsApp")
```

### Problemas comunes

| Problema | Causa | Solución |
|----------|-------|----------|
| "invalid signature" | Secrets no coinciden | Verificar que `WHATSAPP_WEBHOOK_SECRET` sea igual en ambos .env |
| "connection refused" | Central no está corriendo | Iniciar el central en puerto 3050 |
| No se crean conversaciones | Endpoint del central no existe | Verificar que `/api/integrations/whatsapp/webhook` esté implementado |
| Delay muy largo | WHATSAPP_AUTO_REPLY_DELAY alto | Reducir a 1 o 2 segundos |

## 🔄 Cambiar a Producción

Cuando quieras usar **Meta API real** en vez del simulador:

1. **En `/back/central/.env`:**
   - Comentar las variables de TESTING
   - Descomentar las variables de PRODUCCIÓN

2. **Reiniciar Central Backend**

3. **El sistema ahora usará Meta API real**
   - WHATSAPP_URL → `https://graph.facebook.com/v22.0/`
   - WHATSAPP_TOKEN → Token real de Meta
   - WHATSAPP_WEBHOOK_SECRET → Secret real de Meta

## ✅ Checklist de Verificación

Antes de empezar testing:

- [ ] Central corriendo en puerto 3050
- [ ] Variables de entorno de testing activas en central
- [ ] `WHATSAPP_WEBHOOK_SECRET` igual en ambos .env
- [ ] Simulador compilado sin errores
- [ ] Base de datos corriendo
- [ ] Tablas de WhatsApp creadas (migrations)
- [ ] Endpoint `/api/integrations/whatsapp/webhook` implementado

---

**Última actualización:** 2026-02-01
**Modo actual:** TESTING (Simulador Local)
