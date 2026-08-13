package productmatch

const (
	GroupSKUSpacing = "sku_spacing"
	GroupSKUTypo    = "sku_typo"
)

type TypoDetail struct {
	SKU             string
	Label           string
	Tone            string
	Group           string
	CounterpartSKU  string
	CounterpartName string
	ChannelQty      *int
	OwnQty          *int
	FixSide         string
	Pattern         string
	ParentRef       string
	ParentLabel     string
	VariantLabel    string
}

func qtyOrNil(value int, has bool) *int {
	if !has {
		return nil
	}
	v := value
	return &v
}

func FixInstruction(item TypoSuspect, channelLabel string) string {
	donde, correcto := channelLabel, item.ProbabilitySKU
	if item.FixSide == FixInProbability {
		donde, correcto = "Probability", item.ChannelSKU
	}
	switch item.Pattern {
	case PatternSpacing:
		return "sobra un espacio en el SKU de " + donde + ": deberia quedar " + correcto
	case PatternTransposed:
		return "letras trocadas en " + donde + ": deberia decir " + correcto
	default:
		return "un caracter de diferencia y existencias parecidas; en " + donde + " deberia decir " + correcto
	}
}

func TypoDetails(suspects []TypoSuspect, channelLabel string) []TypoDetail {
	out := make([]TypoDetail, 0, len(suspects))
	for _, item := range suspects {
		grupo, tono := GroupSKUTypo, "warn"
		if item.Pattern == PatternSpacing {
			grupo, tono = GroupSKUSpacing, "info"
		}
		out = append(out, TypoDetail{
			SKU:             item.ChannelSKU,
			Label:           FixInstruction(item, channelLabel),
			Tone:            tono,
			Group:           grupo,
			CounterpartSKU:  item.ProbabilitySKU,
			CounterpartName: item.ProbabilityName,
			ChannelQty:      qtyOrNil(item.ChannelQty, item.HasChannelQty),
			OwnQty:          qtyOrNil(item.ProbabilityQty, item.HasProbQty),
			FixSide:         item.FixSide,
			Pattern:         item.Pattern,
			ParentRef:       item.ChannelRef,
			ParentLabel:     item.ChannelName,
			VariantLabel:    item.ChannelLabel,
		})
	}
	return out
}

func CountPattern(suspects []TypoSuspect, pattern string) int {
	total := 0
	for _, item := range suspects {
		if item.Pattern == pattern {
			total++
		}
	}
	return total
}
