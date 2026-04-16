package app

import (
	domainerrors "github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/announcements/internal/domain/entities"
)

func validateDisplayType(dt entities.DisplayType) error {
	switch dt {
	case entities.DisplayTypeModalImage, entities.DisplayTypeModalText, entities.DisplayTypeTicker:
		return nil
	}
	return domainerrors.ErrInvalidDisplayType
}

func validateFrequencyType(ft entities.FrequencyType) error {
	switch ft {
	case entities.FrequencyOnce, entities.FrequencyDaily, entities.FrequencyAlways, entities.FrequencyRequiresAcceptance:
		return nil
	}
	return domainerrors.ErrInvalidFrequencyType
}
