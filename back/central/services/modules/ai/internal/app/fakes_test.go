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
	replies  []*dtos.ModelReply
	err      error
	requests []dtos.ModelRequest
}

var _ ports.IAssistantModel = (*modelFake)(nil)

func (f *modelFake) Reply(_ context.Context, req dtos.ModelRequest) (*dtos.ModelReply, error) {
	snapshot := req
	snapshot.Messages = append([]dtos.ModelMessage(nil), req.Messages...)
	f.requests = append(f.requests, snapshot)
	if f.err != nil {
		return nil, f.err
	}
	if len(f.replies) > 0 {
		index := len(f.requests) - 1
		if index >= len(f.replies) {
			index = len(f.replies) - 1
		}
		return f.replies[index], nil
	}
	return f.reply, nil
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

type readerFake struct {
	unread      []entities.UnreadChats
	orders      []entities.OrderInfo
	businessIDs []uint
	numbers     []string
}

var _ ports.IBusinessDataReader = (*readerFake)(nil)

func (f *readerFake) FindOrders(_ context.Context, businessID uint, number string) ([]entities.OrderInfo, error) {
	f.businessIDs = append(f.businessIDs, businessID)
	f.numbers = append(f.numbers, number)
	return f.orders, nil
}

func (f *readerFake) FindShipments(_ context.Context, businessID uint, _ string) ([]entities.ShipmentInfo, error) {
	f.businessIDs = append(f.businessIDs, businessID)
	return nil, nil
}

func (f *readerFake) ListOrders(_ context.Context, businessID uint, _ dtos.OrderQuery) ([]entities.OrderSummary, int64, error) {
	f.businessIDs = append(f.businessIDs, businessID)
	return nil, 0, nil
}

func (f *readerFake) DescribeIdentity(_ context.Context, _ uint, _ *uint) (*entities.ChatIdentity, error) {
	return &entities.ChatIdentity{UserName: "Ana", BusinessName: "Demo"}, nil
}

func (f *readerFake) CountUnreadWhatsAppChats(_ context.Context) ([]entities.UnreadChats, error) {
	return f.unread, nil
}

func (f *readerFake) SummarizeOrders(_ context.Context, businessID uint, from, to time.Time) (*entities.OrdersOverview, error) {
	f.businessIDs = append(f.businessIDs, businessID)
	return &entities.OrdersOverview{From: from, To: to}, nil
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
	return New(&recommendationFake{}, model, &navigationFake{catalog: sampleCatalog()}, s, nil, nil, nil, nil, log.New()).(*UseCase)
}

func newDataUseCase(model *modelFake, reader *readerFake, catalog *entities.NavigationCatalog) *UseCase {
	uc := New(&recommendationFake{}, model, &navigationFake{catalog: catalog}, &storeFake{}, nil, nil, reader, nil, log.New()).(*UseCase)
	uc.now = func() time.Time { return time.Date(2026, 9, 14, 15, 0, 0, 0, colombia) }
	return uc
}

type alertsFake struct {
	saved []entities.Alert
	last  map[uint]*entities.Alert
}

var _ ports.IAlertRepository = (*alertsFake)(nil)

func (f *alertsFake) SaveAlert(_ context.Context, alert entities.Alert) (bool, error) {
	f.saved = append(f.saved, alert)
	return true, nil
}

func (f *alertsFake) ListAlerts(_ context.Context, _ dtos.AlertQuery) ([]entities.Alert, int64, error) {
	return nil, 0, nil
}

func (f *alertsFake) CountUnread(_ context.Context, _, _ uint) (int64, *entities.Alert, error) {
	return 0, nil, nil
}

func (f *alertsFake) MarkSeen(_ context.Context, _, _ uint, _ time.Time) error { return nil }

func (f *alertsFake) RecentAlerts(_ context.Context, _ uint, _ int) ([]entities.Alert, error) {
	return nil, nil
}

func (f *alertsFake) DeleteAlertsOlderThan(_ context.Context, _ time.Time) (int64, error) {
	return 0, nil
}

func (f *alertsFake) LastAlertOfType(_ context.Context, businessID uint, _ string) (*entities.Alert, error) {
	if f.last == nil {
		return nil, nil
	}
	return f.last[businessID], nil
}
