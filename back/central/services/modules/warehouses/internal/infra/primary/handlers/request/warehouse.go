package request

// CreateWarehouseRequest payload de creación de bodega
type CreateWarehouseRequest struct {
	Name          string `json:"name" binding:"required,min=2,max=255"`
	Code          string `json:"code" binding:"required,min=1,max=50"`
	Address       string `json:"address" binding:"omitempty,max=500"`
	City          string `json:"city" binding:"omitempty,max=100"`
	State         string `json:"state" binding:"omitempty,max=100"`
	Country       string `json:"country" binding:"omitempty,max=50"`
	ZipCode       string `json:"zip_code" binding:"omitempty,max=20"`
	Phone         string `json:"phone" binding:"omitempty,max=50"`
	ContactName   string `json:"contact_name" binding:"omitempty,max=255"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,email,max=255"`
	IsActive      *bool  `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
	IsFulfillment bool   `json:"is_fulfillment"`
}

// UpdateWarehouseRequest payload de actualización de bodega
type UpdateWarehouseRequest struct {
	Name          string `json:"name" binding:"required,min=2,max=255"`
	Code          string `json:"code" binding:"required,min=1,max=50"`
	Address       string `json:"address" binding:"omitempty,max=500"`
	City          string `json:"city" binding:"omitempty,max=100"`
	State         string `json:"state" binding:"omitempty,max=100"`
	Country       string `json:"country" binding:"omitempty,max=50"`
	ZipCode       string `json:"zip_code" binding:"omitempty,max=20"`
	Phone         string `json:"phone" binding:"omitempty,max=50"`
	ContactName   string `json:"contact_name" binding:"omitempty,max=255"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,email,max=255"`
	IsActive      *bool  `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
	IsFulfillment bool   `json:"is_fulfillment"`
}

// CreateLocationRequest payload de creación de ubicación
type CreateLocationRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=255"`
	Code          string `json:"code" binding:"required,min=1,max=50"`
	Type          string `json:"type" binding:"omitempty,oneof=storage picking packing receiving shipping"`
	IsActive      *bool  `json:"is_active"`
	IsFulfillment bool   `json:"is_fulfillment"`
	Capacity      *int   `json:"capacity" binding:"omitempty,min=0"`
}

// UpdateLocationRequest payload de actualización de ubicación
type UpdateLocationRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=255"`
	Code          string `json:"code" binding:"required,min=1,max=50"`
	Type          string `json:"type" binding:"omitempty,oneof=storage picking packing receiving shipping"`
	IsActive      *bool  `json:"is_active"`
	IsFulfillment bool   `json:"is_fulfillment"`
	Capacity      *int   `json:"capacity" binding:"omitempty,min=0"`
}
