package generator

import (
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
)

// OrderGenerator genera órdenes canónicas aleatorias
type OrderGenerator struct{}

// New crea una nueva instancia del generador
func New() domain.IOrderGenerator {
	return &OrderGenerator{}
}
