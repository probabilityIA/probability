# Frontend Central - Dashboard Admin (Next.js)

## Descripción

Dashboard administrativo para gestionar pedidos, productos, usuarios y configuraciones de la plataforma Probability.

## Stack Tecnológico

- **Next.js 16.2** (App Router)
- **React 19.2**
- **TypeScript 5**
- **TailwindCSS 4**
- **pnpm** (package manager)

### Dependencias Principales

| Librería | Uso |
|----------|-----|
| `@heroicons/react` | Iconos |
| `date-fns` | Manipulación de fechas |
| `recharts` | Gráficos y visualizaciones |
| `leaflet` / `react-leaflet` | Mapas |
| `react-day-picker` | Selector de fechas |
| `react-colorful` | Selector de colores |

## Arquitectura por Módulos

El proyecto aplica **Arquitectura Hexagonal adaptada al frontend**:

```
service/
├── domain/           # Tipos, interfaces, contratos
│   ├── types.ts      # Entidades y DTOs
│   └── ports.ts      # Interfaces de repositorio
├── app/              # Lógica de aplicación
│   └── use-cases.ts  # Casos de uso
├── infra/            # Implementaciones
│   ├── repository/   # Llamadas a API
│   │   └── api-repository.ts
│   └── actions/      # Server Actions (Next.js)
│       └── index.ts
└── ui/               # Presentación
    ├── components/   # Componentes React
    ├── hooks/        # Custom hooks
    └── styles/       # Estilos específicos
```

## Estructura del Proyecto

```
src/
├── app/                      # App Router (Next.js 13+)
│   ├── layout.tsx            # Layout raíz
│   ├── page.tsx              # Página principal
│   └── (auth)/               # Grupo de rutas autenticadas
│       ├── layout.tsx        # Layout con sidebar
│       ├── layout-content.tsx
│       ├── login/page.tsx
│       ├── home/page.tsx
│       ├── orders/page.tsx
│       ├── products/page.tsx
│       ├── shipments/page.tsx
│       ├── users/page.tsx
│       ├── roles/page.tsx
│       ├── permissions/page.tsx
│       ├── resources/page.tsx
│       ├── businesses/page.tsx
│       ├── integrations/page.tsx
│       ├── order-status/page.tsx
│       └── notification-config/page.tsx
├── services/                 # Módulos de negocio
│   ├── auth/                 # Autenticación
│   │   ├── login/
│   │   ├── users/
│   │   ├── roles/
│   │   ├── permissions/
│   │   ├── resources/
│   │   └── business/
│   ├── modules/              # Módulos core
│   │   ├── orders/
│   │   ├── products/
│   │   ├── shipments/
│   │   ├── dashboard/
│   │   ├── orderstatus/
│   │   ├── fulfillmentstatus/
│   │   └── notification-config/
│   ├── integrations/         # Integraciones
│   │   └── core/
│   └── transport/
└── shared/                   # Código compartido
    ├── ui/                   # Componentes UI reutilizables
    ├── hooks/                # Hooks compartidos
    ├── contexts/             # React Contexts
    ├── providers/            # Providers
    ├── config/               # Configuración
    └── utils/                # Utilidades
```

## Servicios

### Auth (`services/auth/`)

| Módulo | Descripción |
|--------|-------------|
| `login/` | Autenticación, cambio de contraseña |
| `users/` | CRUD de usuarios |
| `roles/` | Gestión de roles |
| `permissions/` | Gestión de permisos |
| `resources/` | Recursos del sistema |
| `business/` | Gestión de negocios |

### Modules (`services/modules/`)

| Módulo | Descripción |
|--------|-------------|
| `orders/` | Listado y detalle de pedidos |
| `products/` | Catálogo de productos |
| `shipments/` | Gestión de envíos |
| `dashboard/` | Métricas y estadísticas |
| `orderstatus/` | Mapeo de estados de pedidos |
| `notification-config/` | Configuración de notificaciones |

### Integrations (`services/integrations/`)

| Módulo | Descripción |
|--------|-------------|
| `core/` | Gestión de integraciones (Shopify, WhatsApp, etc.) |

## Shared (`shared/`)

### UI Components (`shared/ui/`)

Componentes reutilizables:

- `button.tsx` - Botones
- `input.tsx` - Campos de texto
- `select.tsx` - Selectores
- `table.tsx` - Tablas de datos
- `modal.tsx` - Modales
- `form-modal.tsx` - Modales con formularios
- `confirm-modal.tsx` - Modales de confirmación
- `spinner.tsx` - Loading spinner
- `toast.tsx` - Notificaciones
- `badge.tsx` - Badges
- `accordion.tsx` - Acordeones
- `date-picker.tsx` - Selector de fecha
- `date-range-picker.tsx` - Rango de fechas
- `sidebar.tsx` - Sidebar principal
- `avatar-upload.tsx` - Upload de avatar
- `dynamic-filters.tsx` - Filtros dinámicos

### Hooks (`shared/hooks/`)

- `use-sse.ts` - Server-Sent Events

### Contexts (`shared/contexts/`)

- `sidebar-context.tsx` - Estado del sidebar
- `permissions-context.tsx` - Permisos del usuario

### Providers (`shared/providers/`)

- `theme-provider.tsx` - Tema de la aplicación
- `toast-provider.tsx` - Sistema de notificaciones

### Utils (`shared/utils/`)

- `token-storage.ts` - Gestión de tokens JWT
- `http-logger.ts` - Logger HTTP
- `apply-business-theme.ts` - Tema por negocio
- `sound.ts` - Efectos de sonido

### Config (`shared/config/`)

- `env.ts` - Variables de entorno

## Patrones de Código

### Domain - Ports (Interfaces)

```typescript
// domain/ports.ts
export interface IUserRepository {
    getUsers(params?: GetUsersParams): Promise<PaginatedResponse<User>>;
    getUserById(id: number): Promise<SingleResponse<User>>;
    createUser(data: CreateUserDTO): Promise<SingleResponse<User>>;
    updateUser(id: number, data: UpdateUserDTO): Promise<SingleResponse<User>>;
    deleteUser(id: number): Promise<ActionResponse>;
}
```

### Domain - Types (Entidades)

```typescript
// domain/types.ts
export interface User {
    id: number;
    name: string;
    email: string;
    role: string;
    // ...
}

export interface CreateUserDTO {
    name: string;
    email: string;
    password: string;
}

export interface PaginatedResponse<T> {
    data: T[];
    total: number;
    page: number;
    pageSize: number;
}
```

### Infra - Repository (API Client)

```typescript
// infra/repository/api-repository.ts
export class UserApiRepository implements IUserRepository {
    async getUsers(params?: GetUsersParams): Promise<PaginatedResponse<User>> {
        const response = await fetch(`/api/users?${new URLSearchParams(params)}`);
        return response.json();
    }
    // ...
}
```

### UI - Hooks

```typescript
// ui/hooks/useUsers.ts
export function useUsers() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);

    const fetchUsers = async (params?: GetUsersParams) => {
        setLoading(true);
        const repo = new UserApiRepository();
        const response = await repo.getUsers(params);
        setUsers(response.data);
        setLoading(false);
    };

    return { users, loading, fetchUsers };
}
```

### UI - Components

```typescript
// ui/components/UserList.tsx
'use client';

import { useUsers } from '../hooks/useUsers';
import { Table, Button, Spinner } from '@/shared/ui';

export function UserList() {
    const { users, loading, fetchUsers } = useUsers();

    if (loading) return <Spinner />;

    return (
        <Table
            data={users}
            columns={[...]}
            onRowClick={(user) => {...}}
        />
    );
}
```

## Autenticación

El sistema usa JWT con dos tokens:
- **Session Token**: Token de sesión del usuario
- **Business Token**: Token del negocio seleccionado

```typescript
// shared/config/token-storage.ts
TokenStorage.getSessionToken();
TokenStorage.getBusinessToken();
TokenStorage.getUser();
TokenStorage.getBusinessesData();
TokenStorage.clearSession();
```

## Comandos

```bash
# Instalar dependencias
pnpm install

# Desarrollo
pnpm dev

# Build producción
pnpm build

# Iniciar producción
pnpm start

# Linter
pnpm lint
```

## Variables de Entorno

```env
# API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Otras configuraciones
NEXT_PUBLIC_APP_NAME=Probability Central
```

## Convenciones

1. **Componentes**: PascalCase (`UserList.tsx`)
2. **Hooks**: camelCase con prefijo `use` (`useUsers.ts`)
3. **Tipos/Interfaces**: PascalCase, interfaces con prefijo `I` para ports
4. **Archivos**: kebab-case (`api-repository.ts`)
5. **Estilos**: TailwindCSS (clases utility)
6. **Estado**: React hooks (`useState`, `useEffect`, `useContext`)
7. **Fetching**: Repositorios que implementan interfaces del dominio
8. **Páginas**: Archivos `page.tsx` en `app/` (App Router)

## Rutas Principales

| Ruta | Descripción |
|------|-------------|
| `/login` | Inicio de sesión |
| `/home` | Dashboard principal |
| `/orders` | Gestión de pedidos |
| `/products` | Catálogo de productos |
| `/shipments` | Gestión de envíos |
| `/users` | Gestión de usuarios |
| `/roles` | Gestión de roles |
| `/permissions` | Gestión de permisos |
| `/integrations` | Configuración de integraciones |
