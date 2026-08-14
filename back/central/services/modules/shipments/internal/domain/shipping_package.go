package domain

import "sort"

const (
	ShippingPackageStrategyProductDimensions = "product_dimensions"
	ShippingPackageStrategyStandardBox       = "standard_box"
)

type StandardBox struct {
	Name     string
	Weight   *float64
	Length   *float64
	Width    *float64
	Height   *float64
	MaxItems int
}

type ShippingPackageConfig struct {
	Strategy      string
	StandardBoxes []StandardBox
}

type ProductDimensions struct {
	Weight *float64
	Length *float64
	Width  *float64
	Height *float64
}

func (c *ShippingPackageConfig) SelectBox(totalQuantity int, itemLength, itemWidth, itemHeight float64) *StandardBox {
	if c == nil || len(c.StandardBoxes) == 0 {
		return nil
	}

	var best *StandardBox
	for i := range c.StandardBoxes {
		box := &c.StandardBoxes[i]
		if box.MaxItems < totalQuantity {
			continue
		}
		if !boxFitsItem(box, itemLength, itemWidth, itemHeight) {
			continue
		}
		if best == nil || box.MaxItems < best.MaxItems {
			best = box
		}
	}
	return best
}

func boxFitsItem(box *StandardBox, itemLength, itemWidth, itemHeight float64) bool {
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
