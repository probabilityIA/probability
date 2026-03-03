package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// InvoiceRepository implementa IInvoiceRepository

// NewInvoiceRepository crea un nuevo repositorio de facturas

// Create crea una nueva factura
func (r *Repository) CreateInvoice(ctx context.Context, invoice *entities.Invoice) error {
	model := mappers.InvoiceToModel(invoice)

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Msg("Failed to create invoice")
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	// Actualizar el ID de la entidad
	invoice.ID = model.ID
	invoice.InternalNumber = model.InternalNumber // Generado por BeforeCreate

	return nil
}

// GetByID obtiene una factura por ID
func (r *Repository) GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderID obtiene una factura por ID de orden
func (r *Repository) GetInvoiceByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).Where("order_id = ?", orderID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found for order: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderAndProvider obtiene una factura por orden y proveedor
func (r *Repository) GetInvoiceByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).
		Where("order_id = ? AND invoicing_provider_id = ?", orderID, providerID).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// ListInvoices lista facturas con filtros y paginación
func (r *Repository) ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error) {
	var modelsList []*models.Invoice
	var total int64

	query := r.db.Conn(ctx).Model(&models.Invoice{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok {
		query = query.Where("business_id = ?", businessID)
	}

	if orderID, ok := filters["order_id"].(string); ok {
		query = query.Where("order_id = ?", orderID)
	}

	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}

	if providerID, ok := filters["invoicing_provider_id"].(uint); ok {
		query = query.Where("invoicing_provider_id = ?", providerID)
	}

	if invoiceNumber, ok := filters["invoice_number"].(string); ok && invoiceNumber != "" {
		query = query.Where("invoice_number ILIKE ?", "%"+invoiceNumber+"%")
	}

	if customerName, ok := filters["customer_name"].(string); ok && customerName != "" {
		query = query.Where("customer_name ILIKE ?", "%"+customerName+"%")
	}

	if orderNumber, ok := filters["order_number"].(string); ok && orderNumber != "" {
		query = query.Where("order_id IN (SELECT id FROM orders WHERE (order_number ILIKE ? OR external_id ILIKE ?) AND deleted_at IS NULL)",
			"%"+orderNumber+"%", "%"+orderNumber+"%")
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		query = query.Where("invoices.created_at >= ?", startDate)
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		query = query.Where("invoices.created_at < ?::date + interval '1 day'", endDate)
	}

	// Contar total antes de paginar
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count invoices: %w", err)
	}

	// Paginación
	page := 1
	pageSize := 20

	if p, ok := filters["page"].(int); ok && p > 0 {
		page = p
	}
	if ps, ok := filters["page_size"].(int); ok && ps > 0 {
		pageSize = ps
	}

	offset := (page - 1) * pageSize

	// Ordenar y paginar (con preload del logo del proveedor)
	if err := query.
		Preload("InvoicingIntegration.IntegrationType").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&modelsList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	return mappers.InvoiceListToDomain(modelsList), total, nil
}

// Update actualiza una factura
func (r *Repository) UpdateInvoice(ctx context.Context, invoice *entities.Invoice) error {
	model := mappers.InvoiceToModel(invoice)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to update invoice")
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}

// Delete elimina una factura (soft delete)
func (r *Repository) DeleteInvoice(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.Invoice{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	return nil
}

// GetIssuedInvoicesByDateRange retorna facturas con status "issued" de un negocio en un rango de fechas.
// Usado para comparación en memoria contra el proveedor (sin persistencia adicional).
// dateFrom y dateTo deben estar en formato YYYY-MM-DD.
func (r *Repository) GetIssuedInvoicesByDateRange(ctx context.Context, businessID uint, dateFrom, dateTo string) ([]*entities.Invoice, error) {
	var modelsList []*models.Invoice

	if err := r.db.Conn(ctx).
		Model(&models.Invoice{}).
		Preload("InvoiceItems").
		Where("business_id = ?", businessID).
		Where("status = ?", "issued").
		Where("DATE(created_at AT TIME ZONE 'America/Bogota') BETWEEN ? AND ?", dateFrom, dateTo).
		Order("created_at DESC").
		Find(&modelsList).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Str("date_from", dateFrom).
			Str("date_to", dateTo).
			Msg("Failed to get issued invoices by date range")
		return nil, fmt.Errorf("failed to get issued invoices: %w", err)
	}

	return mappers.InvoiceListToDomain(modelsList), nil
}

// GetOrderCreatedAtsByIDs retorna map[orderID]createdAt para un batch de órdenes.
// Replicado localmente para evitar compartir repositorios entre módulos (módulo orders gestiona tabla orders).
func (r *Repository) GetOrderCreatedAtsByIDs(ctx context.Context, orderIDs []string) (map[string]*time.Time, error) {
	if len(orderIDs) == 0 {
		return map[string]*time.Time{}, nil
	}

	var rows []struct {
		ID        string
		CreatedAt time.Time
	}

	if err := r.db.Conn(ctx).
		Table("orders").
		Select("id, created_at").
		Where("id IN (?)", orderIDs).
		Where("deleted_at IS NULL").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("failed to get order created_at: %w", err)
	}

	result := make(map[string]*time.Time, len(rows))
	for _, row := range rows {
		t := row.CreatedAt
		result[row.ID] = &t
	}
	return result, nil
}

// InvoiceExistsForOrder verifica si existe una factura VÁLIDA para una orden e integración
// Solo considera facturas con status pending, issued o draft (excluye failed y cancelled)
func (r *Repository) InvoiceExistsForOrder(ctx context.Context, orderID string, integrationID uint) (bool, error) {
	var count int64

	// Solo contar facturas válidas (no fallidas ni canceladas)
	// Las facturas con status "failed" o "cancelled" NO bloquean crear una nueva factura
	if err := r.db.Conn(ctx).Model(&models.Invoice{}).
		Where("order_id = ? AND (invoicing_integration_id = ? OR invoicing_provider_id = ?)", orderID, integrationID, integrationID).
		Where("status NOT IN (?)", []string{"failed", "cancelled"}).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check invoice existence: %w", err)
	}

	return count > 0, nil
}
