# Cliente de API Softpymes

Este cliente implementa la interfaz `IInvoicingProviderClient` para integración con la API de Softpymes.

## Características

- ✅ Autenticación con cache de token (60 minutos)
- ✅ Creación de facturas electrónicas
- ✅ Cancelación de facturas
- ✅ Creación de notas de crédito
- ✅ Consulta de estado de facturas
- ✅ Reintentos automáticos (HTTP 429)
- ✅ Manejo de errores robusto
- ✅ Logging estructurado
- ✅ Thread-safe

## Uso

### Inicialización

```go
import (
    "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/providers/softpymes"
    "github.com/secamc93/probability/back/central/shared/log"
)

// Crear cliente
baseURL := "https://api-integracion.softpymes.com.co/app/integration/"
logger := log.New()
client := softpymes.New(baseURL, logger)
```

### Autenticación

```go
credentials := map[string]interface{}{
    "api_key":    "your_api_key",
    "api_secret": "your_api_secret",
}

token, err := client.Authenticate(ctx, credentials)
if err != nil {
    // Manejar error
}
```

### Crear Factura

```go
import "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"

request := &ports.InvoiceRequest{
    Invoice:      invoice,
    InvoiceItems: items,
    Provider:     provider,
    Config:       config,
}

response, err := client.CreateInvoice(ctx, token, request)
if err != nil {
    // Manejar error
}

// response contiene: InvoiceNumber, CUFE, PDFURL, XMLURL, etc.
```

### Cancelar Factura

```go
err := client.CancelInvoice(ctx, token, externalID, "Razón de cancelación")
if err != nil {
    // Manejar error
}
```

### Crear Nota de Crédito

```go
request := &ports.CreditNoteRequest{
    Invoice:    invoice,
    CreditNote: creditNote,
    Provider:   provider,
}

response, err := client.CreateCreditNote(ctx, token, request)
if err != nil {
    // Manejar error
}
```

## Cache de Token

El cliente mantiene automáticamente un cache del token de autenticación:

- **Duración:** 60 minutos (con 5 minutos de margen)
- **Thread-safe:** Usa `sync.RWMutex`
- **Auto-renovación:** Se renueva automáticamente al expirar
- **Limpieza en error:** Se limpia el cache si hay error 401

No es necesario gestionar manualmente el token, el cliente lo hace automáticamente.

## Configuración del Proveedor

El cliente espera la siguiente configuración en `InvoicingProvider.Config`:

```json
{
  "referer": "900123456",      // NIT del facturador
  "branch_code": "001"         // Código de sucursal
}
```

Y las siguientes credenciales en `InvoicingProvider.Credentials`:

```json
{
  "api_key": "your_api_key",
  "api_secret": "your_api_secret"
}
```

## Mapeo de Tipos de Nota de Crédito

| Tipo de Dominio    | Código Softpymes | Descripción      |
|--------------------|------------------|------------------|
| `cancellation`     | `01`             | Anulación        |
| `correction`       | `02`             | Corrección       |
| `full_refund`      | `03`             | Devolución total |
| `partial_refund`   | `03`             | Devolución parcial |

## Manejo de Errores

El cliente retorna errores descriptivos:

- `"authentication token expired"` - Token expirado (limpia cache automáticamente)
- `"authentication failed: [razón]"` - Error de autenticación
- `"invoice creation failed: [razón]"` - Error al crear factura
- `"missing or invalid api_key"` - Credenciales inválidas

## Logging

El cliente registra todas las operaciones importantes:

```go
// Éxito
log.Info().Str("invoice_number", "FV-123").Msg("Invoice created successfully")

// Error
log.Error().Err(err).Msg("Failed to create invoice")
```

## Reintentos

El cliente HTTP compartido (`shared/httpclient`) maneja automáticamente:

- **Rate limiting (429):** 2 reintentos con 3 segundos de espera
- **Timeout:** 30 segundos por request

## Endpoints de API

| Método | Endpoint                     | Descripción             |
|--------|------------------------------|-------------------------|
| POST   | `/get_token`                 | Autenticación           |
| POST   | `/sales_invoice/`            | Crear factura           |
| POST   | `/sales_invoice/cancel`      | Cancelar factura        |
| GET    | `/sales_invoice/status`      | Consultar estado        |
| POST   | `/search/documents/notes/`   | Crear nota de crédito   |

## Documentación de Softpymes

- API Base URL: `https://api-integracion.softpymes.com.co/app/integration/`
- Documentación: `https://docs.softpymes.com.co/api/`

## Testing

Para probar la conexión:

```go
err := client.ValidateCredentials(ctx, credentials)
if err != nil {
    // Credenciales inválidas
}
```

## Arquitectura

El cliente sigue el patrón de **Arquitectura Hexagonal**:

```
ports.IInvoicingProviderClient (interface de dominio)
           ↑
           |
    softpymes.Client (implementación)
           |
           ↓
  shared/httpclient (cliente HTTP compartido)
```

## Notas Importantes

1. **No exponer credenciales:** Las credenciales deben estar encriptadas en la base de datos
2. **Validar configuración:** Verificar que `referer` y `branch_code` estén configurados
3. **Manejo de CUFE:** El CUFE es crítico para trazabilidad tributaria
4. **Verificar PDFs:** Siempre verificar que las URLs de PDF/XML sean accesibles
