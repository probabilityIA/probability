package errors

import "errors"

var (
	ErrWarehouseNotFound   = errors.New("warehouse not found")
	ErrDuplicateCode       = errors.New("a warehouse with this code already exists in your business")
	ErrLocationNotFound    = errors.New("warehouse location not found")
	ErrDuplicateLocCode    = errors.New("a location with this code already exists in this warehouse")
	ErrWarehouseHasStock   = errors.New("warehouse has inventory and cannot be deleted")
	ErrLocationHasStock    = errors.New("location has inventory and cannot be deleted")
)
