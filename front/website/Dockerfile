# Dockerfile en la raíz que apunta al Dockerfile en docker/
# Para compatibilidad con herramientas que esperan Dockerfile en la raíz
FROM --platform=linux/arm64 node:20-alpine AS deps

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci --only=production=false

FROM --platform=linux/arm64 node:20-alpine AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .

ENV NODE_ENV=production
RUN npm run build

FROM --platform=linux/arm64 nginx:alpine AS runner

WORKDIR /app

COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

<<<<<<< HEAD
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1
=======
EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

CMD ["nginx", "-g", "daemon off;"]
