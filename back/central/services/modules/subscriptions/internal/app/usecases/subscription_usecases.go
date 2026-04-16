package usecases

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/repository"
)

type SubscriptionUsecase struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionUsecase(repo repository.SubscriptionRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{repo: repo}
}

// GetBusinessSubscription devuelve la info de suscripción y si está activo o bloqueado
func (uc *SubscriptionUsecase) GetBusinessSubscription(ctx context.Context, businessID uint) (*entities.BusinessSubscription, error) {
	return uc.repo.GetLatestByBusinessID(ctx, businessID)
}

// RegisterSubscriptionPayment lo llamará el Super Admin para habilitar una cuenta y añadir X meses
func (uc *SubscriptionUsecase) RegisterSubscriptionPayment(ctx context.Context, businessID uint, amount float64, monthsToAdd int, paymentRef *string, notes *string) error {
	now := time.Now()
	endDate := now.AddDate(0, monthsToAdd, 0)

	// Crear registro de la suscripción
	sub := &entities.BusinessSubscription{
		BusinessID:       businessID,
		Amount:           amount,
		StartDate:        now,
		EndDate:          endDate,
		Status:           entities.SubscriptionStatusPaid,
		PaymentReference: paymentRef,
		Notes:            notes,
	}

	err := uc.repo.Create(ctx, sub)
	if err != nil {
		return err
	}

	// Actualizar modelo de usuario
	endDateStr := endDate.Format(time.RFC3339)
	return uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusActive, &endDateStr)
}

// DisableBusinessSubscription lo llamará el Super Admin o un cron job para dar de baja a una cuenta
func (uc *SubscriptionUsecase) DisableBusinessSubscription(ctx context.Context, businessID uint) error {
	// Poner fecha de expiración en nulo o en pasado para bloquear, y estado en expired
	now := time.Now()
	endDateStr := now.Format(time.RFC3339)
	return uc.repo.UpdateBusinessSubscriptionStatus(ctx, businessID, entities.BusinessStatusExpired, &endDateStr)
}

// Asegura que todos los negocios estén activos por defecto al iniciar
func (uc *SubscriptionUsecase) EnsureAllBusinessesActive(ctx context.Context) error {
	return uc.repo.EnsureAllBusinessesActive(ctx)
}
