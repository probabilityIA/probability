package orderstatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/ports"
)

// Re-export IRepository para que sea accesible desde otros módulos
type IRepository = ports.IRepository
