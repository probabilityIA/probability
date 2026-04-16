# Probability Central - Mobile App

App movil de Probability Central construida con Flutter. Replica la funcionalidad del dashboard web (Next.js) para dispositivos Android.

## Requisitos

- Flutter SDK 3.10+
- Android SDK (con platform-tools para ADB)
- Dispositivo Android con Depuracion USB o inalambrica habilitada

## Comandos Make (desde la raiz del proyecto)

```bash
# Correr en navegador web (backend local)
make run-mobile-web

# Correr en navegador web (backend produccion)
make run-mobile-web-prod

# Compilar APK de produccion
make build-mobile-web

# Ejecutar tests
make test-mobile

# Instalar dependencias
make install-mobile
```

## Configuracion de desarrollo

### Credenciales de desarrollo

Crear archivo `.env.dev` en `/mobile/mobile_central/` (ignorado por git):

```
DEV_EMAIL=tu_email@ejemplo.com
DEV_PASSWORD=tu_contraseña
```

Estas credenciales pre-llenan el formulario de login solo cuando se corre con `make run-mobile-web` o con `--dart-define`. Nunca se incluyen en builds de produccion.

### Variables de entorno (--dart-define)

| Variable | Descripcion | Default |
|----------|-------------|---------|
| `APP_ENV` | Entorno (development, staging, production) | production |
| `API_BASE_URL` | URL base de la API | segun APP_ENV |
| `DEV_EMAIL` | Email pre-llenado en login | (vacio) |
| `DEV_PASSWORD` | Password pre-llenado en login | (vacio) |

### URLs de API por entorno

| Entorno | URL |
|---------|-----|
| development | `http://10.0.2.2:3050/api/v1` (emulador) o via `API_BASE_URL` |
| staging | `https://staging.probabilityia.com.co/api/v1` |
| production | `https://www.probabilityia.com.co/api/v1` |

## Probar en telefono Android por WiFi

### 1. Preparar el telefono

1. Ir a **Ajustes > Acerca del telefono** y tocar 7 veces en "Numero de compilacion" para activar opciones de desarrollador
2. Ir a **Ajustes > Opciones de desarrollador**
3. Activar **Depuracion inalambrica**

### 2. Vincular el telefono (pairing)

1. En **Depuracion inalambrica**, tocar **Vincular dispositivo con codigo de vinculacion**
2. Anotar el **IP:PUERTO_PAIRING** y el **CODIGO** que aparecen
3. En la terminal:

```bash
# Vincular (usar el puerto y codigo del dialogo de vinculacion)
~/Android/Sdk/platform-tools/adb pair IP:PUERTO_PAIRING CODIGO
# Ejemplo: adb pair 192.168.3.248:40121 393184
```

### 3. Conectar

1. En la pantalla principal de **Depuracion inalambrica** (no el dialogo), buscar **Direccion IP y puerto**
2. Usar ese puerto (diferente al de vinculacion):

```bash
# Conectar (usar el puerto de la pantalla principal)
~/Android/Sdk/platform-tools/adb connect IP:PUERTO_CONEXION
# Ejemplo: adb connect 192.168.3.248:40003
```

3. Verificar conexion:

```bash
~/Android/Sdk/platform-tools/adb devices
# Debe mostrar: 192.168.3.248:40003  device
```

### 4. Correr la app en el telefono

```bash
# Debug mode (con logs en consola) - apunta a produccion con credenciales dev
cd mobile/mobile_central
flutter run -d 192.168.3.248:40003 \
  --dart-define=DEV_EMAIL=$(grep DEV_EMAIL .env.dev | cut -d= -f2) \
  --dart-define=DEV_PASSWORD=$(grep DEV_PASSWORD .env.dev | cut -d= -f2)

# Debug mode apuntando a backend local
flutter run -d 192.168.3.248:40003 \
  --dart-define=APP_ENV=development \
  --dart-define=API_BASE_URL=http://IP_DE_TU_PC:3050/api/v1 \
  --dart-define=DEV_EMAIL=$(grep DEV_EMAIL .env.dev | cut -d= -f2) \
  --dart-define=DEV_PASSWORD=$(grep DEV_PASSWORD .env.dev | cut -d= -f2)
```

### 5. Instalar APK directamente

```bash
# Generar APK de produccion
flutter build apk --release

# Instalar en telefono conectado
~/Android/Sdk/platform-tools/adb install -r build/app/outputs/flutter-apk/app-release.apk
```

### Notas importantes

- El **puerto de vinculacion** (pairing) y el **puerto de conexion** son diferentes
- El codigo de vinculacion expira rapido, usalo inmediatamente
- PC y telefono deben estar en la **misma red WiFi**
- Si el ping al telefono falla (`ping IP_TELEFONO`), verificar que el WiFi no tenga aislamiento de clientes
- Usar `r` para hot reload y `R` para hot restart en la consola de Flutter

## Estructura del proyecto

```
lib/
├── core/                    # Configuracion, red, router, errores
│   ├── config/              # Environment (API URLs, dev credentials)
│   ├── errors/              # Parser de errores (español)
│   ├── network/             # ApiClient (Dio + X-Client-Type)
│   ├── router/              # GoRouter con ShellRoute + Drawer
│   └── storage/             # TokenStorage (FlutterSecureStorage)
├── services/                # Modulos de negocio (arquitectura hexagonal)
│   ├── auth/                # Login, users, roles, permissions, business, resources, actions
│   ├── integrations/        # Core + 26 plataformas (ecommerce, invoicing, pay, transport, messages)
│   └── modules/             # Orders, products, customers, invoicing, dashboard, etc.
├── shared/
│   ├── types/               # PaginatedResponse, Pagination
│   ├── widgets/             # AppShell (drawer), BusinessSelectorWrapper, NetworkAvatar
│   │   └── modules/         # Module screens con TabBar (orders, inventory, delivery, etc.)
│   └── theme/               # AppTheme
└── main.dart                # Entry point con MultiProvider
```

### Arquitectura por modulo

Cada modulo sigue arquitectura hexagonal:

```
modulo/
├── domain/
│   ├── entities.dart        # Entidades (fromJson/toJson)
│   └── ports.dart           # Interfaz del repositorio
├── app/
│   └── use_cases.dart       # Casos de uso (delegacion al repo)
├── infra/
│   └── repository/          # Implementacion HTTP (ApiClient)
└── ui/
    ├── providers/            # ChangeNotifier (estado)
    └── screens/              # Pantallas Flutter
```

## Tests

2126 tests unitarios cubriendo entities, use cases y providers:

```bash
# Todos los tests
cd mobile/mobile_central && flutter test

# Tests de un modulo especifico
flutter test test/services/auth/
flutter test test/services/modules/
flutter test test/services/integrations/
```
