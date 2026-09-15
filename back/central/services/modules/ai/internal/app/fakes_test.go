package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type recommendationFake struct {
	fn    func(origin, destination string) (*entities.Recommendation, error)
	calls [][2]string
}

var _ ports.IRecommendationProvider = (*recommendationFake)(nil)

func (f *recommendationFake) GetRecommendation(_ context.Context, origin, destination string) (*entities.Recommendation, error) {
	f.calls = append(f.calls, [2]string{origin, destination})
	if f.fn != nil {
		return f.fn(origin, destination)
	}
	return &entities.Recommendation{RecommendedCarrier: "SERVIENTREGA"}, nil
}

type modelFake struct {
	reply    *dtos.ModelReply
	err      error
	requests []dtos.ModelRequest
}

var _ ports.IAssistantModel = (*modelFake)(nil)

func (f *modelFake) Reply(_ context.Context, req dtos.ModelRequest) (*dtos.ModelReply, error) {
	f.requests = append(f.requests, req)
	return f.reply, f.err
}

type navigationFake struct {
	catalog *entities.NavigationCatalog
	err     error
}

var _ ports.INavigationCatalog = (*navigationFake)(nil)

func (f *navigationFake) ForUser(_ context.Context, _ dtos.AccessScope) (*entities.NavigationCatalog, error) {
	return f.catalog, f.err
}

type storeFake struct {
	count     int
	consumeEr error
	introSeen bool
}

var _ ports.IAssistantStore = (*storeFake)(nil)

func (f *storeFake) ConsumeMessage(_ context.Context, _ uint, window time.Duration) (*entities.Usage, error) {
	if f.consumeEr != nil {
		return nil, f.consumeEr
	}
	f.count++
	reset := time.Now().Add(window)
	return &entities.Usage{Count: f.count, ResetAt: &reset}, nil
}

func (f *storeFake) GetUsage(_ context.Context, _ uint) (*entities.Usage, error) {
	return &entities.Usage{Count: f.count}, nil
}

func (f *storeFake) IsIntroSeen(_ context.Context, _ uint) (bool, error) {
	return f.introSeen, nil
}

func (f *storeFake) MarkIntroSeen(_ context.Context, _ uint) error {
	f.introSeen = true
	return nil
}

func sampleCatalog() *entities.NavigationCatalog {
	return &entities.NavigationCatalog{
		Allowed: []entities.Destination{
			{Key: "orders", Label: "Ordenes", Route: "/orders", Description: "Pedidos"},
			{Key: "shipments.cod", Label: "Recaudo contra entrega", Route: "/shipments/cod"},
		},
		Denied: []string{"Contabilidad"},
	}
}

func newTestUseCase(model *modelFake, store *storeFake) *UseCase {
	var s ports.IAssistantStore
	if store != nil {
		s = store
	}
	return New(&recommendationFake{}, model, &navigationFake{catalog: sampleCatalog()}, s, log.New()).(*UseCase)
}
