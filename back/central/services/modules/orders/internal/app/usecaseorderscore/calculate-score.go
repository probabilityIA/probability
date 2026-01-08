package usecaseorderscore

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CalculateOrderScore calcula el score de una orden y sus factores negativos
func (uc *UseCaseOrderScore) CalculateOrderScore(order *domain.ProbabilityOrder) (float64, []string) {
	// Start with 100
	score := 100.0

	// Get Static Negative Factors
	staticFactors := uc.GetStaticNegativeFactors(order)
	fmt.Printf("[CalculateOrderScore] Order %s - Email: '%s', Name: '%s', Phone: '%s', Street: '%s', Street2: '%s', Count: %d, Platform: '%s', Factors: %v\n",
		order.OrderNumber, order.CustomerEmail, order.CustomerName, order.CustomerPhone, order.ShippingStreet, order.Address2, order.CustomerOrderCount, order.Platform, staticFactors)

	// Apply penalties for static factors
	// Each factor reduces score by 10 (example weight)
	// Python reference used weights. We will assume 10 per factor for now or matching Python.
	// Mapa de penalizaciones
	criteriaMap := map[string]float64{
		"Email válido":             -10,
		"Nombre y apellido":        -10,
		"Canal de venta":           -10,
		"Teléfono":                 -10,
		"Dirección":                -10,
		"Complemento de dirección": -10,
		"Historial de compra":      -10,
	}

	// Calculate Score based on factors
	for _, factor := range staticFactors {
		if penalty, exists := criteriaMap[factor]; exists {
			score += penalty // Penalty is negative
		}
	}

	// COD Logic
	if uc.IsCODPayment(order) {
		score = score * 0.8 // Apply 20% reduction
		// Add to factors so user knows why it's not 100%
		staticFactors = append(staticFactors, "Pago Contra Entrega")
	}

	// Ensure limits
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Redondear a 2 decimales
	finalScore := float64(int(score*100)) / 100
	fmt.Printf("[CalculateOrderScore] Order %s - Final Score: %.2f\n", order.OrderNumber, finalScore)
	return finalScore, staticFactors
}

// GetStaticNegativeFactors obtiene la lista de factores negativos estáticos
func (uc *UseCaseOrderScore) GetStaticNegativeFactors(order *domain.ProbabilityOrder) []string {
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
				// We need to marshal/unmarshal because RawData is datatypes.JSON (byte array)
				bytes, err := meta.RawData.MarshalJSON()
				if err == nil {
					if err := json.Unmarshal(bytes, &rawData); err == nil {
						if rawData.ShippingAddress.Address2 != "" {
							fmt.Printf("[CalculateOrderScore] Found Address2 in ChannelMetadata: '%s'\n", rawData.ShippingAddress.Address2)
							address2 = rawData.ShippingAddress.Address2
							break
						}
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

func (uc *UseCaseOrderScore) isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	// Regex simple
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// IsCODPayment verifica si el pago es contra entrega (COD)
func (uc *UseCaseOrderScore) IsCODPayment(order *domain.ProbabilityOrder) bool {
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
		bytes, err := order.PaymentDetails.MarshalJSON()
		if err == nil {
			if err := json.Unmarshal(bytes, &details); err == nil {
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
	}

	// 5. Fallback: Check Tags in Metadata
	if order.Metadata != nil {
		var metadata struct {
			Tags interface{} `json:"tags"` // Can be string or array in some cases, usually string in Shopify
		}
		bytes, err := order.Metadata.MarshalJSON()
		if err == nil {
			if err := json.Unmarshal(bytes, &metadata); err == nil {
				if tagsStr, ok := metadata.Tags.(string); ok && tagsStr != "" {
					lower := strings.ToLower(tagsStr)
					if strings.Contains(lower, "cod") || strings.Contains(lower, "contra") {
						return true
					}
				}
			}
		}
	}

	return false
}

// RemoveAccents normaliza el texto eliminando acentos
func (uc *UseCaseOrderScore) RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
