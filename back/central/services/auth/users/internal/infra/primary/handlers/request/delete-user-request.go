package request

// DeleteUserRequest representa la solicitud para eliminar un usuario
type DeleteUserRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
