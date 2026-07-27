package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type shopifyVariant struct {
	ID                int64  `json:"id"`
	ProductID         int64  `json:"product_id"`
	Title             string `json:"title"`
	SKU               string `json:"sku"`
	Price             string `json:"price"`
	InventoryItemID   int64  `json:"inventory_item_id"`
	InventoryQuantity int    `json:"inventory_quantity"`
	InventoryManagem  string `json:"inventory_management"`
}

type shopifyProduct struct {
	ID       int64             `json:"id"`
	Title    string            `json:"title"`
	Handle   string            `json:"handle"`
	Status   string            `json:"status"`
	Variants []*shopifyVariant `json:"variants"`
}

type shopifyLocation struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type catalog struct {
	mu         sync.Mutex
	products   map[int64]*shopifyProduct
	locations  []*shopifyLocation
	levels     map[string]int
	nextProdID int64
	nextVarID  int64
	nextItemID int64
}

var shopifyCatalog = newCatalog()

func newCatalog() *catalog {
	c := &catalog{
		products:   make(map[int64]*shopifyProduct),
		levels:     make(map[string]int),
		nextProdID: 9000000001,
		nextVarID:  8000000001,
		nextItemID: 7000000001,
	}
	c.locations = []*shopifyLocation{
		{ID: 6000000001, Name: "Bodega Mock Principal", Active: true},
		{ID: 6000000002, Name: "Bodega Mock Secundaria", Active: true},
	}
	c.seed("Camiseta Mock Shopify", "SHOP-MOCK-001", "75000.00", 10)
	c.seed("Gorra Mock Shopify", "SHOP-MOCK-002", "35000.00", 25)
	c.seed("Producto Shopify sin SKU", "", "12000.00", 3)
	return c
}

func (c *catalog) seed(title, sku, price string, qty int) *shopifyProduct {
	productID := c.nextProdID
	c.nextProdID++
	variantID := c.nextVarID
	c.nextVarID++
	itemID := c.nextItemID
	c.nextItemID++

	product := &shopifyProduct{
		ID:     productID,
		Title:  title,
		Handle: sku,
		Status: "active",
		Variants: []*shopifyVariant{
			{
				ID:                variantID,
				ProductID:         productID,
				Title:             "Default Title",
				SKU:               sku,
				Price:             price,
				InventoryItemID:   itemID,
				InventoryQuantity: qty,
				InventoryManagem:  "shopify",
			},
		},
	}
	c.products[productID] = product
	c.levels[levelKey(c.locations[0].ID, itemID)] = qty
	return product
}

func levelKey(locationID, inventoryItemID int64) string {
	return strconv.FormatInt(locationID, 10) + ":" + strconv.FormatInt(inventoryItemID, 10)
}

func (c *catalog) sortedProducts() []*shopifyProduct {
	ids := make([]int64, 0, len(c.products))
	for id := range c.products {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[j] < ids[i] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}
	out := make([]*shopifyProduct, 0, len(ids))
	for _, id := range ids {
		out = append(out, c.products[id])
	}
	return out
}

func requireToken(c *gin.Context) bool {
	if c.GetHeader("X-Shopify-Access-Token") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": "[API] Clave de API no valida o token de acceso (verificar encabezado)"})
		return false
	}
	return true
}

func (h *Handler) handleListProducts(c *gin.Context) {
	if !requireToken(c) {
		return
	}

	shopifyCatalog.mu.Lock()
	list := shopifyCatalog.sortedProducts()
	shopifyCatalog.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"products": list})
}

func (h *Handler) handleGetProduct(c *gin.Context) {
	if !requireToken(c) {
		return
	}

	raw := strings.TrimSuffix(c.Param("product_id"), ".json")
	if base, _, ok := strings.Cut(raw, ":"); ok {
		raw = base
	}
	id, _ := strconv.ParseInt(raw, 10, 64)

	shopifyCatalog.mu.Lock()
	product, ok := shopifyCatalog.products[id]
	shopifyCatalog.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"errors": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	if !requireToken(c) {
		return
	}

	var body struct {
		Product struct {
			Title    string `json:"title"`
			Status   string `json:"status"`
			Variants []struct {
				SKU               string `json:"sku"`
				Price             string `json:"price"`
				InventoryQuantity int    `json:"inventory_quantity"`
			} `json:"variants"`
		} `json:"product"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	sku := ""
	price := "0.00"
	qty := 0
	if len(body.Product.Variants) > 0 {
		sku = body.Product.Variants[0].SKU
		if body.Product.Variants[0].Price != "" {
			price = body.Product.Variants[0].Price
		}
		qty = body.Product.Variants[0].InventoryQuantity
	}

	shopifyCatalog.mu.Lock()
	product := shopifyCatalog.seed(body.Product.Title, sku, price, qty)
	shopifyCatalog.mu.Unlock()

	h.logger.Info().Msgf("Shopify mock: producto creado %d (sku %s)", product.ID, sku)
	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func (h *Handler) handleGetLocations(c *gin.Context) {
	if !requireToken(c) {
		return
	}

	shopifyCatalog.mu.Lock()
	list := shopifyCatalog.locations
	shopifyCatalog.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"locations": list})
}

func (h *Handler) handleSetInventoryLevel(c *gin.Context) {
	if !requireToken(c) {
		return
	}

	var body struct {
		LocationID      int64 `json:"location_id"`
		InventoryItemID int64 `json:"inventory_item_id"`
		Available       int   `json:"available"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}
	if body.LocationID == 0 || body.InventoryItemID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": "location_id e inventory_item_id son requeridos"})
		return
	}

	shopifyCatalog.mu.Lock()
	shopifyCatalog.levels[levelKey(body.LocationID, body.InventoryItemID)] = body.Available
	found := false
	for _, p := range shopifyCatalog.products {
		for _, v := range p.Variants {
			if v.InventoryItemID == body.InventoryItemID {
				v.InventoryQuantity = body.Available
				found = true
			}
		}
	}
	shopifyCatalog.mu.Unlock()

	h.logger.Info().Msgf("Shopify mock: inventory_item %d en location %d queda en %d (variante encontrada: %t)",
		body.InventoryItemID, body.LocationID, body.Available, found)

	c.JSON(http.StatusOK, gin.H{
		"inventory_level": gin.H{
			"inventory_item_id": body.InventoryItemID,
			"location_id":       body.LocationID,
			"available":         body.Available,
		},
	})
}
