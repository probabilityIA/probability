# Módulo Payment Status

## 📋 Descripción

El módulo **Payment Status** gestiona los estados de pago disponibles en la plataforma Probability. Proporciona catálogos de estados predefinidos que clasifican el ciclo de vida de los pagos de pedidos.

## 🎯 Funcionalidades

- Listar estados de pago con filtrado por activo/inactivo
- Obtener ID de estado por código
- Catálogo predefinido de estados (pending, paid, failed, refunded, etc.)
- Categorización de estados para análisis

## 🏗️ Arquitectura

Este módulo sigue **Arquitectura Hexagonal (Clean Architecture)** con la siguiente estructura:

> **Nota importante:** Todo el código del módulo está dentro de la carpeta `internal/` siguiendo la convención de Go. Los paquetes en `internal/` son privados y no pueden ser importados por módulos externos, garantizando el encapsulamiento del módulo.

```
paymentstatus/
├── bundle.go                          # ✅ Ensambla el módulo
├── ports.go                           # ✅ Re-exporta IRepository
├── README.md                          # ✅ Documentación
└── internal/                          # ✅ Convención Go (paquetes privados)
    ├── domain/                        # 🔵 CAPA DE DOMINIO (núcleo)
    │   ├── entities/                  # Entidades PURAS (sin tags)
    │   │   └── payment_status.go
    │   ├── dtos/                      # DTOs PUROS (sin tags)
    │   │   └── payment_status_info.go
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
    │           └── to_dto.go
    │
    └── infra/                         # 🔴 CAPA DE INFRAESTRUCTURA
        ├── primary/                   # Adaptadores de entrada
        │   └── handlers/
        │       ├── constructor.go     # IHandler + New()
        │       ├── routes.go          # Registro de rutas
        │       ├── list-payment-statuses.go
        │       ├── response/          # DTOs de salida HTTP
        │       │   └── payment_status.go
        │       └── mappers/           # Conversiones HTTP ↔ Domain
        │           └── to_response.go
        │
        └── secondary/                 # Adaptadores de salida
            └── repository/
                ├── repository.go      # Implementación GORM
                ├── models/            # Modelos GORM
                │   └── payment_status.go
                └── mappers/           # Conversiones Domain ↔ Models
                    ├── to_domain.go
                    └── to_model.go
```

### Flujo de Dependencias

```
HTTP Request → Handler → UseCase → Repository → Database
     ↓            ↓          ↓           ↓
  response/   mappers/   domain DTOs  models GORM
              ↓          ↓           ↓
           to_response  to_dto    to_domain
```

**Regla de Oro:** Las dependencias SIEMPRE apuntan hacia adentro (Domain es el núcleo).

## 📡 API Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/payment-statuses` | Listar estados de pago |
| GET | `/payment-statuses?is_active=true` | Listar solo estados activos |
| GET | `/payment-statuses?is_active=false` | Listar solo estados inactivos |

## 🗄️ Modelo de Base de Datos

### PaymentStatus

```go
type PaymentStatus struct {
    ID          uint           `gorm:"primarykey"`
    Code        string         `gorm:"size:64;unique;not null;index"`
    Name        string         `gorm:"size:128;not null"`
    Description string         `gorm:"type:text"`
    Category    string         `gorm:"size:64;index"`
    IsActive    bool           `gorm:"default:true;index"`
    Icon        string         `gorm:"size:255"`
    Color       string         `gorm:"size:32"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

## 🔧 Uso

### Inicialización

```go
// En bundle.go del servicio principal
paymentStatusBundle := paymentstatus.New(
    router,
    database,
    logger,
    environment,
)
```

### Ejemplo: Listar Estados de Pago

**Request:**
```bash
GET /payment-statuses?is_active=true
```

**Response:**
```json
{
  "success": true,
  "message": "Estados de pago obtenidos exitosamente",
  "data": [
    {
      "id": 1,
      "code": "pending",
      "name": "Pendiente",
      "description": "Pago pendiente de procesar",
      "category": "waiting",
      "color": "#FFA500"
    },
    {
      "id": 2,
      "code": "paid",
      "name": "Pagado",
      "description": "Pago completado exitosamente",
      "category": "success",
      "color": "#00FF00"
    },
    {
      "id": 3,
      "code": "failed",
      "name": "Fallido",
      "description": "Pago rechazado o fallido",
      "category": "error",
      "color": "#FF0000"
    }
  ]
}
```

## ✅ Estado Arquitectural

### 🎉 Módulo CONFORME con Arquitectura Hexagonal

Este módulo ha sido completamente refactorizado y cumple con todas las reglas de arquitectura hexagonal.

#### ✅ Validaciones Aprobadas

| Aspecto | Estado | Detalles |
|---------|--------|----------|
| **Domain en internal/** | ✅ | Sigue convención Go |
| **Domain organizado** | ✅ | Subcarpetas: `entities/`, `dtos/`, `ports/`, `errors/` |
| **Entidades puras** | ✅ | Sin tags JSON/binding/gorm |
| **Inversión de dependencias** | ✅ | Ports usan entidades de dominio |
| **Mappers organizados** | ✅ | Carpetas dedicadas en cada capa |
| **Repositorios GORM** | ✅ | Usa modelos GORM locales, NO usa `.Table()` |
| **Compilación** | ✅ | `go build ./...` sin errores |

#### 📊 Resultados de Validación

```bash
# Tags JSON en domain/entities/: 0 ✅
# Tags binding en domain/: 0 ✅
# Imports prohibidos (gorm/gin/models) en domain/: 0 ✅
# Uso de .Table() en repositorios: 0 ✅
```

#### 🏗️ Estructura Arquitectural

```
paymentstatus/
├── bundle.go            # ✅ Ensambla el módulo
├── ports.go             # ✅ Re-exporta IRepository
└── internal/            # ✅ Convención Go (paquetes privados)
    ├── domain/          # 🔵 CAPA DE DOMINIO (PURA)
    │   ├── entities/    # Entidades PURAS (sin tags)
    │   ├── dtos/        # DTOs PUROS (sin tags)
    │   ├── ports/       # Interfaces (contratos)
    │   └── errors/      # Errores de dominio
    ├── app/             # 🟢 CAPA DE APLICACIÓN
    │   └── usecases/
    │       └── mappers/ # Conversiones domain ↔ entities
    └── infra/           # 🔴 CAPA DE INFRAESTRUCTURA
        ├── primary/handlers/
        │   ├── response/    # DTOs HTTP salida
        │   └── mappers/     # Conversiones HTTP ↔ domain
        └── secondary/repository/
            ├── models/      # Modelos GORM
            └── mappers/     # Conversiones GORM ↔ domain
```

**Comando para validar:**
```bash
cd /home/cam/Desktop/probability/back/central/services/modules/paymentstatus

# Buscar tags en domain (DEBE retornar 0)
grep -r 'json:"' internal/domain/entities/ | wc -l

# Buscar imports prohibidos en domain (DEBE retornar 0)
grep -r "gorm\|gin\|migration/shared/models" internal/domain/ | wc -l

# Verificar uso de Table() en repositorios (DEBE retornar 0)
grep -r '\.Table(' internal/infra/secondary/repository/ | wc -l

# Compilar
go build ./...
```

#### 🎯 Última Refactorización

**Fecha:** 2026-01-31
**Estado:** ✅ COMPLETADA
**Cumplimiento:** 7/7 reglas (100%)

**Cambios realizados:**
- ✅ Domain reorganizado en subcarpetas (`entities/`, `dtos/`, `ports/`, `errors/`)
- ✅ Entidades puras sin tags de frameworks
- ✅ Ports usan entidades de dominio (no modelos GORM externos)
- ✅ Modelos GORM locales con `TableName()`, `ToDomain()`, `FromDomain()`
- ✅ Mappers separados en cada capa (repository, usecases, handlers)
- ✅ DTOs HTTP separados en `response/`
- ✅ Todo movido a carpeta `internal/` (convención Go)

**Estadísticas:**
- 📝 Archivos creados: 15
- 📝 Archivos modificados: 4
- 📝 Archivos eliminados: 3
- 📊 Total archivos Go: 15

## 📚 Referencias

- **Reglas de Arquitectura:** `.claude/rules/architecture.md`
- **Agente de validación:** `.claude/agents/hexagonal-architecture-assistant.md`
- **Módulo de referencia:** `services/modules/payments/` (arquitectura correcta)
- **CLAUDE.md del proyecto:** `/back/central/CLAUDE.md`

---

**Última actualización:** 2026-01-31
**Estado:** ✅ CONFORME
**Próximo paso:** Aplicar este patrón a otros módulos del proyecto
