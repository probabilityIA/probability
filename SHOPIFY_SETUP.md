# Configuración de Shopify para Producción

## 📋 Variables de Entorno Requeridas

Agregar al archivo `/home/ubuntu/probability/.env` en el servidor:

```bash
# 🛒 Shopify Integration Configuration
SHOPIFY_CLIENT_ID=tu_shopify_client_id_aqui
SHOPIFY_CLIENT_SECRET=tu_shopify_client_secret_aqui
SHOPIFY_REDIRECT_URI=https://app.probabilityia.com.co/api/v1/shopify/callback
SHOPIFY_SCOPES=read_products,write_products,read_orders,write_orders,read_customers,write_customers
SHOPIFY_API_VERSION=2024-01
```

## 🔗 URLs de Webhooks de Compliance (OBLIGATORIOS)

Estas URLs deben configurarse en tu Shopify App Dashboard:

### 1. Customer Data Request (GDPR - Solicitud de Datos)
```
https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/customers/data_request
```

### 2. Customer Redact (GDPR - Eliminación de Datos)
```
https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/customers/redact
```

### 3. Shop Redact (Desinstalación de App)
```
https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/shop/redact
```

## 📝 Pasos para Configuración

### 1. Obtener Credenciales de Shopify

1. Ve a tu [Shopify Partners Dashboard](https://partners.shopify.com/)
2. Selecciona tu app o crea una nueva
3. En "App setup", encontrarás:
   - **Client ID** (API key)
   - **Client Secret** (API secret key)

### 2. Configurar OAuth en Shopify

En la configuración de tu app de Shopify:

**Allowed redirection URL(s):**
```
https://app.probabilityia.com.co/api/v1/shopify/callback
```

### 3. Configurar Webhooks de Compliance

En "App setup" → "GDPR webhooks":

| Webhook Topic | Endpoint URL |
|--------------|-------------|
| Customer data request | `https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/customers/data_request` |
| Customer data erasure | `https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/customers/redact` |
| Shop data erasure | `https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/shop/redact` |

### 4. Actualizar .env en Producción

```bash
# SSH al servidor
ssh -i "probability.pem" ubuntu@ec2-3-224-189-33.compute-1.amazonaws.com

# Editar el archivo .env
cd /home/ubuntu/probability
nano .env

# Agregar las variables de Shopify (ver arriba)

# Reiniciar servicios
sudo docker compose down
sudo docker compose up -d
```

### 5. Verificar Configuración

Prueba los endpoints:

```bash
# Test OAuth callback (debe devolver error de parámetros faltantes)
curl https://app.probabilityia.com.co/api/v1/shopify/callback

# Test compliance webhook (debe devolver 401 sin HMAC válido)
curl -X POST https://app.probabilityia.com.co/api/v1/integrations/shopify/webhooks/customers/data_request
```

## ⚠️ Requisitos de Compliance

Según [Shopify Privacy Compliance](https://shopify.dev/docs/apps/build/compliance/privacy-law-compliance):

1. **Respuesta inmediata**: Los webhooks deben responder con `200 OK` inmediatamente
2. **Plazo de procesamiento**: Completar las acciones dentro de **30 días**
3. **Validación HMAC**: Todos los webhooks deben validar la firma HMAC de Shopify
4. **Obligatorio para App Store**: Estos webhooks son REQUERIDOS para publicar apps en Shopify App Store

## 🔒 Scopes Requeridos

Los scopes configurados permiten:

- `read_products, write_products` - Gestión de productos
- `read_orders, write_orders` - Gestión de pedidos
- `read_customers, write_customers` - Gestión de clientes (necesario para compliance)

## 📊 Endpoints de la Integración

| Endpoint | Método | Descripción | Auth |
|----------|--------|-------------|------|
| `/api/v1/integrations/shopify/connect` | POST | Iniciar OAuth | JWT |
| `/api/v1/shopify/callback` | GET | Callback OAuth | State + HMAC |
| `/api/v1/integrations/shopify/webhook` | POST | Webhook general | HMAC |
| `/api/v1/integrations/shopify/webhooks/customers/data_request` | POST | GDPR - Solicitud de datos | HMAC |
| `/api/v1/integrations/shopify/webhooks/customers/redact` | POST | GDPR - Eliminación de datos | HMAC |
| `/api/v1/integrations/shopify/webhooks/shop/redact` | POST | Desinstalación de app | HMAC |

## 🚀 Próximos Pasos

Una vez configurado:

1. [ ] Implementar la lógica de procesamiento asíncrono en los handlers de compliance
2. [ ] Crear jobs/workers para procesar solicitudes GDPR dentro del plazo de 30 días
3. [ ] Implementar sistema de auditoría para compliance
4. [ ] Configurar notificaciones para solicitudes de datos/eliminación
5. [ ] Documentar procesos de retención de datos según requisitos legales

## 📚 Referencias

- [Shopify Privacy Compliance](https://shopify.dev/docs/apps/build/compliance/privacy-law-compliance)
- [Shopify Webhooks](https://shopify.dev/docs/apps/build/webhooks)
- [Shopify OAuth](https://shopify.dev/docs/apps/build/authentication-authorization/access-tokens/authorization-code-grant)
