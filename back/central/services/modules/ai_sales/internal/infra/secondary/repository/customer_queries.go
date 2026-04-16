package repository

import (
	"context"
	"fmt"
	"time"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// SearchCustomers busca clientes por DNI, email, telefono o nombre dentro de un negocio.
// Tabla consultada: clients (gestionada por modulo customers)
// Replicado localmente para evitar compartir repositorios entre modulos
func (r *repository) SearchCustomers(ctx context.Context, businessID uint, query string) ([]domain.CustomerSearchResult, error) {
	var clients []models.Client

	err := r.db.Conn(ctx).
		Where("business_id = ?", businessID).
		Where("deleted_at IS NULL").
		Where("dni = ? OR email = ? OR phone = ? OR name ILIKE ?",
			query, query, query, "%"+query+"%").
		Limit(5).
		Find(&clients).Error

	if err != nil {
		return nil, fmt.Errorf("error searching customers: %w", err)
	}

	results := make([]domain.CustomerSearchResult, 0, len(clients))
	for _, c := range clients {
		email := ""
		if c.Email != nil {
			email = *c.Email
		}
		dni := ""
		if c.Dni != nil {
			dni = *c.Dni
		}
		results = append(results, domain.CustomerSearchResult{
			ID:    c.ID,
			Name:  c.Name,
			Email: email,
			Phone: c.Phone,
			DNI:   dni,
		})
	}

	return results, nil
}

// GetCustomerLastAddress obtiene la ultima direccion de envio del cliente
// desde las ordenes anteriores (campos denormalizados en la tabla orders).
// Tabla consultada: orders (gestionada por modulo orders)
// Replicado localmente para evitar compartir repositorios entre modulos
func (r *repository) GetCustomerLastAddress(ctx context.Context, businessID uint, customerID uint) (*domain.CustomerLastAddress, error) {
	var order models.Order

	err := r.db.Conn(ctx).
		Where("business_id = ? AND customer_id = ? AND shipping_street != ''", businessID, customerID).
		Order("created_at DESC").
		Select("shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code, created_at").
		First(&order).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting customer last address: %w", err)
	}

	return &domain.CustomerLastAddress{
		Street:     order.ShippingStreet,
		City:       order.ShippingCity,
		State:      order.ShippingState,
		Country:    order.ShippingCountry,
		PostalCode: order.ShippingPostalCode,
		OrderDate:  order.CreatedAt.Format(time.DateOnly),
	}, nil
}

// GetWhatsAppIntegrationID obtiene el ID de integracion de WhatsApp para un negocio.
// Tabla consultada: integrations (gestionada por modulo integrations/core)
// Replicado localmente para evitar compartir repositorios entre modulos
func (r *repository) GetWhatsAppIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	var integration models.Integration

	err := r.db.Conn(ctx).
		Where("business_id = ? AND integration_type_id = ? AND is_active = ? AND deleted_at IS NULL",
			businessID, 2, true).
		Select("id").
		First(&integration).Error

	if err == gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("no active WhatsApp integration found for business %d", businessID)
	}
	if err != nil {
		return 0, fmt.Errorf("error getting WhatsApp integration: %w", err)
	}

	return integration.ID, nil
}
