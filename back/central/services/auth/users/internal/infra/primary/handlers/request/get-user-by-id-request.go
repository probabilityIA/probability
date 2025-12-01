package request

// GetUserByIDRequest representa la solicitud para obtener un usuario por ID
type GetUserByIDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
