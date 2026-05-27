package usecaseshipment

import _ "embed"

//go:embed assets/probability-logo.png
var probabilityLogoPNG []byte

//go:embed assets/probability-icon.png
var probabilityIconPNG []byte

type carrierStyle struct {
	BgR, BgG, BgB    int
	TxtR, TxtG, TxtB int
}

var carrierStyles = map[string]carrierStyle{
	"ENVIA":           {220, 35, 35, 255, 255, 255},
	"INTERRAPIDISIMO": {255, 200, 0, 0, 0, 0},
	"COORDINADORA":    {30, 90, 200, 255, 255, 255},
	"ENVIOCLICK":      {30, 90, 200, 255, 255, 255},
	"TCC":             {235, 90, 30, 255, 255, 255},
	"SERVIENTREGA":    {120, 60, 160, 255, 255, 255},
	"99MINUTOS":       {255, 110, 30, 255, 255, 255},
	"DEPRISA":         {220, 30, 60, 255, 255, 255},
}

func styleForCarrier(carrier string) carrierStyle {
	if s, ok := carrierStyles[carrier]; ok {
		return s
	}
	return carrierStyle{60, 60, 60, 255, 255, 255}
}
