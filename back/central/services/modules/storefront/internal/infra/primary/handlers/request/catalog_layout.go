package request

type UpdateCatalogLayoutRequest struct {
	Columns int `json:"columns" binding:"required,min=1,max=6"`
	Rows    int `json:"rows" binding:"required,min=1,max=12"`
}
