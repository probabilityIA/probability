package usecases

import (
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
)

// UseCases contiene los casos de uso del módulo test
type UseCases struct {
	generator domain.IOrderGenerator
	publisher domain.IOrderPublisher
}

// New crea una nueva instancia de UseCases
func New(generator domain.IOrderGenerator, publisher domain.IOrderPublisher) *UseCases {
	return &UseCases{
		generator: generator,
		publisher: publisher,
	}
}
