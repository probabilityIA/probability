package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/infra/secondary/repository/models"
)

// ToDomain convierte un modelo de infra a entidad de dominio
func ToDomain(model *models.PaymentStatus) entities.PaymentStatus {
	return model.ToDomain()
}

// ToDomainList convierte una lista de modelos a entidades de dominio
func ToDomainList(models []models.PaymentStatus) []entities.PaymentStatus {
	result := make([]entities.PaymentStatus, len(models))
	for i, model := range models {
		result[i] = model.ToDomain()
	}
	return result
}
