package dtos

// ListMovementsParams parámetros para listar movimientos de stock (paginado)
type ListMovementsParams struct {
	BusinessID  uint
	ProductID   string
	WarehouseID *uint
	Type        string // filtrar por tipo: inbound, outbound, adjustment, transfer, return, sync
	Page        int
	PageSize    int
}

// Offset calcula el offset para paginación
func (p ListMovementsParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}
