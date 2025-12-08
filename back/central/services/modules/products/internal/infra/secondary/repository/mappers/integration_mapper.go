package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBProductIntegration convierte una integración de producto de dominio a modelo de base de datos
func ToDBProductIntegration(pi *domain.ProductBusinessIntegration) *models.ProductBusinessIntegration {
	if pi == nil {
		return nil
	}
	dbPI := &models.ProductBusinessIntegration{
		Model: gorm.Model{
			ID:        pi.ID,
			CreatedAt: pi.CreatedAt,
			UpdatedAt: pi.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		ProductID:         pi.ProductID,
		BusinessID:        pi.BusinessID,
		IntegrationID:     pi.IntegrationID,
		ExternalProductID: pi.ExternalProductID,
	}
	if pi.DeletedAt != nil {
		dbPI.DeletedAt = gorm.DeletedAt{Time: *pi.DeletedAt, Valid: true}
	}
	return dbPI
}

// ToDomainProductIntegration convierte una integración de producto de base de datos a dominio
func ToDomainProductIntegration(pi *models.ProductBusinessIntegration) *domain.ProductBusinessIntegration {
	if pi == nil {
		return nil
	}
	var deletedAt *time.Time
	if pi.DeletedAt.Valid {
		deletedAt = &pi.DeletedAt.Time
	}

	result := &domain.ProductBusinessIntegration{
		ID:                pi.ID,
		CreatedAt:         pi.CreatedAt,
		UpdatedAt:         pi.UpdatedAt,
		DeletedAt:         deletedAt,
		ProductID:         pi.ProductID,
		BusinessID:        pi.BusinessID,
		IntegrationID:     pi.IntegrationID,
		ExternalProductID: pi.ExternalProductID,
	}

	// Incluir información de la integración si está disponible (Preload)
	// Verificar si Integration está cargado (el ID del modelo GORM será != 0)
	if pi.Integration.ID != 0 {
		result.IntegrationName = pi.Integration.Name
		// Verificar si IntegrationType está cargado
		if pi.Integration.IntegrationType.ID != 0 {
			result.IntegrationType = pi.Integration.IntegrationType.Code
		}
	}

	return result
}

// ToDomainProductIntegrations convierte una lista de integraciones de producto de base de datos a dominio
func ToDomainProductIntegrations(pis []models.ProductBusinessIntegration) []domain.ProductBusinessIntegration {
	if pis == nil {
		return nil
	}
	result := make([]domain.ProductBusinessIntegration, len(pis))
	for i, pi := range pis {
		result[i] = *ToDomainProductIntegration(&pi)
	}
	return result
}
