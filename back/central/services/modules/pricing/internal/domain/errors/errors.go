package errors

import "errors"

var (
	ErrRuleNotFound          = errors.New("pricing rule not found")
	ErrDiscountNotFound      = errors.New("quantity discount not found")
	ErrDuplicateRule         = errors.New("a pricing rule already exists for this client/product combination")
	ErrDuplicateDiscount     = errors.New("a quantity discount already exists for this product/quantity tier")
	ErrInvalidAdjustmentType = errors.New("adjustment_type must be 'percentage' or 'fixed'")
	ErrInvalidMinQuantity    = errors.New("min_quantity must be greater than 0")
	ErrInvalidDiscountPercent = errors.New("discount_percent must be between 0 and 100")
)
