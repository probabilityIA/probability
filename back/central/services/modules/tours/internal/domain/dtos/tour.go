package dtos

import "github.com/secamc93/probability/back/central/services/modules/tours/internal/domain/entities"

type ListProgressInput struct {
	UserID     uint
	BusinessID uint
}

type SaveProgressInput struct {
	UserID     uint
	BusinessID uint
	TourKey    string
	Version    int
	Status     string
	StepIndex  int
}

type ResetInput struct {
	UserID     uint
	BusinessID uint
	TourKey    string
}

type ListProgressResult struct {
	Items []entities.TourProgress
}

type SkipAllTour struct {
	TourKey string
	Version int
}

type SkipAllInput struct {
	UserID     uint
	BusinessID uint
	Tours      []SkipAllTour
}
