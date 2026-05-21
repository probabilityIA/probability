package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

type ClientGroupResponse struct {
	ID          uint      `json:"id"`
	BusinessID  uint      `json:"business_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	IsActive    bool      `json:"is_active"`
	MemberCount int64     `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FromGroupEntity(g *entities.ClientGroup) ClientGroupResponse {
	return ClientGroupResponse{
		ID:          g.ID,
		BusinessID:  g.BusinessID,
		Name:        g.Name,
		Description: g.Description,
		Color:       g.Color,
		IsActive:    g.IsActive,
		MemberCount: g.MemberCount,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

type ClientGroupsListResponse struct {
	Data       []ClientGroupResponse `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

type ClientSummaryResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Dni       string `json:"dni"`
	GroupID   *uint  `json:"group_id"`
	GroupName string `json:"group_name"`
}

func FromClientSummary(c *entities.ClientSummary) ClientSummaryResponse {
	return ClientSummaryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Phone:     c.Phone,
		Dni:       c.Dni,
		GroupID:   c.GroupID,
		GroupName: c.GroupName,
	}
}

type ClientsListResponse struct {
	Data       []ClientSummaryResponse `json:"data"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}

type CatalogPriceRowResponse struct {
	ProductID    string   `json:"product_id"`
	ProductName  string   `json:"product_name"`
	ProductSKU   string   `json:"product_sku"`
	ImageURL     string   `json:"image_url"`
	FamilyImageURL string  `json:"family_image_url"`
	Currency     string   `json:"currency"`
	BasePrice    float64  `json:"base_price"`
	CustomPrice  *float64 `json:"custom_price"`
	Difference   float64  `json:"difference"`
}

func FromCatalogPriceRow(row *entities.CatalogPriceRow) CatalogPriceRowResponse {
	resp := CatalogPriceRowResponse{
		ProductID:     row.ProductID,
		ProductName:   row.ProductName,
		ProductSKU:    row.ProductSKU,
		ImageURL:      row.ImageURL,
		FamilyImageURL: row.FamilyImageURL,
		Currency:      row.Currency,
		BasePrice:     row.BasePrice,
		CustomPrice:   row.CustomPrice,
	}
	if row.CustomPrice != nil {
		resp.Difference = *row.CustomPrice - row.BasePrice
	}
	return resp
}

type CatalogPricesListResponse struct {
	Data       []CatalogPriceRowResponse `json:"data"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
	TotalPages int                       `json:"total_pages"`
}

type EffectivePriceResponse struct {
	ProductID  string  `json:"product_id"`
	BasePrice  float64 `json:"base_price"`
	FinalPrice float64 `json:"final_price"`
	Source     string  `json:"source"`
	GroupID    *uint   `json:"group_id"`
	GroupName  string  `json:"group_name"`
}

func FromEffectivePrice(p *entities.EffectivePrice) EffectivePriceResponse {
	return EffectivePriceResponse{
		ProductID:  p.ProductID,
		BasePrice:  p.BasePrice,
		FinalPrice: p.FinalPrice,
		Source:     p.Source,
		GroupID:    p.GroupID,
		GroupName:  p.GroupName,
	}
}

func TotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		pages++
	}
	return pages
}
