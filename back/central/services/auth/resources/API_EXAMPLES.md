# Ejemplos de API - Resources Module

Base URL: `/api/v1`

## 1. GET /resources - Obtener recursos filtrados y paginados

**URL**: `GET /api/v1/resources`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (default: 1, min: 1)
- `page_size` (int, opcional): Tamaño de página (default: 10, min: 1, max: 100)
- `name` (string, opcional): Filtrar por nombre (búsqueda parcial)
- `description` (string, opcional): Filtrar por descripción (búsqueda parcial)
- `business_type_id` (int, opcional): Filtrar por tipo de business (incluye genéricos) - solo super admin
- `sort_by` (string, opcional): Campo para ordenar (name, created_at, updated_at) (default: name)
- `sort_order` (string, opcional): Orden (asc, desc) (default: asc)

**Nota**: Usuarios normales solo ven recursos de su business_type (del token). Super admins pueden filtrar por cualquier business_type_id.

**Ejemplo Request**:
```
GET /api/v1/resources?page=1&page_size=10&name=user&business_type_id=1&sort_by=name&sort_order=asc
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recursos obtenidos exitosamente",
  "data": {
    "resources": [
      {
        "id": 1,
        "name": "users",
        "description": "Gestión de usuarios del sistema",
        "business_type_id": 11,
        "business_type_name": "Propiedad Horizontal",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": 2,
        "name": "orders",
        "description": "Gestión de pedidos",
        "business_type_id": 0,
        "business_type_name": "Genérico",
        "created_at": "2024-01-16T09:15:00Z",
        "updated_at": "2024-01-16T09:15:00Z"
      },
      {
        "id": 3,
        "name": "inventory",
        "description": "Gestión de inventario",
        "business_type_id": 11,
        "business_type_name": "Propiedad Horizontal",
        "created_at": "2024-01-17T14:20:00Z",
        "updated_at": "2024-01-17T14:20:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Parámetros de consulta inválidos",
  "error": "page_size must be between 1 and 100"
}
```

**Response 401 Unauthorized**:
```json
{
  "success": false,
  "message": "No autorizado",
  "error": "Token de acceso requerido"
}
```

---

## 2. GET /resources/:id - Obtener recurso por ID

**URL**: `GET /api/v1/resources/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del recurso

**Ejemplo Request**:
```
GET /api/v1/resources/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recurso obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "users",
    "description": "Gestión de usuarios del sistema",
    "business_type_id": 11,
    "business_type_name": "Propiedad Horizontal",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "ID de recurso inválido",
  "error": "El ID del recurso debe ser un número válido"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Recurso no encontrado",
  "error": "recurso con ID 999 no encontrado"
}
```

---

## 3. POST /resources - Crear recurso

**URL**: `POST /api/v1/resources`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "reports",
  "description": "Gestión de reportes del sistema",
  "business_type_id": 11
}
```

**Campos del Request**:
- `name` (string, requerido): Nombre único del recurso
- `description` (string, opcional): Descripción del recurso
- `business_type_id` (int, opcional): ID del tipo de business (null o 0 = genérico)

**Nota**: Solo super admins pueden crear recursos.

**Ejemplo Request**:
```
POST /api/v1/resources
Content-Type: application/json

{
  "name": "reports",
  "description": "Gestión de reportes del sistema",
  "business_type_id": 11
}
```

**Ejemplo Request** (recurso genérico):
```
POST /api/v1/resources
Content-Type: application/json

{
  "name": "analytics",
  "description": "Análisis y estadísticas del sistema",
  "business_type_id": null
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Recurso creado exitosamente",
  "data": {
    "id": 5,
    "name": "reports",
    "description": "Gestión de reportes del sistema",
    "business_type_id": 11,
    "business_type_name": "Propiedad Horizontal",
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
  "error": "el nombre del recurso es obligatorio"
}
```

**Response 400 Bad Request** (nombre muy largo):
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "el nombre del recurso no puede exceder 100 caracteres"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden crear recursos",
  "error": "permisos insuficientes"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "success": false,
  "message": "Recurso ya existe",
  "error": "ya existe un recurso con el nombre 'reports'"
}
```

---

## 4. PUT /resources/:id - Actualizar recurso

**URL**: `PUT /api/v1/resources/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del recurso

**Request Body**:
```json
{
  "name": "reports_updated",
  "description": "Gestión actualizada de reportes del sistema",
  "business_type_id": 11
}
```

**Campos del Request** (todos requeridos):
- `name` (string, requerido): Nombre único del recurso
- `description` (string, requerido): Descripción del recurso
- `business_type_id` (int, opcional): ID del tipo de business (null o 0 = genérico)

**Nota**: Solo super admins pueden actualizar recursos.

**Ejemplo Request**:
```
PUT /api/v1/resources/5
Content-Type: application/json

{
  "name": "reports_updated",
  "description": "Gestión actualizada de reportes del sistema",
  "business_type_id": 11
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recurso actualizado exitosamente",
  "data": {
    "id": 5,
    "name": "reports_updated",
    "description": "Gestión actualizada de reportes del sistema",
    "business_type_id": 11,
    "business_type_name": "Propiedad Horizontal",
    "created_at": "2024-01-22T10:00:00Z",
    "updated_at": "2024-01-22T15:45:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "el nombre del recurso es obligatorio"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden actualizar recursos",
  "error": "permisos insuficientes"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Recurso no encontrado",
  "error": "recurso con ID 999 no encontrado"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "success": false,
  "message": "Conflicto con recurso existente",
  "error": "ya existe otro recurso con el nombre 'reports_updated'"
}
```

---

## 5. DELETE /resources/:id - Eliminar recurso

**URL**: `DELETE /api/v1/resources/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del recurso

**Nota**: 
- Solo super admins pueden eliminar recursos.
- La eliminación es permanente y en cascada (elimina permisos asociados).

**Ejemplo Request**:
```
DELETE /api/v1/resources/5
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Recurso eliminado exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "ID de recurso inválido",
  "error": "El ID del recurso debe ser un número válido"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden eliminar recursos",
  "error": "permisos insuficientes"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Recurso no encontrado",
  "error": "recurso con ID 999 no encontrado"
}
```

---

## Respuestas de Error Comunes

**400 Bad Request**:
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "Detalles específicos del error"
}
```

**401 Unauthorized**:
```json
{
  "success": false,
  "message": "No autorizado",
  "error": "Token de acceso requerido"
}
```

**403 Forbidden**:
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden realizar esta acción",
  "error": "permisos insuficientes"
}
```

**404 Not Found**:
```json
{
  "success": false,
  "message": "Recurso no encontrado",
  "error": "recurso con ID {id} no encontrado"
}
```

**409 Conflict**:
```json
{
  "success": false,
  "message": "Recurso ya existe",
  "error": "ya existe un recurso con el nombre '{name}'"
}
```

**500 Internal Server Error**:
```json
{
  "success": false,
  "message": "Error interno del servidor",
  "error": "Detalles del error"
}
```

---

## Notas Importantes

1. **Filtros en GET /resources**: Los filtros se pueden combinar:
   - `?page=1&page_size=20&name=user&description=gestión` - Combina paginación, nombre y descripción
   - `?business_type_id=11&sort_by=created_at&sort_order=desc` - Filtra por tipo de business y ordena por fecha
   - `?name=order&sort_by=name&sort_order=asc` - Busca por nombre y ordena alfabéticamente

2. **Usuarios Normales vs Super Admin**:
   - **Usuarios normales**: Solo ven recursos de su `business_type_id` (obtenido del token)
   - **Super admins**: Pueden ver recursos de cualquier business_type usando `business_type_id` en query params
   - **Solo super admins** pueden crear, actualizar y eliminar recursos

3. **Recursos Genéricos vs Específicos**:
   - Si `business_type_id` es `null` o `0`, el recurso es genérico y aplica a todos los tipos de business
   - Si `business_type_id` tiene un valor, el recurso es específico para ese tipo de business

4. **Nombre Único**: 
   - El nombre del recurso debe ser único en todo el sistema
   - No puede haber dos recursos con el mismo nombre, incluso si son de diferentes business_types

5. **Validaciones**:
   - Nombre: obligatorio, máximo 100 caracteres
   - Descripción: opcional, máximo 500 caracteres
   - Business Type ID: opcional (null = genérico)

6. **Eliminación en Cascada**:
   - Al eliminar un recurso, se eliminan automáticamente todos los permisos asociados a ese recurso
   - Esta operación es permanente y no se puede deshacer

7. **Ordenamiento**:
   - Campos disponibles: `name`, `created_at`, `updated_at`
   - Orden por defecto: `name` ascendente (`asc`)

8. **Paginación**:
   - Por defecto: `page=1`, `page_size=10`
   - Máximo `page_size`: 100
   - La respuesta incluye información completa de paginación (`total`, `page`, `page_size`, `total_pages`)

