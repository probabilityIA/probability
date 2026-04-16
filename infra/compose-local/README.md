# Docker Compose para Desarrollo Local

Este docker-compose levanta los servicios necesarios para desarrollo local:

- **PostgreSQL** (puerto 5432)
- **Redis** (puerto 6379)
- **MinIO** (puertos 9000 API, 9001 Console)
- **RabbitMQ** (puertos 5672 AMQP, 15672 Management UI)

## Uso

### Levantar todos los servicios

```bash
cd infra/compose-local
docker-compose up -d
```

### Ver logs

```bash
docker-compose logs -f
```

### Detener servicios

```bash
docker-compose down
```

### Detener y eliminar volúmenes (⚠️ borra datos)

```bash
docker-compose down -v
```

## Credenciales por defecto

### PostgreSQL
- **Host:** localhost
- **Puerto:** 5432
- **Database:** probability
- **Usuario:** postgres
- **Contraseña:** postgres

### Redis
- **Host:** localhost
- **Puerto:** 6379
- **Sin contraseña** (solo local)

### MinIO
- **API Endpoint:** http://localhost:9000
- **Console:** http://localhost:9001
- **Access Key:** minioadmin
- **Secret Key:** minioadmin

**Primera vez:**
1. Abre http://localhost:9001
2. Login con `minioadmin` / `minioadmin`
3. Crea el bucket `probability`
4. Configura el bucket como público (Access Policy → Public)

### RabbitMQ
- **AMQP:** localhost:5672
- **Management UI:** http://localhost:15672
- **Usuario:** admin
- **Contraseña:** admin

## Variables de entorno para tu aplicación

Agrega estas variables a tu `.env`:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=probability
DB_USER=postgres
DB_PASS=postgres
DB_LOG_LEVEL=info
PGSSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# MinIO (S3)
S3_ENDPOINT=http://localhost:9000
S3_BUCKET=probability
S3_REGION=us-east-1
S3_KEY=minioadmin
S3_SECRET=minioadmin
URL_BASE_DOMAIN_S3=http://localhost:9000/probability

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=admin
```

## Verificar que todo funciona

```bash
# PostgreSQL
docker exec -it postgres_local psql -U postgres -d probability -c "SELECT version();"

# Redis
docker exec -it redis_local redis-cli ping

# MinIO
curl http://localhost:9000/minio/health/live

# RabbitMQ
curl -u admin:admin http://localhost:15672/api/overview
```
