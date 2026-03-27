.PHONY: help migrate run-backend run-frontend run-mock build test clean docker-up docker-down check-docker run-mobile-web test-mobile

# Variables
BACKEND_DIR = back/central
FRONTEND_DIR = front/central
MIGRATION_DIR = back/migration
TESTING_DIR = back/testing
FRONTEND_TESTING_DIR = front/testing
MOBILE_DIR = mobile/mobile_central
DOCKER_LOCAL = infra/compose-local
DOCKER_PROD = infra/compose-prod

# ======================
# Helpers internos
# ======================

check-docker: ## Verificar e iniciar Docker si es necesario (uso interno)
	@echo "🔍 Verificando servicios Docker..."
	@if ! docker ps | grep -q redis_local; then \
		echo "⚠️  Servicios Docker no están corriendo. Iniciando..."; \
		$(MAKE) docker-up; \
		echo "⏳ Esperando que los servicios se inicialicen (5 segundos)..."; \
		sleep 5; \
		echo "✅ Servicios Docker listos"; \
	else \
		echo "✅ Servicios Docker ya están corriendo"; \
	fi

# ======================
# Ayuda
# ======================

help: ## Mostrar esta ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ======================
# Base de Datos
# ======================

migrate: ## Ejecutar migraciones de base de datos
	@echo "🔄 Ejecutando migraciones..."
	cd $(MIGRATION_DIR) && go run cmd/main.go up

migrate-down: ## Revertir última migración
	@echo "⬇️  Revirtiendo migración..."
	cd $(MIGRATION_DIR) && go run cmd/main.go down

migrate-create: ## Crear nueva migración (uso: make migrate-create NAME=nombre_migracion)
	@echo "📝 Creando migración $(NAME)..."
	cd $(MIGRATION_DIR) && go run cmd/main.go create $(NAME)

# ======================
# Backend
# ======================

run-backend: check-docker ## Iniciar servidor backend
	@echo "🚀 Iniciando backend..."
	cd $(BACKEND_DIR) && go run cmd/main.go

build-backend: ## Compilar backend
	@echo "🔨 Compilando backend..."
	cd $(BACKEND_DIR) && go build -o main cmd/main.go

test-backend: ## Ejecutar tests del backend
	@echo "🧪 Ejecutando tests del backend..."
	cd $(BACKEND_DIR) && go test ./...

swagger: ## Generar documentación Swagger
	@echo "📚 Generando Swagger..."
	cd $(BACKEND_DIR) && swag init -g cmd/main.go

# ======================
# Frontend
# ======================

run-frontend: ## Iniciar servidor frontend (dev)
	@echo "🚀 Iniciando frontend..."
	cd $(FRONTEND_DIR) && pnpm dev

build-frontend: ## Compilar frontend (producción)
	@echo "🔨 Compilando frontend..."
	cd $(FRONTEND_DIR) && pnpm build

install-frontend: ## Instalar dependencias del frontend
	@echo "📦 Instalando dependencias del frontend..."
	cd $(FRONTEND_DIR) && pnpm install

lint-frontend: ## Ejecutar linter del frontend
	@echo "🔍 Ejecutando linter..."
	cd $(FRONTEND_DIR) && pnpm lint

# ======================
# Frontend Testing
# ======================

run-frontend-testing: ## Iniciar frontend de testing (dev, puerto 3051)
	@echo "🚀 Iniciando frontend testing..."
	cd $(FRONTEND_TESTING_DIR) && pnpm dev

build-frontend-testing: ## Compilar frontend de testing (producción)
	@echo "🔨 Compilando frontend testing..."
	cd $(FRONTEND_TESTING_DIR) && pnpm build

install-frontend-testing: ## Instalar dependencias del frontend testing
	@echo "📦 Instalando dependencias del frontend testing..."
	cd $(FRONTEND_TESTING_DIR) && pnpm install

# ======================
# Mobile (Flutter)
# ======================

run-mobile-web: ## Iniciar app mobile en web (apunta a backend local localhost:3050, con credenciales dev)
	@echo "🚀 Iniciando Flutter mobile en web (backend: http://localhost:3050)..."
	cd $(MOBILE_DIR) && flutter run -d chrome \
		--dart-define=APP_ENV=development \
		--dart-define=API_BASE_URL=http://localhost:3050/api/v1 \
		--dart-define=DEV_EMAIL=$$(grep DEV_EMAIL .env.dev 2>/dev/null | cut -d= -f2) \
		--dart-define=DEV_PASSWORD=$$(grep DEV_PASSWORD .env.dev 2>/dev/null | cut -d= -f2)

run-mobile-web-prod: ## Iniciar app mobile en web (apunta a producción)
	@echo "🚀 Iniciando Flutter mobile en web (backend: producción)..."
	cd $(MOBILE_DIR) && flutter run -d chrome

build-mobile-web: ## Compilar app mobile para web
	@echo "🔨 Compilando Flutter mobile para web..."
	cd $(MOBILE_DIR) && flutter build web

test-mobile: ## Ejecutar tests del mobile
	@echo "🧪 Ejecutando tests del mobile..."
	cd $(MOBILE_DIR) && flutter test

install-mobile: ## Instalar dependencias del mobile
	@echo "📦 Instalando dependencias del mobile..."
	cd $(MOBILE_DIR) && flutter pub get

# ======================
# Testing (Simuladores + API)
# ======================

run-testing: ## Iniciar Testing Server (mocks + API en puerto 9092)
	@echo "🎭 Iniciando Testing Server..."
	@echo "  📡 Softpymes HTTP en puerto 9090"
	@echo "  📡 EnvioClick HTTP en puerto 9091"
	@echo "  📡 Testing API en puerto 9092"
	@echo "  🎮 CLI interactivo para Shopify/WhatsApp"
	cd $(TESTING_DIR) && go run cmd/main.go

build-testing: ## Compilar monolito de testing
	@echo "🔨 Compilando Testing Server..."
	cd $(TESTING_DIR) && go build -o testing-server cmd/main.go

# Alias para compatibilidad
run-mock: run-testing ## Alias de run-testing

# ======================
# Docker
# ======================

docker-up: ## Iniciar servicios Docker locales (PostgreSQL, Redis, RabbitMQ, MinIO)
	@echo "🐳 Iniciando servicios Docker..."
	docker-compose -f $(DOCKER_LOCAL)/docker-compose.yaml up -d

docker-down: ## Detener servicios Docker locales
	@echo "🛑 Deteniendo servicios Docker..."
	docker-compose -f $(DOCKER_LOCAL)/docker-compose.yaml down

docker-logs: ## Ver logs de servicios Docker
	docker-compose -f $(DOCKER_LOCAL)/docker-compose.yaml logs -f

docker-ps: ## Ver servicios Docker corriendo
	docker-compose -f $(DOCKER_LOCAL)/docker-compose.yaml ps

# ======================
# Desarrollo
# ======================

dev: docker-up ## Iniciar entorno completo de desarrollo
	@echo "🚀 Entorno de desarrollo listo!"
	@echo ""
	@echo "Servicios disponibles:"
	@echo "  PostgreSQL:  localhost:5433"
	@echo "  Redis:       localhost:6379"
	@echo "  RabbitMQ:    localhost:5672 (UI: http://localhost:15672)"
	@echo "  MinIO:       localhost:9000 (UI: http://localhost:9001)"
	@echo ""
	@echo "Para iniciar backend:          make run-backend"
	@echo "Para iniciar frontend:         make run-frontend"
	@echo "Para iniciar testing backend:  make run-testing"
	@echo "Para iniciar testing frontend: make run-frontend-testing"

run-all: ## Iniciar infra + backend + frontend en tmux
	@./scripts/dev-services.sh start all

stop-all: ## Detener todos los servicios y limpiar tmux
	@./scripts/dev-services.sh stop all

svc-status: ## Ver estado de todos los servicios
	@./scripts/dev-services.sh status

svc-restart: ## Reiniciar un servicio (uso: make svc-restart SVC=backend)
	@./scripts/dev-services.sh restart $(SVC)

svc-logs: ## Ver logs de un servicio (uso: make svc-logs SVC=backend LINES=100)
	@./scripts/dev-services.sh logs $(SVC) $(LINES)

kill-zombies: ## Matar procesos huérfanos (Go, Next.js, puertos)
	@./scripts/dev-services.sh kill-zombies

svc-ports: ## Ver puertos en uso del proyecto
	@./scripts/dev-services.sh ports

# ======================
# Tests
# ======================

test: test-backend ## Ejecutar todos los tests

test-integration: ## Ejecutar tests de integración
	@echo "🧪 Ejecutando tests de integración..."
	cd $(TESTING_DIR) && go test ./...

# ======================
# Limpieza
# ======================

clean: ## Limpiar archivos compilados
	@echo "🧹 Limpiando archivos compilados..."
	find . -name "main" -type f -delete
	find . -name "*.exe" -type f -delete
	find . -name "softpymes-server" -type f -delete
	cd $(FRONTEND_DIR) && rm -rf .next out
	cd $(FRONTEND_TESTING_DIR) && rm -rf .next out

clean-docker: ## Limpiar volúmenes de Docker (⚠️ BORRA DATOS)
	@echo "⚠️  ¿Estás seguro? Esto borrará todos los datos de Docker (PostgreSQL, Redis, etc.)"
	@read -p "Presiona Enter para continuar o Ctrl+C para cancelar..."
	docker-compose -f $(DOCKER_LOCAL)/docker-compose.yaml down -v

# ======================
# Git
# ======================

git-status: ## Ver estado de git
	@git status

git-pull: ## Actualizar desde remoto
	@echo "⬇️  Actualizando desde remoto..."
	@git pull origin main

git-push: ## Subir cambios al remoto
	@echo "⬆️  Subiendo cambios..."
	@git push origin $(shell git branch --show-current)

# ======================
# Producción
# ======================

deploy-prod: ## Desplegar en producción (⚠️ usar con cuidado)
	@echo "🚀 Desplegando en producción..."
	@echo "⚠️  Este comando debe ejecutarse en el servidor de producción"
	cd $(DOCKER_PROD) && sudo podman-compose down
	cd $(DOCKER_PROD) && sudo podman-compose up -d

# ======================
# Información
# ======================

info: ## Mostrar información del proyecto
	@echo "📊 Información del Proyecto Probability"
	@echo ""
	@echo "Backend:"
	@echo "  Directorio: $(BACKEND_DIR)"
	@echo "  Framework:  Gin + GORM"
	@echo "  Puerto:     8080 (default)"
	@echo ""
	@echo "Frontend:"
	@echo "  Directorio: $(FRONTEND_DIR)"
	@echo "  Framework:  Next.js 16.1 + React 19"
	@echo "  Puerto:     3000 (default)"
	@echo ""
	@echo "Frontend Testing:"
	@echo "  Directorio: $(FRONTEND_TESTING_DIR)"
	@echo "  Framework:  Next.js 16.1 + React 19"
	@echo "  Puerto:     3051"
	@echo ""
	@echo "Testing Server (Simuladores + API):"
	@echo "  Directorio: $(TESTING_DIR)"
	@echo "  Simuladores: Softpymes, EnvioClick, Shopify, WhatsApp"
	@echo "  Puertos:    9090 (Softpymes), 9091 (EnvioClick), 9092 (Testing API)"
	@echo ""
	@echo "Base de Datos:"
	@echo "  Producción: database-1.capmmoe4cw2e.us-east-1.rds.amazonaws.com:5432"
	@echo "  Local:      localhost:5433"
	@echo ""
	@echo "Variables de entorno importantes:"
	@echo "  SOFTPYMES_API_URL       - URL de la API de Softpymes (testing: http://localhost:9090)"
	@echo "  SOFTPYMES_MOCK_PORT     - Puerto del simulador Softpymes (default: 9090)"

