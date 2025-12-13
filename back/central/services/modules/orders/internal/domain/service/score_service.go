package service

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// CalculateOrderScore calcula el score de una orden y sus factores negativos
func CalculateOrderScore(order *domain.Order) (float64, []string) {
	// Start with 100
	score := 100.0

	// Get Static Negative Factors
	staticFactors := GetStaticNegativeFactors(order)
	fmt.Printf("[DEBUG] Order %s - Street2: '%s', Count: %d, Factors: %v\n", order.OrderNumber, order.ShippingStreet2, order.CustomerOrderCount, staticFactors)

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

	// COD Logic (placeholder) that might add another factor?
	// If IS COD, usually we reduce score further or add a factor?
	// Reference Python: if payment_method == 'cod': probability *= 0.8
	if IsCODPayment(order) {
		score = score * 0.8 // Apply 20% reduction
	}

	// Ensure limits
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Redondear a 2 decimales
	return float64(int(score*100)) / 100, staticFactors
}

// GetStaticNegativeFactors obtiene la lista de factores negativos estáticos
func GetStaticNegativeFactors(order *domain.Order) []string {
	var factors []string

	// 1. Validación de correo
	if !isValidEmail(order.CustomerEmail) {
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
	if order.ShippingStreet2 == "" || len(order.ShippingStreet2) <= 5 {
		factors = append(factors, "Complemento de dirección")
	}

	// 7. Historial de compra
	if order.CustomerOrderCount == 0 {
		factors = append(factors, "Historial de compra")
	}

	return factors
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	// Regex simple
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsCODPayment(order *domain.Order) bool {
	// 1. Check PaymentMethodID if we have a mapping (Placeholder)
	// 2. Check Financial Details (Shopify)
	// "financial_status": "pending" AND gateway is often "manual" or "cash_on_delivery"

	// Simplificación basada en metadata o gateway
	// Si no tenemos acceso fácil a Gateway string (está en Payment struct), revisamos si podemos inferirlo.
	// El struct Order tiene PaymentMethodID.

	// Si tenemos Payments en el struct Order:
	for _, payment := range order.Payments {
		if payment.Gateway != nil {
			gw := strings.ToLower(*payment.Gateway)
			if strings.Contains(gw, "cod") || strings.Contains(gw, "cash") || strings.Contains(gw, "contra") || gw == "manual" {
				return true
			}
		}
	}

	// Si payment_method_id es específico (asumiendo 2=COD por ejemplo) - NO HARDCÓDIAR IDs sin saber.

	return false
}

// RemoveAccents normaliza el texto eliminando acentos
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}
