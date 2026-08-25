package shippingpkg

import "sort"

const DefaultDimCm = 10.0
const DefaultWeightKg = 1.0

type PackageInput struct {
	TotalQuantity  int
	CartWeightKg   float64
	CatalogWeights map[string]float64
	Items          []PackageItem
}

type PackageItem struct {
	SKU      string
	Quantity int
	Length   *float64
	Width    *float64
	Height   *float64
	Weight   *float64
}

type ResolvedPackage struct {
	Weight  float64 `json:"weight"`
	Length  float64 `json:"length"`
	Width   float64 `json:"width"`
	Height  float64 `json:"height"`
	Source  string  `json:"source"`
	BoxName string  `json:"box_name,omitempty"`
}

const (
	PackageSourceStandardBox = "standard_box"
	PackageSourceProduct     = "product_dimensions"
	PackageSourceDefault     = "default"
)

func Resolve(strategy string, boxes []Box, in PackageInput) ResolvedPackage {
	weight := resolveWeight(in)
	maxLength, maxWidth, maxHeight := maxItemDims(in.Items)

	quantity := in.TotalQuantity
	if quantity <= 0 {
		quantity = totalItemQuantity(in.Items)
	}

	if box := SelectBox(strategy, boxes, quantity, maxLength, maxWidth, maxHeight); box != nil {
		out := ResolvedPackage{
			Weight:  weight,
			Length:  dimOrFallback(box.Length, maxLength),
			Width:   dimOrFallback(box.Width, maxWidth),
			Height:  dimOrFallback(box.Height, maxHeight),
			Source:  PackageSourceStandardBox,
			BoxName: box.Name,
		}
		if box.Weight != nil && *box.Weight > out.Weight {
			out.Weight = *box.Weight
		}
		return out
	}

	if maxLength > 0 || maxWidth > 0 || maxHeight > 0 {
		return ResolvedPackage{
			Weight: weight,
			Length: positiveOrDefault(maxLength),
			Width:  positiveOrDefault(maxWidth),
			Height: positiveOrDefault(maxHeight),
			Source: PackageSourceProduct,
		}
	}

	return ResolvedPackage{
		Weight: weight,
		Length: DefaultDimCm,
		Width:  DefaultDimCm,
		Height: DefaultDimCm,
		Source: PackageSourceDefault,
	}
}

func resolveWeight(in PackageInput) float64 {
	if in.CartWeightKg > 0 {
		return in.CartWeightKg
	}

	var total float64
	for _, it := range in.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		if it.Weight != nil && *it.Weight > 0 {
			total += *it.Weight * float64(qty)
			continue
		}
		if w, ok := in.CatalogWeights[it.SKU]; ok && w > 0 {
			total += w * float64(qty)
		}
	}
	if total > 0 {
		return total
	}
	return DefaultWeightKg
}

func maxItemDims(items []PackageItem) (float64, float64, float64) {
	var maxLength, maxWidth, maxHeight float64
	for _, it := range items {
		if it.Length != nil && *it.Length > maxLength {
			maxLength = *it.Length
		}
		if it.Width != nil && *it.Width > maxWidth {
			maxWidth = *it.Width
		}
		if it.Height != nil && *it.Height > maxHeight {
			maxHeight = *it.Height
		}
	}
	return maxLength, maxWidth, maxHeight
}

func totalItemQuantity(items []PackageItem) int {
	total := 0
	for _, it := range items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		total += qty
	}
	return total
}

func dimOrFallback(v *float64, fallback float64) float64 {
	if v != nil && *v > 0 {
		return *v
	}
	return positiveOrDefault(fallback)
}

func positiveOrDefault(v float64) float64 {
	if v > 0 {
		return v
	}
	return DefaultDimCm
}

const (
	StrategyProductDimensions = "product_dimensions"
	StrategyStandardBox       = "standard_box"
)

type Box struct {
	Name      string   `json:"name"`
	Weight    *float64 `json:"weight"`
	Length    *float64 `json:"length"`
	Width     *float64 `json:"width"`
	Height    *float64 `json:"height"`
	MaxItems  int      `json:"max_items"`
	MaxWeight *float64 `json:"max_weight,omitempty"`
}

func SelectBox(strategy string, boxes []Box, totalQuantity int, itemLength, itemWidth, itemHeight float64) *Box {
	if strategy != StrategyStandardBox || len(boxes) == 0 {
		return nil
	}

	var best *Box
	for i := range boxes {
		box := &boxes[i]
		if box.MaxItems > 0 && box.MaxItems < totalQuantity {
			continue
		}
		if !boxFitsItem(box, itemLength, itemWidth, itemHeight) {
			continue
		}
		if best == nil || boxCapacity(box) < boxCapacity(best) {
			best = box
		}
	}
	return best
}

func boxCapacity(b *Box) int {
	if b.MaxItems <= 0 {
		return 1 << 30
	}
	return b.MaxItems
}

func boxFitsItem(box *Box, itemLength, itemWidth, itemHeight float64) bool {
	boxDims := sortedDims(dimOrInf(box.Length), dimOrInf(box.Width), dimOrInf(box.Height))
	itemDims := sortedDims(itemLength, itemWidth, itemHeight)
	return boxDims[0] >= itemDims[0] && boxDims[1] >= itemDims[1] && boxDims[2] >= itemDims[2]
}

func dimOrInf(v *float64) float64 {
	if v == nil {
		return 1e18
	}
	return *v
}

func sortedDims(a, b, c float64) [3]float64 {
	d := []float64{a, b, c}
	sort.Sort(sort.Reverse(sort.Float64Slice(d)))
	return [3]float64{d[0], d[1], d[2]}
}
