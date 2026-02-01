package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/secondary/repository/models"
)

// ToModel convierte una entidad de dominio a modelo de infra
func ToModel(entity entities.PaymentStatus) *models.PaymentStatus {
	return models.FromDomain(entity)
}

// ToModelList convierte una lista de entidades a modelos de infra
func ToModelList(entities []entities.PaymentStatus) []*models.PaymentStatus {
	result := make([]*models.PaymentStatus, len(entities))
	for i, entity := range entities {
		result[i] = models.FromDomain(entity)
	}
	return result
}
