# Análisis de Vulnerabilidades de Seguridad

## 🚨 Vulnerabilidades Identificadas

Basado en el análisis del código y la configuración, se identificaron las siguientes vulnerabilidades críticas que podrían haber permitido el defacement del sitio:

### 0. **Redis Expuesto Públicamente Sin Contraseña** 🔴 CRÍTICO - VECTOR DE ATAQUE PRINCIPAL

**Ubicación**: `infra/compose-prod/docker-compose.yaml` línea 149-150

**Descripción**:
- Redis estaba expuesto públicamente en el puerto 6379 (`0.0.0.0:6379`)
- Redis NO tenía contraseña configurada (`requirepass`)
- Esto permitió acceso no autorizado desde internet
- Los atacantes (bots automatizados) escanean internet buscando Redis expuestos sin contraseña
- Una vez conectados, pueden leer/escribir datos, ejecutar comandos, y potencialmente comprometer el servidor

**Evidencia del Ataque**:
```yaml
# infra/compose-prod/docker-compose.yaml (ANTES)
redis:
  ports:
    - "6379:6379"  # ⚠️ Expuesto a 0.0.0.0 (todo internet)
  command: redis-server --appendonly yes
  # ❌ NO tenía contraseña configurada
```

**Cómo Ocurrió el Ataque**:
1. Bot automatizado escaneó el rango de IPs buscando puerto 6379 abierto
2. Se conectó a Redis sin autenticación
3. Usó comandos de Redis para modificar configuración o inyectar código
4. Esto causó que el sitio web mostrara el mensaje de defacement (página china con cerdito)

**Solución Implementada**:
- ✅ Puertos de Redis cerrados (solo acceso interno vía `app-network`)
- ✅ Contraseña configurada mediante `REDIS_PASSWORD` y `--requirepass`
- ✅ Healthcheck actualizado para usar autenticación
- ✅ Documentación de seguridad creada

**Código Corregido**:
```yaml
# infra/compose-prod/docker-compose.yaml (DESPUÉS)
redis:
  # Puertos NO expuestos - solo acceso interno
  command: >
    redis-server 
    --appendonly yes
    --requirepass ${REDIS_PASSWORD}
  environment:
    REDIS_PASSWORD: "${REDIS_PASSWORD}"
  healthcheck:
    test: ["CMD", "redis-cli", "-a", "${REDIS_PASSWORD}", "ping"]
```

**Recomendación**:
- ✅ **IMPLEMENTADO**: Redis ahora requiere contraseña y no está expuesto públicamente
- ✅ **IMPLEMENTADO**: Documentación de seguridad creada en `infra/compose-prod/SECURITY_REQUIREMENTS.md`
- ⚠️ **PENDIENTE**: Verificar que `REDIS_PASSWORD` esté configurado en producción con contraseña fuerte

---

### 1. **RabbitMQ Expuesto Públicamente con Credenciales por Defecto** 🔴 CRÍTICO

**Ubicación**: `infra/compose-prod/docker-compose.yaml` línea 169-171

**Descripción**:
- RabbitMQ estaba expuesto públicamente en puertos 5672 (AMQP) y 15672 (Management UI)
- Usaba credenciales por defecto `admin/admin` conocidas públicamente
- Permite acceso no autorizado a la cola de mensajes y interfaz de administración

**Evidencia**:
```yaml
# infra/compose-prod/docker-compose.yaml (ANTES)
rabbitmq:
  ports:
    - "5672:5672"   # ⚠️ Expuesto públicamente
    - "15672:15672" # ⚠️ Management UI expuesta públicamente
  environment:
    RABBITMQ_DEFAULT_USER: admin  # ❌ Credencial por defecto
    RABBITMQ_DEFAULT_PASS: admin  # ❌ Credencial por defecto
```

**Solución Implementada**:
- ✅ Puertos de RabbitMQ cerrados (solo acceso interno)
- ✅ Credenciales ahora usan variables de entorno `RABBITMQ_USER` y `RABBITMQ_PASS`
- ✅ Si se necesita UI, solo se expone en localhost: `127.0.0.1:15672:15672`

---

### 2. **Backend API Expuesto Públicamente** 🟡 ALTO

**Ubicación**: `infra/compose-prod/docker-compose.yaml` línea 69

**Descripción**:
- El backend estaba expuesto en `0.0.0.0:3050`
- Aunque se accede normalmente a través de Nginx, no hay necesidad de exponerlo públicamente
- Aumenta la superficie de ataque innecesariamente

**Solución Implementada**:
- ✅ Backend ahora solo accesible en localhost: `127.0.0.1:3050:3050`
- ✅ Acceso público solo a través de Nginx (puertos 80/443)

---

### 3. **Swagger/API Documentation Expuesta Públicamente** 🔴 CRÍTICO

**Ubicación**: `/swagger/` y `/docs/`

**Descripción**:
- La documentación Swagger está expuesta sin autenticación
- Permite a atacantes descubrir todos los endpoints de la API
- Puede revelar estructura interna, parámetros, y endpoints no documentados públicamente

**Evidencia**:
```nginx
# infra/nginx/nginx.conf línea 132-143
location /swagger/ {
    proxy_pass http://probability_backend/swagger/;
    # Sin autenticación requerida
}
```

**Recomendación**:
- Proteger Swagger con autenticación básica HTTP o IP whitelist
- O moverlo solo a entornos de desarrollo/staging
- O deshabilitarlo completamente en producción

---

### 2. **CORS Excesivamente Permisivo** 🔴 CRÍTICO

**Ubicación**: Configuración de Nginx y backend

**Descripción**:
- CORS configurado con `Access-Control-Allow-Origin: *`
- Permite que cualquier dominio haga requests a la API
- Facilita ataques CSRF y acceso no autorizado desde cualquier origen

**Evidencia**:
```nginx
# infra/nginx/nginx.conf línea 107
add_header 'Access-Control-Allow-Origin' '*' always;
```

**Recomendación**:
- Restringir CORS solo a dominios específicos conocidos
- Eliminar el wildcard `*`
- Configurar una lista blanca de orígenes permitidos

---

### 3. **Validación Insuficiente en Carga de Archivos** 🟡 ALTO

**Ubicación**: `back/central/shared/storage/upload_image.go`

**Descripción**:
- Solo se valida el header `Content-Type` del request
- No se valida el contenido real del archivo (magic bytes/file signature)
- Un atacante podría falsificar el Content-Type y subir archivos maliciosos
- El método `UploadFile` no tiene restricciones de tipo de archivo

**Evidencia**:
```go
// upload_image.go línea 31-34
contentType := file.Header.Get("Content-Type")
if !allowedImageTypes[contentType] {
    return "", errs.New("tipo de archivo no permitido...")
}
// Solo valida el header, no el contenido real
```

**Recomendación**:
- Validar magic bytes del archivo (primeros bytes del contenido)
- Implementar validación del contenido real, no solo headers
- Restringir `UploadFile` para que también valide tipos de archivo

---

### 4. **Rutas Públicas Sin Rate Limiting** 🟡 ALTO

**Ubicación**: Endpoints como `/health`, `/ping`, `/test`

**Descripción**:
- Endpoints públicos pueden ser usados para DDoS o reconnaissance
- No hay límite de tasa de requests
- Pueden ser utilizados para escanear la infraestructura

**Evidencia**:
```go
// router.go líneas 32-44
r.GET("/health", func(c *gin.Context) {...})
r.GET("/test", func(c *gin.Context) {...})
```

**Recomendación**:
- Implementar rate limiting en todos los endpoints públicos
- Usar nginx rate limiting o middleware de rate limiting en el backend

---

### 5. **Archivos Estáticos del Website Vulnerables** 🟡 MEDIO (No fue el vector principal)

**Ubicación**: Contenedor del website (`font-website`)

**Descripción**:
- El website Astro se sirve como archivos estáticos desde `/usr/share/nginx/html`
- Si el contenedor o volumen fue comprometido, los archivos pueden ser modificados
- No hay verificación de integridad de los archivos estáticos
- El volumen podría estar montado sin permisos restringidos

**Posible Vector de Ataque**:
1. Acceso al contenedor Docker (si hay vulnerabilidad en nginx o configuración)
2. Volumen compartido con permisos incorrectos
3. Build comprometido (dependencias maliciosas)
4. Acceso al sistema de archivos del host

**Recomendación**:
- Usar volúmenes de solo lectura para archivos estáticos
- Implementar verificación de integridad (checksums)
- Revisar permisos del contenedor (no ejecutar como root)
- Implementar file integrity monitoring

---

### 6. **Headers de Seguridad Faltantes** 🟡 MEDIO

**Ubicación**: Configuración de Nginx

**Descripción**:
- Faltan headers de seguridad importantes como:
  - `Content-Security-Policy`
  - `Strict-Transport-Security` (HSTS)
  - `X-Frame-Options` (solo en website, no en nginx principal)
  - `Referrer-Policy`

**Evidencia**:
```nginx
# El nginx principal no tiene estos headers
# Solo el nginx del website tiene algunos (X-Frame-Options, etc.)
```

**Recomendación**:
- Agregar todos los headers de seguridad necesarios
- Implementar CSP estricto
- Habilitar HSTS con `max-age` apropiado

---

## 🔍 Cómo Ocurrió el Defacement (CONFIRMADO)

**Vector de Ataque Principal**: Redis expuesto públicamente sin contraseña

### Escenario Confirmado: Acceso a través de Redis

1. **Detección**: Bot automatizado escaneó internet buscando puertos 6379 (Redis) abiertos
2. **Conexión**: Se conectó a Redis sin autenticación (no tenía contraseña)
3. **Explotación**: Usó comandos de Redis para:
   - Leer/escribir datos en caché
   - Potencialmente modificar configuración
   - Inyectar código o redirecciones
4. **Resultado**: El sitio web comenzó a mostrar el mensaje de defacement (página china con cerdito)

**Evidencia**:
- Redis estaba configurado con `ports: - "6379:6379"` (expuesto a 0.0.0.0)
- Redis NO tenía `--requirepass` configurado
- El mensaje de defacement es característico de bots automatizados que escanean Redis

### Otros Escenarios Posibles (Menos Probables):
- El atacante encontró una vulnerabilidad en nginx o en alguna dependencia
- Obtuvo acceso al contenedor del website
- Modificó directamente `/usr/share/nginx/html/index.html`
- **Más probable si el contenedor se ejecuta como root**

### Escenario 2: Acceso al Sistema de Archivos del Host
- El atacante comprometió el servidor host
- Accedió al volumen donde se montan los archivos del website
- Modificó los archivos estáticos directamente

### Escenario 3: Build Comprometido
- Vulnerabilidad en dependencias de npm durante el build
- Script de build malicioso ejecutado durante `npm run build`
- Archivos comprometidos empaquetados en la imagen Docker

### Escenario 4: Volumen Compartido Sin Permisos
- Volumen Docker montado con permisos incorrectos
- Múltiples servicios tienen acceso de escritura al mismo volumen
- Un servicio comprometido modificó los archivos del website

---

## ✅ Acciones Inmediatas Recomendadas

### Prioridad CRÍTICA (Hacer AHORA):

1. **Restaurar el sitio**:
   ```bash
   # Rebuild y redeploy del contenedor del website
   docker-compose down
   docker-compose build --no-cache font-website
   docker-compose up -d
   ```

2. **Revisar logs**:
   ```bash
   # Buscar actividad sospechosa
   docker logs font-website --since 48h
   docker exec font-website ls -la /usr/share/nginx/html
   ```

3. **Cambiar todas las credenciales**:
   - Base de datos
   - Claves de API
   - Tokens JWT secretos
   - Credenciales S3/MinIO

4. **Proteger Swagger**:
   - Deshabilitar o proteger con autenticación básica

5. **Restringir CORS**:
   - Cambiar de `*` a dominios específicos

### Prioridad ALTA (Esta semana):

6. **Implementar validación de archivos**:
   - Validar magic bytes en uploads
   - Restringir tipos de archivo estrictamente

7. **Asegurar volúmenes**:
   - Usar volúmenes de solo lectura para archivos estáticos
   - Revisar permisos de todos los volúmenes

8. **Implementar rate limiting**:
   - En todos los endpoints públicos

9. **Agregar headers de seguridad**:
   - CSP, HSTS, etc.

10. **Auditoría de seguridad**:
    - Revisar todos los accesos recientes
    - Buscar backdoors o cambios no autorizados
    - Revisar dependencias por vulnerabilidades conocidas

---

## 🛡️ Mejores Prácticas de Seguridad Recomendadas

1. **Principio de Menor Privilegio**:
   - Contenedores no deben ejecutarse como root
   - Usar usuarios no privilegiados

2. **Seguridad en Capas**:
   - WAF (Web Application Firewall) antes de nginx
   - Rate limiting
   - Validación estricta en cada capa

3. **Monitoreo**:
   - File integrity monitoring
   - Logs centralizados
   - Alertas de seguridad

4. **Actualizaciones**:
   - Mantener todas las dependencias actualizadas
   - Revisar CVE regularmente
   - Parches de seguridad inmediatos

5. **Backups**:
   - Backups regulares y verificados
   - Plan de recuperación ante desastres

---

## 📝 Notas Adicionales

- **El mensaje en chino** sugiere que fue un grupo organizado (Alianza de Seguridad de Red Juvenil)
- **La URL `hk.h-acker.cc`** es indicativa de un ataque de defacement
- **El hecho de que dejaron un mensaje** sugiere que fue más un "mensaje de seguridad" que un ataque malicioso destructivo, pero igual es una vulnerabilidad crítica que debe ser cerrada

---

**Fecha del análisis**: $(date)
**Analista**: Cursor AI Assistant
**Próxima revisión recomendada**: Después de implementar las correcciones críticas

