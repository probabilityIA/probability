#!/bin/bash

# Script para limpiar cache y resolver errores de Astro
echo "🧹 Limpiando cache de Astro..."

# Eliminar directorio dist
if [ -d "dist" ]; then
    echo "📁 Eliminando directorio dist..."
    rm -rf dist
fi

# Eliminar directorio .astro
if [ -d ".astro" ]; then
    echo "📁 Eliminando directorio .astro..."
    rm -rf .astro
fi

# Eliminar cache de node_modules
if [ -d "node_modules/.cache" ]; then
    echo "📁 Eliminando cache de node_modules..."
    rm -rf node_modules/.cache
fi

# Limpiar cache de npm
echo "🧹 Limpiando cache de npm..."
npm cache clean --force

# Limpiar archivos temporales
echo "🧹 Limpiando archivos temporales..."
find . -name "*.tmp*" -type f -delete 2>/dev/null || true

echo "✅ Limpieza completada!"
echo "🚀 Iniciando servidor de desarrollo..."

# Iniciar servidor de desarrollo
npm run dev
