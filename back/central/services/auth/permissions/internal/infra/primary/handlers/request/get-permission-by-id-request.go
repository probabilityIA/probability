package request

// GetPermissionByIDRequest representa la solicitud para obtener un permiso por ID
type GetPermissionByIDRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}
