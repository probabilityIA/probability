package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	errs "github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

func (uc *UseCase) CreateCustomPlan(ctx context.Context, dto dtos.CreateCustomPlanDTO, actorUserID uint) (*entities.SubscriptionType, error) {
	if dto.BusinessID == 0 {
		return nil, errs.ErrBusinessRequired
	}
	if dto.Months <= 0 {
		return nil, errs.ErrInvalidMonths
	}
	if dto.Name == "" || dto.Code == "" || dto.Price < 0 {
		return nil, errs.ErrInvalidSubscriptionType
	}

	for _, code := range dto.ModuleCodes {
		if !moduleregistry.IsValid(code) {
			return nil, errs.ErrInvalidModuleCode
		}
	}

	billingPeriod := dto.BillingPeriod
	if billingPeriod == "" {
		billingPeriod = "monthly"
	}

	businessID := dto.BusinessID
	subType := &entities.SubscriptionType{
		Name:                 dto.Name,
		Code:                 dto.Code,
		Description:          dto.Description,
		Price:                dto.Price,
		BillingPeriod:        billingPeriod,
		Active:               true,
		ModuleCodes:          dto.ModuleCodes,
		MaxEcommerceChannels: dto.MaxEcommerceChannels,
		BusinessID:           &businessID,
		IncludedShipments:    dto.IncludedShipments,
		ShipmentOveragePrice: dto.ShipmentOveragePrice,
		IncludedInvoices:     dto.IncludedInvoices,
		InvoiceOveragePrice:  dto.InvoiceOveragePrice,
		IncludedOrders:       dto.IncludedOrders,
		OrderOveragePrice:    dto.OrderOveragePrice,
	}

	if err := uc.repo.CreateSubscriptionType(ctx, subType); err != nil {
		return nil, err
	}

	if _, err := uc.RegisterPayment(ctx, dtos.RegisterPaymentDTO{
		BusinessID:         dto.BusinessID,
		SubscriptionTypeID: subType.ID,
		Months:             dto.Months,
		PaymentReference:   dto.PaymentReference,
		Notes:              dto.Notes,
	}, actorUserID); err != nil {
		return nil, err
	}

	return subType, nil
}

func (uc *UseCase) ListCustomPlans(ctx context.Context, businessID *uint) ([]entities.SubscriptionType, error) {
	return uc.repo.ListCustomPlans(ctx, businessID)
}
