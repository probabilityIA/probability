# Ejemplos de API - Permissions Module

Base URL: `/api/v1`

## 1. GET /permissions - Obtener todos los permisos

**URL**: `GET /api/v1/permissions`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `business_type_id` (int, opcional): Filtrar por tipo de business (incluye genéricos)
- `name` (string, opcional): Filtrar por nombre de permiso (búsqueda parcial)
- `scope_id` (int, opcional): Filtrar por ID de scope

**Ejemplo Request**:
```
GET /api/v1/permissions?business_type_id=1&name=crear&scope_id=1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Crear usuarios",
      "code": "users:create",
      "description": "Permite crear nuevos usuarios en el sistema",
      "resource": "users",
      "action": "create",
      "resource_id": 3,
      "action_id": 5,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 11,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 2,
      "name": "Editar usuarios",
      "code": "users:update",
      "description": "Permite editar usuarios existentes",
      "resource": "users",
      "action": "update",
      "resource_id": 3,
      "action_id": 6,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    }
  ],
  "total": 2
}
```

**Response 401 Unauthorized**:
```json
{
  "error": "Token de acceso requerido"
}
```

---

## 2. GET /permissions/:id - Obtener permiso por ID

**URL**: `GET /api/v1/permissions/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del permiso

**Ejemplo Request**:
```
GET /api/v1/permissions/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Crear usuarios",
    "code": "users:create",
    "description": "Permite crear nuevos usuarios en el sistema",
    "resource": "users",
    "action": "create",
    "resource_id": 3,
    "action_id": 5,
    "scope_id": 1,
    "scope_name": "Sistema",
    "scope_code": "system",
    "business_type_id": 11,
    "business_type_name": "Propiedad Horizontal"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "error": "ID de permiso inválido"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Permiso no encontrado"
}
```

---

## 3. GET /permissions/scope/:scope_id - Obtener permisos por scope

**URL**: `GET /api/v1/permissions/scope/:scope_id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `scope_id` (int, requerido): ID del scope

**Ejemplo Request**:
```
GET /api/v1/permissions/scope/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Crear usuarios",
      "code": "users:create",
      "description": "Permite crear nuevos usuarios en el sistema",
      "resource": "users",
      "action": "create",
      "resource_id": 3,
      "action_id": 5,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 11,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 3,
      "name": "Eliminar usuarios",
      "code": "users:delete",
      "description": "Permite eliminar usuarios del sistema",
      "resource": "users",
      "action": "delete",
      "resource_id": 3,
      "action_id": 7,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    }
  ],
  "total": 2
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Scope ID de permiso inválido"
}
```

---

## 4. GET /permissions/resource/:resource - Obtener permisos por recurso

**URL**: `GET /api/v1/permissions/resource/:resource`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `resource` (string, requerido): Nombre del recurso (ej: "users", "orders", "inventory")

**Ejemplo Request**:
```
GET /api/v1/permissions/resource/users
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Crear usuarios",
      "code": "users:create",
      "description": "Permite crear nuevos usuarios en el sistema",
      "resource": "users",
      "action": "create",
      "resource_id": 3,
      "action_id": 5,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 11,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 2,
      "name": "Editar usuarios",
      "code": "users:update",
      "description": "Permite editar usuarios existentes",
      "resource": "users",
      "action": "update",
      "resource_id": 3,
      "action_id": 6,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    },
    {
      "id": 3,
      "name": "Eliminar usuarios",
      "code": "users:delete",
      "description": "Permite eliminar usuarios del sistema",
      "resource": "users",
      "action": "delete",
      "resource_id": 3,
      "action_id": 7,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    }
  ],
  "total": 3
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Recurso de permiso inválido"
}
```

---

## 5. POST /permissions - Crear permiso

**URL**: `POST /api/v1/permissions`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Crear usuario",
  "code": "horizontalproperty_createuser",
  "description": "Permite crear nuevos usuarios en el sistema",
  "resource_id": 1,
  "action_id": 1,
  "scope_id": 1,
  "business_type_id": 11
}
```

**Campos del Request**:
- `name` (string, requerido): Nombre del permiso
- `code` (string, opcional): Código del permiso (se genera automáticamente si no se proporciona)
- `description` (string, opcional): Descripción del permiso
- `resource_id` (int, requerido): ID del recurso
- `action_id` (int, requerido): ID de la acción
- `scope_id` (int, requerido): ID del scope
- `business_type_id` (int, opcional): ID del tipo de business (null para genérico)

**Ejemplo Request**:
```
POST /api/v1/permissions
Content-Type: application/json

{
  "name": "Ver reportes",
  "code": "reports:view",
  "description": "Permite ver reportes del sistema",
  "resource_id": 5,
  "action_id": 2,
  "scope_id": 1,
  "business_type_id": null
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Permiso creado exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Datos de entrada inválidos: name is required"
}
```

**Response 409 Conflict** (código duplicado):
```json
{
  "error": "Ya existe un permiso con este código"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "error": "Ya existe un permiso con el nombre 'Crear usuario'. Por favor, use un nombre diferente."
}
```

---

## 6. PUT /permissions/:id - Actualizar permiso

**URL**: `PUT /api/v1/permissions/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del permiso

**Request Body**:
```json
{
  "name": "Crear usuarios",
  "code": "users:create",
  "description": "Permite crear nuevos usuarios en el sistema",
  "resource_id": 1,
  "action_id": 1,
  "scope_id": 1,
  "business_type_id": 11
}
```

**Campos del Request** (todos requeridos):
- `name` (string, requerido): Nombre del permiso
- `code` (string, requerido): Código del permiso
- `description` (string, requerido): Descripción del permiso
- `resource_id` (int, requerido): ID del recurso
- `action_id` (int, requerido): ID de la acción
- `scope_id` (int, requerido): ID del scope
- `business_type_id` (int, opcional): ID del tipo de business (null para genérico)

**Ejemplo Request**:
```
PUT /api/v1/permissions/1
Content-Type: application/json

{
  "name": "Crear usuarios actualizado",
  "code": "users:create",
  "description": "Nueva descripción del permiso",
  "resource_id": 1,
  "action_id": 1,
  "scope_id": 1,
  "business_type_id": null
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Permiso actualizado exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Datos de entrada inválidos: name is required"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Permiso no encontrado"
}
```

**Response 409 Conflict** (código duplicado):
```json
{
  "error": "Ya existe otro permiso con este código"
}
```

---

## 7. DELETE /permissions/:id - Eliminar permiso

**URL**: `DELETE /api/v1/permissions/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del permiso

**Nota**: Esta es una eliminación física permanente. El permiso no se puede recuperar.

**Ejemplo Request**:
```
DELETE /api/v1/permissions/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Permiso eliminado exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "ID de permiso inválido"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Permiso no encontrado"
}
```

---

## Respuestas de Error Comunes

**400 Bad Request**:
```json
{
  "error": "Datos de entrada inválidos"
}
```

**401 Unauthorized**:
```json
{
  "error": "Token de acceso requerido"
}
```

**404 Not Found**:
```json
{
  "error": "Permiso no encontrado"
}
```

**409 Conflict**:
```json
{
  "error": "Ya existe un permiso con este código"
}
```

**500 Internal Server Error**:
```json
{
  "error": "Error interno del servidor"
}
```

---

## Notas Importantes

1. **Código de Permiso**: Si no se proporciona el campo `code` al crear un permiso, se genera automáticamente basado en el nombre y otros campos.

2. **Business Type ID**: 
   - Si se proporciona un `business_type_id`, el permiso es específico para ese tipo de business
   - Si es `null` o `0`, el permiso es genérico y aplica a todos los tipos de business

3. **Eliminación Permanente**: El endpoint DELETE realiza una eliminación física permanente. No hay soft delete.

4. **Filtros**: Los filtros en GET /permissions se pueden combinar:
   - `?business_type_id=1&name=crear` - Filtra por tipo de business Y nombre
   - `?scope_id=1` - Filtra solo por scope
   - `?business_type_id=1&scope_id=1&name=usuario` - Combina todos los filtros

5. **Recursos Comunes**: Algunos recursos comunes pueden ser:
   - `users`
   - `orders`
   - `inventory`
   - `reports`
   - `analytics`
   - `settings`

6. **Scopes Comunes**: Los scopes típicamente incluyen:
   - `1` - Sistema
   - `2` - Business
   - `3` - Módulo

