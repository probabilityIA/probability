package request

// GetRolesByLevelRequest representa la solicitud para obtener roles por nivel
type GetRolesByLevelRequest struct {
	Level int `form:"level" binding:"required,min=1,max=10"`
}
