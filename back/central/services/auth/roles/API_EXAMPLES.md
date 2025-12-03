# Ejemplos de API - Roles Module

Base URL: `/api/v1`

## 1. GET /roles - Obtener todos los roles

**URL**: `GET /api/v1/roles`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `business_type_id` (int, opcional): Filtrar por tipo de business
- `scope_id` (int, opcional): Filtrar por ID de scope
- `is_system` (boolean, opcional): Filtrar por rol de sistema (true/false)
- `name` (string, opcional): Buscar en el nombre del rol (búsqueda parcial)
- `level` (int, opcional): Filtrar por nivel del rol (1-10)

**Nota**: Usuarios normales solo ven roles de su business_type (del token). Super admins pueden filtrar por cualquier business_type_id.

**Ejemplo Request**:
```
GET /api/v1/roles?business_type_id=1&scope_id=1&is_system=false&name=admin&level=2
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Administrador",
      "code": "admin",
      "description": "Rol de administrador del sistema",
      "level": 2,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 2,
      "name": "Gerente",
      "code": "manager",
      "description": "Rol de gerente",
      "level": 3,
      "is_system": false,
      "scope_id": 2,
      "scope_name": "Business",
      "scope_code": "business",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    }
  ],
  "count": 2
}
```

**Response 401 Unauthorized**:
```json
{
  "error": "Token de acceso requerido"
}
```

---

## 2. GET /roles/:id - Obtener rol por ID

**URL**: `GET /api/v1/roles/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del rol

**Ejemplo Request**:
```
GET /api/v1/roles/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Administrador",
    "code": "admin",
    "description": "Rol de administrador del sistema",
    "level": 2,
    "is_system": true,
    "scope_id": 1,
    "scope_name": "Sistema",
    "scope_code": "system",
    "business_type_id": 1,
    "business_type_name": "Propiedad Horizontal"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "error": "ID inválido: id is required"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Rol no encontrado"
}
```

---

## 3. GET /roles/scope/:scope_id - Obtener roles por scope

**URL**: `GET /api/v1/roles/scope/:scope_id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `scope_id` (int, requerido): ID del scope

**Ejemplo Request**:
```
GET /api/v1/roles/scope/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Administrador",
      "code": "admin",
      "description": "Rol de administrador del sistema",
      "level": 2,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 3,
      "name": "Super Admin",
      "code": "super_admin",
      "description": "Rol de super administrador",
      "level": 1,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    }
  ],
  "count": 2
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Scope ID de rol inválido"
}
```

---

## 4. GET /roles/level/:level - Obtener roles por nivel

**URL**: `GET /api/v1/roles/level/:level`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `level` (int, requerido): Nivel del rol (1-10)

**Ejemplo Request**:
```
GET /api/v1/roles/level/2
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Administrador",
      "code": "admin",
      "description": "Rol de administrador del sistema",
      "level": 2,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 4,
      "name": "Coordinador",
      "code": "coordinator",
      "description": "Rol de coordinador",
      "level": 2,
      "is_system": false,
      "scope_id": 2,
      "scope_name": "Business",
      "scope_code": "business",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    }
  ],
  "count": 2
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Nivel inválido: level must be between 1 and 10"
}
```

---

## 5. GET /roles/system - Obtener roles del sistema

**URL**: `GET /api/v1/roles/system`

**Headers**:
```
Authorization: Bearer {token}
```

**Ejemplo Request**:
```
GET /api/v1/roles/system
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Administrador",
      "code": "admin",
      "description": "Rol de administrador del sistema",
      "level": 2,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 1,
      "business_type_name": "Propiedad Horizontal"
    },
    {
      "id": 3,
      "name": "Super Admin",
      "code": "super_admin",
      "description": "Rol de super administrador",
      "level": 1,
      "is_system": true,
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system",
      "business_type_id": 0,
      "business_type_name": "Genérico"
    }
  ],
  "count": 2
}
```

---

## 6. POST /roles - Crear rol

**URL**: `POST /api/v1/roles`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "Administrador",
  "description": "Rol de administrador del sistema",
  "level": 2,
  "is_system": false,
  "scope_id": 1,
  "business_type_id": 1
}
```

**Campos del Request**:
- `name` (string, requerido): Nombre del rol
- `description` (string, requerido): Descripción del rol
- `level` (int, requerido): Nivel del rol (1-10, min: 1, max: 10)
- `is_system` (boolean, requerido): ¿Es rol del sistema?
- `scope_id` (int, requerido): ID del scope
- `business_type_id` (int, requerido): ID del tipo de business

**Ejemplo Request**:
```
POST /api/v1/roles
Content-Type: application/json

{
  "name": "Editor",
  "description": "Rol de editor de contenido",
  "level": 4,
  "is_system": false,
  "scope_id": 2,
  "business_type_id": 1
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Rol creado exitosamente",
  "data": {
    "id": 5,
    "name": "Editor",
    "description": "Rol de editor de contenido",
    "level": 4,
    "is_system": false,
    "scope_id": 2,
    "business_type_id": 1,
    "created_at": "2024-01-22T10:00:00Z",
    "updated_at": "2024-01-22T10:00:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "name is required"
}
```

---

## 7. PUT /roles/:id - Actualizar rol

**URL**: `PUT /api/v1/roles/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del rol

**Request Body** (todos los campos son opcionales):
```json
{
  "name": "Administrador Actualizado",
  "description": "Nueva descripción del rol",
  "level": 3,
  "is_system": false,
  "scope_id": 1,
  "business_type_id": 1
}
```

**Campos del Request** (todos opcionales):
- `name` (string, opcional): Nombre del rol
- `description` (string, opcional): Descripción del rol
- `level` (int, opcional): Nivel del rol (1-10)
- `is_system` (boolean, opcional): ¿Es rol del sistema?
- `scope_id` (int, opcional): ID del scope
- `business_type_id` (int, opcional): ID del tipo de business

**Ejemplo Request**:
```
PUT /api/v1/roles/1
Content-Type: application/json

{
  "name": "Administrador Actualizado",
  "level": 3
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Rol actualizado exitosamente",
  "data": {
    "id": 1,
    "name": "Administrador Actualizado",
    "description": "Rol de administrador del sistema",
    "level": 3,
    "is_system": true,
    "scope_id": 1,
    "business_type_id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-22T11:30:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "level must be between 1 and 10"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Rol no encontrado",
  "error": "rol no encontrado"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "success": false,
  "message": "Ya existe un rol con este nombre. Por favor, use un nombre diferente.",
  "error": "duplicate key value violates unique constraint \"uni_role_name\""
}
```

---

## 8. POST /roles/:id/permissions - Asignar permisos a un rol

**URL**: `POST /api/v1/roles/:id/permissions`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del rol

**Request Body**:
```json
{
  "permission_ids": [1, 2, 3, 5]
}
```

**Campos del Request**:
- `permission_ids` (array de int, requerido): Array de IDs de permisos a asignar (mínimo 1 permiso)

**Nota**: Solo se pueden asignar permisos que pertenezcan al mismo business_type que el rol.

**Ejemplo Request**:
```
POST /api/v1/roles/1/permissions
Content-Type: application/json

{
  "permission_ids": [1, 2, 3, 5, 8]
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Permisos asignados exitosamente al rol",
  "role_id": 1,
  "permission_ids": [1, 2, 3, 5, 8]
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Datos de entrada inválidos"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Rol no encontrado"
}
```

**Response 500 Internal Server Error** (permisos de diferente business_type):
```json
{
  "error": "Algunos permisos no pertenecen al mismo business_type que el rol"
}
```

---

## 9. GET /roles/:id/permissions - Obtener permisos de un rol

**URL**: `GET /api/v1/roles/:id/permissions`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del rol

**Ejemplo Request**:
```
GET /api/v1/roles/1/permissions
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Permisos del rol obtenidos exitosamente",
  "role_id": 1,
  "role_name": "Administrador",
  "permissions": [
    {
      "id": 1,
      "resource": "users",
      "action": "create",
      "description": "Crear usuarios",
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system"
    },
    {
      "id": 2,
      "resource": "users",
      "action": "update",
      "description": "Editar usuarios",
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system"
    },
    {
      "id": 3,
      "resource": "users",
      "action": "delete",
      "description": "Eliminar usuarios",
      "scope_id": 1,
      "scope_name": "Sistema",
      "scope_code": "system"
    }
  ],
  "count": 3
}
```

**Response 400 Bad Request**:
```json
{
  "error": "ID de rol inválido"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Rol no encontrado"
}
```

---

## 10. DELETE /roles/:id/permissions/:permission_id - Eliminar permiso de un rol

**URL**: `DELETE /api/v1/roles/:id/permissions/:permission_id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del rol
- `permission_id` (int, requerido): ID del permiso a eliminar

**Ejemplo Request**:
```
DELETE /api/v1/roles/1/permissions/3
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Permiso eliminado exitosamente del rol"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "ID de rol inválido"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Rol no encontrado"
}
```

**Response 404 Not Found** (permiso no asignado):
```json
{
  "error": "permiso no encontrado"
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
  "error": "Rol no encontrado"
}
```

**409 Conflict**:
```json
{
  "success": false,
  "message": "Ya existe un rol con este nombre. Por favor, use un nombre diferente.",
  "error": "duplicate key value violates unique constraint \"uni_role_name\""
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

1. **Filtros en GET /roles**: Los filtros se pueden combinar:
   - `?business_type_id=1&scope_id=1` - Filtra por tipo de business Y scope
   - `?is_system=true&level=2` - Filtra por sistema Y nivel
   - `?name=admin&business_type_id=1&scope_id=1` - Combina todos los filtros

2. **Usuarios Normales vs Super Admin**:
   - **Usuarios normales**: Solo ven roles de su `business_type_id` (obtenido del token)
   - **Super admins**: Pueden filtrar por cualquier `business_type_id` usando query params

3. **Niveles de Rol**: Los niveles van del 1 al 10, donde:
   - Nivel 1 = Mayor jerarquía (Super Admin)
   - Nivel 10 = Menor jerarquía

4. **Roles del Sistema**: Los roles con `is_system: true` son roles especiales del sistema que generalmente no se pueden modificar.

5. **Asignación de Permisos**: 
   - Solo se pueden asignar permisos que pertenezcan al mismo `business_type_id` que el rol
   - Si el rol es genérico (business_type_id = 0), solo se pueden asignar permisos genéricos

6. **Código de Rol**: El código se genera automáticamente basado en el nombre del rol (normalizado).

7. **Actualización Parcial**: El endpoint PUT permite actualizar solo los campos que se envíen. Los campos no enviados no se modifican.

