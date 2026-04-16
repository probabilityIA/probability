# Ejemplos de API - Users Module

Base URL: `/api/v1`

## 1. GET /users - Obtener usuarios filtrados y paginados

**URL**: `GET /api/v1/users`

**Headers**:
```
Authorization: Bearer {token}
```

**Query Parameters**:
- `page` (int, opcional): Número de página (default: 1, min: 1)
- `page_size` (int, opcional): Tamaño de página (default: 10, min: 1, max: 100)
- `name` (string, opcional): Filtrar por nombre (búsqueda parcial)
- `email` (string, opcional): Filtrar por email (búsqueda parcial)
- `phone` (string, opcional): Filtrar por teléfono (búsqueda parcial, exactamente 10 dígitos)
- `user_ids` (string, opcional): Filtrar por IDs de usuarios separados por comas (ej: "1,2,3")
- `is_active` (boolean, opcional): Filtrar por estado activo
- `role_id` (int, opcional): Filtrar por ID de rol
- `business_id` (int, opcional): Filtrar por ID de business (solo super admin)
- `created_at` (string, opcional): Filtrar por fecha de creación (formato: "YYYY-MM-DD" o "YYYY-MM-DD,YYYY-MM-DD" para rango)
- `sort_by` (string, opcional): Campo para ordenar (id, name, email, phone, is_active, created_at, updated_at) (default: created_at)
- `sort_order` (string, opcional): Orden de clasificación (asc, desc) (default: desc)

**Nota**: Usuarios normales solo ven usuarios de su business (del token). Super admins pueden filtrar por cualquier business_id.

**Ejemplo Request**:
```
GET /api/v1/users?page=1&page_size=10&name=Juan&is_active=true&role_id=2&sort_by=name&sort_order=asc
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Juan Pérez",
      "email": "juan.perez@example.com",
      "phone": "3001234567",
      "avatar_url": "https://s3.amazonaws.com/bucket/avatar1.jpg",
      "is_active": true,
      "is_super_user": false,
      "last_login_at": "2024-01-20T14:30:00Z",
      "business_role_assignments": [
        {
          "business_id": 16,
          "business_name": "Restaurante El Buen Sabor",
          "role_id": 4,
          "role_name": "Gerente"
        },
        {
          "business_id": 21,
          "business_name": "Cafetería Central",
          "role_id": 5,
          "role_name": "Administrador"
        }
      ],
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-20T14:20:00Z"
    },
    {
      "id": 2,
      "name": "María García",
      "email": "maria.garcia@example.com",
      "phone": "3007654321",
      "avatar_url": "",
      "is_active": true,
      "is_super_user": false,
      "last_login_at": null,
      "business_role_assignments": [
        {
          "business_id": 16,
          "business_name": "Restaurante El Buen Sabor",
          "role_id": 6,
          "role_name": "Mesero"
        }
      ],
      "created_at": "2024-01-16T09:15:00Z",
      "updated_at": "2024-01-16T09:15:00Z"
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

**Response 400 Bad Request**:
```json
{
  "error": "Parámetros de filtro inválidos"
}
```

---

## 2. GET /users/:id - Obtener usuario por ID

**URL**: `GET /api/v1/users/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del usuario

**Nota**: Usuarios normales solo pueden ver usuarios de su business. Super admins pueden ver cualquier usuario.

**Ejemplo Request**:
```
GET /api/v1/users/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Juan Pérez",
    "email": "juan.perez@example.com",
    "phone": "3001234567",
    "avatar_url": "https://s3.amazonaws.com/bucket/avatar1.jpg",
    "is_active": true,
    "is_super_user": false,
    "last_login_at": "2024-01-20T14:30:00Z",
    "business_role_assignments": [
      {
        "business_id": 16,
        "business_name": "Restaurante El Buen Sabor",
        "role_id": 4,
        "role_name": "Gerente"
      },
      {
        "business_id": 21,
        "business_name": "Cafetería Central",
        "role_id": 5,
        "role_name": "Administrador"
      }
    ],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-20T14:20:00Z"
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
  "error": "Usuario no encontrado"
}
```

---

## 3. POST /users - Crear usuario

**URL**: `POST /api/v1/users`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Form Data Parameters**:
- `name` (string, requerido): Nombre del usuario (2-100 caracteres)
- `email` (string, requerido): Email válido
- `phone` (string, opcional): Teléfono (exactamente 10 dígitos)
- `is_active` (boolean, opcional): ¿Activo? (default: true)
- `avatarFile` (file, opcional): Imagen de avatar (se sube a S3)
- `business_ids` (string, opcional): IDs de negocios separados por comas (ej: "16,21")

**Nota**: 
- La contraseña se genera automáticamente y se devuelve en la respuesta (solo se muestra una vez).
- Si se envía `business_ids`, el usuario se asocia a esos negocios.
- El formato acepta tanto JSON como multipart/form-data.

**Ejemplo Request** (multipart/form-data):
```
POST /api/v1/users
Content-Type: multipart/form-data

name=Juan Pérez
email=juan.perez@example.com
phone=3001234567
is_active=true
business_ids=16,21
avatarFile=[archivo binario]
```

**Ejemplo Request** (JSON):
```
POST /api/v1/users
Content-Type: application/json

{
  "name": "Juan Pérez",
  "email": "juan.perez@example.com",
  "phone": "3001234567",
  "is_active": true,
  "business_ids": [16, 21]
}
```

**Response 201 Created**:
```json
{
  "success": true,
  "email": "juan.perez@example.com",
  "password": "TempPass123!",
  "message": "Usuario creado exitosamente. La contraseña generada es: TempPass123!"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "El nombre debe tener entre 2 y 100 caracteres"
}
```

**Response 400 Bad Request** (teléfono inválido):
```json
{
  "error": "El teléfono debe tener exactamente 10 dígitos"
}
```

**Response 409 Conflict** (email duplicado):
```json
{
  "error": "El email ya está registrado en el sistema"
}
```

---

## 4. PUT /users/:id - Actualizar usuario

**URL**: `PUT /api/v1/users/:id`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: multipart/form-data
```

**Path Parameters**:
- `id` (int, requerido): ID del usuario

**Form Data Parameters** (todos opcionales):
- `name` (string): Nombre del usuario (2-100 caracteres)
- `email` (string): Email válido
- `phone` (string): Teléfono (exactamente 10 dígitos)
- `is_active` (boolean): ¿Activo?
- `remove_avatar` (boolean): Eliminar avatar actual (true/false)
- `avatarFile` (file): Nueva imagen de avatar (se sube a S3)
- `business_ids` (string): IDs de negocios separados por comas (ej: "16,21") - sustituye relaciones existentes

**Nota**: 
- Usuarios normales solo pueden actualizar usuarios de su business. Super admins pueden actualizar cualquier usuario.
- Si se envía `business_ids`, se sustituyen todas las relaciones de businesses del usuario.
- El formato acepta tanto JSON como multipart/form-data.

**Ejemplo Request** (multipart/form-data):
```
PUT /api/v1/users/1
Content-Type: multipart/form-data

name=Juan Pérez Actualizado
phone=3009876543
is_active=false
business_ids=16,21,25
avatarFile=[archivo binario]
```

**Ejemplo Request** (JSON):
```
PUT /api/v1/users/1
Content-Type: application/json

{
  "name": "Juan Pérez Actualizado",
  "phone": "3009876543",
  "is_active": false,
  "business_ids": [16, 21, 25]
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Juan Pérez Actualizado",
    "email": "juan.perez@example.com",
    "phone": "3009876543",
    "avatar_url": "https://s3.amazonaws.com/bucket/avatar1_updated.jpg",
    "is_active": false,
    "is_super_user": false,
    "last_login_at": "2024-01-20T14:30:00Z",
    "business_role_assignments": [
      {
        "business_id": 16,
        "business_name": "Restaurante El Buen Sabor",
        "role_id": 4,
        "role_name": "Gerente"
      },
      {
        "business_id": 21,
        "business_name": "Cafetería Central",
        "role_id": 5,
        "role_name": "Administrador"
      },
      {
        "business_id": 25,
        "business_name": "Bar El Refugio",
        "role_id": 6,
        "role_name": "Mesero"
      }
    ],
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-22T15:45:00Z"
  }
}
```

**Response 400 Bad Request**:
```json
{
  "error": "El email no tiene un formato válido"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Usuario no encontrado"
}
```

**Response 409 Conflict** (email duplicado):
```json
{
  "error": "El email ya está registrado en el sistema"
}
```

---

## 5. DELETE /users/:id - Eliminar usuario

**URL**: `DELETE /api/v1/users/:id`

**Headers**:
```
Authorization: Bearer {token}
```

**Path Parameters**:
- `id` (int, requerido): ID del usuario

**Ejemplo Request**:
```
DELETE /api/v1/users/1
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Usuario eliminado exitosamente"
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
  "error": "Usuario no encontrado"
}
```

---

## 6. POST /users/:id/assign-role - Asignar roles a usuario en businesses

**URL**: `POST /api/v1/users/:id/assign-role`

**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
```

**Path Parameters**:
- `id` (int, requerido): ID del usuario

**Request Body**:
```json
{
  "assignments": [
    {
      "business_id": 16,
      "role_id": 4
    },
    {
      "business_id": 21,
      "role_id": 5
    }
  ]
}
```

**Campos del Request**:
- `assignments` (array, requerido): Array de asignaciones (mínimo 1)
  - `business_id` (int, requerido): ID del business
  - `role_id` (int, requerido): ID del rol

**Nota**: 
- Solo super admins pueden asignar roles a otros usuarios. Usuarios normales solo pueden asignarse roles a sí mismos.
- El usuario debe estar previamente asociado a cada business.
- Solo se permite un rol por business.
- Cada rol debe ser del mismo tipo de business que su business asociado.

**Ejemplo Request**:
```
POST /api/v1/users/1/assign-role
Content-Type: application/json

{
  "assignments": [
    {
      "business_id": 16,
      "role_id": 4
    },
    {
      "business_id": 21,
      "role_id": 5
    },
    {
      "business_id": 25,
      "role_id": 6
    }
  ]
}
```

**Response 200 OK**:
```json
{
  "success": true,
  "message": "Roles asignados exitosamente al usuario en los businesses"
}
```

**Response 400 Bad Request**:
```json
{
  "error": "Datos de entrada inválidos: assignments is required"
}
```

**Response 400 Bad Request** (sin asignaciones):
```json
{
  "error": "Debe proporcionar al menos una asignación"
}
```

**Response 403 Forbidden** (usuario normal intentando asignar a otro):
```json
{
  "error": "No tienes permisos para asignar roles a otros usuarios"
}
```

**Response 403 Forbidden** (usuario no asociado al business):
```json
{
  "error": "El usuario no está asociado al business con ID 16"
}
```

**Response 403 Forbidden** (rol no corresponde al tipo de business):
```json
{
  "error": "El rol con ID 4 no corresponde al tipo de business del business con ID 16"
}
```

**Response 404 Not Found**:
```json
{
  "error": "Usuario no encontrado"
}
```

**Response 404 Not Found** (business no encontrado):
```json
{
  "error": "Algunos businesses no fueron encontrados"
}
```

**Response 404 Not Found** (rol no encontrado):
```json
{
  "error": "Algunos roles no fueron encontrados"
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

**401 Unauthorized** (business_id no disponible):
```json
{
  "error": "Token inválido: business_id no disponible"
}
```

**403 Forbidden**:
```json
{
  "error": "No tienes permisos para realizar esta acción"
}
```

**404 Not Found**:
```json
{
  "error": "Usuario no encontrado"
}
```

**409 Conflict**:
```json
{
  "error": "El email ya está registrado en el sistema"
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

1. **Filtros en GET /users**: Los filtros se pueden combinar:
   - `?page=1&page_size=20&name=Juan&is_active=true` - Combina paginación, nombre y estado
   - `?role_id=2&business_id=16&sort_by=name&sort_order=asc` - Filtra por rol y business, ordena por nombre
   - `?created_at=2024-01-01,2024-01-31` - Filtra por rango de fechas

2. **Usuarios Normales vs Super Admin**:
   - **Usuarios normales**: Solo ven y gestionan usuarios de su `business_id` (obtenido del token)
   - **Super admins**: Pueden ver y gestionar usuarios de cualquier business usando `business_id` en query params

3. **Creación de Usuario**:
   - La contraseña se genera automáticamente y se devuelve en la respuesta
   - **IMPORTANTE**: La contraseña solo se muestra una vez en la creación
   - Si se envía `business_ids`, el usuario se asocia a esos negocios automáticamente

4. **Actualización de Usuario**:
   - Es una actualización parcial: solo se actualizan los campos enviados
   - Si se envía `business_ids`, se **sustituyen** todas las relaciones de businesses existentes
   - Para eliminar el avatar, enviar `remove_avatar=true`
   - Para actualizar el avatar, enviar un nuevo `avatarFile`

5. **Asignación de Roles**:
   - El usuario debe estar previamente asociado a cada business
   - Solo se permite un rol por business
   - Cada rol debe ser del mismo tipo de business que su business asociado
   - Solo super admins pueden asignar roles a otros usuarios

6. **Formato de Datos**:
   - Los endpoints POST y PUT aceptan tanto `multipart/form-data` como `application/json`
   - Para subir archivos (avatar), usar `multipart/form-data`
   - Para `business_ids` en multipart, usar string separado por comas: `"16,21"`
   - Para `business_ids` en JSON, usar array: `[16, 21]`

7. **Paginación**:
   - Por defecto: `page=1`, `page_size=10`
   - Máximo `page_size`: 100
   - La respuesta incluye información completa de paginación

8. **Ordenamiento**:
   - Campos disponibles: `id`, `name`, `email`, `phone`, `is_active`, `created_at`, `updated_at`
   - Orden por defecto: `created_at` descendente (`desc`)

