package models

import "gorm.io/gorm"

type InventoryLevel struct {
	gorm.Model
	ProductID    string `gorm:"type:varchar(64);not null;uniqueIndex:idx_inventory_level_key,priority:1"`
	WarehouseID  uint   `gorm:"not null;uniqueIndex:idx_inventory_level_key,priority:2"`
	LocationID   *uint  `gorm:"index;uniqueIndex:idx_inventory_level_key,priority:3"`
	LotID        *uint  `gorm:"index;uniqueIndex:idx_inventory_level_key,priority:4"`
	StateID      *uint  `gorm:"index;uniqueIndex:idx_inventory_level_key,priority:5"`
	BusinessID   uint   `gorm:"not null;index"`
	Quantity     int    `gorm:"default:0;not null"`
	ReservedQty  int    `gorm:"default:0;not null"`
	AvailableQty int    `gorm:"default:0;not null"`
	MinStock     *int
	MaxStock     *int
	ReorderPoint *int

	Product   Product            `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse Warehouse          `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Location  *WarehouseLocation `gorm:"foreignKey:LocationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Business  Business           `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Lot       *InventoryLot      `gorm:"foreignKey:LotID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	State     *InventoryState    `gorm:"foreignKey:StateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla
func (InventoryLevel) TableName() string {
	return "inventory_levels"
}

type StockMovement struct {
	gorm.Model
	ProductID      string  `gorm:"type:varchar(64);not null;index"`
	WarehouseID    uint    `gorm:"not null;index"`
	LocationID     *uint   `gorm:"index"`
	LotID          *uint   `gorm:"index"`
	SerialID       *uint   `gorm:"index"`
	FromStateID    *uint   `gorm:"index"`
	ToStateID      *uint   `gorm:"index"`
	UomID          *uint   `gorm:"index"`
	QtyInBaseUom   float64 `gorm:"default:0"`
	BusinessID     uint    `gorm:"not null;index"`
	MovementTypeID uint    `gorm:"not null;index"`
	Reason         string  `gorm:"size:255"`
	Quantity       int     `gorm:"not null"`
	PreviousQty    int     `gorm:"not null"`
	NewQty         int     `gorm:"not null"`
	ReferenceType  *string `gorm:"size:50"`
	ReferenceID    *string `gorm:"size:64"`
	IntegrationID  *uint   `gorm:"index"`
	Notes          string  `gorm:"type:text"`
	CreatedByID    *uint   `gorm:"index"`

	MovementType StockMovementType `gorm:"foreignKey:MovementTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Product      Product           `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Warehouse    Warehouse         `gorm:"foreignKey:WarehouseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Business     Business          `gorm:"foreignKey:BusinessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedBy    *User             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Lot          *InventoryLot     `gorm:"foreignKey:LotID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Serial       *InventorySerial  `gorm:"foreignKey:SerialID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Uom          *UnitOfMeasure    `gorm:"foreignKey:UomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName especifica el nombre de la tabla
func (StockMovement) TableName() string {
	return "stock_movements"
}
