# Documentación de APIs - Core Integrations

Este documento describe todos los endpoints disponibles en el módulo core de integraciones.

**Base URL:** `/api/v1`

**Autenticación:** Todos los endpoints requieren un token JWT válido en el header `Authorization: Bearer <token>`

---

## Tabla de Contenidos

### Integrations
1. [Listar Integraciones](#1-listar-integraciones)
2. [Obtener Integración por ID](#2-obtener-integración-por-id)
3. [Obtener Integración por Tipo](#3-obtener-integración-por-tipo)
4. [Crear Integración](#4-crear-integración)
5. [Actualizar Integración](#5-actualizar-integración)
6. [Eliminar Integración](#6-eliminar-integración)
7. [Probar Integración](#7-probar-integración)
8. [Activar Integración](#8-activar-integración)
9. [Desactivar Integración](#9-desactivar-integración)
10. [Marcar como Default](#10-marcar-como-default)

### Integration Types
11. [Listar Tipos de Integración](#11-listar-tipos-de-integración)
12. [Listar Tipos de Integración Activos](#12-listar-tipos-de-integración-activos)
13. [Obtener Tipo de Integración por ID](#13-obtener-tipo-de-integración-por-id)
14. [Obtener Tipo de Integración por Código](#14-obtener-tipo-de-integración-por-código)
15. [Crear Tipo de Integración](#15-crear-tipo-de-integración)
16. [Actualizar Tipo de Integración](#16-actualizar-tipo-de-integración)
17. [Eliminar Tipo de Integración](#17-eliminar-tipo-de-integración)

---

## Integrations

### 1. Listar Integraciones

Obtiene una lista paginada de integraciones con filtros opcionales.

**URL:** `GET /integrations`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `page` | int | No | Número de página | `1` |
| `page_size` | int | No | Tamaño de página (máx 100) | `10` |
| `integration_type_id` | int | No | Filtrar por ID del tipo de integración | `1` |
| `integration_type_code` | string | No | Filtrar por código del tipo de integración | `whatsapp` |
| `category` | string | No | Filtrar por categoría (`internal` o `external`) | `internal` |
| `business_id` | int | No | Filtrar por business ID | `16` |
| `is_active` | bool | No | Filtrar por estado activo | `true` |
| `search` | string | No | Buscar por nombre o código | `whatsapp` |

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integrations?page=1&page_size=10&integration_type_code=whatsapp&is_active=true" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integraciones obtenidas exitosamente",
  "data": [
    {
      "id": 1,
      "name": "WhatsApp Principal",
      "code": "whatsapp_platform",
      "integration_type_id": 1,
      "integration_type": {
        "id": 1,
        "name": "WhatsApp",
        "code": "whatsapp"
      },
      "category": "internal",
      "business_id": null,
      "is_active": true,
      "is_default": true,
      "config": {
        "phone_number_id": "1234567890"
      },
      "description": "Integración principal de WhatsApp",
      "created_by_id": 1,
      "updated_by_id": null,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "page_size": 10,
  "total_pages": 3
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `400 Bad Request`: Parámetros de consulta inválidos
- `401 Unauthorized`: Token inválido o faltante
- `500 Internal Server Error`: Error interno del servidor

---

### 2. Obtener Integración por ID

Obtiene los detalles de una integración específica por su ID.

**URL:** `GET /integrations/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integrations/1" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración obtenida exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "integration_type_id": 1,
    "integration_type": {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp"
    },
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "1234567890"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

### 3. Obtener Integración por Tipo

Obtiene una integración activa por código de tipo de integración y business ID opcional.

**URL:** `GET /integrations/type/:type`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `type` | string | Código del tipo de integración | `whatsapp` |

**Query Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `business_id` | int | No | ID del business (si no se envía, busca integraciones globales) | `16` |

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integrations/type/whatsapp?business_id=16" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración obtenida exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "integration_type_id": 1,
    "integration_type": {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp"
    },
    "category": "internal",
    "business_id": 16,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "1234567890"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

### 4. Crear Integración

Crea una nueva integración en el sistema. **Solo super administradores pueden crear integraciones.**

**URL:** `POST /integrations`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `name` | string | Sí | Nombre de la integración | `"WhatsApp Principal"` |
| `code` | string | Sí | Código único de la integración | `"whatsapp_platform"` |
| `integration_type_id` | int | Sí | ID del tipo de integración | `1` |
| `category` | string | Sí | Categoría (`internal` o `external`) | `"internal"` |
| `business_id` | int | No | ID del business (null para integraciones globales) | `16` |
| `is_active` | bool | No | Estado activo | `true` |
| `is_default` | bool | No | Marcar como default | `true` |
| `config` | object | No | Configuración flexible (JSON) | `{"phone_number_id": "1234567890"}` |
| `credentials` | object | No | Credenciales (se encriptarán automáticamente) | `{"access_token": "token123"}` |
| `description` | string | No | Descripción | `"Integración principal de WhatsApp"` |

**Ejemplo de Request:**
```bash
curl -X POST "https://api.example.com/api/v1/integrations" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "integration_type_id": 1,
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "1234567890"
    },
    "credentials": {
      "access_token": "EAAxxxxxxxxxxxxx",
      "phone_number_id": "1234567890"
    },
    "description": "Integración principal de WhatsApp"
  }'
```

**Ejemplo de Response (201 Created):**
```json
{
  "success": true,
  "message": "Integración creada exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Principal",
    "code": "whatsapp_platform",
    "integration_type_id": 1,
    "integration_type": {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp"
    },
    "category": "internal",
    "business_id": null,
    "is_active": true,
    "is_default": true,
    "config": {
      "phone_number_id": "1234567890"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": null,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Códigos de Estado:**
- `201 Created`: Integración creada exitosamente
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token inválido o faltante
- `403 Forbidden`: Solo super administradores pueden crear integraciones
- `409 Conflict`: Ya existe una integración con ese código
- `500 Internal Server Error`: Error interno del servidor

---

### 5. Actualizar Integración

Actualiza una integración existente. Todos los campos son opcionales.

**URL:** `PUT /integrations/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Body Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `name` | string | No | Nombre de la integración | `"WhatsApp Actualizado"` |
| `code` | string | No | Código único de la integración | `"whatsapp_platform"` |
| `integration_type_id` | int | No | ID del tipo de integración | `1` |
| `is_active` | bool | No | Estado activo | `true` |
| `is_default` | bool | No | Marcar como default | `true` |
| `config` | object | No | Configuración flexible (JSON) | `{"phone_number_id": "1234567890"}` |
| `credentials` | object | No | Credenciales (se encriptarán automáticamente) | `{"access_token": "new_token"}` |
| `description` | string | No | Descripción | `"Nueva descripción"` |

**Ejemplo de Request:**
```bash
curl -X PUT "https://api.example.com/api/v1/integrations/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp Actualizado",
    "is_active": false,
    "config": {
      "phone_number_id": "9876543210"
    }
  }'
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración actualizada exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Actualizado",
    "code": "whatsapp_platform",
    "integration_type_id": 1,
    "integration_type": {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp"
    },
    "category": "internal",
    "business_id": null,
    "is_active": false,
    "is_default": true,
    "config": {
      "phone_number_id": "9876543210"
    },
    "description": "Integración principal de WhatsApp",
    "created_by_id": 1,
    "updated_by_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Integración actualizada exitosamente
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `409 Conflict`: Ya existe otra integración con ese código
- `500 Internal Server Error`: Error interno del servidor

---

### 6. Eliminar Integración

Elimina una integración del sistema. **Nota:** No se puede eliminar la integración de WhatsApp si es la única de ese tipo.

**URL:** `DELETE /integrations/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X DELETE "https://api.example.com/api/v1/integrations/1" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración eliminada exitosamente"
}
```

**Códigos de Estado:**
- `200 OK`: Integración eliminada exitosamente
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `400 Bad Request`: No se puede eliminar (ej: única integración de WhatsApp)
- `500 Internal Server Error`: Error interno del servidor

---

### 7. Probar Integración

Prueba la conexión de una integración usando su tester registrado.

**URL:** `POST /integrations/:id/test`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X POST "https://api.example.com/api/v1/integrations/1/test" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Test de conexión exitoso"
}
```

**Ejemplo de Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Test de conexión falló",
  "error": "access_token inválido o expirado"
}
```

**Códigos de Estado:**
- `200 OK`: Test de conexión exitoso
- `400 Bad Request`: Test de conexión falló
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

### 8. Activar Integración

Activa una integración que estaba desactivada.

**URL:** `PUT /integrations/:id/activate`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X PUT "https://api.example.com/api/v1/integrations/1/activate" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración activada exitosamente"
}
```

**Códigos de Estado:**
- `200 OK`: Integración activada exitosamente
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

### 9. Desactivar Integración

Desactiva una integración activa.

**URL:** `PUT /integrations/:id/deactivate`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X PUT "https://api.example.com/api/v1/integrations/1/deactivate" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración desactivada exitosamente"
}
```

**Códigos de Estado:**
- `200 OK`: Integración desactivada exitosamente
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

### 10. Marcar como Default

Marca una integración como default para su tipo y business. Esto desmarcará automáticamente otras integraciones del mismo tipo.

**URL:** `PUT /integrations/:id/set-default`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID de la integración | `1` |

**Ejemplo de Request:**
```bash
curl -X PUT "https://api.example.com/api/v1/integrations/1/set-default" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Integración marcada como default exitosamente"
}
```

**Códigos de Estado:**
- `200 OK`: Integración marcada como default exitosamente
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Integración no encontrada
- `500 Internal Server Error`: Error interno del servidor

---

## Integration Types

### 11. Listar Tipos de Integración

Obtiene todos los tipos de integración disponibles en el sistema.

**URL:** `GET /integration-types`

**Headers:**
```
Authorization: Bearer <token>
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integration-types" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipos de integración obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp",
      "description": "Integración con WhatsApp Cloud API",
      "icon": "whatsapp-icon",
      "category": "internal",
      "is_active": true,
      "config_schema": {
        "type": "object",
        "properties": {
          "phone_number_id": {
            "type": "string",
            "required": true
          }
        }
      },
      "credentials_schema": {
        "type": "object",
        "properties": {
          "access_token": {
            "type": "string",
            "required": true
          }
        }
      },
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `500 Internal Server Error`: Error interno del servidor

---

### 12. Listar Tipos de Integración Activos

Obtiene solo los tipos de integración que están activos.

**URL:** `GET /integration-types/active`

**Headers:**
```
Authorization: Bearer <token>
```

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integration-types/active" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipos de integración activos obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "name": "WhatsApp",
      "code": "whatsapp",
      "description": "Integración con WhatsApp Cloud API",
      "icon": "whatsapp-icon",
      "category": "internal",
      "is_active": true,
      "config_schema": {},
      "credentials_schema": {},
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `500 Internal Server Error`: Error interno del servidor

---

### 13. Obtener Tipo de Integración por ID

Obtiene los detalles de un tipo de integración específico por su ID.

**URL:** `GET /integration-types/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID del tipo de integración | `1` |

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integration-types/1" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipo de integración obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp",
    "code": "whatsapp",
    "description": "Integración con WhatsApp Cloud API",
    "icon": "whatsapp-icon",
    "category": "internal",
    "is_active": true,
    "config_schema": {
      "type": "object",
      "properties": {
        "phone_number_id": {
          "type": "string",
          "required": true
        }
      }
    },
    "credentials_schema": {
      "type": "object",
      "properties": {
        "access_token": {
          "type": "string",
          "required": true
        }
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Tipo de integración no encontrado
- `500 Internal Server Error`: Error interno del servidor

---

### 14. Obtener Tipo de Integración por Código

Obtiene los detalles de un tipo de integración específico por su código.

**URL:** `GET /integration-types/code/:code`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `code` | string | Código del tipo de integración | `whatsapp` |

**Ejemplo de Request:**
```bash
curl -X GET "https://api.example.com/api/v1/integration-types/code/whatsapp" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipo de integración obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp",
    "code": "whatsapp",
    "description": "Integración con WhatsApp Cloud API",
    "icon": "whatsapp-icon",
    "category": "internal",
    "is_active": true,
    "config_schema": {},
    "credentials_schema": {},
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Operación exitosa
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Tipo de integración no encontrado
- `500 Internal Server Error`: Error interno del servidor

---

### 15. Crear Tipo de Integración

Crea un nuevo tipo de integración en el sistema.

**URL:** `POST /integration-types`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Body Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `name` | string | Sí | Nombre del tipo de integración | `"WhatsApp"` |
| `code` | string | No | Código único (se genera automáticamente si no se proporciona) | `"whatsapp"` |
| `description` | string | No | Descripción | `"Integración con WhatsApp Cloud API"` |
| `icon` | string | No | Nombre del icono | `"whatsapp-icon"` |
| `category` | string | Sí | Categoría (`internal` o `external`) | `"internal"` |
| `is_active` | bool | No | Estado activo | `true` |
| `config_schema` | object | No | JSON Schema para validar configuración | `{"type": "object", "properties": {...}}` |
| `credentials_schema` | object | No | JSON Schema para validar credenciales | `{"type": "object", "properties": {...}}` |

**Ejemplo de Request:**
```bash
curl -X POST "https://api.example.com/api/v1/integration-types" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Shopify",
    "code": "shopify",
    "description": "Integración con Shopify API",
    "icon": "shopify-icon",
    "category": "external",
    "is_active": true,
    "config_schema": {
      "type": "object",
      "properties": {
        "store_url": {
          "type": "string",
          "required": true
        }
      }
    },
    "credentials_schema": {
      "type": "object",
      "properties": {
        "api_key": {
          "type": "string",
          "required": true
        },
        "api_secret": {
          "type": "string",
          "required": true
        }
      }
    }
  }'
```

**Ejemplo de Response (201 Created):**
```json
{
  "success": true,
  "message": "Tipo de integración creado exitosamente",
  "data": {
    "id": 2,
    "name": "Shopify",
    "code": "shopify",
    "description": "Integración con Shopify API",
    "icon": "shopify-icon",
    "category": "external",
    "is_active": true,
    "config_schema": {
      "type": "object",
      "properties": {
        "store_url": {
          "type": "string",
          "required": true
        }
      }
    },
    "credentials_schema": {
      "type": "object",
      "properties": {
        "api_key": {
          "type": "string",
          "required": true
        },
        "api_secret": {
          "type": "string",
          "required": true
        }
      }
    },
    "created_at": "2024-01-15T12:00:00Z",
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

**Códigos de Estado:**
- `201 Created`: Tipo de integración creado exitosamente
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token inválido o faltante
- `409 Conflict`: Ya existe un tipo de integración con ese nombre o código
- `500 Internal Server Error`: Error interno del servidor

---

### 16. Actualizar Tipo de Integración

Actualiza un tipo de integración existente. Todos los campos son opcionales.

**URL:** `PUT /integration-types/:id`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID del tipo de integración | `1` |

**Body Parameters:**
| Parámetro | Tipo | Requerido | Descripción | Ejemplo |
|-----------|------|-----------|-------------|---------|
| `name` | string | No | Nombre del tipo de integración | `"WhatsApp Actualizado"` |
| `code` | string | No | Código único | `"whatsapp"` |
| `description` | string | No | Descripción | `"Nueva descripción"` |
| `icon` | string | No | Nombre del icono | `"whatsapp-icon"` |
| `category` | string | No | Categoría (`internal` o `external`) | `"internal"` |
| `is_active` | bool | No | Estado activo | `true` |
| `config_schema` | object | No | JSON Schema para validar configuración | `{"type": "object"}` |
| `credentials_schema` | object | No | JSON Schema para validar credenciales | `{"type": "object"}` |

**Ejemplo de Request:**
```bash
curl -X PUT "https://api.example.com/api/v1/integration-types/1" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "WhatsApp Actualizado",
    "description": "Nueva descripción actualizada",
    "is_active": false
  }'
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipo de integración actualizado exitosamente",
  "data": {
    "id": 1,
    "name": "WhatsApp Actualizado",
    "code": "whatsapp",
    "description": "Nueva descripción actualizada",
    "icon": "whatsapp-icon",
    "category": "internal",
    "is_active": false,
    "config_schema": {},
    "credentials_schema": {},
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

**Códigos de Estado:**
- `200 OK`: Tipo de integración actualizado exitosamente
- `400 Bad Request`: Datos de entrada inválidos
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Tipo de integración no encontrado
- `409 Conflict`: Ya existe otro tipo de integración con ese nombre o código
- `500 Internal Server Error`: Error interno del servidor

---

### 17. Eliminar Tipo de Integración

Elimina un tipo de integración del sistema. **Nota:** No se puede eliminar si hay integraciones asociadas.

**URL:** `DELETE /integration-types/:id`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
| Parámetro | Tipo | Descripción | Ejemplo |
|-----------|------|-------------|---------|
| `id` | int | ID del tipo de integración | `1` |

**Ejemplo de Request:**
```bash
curl -X DELETE "https://api.example.com/api/v1/integration-types/1" \
  -H "Authorization: Bearer <token>"
```

**Ejemplo de Response (200 OK):**
```json
{
  "success": true,
  "message": "Tipo de integración eliminado exitosamente"
}
```

**Códigos de Estado:**
- `200 OK`: Tipo de integración eliminado exitosamente
- `401 Unauthorized`: Token inválido o faltante
- `404 Not Found`: Tipo de integración no encontrado
- `400 Bad Request`: No se puede eliminar (hay integraciones asociadas)
- `500 Internal Server Error`: Error interno del servidor

---

## Notas Importantes

### Seguridad
- Todas las credenciales se encriptan automáticamente antes de guardarse en la base de datos usando AES-256-GCM
- Las credenciales nunca se exponen en las respuestas de los endpoints
- Solo los super administradores pueden crear integraciones

### Validaciones
- El código de integración debe ser único por business (o global si `business_id` es null)
- No se puede eliminar la integración de WhatsApp si es la única de ese tipo
- No se puede eliminar un tipo de integración si hay integraciones asociadas
- El nombre del tipo de integración debe ser único

### Filtros
- Los filtros en la lista de integraciones se pueden combinar
- Si no se especifica `business_id`, se muestran tanto integraciones globales como por business
- La búsqueda (`search`) busca en los campos `name` y `code`

### Paginación
- Por defecto: `page=1`, `page_size=10`
- Máximo `page_size`: 100
- La respuesta incluye `total`, `page`, `page_size` y `total_pages`

