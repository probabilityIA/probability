package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/paymentstatus/internal/domain/entities"
)

// IRepository define la interfaz para el almacenamiento de estados de pago
type IRepository interface {
	// GetPaymentStatusByCode obtiene un estado de pago por su código
	GetPaymentStatusByCode(ctx context.Context, code string) (*entities.PaymentStatus, error)

	// GetPaymentStatusIDByCode obtiene el ID de un estado de pago por su código
	GetPaymentStatusIDByCode(ctx context.Context, code string) (*uint, error)

	// ListPaymentStatuses lista todos los estados de pago
	ListPaymentStatuses(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error)
}
