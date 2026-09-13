package dtos

type CreateFlowDTO struct {
	BusinessID     uint   `json:"-"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	RootTemplateID *uint  `json:"root_template_id"`
}

type UpdateFlowDTO struct {
	ID             uint   `json:"-"`
	BusinessID     uint   `json:"-"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	RootTemplateID *uint  `json:"root_template_id"`
	Enabled        *bool  `json:"enabled"`
}
