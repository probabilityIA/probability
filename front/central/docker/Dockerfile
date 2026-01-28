# ============================================
# Dockerfile Multi-Stage para Next.js 15
# Arquitectura: ARM64 (compatible con AWS Graviton)
# ============================================

# ============================================
# Stage 1: Dependencies
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS deps

WORKDIR /app

# Copiar archivos de dependencias
COPY package.json package-lock.json* pnpm-lock.yaml* ./

# Instalar dependencias basado en el lock file disponible
# Los overrides en package.json forzarán styled-jsx@5.1.7 (compatible con React 19.2.3)
# --legacy-peer-deps permite resolver conflictos de peer dependencies (styled-jsx con React 19)
RUN \
  if [ -f pnpm-lock.yaml ]; then \
    corepack enable pnpm && pnpm install --frozen-lockfile; \
  elif [ -f package-lock.json ]; then \
    npm ci --legacy-peer-deps; \
  else \
    npm install --legacy-peer-deps; \
  fi

# ============================================
# Stage 2: Builder
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS builder

WORKDIR /app

# Copiar dependencias instaladas desde stage anterior
COPY --from=deps /app/node_modules ./node_modules

# Copiar código fuente
COPY . .

# Build args para variables de Next.js
# NEXT_PUBLIC_* = Cliente (SSE, dominio público)
# API_BASE_URL = Servidor (Server Actions, red interna Docker)
ARG NEXT_PUBLIC_API_BASE_URL
ARG API_BASE_URL
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}
ENV API_BASE_URL=${API_BASE_URL}

# Variables de entorno para build
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# Limpiar cache de Next.js antes del build
RUN rm -rf .next

# Build de Next.js (limpio, sin cache)
RUN \
  if [ -f pnpm-lock.yaml ]; then \
    corepack enable pnpm && pnpm run build; \
  else \
    npm run build; \
  fi

# ============================================
# Stage 3: Runner (Imagen final)
# ============================================
FROM --platform=linux/arm64 node:20-alpine AS runner

WORKDIR /app

# Crear usuario no-root para seguridad
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copiar archivos públicos
COPY --from=builder /app/public ./public

# Copiar archivos de build
# Next.js genera archivos standalone automáticamente
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Variables de entorno
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# IMPORTANTE: API_BASE_URL debe pasarse en runtime (no en build time)
# Se configura en docker-compose.yml o docker run -e API_BASE_URL=...

# Cambiar a usuario no-root
USER nextjs

# Exponer puerto (3000 para Podman rootless)
EXPOSE 3000

# Comando de inicio
CMD ["node", "server.js"]

