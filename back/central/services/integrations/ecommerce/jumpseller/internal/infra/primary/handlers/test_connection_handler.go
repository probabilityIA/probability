package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type testExistingResponse struct {
	Success   bool   `json:"success"`
	StoreName string `json:"store_name,omitempty"`
	StoreCode string `json:"store_code,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

func (h *jumpsellerHandler) TestExistingConnection(c *gin.Context) {
	integrationID := c.Query("integration_id")
	if integrationID == "" {
		c.JSON(http.StatusBadRequest, testExistingResponse{Success: false, Error: "integration_id es requerido"})
		return
	}

	businessID := c.GetUint("business_id")

	storeInfo, err := h.useCase.TestExistingConnection(c.Request.Context(), integrationID, businessID)
	if err != nil {
		c.JSON(http.StatusOK, testExistingResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, testExistingResponse{
		Success:   true,
		StoreName: storeInfo.Name,
		StoreCode: storeInfo.Code,
		Message:   "Conexion exitosa con la tienda " + storeInfo.Name,
	})
}
