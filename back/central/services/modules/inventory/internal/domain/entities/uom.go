package entities

type UnitOfMeasure struct {
	ID       uint
	Code     string
	Name     string
	Type     string
	IsActive bool
}

type ProductUoM struct {
	ID               uint
	ProductID        string
	UomID            uint
	BusinessID       uint
	ConversionFactor float64
	IsBase           bool
	Barcode          string
	IsActive         bool

	UomCode string
	UomName string
	UomType string
}
