# Probability

Plataforma de gestión de e-commerce multi-tenant que permite centralizar y administrar pedidos, productos, pagos y envíos desde múltiples canales de venta.

## 📚 Documentación

- **[Deployment System](docs/deploy/README.md)** - Sistema de deployment automático CI/CD
  - [Panic/Restart Mechanism](docs/deploy/panic-restart-mechanism.md) - Auto-recovery de contenedores
  - [Workflow Structure](docs/deploy/workflow-structure.md) - Estructura de workflows GitHub Actions
  - [Troubleshooting](docs/deploy/troubleshooting.md) - Guía de resolución de problemas

## Stack Tecnológico

### Backend
- **Go 1.23** con Gin Framework
- **PostgreSQL 15** (GORM)
- **Redis 7** (Cache)
- **RabbitMQ 3** (Cola de mensajes)
- **MinIO/S3** (Almacenamiento)
- **DynamoDB** (NoSQL)
- **Swagger** (Documentación API)

### Frontend
- **Dashboard**: Next.js 16.1, React 19, TailwindCSS 4
- **Website**: Astro 5, Preact, TailwindCSS 4

## Estructura del Proyecto

```
probability/
├── back/
│   ├── central/           # API principal
│   │   ├── cmd/           # Entry point
│   │   ├── services/      # Servicios de negocio
│   │   │   ├── auth/      # Autenticación (users, roles, permissions)
│   │   │   ├── modules/   # Módulos (orders, products, payments, shipments)
│   │   │   └── integrations/  # Integraciones (Shopify, Amazon, Meli, WhatsApp)
│   │   └── shared/        # Utilidades compartidas
│   ├── migration/         # Migraciones de base de datos
│   ├── notify-email/      # Servicio de notificaciones
│   └── integrationTest/   # Tests de integración
├── front/
│   ├── central/           # Dashboard admin (Next.js)
│   └── website/           # Landing page (Astro)
├── infra/
│   ├── compose-local/     # Docker Compose desarrollo
│   ├── compose-prod/      # Docker Compose producción
│   └── nginx/             # Configuración Nginx
└── scripts/               # Scripts de utilidad
```

## Requisitos

- Go 1.23+
- Node.js 20+
- pnpm
- Docker & Docker Compose

## Instalación

### 1. Iniciar servicios de infraestructura

```bash
docker-compose -f infra/compose-local/docker-compose.yaml up -d
```

### 2. Configurar variables de entorno

```bash
cp back/central/.env.example back/central/.env
# Editar .env con las configuraciones necesarias
```

### 3. Ejecutar migraciones

```bash
cd back/migration
go run main.go
```

### 4. Iniciar Backend

```bash
cd back/central
go run cmd/main.go
```

### 5. Iniciar Frontend Dashboard

```bash
cd front/central
pnpm install
pnpm dev
```

### 6. Iniciar Website

```bash
cd front/website
pnpm install
pnpm dev
```

## Puertos de Desarrollo

| Servicio       | Puerto | Descripción            |
|----------------|--------|------------------------|
| Backend API    | 8080   | API REST               |
| Frontend       | 3000   | Dashboard Next.js      |
| Website        | 4321   | Landing Astro          |
| PostgreSQL     | 5433   | Base de datos          |
| Redis          | 6379   | Cache                  |
| RabbitMQ       | 5672   | Cola de mensajes       |
| RabbitMQ UI    | 15672  | Consola de gestión     |
| MinIO API      | 9000   | Almacenamiento         |
| MinIO UI       | 9001   | Consola MinIO          |

## Arquitectura

El backend sigue el patrón de **Arquitectura Hexagonal**:

```
service/
├── internal/
│   ├── domain/        # Entidades y puertos (interfaces)
│   ├── app/           # Casos de uso
│   └── infra/
│       ├── primary/   # Handlers HTTP, consumidores
│       └── secondary/ # Implementaciones (DB, servicios externos)
```

## Integraciones

- **Shopify**: Sincronización de productos y webhooks
- **Amazon**: Marketplace
- **MercadoLibre**: Marketplace LATAM
- **WhatsApp**: Notificaciones a clientes

## Scripts Útiles

```bash
# Verificar cola RabbitMQ
./scripts/check_rabbitmq_queue.sh

# Purgar cola RabbitMQ
./scripts/purge_rabbitmq_queue.sh
```

## Documentación API

La documentación Swagger está disponible en:
```
http://localhost:8080/swagger/index.html
```

Para regenerar la documentación:
```bash
cd back/central
swag init -g cmd/main.go
```
