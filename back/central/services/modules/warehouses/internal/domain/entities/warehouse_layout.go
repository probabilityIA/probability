package entities

import "time"

type LayoutNode struct {
	NodeID   string
	RefType  string
	RefID    uint
	X        float64
	Y        float64
	Width    float64
	Height   float64
	Rotation float64
	Color    string
	Label    string
}

type WarehouseLayout struct {
	ID           uint
	WarehouseID  uint
	BusinessID   uint
	CanvasWidth  float64
	CanvasHeight float64
	GridSize     float64
	Scale        float64
	Nodes        []LayoutNode
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
