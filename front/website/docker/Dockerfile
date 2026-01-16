# ============================================
# Dockerfile Multi-Stage para Astro
# Arquitectura: ARM64 (compatible con AWS Graviton)
# ============================================

# ============================================
# Stage 1: Dependencies
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS deps

WORKDIR /app

# Copiar archivos de dependencias
COPY package.json package-lock.json* ./

# Instalar dependencias
RUN npm ci --only=production=false

# ============================================
# Stage 2: Builder
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS builder

WORKDIR /app

# Copiar dependencias instaladas desde stage anterior
COPY --from=deps /app/node_modules ./node_modules

# Copiar código fuente
COPY . .

# Variables de entorno para build
ENV NODE_ENV=production

# Build de Astro (genera archivos estáticos en dist/)
RUN npm run build

# ============================================
# Stage 3: Runner con Nginx
# ============================================
FROM --platform=linux/arm64 nginx:alpine AS runner

WORKDIR /app

# Copiar archivos estáticos generados por Astro
COPY --from=builder /app/dist /usr/share/nginx/html

# Copiar configuración personalizada de nginx
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

# Crear usuario no-root para nginx (ya existe en nginx:alpine)
# nginx:alpine ya usa un usuario no-root por defecto

# Exponer puerto (8080 para evitar problemas de permisos en rootless)
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Comando de inicio
CMD ["nginx", "-g", "daemon off;"]
