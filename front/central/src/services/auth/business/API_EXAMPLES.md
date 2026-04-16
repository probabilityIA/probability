# Ejemplos de API - Business Module

Base URL: `/api/v1`

## 1. GET /businesses - Obtener lista de negocios

**URL**: `GET /api/v1/businesses`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (default: 1)
- `per_page` (int, opcional): Elementos por página (default: 10, max: 100)
- `name` (string, opcional): Filtrar por nombre de negocio
- `business_type_id` (int, opcional): Filtrar por tipo de negocio
- `is_active` (boolean, opcional): Filtrar por estado activo/inactivo

**Ejemplo Request**:
```
GET /api/v1/businesses?page=1&per_page=10&name=restaurante&is_active=true
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Negocios obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "name": "Restaurante El Buen Sabor",
      "description": "Restaurante de comida tradicional",
      "address": "Calle 123 #45-67",
      "phone": "+573001234567",
      "email": "contacto@elbuensabor.com",
      "website": "https://elbuensabor.com",
      "logo_url": "https://s3.amazonaws.com/bucket/logo.png",
      "primary_color": "#FF5733",
      "secondary_color": "#33FF57",
      "tertiary_color": "#3357FF",
      "quaternary_color": "#FF33F5",
      "navbar_image_url": "https://s3.amazonaws.com/bucket/navbar.png",
      "is_active": true,
      "business_type_id": 1,
      "business_type": "Restaurante",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-20T14:20:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 25,
    "last_page": 3,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## 2. GET /businesses/:id - Obtener negocio por ID

**URL**: `GET /api/v1/businesses/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del negocio

**Ejemplo Request**:
```
GET /api/v1/businesses/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Negocio obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "Restaurante El Buen Sabor",
    "code": "REST001",
    "business_type": {
      "id": 1,
      "name": "Restaurante",
      "code": "RESTAURANT",
      "description": "Negocio de restaurante",
      "icon": "restaurant",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "timezone": "America/Bogota",
    "address": "Calle 123 #45-67",
    "description": "Restaurante de comida tradicional",
    "logo_url": "https://s3.amazonaws.com/bucket/logo.png",
    "primary_color": "#FF5733",
    "secondary_color": "#33FF57",
    "tertiary_color": "#3357FF",
    "quaternary_color": "#FF33F5",
    "navbar_image_url": "https://s3.amazonaws.com/bucket/navbar.png",
    "custom_domain": "elbuensabor.com",
    "is_active": true,
    "enable_delivery": true,
    "enable_pickup": true,
    "enable_reservations": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-20T14:20:00Z"
  }
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "error": "not_found",
  "message": "Negocio no encontrado"
}
```

---

## 3. POST /businesses - Crear negocio

**URL**: `POST /api/v1/businesses`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Form Data Parameters**:
- `name` (string, requerido): Nombre del negocio
- `code` (string, requerido): Código del negocio
- `business_type_id` (int, requerido): ID del tipo de negocio
- `timezone` (string, opcional): Zona horaria
- `address` (string, opcional): Dirección
- `description` (string, opcional): Descripción
- `logo_file` (file, opcional): Logo del negocio
- `primary_color` (string, opcional): Color primario (hex)
- `secondary_color` (string, opcional): Color secundario (hex)
- `tertiary_color` (string, opcional): Color terciario (hex)
- `quaternary_color` (string, opcional): Color cuaternario (hex)
- `navbar_image_file` (file, opcional): Imagen del navbar
- `custom_domain` (string, opcional): Dominio personalizado
- `is_active` (boolean, opcional): ¿Activo? (default: true)
- `enable_delivery` (boolean, opcional): Habilitar delivery (default: false)
- `enable_pickup` (boolean, opcional): Habilitar pickup (default: false)
- `enable_reservations` (boolean, opcional): Habilitar reservas (default: false)

**Ejemplo Request** (multipart/form-data):
```
POST /api/v1/businesses
Content-Type: multipart/form-data

name=Restaurante El Buen Sabor
code=REST001
business_type_id=1
timezone=America/Bogota
address=Calle 123 #45-67
description=Restaurante de comida tradicional
primary_color=#FF5733
secondary_color=#33FF57
enable_delivery=true
enable_pickup=true
logo_file=[archivo binario]
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Negocio creado exitosamente",
  "data": {
    "id": 1,
    "name": "Restaurante El Buen Sabor",
    "description": "Restaurante de comida tradicional",
    "address": "Calle 123 #45-67",
    "phone": "",
    "email": "",
    "website": "",
    "logo_url": "https://s3.amazonaws.com/bucket/logo.png",
    "primary_color": "#FF5733",
    "secondary_color": "#33FF57",
    "tertiary_color": "",
    "quaternary_color": "",
    "navbar_image_url": "",
    "is_active": true,
    "business_type_id": 1,
    "business_type": "Restaurante",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "error": "invalid_request",
  "message": "Datos de entrada inválidos: name is required"
}
```

**Response 409 Conflict** (código duplicado):
```json
{
  "success": false,
  "error": "code_already_exists",
  "message": "El código del negocio ya está en uso"
}
```

---

## 4. PUT /businesses/:id - Actualizar negocio

**URL**: `PUT /api/v1/businesses/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Path Parameters**:
- `id` (int, requerido): ID del negocio

**Form Data Parameters** (todos opcionales):
- `name` (string): Nombre del negocio
- `code` (string): Código del negocio
- `business_type_id` (int): ID del tipo de negocio
- `timezone` (string): Zona horaria
- `address` (string): Dirección
- `description` (string): Descripción
- `logo_file` (file): Logo del negocio
- `primary_color` (string): Color primario
- `secondary_color` (string): Color secundario
- `tertiary_color` (string): Color terciario
- `quaternary_color` (string): Color cuaternario
- `navbar_image_file` (file): Imagen del navbar
- `custom_domain` (string): Dominio personalizado
- `is_active` (boolean): ¿Activo?
- `enable_delivery` (boolean): Habilitar delivery
- `enable_pickup` (boolean): Habilitar pickup
- `enable_reservations` (boolean): Habilitar reservas

**Ejemplo Request**:
```
PUT /api/v1/businesses/1
Content-Type: multipart/form-data

name=Restaurante El Buen Sabor Actualizado
description=Nueva descripción
enable_reservations=true
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Negocio actualizado exitosamente",
  "data": {
    "id": 1,
    "name": "Restaurante El Buen Sabor Actualizado",
    "description": "Nueva descripción",
    "address": "Calle 123 #45-67",
    "phone": "+573001234567",
    "email": "contacto@elbuensabor.com",
    "website": "https://elbuensabor.com",
    "logo_url": "https://s3.amazonaws.com/bucket/logo.png",
    "primary_color": "#FF5733",
    "secondary_color": "#33FF57",
    "tertiary_color": "#3357FF",
    "quaternary_color": "#FF33F5",
    "navbar_image_url": "https://s3.amazonaws.com/bucket/navbar.png",
    "is_active": true,
    "business_type_id": 1,
    "business_type": "Restaurante",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-21T15:45:00Z"
  }
}
```

---

## 5. DELETE /businesses/:id - Eliminar negocio

**URL**: `DELETE /api/v1/businesses/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del negocio

**Ejemplo Request**:
```
DELETE /api/v1/businesses/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Negocio eliminado exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "error": "not_found",
  "message": "Negocio no encontrado"
}
```

---

## 6. GET /businesses/configured-resources - Obtener recursos configurados de negocios

**URL**: `GET /api/v1/businesses/configured-resources`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (default: 1)
- `per_page` (int, opcional): Elementos por página (default: 10, max: 100)
- `business_id` (int, opcional): Filtrar por ID de business (solo super admin)
- `business_type_id` (int, opcional): Filtrar por ID de tipo de business (solo super admin)

**Nota**: Usuarios normales solo ven su propio business (del token). Super admins pueden filtrar.

**Ejemplo Request**:
```
GET /api/v1/businesses/configured-resources?page=1&per_page=10&business_type_id=1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Businesses con recursos configurados obtenidos exitosamente",
  "data": [
    {
      "business_id": 1,
      "resources": [
        {
          "resource_id": 1,
          "resource_name": "users",
          "is_active": true
        },
        {
          "resource_id": 2,
          "resource_name": "orders",
          "is_active": true
        },
        {
          "resource_id": 3,
          "resource_name": "inventory",
          "is_active": false
        }
      ],
      "total": 5,
      "active": 3,
      "inactive": 2
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 15,
    "last_page": 2,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## 7. GET /businesses/:id/configured-resources - Obtener recursos configurados de un negocio específico

**URL**: `GET /api/v1/businesses/:id/configured-resources`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del business

**Nota**: Usuarios normales solo pueden acceder a su propio business. Super admins pueden acceder a cualquier business.

**Ejemplo Request**:
```
GET /api/v1/businesses/1/configured-resources
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Configuración de recursos del business obtenida exitosamente",
  "data": {
    "business_id": 1,
    "resources": [
      {
        "resource_id": 1,
        "resource_name": "users",
        "is_active": true
      },
      {
        "resource_id": 2,
        "resource_name": "orders",
        "is_active": true
      },
      {
        "resource_id": 3,
        "resource_name": "inventory",
        "is_active": false
      },
      {
        "resource_id": 4,
        "resource_name": "reports",
        "is_active": true
      },
      {
        "resource_id": 5,
        "resource_name": "analytics",
        "is_active": false
      }
    ],
    "total": 5,
    "active": 3,
    "inactive": 2
  }
}
```

**Response 403 Forbidden** (usuario normal intentando acceder a otro business):
```json
{
  "success": false,
  "message": "Sin permisos para acceder a esta configuración",
  "error": "Solo puedes acceder a la configuración de tu propio business"
}
```

---

## 8. PUT /businesses/configured-resources/:resource_id/activate - Activar recurso de business

**URL**: `PUT /api/v1/businesses/configured-resources/:resource_id/activate`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `resource_id` (int, requerido): ID del recurso

**Query Parameters**:
- `business_id` (int, requerido para super admin): ID del business

**Nota**: Usuarios normales usan su business_id del token. Super admins deben enviar business_id en query.

**Ejemplo Request** (Super Admin):
```
PUT /api/v1/businesses/configured-resources/3/activate?business_id=1
```

**Ejemplo Request** (Usuario Normal):
```
PUT /api/v1/businesses/configured-resources/3/activate
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recurso activado exitosamente"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "business no encontrado",
  "error": "business no encontrado"
}
```

---

## 9. PUT /businesses/configured-resources/:resource_id/deactivate - Desactivar recurso de business

**URL**: `PUT /api/v1/businesses/configured-resources/:resource_id/deactivate`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `resource_id` (int, requerido): ID del recurso

**Query Parameters**:
- `business_id` (int, requerido para super admin): ID del business

**Ejemplo Request** (Super Admin):
```
PUT /api/v1/businesses/configured-resources/3/deactivate?business_id=1
```

**Ejemplo Request** (Usuario Normal):
```
PUT /api/v1/businesses/configured-resources/3/deactivate
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recurso desactivado exitosamente"
}
```

---

## Business Types Endpoints

## 10. GET /business-types - Obtener tipos de negocio

**URL**: `GET /api/v1/business-types`

**Headers**:
```
Authorization: Bearer {token}
```

**Ejemplo Request**:
```
GET /api/v1/business-types
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Tipos de negocio obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "name": "Restaurante",
      "code": "RESTAURANT",
      "description": "Negocio de restaurante",
      "icon": "restaurant",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Cafetería",
      "code": "CAFE",
      "description": "Negocio de cafetería",
      "icon": "cafe",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## 11. GET /business-types/:id - Obtener tipo de negocio por ID

**URL**: `GET /api/v1/business-types/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del tipo de negocio

**Ejemplo Request**:
```
GET /api/v1/business-types/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Tipo de negocio obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "Restaurante",
    "code": "RESTAURANT",
    "description": "Negocio de restaurante",
    "icon": "restaurant",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## 12. POST /business-types - Crear tipo de negocio

**URL**: `POST /api/v1/business-types`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Bar",
  "code": "BAR",
  "description": "Negocio de bar",
  "icon": "bar",
  "is_active": true
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Tipo de negocio creado exitosamente",
  "data": {
    "id": 3,
    "name": "Bar",
    "code": "BAR",
    "description": "Negocio de bar",
    "icon": "bar",
    "is_active": true,
    "created_at": "2024-01-22T10:00:00Z",
    "updated_at": "2024-01-22T10:00:00Z"
  }
}
```

---

## 13. PUT /business-types/:id - Actualizar tipo de negocio

**URL**: `PUT /api/v1/business-types/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del tipo de negocio

**Request Body**:
```json
{
  "name": "Bar Actualizado",
  "description": "Nueva descripción del bar",
  "is_active": false
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Tipo de negocio actualizado exitosamente",
  "data": {
    "id": 3,
    "name": "Bar Actualizado",
    "code": "BAR",
    "description": "Nueva descripción del bar",
    "icon": "bar",
    "is_active": false,
    "created_at": "2024-01-22T10:00:00Z",
    "updated_at": "2024-01-22T11:30:00Z"
  }
}
```

---

## 14. DELETE /business-types/:id - Eliminar tipo de negocio

**URL**: `DELETE /api/v1/business-types/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del tipo de negocio

**Ejemplo Request**:
```
DELETE /api/v1/business-types/3
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Tipo de negocio eliminado exitosamente"
}
```

---

## Respuestas de Error Comunes

**400 Bad Request**:
```json
{
  "success": false,
  "error": "invalid_request",
  "message": "Datos de entrada inválidos"
}
```

**401 Unauthorized**:
```json
{
  "success": false,
  "error": "unauthorized",
  "message": "Token de acceso requerido"
}
```

**403 Forbidden**:
```json
{
  "success": false,
  "error": "access_denied",
  "message": "No tienes permisos para acceder a este endpoint"
}
```

**404 Not Found**:
```json
{
  "success": false,
  "error": "not_found",
  "message": "Recurso no encontrado"
}
```

**500 Internal Server Error**:
```json
{
  "success": false,
  "error": "internal_error",
  "message": "Error interno del servidor"
}
```

