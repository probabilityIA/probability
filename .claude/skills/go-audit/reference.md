# Go Audit — Reference Patterns (Probability Project)

These are real patterns from the Probability codebase. Use them as the gold standard when auditing.

---

## Correct Patterns

### 1. RabbitMQ Consumer Loop

**File:** `back/central/shared/rabbitmq/rabbitmq.go:187-230`

The canonical consumer pattern: `select` over `ctx.Done()` and the message channel, each consumer on its own dedicated AMQP channel.

```go
go func() {
    for {
        select {
        case <-ctx.Done():
            r.logger.Info().
                Str("queue", queueName).
                Msg("Stopping consumer due to context cancellation")
            return
        case msg, ok := <-msgs:
            if !ok {
                r.logger.Warn().
                    Str("queue", queueName).
                    Msg("Consumer channel closed")
                return
            }
            if err := handler(msg.Body); err != nil {
                msg.Nack(false, true)  // requeue on error
            } else {
                msg.Ack(false)
            }
        }
    }
}()
```

**Key traits:**
- `ctx.Done()` provides clean shutdown
- `ok` check on channel receive detects closed channel
- Nack with requeue on handler error
- Each consumer gets its own AMQP channel (line 160) to avoid "unexpected command received"

---

### 2. RWMutex for Read-Heavy Registry

**File:** `back/central/services/integrations/core/internal/app/usecaseintegrations/provider_registry.go:11-53`

```go
type providerRegistry struct {
    providers map[int]domain.IIntegrationContract
    mu        sync.RWMutex
}

// Write — exclusive lock
func (r *providerRegistry) Register(integrationType int, provider domain.IIntegrationContract) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.providers[integrationType] = provider
}

// Read — shared lock
func (r *providerRegistry) Get(integrationType int) (domain.IIntegrationContract, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    provider, exists := r.providers[integrationType]
    return provider, exists
}

// Read — shared lock
func (r *providerRegistry) ListRegisteredTypes() []int {
    r.mu.RLock()
    defer r.mu.RUnlock()
    // ...
}
```

**Key traits:**
- `RWMutex` (not plain `Mutex`) because reads far outnumber writes
- `RLock` for all read methods, `Lock` only for Register
- `defer` immediately after lock acquisition — always

---

### 3. Token Cache with Mutex + TTL

**File (Softpymes):** `back/central/services/integrations/invoicing/softpymes/internal/infra/secondary/client/token_cache.go:16-53`

```go
type TokenCache struct {
    mu      sync.Mutex
    entries map[string]*tokenEntry
}

type tokenEntry struct {
    token     string
    expiresAt time.Time
}

func (tc *TokenCache) Get(baseURL string) (string, bool) {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    entry, ok := tc.entries[baseURL]
    if !ok || entry.token == "" || time.Now().After(entry.expiresAt) {
        return "", false
    }
    return entry.token, true
}

func (tc *TokenCache) Set(baseURL, token string, expiresIn int) {
    tc.mu.Lock()
    defer tc.mu.Unlock()
    expirationTime := time.Now().Add(time.Duration(expiresIn-300) * time.Second) // 5-min buffer
    tc.entries[baseURL] = &tokenEntry{token: token, expiresAt: expirationTime}
}
```

**Variant (Factus):** `back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/token_cache.go`
- Uses `sync.RWMutex` with `RLock` for reads (higher read frequency)
- Manages both access and refresh tokens with different TTL buffers

**Key traits:**
- TTL buffer (subtract 5 min from actual expiry) prevents using nearly-expired tokens
- `defer Unlock()` always present
- Thread-safe map access

---

### 4. GORM Closure Transaction

**File:** `back/central/services/auth/users/internal/infra/secondary/repository/user_repository.go:128-180`

```go
return db.Transaction(func(tx *gorm.DB) error {
    // Step 1: Delete old records
    if err := tx.Table("business_staff").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
        return err  // auto-rollback
    }
    if err := tx.Table("user_businesses").Where("user_id = ?", userID).Delete(nil).Error; err != nil {
        return err  // auto-rollback
    }

    // Step 2: Insert new records
    if err := tx.CreateInBatches(businessStaffRecords, 100).Error; err != nil {
        return err  // auto-rollback
    }

    return nil  // auto-commit
})
```

**Key traits:**
- Uses GORM's closure-style `db.Transaction()` — no manual Begin/Commit/Rollback
- Every operation returns its error (no silent failures)
- Only DB operations inside the transaction (no external API calls)
- `return nil` triggers commit, `return err` triggers rollback

---

### 5. SSE Keep-Alive with Ticker

**File:** `back/central/services/events/internal/infra/primary/handlers/sse_handler.go:147-169`

```go
func (h *SSEHandler) keepConnectionAlive(w http.ResponseWriter, connectionID string, ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()                           // cleanup on exit

    done := ctx.Done()                            // cache channel reference
    flusher, hasFlusher := w.(http.Flusher)

    for {
        select {
        case <-ticker.C:                          // periodic heartbeat
            h.sendSSEMessage(w, "keep-alive", "ping")
            if hasFlusher {
                flusher.Flush()
            }
        case <-done:                              // client disconnected
            h.eventManager.RemoveConnection(connectionID)
            h.logger.Info(ctx).
                Str("connection_id", connectionID).
                Msg("Cliente SSE desconectado")
            return
        }
    }
}
```

**Key traits:**
- `defer ticker.Stop()` prevents ticker goroutine leak
- `done := ctx.Done()` pre-assigned (avoids repeated method call in hot loop)
- Clean shutdown via `ctx.Done()` when client disconnects
- Resource cleanup (RemoveConnection) on exit

---

### 6. Channel Send with Backpressure

**File:** `back/central/services/events/internal/infra/secondary/sse/publish.go:30-44`

```go
func (m *EventManager) PublishEvent(event entities.Event) {
    select {
    case m.eventChan <- event:                    // non-blocking send
        if m.logger != nil {
            m.logger.Debug(context.Background()).
                Str("event_id", event.ID).
                Msg("Evento enviado al canal")
        }
    default:                                      // backpressure: drop if full
        if m.logger != nil {
            m.logger.Warn(context.Background()).
                Interface("event", event).
                Msg("Canal de eventos lleno, descartando evento")
        }
    }
}
```

**Key traits:**
- `select` with `default` makes the send non-blocking
- If channel is full, event is dropped with a warning (not silently lost)
- Caller never blocks waiting for channel space

---

### 7. Domain Errors — Grouped Sentinels

**File:** `back/central/services/modules/invoicing/internal/domain/errors/errors.go`

The invoicing module is the best example — errors grouped by HTTP status semantics:

```go
package errors

import "errors"

var (
    // 400 - Validation
    ErrInvalidInvoiceData  = errors.New("invalid invoice data")
    ErrMissingRequiredField = errors.New("missing required field")

    // 404 - Not Found
    ErrInvoiceNotFound = errors.New("invoice not found")
    ErrConfigNotFound  = errors.New("invoicing config not found")

    // 409 - Conflict
    ErrInvoiceAlreadyExists = errors.New("invoice already exists for this order")
    ErrConfigAlreadyExists  = errors.New("invoicing config already exists")

    // 422 - Business Rule
    ErrOrderNotInvoiceable = errors.New("order is not invoiceable")
    ErrMaxRetriesExceeded  = errors.New("max retries exceeded")

    // 401/429 - Provider
    ErrProviderUnauthorized    = errors.New("provider unauthorized")
    ErrProviderRateLimitExceeded = errors.New("provider rate limit exceeded")
)
```

**Key traits:**
- Sentinels grouped by HTTP status category (comments indicating 400/404/409/422)
- English internal messages (handler translates to Spanish)
- Every possible failure case has its own sentinel
- Uses `errors.New()` — pure, no tags, no dependencies

---

### 8. Centralized Error Handler in Handlers

**File:** `back/central/services/modules/invoicing/internal/infra/primary/handlers/error_handler.go`

The gold standard for mapping domain errors to HTTP responses:

```go
func handleDomainError(c *gin.Context, err error, code string) {
    status, message := resolveInvoicingError(err)
    c.JSON(status, response.Error{
        Error:   code,       // machine-readable: "INVOICE_NOT_FOUND"
        Message: message,    // human-readable Spanish: "Factura no encontrada"
    })
}

func resolveInvoicingError(err error) (int, string) {
    switch {
    // 404 Not Found
    case errors.Is(err, invoicingErrors.ErrInvoiceNotFound):
        return http.StatusNotFound, "Factura no encontrada"
    case errors.Is(err, invoicingErrors.ErrConfigNotFound):
        return http.StatusNotFound, "Configuracion de facturacion no encontrada"

    // 409 Conflict
    case errors.Is(err, invoicingErrors.ErrInvoiceAlreadyExists):
        return http.StatusConflict, "Ya existe una factura para esta orden"
    case errors.Is(err, invoicingErrors.ErrConfigAlreadyExists):
        return http.StatusConflict, "Ya existe una configuracion de facturacion"

    // 422 Unprocessable Entity
    case errors.Is(err, invoicingErrors.ErrOrderNotInvoiceable):
        return http.StatusUnprocessableEntity, "La orden no es facturable"
    case errors.Is(err, invoicingErrors.ErrMaxRetriesExceeded):
        return http.StatusUnprocessableEntity, "Se ha superado el numero maximo de reintentos"

    // 400 Bad Request
    case errors.Is(err, invoicingErrors.ErrInvalidInvoiceData):
        return http.StatusBadRequest, "Datos de factura invalidos"

    // 401 Unauthorized
    case errors.Is(err, invoicingErrors.ErrProviderUnauthorized):
        return http.StatusUnauthorized, "No autorizado con el proveedor de facturacion"

    // 429 Too Many Requests
    case errors.Is(err, invoicingErrors.ErrProviderRateLimitExceeded):
        return http.StatusTooManyRequests, "Limite de solicitudes del proveedor excedido"

    // 500 default (catch-all)
    default:
        return http.StatusInternalServerError, "Error interno del servidor"
    }
}
```

**Key traits:**
- Centralized in one file — all error mapping in one place
- `errors.Is()` for safe sentinel comparison (works with wrapped errors)
- Spanish messages for all responses
- Machine-readable `Error` code + human-readable `Message`
- Always has a `default` catch-all for unexpected errors (never panics)

---

### 9. Repository Error Mapping to Domain

**File:** `back/central/services/modules/orders/internal/infra/secondary/repository/repository.go`

```go
// Pattern 1: Not-found → domain sentinel
if err := r.db.Conn(ctx).First(&order).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, domainerrors.ErrOrderNotFound  // mapped to domain
    }
    return nil, fmt.Errorf("error al buscar orden: %w", err)  // wrapped with context
}

// Pattern 2: Duplicate key → domain sentinel
if err := r.db.Conn(ctx).Create(dbOrder).Error; err != nil {
    errMsg := err.Error()
    if strings.Contains(errMsg, "duplicate key value violates unique constraint") &&
        strings.Contains(errMsg, "idx_integration_external_id") {
        return domainerrors.ErrOrderAlreadyExists
    }
    return fmt.Errorf("error al crear orden: %w", err)
}

// Pattern 3: Not-found is normal (return nil, nil)
if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, nil  // caller checks nil — not an error condition
}
```

**Key traits:**
- `gorm.ErrRecordNotFound` always mapped to domain sentinel (never leaks raw)
- SQLSTATE 23505 duplicate detection with `strings.Contains`
- Raw DB errors wrapped with `fmt.Errorf` for context
- `errors.Is()` used (not `==`)

---

### 10. Use Case Error Logging Pattern

**File:** `back/central/services/modules/customers/internal/app/delete_customer.go`

```go
func (uc *UseCase) DeleteClient(ctx context.Context, businessID, clientID uint) error {
    // Step 1: Verify exists (returns ErrClientNotFound if not)
    _, err := uc.repo.GetByID(ctx, businessID, clientID)
    if err != nil {
        uc.log.Error(ctx).Err(err).
            Uint("client_id", clientID).
            Msg("Error al obtener cliente para eliminar")
        return err  // domain error from repo (ErrClientNotFound)
    }

    // Step 2: Business rule check
    orderCount, _, _, err := uc.repo.GetOrderStats(ctx, clientID)
    if err != nil {
        uc.log.Error(ctx).Err(err).
            Uint("client_id", clientID).
            Msg("Error al obtener estadisticas de ordenes")
        return err
    }
    if orderCount > 0 {
        return fmt.Errorf("%w: tiene %d orden(es)", domainerrors.ErrClientHasOrders, orderCount)
        // errors.Is(err, ErrClientHasOrders) == true because of %w wrapping
    }

    // Step 3: Execute
    if err := uc.repo.Delete(ctx, businessID, clientID); err != nil {
        uc.log.Error(ctx).Err(err).
            Uint("client_id", clientID).
            Msg("Error al eliminar cliente")
        return err
    }

    return nil
}
```

**Key traits:**
- Every `if err != nil` has a log before returning
- Logs include context fields (IDs, relevant data)
- Domain sentinels wrapped with `%w` to add context while preserving `errors.Is()`
- Business rule violations use domain sentinels (not generic errors)

---

## Anti-Patterns (Known Issues)

### Anti-Pattern 1: Discarded Error in Goroutine

**File:** `back/central/services/modules/inventory/internal/infra/secondary/redis/event_publisher.go:58-70`

```go
// BUG: error silently discarded
go func() {
    _ = rabbitmq.PublishEvent(context.Background(), p.queue, rabbitmq.EventEnvelope{
        Type:       event.EventType,
        Category:   "inventory",
        BusinessID: event.BusinessID,
        Data:       event.Data,
        Metadata: map[string]interface{}{
            "order_id":     event.OrderID,
            "warehouse_id": event.WarehouseID,
        },
    })
}()
```

**Also in:** `back/central/services/modules/invoicing/internal/infra/secondary/queue/event_publisher.go` (lines 42-55, 65-78, 87-102, 111-126 — all four `PublishInvoice*` methods)

**Problem:** If RabbitMQ publish fails, the error is completely unobserved. Events are silently lost.

**Fix:**
```go
go func() {
    if err := rabbitmq.PublishEvent(context.Background(), p.queue, envelope); err != nil {
        p.logger.Error(context.Background()).
            Err(err).
            Str("event_type", event.EventType).
            Msg("Failed to publish event to RabbitMQ")
    }
}()
```

---

### Anti-Pattern 2: context.Background() in Request-Scoped Goroutine

**File:** `back/central/services/modules/orders/internal/app/usecasecreateorder/publish_events.go:86-94`

```go
// Uses context.Background() — loses request tracing
go func() {
    if err := uc.rabbitEventPublisher.PublishOrderEvent(context.Background(), event, order); err != nil {
        uc.logger.Error(context.Background()).
            Err(err).
            Msg("Error al publicar evento a RabbitMQ")
    }
}()
```

**Also at lines 108-123** (score calculation with explicit `bgCtx := context.Background()`).

**Also in:** `back/central/services/modules/shipments/internal/infra/secondary/queue/sse_publisher.go:94-108`

**Nuance:** For fire-and-forget background publishing, `context.Background()` is sometimes intentional — the publish should survive request cancellation. But tracing context (business_id, user_id) is lost. The recommended approach is to create a detached context that carries tracing values but no cancellation:

```go
// Intentional: detach from request lifecycle but keep trace IDs
bgCtx := context.WithoutCancel(ctx) // Go 1.21+
go func() {
    if err := uc.rabbitEventPublisher.PublishOrderEvent(bgCtx, event, order); err != nil {
        uc.logger.Error(bgCtx).Err(err).Msg("Failed to publish event")
    }
}()
```

---

### Anti-Pattern 3: String Comparison Instead of errors.Is()

**File:** `back/central/services/modules/orders/internal/infra/primary/handlers/get-order.go:37-46`

```go
// BUG: fragile string comparison — breaks if error message changes
if err.Error() == "order not found" {
    c.JSON(http.StatusNotFound, gin.H{
        "success": false,
        "message": "Orden no encontrada",
        "error":   err.Error(),
    })
    return
}
```

**Problem:** If the domain error message changes even slightly, this comparison breaks silently and falls through to 500.

**Fix:**
```go
if errors.Is(err, domainerrors.ErrOrderNotFound) {
    c.JSON(http.StatusNotFound, gin.H{
        "error":   "ORDER_NOT_FOUND",
        "message": "Orden no encontrada",
    })
    return
}
```

---

### Anti-Pattern 4: Leaking Raw Error to Frontend

**File:** `back/central/services/modules/customers/internal/infra/primary/handlers/create_customer.go:36-43`

```go
// BUG: err.Error() leaks internal English error message to the user
if errors.Is(err, domainerrors.ErrDuplicateEmail) || errors.Is(err, domainerrors.ErrDuplicateDni) {
    c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
    return
}
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
```

**Problems:**
1. `err.Error()` returns the internal English message — the user sees "a client with this email already exists"
2. On 500, the raw DB/Go error is exposed — potential information leak
3. No machine-readable error code for frontend conditional logic

**Fix:**
```go
switch {
case errors.Is(err, domainerrors.ErrDuplicateEmail):
    c.JSON(http.StatusConflict, gin.H{
        "error":   "DUPLICATE_EMAIL",
        "message": "Ya existe un cliente con este correo electronico",
    })
    return
case errors.Is(err, domainerrors.ErrDuplicateDni):
    c.JSON(http.StatusConflict, gin.H{
        "error":   "DUPLICATE_DNI",
        "message": "Ya existe un cliente con este documento de identidad",
    })
    return
default:
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "INTERNAL_ERROR",
        "message": "Error interno del servidor",
    })
    return
}
```

---

### Anti-Pattern 5: Inconsistent Error Response Formats

Found across the codebase — three different JSON shapes for errors:

```json
// Format A (orders): success flag + message + error
{"success": false, "message": "Orden no encontrada", "error": "order not found"}

// Format B (invoicing): error code + message (CORRECT)
{"error": "INVOICE_NOT_FOUND", "message": "Factura no encontrada"}

// Format C (customers): just error string
{"error": "client not found"}
```

**Problem:** Frontend has to handle 3 different error shapes. Can't build a generic error handler.

**Recommended standard (Format B):**
```json
{"error": "ERROR_CODE", "message": "Mensaje descriptivo en espanol"}
```

---

### Anti-Pattern 6: Use Case Returns Error Without Logging

```go
// BUG: error returned but never logged — invisible in production
func (uc *UseCase) GetItem(ctx context.Context, id uint) (*entities.Item, error) {
    item, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err  // no log — if this fails, nobody knows
    }
    return item, nil
}
```

**Fix:**
```go
func (uc *UseCase) GetItem(ctx context.Context, id uint) (*entities.Item, error) {
    item, err := uc.repo.GetByID(ctx, id)
    if err != nil {
        uc.log.Error(ctx).Err(err).Uint("id", id).Msg("Error al obtener item")
        return nil, err
    }
    return item, nil
}
```

---

## Quick Reference: What to Flag

| Pattern | Severity | Category |
|---------|----------|----------|
| `go func()` with infinite loop, no `ctx.Done()` | CRITICAL | Goroutine Lifecycle |
| `_ = fn()` in goroutine (network/DB/IO) | CRITICAL | Error Handling |
| `Lock()` without `defer Unlock()` | CRITICAL | Mutex & Sync |
| Mutex passed by value | CRITICAL | Mutex & Sync |
| External API call inside DB transaction | CRITICAL | Transactions |
| Missing error return inside `db.Transaction` closure | CRITICAL | Transactions |
| Channel close from receiver side | CRITICAL | Channels |
| Resource created in loop without per-iteration cleanup | CRITICAL | Resource Cleanup |
| `defer` inside a loop body | CRITICAL | Resource Cleanup |
| `context.Background()` in goroutine without comment | WARNING | Context Propagation |
| Fire-and-forget goroutine with no error logging | WARNING | Goroutine Lifecycle |
| `sync.Mutex` for read-heavy access | WARNING | Mutex & Sync |
| Channel send without `select`/`default` | WARNING | Channels |
| Transaction scope too wide | WARNING | Transactions |
| Unjustified channel buffer size | INFO | Channels |
| Intentional `context.Background()` with error logging | INFO | Context Propagation |
| Domain error not checked in handler (falls to 500) | CRITICAL | Error Handling |
| `err.Error()` exposed directly in HTTP response | CRITICAL | Error Handling |
| `gorm.ErrRecordNotFound` leaks to handler unmapped | CRITICAL | Error Handling |
| Error returned without any logging in the chain | CRITICAL | Error Handling |
| `err.Error() == "string"` instead of `errors.Is()` | WARNING | Error Handling |
| Error message in English in HTTP response | WARNING | Error Handling |
| Wrong HTTP status code for error type | WARNING | Error Handling |
| No centralized `error_handler.go` (inline scattered) | WARNING | Error Handling |
| Inconsistent error JSON format across handlers | WARNING | Error Handling |
| Generic `fmt.Errorf()` instead of domain sentinel | WARNING | Error Handling |
| Missing machine-readable error code in response | INFO | Error Handling |
