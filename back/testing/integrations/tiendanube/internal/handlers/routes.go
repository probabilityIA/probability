package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.TrimSpace(c.GetHeader("User-Agent")) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User-Agent header is required"})
			return
		}

		token := bearerToken(c.GetHeader("Authentication"))
		if token == "" {
			token = bearerToken(c.GetHeader("Authorization"))
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Unauthorized", "description": "missing access token"})
			return
		}

		if c.Param("store_id") == "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}

		c.Next()
	}
}

func bearerToken(raw string) string {
	parts := strings.Fields(strings.TrimSpace(raw))
	if len(parts) != 2 {
		return ""
	}
	if !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}

func (h *Handler) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "tiendanube-mock"})
}

func (h *Handler) handleStore(c *gin.Context) {
	storeID, _ := strconv.ParseInt(c.Param("store_id"), 10, 64)
	c.JSON(http.StatusOK, gin.H{
		"id":            storeID,
		"name":          localized{"es": "Tienda Mock Probability"},
		"url":           "https://mock-probability.mitiendanube.com",
		"country":       "CO",
		"currency":      "COP",
		"main_language": "es",
		"email":         "mock@probabilityia.com.co",
	})
}

func (h *Handler) handleListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "30"))
	if perPage < 1 || perPage > 200 {
		perPage = 30
	}

	h.mu.Lock()
	todos := h.snapshotLocked()
	h.mu.Unlock()

	if ids := strings.TrimSpace(c.Query("ids")); ids != "" {
		buscados := make(map[int64]bool)
		for _, raw := range strings.Split(ids, ",") {
			if n, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64); err == nil {
				buscados[n] = true
			}
		}
		filtrados := make([]product, 0, len(buscados))
		for _, p := range todos {
			if buscados[p.ID] {
				filtrados = append(filtrados, p)
			}
		}
		todos = filtrados
	}

	if q := strings.ToLower(strings.TrimSpace(c.Query("q"))); q != "" {
		filtrados := make([]product, 0)
		for _, p := range todos {
			if matchesQuery(p, q) {
				filtrados = append(filtrados, p)
			}
		}
		todos = filtrados
	}

	total := len(todos)
	inicio := (page - 1) * perPage
	if inicio > total {
		inicio = total
	}
	fin := inicio + perPage
	if fin > total {
		fin = total
	}

	c.Header("x-total-count", strconv.Itoa(total))
	c.Header("x-rate-limit-limit", "40")
	c.Header("x-rate-limit-remaining", "39")
	c.JSON(http.StatusOK, todos[inicio:fin])
}

func matchesQuery(p product, q string) bool {
	for _, value := range p.Name {
		if strings.Contains(strings.ToLower(value), q) {
			return true
		}
	}
	for _, v := range p.Variants {
		if strings.Contains(strings.ToLower(v.SKU), q) {
			return true
		}
		if strings.Contains(strings.ToLower(v.Barcode), q) {
			return true
		}
	}
	return false
}

func (h *Handler) handleGetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found, ok := h.products[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}
	c.JSON(http.StatusOK, *found)
}

type createVariantPayload struct {
	SKU     string      `json:"sku"`
	Barcode string      `json:"barcode"`
	Price   interface{} `json:"price"`
	Stock   interface{} `json:"stock"`
	Weight  *float64    `json:"weight"`
	Depth   *float64    `json:"depth"`
	Width   *float64    `json:"width"`
	Height  *float64    `json:"height"`
}

type createProductPayload struct {
	Name        map[string]string      `json:"name"`
	Description map[string]string      `json:"description"`
	Published   *bool                  `json:"published"`
	Variants    []createVariantPayload `json:"variants"`
	Images      []struct {
		Src string `json:"src"`
	} `json:"images"`
}

func (h *Handler) handleCreateProduct(c *gin.Context) {
	if !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"code": 415, "message": "Unsupported Media Type"})
		return
	}

	var payload createProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problems parsing JSON"})
		return
	}

	if len(payload.Name) == 0 || strings.TrimSpace(firstValue(payload.Name)) == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"name": []string{"can't be blank"}})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	nuevo := &product{
		ID:          h.nextProductID,
		Name:        payload.Name,
		Description: payload.Description,
		Published:   payload.Published == nil || *payload.Published,
		Variants:    make([]variant, 0, 1),
		Images:      make([]image, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	h.nextProductID++

	if len(payload.Variants) == 0 {
		payload.Variants = []createVariantPayload{{}}
	}

	for _, v := range payload.Variants {
		if sku := strings.TrimSpace(v.SKU); sku != "" && h.skuExistsLocked(sku) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"sku": []string{"has already been taken"}})
			return
		}
		nueva := variant{
			ID:        h.nextVariantID,
			ProductID: nuevo.ID,
			SKU:       strings.TrimSpace(v.SKU),
			Barcode:   strings.TrimSpace(v.Barcode),
			Price:     money(toFloat(v.Price)),
			Weight:    money(deref(v.Weight)),
			Depth:     money(deref(v.Depth)),
			Width:     money(deref(v.Width)),
			Height:    money(deref(v.Height)),
		}
		h.nextVariantID++
		if v.Stock == nil {
			nueva.Stock = nil
			nueva.StockManagement = false
		} else {
			nueva.Stock = int(toFloat(v.Stock))
			nueva.StockManagement = true
		}
		nuevo.Variants = append(nuevo.Variants, nueva)
	}

	for _, img := range payload.Images {
		if strings.TrimSpace(img.Src) == "" {
			continue
		}
		nuevo.Images = append(nuevo.Images, image{
			ID:        h.nextImageID,
			ProductID: nuevo.ID,
			Src:       img.Src,
			Position:  len(nuevo.Images) + 1,
		})
		h.nextImageID++
	}

	h.products[nuevo.ID] = nuevo
	h.order = append(h.order, nuevo.ID)

	c.JSON(http.StatusCreated, *nuevo)
}

type updateProductPayload struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
}

func (h *Handler) handleUpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}

	var payload updateProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problems parsing JSON"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found, ok := h.products[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}

	if len(payload.Name) > 0 {
		found.Name = payload.Name
	}
	if len(payload.Description) > 0 {
		found.Description = payload.Description
	}
	found.UpdatedAt = time.Now().Format(time.RFC3339)

	c.JSON(http.StatusOK, *found)
}

type updateVariantPayload struct {
	Price   interface{} `json:"price"`
	Barcode *string     `json:"barcode"`
	Stock   interface{} `json:"stock"`
	Weight  *float64    `json:"weight"`
	Depth   *float64    `json:"depth"`
	Width   *float64    `json:"width"`
	Height  *float64    `json:"height"`
}

func (h *Handler) handleUpdateVariant(c *gin.Context) {
	productID, perr := strconv.ParseInt(c.Param("id"), 10, 64)
	variantID, verr := strconv.ParseInt(c.Param("vid"), 10, 64)
	if perr != nil || verr != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}

	raw := make(map[string]interface{})
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problems parsing JSON"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	found, ok := h.products[productID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
		return
	}

	for i := range found.Variants {
		if found.Variants[i].ID != variantID {
			continue
		}

		if value, present := raw["price"]; present && value != nil {
			found.Variants[i].Price = money(toFloat(value))
		}
		if value, present := raw["barcode"]; present {
			if text, ok := value.(string); ok {
				found.Variants[i].Barcode = strings.TrimSpace(text)
			}
		}
		for key, target := range map[string]*string{
			"weight": &found.Variants[i].Weight,
			"depth":  &found.Variants[i].Depth,
			"width":  &found.Variants[i].Width,
			"height": &found.Variants[i].Height,
		} {
			if value, present := raw[key]; present && value != nil {
				*target = money(toFloat(value))
			}
		}
		if value, present := raw["stock"]; present {
			if value == nil || value == "" {
				found.Variants[i].Stock = nil
				found.Variants[i].StockManagement = false
			} else {
				found.Variants[i].Stock = int(toFloat(value))
				found.Variants[i].StockManagement = true
			}
		}

		found.UpdatedAt = time.Now().Format(time.RFC3339)
		c.JSON(http.StatusOK, found.Variants[i])
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
}

func (h *Handler) handleReset(c *gin.Context) {
	h.mu.Lock()
	h.products = make(map[int64]*product)
	h.order = make([]int64, 0)
	h.nextProductID = 2001
	h.nextVariantID = 9001
	h.nextImageID = 7001
	h.seedDefaultsLocked()
	total := len(h.order)
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"success": true, "products": total})
}

func (h *Handler) handleSeedProducts(c *gin.Context) {
	var payload struct {
		Products []struct {
			Name  string  `json:"name"`
			SKU   string  `json:"sku"`
			Price float64 `json:"price"`
			Stock *int    `json:"stock"`
		} `json:"products"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Problems parsing JSON"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	creados := 0
	for _, item := range payload.Products {
		if strings.TrimSpace(item.SKU) == "" || h.skuExistsLocked(item.SKU) {
			continue
		}
		h.addProductLocked(item.Name, item.SKU, item.Price, item.Stock)
		creados++
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "created": creados, "total": len(h.order)})
}

func firstValue(m map[string]string) string {
	for _, lang := range []string{"es", "pt", "en"} {
		if v, ok := m[lang]; ok {
			return v
		}
	}
	for _, v := range m {
		return v
	}
	return ""
}

func toFloat(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0
		}
		return parsed
	}
	return 0
}

func deref(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func (h *Handler) snapshotLocked() []product {
	out := make([]product, 0, len(h.order))
	for _, id := range h.order {
		if p, ok := h.products[id]; ok {
			out = append(out, *p)
		}
	}
	return out
}

func (h *Handler) skuExistsLocked(sku string) bool {
	limpio := strings.ToLower(strings.TrimSpace(sku))
	for _, p := range h.products {
		for _, v := range p.Variants {
			if strings.ToLower(v.SKU) == limpio {
				return true
			}
		}
	}
	return false
}
