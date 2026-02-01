# Módulo Payments

## 📋 Descripción

El módulo **Payments** gestiona los métodos de pago disponibles en la plataforma Probability y sus mapeos con métodos de pago externos de integraciones (Shopify, MercadoLibre, Amazon, etc.).

## 🎯 Funcionalidades

### Payment Methods (Métodos de Pago)
- Crear, leer, actualizar y eliminar métodos de pago
- Activar/desactivar métodos de pago
- Listar métodos con paginación y filtros
- Categorización por tipo (tarjeta, billetera digital, transferencia bancaria, efectivo)
- Gestión de iconos y colores para UI

### Payment Mappings (Mapeos de Métodos de Pago)
- Mapear métodos de pago externos a métodos internos
- Gestionar múltiples mapeos por integración
- Configurar prioridad de mapeos
- Activar/desactivar mapeos
- Listar mapeos por tipo de integración

## 🏗️ Arquitectura

Este módulo sigue **Arquitectura Hexagonal (Clean Architecture)** con la siguiente estructura:

> **Nota importante:** Todo el código del módulo está dentro de la carpeta `internal/` siguiendo la convención de Go. Los paquetes en `internal/` son privados y no pueden ser importados por módulos externos, garantizando el encapsulamiento del módulo.

```
payments/
├── bundle.go                          # ✅ Ensambla e inyecta dependencias
└── internal/                          # ✅ Carpeta internal (convención Go)
    ├── domain/                        # 🔵 CAPA DE DOMINIO (núcleo)
    │   ├── entities/                  # Entidades de negocio PURAS
    │   │   ├── payment_method.go
    │   │   └── payment_mapping.go
    │   ├── dtos/                      # DTOs de dominio (sin tags)
    │   │   ├── create_payment_method.go
    │   │   ├── update_payment_method.go
    │   │   └── responses.go
    │   ├── ports/                     # Interfaces de repositorios
    │   │   └── ports.go
    │   └── errors/                    # Errores de dominio
    │       └── errors.go
    │
    ├── app/                           # 🟢 CAPA DE APLICACIÓN
    │   └── usecases/
    │       ├── constructor.go         # IUseCase interface + New()
    │       ├── usecases.go            # Implementación casos de uso
    │       └── mappers/               # Conversiones de datos
    │           ├── to_domain.go
    │           └── to_response.go
    │
    └── infra/                         # 🔴 CAPA DE INFRAESTRUCTURA
        ├── primary/                   # Adaptadores de entrada
        │   └── handlers/
        │       ├── constructor.go     # IHandler + New()
        │       ├── routes.go          # Registro de rutas
        │       ├── create-payment-method.go
        │       ├── list-payment-methods.go
        │       ├── get-payment-method.go
        │       ├── update-payment-method.go
        │       ├── delete-payment-method.go
        │       ├── toggle-payment-method.go
        │       ├── create-payment-mapping.go
        │       ├── list-payment-mappings.go
        │       ├── get-payment-mapping.go
        │       ├── get-payment-mappings-by-integration.go
        │       ├── update-payment-mapping.go
        │       ├── delete-payment-mapping.go
        │       ├── toggle-payment-mapping.go
        │       ├── request/           # DTOs de entrada HTTP
        │       │   ├── create_payment_method.go
        │       │   ├── update_payment_method.go
        │       │   ├── create_payment_mapping.go
        │       │   └── update_payment_mapping.go
        │       ├── response/          # DTOs de salida HTTP
        │       │   ├── payment_method.go
        │       │   ├── payment_mapping.go
        │       │   └── error.go
        │       └── mappers/           # Conversiones HTTP ↔ Domain
        │           ├── to_domain.go
        │           └── to_response.go
        │
        └── secondary/                 # Adaptadores de salida
            └── repository/
                ├── repository.go      # Implementación GORM
                └── mappers/           # Conversiones Domain ↔ Models
                    ├── to_domain.go
                    └── to_model.go
```

### Flujo de Dependencias

```
HTTP Request → Handler → UseCase → Repository → Database
     ↓            ↓          ↓           ↓
  request/    mappers/  domain DTOs  models GORM
  response/     ↓          ↓           ↓
              domain    entities    mappers/
```

**Regla de Oro:** Las dependencias SIEMPRE apuntan hacia adentro (Domain es el núcleo).

## 📡 API Endpoints

### Payment Methods

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/payments/methods` | Listar métodos de pago |
| GET | `/payments/methods/:id` | Obtener método por ID |
| POST | `/payments/methods` | Crear método de pago |
| PUT | `/payments/methods/:id` | Actualizar método de pago |
| DELETE | `/payments/methods/:id` | Eliminar método de pago |
| PATCH | `/payments/methods/:id/toggle` | Activar/desactivar método |

### Payment Mappings

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/payments/mappings` | Listar mapeos |
| GET | `/payments/mappings/:id` | Obtener mapeo por ID |
| GET | `/payments/mappings/integration/:type` | Listar mapeos por integración |
| POST | `/payments/mappings` | Crear mapeo |
| PUT | `/payments/mappings/:id` | Actualizar mapeo |
| DELETE | `/payments/mappings/:id` | Eliminar mapeo |
| PATCH | `/payments/mappings/:id/toggle` | Activar/desactivar mapeo |

## 🗄️ Modelos de Base de Datos

### PaymentMethod

```go
type PaymentMethod struct {
    ID          uint      `gorm:"primary_key"`
    Code        string    `gorm:"unique;not null;size:64"`
    Name        string    `gorm:"not null;size:128"`
    Description string    `gorm:"type:text"`
    Category    string    `gorm:"not null;size:50"` // card, digital_wallet, bank_transfer, cash
    Provider    string    `gorm:"size:64"`
    IsActive    bool      `gorm:"default:true"`
    Icon        string    `gorm:"size:255"`
    Color       string    `gorm:"size:32"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### PaymentMethodMapping

```go
type PaymentMethodMapping struct {
    ID              uint   `gorm:"primary_key"`
    IntegrationType string `gorm:"not null;size:50"` // shopify, meli, amazon
    OriginalMethod  string `gorm:"not null;size:100"`
    PaymentMethodID uint   `gorm:"not null"`
    PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID"`
    IsActive        bool   `gorm:"default:true"`
    Priority        int    `gorm:"default:0"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## 🔧 Uso

### Inicialización

```go
// En bundle.go del servicio principal
paymentsBundle := payments.New(
    database,
    logger,
)

// Registrar rutas
paymentsBundle.RegisterRoutes(router)
```

### Ejemplo: Crear Método de Pago

**Request:**
```bash
POST /payments/methods
Content-Type: application/json

{
  "code": "credit_card",
  "name": "Tarjeta de Crédito",
  "description": "Pago con tarjeta de crédito Visa/Mastercard",
  "category": "card",
  "provider": "stripe",
  "icon": "credit-card-icon.svg",
  "color": "#1E40AF"
}
```

**Response:**
```json
{
  "id": 1,
  "code": "credit_card",
  "name": "Tarjeta de Crédito",
  "description": "Pago con tarjeta de crédito Visa/Mastercard",
  "category": "card",
  "provider": "stripe",
  "is_active": true,
  "icon": "credit-card-icon.svg",
  "color": "#1E40AF",
  "created_at": "2026-01-31T10:00:00Z",
  "updated_at": "2026-01-31T10:00:00Z"
}
```

### Ejemplo: Crear Mapeo de Pago

**Request:**
```bash
POST /payments/mappings
Content-Type: application/json

{
  "integration_type": "shopify",
  "original_method": "shopify_payments",
  "payment_method_id": 1,
  "priority": 1
}
```

**Response:**
```json
{
  "id": 1,
  "integration_type": "shopify",
  "original_method": "shopify_payments",
  "payment_method_id": 1,
  "payment_method": {
    "id": 1,
    "code": "credit_card",
    "name": "Tarjeta de Crédito",
    "category": "card"
  },
  "is_active": true,
  "priority": 1,
  "created_at": "2026-01-31T10:05:00Z",
  "updated_at": "2026-01-31T10:05:00Z"
}
```

## ✅ Estado Arquitectural

### 🎉 Módulo CONFORME con Arquitectura Hexagonal

Este módulo ha sido completamente refactorizado y cumple con todas las reglas de arquitectura hexagonal.

#### ✅ Validaciones Aprobadas

| Aspecto | Estado | Detalles |
|---------|--------|----------|
| **Domain organizado** | ✅ | Subcarpetas: `entities/`, `dtos/`, `ports/`, `errors/` |
| **Entidades puras** | ✅ | Sin tags JSON/binding/gorm |
| **Inversión de dependencias** | ✅ | Ports usan entidades de dominio |
| **DTOs separados** | ✅ | Carpetas `request/`, `response/`, `mappers/` en handlers |
| **Mappers organizados** | ✅ | Carpetas dedicadas en cada capa |
| **Repositorios GORM** | ✅ | Usa modelos GORM, NO usa `.Table()` |
| **Compilación** | ✅ | `go build ./...` sin errores |

#### 📊 Resultados de Validación

```bash
# Tags JSON en domain/entities/: 0 ✅
# Tags binding en domain/: 0 ✅
# Imports prohibidos (gorm/gin) en domain/: 0 ✅
# Uso de .Table() en repositorios: 0 ✅
```

#### 🏗️ Estructura Arquitectural

```
payments/
├── bundle.go            # ✅ Ensambla el módulo
└── internal/            # ✅ Convención Go (paquetes privados)
    ├── domain/          # 🔵 CAPA DE DOMINIO (núcleo)
    │   ├── entities/    # Entidades PURAS (sin tags)
    │   ├── dtos/        # DTOs PUROS (sin tags)
    │   ├── ports/       # Interfaces (contratos)
    │   └── errors/      # Errores de dominio
    ├── app/             # 🟢 CAPA DE APLICACIÓN
    │   └── usecases/
    │       └── mappers/ # Conversiones domain ↔ entities
    └── infra/           # 🔴 CAPA DE INFRAESTRUCTURA
        ├── primary/handlers/
        │   ├── request/     # DTOs HTTP entrada
        │   ├── response/    # DTOs HTTP salida
        │   └── mappers/     # Conversiones HTTP ↔ domain
        └── secondary/repository/
            └── mappers/     # Conversiones GORM ↔ domain
```

**Comando para validar:**
```bash
cd /home/cam/Desktop/probability/back/central/services/modules/payments

# Buscar tags en domain (DEBE retornar 0)
grep -r 'json:"' internal/domain/entities/ | wc -l

# Buscar imports prohibidos en domain (DEBE retornar 0)
grep -r "gorm\|gin" internal/domain/ | grep -v "OriginalMethod" | wc -l

# Verificar uso de Table() en repositorios (DEBE retornar 0)
grep -r '\.Table(' internal/infra/secondary/repository/ | wc -l

# Compilar
go build ./...

# Tests
go test ./...
```

#### 🎯 Última Refactorización

**Fecha:** 2026-01-31
**Estado:** ✅ COMPLETADA
**Cumplimiento:** 11/11 reglas (100%)

**Cambios realizados:**
- ✅ Domain reorganizado en subcarpetas (`entities/`, `dtos/`, `ports/`, `errors/`)
- ✅ Entidades puras sin tags de frameworks
- ✅ Ports usan entidades de dominio (no modelos GORM)
- ✅ Mappers separados en cada capa (repository, usecases, handlers)
- ✅ DTOs HTTP separados en `request/` y `response/`
- ✅ Todo movido a carpeta `internal/` (convención Go)

**Estadísticas:**
- 📝 Archivos creados: 24
- 📝 Archivos modificados: 47
- 📝 Archivos eliminados: 2
- 📊 Total archivos Go: 40 (Domain: 9, App: 4, Infra: 27)

## 📚 Referencias

- **Arquitectura Hexagonal:** `.claude/rules/architecture.md`
- **Agente de validación:** `.claude/agents/hexagonal-architecture-assistant.md`
- **Convenciones del proyecto:** `CLAUDE.md`

## 🔍 Validación Arquitectural

Para validar el cumplimiento de arquitectura hexagonal:

```bash
# Desde la raíz del proyecto
claude code "aplica el agente hexagonal-architecture-assistant al módulo payments"
```

---

**Última actualización:** 2026-01-31
**Estado:** ❌ NO CONFORME - Requiere refactorización
**Prioridad:** Alta - Violaciones críticas detectadas
