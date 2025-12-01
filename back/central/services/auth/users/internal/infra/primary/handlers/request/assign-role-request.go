package request

// BusinessRoleAssignmentItem representa una asignación individual de rol a business
type BusinessRoleAssignmentItem struct {
	BusinessID uint `json:"business_id" binding:"required,min=1" example:"16"`
	RoleID     uint `json:"role_id" binding:"required,min=1" example:"4"`
}

// AssignRoleToUserBusinessRequest representa la solicitud para asignar roles a un usuario en múltiples businesses
type AssignRoleToUserBusinessRequest struct {
	Assignments []BusinessRoleAssignmentItem `json:"assignments" binding:"required,min=1,dive" example:"[{\"business_id\":16,\"role_id\":4},{\"business_id\":21,\"role_id\":5}]"`
}
