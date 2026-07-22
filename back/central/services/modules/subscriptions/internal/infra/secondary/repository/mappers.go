package repository

import (
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

func marshalModuleCodes(codes []string) datatypes.JSON {
	if codes == nil {
		codes = []string{}
	}
	raw, _ := json.Marshal(codes)
	return datatypes.JSON(raw)
}

func unmarshalModuleCodes(raw datatypes.JSON) []string {
	var codes []string
	if len(raw) == 0 {
		return codes
	}
	_ = json.Unmarshal(raw, &codes)
	return codes
}

func subscriptionTypeToEntity(m *models.SubscriptionType) *entities.SubscriptionType {
	return &entities.SubscriptionType{
		ID:                   m.ID,
		Name:                 m.Name,
		Code:                 m.Code,
		Description:          m.Description,
		Price:                m.Price,
		BillingPeriod:        m.BillingPeriod,
		Active:               m.Active,
		ModuleCodes:          unmarshalModuleCodes(m.Features),
		MaxEcommerceChannels: m.MaxEcommerceChannels,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func subscriptionToEntity(m *models.BusinessSubscription) *entities.BusinessSubscription {
	typeName := ""
	subTypeID := uint(0)
	if m.SubscriptionType != nil {
		typeName = m.SubscriptionType.Name
	}
	if m.SubscriptionTypeID != nil {
		subTypeID = *m.SubscriptionTypeID
	}

	months := 0
	if m.Months != nil {
		months = *m.Months
	}

	return &entities.BusinessSubscription{
		ID:                   m.ID,
		BusinessID:           m.BusinessID,
		SubscriptionTypeID:   subTypeID,
		SubscriptionTypeName: typeName,
		Months:               months,
		Amount:               m.Amount,
		StartDate:            m.StartDate,
		EndDate:              m.EndDate,
		Status:               m.Status,
		PaymentReference:     m.PaymentReference,
		Notes:                m.Notes,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func overrideToEntity(m *models.BusinessModuleOverride) *entities.BusinessModuleOverride {
	return &entities.BusinessModuleOverride{
		ID:              m.ID,
		BusinessID:      m.BusinessID,
		ModuleCode:      m.ModuleCode,
		GrantedByUserID: m.GrantedByUserID,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt,
	}
}
