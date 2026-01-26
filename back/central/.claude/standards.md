# Estándares Go - Central Reserve

## 1. Manejo de Errores

### Errores de Dominio
Definir en `internal/domain/errors.go`:
```go
var (
    ErrUserNotFound = errors.New("usuario no encontrado")
    ErrEmailExists  = errors.New("el email ya está registrado")
)
```

### Uso de errors.Is() en Handlers
```go
switch {
case errors.Is(err, domain.ErrUserNotFound):
    c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
case errors.Is(err, domain.ErrEmailExists):
    c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
default:
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno"})
}
```

### Wrapping de Errores
```go
return fmt.Errorf("error al crear usuario: %w", err)
```

## 2. Logging con Zerolog

### Patrón Estándar
```go
// Info con campos
h.logger.Info().
    Str("email", req.Email).
    Uint("user_id", userID).
    Msg("Usuario creado exitosamente")

// Error con error original
h.logger.Error().Err(err).Msg("Error al crear usuario")
```

### Logger Contextual
```go
// En constructor del handler
contextualLogger := logger.WithModule("usuarios")
```

## 3. Interfaces (Puertos)

### Ubicación
Siempre en `internal/domain/ports.go`

### Convención de Nombres
- Prefijo `I`: `IUserRepository`, `ILogger`
- Interfaces de caso de uso: `Iapp`
- Interfaces de handler: `Ihandlers`

### Ejemplo
```go
type IUserRepository interface {
    GetUserByID(ctx context.Context, id uint) (*User, error)
    CreateUser(ctx context.Context, user UserEntity) (uint, error)
}
```

## 4. Constructores

### Patrón
```go
func New(deps...) InterfaceType {
    return &implementation{deps}
}
```

### Ejemplo
```go
func New(db db.IDatabase, logger log.ILogger) domain.IUserRepository {
    return &Repository{database: db, logger: logger}
}
```

## 5. DTOs

### Ubicación
- Request: `infra/primary/handlers/request/`
- Response: `infra/primary/handlers/response/`
- Domain: `internal/domain/dtos.go`

### Mappers
En `infra/primary/handlers/mapper/`:
```go
func ToCreateUserDTO(req CreateUserRequest) domain.CreateUserDTO
func ToUserResponse(dto domain.UserDTO) response.UserResponse
```

## 6. Contextos

### Propagación
Siempre pasar `ctx context.Context` como primer parámetro:
```go
func (uc *UseCase) CreateUser(ctx context.Context, dto CreateDTO) error {
    return uc.repo.Create(ctx, entity)
}
```

## 7. Convenciones de Archivos

| Tipo | Patrón | Ejemplo |
|------|--------|---------|
| Caso de uso | kebab-case | `create-user.go` |
| Repository | snake_case | `user_repository.go` |
| Handler | kebab-case | `create-user.go` |
| Constructor | - | `constructor.go` |
| Router | - | `router.go` |
