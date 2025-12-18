# 🐳 Docker - Probability Website Frontend

Documentación para construir y desplegar la imagen Docker del frontend Website (Astro) para ARM64.

## 📋 Requisitos Previos

- **Docker** 20.10 o superior con BuildKit habilitado
- **Docker Buildx** para builds multi-arquitectura
- **AWS CLI** configurado con credenciales válidas

## 🏗️ Arquitectura

La imagen está optimizada para **ARM64 (AWS Graviton)** y utiliza:
- **Base**: Node.js 20 Alpine (build) + Nginx Alpine (runtime)
- **Multi-stage build**: Reduce el tamaño final de la imagen
- **Static Site**: Astro genera archivos estáticos optimizados
- **Nginx**: Servidor web ligero para servir archivos estáticos
- **Non-root user**: Nginx Alpine ya ejecuta como usuario no-root

### 🌐 Arquitectura de Red

```
┌─────────────────────────────────────────────────────────────┐
│                    SERVIDOR PRODUCCIÓN                       │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │         Red Interna Docker: probability-network        │ │
│  │                                                         │ │
│  │  ┌──────────────────────┐                               │ │
│  │  │   Website           │                               │ │
│  │  │   (Astro + Nginx)   │                               │ │
│  │  │   Interno: 80       │                               │ │
│  │  │   Host: 8080        │                               │ │
│  │  │   (8080:80)         │                               │ │
│  │  └──────────────────────┘                               │ │
│  │         │                                                │ │
│  └─────────│────────────────────────────────────────────────┘ │
│            │                                                    │
│            │ HTTP/HTTPS                                         │
│            ▼                                                    │
│   https://probabilityia.com.co                                 │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

## 🚀 Despliegue a Producción

### Desplegar a ECR Público

```bash
# Desde el directorio raíz del proyecto (front/website)
./script/deploy.sh
```

O con una versión específica:

```bash
./script/deploy.sh v1.0.0
```

Este script:
1. ✅ Verifica dependencias (Docker, AWS CLI, Buildx)
2. 📦 Instala dependencias de Node.js
3. 🔨 Construye la imagen para ARM64
4. 🏷️ Crea tags descriptivos (website-latest, website-TIMESTAMP)
5. 🔐 Hace login a ECR público
6. ⬆️ Sube la imagen a ECR

## 📦 Usar la Imagen desde ECR

### Pull de la Imagen

```bash
# Login a ECR público
aws ecr-public get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin public.ecr.aws

# Pull de la imagen
docker pull public.ecr.aws/c1l9h7c9/probability:website-latest
```

### Ejecutar en Servidor ARM64

```bash
# Ejecución básica
docker run -d \
  --name probability-website \
  --restart unless-stopped \
  -p 8080:80 \
  public.ecr.aws/c1l9h7c9/probability:website-latest
```

**NOTAS:**
- Puerto interno: `80` (Nginx escucha en puerto 80)
- Puerto expuesto: `8080` (acceso desde el host)
- La imagen incluye todos los archivos estáticos generados por Astro
- Nginx sirve los archivos con compresión gzip y cache optimizado

### Ejecutar con Docker Compose

```bash
# Desde el directorio raíz del proyecto
docker-compose up -d
```

## 📊 Métricas de la Imagen

- **Tamaño final**: ~30-50 MB (comprimido)
- **Arquitectura**: linux/arm64
- **Base image**: nginx:alpine (runtime)
- **Usuario**: nginx (non-root, ya incluido en nginx:alpine)

## 🔍 Troubleshooting

### Build Falla en Simulación ARM64

Si el build de ARM64 falla en un sistema x86/amd64:

```bash
# Verificar que buildx esté instalado
docker buildx version

# Crear nuevo builder
docker buildx create --name multiarch-builder --driver docker-container --use

# Listar plataformas disponibles
docker buildx inspect --bootstrap
```

### Imagen No Inicia

Ver logs del contenedor:
```bash
docker logs -f probability-website
```

Entrar al contenedor:
```bash
docker exec -it probability-website sh
```

Verificar que los archivos estén presentes:
```bash
docker exec -it probability-website ls -la /usr/share/nginx/html
```

### Nginx No Sirve Archivos

Verificar configuración de nginx:
```bash
docker exec -it probability-website cat /etc/nginx/conf.d/default.conf
```

Probar nginx:
```bash
docker exec -it probability-website nginx -t
```

## 🏷️ Tags Disponibles en ECR

- `website-latest`: Última versión estable
- `website-YYYYMMDD-HHMMSS`: Versión con timestamp
- `website-vX.Y.Z`: Versiones específicas

Ver todos los tags:
```
https://gallery.ecr.aws/c1l9h7c9/probability
```

## 📝 Notas Importantes

1. **Static Site**: Astro genera un sitio estático, no necesita Node.js en runtime
2. **Multi-Stage Build**: Reduce el tamaño final eliminando dependencias de desarrollo
3. **ARM64 Native**: La imagen está compilada nativamente para ARM64 (AWS Graviton)
4. **Security**: Nginx Alpine ejecuta como usuario no-root por defecto
5. **Cache**: Docker usa caché de capas para builds más rápidos
6. **Gzip**: Nginx comprime automáticamente las respuestas
7. **Healthcheck**: Incluido para monitoreo de salud del contenedor

## 🔗 Enlaces Útiles

- [Astro Docker Deployment](https://docs.astro.build/en/guides/deploy/docker/)
- [Docker Buildx Multi-platform](https://docs.docker.com/build/building/multi-platform/)
- [AWS ECR Public Gallery](https://gallery.ecr.aws/c1l9h7c9/probability)
- [AWS Graviton](https://aws.amazon.com/ec2/graviton/)
- [Nginx Alpine](https://hub.docker.com/_/nginx)

## 📞 Soporte

Para problemas con el despliegue, contacta al equipo de DevOps.
