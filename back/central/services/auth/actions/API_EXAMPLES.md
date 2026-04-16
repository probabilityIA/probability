# Ejemplos de API - Actions Module

Base URL: `/api/v1`

## 1. GET /actions - Obtener actions filtrados y paginados

**URL**: `GET /api/v1/actions`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (default: 1, min: 1)
- `page_size` (int, opcional): Tamaño de página (default: 10, min: 1, max: 100)
- `name` (string, opcional): Filtrar por nombre (búsqueda parcial)

**Ejemplo Request**:
```
GET /api/v1/actions?page=1&page_size=10&name=create
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Actions obtenidos exitosamente",
  "data": {
    "actions": [
      {
        "id": 1,
        "name": "create",
        "description": "Permite crear nuevos registros",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": 2,
        "name": "read",
        "description": "Permite leer registros",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": 3,
        "name": "update",
        "description": "Permite actualizar registros existentes",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
      },
      {
        "id": 4,
        "name": "delete",
        "description": "Permite eliminar registros",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z"
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

## 2. GET /actions/:id - Obtener action por ID

**URL**: `GET /api/v1/actions/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del action

**Ejemplo Request**:
```
GET /api/v1/actions/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Action obtenido exitosamente",
  "data": {
    "id": 1,
    "name": "create",
    "description": "Permite crear nuevos registros",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "ID de action inválido",
  "error": "El ID del action debe ser un número válido"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Action no encontrado",
  "error": "action con ID 999 no encontrado"
}
```

---

## 3. POST /actions - Crear action

**URL**: `POST /api/v1/actions`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body**:
```json
{
  "name": "approve",
  "description": "Permite aprobar solicitudes"
}
```

**Campos del Request**:
- `name` (string, requerido): Nombre único del action (máximo 20 caracteres)
- `description` (string, opcional): Descripción del action (máximo 255 caracteres)

**Nota**: Solo super admins pueden crear actions.

**Ejemplo Request**:
```
POST /api/v1/actions
Content-Type: application/json

{
  "name": "approve",
  "description": "Permite aprobar solicitudes"
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "message": "Action creado exitosamente",
  "data": {
    "id": 5,
    "name": "approve",
    "description": "Permite aprobar solicitudes",
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
  "error": "el nombre del action es obligatorio"
}
```

**Response 400 Bad Request** (nombre muy largo):
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "el nombre del action no puede exceder 20 caracteres"
}
```

**Response 400 Bad Request** (descripción muy larga):
```json
{
  "success": false,
  "message": "Datos de entrada inválidos",
  "error": "la descripción del action no puede exceder 255 caracteres"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden crear actions",
  "error": "permisos insuficientes"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "success": false,
  "message": "Action ya existe",
  "error": "ya existe un action con el nombre 'approve'"
}
```

---

## 4. PUT /actions/:id - Actualizar action

**URL**: `PUT /api/v1/actions/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del action

**Request Body**:
```json
{
  "name": "approve_updated",
  "description": "Permite aprobar y rechazar solicitudes"
}
```

**Campos del Request** (todos requeridos):
- `name` (string, requerido): Nombre único del action (máximo 20 caracteres)
- `description` (string, requerido): Descripción del action (máximo 255 caracteres)

**Nota**: Solo super admins pueden actualizar actions.

**Ejemplo Request**:
```
PUT /api/v1/actions/5
Content-Type: application/json

{
  "name": "approve_updated",
  "description": "Permite aprobar y rechazar solicitudes"
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Action actualizado exitosamente",
  "data": {
    "id": 5,
    "name": "approve_updated",
    "description": "Permite aprobar y rechazar solicitudes",
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
  "error": "el nombre del action es obligatorio"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden actualizar actions",
  "error": "permisos insuficientes"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Action no encontrado",
  "error": "action con ID 999 no encontrado"
}
```

**Response 409 Conflict** (nombre duplicado):
```json
{
  "success": false,
  "message": "Conflicto con action existente",
  "error": "ya existe otro action con el nombre 'approve_updated'"
}
```

---

## 5. DELETE /actions/:id - Eliminar action

**URL**: `DELETE /api/v1/actions/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del action

**Nota**: 
- Solo super admins pueden eliminar actions.
- No se puede eliminar un action si tiene permisos asociados.

**Ejemplo Request**:
```
DELETE /api/v1/actions/5
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Action eliminado exitosamente"
}
```

**Response 400 Bad Request**:
```json
{
  "success": false,
  "message": "ID de action inválido",
  "error": "El ID del action debe ser un número válido"
}
```

**Response 403 Forbidden** (usuario no super admin):
```json
{
  "success": false,
  "message": "Solo los super usuarios pueden eliminar actions",
  "error": "permisos insuficientes"
}
```

**Response 404 Not Found**:
```json
{
  "success": false,
  "message": "Action no encontrado",
  "error": "action con ID 999 no encontrado"
}
```

**Response 409 Conflict** (tiene permisos asociados):
```json
{
  "success": false,
  "message": "no se puede eliminar el action porque tiene permisos asociados",
  "error": "no se puede eliminar el action porque tiene permisos asociados"
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
  "message": "Action no encontrado",
  "error": "action con ID {id} no encontrado"
}
```

**409 Conflict**:
```json
{
  "success": false,
  "message": "Action ya existe",
  "error": "ya existe un action con el nombre '{name}'"
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

1. **Filtros en GET /actions**: 
   - `?page=1&page_size=20&name=create` - Combina paginación y búsqueda por nombre
   - `?name=update` - Busca actions que contengan "update" en el nombre

2. **Solo Super Admins**: 
   - Solo los super usuarios pueden crear, actualizar y eliminar actions
   - Todos los usuarios autenticados pueden leer (GET) actions

3. **Nombre Único**: 
   - El nombre del action debe ser único en todo el sistema
   - No puede haber dos actions con el mismo nombre

4. **Validaciones**:
   - Nombre: obligatorio, máximo 20 caracteres
   - Descripción: opcional, máximo 255 caracteres

5. **Eliminación con Restricciones**:
   - No se puede eliminar un action si tiene permisos asociados
   - Primero se deben eliminar o actualizar todos los permisos que usan ese action

6. **Actions Comunes**:
   - `create` - Crear nuevos registros
   - `read` - Leer registros
   - `update` - Actualizar registros existentes
   - `delete` - Eliminar registros
   - `approve` - Aprobar solicitudes
   - `reject` - Rechazar solicitudes
   - `view` - Ver información
   - `export` - Exportar datos
   - `import` - Importar datos

7. **Paginación**:
   - Por defecto: `page=1`, `page_size=10`
   - Máximo `page_size`: 100
   - La respuesta incluye información completa de paginación (`total`, `page`, `page_size`, `total_pages`)

8. **Sin Business Type**:
   - A diferencia de los recursos, los actions son genéricos y no están asociados a un tipo de business específico
   - Todos los actions están disponibles para todos los tipos de business

