package app

import (
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/warehouses/internal/domain/entities"
)

func toEntityStandardBoxes(boxes []dtos.StandardBoxDTO) []entities.StandardBox {
	result := make([]entities.StandardBox, len(boxes))
	for i, b := range boxes {
		result[i] = entities.StandardBox{
			Name:     b.Name,
			Weight:   b.Weight,
			Length:   b.Length,
			Width:    b.Width,
			Height:   b.Height,
			MaxItems: b.MaxItems,
		}
	}
	return result
}
