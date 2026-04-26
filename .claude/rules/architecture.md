# Arquitectura Hexagonal

Dependencias: Infra -> App -> Domain. Domain nunca conoce capas externas.

## Backend (Go)

```
service/
├── bundle.go
└── internal/
    ├── domain/
    │   ├── entities/
    │   ├── dtos/
    │   ├── ports/
    │   └── errors/
    ├── app/
    │   ├── constructor.go
    │   ├── usecase.go
    │   ├── request/
    │   ├── response/
    │   └── mappers/
    └── infra/
        ├── primary/handlers/
        │   ├── constructor.go
        │   ├── routes.go
        │   ├── create_handler.go
        │   ├── request/
        │   ├── response/
        │   └── mappers/
        └── secondary/repository/
            ├── constructor.go
            ├── repository.go
            └── mappers/
```

**Reglas:**
- Domain: sin tags (`json:`, `gorm:`). Imports: solo `context`, `time`, `errors`, `fmt`, `uuid`. Prohibidos: gorm, gin, net/http
- Un `constructor.go` + una funcion `New()` por carpeta. Una interfaz publica.
- Handlers: un metodo por archivo + `routes.go` con `RegisterRoutes()`
- Repos: modelos de `migration/shared/models` (NUNCA `.Table("users")`, NUNCA `models/` local en modulos)
- Ports SOLO en `internal/domain/ports/`, nunca en raiz del modulo
- GET listas: paginacion obligatoria. `PaginationParams`/`PaginatedResponse` en `domain/dtos/`. Default page=1, pageSize=10, max=100. Response: `{ data, total, page, page_size, total_pages }`. Excepcion: lookups por ID, catalogos < 50 registros.

**Violaciones criticas:** domain con tags/gorm/gin | multiples constructores | DTOs en handlers | `.Table()` | models/ local | ports.go en raiz | GET sin paginacion

## Frontend (Next.js)

```
services/[module]/
├── domain/
│   ├── types.ts
│   └── ports.ts
├── app/
│   └── use-cases.ts
├── infra/
│   ├── repository/
│   └── actions/
└── ui/
    ├── components/
    └── hooks/
```

**Decision:**
- Leer datos -> Server Component
- Mutacion -> Server Action (`'use server'` + `revalidatePath()`)
- Interactividad/tiempo real -> Client Component (+ WebSocket/SSE si aplica)

**GET listas:** paginacion en URL `?page=1&page_size=10`, `cache: 'no-store'`. `PaginationParams`/`PaginatedResponse` en `domain/types.ts`.

**Violaciones criticas:** fetch en domain | Client Component con fetch | useEffect para fetch inicial | mutaciones sin Server Actions | listados sin paginacion
