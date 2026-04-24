package entities

type InventoryState struct {
	ID          uint
	Code        string
	Name        string
	Description string
	IsTerminal  bool
	IsActive    bool
}
