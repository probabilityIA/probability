# Módulo de Simulación de SoftPymes

Módulo de pruebas de integración para el sistema de facturación electrónica de SoftPymes, siguiendo arquitectura hexagonal.

## 🏗️ Arquitectura Hexagonal

```
softpymes/
├── bundle.go                        # Punto de entrada del módulo
└── internal/
    ├── domain/                      # Capa de dominio (sin dependencias externas)
    │   ├── entities.go              # Entidades (Invoice, CreditNote, AuthToken)
    │   ├── ports.go                 # Interfaces (IAPIClient)
    │   └── repository.go            # Repositorio en memoria
    └── app/
        └── usecases/                # Casos de uso
            ├── constructor.go       # Constructor de use cases
            └── api_simulator.go     # Lógica de simulación
```

## 🎯 Funcionalidades

### 1. Autenticación
Simular autenticación con SoftPymes API:

```go
token, err := softpymesIntegration.SimulateAuth("api_key", "api_secret", "https://tutienda.com")
```

### 2. Creación de Facturas
Simular creación de facturas electrónicas:

```go
invoiceData := map[string]interface{}{
    "order_id": "ORD-001",
    "customer": map[string]interface{}{
        "name": "Juan Pérez",
        "email": "juan@example.com",
        "nit": "123456789",
    },
    "items": []interface{}{...},
    "total": 100000.0,
}

invoice, err := softpymesIntegration.SimulateInvoice(token, invoiceData)
```

### 3. Notas de Crédito
Simular creación de notas de crédito:

```go
creditNoteData := map[string]interface{}{
    "invoice_id": "<external_id>",
    "amount": 50000.0,
    "reason": "Devolución de producto",
    "note_type": "partial",
}

creditNote, err := softpymesIntegration.SimulateCreditNote(token, creditNoteData)
```

### 4. Listar Documentos
Ver todas las facturas y notas de crédito simuladas:

```go
repo := softpymesIntegration.GetRepository()
invoices := repo.GetAllInvoices()
creditNotes := repo.GetAllCreditNotes()
```

## 📊 Entidades del Dominio

### Invoice (Factura)
```go
type Invoice struct {
    ID            string    // UUID
    InvoiceNumber string    // SPY-1001, SPY-1002, etc.
    ExternalID    string    // UUID (usado para referencias)
    OrderID       string    // ID de la orden origen
    CustomerName  string
    CustomerEmail string
    CustomerNIT   string
    Total         float64
    Currency      string    // "COP"
    Items         []InvoiceItem
    InvoiceURL    string    // https://softpymes-mock.local/invoices/{id}
    PDFURL        string    // https://softpymes-mock.local/invoices/{id}.pdf
    XMLURL        string    // https://softpymes-mock.local/invoices/{id}.xml
    CUFE          string    // CUFE-{uuid}
    IssuedAt      time.Time
    CreatedAt     time.Time
}
```

### CreditNote (Nota de Crédito)
```go
type CreditNote struct {
    ID               string    // UUID
    CreditNoteNumber string    // NC-2001, NC-2002, etc.
    ExternalID       string    // UUID
    InvoiceID        string    // External ID de la factura
    Amount           float64   // Monto a acreditar
    Reason           string    // Razón de la nota
    NoteType         string    // "total" o "partial"
    NoteURL          string
    PDFURL           string
    XMLURL           string
    CUFE             string
    IssuedAt         time.Time
    CreatedAt        time.Time
}
```

## 🚀 Uso desde el Menú Interactivo

```bash
cd /back/integrationTest
go run cmd/main.go
```

**Opciones disponibles:**

```
📄 SOFTPYMES (Facturación):
11. Simular autenticación
12. Simular creación de factura
13. Simular nota de crédito
14. Listar facturas almacenadas
```

### Ejemplo de Flujo

#### 1. Autenticar (Opción 11)

```
Opción: 11
API Key: test_key_123
API Secret: test_secret_456
Referer (ej: https://tutienda.com): https://mitienda.com.co

✅ Token generado: spy_token_a1b2c3d4
💡 Guarda este token para crear facturas
```

#### 2. Crear Factura (Opción 12)

```
Opción: 12
Token (obtenido en opción 11): spy_token_a1b2c3d4
Order ID (ej: ORD-001): ORD-001
Nombre cliente: Juan Pérez
Email cliente: juan@example.com
NIT cliente: 123456789
Total (ej: 100000): 100000

✅ Factura creada:
  Número: SPY-1001
  CUFE: CUFE-a1b2c3d4e5f6g7h8
  Total: $100000.00 COP
  PDF: https://softpymes-mock.local/invoices/{uuid}.pdf
```

#### 3. Crear Nota de Crédito (Opción 13)

```
Opción: 13
Token: spy_token_a1b2c3d4
Invoice ID (external_id de la factura): <uuid-de-la-factura>
Monto a acreditar: 50000
Razón (ej: Devolución de producto): Producto defectuoso
Tipo (total/partial): partial

✅ Nota de crédito creada:
  Número: NC-2001
  CUFE: CUFE-NC-i9j8k7l6m5n4
  Monto: $50000.00
  Tipo: partial
  PDF: https://softpymes-mock.local/credit-notes/{uuid}.pdf
```

#### 4. Listar Documentos (Opción 14)

```
Opción: 14

📄 Facturas almacenadas (3):
  1. SPY-1001 - ORD-001 - $100000.00 COP - Cliente: Juan Pérez
  2. SPY-1002 - ORD-002 - $250000.00 COP - Cliente: María López
  3. SPY-1003 - ORD-003 - $75000.00 COP - Cliente: Carlos Gómez

💳 Notas de crédito almacenadas (1):
  1. NC-2001 - Factura: <uuid> - $50000.00 - Tipo: partial
```

## 📋 Numeración de Documentos

| Documento | Formato | Secuencia | Ejemplo |
|-----------|---------|-----------|---------|
| **Factura** | SPY-NNNN | Inicia en 1001 | SPY-1001, SPY-1002 |
| **Nota de Crédito** | NC-NNNN | Inicia en 2001 | NC-2001, NC-2002 |

## 🔐 Autenticación

El simulador genera tokens ficticios:
- **Formato:** `spy_token_{random_8_chars}`
- **Expiración:** 1 hora desde creación
- **Validación:** El token debe existir y no estar expirado

**Ejemplo de token:** `spy_token_a1b2c3d4`

## 🧪 Escenarios de Testing

### Escenario 1: Facturar Orden Completa

1. **Central crea orden** en el sistema
2. **Opción 11:** Autenticar en SoftPymes
3. **Opción 12:** Crear factura con datos de la orden
4. **Verificar en BD:** Factura guardada con CUFE y URLs

### Escenario 2: Devolución Parcial

1. **Crear factura** (Opción 12)
2. **Obtener external_id** de la factura (Opción 14)
3. **Opción 13:** Crear nota de crédito parcial
4. **Verificar:** Nota vinculada a la factura correcta

### Escenario 3: Token Expirado

1. **Autenticar** y obtener token
2. **Esperar >1 hora** (o modificar expiración en código)
3. **Intentar crear factura** con token expirado
4. **Resultado esperado:** Error "token expired"

## 🔄 Integración con Sistema Real

El sistema real debe:

1. ✅ Llamar a `/auth` para obtener token
2. ✅ Usar token en header `Authorization: Bearer {token}`
3. ✅ Guardar `external_id`, `invoice_number`, `CUFE` en BD
4. ✅ Almacenar URLs de PDF y XML para descarga
5. ✅ Manejar tokens expirados (401) y re-autenticar

## 📊 Respuestas del Simulador

### Autenticación Exitosa

```json
{
  "token": "spy_token_a1b2c3d4",
  "expires_at": "2026-02-01T23:30:00Z"
}
```

### Factura Creada

```json
{
  "success": true,
  "message": "Invoice created successfully",
  "invoice_number": "SPY-1001",
  "external_id": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
  "invoice_url": "https://softpymes-mock.local/invoices/{uuid}",
  "pdf_url": "https://softpymes-mock.local/invoices/{uuid}.pdf",
  "xml_url": "https://softpymes-mock.local/invoices/{uuid}.xml",
  "cufe": "CUFE-a1b2c3d4e5f6g7h8",
  "issued_at": "2026-02-01T22:30:00Z"
}
```

### Nota de Crédito Creada

```json
{
  "success": true,
  "message": "Credit note created successfully",
  "credit_note_number": "NC-2001",
  "external_id": "b2c3d4e5-f6g7-h8i9-j0k1-l2m3n4o5p6q7",
  "note_url": "https://softpymes-mock.local/credit-notes/{uuid}",
  "pdf_url": "https://softpymes-mock.local/credit-notes/{uuid}.pdf",
  "xml_url": "https://softpymes-mock.local/credit-notes/{uuid}.xml",
  "cufe": "CUFE-NC-i9j8k7l6m5n4",
  "issued_at": "2026-02-01T22:35:00Z"
}
```

## 🐛 Errores Simulados

| Error | Condición | Mensaje |
|-------|-----------|---------|
| **Invalid credentials** | API key o secret vacío | "invalid credentials" |
| **Invalid token** | Token no existe | "invalid token" |
| **Token expired** | Token > 1 hora | "token expired" |
| **Invoice not found** | Invoice ID no existe (nota de crédito) | "invoice not found: {id}" |

## 📝 Notas Importantes

1. **Repositorio en Memoria**: Los documentos se pierden al reiniciar
2. **No Valida API Real**: Solo simula respuestas, no llama a SoftPymes real
3. **URLs Ficticias**: Los links generados son mock (no descargan PDFs reales)
4. **CUFE Simulados**: Los CUFEs son UUIDs, no CUFEs DIAN reales
5. **Same Process**: Corre en el mismo proceso que otros simuladores

## ✅ Arquitectura Hexagonal Verificada

| Capa | Ubicación | ✅ Sin Dependencias Externas |
|------|-----------|----------------------------|
| **Domain** | `internal/domain/` | ✅ Solo `time`, `uuid`, `sync` |
| **Application** | `internal/app/usecases/` | ✅ Solo depende de domain |

## 🔗 Uso en Código

```go
// Inicializar
softpymesIntegration := softpymes.New(logger)

// Autenticar
token, err := softpymesIntegration.SimulateAuth("key", "secret", "https://site.com")

// Crear factura
invoice, err := softpymesIntegration.SimulateInvoice(token, invoiceData)

// Listar facturas
repo := softpymesIntegration.GetRepository()
invoices := repo.GetAllInvoices()
```

---

**Implementado siguiendo:**
- ✅ Arquitectura Hexagonal
- ✅ Domain sin dependencias externas
- ✅ Ports and Adapters pattern
- ✅ Repository pattern (in-memory)
- ✅ Use Cases para lógica de negocio
- ✅ Numeración secuencial de documentos
- ✅ Validación de tokens
- ✅ CUFEs simulados
