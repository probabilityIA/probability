package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/infra/primary/handlers/response"
)

func (h *Handlers) ListCustomerAddresses(c *gin.Context) {
	businessID, ok := h.resolveBusinessID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_id is required"})
		return
	}

	customerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || customerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := dtos.ListCustomerAddressesParams{
		CustomerID: uint(customerID),
		BusinessID: businessID,
		Page:       page,
		PageSize:   pageSize,
	}

	addresses, total, err := h.uc.ListCustomerAddresses(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := make([]response.CustomerAddressResponse, len(addresses))
	for i, a := range addresses {
		data[i] = response.AddressFromEntity(&a)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, response.CustomerAddressListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}
