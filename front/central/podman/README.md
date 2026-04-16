# 🐳 Podman - Probability Frontend Central

Configuración completa para usar Podman en lugar de Docker.

## 📁 Estructura

```
podman/
├── deploy-podman.sh       # Script de despliegue a producción (ECR)
├── PODMAN_MIGRATION.md    # Guía completa de migración
└── README.md              # Este archivo

Nota: Podman usa el mismo docker/Dockerfile que Docker
```

## 🚀 Inicio Rápido

### Desarrollo Local

```bash
# Desde front/central/
cd front/central

# Build y run manual
podman build --platform linux/arm64 -f docker/Dockerfile -t probability-front-central:latest .
podman run -d --name probability-frontend -p 3000:80 probability-front-central:latest
```

### Despliegue a Producción

```bash
# Desde front/central/
cd front/central
./podman/deploy-podman.sh latest
```

### Build Manual

```bash
# Desde front/central/
cd front/central
podman build --platform linux/arm64 -f docker/Dockerfile -t probability-front-central:latest .
```

## 📝 Archivos

### Dockerfile
Podman usa el mismo `docker/Dockerfile` que Docker, ya que es completamente compatible. No necesitamos un archivo separado.

### deploy-podman.sh
Script de despliegue que:
- Construye la imagen para ARM64
- La etiqueta para ECR
- La sube a AWS ECR público

## 🔄 Diferencias con Docker

| Aspecto | Docker | Podman |
|---------|--------|--------|
| Daemon | Requerido | No requerido |
| Root | Requerido por defecto | Rootless por defecto |
| Build multi-arch | Requiere buildx | Soporte nativo |
| Archivo de build | Dockerfile | Dockerfile (mismo archivo) |

## 📚 Documentación

Para más detalles, consulta:
- `PODMAN_MIGRATION.md` - Guía completa de migración
- [Documentación oficial de Podman](https://podman.io/getting-started/)

## ✅ Ventajas de Podman

1. **Sin daemon**: Más ligero, no requiere servicio corriendo
2. **Rootless**: Ejecuta contenedores sin privilegios root
3. **Compatible**: Usa los mismos formatos que Docker
4. **Seguro**: Mejor aislamiento y seguridad por defecto
