package dtos

// GetProductInventoryParams parámetros para obtener inventario de un producto en todas las bodegas
type GetProductInventoryParams struct {
	ProductID  string
	BusinessID uint
}

// ListWarehouseInventoryParams parámetros para listar inventario de una bodega (paginado)
type ListWarehouseInventoryParams struct {
	WarehouseID uint
	BusinessID  uint
	Search      string // buscar por nombre o SKU del producto
	LowStock    bool   // solo productos con stock bajo
	Page        int
	PageSize    int
}

// Offset calcula el offset para paginación
func (p ListWarehouseInventoryParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.PageSize
}
