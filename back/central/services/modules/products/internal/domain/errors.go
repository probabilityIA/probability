package domain

import "errors"

var (
	ErrProductNotFound            = errors.New("product not found")
	ErrProductAlreadyExists       = errors.New("product with this SKU already exists for this business")
	ErrInvalidProductData         = errors.New("invalid product data")
	ErrProductFamilyNotFound      = errors.New("product family not found")
	ErrVariantAlreadyExists       = errors.New("product variant with the same attributes already exists for this family")
	ErrProductIntegrationNotFound = errors.New("product integration mapping not found")
	ErrFamilyHasActiveVariants    = errors.New("cannot delete family with active variants")
)
