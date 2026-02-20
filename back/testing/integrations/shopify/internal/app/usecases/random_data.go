package usecases

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/shopify/internal/domain"
)

var (
	firstNames = []string{
		"María", "Juan", "Ana", "Carlos", "Laura", "Pedro", "Sofía", "Luis",
		"Carmen", "Miguel", "Isabel", "Javier", "Patricia", "Francisco", "Lucía",
		"Antonio", "Elena", "Manuel", "Marta", "José", "Cristina", "David",
		"Paula", "Daniel", "Andrea", "Alejandro", "Natalia", "Roberto",
	}

	lastNames = []string{
		"García", "Rodríguez", "González", "Fernández", "López", "Martínez",
		"Sánchez", "Pérez", "Gómez", "Martín", "Jiménez", "Ruiz", "Hernández",
		"Díaz", "Moreno", "Álvarez", "Muñoz", "Romero", "Alonso", "Gutiérrez",
		"Navarro", "Torres", "Domínguez", "Vázquez", "Ramos", "Gil", "Ramírez",
	}

	cities = []string{
		"Bogotá", "Medellín", "Cali", "Barranquilla", "Cartagena", "Bucaramanga",
		"Pereira", "Santa Marta", "Manizales", "Armenia", "Villavicencio",
		"Pasto", "Ibagué", "Neiva", "Valledupar", "Montería", "Sincelejo",
	}

	provinces = []string{
		"Cundinamarca", "Antioquia", "Valle del Cauca", "Atlántico", "Bolívar",
		"Santander", "Risaralda", "Magdalena", "Caldas", "Quindío", "Meta",
		"Nariño", "Tolima", "Huila", "Cesar", "Córdoba", "Sucre",
	}

	streets = []string{
		"Calle", "Carrera", "Avenida", "Diagonal", "Transversal",
	}

	// Productos reales del catálogo de Softpymes (itemCode debe coincidir)
	productNames = []string{
		"Proteína Aislada (ISO) - 2 Lb (910g) - Vainilla",
		"Proteína Whey - 2lb (910g) - Vainilla",
		"Proteína Whey - 2lb (910g) - Chocolate",
		"Creatina Monohidrato - 300g",
		"Proteína Aislada (ISO) - 2 Lb (910g) - Chocolate",
		"Proteína Vegetal - 2 Libras (910g) - Vainilla",
		"Proteína Vegetal - 2 Libras (910g) - Chocolate",
		"Multivitaminico - Gomas",
		"Omega 3 + prebioticos - Gomas",
		"Citrato de Magnesio Limon - 210g",
		"BCAAs sabor limon mandarino - 300g",
		"PR - 600g",
		"Colágeno Hidrolizado - 300g",
		"Pancakes de Proteina (770)",
		"Creatina Monohidrato - 100g",
	}

	// SKUs que coinciden con los itemCode de Softpymes
	productSKUs = []string{
		"PT01001", "PT01002", "PT01003", "PT01004", "PT01005",
		"PT01006", "PT01007", "PT02038", "PT02039", "PT02041",
		"PT02043", "PT02044", "PT01015", "PT01016", "PT02050",
	}

	vendors = []string{
		"Mi Tienda", "Fashion Store", "Tech Solutions", "Home & Living",
		"Sports World", "Electronics Plus", "Furniture House", "Style Co",
		"Digital Store", "Quality Goods",
	}

	currencies = []string{"COP"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomDataGenerator genera datos aleatorios para órdenes
type RandomDataGenerator struct{}

// NewRandomDataGenerator crea un nuevo generador de datos aleatorios
func NewRandomDataGenerator() *RandomDataGenerator {
	return &RandomDataGenerator{}
}

// GenerateCustomer genera un cliente aleatorio
func (g *RandomDataGenerator) GenerateCustomer() *domain.Customer {
	// Usar datos específicos de Carlos Arrieta para testing
	firstName := "Carlos"
	lastName := "Arrieta"
	email := "carlos.arrieta@example.com"
	phone := "+573023406789"

	now := time.Now()
	createdAt := now.AddDate(0, 0, -rand.Intn(365))

	return &domain.Customer{
		ID:                        int64(rand.Intn(999999999) + 100000000),
		Email:                     email,
		AcceptsMarketing:          true,
		CreatedAt:                 createdAt,
		UpdatedAt:                 now,
		FirstName:                 firstName,
		LastName:                  lastName,
		State:                     "enabled",
		Note:                      nil,
		VerifiedEmail:             true,
		MultipassIdentifier:       nil,
		TaxExempt:                 false,
		Phone:                     g.stringPtr(phone),
		Tags:                      "",
		Currency:                  "COP",
		AcceptsMarketingUpdatedAt: g.timePtr(createdAt),
		MarketingOptInLevel:       nil,
		AdminGraphQLAPIID:         fmt.Sprintf("gid://shopify/Customer/%d", rand.Intn(999999999)+100000000),
		DefaultAddress:            g.GenerateAddress(),
	}
}

// GenerateAddress genera una dirección aleatoria
func (g *RandomDataGenerator) GenerateAddress() *domain.Address {
	city := cities[rand.Intn(len(cities))]
	province := provinces[rand.Intn(len(provinces))]
	streetType := streets[rand.Intn(len(streets))]
	streetNumber := rand.Intn(200) + 1
	address1 := fmt.Sprintf("%s %d", streetType, streetNumber)

	var address2 *string
	if rand.Float32() < 0.4 {
		apt := fmt.Sprintf("Apto %d", rand.Intn(500)+1)
		address2 = &apt
	}

	return &domain.Address{
		FirstName:    g.randomChoice(firstNames),
		LastName:     g.randomChoice(lastNames),
		Company:      nil,
		Address1:     address1,
		Address2:     address2,
		City:         city,
		Province:     province,
		Country:      "Colombia",
		Zip:          fmt.Sprintf("%05d", rand.Intn(99999)),
		Phone:        g.stringPtr(fmt.Sprintf("+57%d", rand.Intn(9999999999)+3000000000)),
		Name:         "",
		ProvinceCode: g.stringPtr(getProvinceCode(province)),
		CountryCode:  "CO",
		CountryName:  g.stringPtr("Colombia"),
		Default:      nil,
		Latitude:     g.floatPtr(4.0 + rand.Float64()*2.0),
		Longitude:    g.floatPtr(-76.0 - rand.Float64()*2.0),
	}
}

// GenerateLineItems genera items de línea aleatorios
func (g *RandomDataGenerator) GenerateLineItems(count int) []domain.LineItem {
	items := make([]domain.LineItem, 0, count)

	for i := 0; i < count; i++ {
		// Emparejar nombre y SKU del mismo índice (catálogo Softpymes)
		idx := rand.Intn(len(productSKUs))
		productName := productNames[idx]
		sku := productSKUs[idx]
		// Precios realistas en COP (entre 25.300 y 152.900 como en catálogo Softpymes)
		price := float64(rand.Intn(127600)+25300) / 1.0
		quantity := rand.Intn(3) + 1

		item := domain.LineItem{
			ID:                         int64(rand.Intn(999999999) + 100000000),
			AdminGraphQLAPIID:          fmt.Sprintf("gid://shopify/LineItem/%d", rand.Intn(999999999)+100000000),
			FulfillableQuantity:        quantity,
			FulfillmentService:         g.stringPtr("manual"),
			FulfillmentStatus:          nil,
			GiftCard:                   false,
			Grams:                      rand.Intn(5000) + 100,
			Name:                       fmt.Sprintf("%s - %s", productName, g.randomChoice([]string{"S", "M", "L", "XL"})),
			Price:                      fmt.Sprintf("%.2f", price),
			PriceSet:                   g.generateMoneySet(fmt.Sprintf("%.2f", price), g.randomChoice(currencies)),
			ProductExists:              true,
			ProductID:                  int64(rand.Intn(999999999) + 100000000),
			Properties:                 []domain.Property{},
			Quantity:                   quantity,
			RequiresShipping:           rand.Float32() < 0.8,
			SKU:                        sku,
			Taxable:                    true,
			Title:                      productName,
			TotalDiscount:              "0.00",
			TotalDiscountSet:           g.generateMoneySet("0.00", g.randomChoice(currencies)),
			VariantID:                  int64(rand.Intn(999999999) + 100000000),
			VariantInventoryManagement: g.stringPtr("shopify"),
			VariantTitle:               g.stringPtr(g.randomChoice([]string{"S", "M", "L", "XL", "Default"})),
			Vendor:                     g.stringPtr(vendors[rand.Intn(len(vendors))]),
			TaxLines:                   []domain.TaxLine{},
			Duties:                     []domain.Duty{},
			DiscountAllocations:        []domain.DiscountAllocation{},
		}

		items = append(items, item)
	}

	return items
}

// GenerateOrderNumber genera un número de orden único
func (g *RandomDataGenerator) GenerateOrderNumber() string {
	return fmt.Sprintf("#%d", rand.Intn(99999)+1000)
}

// randomChoice elige un elemento aleatorio de un slice
func (g *RandomDataGenerator) randomChoice(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

// Helper functions

func toLower(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func getProvinceCode(province string) string {
	codes := map[string]string{
		"Cundinamarca":    "DC",
		"Antioquia":       "ANT",
		"Valle del Cauca": "VAL",
		"Atlántico":       "ATL",
		"Bolívar":         "BOL",
		"Santander":       "SAN",
		"Risaralda":       "RIS",
		"Magdalena":       "MAG",
		"Caldas":          "CAL",
		"Quindío":         "QUI",
		"Meta":            "MET",
		"Nariño":          "NAR",
		"Tolima":          "TOL",
		"Huila":           "HUI",
		"Cesar":           "CES",
		"Córdoba":         "COR",
		"Sucre":           "SUC",
	}
	if code, ok := codes[province]; ok {
		return code
	}
	return "DC"
}

func (g *RandomDataGenerator) generateMoneySet(amount, currency string) *domain.MoneySet {
	return &domain.MoneySet{
		ShopMoney: domain.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
		PresentmentMoney: domain.Money{
			Amount:       amount,
			CurrencyCode: currency,
		},
	}
}

func (g *RandomDataGenerator) stringPtr(s string) *string {
	return &s
}

func (g *RandomDataGenerator) floatPtr(f float64) *float64 {
	return &f
}

func (g *RandomDataGenerator) timePtr(t time.Time) *time.Time {
	return &t
}













