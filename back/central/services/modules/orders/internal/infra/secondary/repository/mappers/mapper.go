package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBOrder convierte una orden de dominio a modelo de base de datos
func ToDBOrder(o *domain.Order) *models.Order {
	if o == nil {
		return nil
	}

	return &models.Order{
		ID:                  o.ID,
		CreatedAt:           o.CreatedAt,
		UpdatedAt:           o.UpdatedAt,
		DeletedAt:           o.DeletedAt,
		BusinessID:          o.BusinessID,
		IntegrationID:       o.IntegrationID,
		IntegrationType:     o.IntegrationType,
		Platform:            o.Platform,
		ExternalID:          o.ExternalID,
		OrderNumber:         o.OrderNumber,
		InternalNumber:      o.InternalNumber,
		Subtotal:            o.Subtotal,
		Tax:                 o.Tax,
		Discount:            o.Discount,
		ShippingCost:        o.ShippingCost,
		TotalAmount:         o.TotalAmount,
		Currency:            o.Currency,
		CodTotal:            o.CodTotal,
		CustomerID:          o.CustomerID,
		CustomerName:        o.CustomerName,
		CustomerEmail:       o.CustomerEmail,
		CustomerPhone:       o.CustomerPhone,
		CustomerDNI:         o.CustomerDNI,
		ShippingStreet:      o.ShippingStreet,
		ShippingCity:        o.ShippingCity,
		ShippingState:       o.ShippingState,
		ShippingCountry:     o.ShippingCountry,
		ShippingPostalCode:  o.ShippingPostalCode,
		ShippingLat:         o.ShippingLat,
		ShippingLng:         o.ShippingLng,
		PaymentMethodID:     o.PaymentMethodID,
		IsPaid:              o.IsPaid,
		PaidAt:              o.PaidAt,
		TrackingNumber:      o.TrackingNumber,
		TrackingLink:        o.TrackingLink,
		GuideID:             o.GuideID,
		GuideLink:           o.GuideLink,
		DeliveryDate:        o.DeliveryDate,
		DeliveredAt:         o.DeliveredAt,
		DeliveryProbability: o.DeliveryProbability,
		WarehouseID:         o.WarehouseID,
		WarehouseName:       o.WarehouseName,
		DriverID:            o.DriverID,
		DriverName:          o.DriverName,
		IsLastMile:          o.IsLastMile,
		Weight:              o.Weight,
		Height:              o.Height,
		Width:               o.Width,
		Length:              o.Length,
		Boxes:               o.Boxes,
		OrderTypeID:         o.OrderTypeID,
		OrderTypeName:       o.OrderTypeName,
		Status:              o.Status,
		OriginalStatus:      o.OriginalStatus,
		Notes:               o.Notes,
		Coupon:              o.Coupon,
		Approved:            o.Approved,
		UserID:              o.UserID,
		UserName:            o.UserName,
		Invoiceable:         o.Invoiceable,
		InvoiceURL:          o.InvoiceURL,
		InvoiceID:           o.InvoiceID,
		InvoiceProvider:     o.InvoiceProvider,
		OrderStatusURL:      o.OrderStatusURL,
		Items:               o.Items,
		Metadata:            o.Metadata,
		FinancialDetails:    o.FinancialDetails,
		ShippingDetails:     o.ShippingDetails,
		PaymentDetails:      o.PaymentDetails,
		FulfillmentDetails:  o.FulfillmentDetails,
		OccurredAt:          o.OccurredAt,
		ImportedAt:          o.ImportedAt,
		OrderItems:          ToDBOrderItems(o.OrderItems),
		Addresses:           ToDBAddresses(o.Addresses),
		Payments:            ToDBPayments(o.Payments),
		Shipments:           ToDBShipments(o.Shipments),
		ChannelMetadata:     ToDBChannelMetadataList(o.ChannelMetadata),
	}
}

// ToDomainOrder convierte una orden de base de datos a dominio
func ToDomainOrder(o *models.Order) *domain.Order {
	if o == nil {
		return nil
	}

	return &domain.Order{
		ID:                  o.ID,
		CreatedAt:           o.CreatedAt,
		UpdatedAt:           o.UpdatedAt,
		DeletedAt:           o.DeletedAt,
		BusinessID:          o.BusinessID,
		IntegrationID:       o.IntegrationID,
		IntegrationType:     o.IntegrationType,
		Platform:            o.Platform,
		ExternalID:          o.ExternalID,
		OrderNumber:         o.OrderNumber,
		InternalNumber:      o.InternalNumber,
		Subtotal:            o.Subtotal,
		Tax:                 o.Tax,
		Discount:            o.Discount,
		ShippingCost:        o.ShippingCost,
		TotalAmount:         o.TotalAmount,
		Currency:            o.Currency,
		CodTotal:            o.CodTotal,
		CustomerID:          o.CustomerID,
		CustomerName:        o.CustomerName,
		CustomerEmail:       o.CustomerEmail,
		CustomerPhone:       o.CustomerPhone,
		CustomerDNI:         o.CustomerDNI,
		ShippingStreet:      o.ShippingStreet,
		ShippingCity:        o.ShippingCity,
		ShippingState:       o.ShippingState,
		ShippingCountry:     o.ShippingCountry,
		ShippingPostalCode:  o.ShippingPostalCode,
		ShippingLat:         o.ShippingLat,
		ShippingLng:         o.ShippingLng,
		PaymentMethodID:     o.PaymentMethodID,
		IsPaid:              o.IsPaid,
		PaidAt:              o.PaidAt,
		TrackingNumber:      o.TrackingNumber,
		TrackingLink:        o.TrackingLink,
		GuideID:             o.GuideID,
		GuideLink:           o.GuideLink,
		DeliveryDate:        o.DeliveryDate,
		DeliveredAt:         o.DeliveredAt,
		DeliveryProbability: o.DeliveryProbability,
		WarehouseID:         o.WarehouseID,
		WarehouseName:       o.WarehouseName,
		DriverID:            o.DriverID,
		DriverName:          o.DriverName,
		IsLastMile:          o.IsLastMile,
		Weight:              o.Weight,
		Height:              o.Height,
		Width:               o.Width,
		Length:              o.Length,
		Boxes:               o.Boxes,
		OrderTypeID:         o.OrderTypeID,
		OrderTypeName:       o.OrderTypeName,
		Status:              o.Status,
		OriginalStatus:      o.OriginalStatus,
		Notes:               o.Notes,
		Coupon:              o.Coupon,
		Approved:            o.Approved,
		UserID:              o.UserID,
		UserName:            o.UserName,
		Invoiceable:         o.Invoiceable,
		InvoiceURL:          o.InvoiceURL,
		InvoiceID:           o.InvoiceID,
		InvoiceProvider:     o.InvoiceProvider,
		OrderStatusURL:      o.OrderStatusURL,
		Items:               o.Items,
		Metadata:            o.Metadata,
		FinancialDetails:    o.FinancialDetails,
		ShippingDetails:     o.ShippingDetails,
		PaymentDetails:      o.PaymentDetails,
		FulfillmentDetails:  o.FulfillmentDetails,
		OccurredAt:          o.OccurredAt,
		ImportedAt:          o.ImportedAt,
		OrderItems:          ToDomainOrderItems(o.OrderItems),
		Addresses:           ToDomainAddresses(o.Addresses),
		Payments:            ToDomainPayments(o.Payments),
		Shipments:           ToDomainShipments(o.Shipments),
		ChannelMetadata:     ToDomainChannelMetadataList(o.ChannelMetadata),
	}
}

// ToDBOrderItems convierte una lista de items de dominio a base de datos
func ToDBOrderItems(items []domain.OrderItem) []models.OrderItem {
	if items == nil {
		return nil
	}
	result := make([]models.OrderItem, len(items))
	for i, item := range items {
		result[i] = models.OrderItem{
			Model: gorm.Model{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: gorm.DeletedAt{}, // Handle DeletedAt if needed
			},
			OrderID:           item.OrderID,
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			UnitPrice:         item.UnitPrice,
			TotalPrice:        item.TotalPrice,
			Currency:          item.Currency,
			Discount:          item.Discount,
			Tax:               item.Tax,
			TaxRate:           item.TaxRate,
			FulfillmentStatus: item.FulfillmentStatus,
			Metadata:          item.Metadata,
		}
		if item.DeletedAt != nil {
			result[i].DeletedAt = gorm.DeletedAt{Time: *item.DeletedAt, Valid: true}
		}
	}
	return result
}

// ToDomainOrderItems convierte una lista de items de base de datos a dominio
func ToDomainOrderItems(items []models.OrderItem) []domain.OrderItem {
	if items == nil {
		return nil
	}
	result := make([]domain.OrderItem, len(items))
	for i, item := range items {
		var deletedAt *time.Time
		if item.DeletedAt.Valid {
			deletedAt = &item.DeletedAt.Time
		}
		// Mapear campos básicos del modelo
		domainItem := domain.OrderItem{
			ID:                item.ID,
			CreatedAt:         item.CreatedAt,
			UpdatedAt:         item.UpdatedAt,
			DeletedAt:         deletedAt,
			OrderID:           item.OrderID,
			ProductID:         item.ProductID,
			VariantID:         item.VariantID,
			Quantity:          item.Quantity,
			UnitPrice:         item.UnitPrice,
			TotalPrice:        item.TotalPrice,
			Currency:          item.Currency,
			Discount:          item.Discount,
			Tax:               item.Tax,
			TaxRate:           item.TaxRate,
			FulfillmentStatus: item.FulfillmentStatus,
			Metadata:          item.Metadata,
		}

		// Si el Product está preloaded, obtener información del producto
		// Verificar si Product fue preloaded (ID no es cero)
		if item.Product.ID != "" && item.Product.SKU != "" {
			domainItem.ProductSKU = item.Product.SKU
			domainItem.ProductName = item.Product.Name
			// ProductTitle no existe en el modelo Product, se deja vacío
		}

		result[i] = domainItem
	}
	return result
}

// ToDBAddresses convierte una lista de direcciones de dominio a base de datos
func ToDBAddresses(addresses []domain.Address) []models.Address {
	if addresses == nil {
		return nil
	}
	result := make([]models.Address, len(addresses))
	for i, addr := range addresses {
		result[i] = models.Address{
			Model: gorm.Model{
				ID:        addr.ID,
				CreatedAt: addr.CreatedAt,
				UpdatedAt: addr.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			Type:         addr.Type,
			OrderID:      addr.OrderID,
			FirstName:    addr.FirstName,
			LastName:     addr.LastName,
			Company:      addr.Company,
			Phone:        addr.Phone,
			Street:       addr.Street,
			Street2:      addr.Street2,
			City:         addr.City,
			State:        addr.State,
			Country:      addr.Country,
			PostalCode:   addr.PostalCode,
			Latitude:     addr.Latitude,
			Longitude:    addr.Longitude,
			Instructions: addr.Instructions,
			IsDefault:    addr.IsDefault,
			Metadata:     addr.Metadata,
		}
		if addr.DeletedAt != nil {
			result[i].DeletedAt = gorm.DeletedAt{Time: *addr.DeletedAt, Valid: true}
		}
	}
	return result
}

// ToDomainAddresses convierte una lista de direcciones de base de datos a dominio
func ToDomainAddresses(addresses []models.Address) []domain.Address {
	if addresses == nil {
		return nil
	}
	result := make([]domain.Address, len(addresses))
	for i, addr := range addresses {
		var deletedAt *time.Time
		if addr.DeletedAt.Valid {
			deletedAt = &addr.DeletedAt.Time
		}
		result[i] = domain.Address{
			ID:           addr.ID,
			CreatedAt:    addr.CreatedAt,
			UpdatedAt:    addr.UpdatedAt,
			DeletedAt:    deletedAt,
			Type:         addr.Type,
			OrderID:      addr.OrderID,
			FirstName:    addr.FirstName,
			LastName:     addr.LastName,
			Company:      addr.Company,
			Phone:        addr.Phone,
			Street:       addr.Street,
			Street2:      addr.Street2,
			City:         addr.City,
			State:        addr.State,
			Country:      addr.Country,
			PostalCode:   addr.PostalCode,
			Latitude:     addr.Latitude,
			Longitude:    addr.Longitude,
			Instructions: addr.Instructions,
			IsDefault:    addr.IsDefault,
			Metadata:     addr.Metadata,
		}
	}
	return result
}

// ToDBPayments convierte una lista de pagos de dominio a base de datos
func ToDBPayments(payments []domain.Payment) []models.Payment {
	if payments == nil {
		return nil
	}
	result := make([]models.Payment, len(payments))
	for i, p := range payments {
		result[i] = models.Payment{
			Model: gorm.Model{
				ID:        p.ID,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			OrderID:          p.OrderID,
			PaymentMethodID:  p.PaymentMethodID,
			Amount:           p.Amount,
			Currency:         p.Currency,
			ExchangeRate:     p.ExchangeRate,
			Status:           p.Status,
			PaidAt:           p.PaidAt,
			ProcessedAt:      p.ProcessedAt,
			TransactionID:    p.TransactionID,
			PaymentReference: p.PaymentReference,
			Gateway:          p.Gateway,
			RefundAmount:     p.RefundAmount,
			RefundedAt:       p.RefundedAt,
			FailureReason:    p.FailureReason,
			Metadata:         p.Metadata,
		}
		if p.DeletedAt != nil {
			result[i].DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
		}
	}
	return result
}

// ToDomainPayments convierte una lista de pagos de base de datos a dominio
func ToDomainPayments(payments []models.Payment) []domain.Payment {
	if payments == nil {
		return nil
	}
	result := make([]domain.Payment, len(payments))
	for i, p := range payments {
		var deletedAt *time.Time
		if p.DeletedAt.Valid {
			deletedAt = &p.DeletedAt.Time
		}
		result[i] = domain.Payment{
			ID:               p.ID,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
			DeletedAt:        deletedAt,
			OrderID:          p.OrderID,
			PaymentMethodID:  p.PaymentMethodID,
			Amount:           p.Amount,
			Currency:         p.Currency,
			ExchangeRate:     p.ExchangeRate,
			Status:           p.Status,
			PaidAt:           p.PaidAt,
			ProcessedAt:      p.ProcessedAt,
			TransactionID:    p.TransactionID,
			PaymentReference: p.PaymentReference,
			Gateway:          p.Gateway,
			RefundAmount:     p.RefundAmount,
			RefundedAt:       p.RefundedAt,
			FailureReason:    p.FailureReason,
			Metadata:         p.Metadata,
		}
	}
	return result
}

// ToDBShipments convierte una lista de envíos de dominio a base de datos
func ToDBShipments(shipments []domain.Shipment) []models.Shipment {
	if shipments == nil {
		return nil
	}
	result := make([]models.Shipment, len(shipments))
	for i, s := range shipments {
		result[i] = models.Shipment{
			Model: gorm.Model{
				ID:        s.ID,
				CreatedAt: s.CreatedAt,
				UpdatedAt: s.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			OrderID:           s.OrderID,
			TrackingNumber:    s.TrackingNumber,
			TrackingURL:       s.TrackingURL,
			Carrier:           s.Carrier,
			CarrierCode:       s.CarrierCode,
			GuideID:           s.GuideID,
			GuideURL:          s.GuideURL,
			Status:            s.Status,
			ShippedAt:         s.ShippedAt,
			DeliveredAt:       s.DeliveredAt,
			ShippingAddressID: s.ShippingAddressID,
			ShippingCost:      s.ShippingCost,
			InsuranceCost:     s.InsuranceCost,
			TotalCost:         s.TotalCost,
			Weight:            s.Weight,
			Height:            s.Height,
			Width:             s.Width,
			Length:            s.Length,
			WarehouseID:       s.WarehouseID,
			WarehouseName:     s.WarehouseName,
			DriverID:          s.DriverID,
			DriverName:        s.DriverName,
			IsLastMile:        s.IsLastMile,
			EstimatedDelivery: s.EstimatedDelivery,
			DeliveryNotes:     s.DeliveryNotes,
			Metadata:          s.Metadata,
		}
		if s.DeletedAt != nil {
			result[i].DeletedAt = gorm.DeletedAt{Time: *s.DeletedAt, Valid: true}
		}
	}
	return result
}

// ToDomainShipments convierte una lista de envíos de base de datos a dominio
func ToDomainShipments(shipments []models.Shipment) []domain.Shipment {
	if shipments == nil {
		return nil
	}
	result := make([]domain.Shipment, len(shipments))
	for i, s := range shipments {
		var deletedAt *time.Time
		if s.DeletedAt.Valid {
			deletedAt = &s.DeletedAt.Time
		}
		result[i] = domain.Shipment{
			ID:                s.ID,
			CreatedAt:         s.CreatedAt,
			UpdatedAt:         s.UpdatedAt,
			DeletedAt:         deletedAt,
			OrderID:           s.OrderID,
			TrackingNumber:    s.TrackingNumber,
			TrackingURL:       s.TrackingURL,
			Carrier:           s.Carrier,
			CarrierCode:       s.CarrierCode,
			GuideID:           s.GuideID,
			GuideURL:          s.GuideURL,
			Status:            s.Status,
			ShippedAt:         s.ShippedAt,
			DeliveredAt:       s.DeliveredAt,
			ShippingAddressID: s.ShippingAddressID,
			ShippingCost:      s.ShippingCost,
			InsuranceCost:     s.InsuranceCost,
			TotalCost:         s.TotalCost,
			Weight:            s.Weight,
			Height:            s.Height,
			Width:             s.Width,
			Length:            s.Length,
			WarehouseID:       s.WarehouseID,
			WarehouseName:     s.WarehouseName,
			DriverID:          s.DriverID,
			DriverName:        s.DriverName,
			IsLastMile:        s.IsLastMile,
			EstimatedDelivery: s.EstimatedDelivery,
			DeliveryNotes:     s.DeliveryNotes,
			Metadata:          s.Metadata,
		}
	}
	return result
}

// ToDBChannelMetadataList convierte una lista de metadata de canal de dominio a base de datos
func ToDBChannelMetadataList(metadata []domain.OrderChannelMetadata) []models.OrderChannelMetadata {
	if metadata == nil {
		return nil
	}
	result := make([]models.OrderChannelMetadata, len(metadata))
	for i, m := range metadata {
		result[i] = models.OrderChannelMetadata{
			Model: gorm.Model{
				ID:        m.ID,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
				DeletedAt: gorm.DeletedAt{},
			},
			OrderID:       m.OrderID,
			ChannelSource: m.ChannelSource,
			IntegrationID: m.IntegrationID,
			RawData:       m.RawData,
			Version:       m.Version,
			ReceivedAt:    m.ReceivedAt,
			ProcessedAt:   m.ProcessedAt,
			IsLatest:      m.IsLatest,
			LastSyncedAt:  m.LastSyncedAt,
			SyncStatus:    m.SyncStatus,
		}
		if m.DeletedAt != nil {
			result[i].DeletedAt = gorm.DeletedAt{Time: *m.DeletedAt, Valid: true}
		}
	}
	return result
}

// ToDomainChannelMetadataList convierte una lista de metadata de canal de base de datos a dominio
func ToDomainChannelMetadataList(metadata []models.OrderChannelMetadata) []domain.OrderChannelMetadata {
	if metadata == nil {
		return nil
	}
	result := make([]domain.OrderChannelMetadata, len(metadata))
	for i, m := range metadata {
		var deletedAt *time.Time
		if m.DeletedAt.Valid {
			deletedAt = &m.DeletedAt.Time
		}
		result[i] = domain.OrderChannelMetadata{
			ID:            m.ID,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
			DeletedAt:     deletedAt,
			OrderID:       m.OrderID,
			ChannelSource: m.ChannelSource,
			IntegrationID: m.IntegrationID,
			RawData:       m.RawData,
			Version:       m.Version,
			ReceivedAt:    m.ReceivedAt,
			ProcessedAt:   m.ProcessedAt,
			IsLatest:      m.IsLatest,
			LastSyncedAt:  m.LastSyncedAt,
			SyncStatus:    m.SyncStatus,
		}
	}
	return result
}

// ToDBChannelMetadata convierte metadata de canal de dominio a base de datos (singular)
func ToDBChannelMetadata(m *domain.OrderChannelMetadata) *models.OrderChannelMetadata {
	if m == nil {
		return nil
	}

	result := &models.OrderChannelMetadata{
		Model: gorm.Model{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		OrderID:       m.OrderID,
		ChannelSource: m.ChannelSource,
		IntegrationID: m.IntegrationID,
		RawData:       m.RawData,
		Version:       m.Version,
		ReceivedAt:    m.ReceivedAt,
		ProcessedAt:   m.ProcessedAt,
		IsLatest:      m.IsLatest,
		LastSyncedAt:  m.LastSyncedAt,
		SyncStatus:    m.SyncStatus,
	}
	if m.DeletedAt != nil {
		result.DeletedAt = gorm.DeletedAt{Time: *m.DeletedAt, Valid: true}
	}
	return result
}

// ToDomainChannelMetadata convierte metadata de canal de base de datos a dominio (singular)
func ToDomainChannelMetadata(m *models.OrderChannelMetadata) *domain.OrderChannelMetadata {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &domain.OrderChannelMetadata{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     deletedAt,
		OrderID:       m.OrderID,
		ChannelSource: m.ChannelSource,
		IntegrationID: m.IntegrationID,
		RawData:       m.RawData,
		Version:       m.Version,
		ReceivedAt:    m.ReceivedAt,
		ProcessedAt:   m.ProcessedAt,
		IsLatest:      m.IsLatest,
		LastSyncedAt:  m.LastSyncedAt,
		SyncStatus:    m.SyncStatus,
	}
}
