package request

type SaveClientGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type AddGroupMembersRequest struct {
	ClientIDs []uint `json:"client_ids" binding:"required,min=1"`
}

type CatalogPriceItemRequest struct {
	ProductID string   `json:"product_id" binding:"required"`
	Price     *float64 `json:"price"`
}

type SaveCatalogPricesRequest struct {
	ClientGroupID *uint                     `json:"client_group_id"`
	ClientID      *uint                     `json:"client_id"`
	Items         []CatalogPriceItemRequest `json:"items" binding:"required"`
}
