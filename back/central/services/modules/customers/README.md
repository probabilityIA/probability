# Módulo `customers`

Gestión de clientes del negocio. Permite crear, consultar, actualizar y eliminar clientes, con estadísticas básicas de sus órdenes.

## Estructura

```
customers/
├── bundle.go                                       # Entry point del módulo
└── internal/
    ├── domain/
    │   ├── entities/client.go                      # Entidad Customer (sin tags)
    │   ├── dtos/client.go                          # ListCustomersParams, CreateCustomerDTO, UpdateCustomerDTO
    │   ├── ports/ports.go                          # IRepository interface
    │   └── errors/errors.go                        # Errores del dominio
    ├── app/
    │   ├── constructor.go                          # IUseCase + New()
    │   ├── create_customer.go
    │   ├── get_customer.go                         # Incluye stats de órdenes
    │   ├── list_customers.go
    │   ├── update_customer.go
    │   └── delete_customer.go                      # Bloquea si tiene órdenes
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go                      # IHandlers + New()
        │   ├── routes.go                           # RegisterRoutes()
        │   ├── create_customer.go
        │   ├── get_customer.go
        │   ├── list_customers.go
        │   ├── update_customer.go
        │   ├── delete_customer.go
        │   ├── request/client.go                   # CreateCustomerRequest, UpdateCustomerRequest
        │   └── response/client.go                  # CustomerResponse, CustomerDetailResponse
        └── secondary/repository/
            └── repository.go                       # CRUD + GetOrderStats
```

## Endpoints

| Método   | Ruta              | Descripción                                            |
|----------|-------------------|--------------------------------------------------------|
| `GET`    | `/api/v1/customers`     | Listar clientes (paginado, filtro por `?search=`) |
| `GET`    | `/api/v1/customers/:id` | Detalle del cliente + estadísticas de órdenes    |
| `POST`   | `/api/v1/customers`     | Crear cliente                                    |
| `PUT`    | `/api/v1/customers/:id` | Actualizar cliente                               |
| `DELETE` | `/api/v1/customers/:id` | Eliminar cliente (soft delete)                   |

### Query params de `GET /customers`

| Parámetro   | Tipo   | Default | Descripción                               |
|-------------|--------|---------|-------------------------------------------|
| `page`      | int    | 1       | Número de página                          |
| `page_size` | int    | 20      | Registros por página (máx. 100)           |
| `search`    | string | —       | Búsqueda por nombre, email o teléfono     |

### Respuesta `GET /customers` (200)

```json
{
  "data": [
    {
      "id": 1,
      "business_id": 5,
      "name": "Juan Pérez",
      "email": "juan@example.com",
      "phone": "3001234567",
      "dni": "12345678",
      "created_at": "2026-02-26T10:00:00Z",
      "updated_at": "2026-02-26T10:00:00Z"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 20,
  "total_pages": 3
}
```

### Respuesta `GET /customers/:id` (200)

Incluye estadísticas calculadas desde la tabla `orders`:

```json
{
  "id": 1,
  "business_id": 5,
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "phone": "3001234567",
  "dni": "12345678",
  "created_at": "2026-02-26T10:00:00Z",
  "updated_at": "2026-02-26T10:00:00Z",
  "order_count": 7,
  "total_spent": 350000,
  "last_order_at": "2026-02-20T15:30:00Z"
}
```

### Request `POST /customers`

```json
{
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "phone": "3001234567",
  "dni": "12345678"
}
```

> `name` es requerido (mín. 2 caracteres). `email`, `phone` y `dni` son opcionales.

## Errores

| Error                  | HTTP | Descripción                                          |
|------------------------|------|------------------------------------------------------|
| `ErrClientNotFound`    | 404  | Cliente no encontrado en el negocio                 |
| `ErrDuplicateEmail`    | 409  | Ya existe un cliente con ese email en el negocio    |
| `ErrDuplicateDni`      | 409  | Ya existe un cliente con ese DNI en el negocio      |
| `ErrClientHasOrders`   | 409  | El cliente tiene órdenes y no puede eliminarse      |

## Modelo en base de datos

Usa el modelo `Client` de `migration/shared/models/models.go`. **No hay migración pendiente** — la tabla `clients` ya existe.

```go
type Client struct {
    gorm.Model
    BusinessID uint    // FK → businesses.id
    Name       string
    Email      string  // unique por negocio
    Phone      string
    Dni        *string // unique por negocio (nullable)
}
```

## Notas de arquitectura

- **Aislamiento de repositorios**: `GetOrderStats` consulta la tabla `orders` directamente (`WHERE customer_id = ?`) sin importar el repositorio del módulo `orders`.
- **Multi-tenant**: todos los queries filtran por `business_id` extraído del JWT.
- **Soft delete**: usa `gorm.Model` → los registros eliminados quedan en la tabla con `deleted_at`.
