package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type mockOrderProduct struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	VariantID int64  `json:"variant_id"`
	Name      string `json:"name"`
	SKU       string `json:"sku"`
	Price     string `json:"price"`
	Quantity  int    `json:"quantity"`
	Weight    string `json:"weight"`
	ImageURL  string `json:"image_url"`
}

type mockOrderCustomer struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Identification string `json:"identification"`
}

type mockOrderAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Number    string `json:"number"`
	Locality  string `json:"locality"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	Zipcode   string `json:"zipcode"`
}

type mockOrder struct {
	ID                   int64              `json:"id"`
	Number               int64              `json:"number"`
	Token                string             `json:"token"`
	StoreID              string             `json:"store_id"`
	Currency             string             `json:"currency"`
	Status               string             `json:"status"`
	PaymentStatus        string             `json:"payment_status"`
	ShippingStatus       string             `json:"shipping_status"`
	Subtotal             string             `json:"subtotal"`
	Discount             string             `json:"discount"`
	Total                string             `json:"total"`
	ShippingCostCustomer string             `json:"shipping_cost_customer"`
	ShippingCostOwner    string             `json:"shipping_cost_owner"`
	Weight               string             `json:"weight"`
	ContactName          string             `json:"contact_name"`
	ContactEmail         string             `json:"contact_email"`
	ContactPhone         string             `json:"contact_phone"`
	Gateway              string             `json:"gateway"`
	GatewayName          string             `json:"gateway_name"`
	Note                 string             `json:"note"`
	ShippingOption       string             `json:"shipping_option"`
	ShippingAddress      mockOrderAddress   `json:"shipping_address"`
	Customer             mockOrderCustomer  `json:"customer"`
	Products             []mockOrderProduct `json:"products"`
	CreatedAt            string             `json:"created_at"`
	UpdatedAt            string             `json:"updated_at"`
	PaidAt               string             `json:"paid_at"`
	CancelledAt          string             `json:"cancelled_at"`
	ClosedAt             string             `json:"closed_at"`
}

type orderSeed struct {
	status         string
	paymentStatus  string
	shippingStatus string
	contact        string
	sku            string
	name           string
	price          float64
	quantity       int
	diasAtras      int
}

var defaultOrderSeed = []orderSeed{
	{status: "open", paymentStatus: "pending", shippingStatus: "unpacked", contact: "Laura Gomez", sku: "8013XL", name: "Camiseta Padel White - XL / Blanco", price: 89000, quantity: 1, diasAtras: 1},
	{status: "open", paymentStatus: "paid", shippingStatus: "unpacked", contact: "Andres Pardo", sku: "8019M", name: "Pantaloneta Padel Negra - M / Negro", price: 75000, quantity: 2, diasAtras: 2},
	{status: "open", paymentStatus: "paid", shippingStatus: "shipped", contact: "Marcela Rios", sku: "52205", name: "Proteina Whey 2lb", price: 180000, quantity: 1, diasAtras: 4},
	{status: "closed", paymentStatus: "paid", shippingStatus: "delivered", contact: "Julian Mesa", sku: "8227L", name: "Camiseta Motion Black - L / Negro", price: 92000, quantity: 3, diasAtras: 7},
	{status: "cancelled", paymentStatus: "voided", shippingStatus: "unpacked", contact: "Sofia Lozano", sku: "8013XL", name: "Camiseta Padel White - XL / Blanco", price: 89000, quantity: 1, diasAtras: 9},
	{status: "open", paymentStatus: "refunded", shippingStatus: "delivered", contact: "Camilo Vera", sku: "TN-ONLY-001", name: "Producto solo en Tiendanube", price: 45000, quantity: 1, diasAtras: 12},
}

func (h *Handler) seedOrdersLocked() {
	for _, seed := range defaultOrderSeed {
		h.addOrderLocked(seed)
	}
}

func (h *Handler) addOrderLocked(seed orderSeed) *mockOrder {
	creada := time.Now().AddDate(0, 0, -seed.diasAtras)
	id := h.nextOrderID
	h.nextOrderID++

	total := seed.price * float64(seed.quantity)

	orden := &mockOrder{
		ID:                   id,
		Number:               id,
		Token:                "mock-token-" + strconv.FormatInt(id, 10),
		StoreID:              MockStoreID,
		Currency:             "COP",
		Status:               seed.status,
		PaymentStatus:        seed.paymentStatus,
		ShippingStatus:       seed.shippingStatus,
		Subtotal:             money(total),
		Discount:             "0.00",
		Total:                money(total),
		ShippingCostCustomer: "0.00",
		ShippingCostOwner:    "0.00",
		Weight:               "1.00",
		ContactName:          seed.contact,
		ContactEmail:         strings.ToLower(strings.ReplaceAll(seed.contact, " ", ".")) + "@example.com",
		ContactPhone:         "3001234567",
		Gateway:              "mock",
		GatewayName:          "Pago Mock",
		ShippingOption:       "Envio estandar",
		ShippingAddress: mockOrderAddress{
			Name:     seed.contact,
			Phone:    "3001234567",
			Address:  "Calle 100",
			Number:   "10-20",
			Locality: "Chapinero",
			City:     "Bogota",
			Province: "Bogota D.C.",
			Country:  "CO",
			Zipcode:  "110111",
		},
		Customer: mockOrderCustomer{
			ID:             id,
			Name:           seed.contact,
			Email:          strings.ToLower(strings.ReplaceAll(seed.contact, " ", ".")) + "@example.com",
			Phone:          "3001234567",
			Identification: "10" + strconv.FormatInt(id, 10),
		},
		Products: []mockOrderProduct{{
			ID:        id * 10,
			ProductID: 2001,
			VariantID: 9001,
			Name:      seed.name,
			SKU:       seed.sku,
			Price:     money(seed.price),
			Quantity:  seed.quantity,
			Weight:    "1.00",
		}},
		CreatedAt: creada.Format(time.RFC3339),
		UpdatedAt: creada.Format(time.RFC3339),
	}

	if seed.paymentStatus == "paid" || seed.paymentStatus == "refunded" {
		orden.PaidAt = creada.Format(time.RFC3339)
	}
	if seed.status == "cancelled" {
		orden.CancelledAt = creada.Format(time.RFC3339)
	}
	if seed.status == "closed" {
		orden.ClosedAt = creada.Format(time.RFC3339)
	}

	h.orders[id] = orden
	h.orderIDs = append(h.orderIDs, id)
	return orden
}

func (h *Handler) handleListOrders(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if perPage <= 0 {
		perPage = 50
	}

	desde := parseMockDate(c.Query("created_at_min"))
	hasta := parseMockDate(c.Query("created_at_max"))
	status := strings.TrimSpace(c.Query("status"))

	filtradas := make([]*mockOrder, 0, len(h.orderIDs))
	for _, id := range h.orderIDs {
		orden := h.orders[id]
		if orden == nil {
			continue
		}
		if status != "" && !strings.EqualFold(orden.Status, status) {
			continue
		}
		creada := parseMockDate(orden.CreatedAt)
		if desde != nil && creada != nil && creada.Before(*desde) {
			continue
		}
		if hasta != nil && creada != nil && creada.After(*hasta) {
			continue
		}
		filtradas = append(filtradas, orden)
	}

	inicio := (page - 1) * perPage
	if inicio >= len(filtradas) {
		c.JSON(http.StatusOK, []mockOrder{})
		return
	}
	fin := inicio + perPage
	if fin > len(filtradas) {
		fin = len(filtradas)
	}

	c.JSON(http.StatusOK, filtradas[inicio:fin])
}

func (h *Handler) handleGetOrder(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}
	orden, ok := h.orders[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, orden)
}

func (h *Handler) handleSeedOrders(c *gin.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.seedOrdersLocked()
	c.JSON(http.StatusOK, gin.H{"orders": len(h.orderIDs)})
}

func parseMockDate(value string) *time.Time {
	clean := strings.TrimSpace(value)
	if clean == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if parsed, err := time.Parse(layout, clean); err == nil {
			return &parsed
		}
	}
	return nil
}
