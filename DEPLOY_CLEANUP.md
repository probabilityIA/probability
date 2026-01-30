# 🚀 Deploy y Limpieza Automática de Imágenes

## 📋 Problema Resuelto

Las imágenes Docker/Podman se acumulaban en el servidor después de cada deploy, consumiendo espacio en disco innecesariamente.

## ✅ Solución Implementada

Se crearon dos scripts mejorados:

### 1. `update_services_improved.zsh` - Deploy con Limpieza Automática

**Ubicación:** `/home/ubuntu/probability/update_services_improved.zsh`

**Qué hace:**
- ✅ Descarga las imágenes más recientes desde ECR
- ✅ Actualiza los servicios (back-central, font-central, font-website)
- ✅ Verifica que todos los servicios estén corriendo
- ✅ **ELIMINA automáticamente** todas las imágenes antiguas no utilizadas
- ✅ Muestra estadísticas de espacio liberado

**Características:**
- Limpieza segura: Solo elimina imágenes después de verificar que los servicios están corriendo
- Muestra espacio en disco antes y después
- Cuenta cuántas imágenes se eliminaron
- Recarga Nginx automáticamente

**Uso:**
```bash
cd /home/ubuntu/probability
./update_services_improved.zsh
```

### 2. `cleanup_images.sh` - Limpieza Manual de Imágenes

**Ubicación:** `/home/ubuntu/probability/cleanup_images.sh`

**Qué hace:**
- 🗑️ Elimina contenedores detenidos
- 🗑️ Elimina imágenes dangling (`<none>`)
- 🗑️ Elimina TODAS las imágenes no utilizadas (incluso con tags)
- 🗑️ Elimina volúmenes huérfanos
- 🗑️ Elimina redes no utilizadas
- 🗑️ Limpia build cache (solo Docker)

**Uso:**
```bash
# Con confirmación
./cleanup_images.sh

# Sin confirmación (forzado)
./cleanup_images.sh --force
```

## 📊 Comparación de Scripts

| Característica | `update_services.zsh` (Antiguo) | `update_services_improved.zsh` (Nuevo) |
|----------------|----------------------------------|----------------------------------------|
| Actualiza servicios | ✅ | ✅ |
| Elimina dangling images | ✅ | ✅ |
| Elimina imágenes antiguas con tags | ❌ | ✅ |
| Muestra espacio liberado | ❌ | ✅ |
| Muestra estadísticas | ❌ | ✅ |
| Verifica salud de servicios | ✅ | ✅ (mejorado) |
| Limpieza segura | ❌ | ✅ |

## 🎯 Comandos Útiles

### Ver imágenes actuales
```bash
docker images
```

### Ver espacio usado por Docker
```bash
docker system df
```

### Ver cuántas imágenes hay
```bash
docker images -q | wc -l
```

### Limpiar TODO (⚠️ CUIDADO)
```bash
docker system prune -a --volumes -f
```

## 📝 Flujo de Trabajo Recomendado

### Deploy Normal (Recomendado)
```bash
cd /home/ubuntu/probability
./update_services_improved.zsh
```

Este script:
1. Descarga nuevas imágenes
2. Actualiza servicios
3. Verifica que todo funcione
4. Limpia imágenes antiguas automáticamente

### Limpieza Manual (Opcional)
Si necesitas limpiar sin hacer deploy:
```bash
cd /home/ubuntu/probability
./cleanup_images.sh
```

## 🔍 Verificación Post-Deploy

Después de ejecutar el script, verifica:

```bash
# Ver servicios corriendo
docker compose ps

# Ver logs de un servicio específico
docker compose logs -f back-central

# Ver espacio en disco
df -h

# Ver imágenes restantes
docker images
```

## 🛡️ Seguridad del Script

El script mejorado (`update_services_improved.zsh`) tiene protecciones:

1. ✅ **Solo limpia si los servicios están corriendo**: Si algún servicio falla, NO se ejecuta la limpieza
2. ✅ **No elimina imágenes en uso**: Docker/Podman automáticamente protege imágenes de contenedores activos
3. ✅ **Muestra qué se eliminó**: Transparencia total sobre qué imágenes se removieron

## 📈 Beneficios

### Antes
```
Imágenes acumuladas: 50+
Espacio usado: 30GB+
Deploy manual + limpieza manual
```

### Ahora
```
Imágenes: Solo las necesarias (3-5)
Espacio usado: ~5-10GB
Deploy con limpieza automática ✨
```

## ⚠️ Notas Importantes

1. **Script antiguo todavía disponible**: `update_services.zsh` sigue funcionando si prefieres no eliminar imágenes automáticamente
2. **Compatibilidad**: Los scripts funcionan tanto con Docker como con Podman
3. **Sin interrupción**: Los servicios activos NUNCA se detienen durante la limpieza
4. **Reversible**: Si necesitas una imagen antigua, siempre puedes volver a descargarla desde ECR

## 🔄 Migración

Para empezar a usar el nuevo script:

```bash
# Opción 1: Renombrar el antiguo como backup
mv update_services.zsh update_services_old.zsh
mv update_services_improved.zsh update_services.zsh

# Opción 2: Usar el nuevo directamente
./update_services_improved.zsh
```

## 📞 Troubleshooting

### Si la limpieza elimina demasiado
```bash
# Ver qué imágenes están en uso
docker ps -a

# Las imágenes de contenedores activos NUNCA se eliminan
# Solo se eliminan imágenes sin contenedores asociados
```

### Si necesitas una imagen antigua
```bash
# Simplemente vuelve a hacer pull
docker compose pull <servicio>
```

### Si el script falla
```bash
# Ver logs completos
./update_services_improved.zsh 2>&1 | tee deploy.log

# El script se detendrá si algo falla (set -e)
```

## 🎉 Resultado

Ahora cada vez que hagas deploy:
- ✅ Servicios actualizados
- ✅ Imágenes antiguas eliminadas automáticamente
- ✅ Espacio en disco liberado
- ✅ Sin intervención manual necesaria

¡Deploy limpio y automático! 🚀
