# 🔐 Guía de Seguridad SSH Después del Incidente

## ⚠️ ¿Es Peligroso Entrar por SSH?

**Respuesta corta**: SÍ, puede ser peligroso, pero es **NECESARIO** para investigar y reparar. Sigue estos pasos con precaución.

## 🛡️ Precauciones ANTES de Conectarte

### 1. **Cambiar tu Contraseña SSH ANTES de Conectarte** (si es posible)

Si tienes acceso a un panel de control o a otra forma de gestionar el servidor:
- Cambia tu contraseña SSH
- Si usas autenticación por clave, considera rotar tus claves SSH

### 2. **Revisar tus Claves SSH Locales**

```bash
# En tu máquina LOCAL, antes de conectarte
ls -la ~/.ssh/
# Verifica que solo TÚ tengas acceso
chmod 600 ~/.ssh/id_rsa
chmod 644 ~/.ssh/id_rsa.pub
```

### 3. **Usar una Conexión VPN o Red Segura**

- Evita conectarte desde redes públicas (cafés, aeropuertos)
- Usa una VPN si es posible
- Preferiblemente desde una red privada de confianza

### 4. **Anotar TODA tu Sesión**

```bash
# Antes de conectarte, configura logging
script -a ssh-session-$(date +%Y%m%d-%H%M%S).log
# Ahora conecta por SSH
ssh usuario@servidor
# Cuando termines, escribe 'exit' dos veces (una para SSH, otra para script)
```

## 🔍 Pasos Seguros para Conectarte

### Paso 1: Conectarte con Logging Habilitado

```bash
# Conecta con verbose para ver detalles de la conexión
ssh -v usuario@servidor

# O con más verbosidad para debugging
ssh -vvv usuario@servidor
```

**Observa**:
- ¿La clave del host cambió? (fingerprint warning)
- ¿Hay mensajes sospechosos?
- ¿El banner de bienvenida cambió?

### Paso 2: Revisar INMEDIATAMENTE al Conectarte

**NO ejecutes nada más hasta revisar esto:**

```bash
# 1. Verificar último acceso y sesiones activas
who
w
last

# 2. Verificar historial de comandos recientes
history | tail -50

# 3. Verificar procesos sospechosos
ps aux | grep -E '(nc|netcat|python|perl|wget|curl|bash|sh)' | grep -v grep

# 4. Verificar conexiones de red activas
netstat -tulpn
# O si no está disponible:
ss -tulpn

# 5. Verificar archivos modificados recientemente
find /var/www /usr/share/nginx/html /home -type f -mtime -7 -ls 2>/dev/null | head -20

# 6. Verificar usuarios nuevos o cambios de permisos
cat /etc/passwd | grep -E '(/bin/bash|/bin/sh)'
```

### Paso 3: Revisar Archivos Críticos

```bash
# Archivos de configuración SSH
ls -la /etc/ssh/sshd_config
cat /etc/ssh/sshd_config | grep -E '(PermitRootLogin|PasswordAuthentication|PubkeyAuthentication)'

# Archivos de crontab (tareas programadas maliciosas)
crontab -l
sudo crontab -l
ls -la /etc/cron.*
cat /etc/crontab

# Archivos .bashrc, .profile (podrían tener backdoors)
cat ~/.bashrc
cat ~/.profile
cat ~/.bash_profile
```

### Paso 4: Verificar Docker y Contenedores

```bash
# Listar contenedores activos
docker ps -a

# Ver logs de contenedores sospechosos
docker logs font-website --since 72h | tail -100

# Verificar archivos dentro del contenedor del website
docker exec font-website ls -la /usr/share/nginx/html/

# Verificar si hay contenedores nuevos o modificados
docker ps --format "table {{.ID}}\t{{.Image}}\t{{.CreatedAt}}\t{{.Status}}"

# Verificar volúmenes
docker volume ls
```

### Paso 5: Buscar Backdoors y Archivos Maliciosos

```bash
# Buscar archivos PHP sospechosos (si hay PHP)
find /var/www /usr/share/nginx -name "*.php" -mtime -7

# Buscar archivos con permisos sospechosos
find /var/www /usr/share/nginx -type f -perm -o+w -ls

# Buscar archivos ocultos
find /var/www /usr/share/nginx -name ".*" -type f

# Buscar archivos con extensiones sospechosas
find /var/www /usr/share/nginx -type f \( -name "*.sh" -o -name "*.py" -o -name "*.pl" \) -mtime -7

# Verificar archivos index.html/index.php modificados
find / -name "index.html" -o -name "index.php" 2>/dev/null | xargs ls -la | grep "$(date +%Y-%m-%d)\|$(date -d '1 day ago' +%Y-%m-%d)"
```

## 🚨 Señales de Alerta (Si Encuentras Esto, El Servidor Está Comprometido)

1. **Nuevos usuarios** en `/etc/passwd` que no reconoces
2. **Procesos desconocidos** ejecutándose
3. **Conexiones de red** a IPs sospechosas
4. **Archivos modificados** recientemente que no deberían cambiar
5. **Crontabs** con comandos que no reconoces
6. **Servicios nuevos** ejecutándose (ver con `systemctl list-units`)
7. **Cambios en archivos de configuración** SSH
8. **Claves SSH nuevas** en `~/.ssh/authorized_keys` que no reconoces

## ✅ Acciones Inmediatas SI Encuentras Compromiso

### Si el Servidor Está Claramente Comprometido:

1. **NO cierres la sesión SSH todavía**
2. **Documenta TODO**:
   ```bash
   # Capturar estado actual
   ps aux > procesos-actuales.txt
   netstat -tulpn > conexiones-actuales.txt
   history > historial-comandos.txt
   ```

3. **Desconecta el servidor de la red** (si puedes):
   - Cierra los puertos críticos
   - O detén los servicios Docker

4. **Toma un snapshot/backup** del servidor antes de limpiarlo

5. **NO elimines archivos todavía** - necesitas evidencia

6. **Contacta a tu proveedor** o administrador de sistemas

### Si el Servidor Parece Limpio:

1. **Restaurar archivos del website**:
   ```bash
   # Rebuild del contenedor del website
   cd /ruta/al/proyecto
   docker-compose down font-website
   docker-compose build --no-cache font-website
   docker-compose up -d font-website
   ```

2. **Verificar integridad**:
   ```bash
   # Verificar que el index.html es el correcto
   docker exec font-website cat /usr/share/nginx/html/index.html | head -20
   ```

3. **Cambiar TODAS las credenciales**:
   - Contraseñas de base de datos
   - Secrets de JWT
   - Claves de API
   - Credenciales S3/MinIO

## 🔒 Mejoras de Seguridad SSH Recomendadas

Una vez que hayas limpiado el servidor:

### 1. Deshabilitar Login Root

```bash
sudo nano /etc/ssh/sshd_config
# Cambiar:
PermitRootLogin no
PasswordAuthentication no  # Solo usar claves SSH
PubkeyAuthentication yes
```

### 2. Cambiar Puerto SSH (Opcional pero Recomendado)

```bash
# En /etc/ssh/sshd_config
Port 2222  # Cambiar de 22 a otro puerto
```

### 3. Configurar Fail2Ban

```bash
# Instalar fail2ban
sudo apt-get update
sudo apt-get install fail2ban -y

# Configurar para SSH
sudo nano /etc/fail2ban/jail.local
```

Contenido de `jail.local`:
```ini
[sshd]
enabled = true
port = 22  # O el puerto que uses
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
```

```bash
sudo systemctl restart fail2ban
```

### 4. Usar Solo Autenticación por Claves

```bash
# Generar clave SSH en tu máquina local (si no tienes)
ssh-keygen -t ed25519 -C "tu-email@ejemplo.com"

# Copiar clave al servidor
ssh-copy-id usuario@servidor

# Luego deshabilitar contraseñas en sshd_config
```

### 5. Configurar IP Whitelist (Si es Posible)

Si solo necesitas acceso desde IPs específicas:
```bash
# En /etc/ssh/sshd_config
AllowUsers usuario@IP_PERMITIDA
# O usar ufw/firewall
sudo ufw allow from TU_IP to any port 22
```

## 📋 Checklist de Seguridad Post-Incidente

- [ ] Cambiar todas las contraseñas y credenciales
- [ ] Revisar logs de acceso SSH (`/var/log/auth.log`)
- [ ] Revisar logs de nginx (`/var/log/nginx/`)
- [ ] Revisar logs de Docker
- [ ] Verificar integridad de archivos críticos
- [ ] Implementar Fail2Ban
- [ ] Configurar autenticación solo por claves SSH
- [ ] Deshabilitar Swagger público
- [ ] Restringir CORS
- [ ] Implementar rate limiting
- [ ] Configurar monitoreo de archivos (file integrity monitoring)
- [ ] Hacer backup completo del servidor
- [ ] Documentar el incidente

## 🆘 Si No Puedes Conectarte o Sientes que Es Muy Peligroso

1. **Contacta a tu proveedor de hosting** inmediatamente
2. **Pide que cambien las credenciales SSH** desde el panel
3. **Solicita un snapshot del servidor** antes de hacer cambios
4. **Considera contratar un experto en seguridad** si no te sientes cómodo

## 📝 Notas Importantes

- **Mantén esta sesión documentada** - guarda todos los logs
- **No entres desde múltiples sesiones** - usa solo una conexión
- **Si algo se ve sospechoso, desconéctate inmediatamente**
- **No ejecutes comandos que no entiendas completamente**

---

**⚠️ ADVERTENCIA CRÍTICA**: También encontré que tienes **credenciales de base de datos en texto plano** en `.vscode/settings.json`. Esto es un riesgo de seguridad grave. Cambia esas credenciales INMEDIATAMENTE.

