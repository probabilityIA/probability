package mapper

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/core/internal/infra/primary/handlers/handlerintegrationtype/request"
	"github.com/stretchr/testify/assert"
)

// ============================================
// ToIntegrationCategoryResponse
// ============================================

func TestToIntegrationCategoryResponse_MapeaCorrectamente(t *testing.T) {
	// Arrange
	categoria := &domain.IntegrationCategory{
		ID:           1,
		Code:         "ecommerce",
		Name:         "Ecommerce",
		Description:  "Plataformas de comercio electrónico",
		Icon:         "shopping-cart",
		Color:        "#FF6B00",
		DisplayOrder: 1,
		IsActive:     true,
		IsVisible:    true,
	}

	// Act
	resp := ToIntegrationCategoryResponse(categoria)

	// Assert
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "ecommerce", resp.Code)
	assert.Equal(t, "Ecommerce", resp.Name)
	assert.Equal(t, "Plataformas de comercio electrónico", resp.Description)
	assert.Equal(t, "shopping-cart", resp.Icon)
	assert.Equal(t, "#FF6B00", resp.Color)
	assert.Equal(t, 1, resp.DisplayOrder)
	assert.True(t, resp.IsActive)
	assert.True(t, resp.IsVisible)
}

func TestToIntegrationCategoryResponse_CamposVacios(t *testing.T) {
	// Arrange
	categoria := &domain.IntegrationCategory{
		ID:   2,
		Code: "messaging",
		Name: "Mensajería",
	}

	// Act
	resp := ToIntegrationCategoryResponse(categoria)

	// Assert — campos sin valor quedan en zero value
	assert.Equal(t, uint(2), resp.ID)
	assert.Equal(t, "messaging", resp.Code)
	assert.Equal(t, "", resp.Icon)
	assert.Equal(t, "", resp.Color)
	assert.Equal(t, 0, resp.DisplayOrder)
}

// ============================================
// ToIntegrationTypeResponse
// ============================================

func TestToIntegrationTypeResponse_SinCategoria(t *testing.T) {
	// Arrange
	tipo := &domain.IntegrationType{
		ID:          1,
		Name:        "Shopify",
		Code:        "shopify",
		Description: "Plataforma de ecommerce",
		Icon:        "shopify-icon",
		ImageURL:    "",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Act
	resp := ToIntegrationTypeResponse(tipo, "")

	// Assert
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, "Shopify", resp.Name)
	assert.Equal(t, "shopify", resp.Code)
	assert.Equal(t, "Plataforma de ecommerce", resp.Description)
	assert.True(t, resp.IsActive)
	assert.Equal(t, "", resp.ImageURL)
	assert.Nil(t, resp.Category)
}

func TestToIntegrationTypeResponse_ConCategoria(t *testing.T) {
	// Arrange
	categoria := &domain.IntegrationCategory{
		ID:    5,
		Code:  "ecommerce",
		Name:  "Ecommerce",
		Icon:  "cart",
		Color: "#FF6B00",
	}
	tipo := &domain.IntegrationType{
		ID:        1,
		Name:      "Shopify",
		Code:      "shopify",
		Category:  categoria,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationTypeResponse(tipo, "")

	// Assert
	assert.NotNil(t, resp.Category)
	assert.Equal(t, uint(5), resp.Category.ID)
	assert.Equal(t, "ecommerce", resp.Category.Code)
	assert.Equal(t, "Ecommerce", resp.Category.Name)
}

func TestToIntegrationTypeResponse_ImageURLRelativaSeCombina(t *testing.T) {
	// Arrange
	tipo := &domain.IntegrationType{
		ID:        1,
		Name:      "Shopify",
		Code:      "shopify",
		ImageURL:  "logos/shopify.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationTypeResponse(tipo, "https://cdn.example.com")

	// Assert
	assert.Equal(t, "https://cdn.example.com/logos/shopify.png", resp.ImageURL)
}

func TestToIntegrationTypeResponse_ImageURLAbsolutaNoSeCombina(t *testing.T) {
	// Arrange
	tipo := &domain.IntegrationType{
		ID:        1,
		Name:      "Shopify",
		Code:      "shopify",
		ImageURL:  "https://external.cdn.com/shopify.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	resp := ToIntegrationTypeResponse(tipo, "https://cdn.example.com")

	// Assert — URL absoluta no se combina con la base
	assert.Equal(t, "https://external.cdn.com/shopify.png", resp.ImageURL)
}

func TestToIntegrationTypeResponse_SinImageURLBase_UsaDirectamente(t *testing.T) {
	// Arrange
	tipo := &domain.IntegrationType{
		ID:        1,
		Name:      "Factus",
		Code:      "factus",
		ImageURL:  "logos/factus.png",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act — sin base URL
	resp := ToIntegrationTypeResponse(tipo, "")

	// Assert — sin base, usa la ruta tal cual
	assert.Equal(t, "logos/factus.png", resp.ImageURL)
}

func TestToIntegrationTypeResponse_FormatoDeFechas(t *testing.T) {
	// Arrange
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	tipo := &domain.IntegrationType{
		ID:        1,
		Name:      "Test",
		Code:      "test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	resp := ToIntegrationTypeResponse(tipo, "")

	// Assert — formato RFC3339 con zona horaria
	assert.Equal(t, "2024-06-15T10:30:00Z", resp.CreatedAt)
	assert.Equal(t, "2024-06-15T10:30:00Z", resp.UpdatedAt)
}

// ============================================
// ToCreateIntegrationTypeDTO
// ============================================

func TestToCreateIntegrationTypeDTO_CamposBasicos(t *testing.T) {
	// Arrange
	req := request.CreateIntegrationTypeRequest{
		Name:       "Siigo",
		Code:       "siigo",
		CategoryID: 1,
		IsActive:   true,
	}

	// Act
	dto := ToCreateIntegrationTypeDTO(req)

	// Assert
	assert.Equal(t, "Siigo", dto.Name)
	assert.Equal(t, "siigo", dto.Code)
	assert.True(t, dto.IsActive)
}

func TestToCreateIntegrationTypeDTO_ConConfigSchema(t *testing.T) {
	// Arrange
	configSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"shop_domain": map[string]interface{}{"type": "string"},
		},
	}
	req := request.CreateIntegrationTypeRequest{
		Name:         "Shopify",
		Code:         "shopify",
		ConfigSchema: configSchema,
	}

	// Act
	dto := ToCreateIntegrationTypeDTO(req)

	// Assert — ConfigSchema se serializa a JSON
	assert.NotNil(t, dto.ConfigSchema)
	assert.Greater(t, len(dto.ConfigSchema), 0)
}

func TestToCreateIntegrationTypeDTO_SinConfigSchema(t *testing.T) {
	// Arrange
	req := request.CreateIntegrationTypeRequest{
		Name: "Test",
		Code: "test",
	}

	// Act
	dto := ToCreateIntegrationTypeDTO(req)

	// Assert — sin schema, queda en nil/vacío
	assert.Nil(t, dto.ConfigSchema)
}

// ============================================
// ToUpdateIntegrationTypeDTO
// ============================================

func TestToUpdateIntegrationTypeDTO_SoloCamposPresentes(t *testing.T) {
	// Arrange
	nombre := "Nombre Actualizado"
	req := request.UpdateIntegrationTypeRequest{
		Name: &nombre,
	}

	// Act
	dto := ToUpdateIntegrationTypeDTO(req)

	// Assert
	assert.NotNil(t, dto.Name)
	assert.Equal(t, "Nombre Actualizado", *dto.Name)
	assert.Nil(t, dto.Code)
	assert.Nil(t, dto.IsActive)
}

func TestToUpdateIntegrationTypeDTO_TodosLosCampos(t *testing.T) {
	// Arrange
	nombre := "Shopify Pro"
	codigo := "shopify_pro"
	desc := "Versión mejorada"
	icono := "shopify-pro"
	categoriaID := uint(3)
	activo := false
	removeImage := true
	configSchema := map[string]interface{}{"version": 2}
	configSchemaPtr := &configSchema

	req := request.UpdateIntegrationTypeRequest{
		Name:         &nombre,
		Code:         &codigo,
		Description:  &desc,
		Icon:         &icono,
		CategoryID:   &categoriaID,
		IsActive:     &activo,
		ConfigSchema: configSchemaPtr,
		RemoveImage:  &removeImage,
	}

	// Act
	dto := ToUpdateIntegrationTypeDTO(req)

	// Assert
	assert.Equal(t, "Shopify Pro", *dto.Name)
	assert.Equal(t, "shopify_pro", *dto.Code)
	assert.Equal(t, "Versión mejorada", *dto.Description)
	assert.Equal(t, "shopify-pro", *dto.Icon)
	assert.Equal(t, uint(3), *dto.CategoryID)
	assert.Equal(t, false, *dto.IsActive)
	assert.NotNil(t, dto.ConfigSchema)
	assert.True(t, dto.RemoveImage)
}

func TestToUpdateIntegrationTypeDTO_SinCampos_DTOVacio(t *testing.T) {
	// Arrange
	req := request.UpdateIntegrationTypeRequest{}

	// Act
	dto := ToUpdateIntegrationTypeDTO(req)

	// Assert — todo queda en nil / zero value
	assert.Nil(t, dto.Name)
	assert.Nil(t, dto.Code)
	assert.Nil(t, dto.IsActive)
	assert.Nil(t, dto.ConfigSchema)
	assert.False(t, dto.RemoveImage)
}
