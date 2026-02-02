# Tests Unitarios - Casos de Uso de Orders

## Resumen

Tests unitarios completos para los casos de uso del módulo de orders siguiendo arquitectura hexagonal y mejores prácticas de Go.

### Cobertura

**Cobertura actual: 60.6%** de statements

### Archivos de Test

- `create-order_test.go` - Tests de creación de órdenes
- `get-order_test.go` - Tests de obtención de órdenes por ID
- `list-orders_test.go` - Tests de listado paginado de órdenes
- `update-order_test.go` - Tests de actualización de órdenes

### Total de Tests

**32 tests ejecutándose exitosamente**

---

## Estructura de Mocks

Los mocks están ubicados en `/internal/mocks/` siguiendo la regla de aislamiento de dependencias:

### Mocks Disponibles

1. **RepositoryMock** (`repository_mock.go`)
   - Mockea todas las operaciones del repositorio de órdenes
   - Incluye métodos CRUD, validación, catálogo y consultas a estados

2. **EventPublisherMock** (`event_publisher_mock.go`)
   - Mockea publicación de eventos a Redis

3. **RabbitPublisherMock** (`rabbit_publisher_mock.go`)
   - Mockea publicación de eventos a RabbitMQ

4. **LoggerMock** (`logger_mock.go`)
   - Mockea el logger estructurado (zerolog)

5. **ScoreUseCaseMock** (`score_usecase_mock.go`)
   - Mockea el caso de uso de cálculo de score

---

## Tests de CreateOrder

### Casos Cubiertos

| Test | Descripción | Validaciones |
|------|-------------|--------------|
| `TestCreateOrder_Success` | Creación exitosa de orden | - Orden no existe previamente<br>- Se guarda correctamente<br>- Se generan IDs<br>- Se retorna response completo |
| `TestCreateOrder_OrderAlreadyExists` | Orden duplicada por ExternalID | - Detecta orden existente<br>- Retorna error apropiado<br>- NO llama CreateOrder |
| `TestCreateOrder_OrderExistsCheckError` | Error en validación de existencia | - Maneja error de BD<br>- Propaga error correctamente |
| `TestCreateOrder_RepositoryError` | Error al guardar en BD | - Maneja error de inserción<br>- Retorna error envuelto |
| `TestCreateOrder_ValidatesRequiredFields` | Mapeo correcto de campos | - Todos los campos se mapean<br>- IDs, montos, cliente, etc. |

### Ejemplo de Uso

```go
func TestCreateOrder_Success(t *testing.T) {
    // Arrange - Configurar mocks
    mockRepo := new(mocks.RepositoryMock)
    mockRedisPublisher := new(mocks.EventPublisherMock)
    // ...

    // Configurar expectativas
    mockRepo.On("OrderExists", ctx, req.ExternalID, req.IntegrationID).
        Return(false, nil)
    mockRepo.On("CreateOrder", ctx, mock.AnythingOfType("*entities.ProbabilityOrder")).
        Return(nil)

    // Act - Ejecutar
    result, err := useCase.CreateOrder(ctx, req)

    // Assert - Verificar
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

---

## Tests de GetOrderByID

### Casos Cubiertos

| Test | Descripción | Validaciones |
|------|-------------|--------------|
| `TestGetOrderByID_Success` | Obtención exitosa por ID | - Retorna orden completa<br>- Mapeo correcto a response |
| `TestGetOrderByID_EmptyID` | ID vacío | - Valida ID requerido<br>- NO llama al repositorio |
| `TestGetOrderByID_NotFound` | Orden no encontrada | - Maneja error not found<br>- Retorna nil |
| `TestGetOrderByID_DatabaseError` | Error de conexión BD | - Propaga error de BD |
| `TestGetOrderByID_WithCompleteData` | Orden con todos los campos | - Mapea campos opcionales<br>- Incluye tracking, pago, etc. |

---

## Tests de ListOrders

### Casos Cubiertos

| Test | Descripción | Validaciones |
|------|-------------|--------------|
| `TestListOrders_Success` | Listado paginado exitoso | - Retorna órdenes<br>- Metadata de paginación correcta<br>- Cálculo de páginas |
| `TestListOrders_EmptyResult` | Sin resultados | - Maneja lista vacía<br>- Total = 0, páginas = 0 |
| `TestListOrders_PaginationValidation` | Validación de parámetros | **6 sub-tests**:<br>- Página negativa → 1<br>- Página 0 → 1<br>- PageSize negativo → 10<br>- PageSize 0 → 10<br>- PageSize > 100 → 10<br>- Valores válidos no cambian |
| `TestListOrders_RepositoryError` | Error de BD | - Maneja timeout/error |
| `TestListOrders_TotalPagesCalculation` | Cálculo de páginas | **5 sub-tests**:<br>- Sin registros<br>- Menos que pageSize<br>- Exacto una página<br>- Múltiples páginas<br>- Exacto múltiples |
| `TestListOrders_WithFilters` | Filtros complejos | - Pasa filtros al repo<br>- Múltiples criterios |

### Validación de Paginación

Los tests validan que se cumplen las reglas de paginación:

```go
// Límites aplicados
page < 1        → page = 1
pageSize < 1    → pageSize = 10
pageSize > 100  → pageSize = 10
```

---

## Tests de UpdateOrder

### Casos Cubiertos

| Test | Descripción | Validaciones |
|------|-------------|--------------|
| `TestUpdateOrder_Success` | Actualización exitosa | - Solo actualiza campos enviados<br>- Recalcula score<br>- Publica eventos |
| `TestUpdateOrder_EmptyID` | ID vacío | - Valida ID requerido |
| `TestUpdateOrder_OrderNotFound` | Orden no existe | - Maneja not found |
| `TestUpdateOrder_StatusChange_PublishesStatusEvent` | Cambio de estado | - Publica 2 eventos:<br>&nbsp;&nbsp;1. order.updated<br>&nbsp;&nbsp;2. order.status_changed |
| `TestUpdateOrder_PartialUpdate` | Update parcial | - Solo modifica campos enviados<br>- Preserva campos no enviados |
| `TestUpdateOrder_RepositoryError` | Error al guardar | - Maneja constraint violation |
| `TestUpdateOrder_ConfirmationStatus` | Estados de confirmación | **3 sub-tests**:<br>- "yes" → isConfirmed = true<br>- "no" → isConfirmed = false<br>- "pending" → isConfirmed = nil |

### Lógica Especial Testeada

**Recálculo de Score:**
- Se llama después de cada update
- No bloqueante (no falla si score falla)
- Recarga orden si score es exitoso

**Eventos Duales:**
- Publica a Redis (pub/sub)
- Publica a RabbitMQ (queues)

---

## Ejecutar Tests

### Todos los tests

```bash
cd /home/cam/Desktop/probability/back/central/services/modules/orders
go test ./internal/app/usecaseorder/... -v
```

### Con cobertura

```bash
go test ./internal/app/usecaseorder -cover
# Output: coverage: 60.6% of statements
```

### Cobertura detallada

```bash
go test ./internal/app/usecaseorder -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test específico

```bash
go test ./internal/app/usecaseorder -run TestCreateOrder_Success -v
```

### Tests por patrón

```bash
# Todos los tests de CreateOrder
go test ./internal/app/usecaseorder -run "TestCreateOrder" -v

# Todos los tests de validación
go test ./internal/app/usecaseorder -run "Validation" -v
```

---

## Convenciones Usadas

### Nomenclatura

```
Test{UseCase}_{Method}_{Scenario}

Ejemplos:
- TestCreateOrder_Success
- TestListOrders_EmptyResult
- TestUpdateOrder_PartialUpdate
```

### Patrón AAA

Todos los tests siguen el patrón **Arrange - Act - Assert**:

```go
func TestExample(t *testing.T) {
    // Arrange - Configurar mocks y datos
    mockRepo := new(mocks.RepositoryMock)
    mockRepo.On("Method", args).Return(result, nil)

    // Act - Ejecutar el método a testear
    result, err := useCase.Method(ctx, input)

    // Assert - Verificar resultados
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    mockRepo.AssertExpectations(t)
}
```

### Table-Driven Tests

Para casos con múltiples escenarios similares:

```go
tests := []struct {
    name     string
    input    int
    expected int
}{
    {"caso 1", 1, 10},
    {"caso 2", 2, 20},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

---

## Mejores Prácticas Aplicadas

1. **Independencia**: Cada test es independiente y no depende de otros
2. **Mocks Limpios**: Uso de testify/mock para expectativas claras
3. **Errores Testeados**: Tanto casos felices como casos de error
4. **Nombres Descriptivos**: Los nombres de test documentan el comportamiento
5. **Sin BD Real**: Todos los tests son unitarios, no hay integración con BD
6. **Fast**: Los 32 tests se ejecutan en ~20ms

---

## Próximos Pasos (Sugerencias)

### Tests Pendientes

- `delete-order_test.go` - Tests de eliminación
- `get-order-raw_test.go` - Tests de obtención de metadata cruda
- `request-confirmation_test.go` - Tests de solicitud de confirmación

### Tests de Handlers

Generar tests para la capa de infraestructura HTTP:
- `internal/infra/primary/handlers/*_test.go`

### Tests de Integración

Tests que usen BD real:
- `internal/infra/secondary/repository/*_integration_test.go`

### Aumentar Cobertura

Áreas para mejorar cobertura (actualmente 60.6%):
- Validaciones complejas en DTOs
- Lógica de mappers
- Edge cases en helpers

---

**Generado**: 2026-02-02
**Módulo**: orders
**Framework**: Go 1.23 + testify/mock
