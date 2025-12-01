package request

// DeletePermissionRequest representa la solicitud para eliminar un permiso
type DeletePermissionRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
