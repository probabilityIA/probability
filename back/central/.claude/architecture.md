# Arquitectura Hexagonal - Central Reserve

## Estructura del Proyecto

```
central-reserve/
├── cmd/                    # Punto de entrada
│   ├── main.go            # Bootstrap
│   └── internal/
│       ├── server/init.go # Inicialización
│       └── routes/        # Enrutamiento HTTP
├── services/              # Módulos de negocio
│   ├── auth/              # Autenticación
│   ├── horizontalproperty/# Propiedades
│   └── restaurants/       # Restaurantes
└── shared/                # Código compartido
    ├── db/                # Abstracción BD
    ├── log/               # Logger (zerolog)
    ├── jwt/               # Tokens JWT
    ├── errs/              # Errores custom
    └── storage/           # S3/AWS
```

## Estructura de un Servicio (Arquitectura Hexagonal)

```
services/{servicio}/{subdominio}/
├── internal/
│   ├── domain/           # DOMINIO (núcleo)
│   │   ├── entities.go   # Entidades de negocio
│   │   ├── ports.go      # Interfaces/contratos
│   │   ├── errors.go     # Errores del dominio
│   │   └── dtos.go       # Data Transfer Objects
│   │
│   ├── app/              # APLICACIÓN (casos de uso)
│   │   ├── constructor.go
│   │   └── {accion}.go   # create-user.go, etc.
│   │
│   └── infra/            # INFRAESTRUCTURA
│       ├── primary/      # Adaptadores entrada (HTTP)
│       │   └── handlers/
│       │       ├── constructor.go
│       │       ├── router.go
│       │       ├── mapper/
│       │       ├── request/
│       │       └── response/
│       │
│       └── secondary/    # Adaptadores salida (BD)
│           └── repository/
│
└── bundle.go             # Inicializador del módulo
```

## Flujo de una Petición

```
HTTP Request → Middleware JWT → Handler (infra/primary)
                                    ↓
                              UseCase (app)
                                    ↓
                            Repository (infra/secondary)
                                    ↓
                              Domain (entities)
                                    ↓
                              Response JSON
```

## Inyección de Dependencias

```go
// En bundle.go de cada servicio
func New(db db.IDatabase, logger log.ILogger, router *gin.RouterGroup) {
    // 1. Repositorios (infra/secondary)
    repo := repository.New(db, logger)

    // 2. Casos de uso (app)
    useCase := app.New(repo, logger)

    // 3. Handlers (infra/primary)
    handler := handlers.New(useCase, logger)

    // 4. Registrar rutas
    handler.RegisterRoutes(router)
}
```

## Reglas de Dependencia

- Domain NO depende de nada externo
- App depende SOLO de Domain (interfaces)
- Infra implementa las interfaces de Domain
- Las dependencias fluyen HACIA el centro (Domain)

## Reglas de Organización de Archivos

Ver documento completo en `.claude/file-organization-rules.md`

### Resumen:

**Handlers DEBEN tener**:
- `handlers/request/` - DTOs de entrada
- `handlers/response/` - DTOs de salida
- `handlers/mappers/` - Conversiones request/response ↔ domain

**Repositorios DEBEN tener**:
- `repository/mappers/` - Conversiones models ↔ domain
  - `to_domain.go` - models → domain
  - `to_model.go` - domain → models

**Regla**: NO definir mappers inline en archivos de handlers o repositorios. Centralizar en carpetas `mappers/`.
