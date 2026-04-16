package mapper

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrations/request"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

// ============================================
// ToCreateIntegrationDTO
// ============================================

func TestToCreateIntegrationDTO_CamposBasicos(t *testing.T) {
	// Arrange
	businessID := uint(10)
	req := request.CreateIntegrationRequest{
		Name:              "Shopify Store",
		Code:              "shopify_store",
		IntegrationTypeID: 1,
		IsActive:          true,
	}
	createdByID := uint(5)

	// Act
	dto := ToCreateIntegrationDTO(req, createdByID, 0)

	// Assert
	assert.Equal(t, "Shopify Store", dto.Name)
	assert.Equal(t, "shopify_store", dto.Code)
	assert.Equal(t, uint(1), dto.IntegrationTypeID)
	assert.Equal(t, true, dto.IsActive)
	assert.Equal(t, uint(5), dto.CreatedByID)
	_ = businessID
}

func TestToCreateIntegrationDTO_BusinessIDDeContextoSobreescribe(t *testing.T) {
	// Arrange
	originalBizID := uint(1)
	req := request.CreateIntegrationRequest{
		Name:              "Test",
		Code:              "test",
		IntegrationTypeID: 1,
		BusinessID:        &originalBizID,
	}

	// Act — businessID de contexto = 20 sobreescribe el del request
	dto := ToCreateIntegrationDTO(req, 99, 20)

	// Assert — el BusinessID del contexto tiene prioridad
	assert.NotNil(t, dto.BusinessID)
	assert.Equal(t, uint(20), *dto.BusinessID)
}

func TestToCreateIntegrationDTO_SinBusinessIDDeContexto_UsaElDelRequest(t *testing.T) {
	// Arrange
	bizID := uint(7)
	req := request.CreateIntegrationRequest{
		Name:              "Test",
		Code:              "test",
		IntegrationTypeID: 1,
		BusinessID:        &bizID,
	}

	// Act — businessID de contexto = 0 no sobreescribe
	dto := ToCreateIntegrationDTO(req, 1, 0)

	// Assert — mantiene el del request
	assert.NotNil(t, dto.BusinessID)
	assert.Equal(t, uint(7), *dto.BusinessID)
}

func TestToCreateIntegrationDTO_ConConfig(t *testing.T) {
	// Arrange
	req := request.CreateIntegrationRequest{
		Name:              "Test",
		Code:              "test",
		IntegrationTypeID: 1,
		Config:            map[string]interface{}{"shop_domain": "mi-tienda.myshopify.com"},
	}

	// Act
	dto := ToCreateIntegrationDTO(req, 1, 0)

	// Assert — Config se serializa a JSON
	assert.NotNil(t, dto.Config)
	assert.Greater(t, len(dto.Config), 0)
}

// ============================================
// ToUpdateIntegrationDTO
// ============================================

func TestToUpdateIntegrationDTO_SoloCamposPresentes(t *testing.T) {
	// Arrange
	nombre := "Nuevo Nombre"
	req := request.UpdateIntegrationRequest{
		Name: &nombre,
	}

	// Act
	dto := ToUpdateIntegrationDTO(req, 5)

	// Assert
	assert.Equal(t, uint(5), dto.UpdatedByID)
	assert.NotNil(t, dto.Name)
	assert.Equal(t, "Nuevo Nombre", *dto.Name)
	assert.Nil(t, dto.Code)
	assert.Nil(t, dto.IsActive)
}

func TestToUpdateIntegrationDTO_TodosLosCampos(t *testing.T) {
	// Arrange
	nombre := "Actualizado"
	codigo := "codigo_nuevo"
	tipoID := uint(3)
	storeID := "nueva-tienda.myshopify.com"
	activo := true
	defecto := false
	desc := "Nueva descripción"
	config := map[string]interface{}{"key": "value"}
	creds := map[string]interface{}{"api_key": "abc123"}

	req := request.UpdateIntegrationRequest{
		Name:              &nombre,
		Code:              &codigo,
		IntegrationTypeID: &tipoID,
		StoreID:           &storeID,
		IsActive:          &activo,
		IsDefault:         &defecto,
		Description:       &desc,
		Config:            &config,
		Credentials:       &creds,
	}

	// Act
	dto := ToUpdateIntegrationDTO(req, 10)

	// Assert
	assert.Equal(t, "Actualizado", *dto.Name)
	assert.Equal(t, "codigo_nuevo", *dto.Code)
	assert.Equal(t, uint(3), *dto.IntegrationTypeID)
	assert.Equal(t, "nueva-tienda.myshopify.com", *dto.StoreID)
	assert.Equal(t, true, *dto.IsActive)
	assert.Equal(t, false, *dto.IsDefault)
	assert.Equal(t, "Nueva descripción", *dto.Description)
	assert.NotNil(t, dto.Config)
	assert.NotNil(t, dto.Credentials)
}

// ============================================
// ToIntegrationResponse
// ============================================

func TestToIntegrationResponse_SinIntegrationType(t *testing.T) {
	// Arrange
	integracion := &domain.Integration{
		ID:                1,
		Name:              "Mi Integración",
		Code:              "mi_int",
		IntegrationTypeID: 2,
		IsActive:          true,
		Config:            datatypes.JSON(`{"shop_domain":"test.myshopify.com"}`),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Act
	resp := ToIntegrationResponse(integracion, "")

	// Assert
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Mi Integración", resp.Name)
	assert.Equal(t, "mi_int", resp.Code)
	assert.True(t, resp.IsActive)
	assert.Nil(t, resp.IntegrationType)
}

func TestToIntegrationResponse_ConIntegrationTypeYCategoria(t *testing.T) {
	// Arrange
	categoria := &domain.IntegrationCategory{
		ID:    5,
		Code:  "ecommerce",
		Name:  "Ecommerce",
		Icon:  "shopping-cart",
		Color: "#FF6B00",
	}
	integracion := &domain.Integration{
		ID:                1,
		Name:              "Shopify",
		Code:              "shopify",
		IntegrationTypeID: 1,
		IntegrationType: &domain.IntegrationType{
			ID:       1,
			Name:     "Shopify",
			Code:     "shopify",
			ImageURL: "logos/shopify.png",
			Category: categoria,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationResponse(integracion, "https://cdn.example.com")

	// Assert
	assert.NotNil(t, resp.IntegrationType)
	assert.Equal(t, uint(1), resp.IntegrationType.ID)
	assert.Equal(t, "ecommerce", resp.Category)
	assert.Equal(t, "Ecommerce", resp.CategoryName)
	assert.Equal(t, "#FF6B00", resp.CategoryColor)
	// La URL debe ser la base + el path
	assert.Equal(t, "https://cdn.example.com/logos/shopify.png", resp.IntegrationType.ImageURL)
}

func TestToIntegrationResponse_ImageURLAbsolutaNoSeCombina(t *testing.T) {
	// Arrange
	integracion := &domain.Integration{
		ID:   1,
		Name: "Test",
		IntegrationType: &domain.IntegrationType{
			ID:       1,
			Name:     "Test",
			ImageURL: "https://external.cdn.com/logo.png",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationResponse(integracion, "https://cdn.example.com")

	// Assert — URL absoluta se usa directamente
	assert.Equal(t, "https://external.cdn.com/logo.png", resp.IntegrationType.ImageURL)
}

func TestToIntegrationResponse_ConfigInvalidaRetornaMapaVacio(t *testing.T) {
	// Arrange
	integracion := &domain.Integration{
		ID:        1,
		Name:      "Test",
		Config:    datatypes.JSON(`invalid json`),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationResponse(integracion, "")

	// Assert — config inválida retorna mapa vacío, no falla
	assert.NotNil(t, resp.Config)
	assert.Empty(t, resp.Config)
}

// ============================================
// ToIntegrationListResponse
// ============================================

func TestToIntegrationListResponse_CalculaPaginacion(t *testing.T) {
	// Arrange
	integraciones := []*domain.Integration{
		{ID: 1, Name: "Int 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Int 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Act — 25 total, 10 por página = 3 páginas
	resp := ToIntegrationListResponse(integraciones, 25, 1, 10, "")

	// Assert
	assert.True(t, resp.Success)
	assert.Equal(t, int64(25), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, 3, resp.TotalPages)
	assert.Len(t, resp.Data, 2)
}

func TestToIntegrationListResponse_PageSizeCeroUsaDefault(t *testing.T) {
	// Arrange
	integraciones := []*domain.Integration{
		{ID: 1, Name: "Int 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	// Act — pageSize 0 usa default 10
	resp := ToIntegrationListResponse(integraciones, 15, 1, 0, "")

	// Assert
	assert.Equal(t, 10, resp.PageSize)
	assert.Equal(t, 2, resp.TotalPages) // ceil(15/10) = 2
}

func TestToIntegrationListResponse_TotalExactamenteDivisible(t *testing.T) {
	// Arrange
	integraciones := []*domain.Integration{}

	// Act — 20 total, 10 por página = exactamente 2 páginas (sin resto)
	resp := ToIntegrationListResponse(integraciones, 20, 1, 10, "")

	// Assert
	assert.Equal(t, 2, resp.TotalPages)
}

// ============================================
// ToIntegrationFilters
// ============================================

func TestToIntegrationFilters_CamposBasicos(t *testing.T) {
	// Arrange
	bizID := uint(5)
	activo := true
	req := request.GetIntegrationsRequest{
		Page:       2,
		PageSize:   20,
		BusinessID: &bizID,
		IsActive:   &activo,
	}

	// Act
	filters := ToIntegrationFilters(req)

	// Assert
	assert.Equal(t, 2, filters.Page)
	assert.Equal(t, 20, filters.PageSize)
	assert.Equal(t, &bizID, filters.BusinessID)
	assert.Equal(t, &activo, filters.IsActive)
}

func TestToIntegrationFilters_PerPageComoFallback(t *testing.T) {
	// Arrange — usa per_page cuando page_size no está
	req := request.GetIntegrationsRequest{
		Page:    1,
		PerPage: 15,
	}

	// Act
	filters := ToIntegrationFilters(req)

	// Assert
	assert.Equal(t, 15, filters.PageSize)
}

func TestToIntegrationFilters_SinPageSizeNiPerPage_UsaDefault(t *testing.T) {
	// Arrange
	req := request.GetIntegrationsRequest{
		Page: 1,
	}

	// Act
	filters := ToIntegrationFilters(req)

	// Assert — default de 10
	assert.Equal(t, 10, filters.PageSize)
}
