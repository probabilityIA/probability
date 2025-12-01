package request

// GetUsersRequest representa la solicitud para obtener usuarios filtrados
type GetUsersRequest struct {
	Page       int    `form:"page" binding:"min=1" default:"1"`
	PageSize   int    `form:"page_size" binding:"min=1,max=100" default:"10"`
	Name       string `form:"name" binding:"omitempty,max=100"`
	Email      string `form:"email" binding:"omitempty,email"`
	Phone      string `form:"phone" binding:"omitempty,len=10"`
	UserIDs    string `form:"user_ids" binding:"omitempty"` // IDs separados por comas
	IsActive   *bool  `form:"is_active"`
	RoleID     *uint  `form:"role_id" binding:"omitempty,min=1"`
	BusinessID *uint  `form:"business_id" binding:"omitempty,min=1"`
	CreatedAt  string `form:"created_at" binding:"omitempty,datetime=2006-01-02"` // formato: "2024-01-01" o "2024-01-01,2024-12-31"
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=id name email phone is_active created_at updated_at" default:"created_at"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc" default:"desc"`
}
