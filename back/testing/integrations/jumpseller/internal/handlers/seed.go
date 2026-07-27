package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type seedProductRequest struct {
	ID    string  `json:"id"`
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
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

	h.mu.Lock()
	if body.Reset {
		h.products = make(map[int64]*product)
	}
	for _, p := range body.Products {
		id, err := strconv.ParseInt(p.ID, 10, 64)
		if err != nil {
			continue
		}
		price := p.Price
		if price == 0 {
			price = 10000
		}
		h.products[id] = &product{
			ID:     id,
			Name:   p.Name,
			SKU:    p.SKU,
			Price:  price,
			Stock:  p.Stock,
			Status: "available",
			Images: []productImage{{ID: id, URL: "https://mock.probability.test/jumpseller/" + p.SKU + ".jpg"}},
		}
		created++
	}
	total := len(h.products)
	h.mu.Unlock()

	h.logger.Info().Msgf("Jumpseller mock: seed con %d productos (total %d)", created, total)

	c.JSON(http.StatusOK, gin.H{"seeded_products": created, "total_products": total})
}
