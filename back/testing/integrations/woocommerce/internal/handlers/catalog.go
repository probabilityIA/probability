package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) seedCatalog() {
	simple := []*wooProduct{
		{ID: 101, Name: "Camiseta Mock", SKU: "MOCK-SKU-001", Type: "simple", Price: "75000.00", StockQuantity: intPtr(10), ManageStock: true, Status: "publish"},
		{ID: 102, Name: "Gorra Mock", SKU: "MOCK-SKU-002", Type: "simple", Price: "35000.00", StockQuantity: intPtr(25), ManageStock: true, Status: "publish"},
		{ID: 103, Name: "Producto sin SKU", SKU: "", Type: "simple", Price: "12000.00", StockQuantity: intPtr(5), ManageStock: true, Status: "publish"},
	}
	for _, p := range simple {
		if p.SKU != "" {
			p.Images = []wooImage{{ID: p.ID, Src: mockImageURL(p.SKU)}}
		}
		h.products[p.ID] = p
	}

	variable := &wooProduct{ID: 110, Name: "Pantalon Mock", SKU: "MOCK-PANT", Type: "variable", Price: "90000.00", ManageStock: false, Status: "publish", Images: []wooImage{{ID: 110, Src: mockImageURL("MOCK-PANT")}}}
	h.products[variable.ID] = variable
	h.variations[variable.ID] = []*wooVariation{
		{ID: 1101, SKU: "MOCK-PANT-S", Price: "90000.00", StockQuantity: intPtr(4), ManageStock: true},
		{ID: 1102, SKU: "MOCK-PANT-M", Price: "90000.00", StockQuantity: intPtr(7), ManageStock: true},
	}

	h.nextProdID = 200
}

func intPtr(v int) *int {
	return &v
}

func (h *Handler) handleListProducts(c *gin.Context) {
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if perPage <= 0 {
		perPage = 10
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	h.mu.Lock()
	all := make([]*wooProduct, 0, len(h.products))
	for _, id := range sortedProductIDs(h.products) {
		all = append(all, h.products[id])
	}
	h.mu.Unlock()

	start := (page - 1) * perPage
	if start >= len(all) {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	end := start + perPage
	if end > len(all) {
		end = len(all)
	}

	c.Header("X-WP-Total", strconv.Itoa(len(all)))
	c.JSON(http.StatusOK, all[start:end])
}

func (h *Handler) handleGetProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	h.mu.Lock()
	product, ok := h.products[id]
	h.mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": "woocommerce_rest_product_invalid_id", "message": "producto no existe"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *Handler) handleListVariations(c *gin.Context) {
	parentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	h.mu.Lock()
	list := h.variations[parentID]
	h.mu.Unlock()

	if page > 1 {
		c.JSON(http.StatusOK, []gin.H{})
		return
	}
	if list == nil {
		list = []*wooVariation{}
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) handleUpdateVariation(c *gin.Context) {
	parentID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	variationID, _ := strconv.ParseInt(c.Param("vid"), 10, 64)

	var body struct {
		ManageStock   *bool `json:"manage_stock"`
		StockQuantity *int  `json:"stock_quantity"`
	}
	_ = c.ShouldBindJSON(&body)

	h.mu.Lock()
	var target *wooVariation
	for _, v := range h.variations[parentID] {
		if v.ID == variationID {
			target = v
			break
		}
	}
	if target != nil {
		if body.ManageStock != nil {
			target.ManageStock = *body.ManageStock
		}
		if body.StockQuantity != nil {
			target.StockQuantity = body.StockQuantity
		}
	}
	h.mu.Unlock()

	if target == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "woocommerce_rest_product_invalid_id", "message": "variacion no existe"})
		return
	}

	h.logger.Info().Msgf("WooCommerce mock: variacion %d del producto %d con stock %v", variationID, parentID, target.StockQuantity)
	c.JSON(http.StatusOK, target)
}

func sortedProductIDs(products map[int64]*wooProduct) []int64 {
	ids := make([]int64, 0, len(products))
	for id := range products {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			if ids[j] < ids[i] {
				ids[i], ids[j] = ids[j], ids[i]
			}
		}
	}
	return ids
}
