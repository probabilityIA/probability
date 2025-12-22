# Requisitos de Seguridad - Producción

## ⚠️ ADVERTENCIAS CRÍTICAS

Este archivo documenta los requisitos de seguridad para el despliegue en producción. **LEE ESTO ANTES DE DESPLEGAR**.

## 🔐 Variables de Entorno Requeridas

### Redis - CONTRASEÑA OBLIGATORIA

**CRÍTICO**: Redis DEBE tener una contraseña configurada en producción.

```bash
# En tu archivo .env de producción, DEBES configurar:
REDIS_PASSWORD=tu_contraseña_muy_segura_aqui_minimo_32_caracteres
```

**¿Por qué?**
- Redis sin contraseña es accesible por cualquiera desde internet si el puerto está expuesto
- Esto fue el vector de ataque que causó el defacement del sitio
- Incluso si el puerto no está expuesto, es una buena práctica de seguridad

**Generar contraseña segura:**
```bash
# Opción 1: OpenSSL
openssl rand -base64 32

# Opción 2: Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"
```

### RabbitMQ - Cambiar Credenciales por Defecto

**CRÍTICO**: NO uses las credenciales por defecto `admin/admin` en producción.

```bash
# En tu archivo .env de producción, DEBES cambiar:
RABBITMQ_USER=tu_usuario_seguro
RABBITMQ_PASS=tu_contraseña_muy_segura_aqui
```

**¿Por qué?**
- Las credenciales por defecto son conocidas públicamente
- RabbitMQ expone una interfaz de administración que puede ser explotada
- Cambiar las credenciales es esencial para seguridad

### JWT Secret - Contraseña Fuerte

```bash
JWT_SECRET=tu_jwt_secret_muy_seguro_minimo_64_caracteres
```

**Generar JWT secret seguro:**
```bash
openssl rand -hex 32
```

### Encryption Key - 32 Caracteres

```bash
ENCRYPTION_KEY=tu_clave_de_encriptacion_exactamente_32_caracteres
```

## 🚫 Puertos NO Expuestos

Los siguientes servicios **NO deben tener puertos expuestos públicamente** en producción:

- ✅ **Redis (6379)**: Solo acceso interno vía red Docker
- ✅ **RabbitMQ AMQP (5672)**: Solo acceso interno vía red Docker  
- ✅ **RabbitMQ Management UI (15672)**: Solo localhost si se necesita (127.0.0.1:15672:15672)
- ✅ **Backend API (3050)**: Solo localhost para debugging (127.0.0.1:3050:3050)

**Puertos que SÍ deben estar expuestos:**
- ✅ **Nginx HTTP (80)**: Acceso público
- ✅ **Nginx HTTPS (443)**: Acceso público

## 🔍 Verificación Post-Despliegue

Después de desplegar, verifica que los puertos no estén expuestos:

```bash
# Verificar puertos expuestos
docker ps --format "table {{.Names}}\t{{.Ports}}"

# Verificar que Redis NO está accesible desde fuera
# Esto DEBE fallar:
redis-cli -h TU_IP_SERVIDOR -p 6379 ping

# Verificar que RabbitMQ NO está accesible desde fuera
# Esto DEBE fallar:
telnet TU_IP_SERVIDOR 5672
```

## 📋 Checklist de Seguridad Pre-Despliegue

- [ ] `REDIS_PASSWORD` configurado con contraseña fuerte (mínimo 32 caracteres)
- [ ] `RABBITMQ_USER` cambiado de "admin"
- [ ] `RABBITMQ_PASS` cambiado de "admin" a contraseña fuerte
- [ ] `JWT_SECRET` configurado con valor seguro (mínimo 64 caracteres)
- [ ] `ENCRYPTION_KEY` configurado con exactamente 32 caracteres
- [ ] `DB_PASSWORD` configurado con contraseña fuerte
- [ ] Puertos de Redis, RabbitMQ y Backend NO expuestos públicamente
- [ ] Archivo `.env` NO está en el repositorio (verificar .gitignore)
- [ ] Certificados SSL configurados correctamente
- [ ] Nginx configurado con headers de seguridad

## 🛡️ Mejores Prácticas Adicionales

1. **Rotar contraseñas regularmente**: Cambia todas las contraseñas cada 90 días
2. **Monitoreo**: Configura alertas para intentos de acceso fallidos
3. **Backups**: Asegura backups regulares de la base de datos
4. **Logs**: Revisa logs regularmente para actividad sospechosa
5. **Actualizaciones**: Mantén todas las imágenes Docker actualizadas
6. **Firewall**: Configura un firewall para bloquear acceso no autorizado

## 📞 Contacto de Seguridad

Si encuentras una vulnerabilidad, reporta inmediatamente al equipo de seguridad.

---

**Última actualización**: Después del incidente de seguridad del [fecha]
**Versión**: 1.0

