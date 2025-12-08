package domain

import (
	"time"
)

// ───────────────────────────────────────────
//
//	EVENT TYPES
//
// ───────────────────────────────────────────

// EventType define los tipos de eventos que pueden ocurrir
type EventType string

const (
	EventTypeInventorySyncStarted   EventType = "inventory_sync_started"
	EventTypeInventorySyncProgress  EventType = "inventory_sync_progress"
	EventTypeInventorySyncCompleted EventType = "inventory_sync_completed"
	EventTypeInventorySyncFailed    EventType = "inventory_sync_failed"
	EventTypeProductSynced          EventType = "product_synced"
	EventTypeProductFailed          EventType = "product_failed"
	EventTypeConnectionEstablished  EventType = "connection_established"
	EventTypeBatchStarted           EventType = "batch_started"
	EventTypeBatchCompleted         EventType = "batch_completed"
)

// ───────────────────────────────────────────
//
//	SYNC STATUS
//
// ───────────────────────────────────────────

// SyncStatus define los estados de sincronización
type SyncStatus string

const (
	SyncStatusStarted    SyncStatus = "started"
	SyncStatusInProgress SyncStatus = "in_progress"
	SyncStatusCompleted  SyncStatus = "completed"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusPaused     SyncStatus = "paused"
)

// ───────────────────────────────────────────
//
//	EVENT STRUCTURES
//
// ───────────────────────────────────────────

// Event es la estructura base para todos los eventos
type Event struct {
	ID            string
	Type          EventType
	IntegrationID int64
	BusinessID    string
	Timestamp     time.Time
	Data          interface{}
	Metadata      map[string]interface{}
}

// InventorySyncEvent representa un evento de sincronización de inventario
type InventorySyncEvent struct {
	Status   SyncStatus
	Progress SyncProgress
	Error    string
	Details  SyncDetails
}

// ProductSyncEvent representa un evento de sincronización de producto
type ProductSyncEvent struct {
	ProductID string
	Status    string
	Error     string
	SyncedAt  time.Time
	Details   ProductSyncInfo
}

// ConnectionEvent representa un evento de conexión establecida
type ConnectionEvent struct {
	Message string
}

// BatchEvent representa un evento de batch
type BatchEvent struct {
	BatchNumber     int
	TotalBatches    int
	ProductsInBatch int
	Status          string
}

// ───────────────────────────────────────────
//
//	SYNC STRUCTURES
//
// ───────────────────────────────────────────

// SyncProgress representa el progreso de una sincronización
type SyncProgress struct {
	TotalProducts   int
	SyncedProducts  int
	FailedProducts  int
	CurrentProduct  int
	Percentage      float64
	CurrentBatch    int
	TotalBatches    int
	ProcessingSpeed float64 // productos por segundo
}

// SyncDetails contiene los detalles de una sincronización
type SyncDetails struct {
	IntegrationType string
	WarehouseID     *int
	BatchSize       int
	ParallelWorkers int
	StartTime       time.Time
	EstimatedEnd    *time.Time
	CurrentProduct  *ProductSyncInfo
	LastSynced      *ProductSyncInfo
	FailedProducts  []ProductSyncInfo
}

// ProductSyncInfo contiene información sobre un producto sincronizado
type ProductSyncInfo struct {
	ProductID   string
	SKU         string
	ExternalID  string
	Quantity    int
	WarehouseID *int
	Status      string
	Error       string
	SyncedAt    time.Time
	Duration    float64 // duración en milisegundos
}

// SyncSession representa una sesión de sincronización de inventario
type SyncSession struct {
	IntegrationID   int64
	BusinessID      string
	TotalProducts   int
	IntegrationType string
	WarehouseID     *int
	BatchSize       int
	TotalBatches    int
	CurrentBatch    int
	ParallelWorkers int
	CurrentProduct  string
	LastSynced      *ProductSyncInfo
	FailedProducts  []ProductSyncInfo
	ProcessingSpeed float64
	LastSpeedUpdate time.Time
	StartTime       time.Time
	EndTime         *time.Time
	Status          SyncStatus
	Error           string
}

// NewSyncSession crea una nueva sesión de sincronización
func NewSyncSession(integrationID int64, businessID string, totalProducts int, integrationType string, warehouseID *int, batchSize int, parallelWorkers int) *SyncSession {
	now := time.Now()

	return &SyncSession{
		IntegrationID:   integrationID,
		BusinessID:      businessID,
		TotalProducts:   totalProducts,
		IntegrationType: integrationType,
		WarehouseID:     warehouseID,
		BatchSize:       batchSize,
		TotalBatches:    (totalProducts + batchSize - 1) / batchSize, // Redondear hacia arriba
		CurrentBatch:    0,
		ParallelWorkers: parallelWorkers,
		CurrentProduct:  "",
		LastSynced:      nil,
		FailedProducts:  make([]ProductSyncInfo, 0),
		ProcessingSpeed: 0.0,
		LastSpeedUpdate: now,
		StartTime:       now,
		EndTime:         nil,
		Status:          SyncStatusStarted,
		Error:           "",
	}
}

// GetProgress retorna el progreso actual de la sincronización
func (s *SyncSession) GetProgress() SyncProgress {
	syncedCount := 0
	failedCount := len(s.FailedProducts)

	if s.LastSynced != nil {
		syncedCount = 1 // Por simplicidad, podríamos mantener un contador real
	}

	percentage := 0.0
	if s.TotalProducts > 0 {
		percentage = float64(syncedCount+failedCount) / float64(s.TotalProducts) * 100.0
	}

	return SyncProgress{
		TotalProducts:   s.TotalProducts,
		SyncedProducts:  syncedCount,
		FailedProducts:  failedCount,
		CurrentProduct:  0, // Cambiar a string si es necesario
		Percentage:      percentage,
		CurrentBatch:    s.CurrentBatch,
		TotalBatches:    s.TotalBatches,
		ProcessingSpeed: s.ProcessingSpeed,
	}
}

// GetDetails retorna los detalles de la sesión
func (s *SyncSession) GetDetails() SyncDetails {
	return SyncDetails{
		IntegrationType: s.IntegrationType,
		WarehouseID:     s.WarehouseID,
		BatchSize:       s.BatchSize,
		ParallelWorkers: s.ParallelWorkers,
		StartTime:       s.StartTime,
		EstimatedEnd:    s.getEstimatedEndTime(),
		CurrentProduct:  s.LastSynced,
		LastSynced:      s.LastSynced,
		FailedProducts:  s.FailedProducts,
	}
}

// SetCurrentProduct establece el producto actualmente siendo procesado
func (s *SyncSession) SetCurrentProduct(productID, sku, externalID string, warehouseID *int, quantity int) {
	s.CurrentProduct = productID
	s.LastSynced = &ProductSyncInfo{
		ProductID:   productID,
		SKU:         sku,
		ExternalID:  externalID,
		Quantity:    quantity,
		WarehouseID: warehouseID,
		Status:      string(SyncStatusInProgress),
		SyncedAt:    time.Now(),
	}
}

// ProductSynced marca un producto como sincronizado exitosamente
func (s *SyncSession) ProductSynced(productID, sku, externalID string, quantity int, warehouseID *int, duration time.Duration) {
	s.LastSynced = &ProductSyncInfo{
		ProductID:   productID,
		SKU:         sku,
		ExternalID:  externalID,
		Quantity:    quantity,
		WarehouseID: warehouseID,
		Status:      string(SyncStatusCompleted),
		Duration:    float64(duration.Milliseconds()),
		SyncedAt:    time.Now(),
	}

	s.updateProcessingSpeed()
}

// ProductFailed marca un producto como fallido en la sincronización
func (s *SyncSession) ProductFailed(productID, sku, externalID string, warehouseID *int, errorMsg string, duration time.Duration) {
	failedProduct := ProductSyncInfo{
		ProductID:   productID,
		SKU:         sku,
		ExternalID:  externalID,
		WarehouseID: warehouseID,
		Status:      string(SyncStatusFailed),
		Error:       errorMsg,
		Duration:    float64(duration.Milliseconds()),
		SyncedAt:    time.Now(),
	}

	s.FailedProducts = append(s.FailedProducts, failedProduct)
	s.updateProcessingSpeed()
}

// StartBatch inicia un nuevo lote de sincronización
func (s *SyncSession) StartBatch(batchNumber, productsInBatch int) {
	s.CurrentBatch = batchNumber
	s.Status = SyncStatusInProgress
}

// CompleteBatch completa un lote de sincronización
func (s *SyncSession) CompleteBatch(batchNumber int) {
	if s.CurrentBatch == s.TotalBatches {
		s.Status = SyncStatusCompleted
		s.EndTime = &[]time.Time{time.Now()}[0]
	}
}

// updateProcessingSpeed actualiza la velocidad de procesamiento
func (s *SyncSession) updateProcessingSpeed() {
	now := time.Now()
	elapsed := now.Sub(s.LastSpeedUpdate).Seconds()

	if elapsed > 0 {
		// Calcular productos por segundo en el último período
		recentProducts := 1 // Por simplicidad
		s.ProcessingSpeed = float64(recentProducts) / elapsed
		s.LastSpeedUpdate = now
	}
}

// getEstimatedEndTime calcula el tiempo estimado de finalización
func (s *SyncSession) getEstimatedEndTime() *time.Time {
	if s.ProcessingSpeed <= 0 {
		return nil
	}

	remainingProducts := s.TotalProducts - len(s.FailedProducts)
	if s.LastSynced != nil {
		remainingProducts--
	}

	if remainingProducts <= 0 {
		return nil
	}

	estimatedSeconds := float64(remainingProducts) / s.ProcessingSpeed
	estimatedEnd := s.StartTime.Add(time.Duration(estimatedSeconds) * time.Second)

	return &estimatedEnd
}

// IsCompleted verifica si la sincronización está completada
func (s *SyncSession) IsCompleted() bool {
	return s.Status == SyncStatusCompleted || s.Status == SyncStatusFailed
}

// GetSuccessRate retorna la tasa de éxito de la sincronización
func (s *SyncSession) GetSuccessRate() float64 {
	if s.TotalProducts == 0 {
		return 0.0
	}

	successCount := s.TotalProducts - len(s.FailedProducts)
	return float64(successCount) / float64(s.TotalProducts) * 100.0
}

// ───────────────────────────────────────────
//
//	REQUEST/RESPONSE DTOs
//
// ───────────────────────────────────────────

// ConnectionRequest representa una solicitud de conexión
type ConnectionRequest struct {
	IntegrationID int64
	Connection    interface{} // http.ResponseWriter
}

// ConnectionInfo representa información de conexiones
type ConnectionInfo struct {
	IntegrationID     int64
	ActiveConnections int
	TotalConnections  int
	ConnectedAt       time.Time
	LastActivity      time.Time
	Status            string
}

// PublishEventResult representa el resultado de publicar un evento
type PublishEventResult struct {
	Recipients int
	Errors     []string
}

// ConnectionResult representa el resultado de una operación de conexión
type ConnectionResult struct {
	IntegrationID int64
	Operation     string
	Message       string
	Errors        []string
}

// SSEConnectionRequest representa una solicitud de conexión SSE (legacy - mantener para compatibilidad)
type SSEConnectionRequest struct {
	IntegrationID int64
	BusinessID    string
}

// SSEStatus representa el estado de las conexiones SSE
type SSEStatus struct {
	IntegrationID     int64
	ActiveConnections int
	LastActivity      time.Time
	Status            string
}
