# 🚀 Migración de Docker a Podman

Esta guía explica cómo migrar de Docker a Podman para el frontend de Probability Central.

## 📁 Estructura de Archivos

Toda la configuración de Podman está en la carpeta `podman/`:
- `deploy-podman.sh` - Script de despliegue a producción
- `PODMAN_MIGRATION.md` - Esta guía

**Nota**: Podman usa el mismo `docker/Dockerfile` que Docker, ya que Podman es completamente compatible con Dockerfiles.

## 📋 ¿Qué es Podman?

Podman es una alternativa a Docker que:
- ✅ **No requiere daemon**: Ejecuta contenedores sin privilegios root
- ✅ **Compatible con Docker**: Usa los mismos Dockerfiles y docker-compose.yml
- ✅ **Más seguro**: Ejecuta contenedores como usuario no-root por defecto
- ✅ **Mismo rendimiento**: Similar a Docker en velocidad y uso de recursos

## 🔧 Instalación

### Linux (Ubuntu/Debian)
```bash
sudo apt-get update
sudo apt-get install -y podman
```

### macOS
```bash
brew install podman
```

### Verificar instalación
```bash
podman --version
```

## 🔄 Diferencias Clave

### 1. Comandos Básicos

| Docker | Podman |
|--------|--------|
| `docker build` | `podman build` |
| `docker run` | `podman run` |
| `docker ps` | `podman ps` |

### 2. Build Multi-Arquitectura

**Docker (requiere buildx):**
```bash
docker buildx build --platform linux/arm64 ...
```

**Podman (soporte nativo):**
```bash
podman build --platform linux/arm64 ...
```

### 3. Dockerfile

✅ **Podman es completamente compatible con Dockerfiles**

En este proyecto usamos el mismo `docker/Dockerfile` para ambos:
- `docker/Dockerfile` - Usado tanto por Docker como por Podman

## 📝 Uso

### Desarrollo Local

Build y ejecución manual:
```bash
# Desde front/central/
cd front/central

# Build de la imagen
podman build --platform linux/arm64 -f docker/Dockerfile -t probability-front-central:latest .

# Ejecutar contenedor
podman run -d \
    --name probability-frontend \
    -p 3000:80 \
    -e NODE_ENV=production \
    -e NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1 \
    probability-front-central:latest
```

### Despliegue a Producción

Usa el script de deploy para Podman (desde `front/central/`):
```bash
./podman/deploy-podman.sh [version]
```

Ejemplo:
```bash
cd front/central
./podman/deploy-podman.sh latest
./podman/deploy-podman.sh v1.0.0
```

## 🔐 Configuración de Redes

Podman crea redes de forma similar a Docker:

```bash
# Crear red (si no existe)
podman network create app-network

# Ver redes
podman network ls

# Inspeccionar red
podman network inspect app-network
```

## 🐛 Troubleshooting

### Error: "cannot find runtime"
```bash
# Inicializar Podman (solo primera vez)
podman machine init
podman machine start
```

### Error: "permission denied"
```bash
# Podman puede ejecutarse sin root, pero si necesitas privilegios:
sudo podman run ...
```

### Error: "network not found"
```bash
# Crear la red manualmente
podman network create probability-network
```


## 📊 Comparación de Comandos

### Build
```bash
# Docker (desde front/central/)
docker buildx build --platform linux/arm64 -f docker/Dockerfile -t image:tag .

# Podman (desde front/central/) - Usa el mismo Dockerfile
podman build --platform linux/arm64 -f docker/Dockerfile -t image:tag .
```

### Run
```bash
# Docker
docker run -d --name container -p 8080:80 image:tag

# Podman
podman run -d --name container -p 8080:80 image:tag
```


## ✅ Checklist de Migración

- [x] Carpeta `podman/` creada con toda la configuración
- [x] Script de deploy para Podman creado (`podman/deploy-podman.sh`)
- [x] Configurado para usar el mismo `docker/Dockerfile` que Docker
- [ ] Instalar Podman en servidor de producción
- [ ] Probar build local con Podman: `podman build -f docker/Dockerfile -t test:latest .`
- [ ] Probar despliegue con `./podman/deploy-podman.sh latest`
- [ ] Actualizar documentación de CI/CD si aplica

## 🔗 Recursos

- [Documentación oficial de Podman](https://podman.io/getting-started/)
- [Podman vs Docker](https://podman.io/what-is-podman/)

## 💡 Notas Importantes

1. **Mismo Dockerfile**: Podman usa el mismo `docker/Dockerfile` que Docker, ya que es completamente compatible.
2. **Sin daemon**: Podman no requiere un daemon corriendo, lo que lo hace más ligero.
3. **Rootless por defecto**: Los contenedores se ejecutan como usuario no-root, mejorando la seguridad.
4. **Mismo rendimiento**: Podman tiene un rendimiento similar a Docker.
5. **Estructura organizada**: Toda la configuración de Podman está en `podman/` para mantener el proyecto organizado.

## 🎯 Próximos Pasos

1. Instalar Podman en tu máquina local
2. Probar el build (desde `front/central/`): 
   ```bash
   cd front/central
   podman build -f docker/Dockerfile -t test:latest .
   ```
3. Probar el deploy: 
   ```bash
   cd front/central
   ./podman/deploy-podman.sh latest
   ```
4. Actualizar scripts de CI/CD si es necesario
