package app

import (
	"encoding/json"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Category weights
const (
	weightDataQuality     = 0.30
	weightPurchaseHistory = 0.25
	weightLogistics       = 0.20
	weightOrderChars      = 0.15
	weightPaymentRisk     = 0.10
)

// CalculateOrderScore calcula el score de una orden usando un sistema de categorias ponderadas.
// Retorna el score final, la lista plana de factores negativos, y el desglose completo.
func (uc *UseCaseScore) CalculateOrderScore(order *entities.ScoreOrder) (float64, []string, *entities.ScoreBreakdown) {
	var categories []entities.CategoryResult
	var allFactors []string

	// --- Category 1: Data Quality (30%) ---
	dqScore, dqFactors := uc.scoreDataQuality(order)
	dqWeighted := dqScore * weightDataQuality
	categories = append(categories, entities.CategoryResult{
		Name:          "Calidad de datos",
		Weight:        weightDataQuality,
		RawScore:      math.Round(dqScore*100) / 100,
		WeightedScore: math.Round(dqWeighted*100) / 100,
		Factors:       dqFactors,
	})
	allFactors = append(allFactors, dqFactors...)

	// --- Category 2: Purchase History (25%) ---
	phScore, phFactors := uc.scorePurchaseHistory(order)
	phWeighted := phScore * weightPurchaseHistory
	categories = append(categories, entities.CategoryResult{
		Name:          "Historial de compra",
		Weight:        weightPurchaseHistory,
		RawScore:      math.Round(phScore*100) / 100,
		WeightedScore: math.Round(phWeighted*100) / 100,
		Factors:       phFactors,
	})
	allFactors = append(allFactors, phFactors...)

	// --- Category 3: Logistics (20%) ---
	logScore, logFactors := uc.scoreLogistics(order)
	logWeighted := logScore * weightLogistics
	categories = append(categories, entities.CategoryResult{
		Name:          "Logistica",
		Weight:        weightLogistics,
		RawScore:      math.Round(logScore*100) / 100,
		WeightedScore: math.Round(logWeighted*100) / 100,
		Factors:       logFactors,
	})
	allFactors = append(allFactors, logFactors...)

	// --- Category 4: Order Characteristics (15%) ---
	ocScore, ocFactors := uc.scoreOrderCharacteristics(order)
	ocWeighted := ocScore * weightOrderChars
	categories = append(categories, entities.CategoryResult{
		Name:          "Caracteristicas del pedido",
		Weight:        weightOrderChars,
		RawScore:      math.Round(ocScore*100) / 100,
		WeightedScore: math.Round(ocWeighted*100) / 100,
		Factors:       ocFactors,
	})
	allFactors = append(allFactors, ocFactors...)

	// --- Category 5: Payment Risk (10%) ---
	prScore, prFactors := uc.scorePaymentRisk(order)
	prWeighted := prScore * weightPaymentRisk
	categories = append(categories, entities.CategoryResult{
		Name:          "Riesgo de pago",
		Weight:        weightPaymentRisk,
		RawScore:      math.Round(prScore*100) / 100,
		WeightedScore: math.Round(prWeighted*100) / 100,
		Factors:       prFactors,
	})
	allFactors = append(allFactors, prFactors...)

	// --- Final score ---
	finalScore := dqWeighted + phWeighted + logWeighted + ocWeighted + prWeighted

	// Clamp to [0, 100]
	if finalScore < 0 {
		finalScore = 0
	}
	if finalScore > 100 {
		finalScore = 100
	}

	// Round to 2 decimals
	finalScore = math.Round(finalScore*100) / 100

	breakdown := &entities.ScoreBreakdown{
		FinalScore:      finalScore,
		Categories:      categories,
		NegativeFactors: allFactors,
	}

	return finalScore, allFactors, breakdown
}

// scoreLogistics calculates Category 3: Logistics score (0-100)
func (uc *UseCaseScore) scoreLogistics(order *entities.ScoreOrder) (float64, []string) {
	var factors []string
	score := 100.0

	// Sub-signal 1: Delivery history failure rate (60%)
	deliveryScore := 100.0
	if order.DeliveryHistory != nil && order.DeliveryHistory.TotalShipments > 0 {
		failRate := float64(order.DeliveryHistory.FailedShipments) / float64(order.DeliveryHistory.TotalShipments) * 100
		deliveryScore = tierScoreDesc(failRate, []tier{
			{0, 100}, {10, 70}, {25, 40}, {50, 10},
		})
		if failRate >= 10 {
			factors = append(factors, "Alta tasa de envios fallidos del cliente")
		}
	}

	// Sub-signal 2: Distinct addresses (20%)
	addressScore := 100.0
	if order.CustomerHistory != nil && order.CustomerHistory.DistinctAddresses > 5 {
		addressScore = 50.0
		factors = append(factors, "Multiples direcciones de envio distintas")
	}

	// Sub-signal 3: Weight anomaly (20%)
	weightScore := 100.0
	if order.Weight != nil && *order.Weight > 50.0 {
		weightScore = 60.0
		factors = append(factors, "Peso del pedido inusualmente alto")
	}

	score = deliveryScore*0.60 + addressScore*0.20 + weightScore*0.20

	return math.Round(score*100) / 100, factors
}

// scoreOrderCharacteristics calculates Category 4: Order Characteristics score (0-100)
func (uc *UseCaseScore) scoreOrderCharacteristics(order *entities.ScoreOrder) (float64, []string) {
	var factors []string
	score := 100.0

	// Sub-signal 1: Order value range (30%)
	valueScore := 100.0
	if order.TotalAmount > 2000000 {
		valueScore = 50.0
		factors = append(factors, "Valor del pedido muy alto")
	} else if order.TotalAmount > 1000000 {
		valueScore = 75.0
	} else if order.TotalAmount <= 0 {
		valueScore = 30.0
		factors = append(factors, "Valor del pedido es cero o negativo")
	}

	// Sub-signal 2: Item count (20%)
	itemScore := 100.0
	if order.OrderItemCount > 20 {
		itemScore = 50.0
		factors = append(factors, "Cantidad de items inusualmente alta")
	} else if order.OrderItemCount > 10 {
		itemScore = 75.0
	}

	// Sub-signal 3: Coupon usage (15%)
	couponScore := 100.0
	if order.Coupon != nil && *order.Coupon != "" {
		couponScore = 80.0
	}

	// Sub-signal 4: Confirmation status (20%)
	confirmScore := 100.0
	if order.IsConfirmed != nil && !*order.IsConfirmed {
		confirmScore = 50.0
		factors = append(factors, "Pedido no confirmado")
	}

	// Sub-signal 5: COD penalty (15%)
	codScore := 100.0
	if uc.IsCODPayment(order) {
		codScore = 60.0
		factors = append(factors, "Pago Contra Entrega")
	}

	score = valueScore*0.30 + itemScore*0.20 + couponScore*0.15 + confirmScore*0.20 + codScore*0.15

	return math.Round(score*100) / 100, factors
}

// scorePaymentRisk calculates Category 5: Payment Risk score (0-100)
func (uc *UseCaseScore) scorePaymentRisk(order *entities.ScoreOrder) (float64, []string) {
	var factors []string
	score := 100.0

	// Sub-signal 1: Payment status (50%)
	paymentScore := 100.0
	if !order.IsPaid {
		paymentScore = 40.0
		factors = append(factors, "Pedido no pagado")
	}

	// Sub-signal 2: COD risk (30%)
	codRiskScore := 100.0
	if uc.IsCODPayment(order) {
		codRiskScore = 50.0
		// COD factor is already added in order characteristics, avoid duplicate
	}

	// Sub-signal 3: Customer COD history (20%)
	codHistoryScore := 100.0
	if order.CustomerHistory != nil && order.CustomerHistory.TotalOrders > 0 {
		codRate := float64(order.CustomerHistory.CODOrderCount) / float64(order.CustomerHistory.TotalOrders) * 100
		if codRate >= 80 {
			codHistoryScore = 40.0
			factors = append(factors, "Cliente con alta proporcion de pedidos contra entrega")
		} else if codRate >= 50 {
			codHistoryScore = 70.0
		}
	}

	score = paymentScore*0.50 + codRiskScore*0.30 + codHistoryScore*0.20

	return math.Round(score*100) / 100, factors
}

// GetStaticNegativeFactors obtiene la lista de factores negativos estáticos
func (uc *UseCaseScore) GetStaticNegativeFactors(order *entities.ScoreOrder) []string {
	var factors []string

	// 1. Validación de correo
	if !uc.isValidEmail(order.CustomerEmail) {
		factors = append(factors, "Email válido")
	}

	// 2. Nombre y apellido
	if order.CustomerName == "" || !strings.Contains(strings.TrimSpace(order.CustomerName), " ") {
		factors = append(factors, "Nombre y apellido")
	}

	// 3. Canal de venta (Platform)
	if order.Platform == "" {
		factors = append(factors, "Canal de venta")
	}

	// 4. Teléfono
	if order.CustomerPhone == "" {
		factors = append(factors, "Teléfono")
	}

	// 5. Dirección (Longitud mínima)
	if len(order.ShippingStreet) <= 5 {
		factors = append(factors, "Dirección")
	}

	// 6. Complemento de dirección
	address2 := order.Address2

	// Fallback 1: Si el campo transitorio está vacío, buscar en las direcciones relacionadas
	if address2 == "" && len(order.Addresses) > 0 {
		for _, addr := range order.Addresses {
			// Priorizar dirección de envío o usar la primera disponible si no tiene tipo
			if (addr.Type == "shipping" || addr.Type == "") && addr.Street2 != "" {
				address2 = addr.Street2
				break
			}
		}
	}

	// Fallback 2: Check ChannelMetadata RawData (Crucial for Shopify legacy/unmapped orders)
	if (address2 == "" || len(address2) < 2) && len(order.ChannelMetadata) > 0 {
		for _, meta := range order.ChannelMetadata {
			// Try to parse RawData as Shopify order structure
			if len(meta.RawData) > 0 {
				var rawData struct {
					ShippingAddress struct {
						Address2 string `json:"address2"`
					} `json:"shipping_address"`
				}
				// We need to unmarshal because RawData is []byte
				if err := json.Unmarshal(meta.RawData, &rawData); err == nil {
					if rawData.ShippingAddress.Address2 != "" {
						address2 = rawData.ShippingAddress.Address2
						break
					}
				}
			}
		}
	}

	// Validar solo si realmente está vacío o es muy corto (ajustado a < 2 para evitar falsos positivos con "Int 2")
	if address2 == "" || len(address2) < 2 {
		factors = append(factors, "Complemento de dirección")
	}

	// 7. Historial de compra
	if order.CustomerOrderCount == 0 {
		factors = append(factors, "Historial de compra")
	}

	return factors
}

func (uc *UseCaseScore) isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	// Regex simple
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// IsCODPayment verifica si el pago es contra entrega (COD)
func (uc *UseCaseScore) IsCODPayment(order *entities.ScoreOrder) bool {
	// 1. Check PaymentMethodID if we have a mapping (Placeholder)
	// 2. Check Financial Details (Shopify)

	keywords := []string{"cod", "cash", "contra"}

	// Check Payments slice first
	if len(order.Payments) > 0 {
		for _, payment := range order.Payments {
			if payment.Gateway != nil {
				gw := strings.ToLower(*payment.Gateway)
				for _, kw := range keywords {
					if strings.Contains(gw, kw) {
						return true
					}
				}
			}
		}
	}

	// 3. Fallback: Check COD Total if already calculated
	if order.CodTotal != nil && *order.CodTotal > 0 {
		return true
	}

	// 4. Fallback: Check PaymentDetails JSONB (crucial for Shopify if Payments slice is empty)
	if order.PaymentDetails != nil {
		var details struct {
			Gateway             string   `json:"gateway"`
			PaymentGatewayNames []string `json:"payment_gateway_names"`
		}

		// Unmarshal only what we need
		if err := json.Unmarshal(order.PaymentDetails, &details); err == nil {
			// Check single gateway
			if details.Gateway != "" {
				gw := strings.ToLower(details.Gateway)
				for _, kw := range keywords {
					if strings.Contains(gw, kw) {
						return true
					}
				}
			}
			// Check gateway names array
			for _, name := range details.PaymentGatewayNames {
				gw := strings.ToLower(name)
				for _, kw := range keywords {
					if strings.Contains(gw, kw) {
						return true
					}
				}
			}
		}
	}

	// 5. Fallback: Check Tags in Metadata
	if order.Metadata != nil {
		var metadata struct {
			Tags interface{} `json:"tags"` // Can be string or array in some cases, usually string in Shopify
		}
		if err := json.Unmarshal(order.Metadata, &metadata); err == nil {
			if tagsStr, ok := metadata.Tags.(string); ok && tagsStr != "" {
				lower := strings.ToLower(tagsStr)
				if strings.Contains(lower, "cod") || strings.Contains(lower, "contra") {
					return true
				}
			}
		}
	}

	return false
}

// RemoveAccents normaliza el texto eliminando acentos
func (uc *UseCaseScore) RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
