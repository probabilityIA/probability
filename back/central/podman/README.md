# 🐳 Podman - Probability Backend Central

Configuración para usar Podman en lugar de Docker.

## 📁 Estructura

```
podman/
+-- deploy-podman.sh       # Script de despliegue a producción (ECR)
+-- README.md              # Este archivo

Nota: Podman usa el mismo docker/Dockerfile que Docker
```

## 🚀 Uso

### Despliegue a Producción

```bash
# Desde back/central/
cd back/central
./podman/deploy-podman.sh latest
```

### Build Manual

```bash
# Desde back/central/
cd back/central

# Build de la imagen (contexto es el directorio padre para incluir migration)
podman build --platform linux/arm64 -f docker/Dockerfile -t probability-back-central:latest ..
```

### Ejecutar Localmente

```bash
podman run -d \
    --name probability-back-central \
    --env-file .env \
    -p 8080:8080 \
    probability-back-central:latest
```

## 📝 Notas

- El script configura automáticamente la emulación QEMU si estás en x86_64
- Usa el mismo `docker/Dockerfile` que Docker
- El contexto de build es el directorio padre (`..`) para incluir el módulo `migration`

## 🔄 Diferencias con Docker

| Aspecto | Docker | Podman |
|---------|--------|--------|
| Daemon | Requerido | No requerido |
| Root | Requerido por defecto | Rootless por defecto |
| Build multi-arch | Requiere buildx | Soporte nativo + QEMU |
| Archivo de build | Dockerfile | Dockerfile (mismo archivo) |

## ✅ Ventajas de Podman

1. **Sin daemon**: Más ligero, no requiere servicio corriendo
2. **Rootless**: Ejecuta contenedores sin privilegios root
3. **Compatible**: Usa los mismos Dockerfiles que Docker
4. **Seguro**: Mejor aislamiento y seguridad por defecto
