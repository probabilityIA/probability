# Ejemplos de API - Integrations Core Module

Base URL: `/api/v1`

## 1. GET /integrations - Obtener todas las integraciones

**URL**: `GET /api/v1/integrations`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (por defecto: 1)
- `page_size` (int, opcional): Tamaño de página (por defecto: 10, máximo: 100)
- `type` (string, opcional): Filtrar por tipo de integración (whatsapp, shopify, mercado_libre)
- `category` (string, opcional): Filtrar por categoría (internal, external)
- `business_id` (int, opcional): Filtrar por business ID (NULL para integraciones globales)
- `is_active` (bool, opcional): Filtrar por estado activo
- `search` (string, opcional): Buscar por nombre o código (búsqueda parcial)

**Ejemplo Request**:
```
GET /api/v1/integrations?page=1&page_size=10&type=whatsapp&is_active=true
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integraciones obtenidas exitosamente",
  "data": [
    {
      "id": 1,
      "name": "WhatsApp Principal",
      "code": "whatsapp_platform",
      "type": "whatsapp",
      "category": "internal",
      "business_id": null,
      "is_active": true,
      "is_default": true,
      "config": {
        "phone_number_id": "123456789",
        "webhook_url": "https://api.example.com/webhooks/whatsapp",
        "template_language": "es",
        "default_country_code": "+57"
      },
      "description": "Integración principal de WhatsApp",
      "created_by_id": 1,
      "updated_by_id": null,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "name": "Shopify Store 1",
      "code": "shopify_store_1",
      "type": "shopify",
      "category": "external",
      "business_id": 16,
      "is_active": true,
      "is_default": true,
      "config": {
        "store_name": "mi-tienda",
        "api_version": "2024-01",
        "timezone": "America/Bogota"
      },
      "description": "Integración de Shopify para Business 16",
      "created_by_id": 1,
      "updated_by_id": 1,
      "created_at": "2024-01-16T14:20:00Z",
      "updated_at": "2024-01-16T15:30:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 10,
  "total_pages": 1
}
```

**Response 401 Unauthorized**:
```json
{
  "success": false,
  "message": "Token de autorización requerido",
  "error": "Token inválido"
}
```

**Permisos**: Requiere autenticación JWT válida

---

## 2. GET /integrations/:id - Obtener integración por ID

**URL**: `GET /api/v1/integrations/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
GET /api/v1/integrations/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración obtenida exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "type": "whatsapp",
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "123456789",
      "webhook_url": "https://api.example.com/webhooks/whatsapp",
      "template_language": "es",
      "default_country_code": "+57",
      "api_version": "v18.0"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Integración no encontrada",
  "error": "integración con ID 999 no encontrada"
}
```

**Permisos**: Requiere autenticación JWT válida

---

## 3. GET /integrations/type/:type - Obtener integración por tipo

**URL**: `GET /api/v1/integrations/type/:type`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `type` (string, requerido): Tipo de integración (whatsapp, shopify, mercado_libre)

**Query Parameters**:
- `business_id` (int, opcional): ID del business (para integraciones por business). Si no se especifica, busca integraciones globales.

**Ejemplo Request**:
```
GET /api/v1/integrations/type/whatsapp
GET /api/v1/integrations/type/shopify?business_id=16
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración obtenida exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "type": "whatsapp",
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "123456789",
      "webhook_url": "https://api.example.com/webhooks/whatsapp",
      "template_language": "es"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Integración no encontrada",
  "error": "integración activa de tipo 'whatsapp' no encontrada"
}
```

**Permisos**: Requiere autenticación JWT válida

**Nota**: Este endpoint es útil para obtener la configuración de una integración específica (por ejemplo, WhatsApp que es global).

---

## 4. POST /integrations - Crear nueva integración

**URL**: `POST /api/v1/integrations`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Body Parameters**:
- `name` (string, requerido): Nombre de la integración
- `code` (string, requerido): Código único de la integración
- `type` (string, requerido): Tipo de integración (whatsapp, shopify, mercado_libre)
- `category` (string, requerido): Categoría (internal, external)
- `business_id` (int, opcional): ID del business (NULL para integraciones globales como WhatsApp)
- `is_active` (bool, opcional): Si está activa (por defecto: true)
- `is_default` (bool, opcional): Si es la integración por defecto (por defecto: false)
- `config` (object, opcional): Configuración flexible en JSON (no contiene información sensible)
- `credentials` (object, opcional): Credenciales que se encriptarán automáticamente
- `description` (string, opcional): Descripción de la integración

**Ejemplo Request**:
```json
{
  "name": "WhatsApp Principal",
  "code": "whatsapp_platform",
  "type": "whatsapp",
  "category": "internal",
  "business_id": null,
  "is_active": true,
  "is_default": true,
  "config": {
    "phone_number_id": "123456789",
    "webhook_url": "https://api.example.com/webhooks/whatsapp",
    "template_language": "es",
    "default_country_code": "+57",
    "api_version": "v18.0"
  },
  "credentials": {
    "access_token": "EAAxxxxxxxxxxxx"
  },
  "description": "Integración principal de WhatsApp para toda la plataforma"
}
```

**Ejemplo Request (Shopify por business)**:
```json
{
  "name": "Shopify Store Principal",
  "code": "shopify_main",
  "type": "shopify",
  "category": "external",
  "business_id": 16,
  "is_active": true,
  "is_default": true,
  "config": {
    "store_name": "mi-tienda",
    "api_version": "2024-01",
    "webhook_url": "https://api.example.com/webhooks/shopify",
    "timezone": "America/Bogota"
  },
  "credentials": {
    "access_token": "shpat_xxxxxxxxxxxx",
    "api_secret": "secret_xxxxxxxxxxxx"
  },
  "description": "Integración de Shopify para el business 16"
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Integración creada exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "type": "whatsapp",
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "123456789",
      "webhook_url": "https://api.example.com/webhooks/whatsapp",
      "template_language": "es",
      "default_country_code": "+57",
      "api_version": "v18.0"
    },
    "description": "Integración principal de WhatsApp para toda la plataforma",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "Key: 'CreateIntegrationRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

**Response 409 Conflict**:
```json
{
  "success": false,
  "message": "Error al crear integración",
  "error": "ya existe una integración con el código 'whatsapp_platform'"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden crear integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

**Notas**:
- Las credenciales se encriptan automáticamente antes de guardarse
- WhatsApp debe tener `business_id: null` (es global)
- El código debe ser único por business (o global si `business_id` es null)

---

## 5. PUT /integrations/:id - Actualizar integración

**URL**: `PUT /api/v1/integrations/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Body Parameters** (todos opcionales):
- `name` (string, opcional): Nuevo nombre
- `code` (string, opcional): Nuevo código
- `is_active` (bool, opcional): Estado activo
- `is_default` (bool, opcional): Si es default
- `config` (object, opcional): Nueva configuración
- `credentials` (object, opcional): Nuevas credenciales (se encriptarán)
- `description` (string, opcional): Nueva descripción

**Ejemplo Request**:
```json
{
  "name": "WhatsApp Actualizado",
  "is_active": false,
  "config": {
    "phone_number_id": "987654321",
    "template_language": "en"
  },
  "credentials": {
    "access_token": "EAAyyyyyyyyyyyy"
  },
  "description": "Integración actualizada"
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración actualizada exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Actualizado",
    "code": "whatsapp_platform",
    "type": "whatsapp",
    "category": "internal",
    "business_id": null,
    "is_active": false,
    "is_default": true,
    "config": {
      "phone_number_id": "987654321",
      "template_language": "en"
    },
    "description": "Integración actualizada",
    "created_by_id": 1,
    "updated_by_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z"
  }
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al actualizar integración",
  "error": "integración no encontrada"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden actualizar integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

**Notas**:
- Solo se actualizan los campos enviados
- Si se envía `is_default: true`, automáticamente se desmarca la otra integración del mismo tipo como default
- Las credenciales se encriptan automáticamente si se actualizan

---

## 6. DELETE /integrations/:id - Eliminar integración

**URL**: `DELETE /api/v1/integrations/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
DELETE /api/v1/integrations/2
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración eliminada exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al eliminar integración",
  "error": "integración no encontrada"
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Error al eliminar integración",
  "error": "no se puede eliminar la integración de WhatsApp. Solo se puede desactivar"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden eliminar integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

**Notas**:
- No se puede eliminar la integración de WhatsApp (solo desactivar)
- Al eliminar, se eliminan también las credenciales encriptadas

---

## 7. POST /integrations/:id/test - Probar conexión de integración

**URL**: `POST /api/v1/integrations/:id/test`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
POST /api/v1/integrations/1/test
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Conexión probada exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Error al probar integración",
  "error": "access_token requerido en las credenciales"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al probar integración",
  "error": "integración no encontrada"
}
```

**Response 500 Internal Server Error**:
```json
{
  "success": false,
  "message": "Error al probar integración",
  "error": "no hay tester registrado para tipo whatsapp"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden probar integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

**Notas**:
- Este endpoint usa el tester registrado para cada tipo de integración
- Si no hay tester registrado, hace una validación básica de credenciales
- No modifica la integración, solo prueba la conexión

---

## 8. PUT /integrations/:id/activate - Activar integración

**URL**: `PUT /api/v1/integrations/:id/activate`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
PUT /api/v1/integrations/1/activate
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración activada exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al activar integración",
  "error": "integración no encontrada"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden activar integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

---

## 9. PUT /integrations/:id/deactivate - Desactivar integración

**URL**: `PUT /api/v1/integrations/:id/deactivate`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
PUT /api/v1/integrations/1/deactivate
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración desactivada exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al desactivar integración",
  "error": "integración no encontrada"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden desactivar integraciones",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

---

## 10. PUT /integrations/:id/set-default - Marcar como integración por defecto

**URL**: `PUT /api/v1/integrations/:id/set-default`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID de la integración

**Ejemplo Request**:
```
PUT /api/v1/integrations/2/set-default
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Integración marcada como default exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Error al marcar integración como default",
  "error": "integración no encontrada"
}
```

**Response 403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden marcar integraciones como default",
  "error": "permisos insuficientes"
}
```

**Permisos**: Requiere ser Super Admin

**Notas**:
- Al marcar una integración como default, automáticamente se desmarcan las demás del mismo tipo y business
- Solo puede haber una integración default por tipo y business

---

## Tipos de Integración

### WhatsApp
- **Tipo**: `whatsapp`
- **Categoría**: `internal`
- **Business ID**: `null` (siempre global, una sola para toda la plataforma)
- **Config requerida**:
  ```json
  {
    "phone_number_id": "123456789",
    "webhook_url": "https://api.example.com/webhooks/whatsapp",
    "template_language": "es"
  }
  ```
- **Credenciales requeridas**:
  ```json
  {
    "access_token": "EAAxxxxxxxxxxxx"
  }
  ```

### Shopify
- **Tipo**: `shopify`
- **Categoría**: `external`
- **Business ID**: Requerido (puede haber múltiples por business)
- **Config requerida**:
  ```json
  {
    "store_name": "mi-tienda",
    "api_version": "2024-01",
    "timezone": "America/Bogota"
  }
  ```
- **Credenciales requeridas**:
  ```json
  {
    "access_token": "shpat_xxxxxxxxxxxx",
    "api_secret": "secret_xxxxxxxxxxxx"
  }
  ```

### Mercado Libre
- **Tipo**: `mercado_libre`
- **Categoría**: `external`
- **Business ID**: Requerido (puede haber múltiples por business)
- **Config requerida**:
  ```json
  {
    "app_id": "123456789",
    "redirect_uri": "https://api.example.com/oauth/mercadolibre/callback",
    "country": "CO"
  }
  ```
- **Credenciales requeridas**:
  ```json
  {
    "client_id": "123456789",
    "client_secret": "secret_xxxxxxxxxxxx",
    "refresh_token": "TG-xxxxxxxxxxxx"
  }
  ```

---

## Reglas de Negocio

1. **WhatsApp es global**: Solo puede haber UNA integración de tipo `whatsapp` con `business_id = null`
2. **Códigos únicos**: El código debe ser único por business (o global si `business_id` es null)
3. **Encriptación automática**: Las credenciales se encriptan automáticamente usando AES-256-GCM
4. **Default único**: Solo puede haber una integración default por tipo y business
5. **Eliminación de WhatsApp**: No se puede eliminar WhatsApp, solo desactivar
6. **Testing**: Cada integración debe registrar su tester para poder probar la conexión

---

## Seguridad

- **Autenticación**: Todos los endpoints requieren JWT válido
- **Autorización**: Solo Super Admins pueden crear, actualizar, eliminar y probar integraciones
- **Encriptación**: Las credenciales se encriptan automáticamente antes de guardarse
- **Exposición**: Las credenciales NUNCA se exponen en las respuestas HTTP

---

## Ejemplos de Uso Común

### Obtener configuración de WhatsApp (para uso interno)
```bash
GET /api/v1/integrations/type/whatsapp
Authorization: Bearer {token}
```

### Crear integración de Shopify para un business
```bash
POST /api/v1/integrations
Authorization: Bearer {token}
Content-Type: application/json

{
  "name": "Shopify Store 1",
  "code": "shopify_main",
  "type": "shopify",
  "category": "external",
  "business_id": 16,
  "is_active": true,
  "is_default": true,
  "config": {
    "store_name": "mi-tienda",
    "api_version": "2024-01"
  },
  "credentials": {
    "access_token": "shpat_xxxxxxxxxxxx"
  }
}
```

### Probar conexión antes de activar
```bash
POST /api/v1/integrations/2/test
Authorization: Bearer {token}
```

### Desactivar integración temporalmente
```bash
PUT /api/v1/integrations/2/deactivate
Authorization: Bearer {token}
```

