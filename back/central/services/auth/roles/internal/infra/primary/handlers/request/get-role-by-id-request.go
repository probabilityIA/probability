package request

// GetRoleByIDRequest representa la solicitud para obtener un rol por ID
type GetRoleByIDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
