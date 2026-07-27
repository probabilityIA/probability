package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type seedProductRequest struct {
	ID    string `json:"id"`
	SKU   string `json:"sku"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Stock int    `json:"stock"`
}

func (h *Handler) handleSeedProducts(c *gin.Context) {
	var body struct {
		Reset    bool                 `json:"reset"`
		Products []seedProductRequest `json:"products"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	created := 0
	variations := 0

	h.mu.Lock()
	if body.Reset {
		h.products = make(map[int64]*wooProduct)
		h.variations = make(map[int64][]*wooVariation)
	}
	for _, p := range body.Products {
		price := p.Price
		if price == "" {
			price = "10000.00"
		}
		stock := p.Stock

		if parent, variation, ok := strings.Cut(p.ID, ":"); ok {
			parentID, err1 := strconv.ParseInt(parent, 10, 64)
			variationID, err2 := strconv.ParseInt(variation, 10, 64)
			if err1 != nil || err2 != nil {
				continue
			}
			if _, exists := h.products[parentID]; !exists {
				h.products[parentID] = &wooProduct{
					ID: parentID, Name: p.Name, SKU: "", Type: "variable",
					Price: price, ManageStock: false, Status: "publish",
				}
			}
			h.variations[parentID] = append(h.variations[parentID], &wooVariation{
				ID: variationID, SKU: p.SKU, Price: price, StockQuantity: intPtr(stock), ManageStock: true,
			})
			variations++
			continue
		}

		id, err := strconv.ParseInt(p.ID, 10, 64)
		if err != nil {
			continue
		}
		h.products[id] = &wooProduct{
			ID: id, Name: p.Name, SKU: p.SKU, Type: "simple",
			Price: price, StockQuantity: intPtr(stock), ManageStock: true, Status: "publish",
		}
		created++
	}
	total := len(h.products)
	h.mu.Unlock()

	h.logger.Info().Msgf("WooCommerce mock: seed con %d productos y %d variaciones (total %d)", created, variations, total)

	c.JSON(http.StatusOK, gin.H{
		"seeded_products":   created,
		"seeded_variations": variations,
		"total_products":    total,
	})
}
