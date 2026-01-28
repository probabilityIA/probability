# Configuración del Webhook de WhatsApp Business

## ✅ Implementación Completada

El webhook de WhatsApp ya está completamente implementado en el backend y frontend. Ahora puedes ver la URL del webhook directamente desde la interfaz de administración.

---

## 📋 Cómo Acceder al Webhook

### 1. **Desde el Dashboard (Recomendado)**

1. Inicia sesión en el dashboard admin: `http://localhost:3000`
2. Ve a **Integraciones** → **Ver/Editar WhatsApp**
3. En la sección **"🔗 Configuración del Webhook"** verás:
   - ✅ URL del webhook (con botón de copiar)
   - ✅ Eventos a suscribir
   - ✅ Verify Token para Meta
   - ✅ Instrucciones completas

### 2. **Variables de Entorno**

El webhook se construye a partir de la variable:

```env
# En .env
WEBHOOK_BASE_URL=http://localhost:3050
```

**En producción**, cambia esto por tu dominio público:
```env
WEBHOOK_BASE_URL=https://api.probability.com
```

---

## 🔧 Configurar el Webhook en Meta Business Manager

### Paso 1: Exponer tu Backend (Solo para Desarrollo Local)

Si estás en desarrollo local, necesitas exponer tu backend con **ngrok**:

```bash
# Instalar ngrok (si no lo tienes)
# https://ngrok.com/download

# Exponer puerto 3050
ngrok http 3050
```

**Ngrok te dará una URL** tipo: `https://abc123.ngrok.io`

**Actualiza tu .env:**
```env
WEBHOOK_BASE_URL=https://abc123.ngrok.io
```

**Reinicia el backend:**
```bash
cd /home/cam/Desktop/probability/back/central
go run cmd/main.go
```

### Paso 2: Configurar en Meta Business Manager

1. Ve a: https://business.facebook.com/

2. Selecciona tu **WhatsApp Business Account**

3. Ve a **API Setup** → **Configuration**

4. En la sección **Webhook**, haz clic en **Edit**

5. **Configure webhooks**:

   **Callback URL:**
   ```
   https://abc123.ngrok.io/integrations/whatsapp/webhook
   ```
   (O la URL que veas en el dashboard)

   **Verify Token:**
   ```
   probability_whatsapp_verify_token_2026_secure
   ```

6. **Subscribe to fields** (selecciona estos campos):
   - ✅ `messages` (mensajes entrantes y respuestas de botones)
   - ✅ `message_template_status_update` (estado de plantillas)

7. Haz clic en **Verify and Save**

   Meta enviará un request GET al webhook para verificarlo.

   Deberías ver un mensaje de **"Success"** ✅

---

## 🔍 Verificar que el Webhook Funciona

### 1. Ver Logs en Tiempo Real

```bash
tail -f /tmp/backend_new.log | grep -i "webhook"
```

**Logs esperados cuando Meta verifica:**
```
[Webhook Handler] - solicitud de verificación de webhook
mode: subscribe
token: probability_whatsapp_verify_token_2026_secure
[Webhook Handler] - webhook verificado exitosamente
```

### 2. Enviar un Mensaje de Prueba

Publica un evento de orden:

```bash
curl -u admin:admin -X POST http://localhost:15672/api/exchanges/%2F/amq.default/publish \
  -H "Content-Type: application/json" \
  -d '{
    "properties": {},
    "routing_key": "orders.confirmation.requested",
    "payload": "{\"event_type\":\"order.confirmation_requested\",\"order_id\":\"test-webhook-001\",\"order_number\":\"TEST-WEBHOOK-001\",\"business_id\":1,\"customer_name\":\"Tu Nombre\",\"customer_phone\":\"+TU_NUMERO_AQUI\",\"customer_email\":\"test@example.com\",\"total_amount\":50000,\"currency\":\"COP\",\"items_summary\":\"1x Producto Test\",\"shipping_address\":\"Dirección de Prueba\",\"payment_method\":\"Contraentrega\",\"integration_id\":2,\"platform\":\"test\",\"timestamp\":1738033000}",
    "payload_encoding": "string"
  }'
```

**Reemplaza `+TU_NUMERO_AQUI`** con tu número de WhatsApp real.

### 3. Verificar Estados de Mensajes

Después de enviar el mensaje, el webhook recibirá automáticamente eventos de estado:

```bash
tail -f /tmp/backend_new.log | grep -E "status|delivered|read"
```

**Logs esperados:**
```
[WhatsApp Webhook] - procesando cambios de estado de mensajes
message_id: wamid.xxx
status: sent
[WhatsApp Webhook] - estado de mensaje actualizado exitosamente

message_id: wamid.xxx
status: delivered
[WhatsApp Webhook] - estado de mensaje actualizado exitosamente

message_id: wamid.xxx
status: read
[WhatsApp Webhook] - estado de mensaje actualizado exitosamente
```

### 4. Verificar en Base de Datos

```sql
-- Ver últimos mensajes con sus estados
SELECT
  message_id,
  phone_number,
  template_name,
  status,
  delivered_at,
  read_at,
  created_at
FROM whatsapp_message_logs
ORDER BY created_at DESC
LIMIT 10;
```

**Estados posibles:**
- `sent` → Enviado a WhatsApp
- `delivered` → Entregado al teléfono del usuario
- `read` → Usuario leyó el mensaje
- `failed` → Mensaje falló

---

## 🎯 Eventos que Maneja el Webhook

### 1. **Cambios de Estado de Mensajes** (`statuses`)

Notifica cuando un mensaje cambia de estado:
- ✅ **sent** → WhatsApp aceptó el mensaje
- ✅ **delivered** → Mensaje entregado al usuario
- ✅ **read** → Usuario leyó el mensaje
- ❌ **failed** → Mensaje falló (con detalles del error)

**Actualización automática en BD:**
```go
// El webhook actualiza automáticamente whatsapp_message_logs:
- status
- delivered_at
- read_at
```

### 2. **Mensajes Entrantes** (`messages`)

Recibe respuestas del usuario:
- ✅ **Botones** (quick_reply) → "Confirmar pedido", "No confirmar"
- ✅ **Texto** → Mensajes de texto del usuario
- ✅ **Interactivos** → Listas, botones interactivos

**Flujo de Conversación:**
```
START → Usuario presiona "Confirmar pedido"
     → Publica evento a RabbitMQ: orders.whatsapp.confirmed
     → Orders module actualiza orden: is_confirmed = true

START → Usuario presiona "No confirmar"
     → Muestra menú de opciones (novedad, cancelar, asesor)
```

### 3. **Estado de Plantillas** (`message_template_status_update`)

Notifica cambios en plantillas:
- Aprobada
- Rechazada
- En revisión

---

## 🐛 Troubleshooting

### Problema 1: "Webhook verification failed"

**Causa:** El verify token no coincide.

**Solución:**
1. Verifica que en Meta uses: `probability_whatsapp_verify_token_2026_secure`
2. Verifica que en `.env` tengas:
   ```env
   WHATSAPP_VERIFY_TOKEN="probability_whatsapp_verify_token_2026_secure"
   ```
3. Reinicia el backend

### Problema 2: "Connection refused" al verificar

**Causa:** La URL del webhook no es accesible desde internet.

**Solución:**
- **Local**: Usa ngrok para exponer tu backend
- **Producción**: Asegúrate que el dominio sea público y tenga SSL (HTTPS)

### Problema 3: "Signature validation failed"

**Causa:** El HMAC secret no coincide.

**Solución:**
1. Verifica que en `.env` tengas:
   ```env
   WHATSAPP_WEBHOOK_SECRET="probability_webhook_secret_hmac_sha256_2026"
   ```
2. El secret debe coincidir con el que configuraste en Meta (App Secret)

### Problema 4: No llegan estados de mensajes

**Causa:** No estás suscrito al campo `messages` en Meta.

**Solución:**
1. Ve a Meta Business Manager → WhatsApp → Configuration → Webhooks
2. Asegúrate que esté seleccionado:
   - ✅ `messages`
   - ✅ `message_template_status_update`

### Problema 5: Logs muestran errores de firma

**Logs:**
```
[Webhook Handler] - firma inválida
expected: abc123
calculated: xyz789
```

**Solución:**
El App Secret de Meta no coincide con `WHATSAPP_WEBHOOK_SECRET`.

1. Ve a Meta Business Manager → Settings → App Settings
2. Copia el **App Secret**
3. Actualiza `.env`:
   ```env
   WHATSAPP_WEBHOOK_SECRET="<TU_APP_SECRET_DE_META>"
   ```

---

## 📊 Monitoreo en Producción

### 1. Logs Estructurados

```bash
# Ver todos los eventos de webhook
tail -f /var/log/probability/backend.log | jq 'select(.module == "whatsapp-webhook")'

# Ver solo errores
tail -f /var/log/probability/backend.log | jq 'select(.level == "error" and .module == "whatsapp-webhook")'
```

### 2. Métricas Recomendadas (Prometheus)

```
# Contador de webhooks recibidos
whatsapp_webhook_received_total{type="messages"}
whatsapp_webhook_received_total{type="statuses"}

# Contador de errores
whatsapp_webhook_errors_total{error_type="invalid_signature"}
whatsapp_webhook_errors_total{error_type="processing_failed"}

# Latencia de procesamiento
whatsapp_webhook_processing_duration_seconds
```

### 3. Alertas Recomendadas

- ⚠️ Tasa de errores > 5% en últimos 5 minutos
- ⚠️ Webhooks con firma inválida > 10 en última hora
- ⚠️ Latencia de procesamiento > 3 segundos

---

## 🔐 Seguridad

### 1. Validación de Firma HMAC-SHA256

El webhook **valida automáticamente** que los requests vengan de Meta:

```go
// El header X-Hub-Signature-256 debe coincidir con:
HMAC-SHA256(payload, WHATSAPP_WEBHOOK_SECRET)
```

### 2. SSL/TLS Requerido

Meta **solo envía webhooks a URLs HTTPS** (excepto localhost en desarrollo).

### 3. Tokens Seguros

**En producción**, usa tokens fuertes:

```env
# Generar tokens seguros:
WHATSAPP_VERIFY_TOKEN=$(openssl rand -base64 32)
WHATSAPP_WEBHOOK_SECRET=$(openssl rand -base64 32)
```

---

## ✅ Checklist de Configuración

- [ ] Variable `WEBHOOK_BASE_URL` configurada en `.env`
- [ ] Backend corriendo y accesible desde internet
- [ ] Webhook configurado en Meta Business Manager
- [ ] Verify Token: `probability_whatsapp_verify_token_2026_secure`
- [ ] Campos suscritos: `messages` y `message_template_status_update`
- [ ] Webhook verificado exitosamente (✓ Success en Meta)
- [ ] Mensaje de prueba enviado
- [ ] Estados recibidos en logs (`sent`, `delivered`, `read`)
- [ ] Base de datos actualizada correctamente

---

## 📚 Referencias

- [WhatsApp Business API Webhooks](https://business.whatsapp.com/blog/how-to-use-webhooks-from-whatsapp-business-api/)
- [Receiving Messages - WhatsApp SDK](https://whatsapp.github.io/WhatsApp-Nodejs-SDK/receivingMessages/)
- [Webhook Configuration - 360Dialog](https://docs.360dialog.com/docs/waba-messaging/webhook)

---

**Implementado por**: Claude Code
**Fecha**: 2026-01-27
**Branch**: feature/whatsapp-bidirectional-integration
