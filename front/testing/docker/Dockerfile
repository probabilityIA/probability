# ============================================
# Dockerfile Multi-Stage para Testing Frontend (Next.js)
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

# Build args para URLs de API (baked at build time para rewrites y basePath)
ARG CENTRAL_API_URL
ARG TESTING_API_URL
ARG BASE_PATH=/testing
ENV CENTRAL_API_URL=${CENTRAL_API_URL}
ENV TESTING_API_URL=${TESTING_API_URL}
ENV BASE_PATH=${BASE_PATH}

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

# Copiar archivos de build
# Next.js genera archivos standalone automáticamente
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Variables de entorno
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3051
ENV HOSTNAME="0.0.0.0"

# Copiar script de startup con verificación de backend (antes de USER)
COPY docker/startup.sh /app/startup.sh
RUN chmod +x /app/startup.sh && chown nextjs:nodejs /app/startup.sh

# Cambiar a usuario no-root
USER nextjs

# Exponer puerto
EXPOSE 3051

# Comando de inicio con verificación de backend
CMD ["/app/startup.sh"]
