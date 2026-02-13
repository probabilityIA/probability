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
# Stage 3: Runner con Node.js (SSR)
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS runner

WORKDIR /app

# Copiar solo las dependencias necesarias en producción
COPY --from=builder /app/node_modules ./node_modules

# Copiar la aplicación compilada (server + client)
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./package.json

# Variables de entorno
ENV NODE_ENV=production
ENV HOST=0.0.0.0
ENV PORT=80

# Crear usuario no-root
RUN addgroup -g 1001 -S nodejs && adduser -S nodejs -u 1001

# Cambiar propiedad de la carpeta
RUN chown -R nodejs:nodejs /app

USER nodejs

# Exponer puerto
EXPOSE 80

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "require('http').get('http://localhost:80/website/', (r) => {if (r.statusCode !== 200) throw new Error(r.statusCode)})" || exit 1

# Comando de inicio - Ejecutar el servidor Astro SSR
CMD ["node", "./dist/server/entry.mjs"]
