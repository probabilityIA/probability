package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
<<<<<<< HEAD
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/ports"
=======
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// InvoiceRepository implementa IInvoiceRepository
<<<<<<< HEAD
type InvoiceRepository struct {
	*Repository
}

// NewInvoiceRepository crea un nuevo repositorio de facturas
func NewInvoiceRepository(repo *Repository) ports.IInvoiceRepository {
	return &InvoiceRepository{Repository: repo}
}

// Create crea una nueva factura
func (r *InvoiceRepository) Create(ctx context.Context, invoice *entities.Invoice) error {
=======

// NewInvoiceRepository crea un nuevo repositorio de facturas

// Create crea una nueva factura
func (r *Repository) CreateInvoice(ctx context.Context, invoice *entities.Invoice) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
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
<<<<<<< HEAD
func (r *InvoiceRepository) GetByID(ctx context.Context, id uint) (*entities.Invoice, error) {
=======
func (r *Repository) GetInvoiceByID(ctx context.Context, id uint) (*entities.Invoice, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.Invoice

	if err := r.db.Conn(ctx).First(&model, id).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderID obtiene una factura por ID de orden
<<<<<<< HEAD
func (r *InvoiceRepository) GetByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error) {
=======
func (r *Repository) GetInvoiceByOrderID(ctx context.Context, orderID string) (*entities.Invoice, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.Invoice

	if err := r.db.Conn(ctx).Where("order_id = ?", orderID).First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found for order: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

// GetByOrderAndProvider obtiene una factura por orden y proveedor
<<<<<<< HEAD
func (r *InvoiceRepository) GetByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error) {
=======
func (r *Repository) GetInvoiceByOrderAndProvider(ctx context.Context, orderID string, providerID uint) (*entities.Invoice, error) {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	var model models.Invoice

	if err := r.db.Conn(ctx).
		Where("order_id = ? AND invoicing_provider_id = ?", orderID, providerID).
		First(&model).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	return mappers.InvoiceToDomain(&model), nil
}

<<<<<<< HEAD
// List lista facturas con filtros
func (r *InvoiceRepository) List(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, error) {
	var modelsList []*models.Invoice
=======
// ListInvoices lista facturas con filtros y paginación
func (r *Repository) ListInvoices(ctx context.Context, filters map[string]interface{}) ([]*entities.Invoice, int64, error) {
	var modelsList []*models.Invoice
	var total int64
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e

	query := r.db.Conn(ctx).Model(&models.Invoice{})

	// Aplicar filtros
	if businessID, ok := filters["business_id"].(uint); ok {
		query = query.Where("business_id = ?", businessID)
	}

<<<<<<< HEAD
=======
	if orderID, ok := filters["order_id"].(string); ok {
		query = query.Where("order_id = ?", orderID)
	}

>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if status, ok := filters["status"].(string); ok {
		query = query.Where("status = ?", status)
	}

	if providerID, ok := filters["invoicing_provider_id"].(uint); ok {
		query = query.Where("invoicing_provider_id = ?", providerID)
	}

<<<<<<< HEAD
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
=======
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

	// Ordenar y paginar
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&modelsList).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	return mappers.InvoiceListToDomain(modelsList), total, nil
}

// Update actualiza una factura
func (r *Repository) UpdateInvoice(ctx context.Context, invoice *entities.Invoice) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	model := mappers.InvoiceToModel(invoice)

	if err := r.db.Conn(ctx).Save(model).Error; err != nil {
		r.log.Error(ctx).Err(err).Uint("invoice_id", invoice.ID).Msg("Failed to update invoice")
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}

// Delete elimina una factura (soft delete)
<<<<<<< HEAD
func (r *InvoiceRepository) Delete(ctx context.Context, id uint) error {
=======
func (r *Repository) DeleteInvoice(ctx context.Context, id uint) error {
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
	if err := r.db.Conn(ctx).Delete(&models.Invoice{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}

	return nil
}

<<<<<<< HEAD
// ExistsForOrder verifica si existe una factura para una orden y proveedor
func (r *InvoiceRepository) ExistsForOrder(ctx context.Context, orderID string, providerID uint) (bool, error) {
	var count int64

	if err := r.db.Conn(ctx).Model(&models.Invoice{}).
		Where("order_id = ? AND invoicing_provider_id = ?", orderID, providerID).
=======
// InvoiceExistsForOrder verifica si existe una factura VÁLIDA para una orden e integración
// Solo considera facturas con status pending, issued o draft (excluye failed y cancelled)
func (r *Repository) InvoiceExistsForOrder(ctx context.Context, orderID string, integrationID uint) (bool, error) {
	var count int64

	// Solo contar facturas válidas (no fallidas ni canceladas)
	// Las facturas con status "failed" o "cancelled" NO bloquean crear una nueva factura
	if err := r.db.Conn(ctx).Model(&models.Invoice{}).
		Where("order_id = ? AND (invoicing_integration_id = ? OR invoicing_provider_id = ?)", orderID, integrationID, integrationID).
		Where("status NOT IN (?)", []string{"failed", "cancelled"}).
>>>>>>> 7b7c2054fa8e6cf0840b58d299ba6b7ca4e6b49e
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check invoice existence: %w", err)
	}

	return count > 0, nil
}
