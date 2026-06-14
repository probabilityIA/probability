package request

type LayoutNodeRequest struct {
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

type SaveLayoutRequest struct {
	CanvasWidth  float64             `json:"canvas_width"`
	CanvasHeight float64             `json:"canvas_height"`
	GridSize     float64             `json:"grid_size"`
	Nodes        []LayoutNodeRequest `json:"nodes"`
}
