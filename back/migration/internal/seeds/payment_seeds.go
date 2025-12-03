package seeds

import (
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// SeedPaymentMethods crea los métodos de pago iniciales
func SeedPaymentMethods(db *gorm.DB) error {
	paymentMethods := []models.PaymentMethod{
		{
			Code:        "credit_card",
			Name:        "Tarjeta de Crédito",
			Description: "Pago con tarjeta de crédito",
			Category:    "card",
			IsActive:    true,
			Icon:        "/icons/credit-card.svg",
			Color:       "#4F46E5",
		},
		{
			Code:        "debit_card",
			Name:        "Tarjeta de Débito",
			Description: "Pago con tarjeta de débito",
			Category:    "card",
			IsActive:    true,
			Icon:        "/icons/debit-card.svg",
			Color:       "#10B981",
		},
		{
			Code:        "paypal",
			Name:        "PayPal",
			Description: "Pago a través de PayPal",
			Category:    "digital_wallet",
			Provider:    "paypal",
			IsActive:    true,
			Icon:        "/icons/paypal.svg",
			Color:       "#0070BA",
		},
		{
			Code:        "bank_transfer",
			Name:        "Transferencia Bancaria",
			Description: "Transferencia bancaria directa",
			Category:    "bank_transfer",
			IsActive:    true,
			Icon:        "/icons/bank.svg",
			Color:       "#6366F1",
		},
		{
			Code:        "cash",
			Name:        "Efectivo",
			Description: "Pago en efectivo",
			Category:    "cash",
			IsActive:    true,
			Icon:        "/icons/cash.svg",
			Color:       "#059669",
		},
		{
			Code:        "cod",
			Name:        "Contra Entrega",
			Description: "Pago contra entrega (COD)",
			Category:    "cash",
			IsActive:    true,
			Icon:        "/icons/cod.svg",
			Color:       "#F59E0B",
		},
		{
			Code:        "mercadopago",
			Name:        "Mercado Pago",
			Description: "Pago a través de Mercado Pago",
			Category:    "digital_wallet",
			Provider:    "mercadopago",
			IsActive:    true,
			Icon:        "/icons/mercadopago.svg",
			Color:       "#00B1EA",
		},
		{
			Code:        "stripe",
			Name:        "Stripe",
			Description: "Pago a través de Stripe",
			Category:    "digital_wallet",
			Provider:    "stripe",
			IsActive:    true,
			Icon:        "/icons/stripe.svg",
			Color:       "#635BFF",
		},
	}

	for _, pm := range paymentMethods {
		if err := db.Where("code = ?", pm.Code).FirstOrCreate(&pm).Error; err != nil {
			return err
		}
	}

	return nil
}

// SeedShopifyMappings crea los mapeos para Shopify
func SeedShopifyMappings(db *gorm.DB) error {
	// Mapeos de métodos de pago de Shopify
	paymentMappings := []struct {
		OriginalMethod string
		PaymentCode    string
	}{
		{"shopify_payments", "credit_card"},
		{"paypal", "paypal"},
		{"manual", "bank_transfer"},
		{"cash_on_delivery", "cod"},
		{"stripe", "stripe"},
	}

	for _, mapping := range paymentMappings {
		var paymentMethod models.PaymentMethod
		if err := db.Where("code = ?", mapping.PaymentCode).First(&paymentMethod).Error; err != nil {
			continue // Skip if payment method doesn't exist
		}

		pmMapping := models.PaymentMethodMapping{
			IntegrationType: "shopify",
			OriginalMethod:  mapping.OriginalMethod,
			PaymentMethodID: paymentMethod.ID,
			IsActive:        true,
			Priority:        0,
		}

		if err := db.Where("integration_type = ? AND original_method = ?",
			pmMapping.IntegrationType, pmMapping.OriginalMethod).
			FirstOrCreate(&pmMapping).Error; err != nil {
			return err
		}
	}

	// Mapeos de estados de Shopify
	statusMappings := []models.OrderStatusMapping{
		{
			IntegrationType: "shopify",
			OriginalStatus:  "pending",
			MappedStatus:    "pending",
			IsActive:        true,
			Description:     "Orden pendiente de pago",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "authorized",
			MappedStatus:    "pending",
			IsActive:        true,
			Description:     "Pago autorizado, pendiente de captura",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "paid",
			MappedStatus:    "processing",
			IsActive:        true,
			Description:     "Orden pagada, en proceso",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "partially_paid",
			MappedStatus:    "on_hold",
			IsActive:        true,
			Description:     "Orden parcialmente pagada",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "fulfilled",
			MappedStatus:    "shipped",
			IsActive:        true,
			Description:     "Orden cumplida/enviada",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "partially_fulfilled",
			MappedStatus:    "processing",
			IsActive:        true,
			Description:     "Orden parcialmente cumplida",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "cancelled",
			MappedStatus:    "cancelled",
			IsActive:        true,
			Description:     "Orden cancelada",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "refunded",
			MappedStatus:    "refunded",
			IsActive:        true,
			Description:     "Orden reembolsada",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "partially_refunded",
			MappedStatus:    "refunded",
			IsActive:        true,
			Description:     "Orden parcialmente reembolsada",
		},
		{
			IntegrationType: "shopify",
			OriginalStatus:  "voided",
			MappedStatus:    "cancelled",
			IsActive:        true,
			Description:     "Orden anulada",
		},
	}

	for _, sm := range statusMappings {
		if err := db.Where("integration_type = ? AND original_status = ?",
			sm.IntegrationType, sm.OriginalStatus).
			FirstOrCreate(&sm).Error; err != nil {
			return err
		}
	}

	return nil
}

// SeedWhatsAppMappings crea los mapeos para WhatsApp
func SeedWhatsAppMappings(db *gorm.DB) error {
	// Mapeos de métodos de pago de WhatsApp
	paymentMappings := []struct {
		OriginalMethod string
		PaymentCode    string
	}{
		{"cash", "cash"},
		{"transfer", "bank_transfer"},
		{"card", "credit_card"},
	}

	for _, mapping := range paymentMappings {
		var paymentMethod models.PaymentMethod
		if err := db.Where("code = ?", mapping.PaymentCode).First(&paymentMethod).Error; err != nil {
			continue
		}

		pmMapping := models.PaymentMethodMapping{
			IntegrationType: "whatsapp",
			OriginalMethod:  mapping.OriginalMethod,
			PaymentMethodID: paymentMethod.ID,
			IsActive:        true,
			Priority:        0,
		}

		if err := db.Where("integration_type = ? AND original_method = ?",
			pmMapping.IntegrationType, pmMapping.OriginalMethod).
			FirstOrCreate(&pmMapping).Error; err != nil {
			return err
		}
	}

	// Mapeos de estados de WhatsApp
	statusMappings := []models.OrderStatusMapping{
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "received",
			MappedStatus:    "pending",
			IsActive:        true,
			Description:     "Orden recibida por WhatsApp",
		},
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "confirmed",
			MappedStatus:    "processing",
			IsActive:        true,
			Description:     "Orden confirmada",
		},
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "preparing",
			MappedStatus:    "processing",
			IsActive:        true,
			Description:     "Orden en preparación",
		},
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "ready",
			MappedStatus:    "shipped",
			IsActive:        true,
			Description:     "Orden lista para envío",
		},
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "delivered",
			MappedStatus:    "delivered",
			IsActive:        true,
			Description:     "Orden entregada",
		},
		{
			IntegrationType: "whatsapp",
			OriginalStatus:  "cancelled",
			MappedStatus:    "cancelled",
			IsActive:        true,
			Description:     "Orden cancelada",
		},
	}

	for _, sm := range statusMappings {
		if err := db.Where("integration_type = ? AND original_status = ?",
			sm.IntegrationType, sm.OriginalStatus).
			FirstOrCreate(&sm).Error; err != nil {
			return err
		}
	}

	return nil
}
