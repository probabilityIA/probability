package domain

type GuideFormat struct {
	ID          uint    `json:"id"`
	Carrier     string  `json:"carrier"`
	Code        string  `json:"code"`
	Label       string  `json:"label"`
	WidthCm     float64 `json:"width_cm"`
	HeightCm    float64 `json:"height_cm"`
	Adhesive    bool    `json:"adhesive"`
	Strategy    string  `json:"strategy"`
	CropLLxFrac float64 `json:"crop_llx_frac"`
	CropLLyFrac float64 `json:"crop_lly_frac"`
	CropURxFrac float64 `json:"crop_urx_frac"`
	CropURyFrac float64 `json:"crop_ury_frac"`
	SourcePage  int     `json:"source_page"`
	IsDefault   bool    `json:"is_default"`
	SortOrder   int     `json:"sort_order"`
}

const (
	GuideStrategyPassthrough = "passthrough"
	GuideStrategyCrop        = "crop"
	GuideStrategyResize      = "resize"
	GuideStrategyRebuild     = "rebuild"

	CarrierUniversal = "*"
)
