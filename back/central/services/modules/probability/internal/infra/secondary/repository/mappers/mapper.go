package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func OrderToScoreOrder(m *models.Order) *entities.ScoreOrder {
	order := &entities.ScoreOrder{
		ID:             m.ID,
		OrderNumber:    m.OrderNumber,
		IntegrationID:  m.IntegrationID,
		CustomerEmail:  m.CustomerEmail,
		CustomerName:   m.CustomerName,
		Platform:       m.Platform,
		CustomerPhone:  m.CustomerPhone,
		ShippingStreet: m.ShippingStreet,
		CodTotal:       m.CodTotal,
	}

	// BusinessID
	if m.BusinessID != nil {
		order.BusinessID = m.BusinessID
	}

	// CustomerID
	if m.CustomerID != nil {
		order.CustomerID = m.CustomerID
	}

	// JSONB fields
	if m.Metadata != nil {
		order.Metadata = []byte(m.Metadata)
	}
	if m.PaymentDetails != nil {
		order.PaymentDetails = []byte(m.PaymentDetails)
	}

	// Payments
	for _, p := range m.Payments {
		order.Payments = append(order.Payments, entities.ScorePayment{
			Gateway: p.Gateway,
		})
	}

	// Addresses
	for _, a := range m.Addresses {
		order.Addresses = append(order.Addresses, entities.ScoreAddress{
			Type:    a.Type,
			Street2: a.Street2,
		})
	}

	// ChannelMetadata
	for _, cm := range m.ChannelMetadata {
		if cm.RawData != nil {
			order.ChannelMetadata = append(order.ChannelMetadata, entities.ScoreChannelMetadata{
				RawData: []byte(cm.RawData),
			})
		}
	}

	return order
}
