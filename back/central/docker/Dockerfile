# Etapa 1: Build con Go
FROM --platform=linux/arm64 golang:1.23-alpine AS builder

# Variables de compilación
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# Crear usuario no root
RUN adduser -D -g '' appuser

# Usaremos un workspace que contenga ambos módulos: central y migration
WORKDIR /workspace/central

# Copiar go.mod y go.sum del módulo central
COPY central/go.mod central/go.sum ./

# Copiar el módulo migration en la ruta esperada por el replace ../migration
COPY migration /workspace/migration

# Descargar dependencias (ahora que ../migration existe)
RUN go mod download

# Copiar todo el código fuente del módulo central
COPY central .

# Compilar el servicio central
RUN go build -ldflags="-w -s" -o probability-back-central ./cmd/main.go

# Etapa 2: Imagen final minimalista
FROM --platform=linux/arm64 alpine:3.19

# Variables de entorno
ENV TZ=America/Bogota

# Crear usuario no root
RUN adduser -D -g '' appuser

WORKDIR /app

# Copiar binario desde build stage
COPY --from=builder /workspace/central/probability-back-central .

# Instalar certificados SSL y configurar zona horaria
RUN apk update && apk add --no-cache \
    ca-certificates \
    tzdata && \
    cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone && \
    rm -rf /var/cache/apk/*

# Cambiar a usuario no root
USER appuser

# Exponer puerto
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Ejecutar la app
CMD ["./probability-back-central"]
