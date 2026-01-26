# Reglas de Organización de Archivos - Arquitectura Hexagonal

## 🎯 Objetivo

Mantener una estructura consistente y predecible en todos los módulos, facilitando la navegación del código y el mantenimiento.

---

## 📁 Estructura Obligatoria de Handlers

Los handlers (adaptadores primarios HTTP) DEBEN seguir esta estructura:

```
infra/primary/handlers/
├── constructor.go          # Constructor del handler
├── router.go              # Registro de rutas
├── {accion}.go            # Handlers individuales (create-user.go, etc.)
├── request/               # ✅ OBLIGATORIO - DTOs de entrada
│   ├── create-user.go
│   ├── update-user.go
│   └── filters.go
├── response/              # ✅ OBLIGATORIO - DTOs de salida
│   ├── user.go
│   ├── paginated.go
│   └── error.go
└── mappers/               # ✅ OBLIGATORIO - Conversiones request/response ↔ domain
    ├── request.go         # Mappers de request → domain
    └── response.go        # Mappers de domain → response
```

### Reglas:

1. **`request/`** - OBLIGATORIO
   - Contiene TODOS los DTOs de entrada (structs que reciben datos del cliente)
   - Validaciones con tags: `json`, `validate`, `binding`
   - Naming: `{accion}.go` o `{entidad}.go`
   - Ejemplo: `CreateUserRequest`, `UpdateUserRequest`, `UserFiltersRequest`

2. **`response/`** - OBLIGATORIO
   - Contiene TODOS los DTOs de salida (structs que se retornan al cliente)
   - Solo campos que se exponen en la API
   - Naming: `{entidad}.go`
   - Ejemplo: `UserResponse`, `PaginatedUsersResponse`

3. **`mappers/`** - OBLIGATORIO
   - Todas las funciones de conversión entre capas
   - `request.go`: Funciones `ToXXXDTO(req) domain.DTO`
   - `response.go`: Funciones `ToXXXResponse(domain) response.XXX`
   - NO mezclar lógica de negocio, solo conversión de estructuras

### Ejemplo de Uso:

```go
// En handler create-user.go
package handlers

import (
    "central_reserve/services/auth/users/internal/domain"
    "central_reserve/services/auth/users/internal/infra/primary/handlers/mappers"
    "central_reserve/services/auth/users/internal/infra/primary/handlers/request"
    "central_reserve/services/auth/users/internal/infra/primary/handlers/response"
)

func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. Parsear request
    var req request.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 2. Convertir a DTO de dominio usando mapper
    dto := mappers.ToCreateUserDTO(req)

    // 3. Llamar caso de uso
    user, err := h.useCase.CreateUser(c.Request.Context(), dto)
    if err != nil {
        // manejar error
        return
    }

    // 4. Convertir a response usando mapper
    resp := mappers.ToUserResponse(user)

    c.JSON(201, resp)
}
```

---

## 📁 Estructura Obligatoria de Repositorios

Los repositorios (adaptadores secundarios) DEBEN seguir esta estructura:

```
infra/secondary/repository/
├── constructor.go              # Constructor del repositorio
├── {entidad}_repository.go     # Implementación del repositorio
└── mappers/                    # ✅ OBLIGATORIO - Conversiones models ↔ domain
    ├── to_domain.go            # Mappers de models → domain entities
    └── to_model.go             # Mappers de domain entities → models
```

### Reglas:

1. **`mappers/`** - OBLIGATORIO
   - Contiene TODAS las funciones de conversión entre modelos de DB y entidades de dominio
   - `to_domain.go`: Funciones `MapXXXToDomain(model) *domain.Entity`
   - `to_model.go`: Funciones `MapXXXToModel(entity) *models.Model`
   - Centraliza la lógica de mapeo, evitando duplicación

2. **Naming de funciones**:
   - `Map{Entidad}ToDomain(m *models.Model) *domain.Entity`
   - `Map{Entidad}ToModel(e *domain.Entity) *models.Model`
   - Ejemplo: `MapUserToDomain`, `MapVisitToModel`

3. **Ubicación de mappers inline**:
   - ❌ NO definir funciones `mapXXX()` directamente en el repositorio
   - ✅ SÍ extraer a `mappers/`
   - Excepción: Conversiones triviales de 1-2 líneas pueden quedar inline

### Ejemplo de Estructura:

```go
// repository/mappers/to_domain.go
package mappers

import (
    "central_reserve/services/horizontalproperty/visit/internal/domain"
    "dbpostgres/app/infra/models"
)

func MapVisitToDomain(m *models.Visit) *domain.Visit {
    return &domain.Visit{
        ID:             m.ID,
        BusinessID:     m.BusinessID,
        VisitorID:      m.VisitorID,
        // ... resto de campos
    }
}

func MapVisitTypeToDomain(m *models.VisitType) *domain.VisitType {
    return &domain.VisitType{
        ID:                    m.ID,
        Name:                  m.Name,
        RequiresAuthorization: m.RequiresAuthorization,
        // ... resto
    }
}
```

```go
// repository/mappers/to_model.go
package mappers

func MapVisitToModel(v *domain.Visit) *models.Visit {
    return &models.Visit{
        ID:             v.ID,
        BusinessID:     v.BusinessID,
        VisitorID:      v.VisitorID,
        // ... resto de campos
    }
}
```

```go
// repository/visit_repository.go
package repository

import (
    "central_reserve/services/horizontalproperty/visit/internal/domain"
    "central_reserve/services/horizontalproperty/visit/internal/infra/secondary/repository/mappers"
)

func (r *VisitRepository) GetVisitByID(ctx context.Context, id uint) (*domain.Visit, error) {
    var visit models.Visit
    if err := r.db.Conn(ctx).Preload("VisitStatus").First(&visit, id).Error; err != nil {
        return nil, err
    }

    // Usar mapper de la carpeta mappers/
    return mappers.MapVisitToDomain(&visit), nil
}

func (r *VisitRepository) CreateVisit(ctx context.Context, visit *domain.Visit) (*domain.Visit, error) {
    // Usar mapper para convertir a modelo
    model := mappers.MapVisitToModel(visit)

    if err := r.db.Conn(ctx).Create(model).Error; err != nil {
        return nil, err
    }

    // Convertir de vuelta a dominio
    return mappers.MapVisitToDomain(model), nil
}
```

---

## ✅ Validaciones de Arquitectura

Un módulo cumple con las reglas de organización si:

### Handlers:
- [ ] Existe carpeta `handlers/request/` con al menos 1 archivo
- [ ] Existe carpeta `handlers/response/` con al menos 1 archivo
- [ ] Existe carpeta `handlers/mappers/` con al menos 1 archivo
- [ ] Los handlers usan los DTOs de `request/` y `response/`
- [ ] Los handlers usan los mappers de `mappers/`
- [ ] NO hay mappers inline en archivos de handlers (excepto triviales)

### Repositorios:
- [ ] Existe carpeta `repository/mappers/` con archivos `to_domain.go` y/o `to_model.go`
- [ ] Los repositorios usan los mappers de `mappers/`
- [ ] NO hay funciones `mapXXX()` definidas directamente en el archivo del repositorio
- [ ] Todas las conversiones `models ↔ domain` pasan por mappers

---

## 🚨 Violaciones Comunes

### ❌ Violación: Mappers inline en repositorio

```go
// ❌ MAL - Función de mapeo inline en visit_repository.go
func mapVisitToDomain(m *models.Visit) *domain.Visit {
    return &domain.Visit{...}
}

func (r *VisitRepository) GetVisitByID(ctx, id) (*domain.Visit, error) {
    var visit models.Visit
    // ...
    return mapVisitToDomain(&visit), nil
}
```

### ✅ Corrección:

```go
// ✅ BIEN - Mover a repository/mappers/to_domain.go
package mappers

func MapVisitToDomain(m *models.Visit) *domain.Visit {
    return &domain.Visit{...}
}

// En repository/visit_repository.go
import "central_reserve/.../repository/mappers"

func (r *VisitRepository) GetVisitByID(ctx, id) (*domain.Visit, error) {
    var visit models.Visit
    // ...
    return mappers.MapVisitToDomain(&visit), nil
}
```

### ❌ Violación: Carpetas request/response faltantes

```
handlers/
├── create-user.go         # ❌ Define structs inline
├── update-user.go
└── router.go
```

### ✅ Corrección:

```
handlers/
├── create-user.go         # ✅ Importa de request/response
├── request/
│   ├── create-user.go
│   └── update-user.go
├── response/
│   └── user.go
└── mappers/
    ├── request.go
    └── response.go
```

---

## 📝 Checklist de Migración

Para adaptar un módulo existente a estas reglas:

1. **Handlers**:
   - [ ] Crear carpeta `handlers/request/`
   - [ ] Mover/crear DTOs de request
   - [ ] Crear carpeta `handlers/response/`
   - [ ] Mover/crear DTOs de response
   - [ ] Crear carpeta `handlers/mappers/`
   - [ ] Crear `mappers/request.go` y `mappers/response.go`
   - [ ] Mover funciones de conversión a mappers
   - [ ] Actualizar imports en handlers

2. **Repositorios**:
   - [ ] Crear carpeta `repository/mappers/`
   - [ ] Crear `mappers/to_domain.go`
   - [ ] Crear `mappers/to_model.go`
   - [ ] Mover funciones `mapXXXToDomain` a `to_domain.go`
   - [ ] Mover funciones `mapXXXToModel` a `to_model.go`
   - [ ] Actualizar imports en repositorio
   - [ ] Eliminar funciones inline del repositorio

---

## 🎓 Beneficios

1. **Consistencia**: Todos los módulos siguen la misma estructura
2. **Navegación**: Fácil encontrar dónde están los DTOs y mappers
3. **Reutilización**: Los mappers centralizados pueden compartirse
4. **Testing**: Más fácil testear mappers aisladamente
5. **Mantenimiento**: Cambios en DTOs se localizan en un solo lugar
6. **Onboarding**: Nuevos desarrolladores entienden la estructura rápidamente
