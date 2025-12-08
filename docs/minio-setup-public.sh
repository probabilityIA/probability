#!/bin/bash

# Script para hacer público el bucket de MinIO
# Requiere tener mc (MinIO Client) instalado

# Cargar variables de entorno si existe .env
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Configurar valores (usa variables de entorno o valores por defecto)
MINIO_ENDPOINT="${S3_ENDPOINT:-http://localhost:9000}"
MINIO_ACCESS_KEY="${S3_KEY:-minioadmin}"  # Default de MinIO
MINIO_SECRET_KEY="${S3_SECRET:-minioadmin}"  # Default de MinIO
BUCKET_NAME="${S3_BUCKET:-probability}"

# Verificar que mc esté instalado
if ! command -v mc &> /dev/null; then
    echo "❌ Error: mc (MinIO Client) no está instalado"
    echo "💡 Instala mc desde: https://min.io/docs/minio/linux/reference/minio-mc.html"
    exit 1
fi

echo "🔧 Configurando MinIO para acceso público..."
echo "   Endpoint: ${MINIO_ENDPOINT}"
echo "   Bucket: ${BUCKET_NAME}"
echo ""

# Configurar alias
echo "📝 Configurando alias de MinIO..."
mc alias set local ${MINIO_ENDPOINT} ${MINIO_ACCESS_KEY} ${MINIO_SECRET_KEY} 2>/dev/null || {
    echo "⚠️  Alias ya existe o hay un error. Continuando..."
}

# Hacer el bucket público (permite lectura pública de objetos)
echo "🌐 Configurando bucket como público..."
mc anonymous set public local/${BUCKET_NAME} 2>&1

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Bucket ${BUCKET_NAME} configurado como público"
    echo ""
    echo "📋 Para verificar, prueba acceder a:"
    echo "   ${MINIO_ENDPOINT}/${BUCKET_NAME}/businesslogo/tu-imagen.jpg"
    echo ""
    echo "💡 También puedes verificar desde la consola web:"
    echo "   http://localhost:9001"
else
    echo ""
    echo "⚠️  Hubo un error al configurar el bucket como público"
    echo "💡 Puedes hacerlo manualmente desde la consola web:"
    echo "   1. Abre http://localhost:9001"
    echo "   2. Selecciona el bucket ${BUCKET_NAME}"
    echo "   3. Ve a 'Access Policy' y selecciona 'Public'"
fi
