package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

type LayoutNodeResponse struct {
	NodeID   string  `json:"node_id"`
	RefType  string  `json:"ref_type"`
	RefID    uint    `json:"ref_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	Color    string  `json:"color"`
	Label    string  `json:"label"`
}

type LayoutResponse struct {
	WarehouseID  uint                 `json:"warehouse_id"`
	CanvasWidth  float64              `json:"canvas_width"`
	CanvasHeight float64              `json:"canvas_height"`
	GridSize     float64              `json:"grid_size"`
	Nodes        []LayoutNodeResponse `json:"nodes"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

func LayoutFromEntity(l *entities.WarehouseLayout) LayoutResponse {
	nodes := make([]LayoutNodeResponse, len(l.Nodes))
	for i, n := range l.Nodes {
		nodes[i] = LayoutNodeResponse{
			NodeID:   n.NodeID,
			RefType:  n.RefType,
			RefID:    n.RefID,
			X:        n.X,
			Y:        n.Y,
			Width:    n.Width,
			Height:   n.Height,
			Rotation: n.Rotation,
			Color:    n.Color,
			Label:    n.Label,
		}
	}
	return LayoutResponse{
		WarehouseID:  l.WarehouseID,
		CanvasWidth:  l.CanvasWidth,
		CanvasHeight: l.CanvasHeight,
		GridSize:     l.GridSize,
		Nodes:        nodes,
		UpdatedAt:    l.UpdatedAt,
	}
}
