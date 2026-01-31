package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// InvoiceRepository implementa IInvoiceRepository
type InvoiceRepository struct {
	*Repository
}

// NewInvoiceRepository crea un nuevo repositorio de facturas
func NewInvoiceRepository(repo *Repository) ports.IInvoiceRepository {
	return &InvoiceRepository{Repository: repo}
}

// Create crea una nueva factura
func (r *InvoiceRepository) Create(ctx context.Context, invoice *entities.Invoice) error {
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
func (r *InvoiceRepository) GetByID(ctx context.Context, id uint) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderID obtiene una factura por ID de orden
func (r *InvoiceRepository) GetByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).Where("order_id = ?", orderID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found for order: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderAndProvider obtiene una factura por orden y proveedor
func (r *InvoiceRepository) GetByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error) {
	var model models.Invoice

	if err := r.db.Conn(ctx).
		Where("order_id = ? AND invoicing_provider_id = ?", orderID, providerID).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// List lista facturas con filtros
func (r *InvoiceRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error) {
	var modelsList []*models.Invoice

	query := r.db.Conn(ctx).Model(&models.Invoice{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok {
		query = query.Where("business_id = ?", businessID)
	}

	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}

	if providerID, ok := filters["invoicing_provider_id"].(uint); ok {
		query = query.Where("invoicing_provider_id = ?", providerID)
	}

	// Ordenar por más reciente
	query = query.Order("created_at DESC")

	// Aplicar limit si existe
	if limit, ok := filters["limit"].(int); ok {
		query = query.Limit(limit)
	}

	if err := query.Find(&modelsList).Error; err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}

	return mappers.InvoiceListToDomain(modelsList), nil
}

// Update actualiza una factura
func (r *InvoiceRepository) Update(ctx context.Context, invoice *entities.Invoice) error {
	model := mappers.InvoiceToModel(invoice)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to update invoice")
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}

// Delete elimina una factura (soft delete)
func (r *InvoiceRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.Conn(ctx).Delete(&models.Invoice{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	return nil
}

// ExistsForOrder verifica si existe una factura para una orden y proveedor
func (r *InvoiceRepository) ExistsForOrder(ctx context.Context, orderID string, providerID uint) (bool, error) {
	var count int64

	if err := r.db.Conn(ctx).Model(&models.Invoice{}).
		Where("order_id = ? AND invoicing_provider_id = ?", orderID, providerID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check invoice existence: %w", err)
	}

	return count > 0, nil
}
