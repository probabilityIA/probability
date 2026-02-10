# Simulador de Softpymes

Servidor mock que simula la API de Softpymes para pruebas de integración.

## 🚀 Inicio Rápido

### 1. Compilar el servidor (si no existe el binario)

```bash
go build -o softpymes-server server/main.go
```

### 2. Iniciar el servidor

```bash
# Usando el script (puerto 9090 por defecto)
./start-mock.sh

# O con puerto personalizado
./start-mock.sh 8082

# O manualmente
SOFTPYMES_MOCK_PORT=9090 ./softpymes-server
```

### 3. Configurar el backend para usar el mock

```bash
# En el directorio back/central
export SOFTPYMES_API_URL=http://localhost:9090

# Iniciar el backend
go run cmd/main.go
```

## 📋 Endpoints Disponibles

### 1. Autenticación
```bash
POST /oauth/integration/login/
Content-Type: application/json
Referer: {referer}

{
  "apiKey": "cualquier-valor",
  "apiSecret": "cualquier-valor"
}

# Respuesta
{
  "accessToken": "spy_token_xxxxxxxx",
  "expiresInMin": 60,
  "tokenType": "Bearer"
}
```

### 2. Crear Factura (SIN URLs de PDF/XML)
```bash
POST /app/integration/sales_invoice/
Authorization: Bearer {token}
Content-Type: application/json
```

### 3. Buscar Documentos (retorna JSON completo)
```bash
POST /app/integration/search/documents/
Authorization: Bearer {token}
Content-Type: application/json
```

### 4. Health Check
```bash
GET /health
```

## 🔄 Cambiar entre Mock y API Real

```bash
# Usar el simulador local
export SOFTPYMES_API_URL=http://localhost:9090

# Usar la API real de Softpymes
export SOFTPYMES_API_URL=https://api-integracion.softpymes.com.co
```

## ⚠️ Notas Importantes

1. **Autenticación**: El simulador siempre retorna 200 (no valida credenciales reales)
2. **Sin PDF/XML URLs**: Replica la limitación real de la API - no retorna URLs en la respuesta de creación
3. **Búsqueda**: Retorna el JSON completo del documento (igual que la API real)
4. **Tokens**: Los tokens son válidos solo para esta sesión
5. **Almacenamiento en memoria**: Los datos se pierden al reiniciar

## 🧪 Prueba Rápida

```bash
# 1. Autenticar
TOKEN=$(curl -s -X POST http://localhost:9090/oauth/integration/login/ \
  -H "Content-Type: application/json" \
  -H "Referer: 901497840" \
  -d '{"apiKey":"test","apiSecret":"secret"}' | jq -r '.accessToken')

echo "Token: $TOKEN"

# 2. Crear factura
curl -s -X POST http://localhost:9090/app/integration/sales_invoice/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "customer": {"name":"Test","email":"test@test.com","nit":"123"},
    "items": [{"itemCode":"P1","description":"Product 1","quantity":1,"unitPrice":100}],
    "total": 100
  }' | jq .

# 3. Buscar documentos
curl -s -X POST http://localhost:9090/app/integration/search/documents/ \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"dateFrom":"2026-02-01","dateTo":"2026-02-28"}' | jq .
```
