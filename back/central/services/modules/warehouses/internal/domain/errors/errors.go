package errors

import "errors"

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrDuplicateCode     = errors.New("a warehouse with this code already exists in your business")
	ErrLocationNotFound  = errors.New("warehouse location not found")
	ErrDuplicateLocCode  = errors.New("a location with this code already exists in this warehouse")
	ErrWarehouseHasStock = errors.New("warehouse has inventory and cannot be deleted")
	ErrLocationHasStock  = errors.New("location has inventory and cannot be deleted")

	ErrZoneNotFound         = errors.New("warehouse zone not found")
	ErrDuplicateZoneCode    = errors.New("a zone with this code already exists in this warehouse")
	ErrAisleNotFound        = errors.New("warehouse aisle not found")
	ErrDuplicateAisleCode   = errors.New("an aisle with this code already exists in this zone")
	ErrRackNotFound         = errors.New("warehouse rack not found")
	ErrDuplicateRackCode    = errors.New("a rack with this code already exists in this aisle")
	ErrRackLevelNotFound    = errors.New("rack level not found")
	ErrDuplicateLevelCode   = errors.New("a level with this code already exists in this rack")
	ErrInvalidHierarchy     = errors.New("invalid hierarchy reference")
	ErrLocationOverCapacity = errors.New("location volume or weight exceeds configured capacity")
)
