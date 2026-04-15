package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

type mockRepo struct {
	getCustomerSummaryFn        func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error)
	listCustomerAddressesFn     func(ctx context.Context, params dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error)
	listCustomerProductsFn      func(ctx context.Context, params dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error)
	listCustomerOrderItemsFn    func(ctx context.Context, params dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error)
	upsertCustomerSummaryFn     func(ctx context.Context, s *entities.CustomerSummary) error
	upsertCustomerAddressFn     func(ctx context.Context, a *entities.CustomerAddress) error
	upsertCustomerProductFn     func(ctx context.Context, p *entities.CustomerProductHistory) error
	upsertCustomerOrderItemFn   func(ctx context.Context, i *entities.CustomerOrderItem) error
	updateOrderItemsStatusFn    func(ctx context.Context, orderID string, status string) error
	findClientByPhoneFn         func(ctx context.Context, businessID uint, phone string) (*entities.Client, error)
	getByIDFn                   func(ctx context.Context, businessID, clientID uint) (*entities.Client, error)
}

func (m *mockRepo) Create(_ context.Context, c *entities.Client) (*entities.Client, error)         { return c, nil }
func (m *mockRepo) GetByID(ctx context.Context, bID, cID uint) (*entities.Client, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, bID, cID)
	}
	return &entities.Client{ID: cID, BusinessID: bID, Name: "Test"}, nil
}
func (m *mockRepo) List(_ context.Context, _ dtos.ListClientsParams) ([]entities.Client, int64, error) {
	return nil, 0, nil
}
func (m *mockRepo) Update(_ context.Context, c *entities.Client) (*entities.Client, error) { return c, nil }
func (m *mockRepo) Delete(_ context.Context, _, _ uint) error                              { return nil }
func (m *mockRepo) ExistsByEmail(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
	return false, nil
}
func (m *mockRepo) ExistsByDni(_ context.Context, _ uint, _ string, _ *uint) (bool, error) {
	return false, nil
}

func (m *mockRepo) GetCustomerSummary(ctx context.Context, bID, cID uint) (*entities.CustomerSummary, error) {
	if m.getCustomerSummaryFn != nil {
		return m.getCustomerSummaryFn(ctx, bID, cID)
	}
	return nil, nil
}

func (m *mockRepo) ListCustomerAddresses(ctx context.Context, p dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error) {
	if m.listCustomerAddressesFn != nil {
		return m.listCustomerAddressesFn(ctx, p)
	}
	return nil, 0, nil
}

func (m *mockRepo) ListCustomerProducts(ctx context.Context, p dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error) {
	if m.listCustomerProductsFn != nil {
		return m.listCustomerProductsFn(ctx, p)
	}
	return nil, 0, nil
}

func (m *mockRepo) ListCustomerOrderItems(ctx context.Context, p dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error) {
	if m.listCustomerOrderItemsFn != nil {
		return m.listCustomerOrderItemsFn(ctx, p)
	}
	return nil, 0, nil
}

func (m *mockRepo) UpsertCustomerSummary(ctx context.Context, s *entities.CustomerSummary) error {
	if m.upsertCustomerSummaryFn != nil {
		return m.upsertCustomerSummaryFn(ctx, s)
	}
	return nil
}

func (m *mockRepo) UpsertCustomerAddress(ctx context.Context, a *entities.CustomerAddress) error {
	if m.upsertCustomerAddressFn != nil {
		return m.upsertCustomerAddressFn(ctx, a)
	}
	return nil
}

func (m *mockRepo) UpsertCustomerProductHistory(ctx context.Context, p *entities.CustomerProductHistory) error {
	if m.upsertCustomerProductFn != nil {
		return m.upsertCustomerProductFn(ctx, p)
	}
	return nil
}

func (m *mockRepo) UpsertCustomerOrderItem(ctx context.Context, i *entities.CustomerOrderItem) error {
	if m.upsertCustomerOrderItemFn != nil {
		return m.upsertCustomerOrderItemFn(ctx, i)
	}
	return nil
}

func (m *mockRepo) UpdateOrderItemsStatus(ctx context.Context, orderID string, status string) error {
	if m.updateOrderItemsStatusFn != nil {
		return m.updateOrderItemsStatusFn(ctx, orderID, status)
	}
	return nil
}

func (m *mockRepo) FindClientByPhone(ctx context.Context, bID uint, phone string) (*entities.Client, error) {
	if m.findClientByPhoneFn != nil {
		return m.findClientByPhoneFn(ctx, bID, phone)
	}
	return nil, nil
}

func (m *mockRepo) FindClientByDNI(_ context.Context, _ uint, _ string) (*entities.Client, error) {
	return nil, nil
}

func (m *mockRepo) FindClientByEmail(_ context.Context, _ uint, _ string) (*entities.Client, error) {
	return nil, nil
}

func (m *mockRepo) UpdateClientFields(_ context.Context, _ uint, _ map[string]any) error {
	return nil
}

func testLogger() log.ILogger {
	nop := zerolog.Nop()
	return log.NewFromZerolog(nop)
}

func newTestUseCase(repo *mockRepo) *UseCase {
	return &UseCase{repo: repo, log: testLogger()}
}

func TestGetCustomerSummary(t *testing.T) {
	now := time.Now()
	expected := &entities.CustomerSummary{
		ID: 1, CustomerID: 10, BusinessID: 1, TotalOrders: 5,
		DeliveredOrders: 3, TotalSpent: 100000, AvgTicket: 20000,
		FirstOrderAt: &now, LastOrderAt: &now,
	}

	repo := &mockRepo{
		getCustomerSummaryFn: func(_ context.Context, bID, cID uint) (*entities.CustomerSummary, error) {
			if bID != 1 || cID != 10 {
				t.Fatalf("expected businessID=1, customerID=10, got %d, %d", bID, cID)
			}
			return expected, nil
		},
	}
	uc := newTestUseCase(repo)

	result, err := uc.GetCustomerSummary(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalOrders != 5 {
		t.Errorf("expected TotalOrders=5, got %d", result.TotalOrders)
	}
	if result.TotalSpent != 100000 {
		t.Errorf("expected TotalSpent=100000, got %f", result.TotalSpent)
	}
}

func TestGetCustomerSummary_Error(t *testing.T) {
	repo := &mockRepo{
		getCustomerSummaryFn: func(_ context.Context, _, _ uint) (*entities.CustomerSummary, error) {
			return nil, errors.New("not found")
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.GetCustomerSummary(context.Background(), 1, 999)
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected 'not found' error, got %v", err)
	}
}

func TestListCustomerAddresses(t *testing.T) {
	expected := []entities.CustomerAddress{
		{ID: 1, CustomerID: 10, Street: "Calle 100", City: "Bogota"},
		{ID: 2, CustomerID: 10, Street: "Cra 50", City: "Medellin"},
	}
	repo := &mockRepo{
		listCustomerAddressesFn: func(_ context.Context, p dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error) {
			if p.CustomerID != 10 || p.BusinessID != 1 {
				t.Fatalf("wrong params: %+v", p)
			}
			return expected, 2, nil
		},
	}
	uc := newTestUseCase(repo)

	params := dtos.ListCustomerAddressesParams{CustomerID: 10, BusinessID: 1, Page: 1, PageSize: 10}
	result, total, err := uc.ListCustomerAddresses(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(result))
	}
	if total != 2 {
		t.Errorf("expected total=2, got %d", total)
	}
}

func TestListCustomerProducts(t *testing.T) {
	expected := []entities.CustomerProductHistory{
		{ID: 1, ProductID: "PRD_001", ProductName: "Creatina", TimesOrdered: 3, TotalSpent: 150000},
	}
	repo := &mockRepo{
		listCustomerProductsFn: func(_ context.Context, p dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error) {
			if p.CustomerID != 10 {
				t.Fatalf("wrong customerID: %d", p.CustomerID)
			}
			return expected, 1, nil
		},
	}
	uc := newTestUseCase(repo)

	params := dtos.ListCustomerProductsParams{CustomerID: 10, BusinessID: 1, Page: 1, PageSize: 10}
	result, total, err := uc.ListCustomerProducts(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].ProductName != "Creatina" {
		t.Errorf("unexpected result: %+v", result)
	}
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
}

func TestListCustomerOrderItems(t *testing.T) {
	expected := []entities.CustomerOrderItem{
		{ID: 1, OrderID: "abc-123", ProductName: "Creatina", Quantity: 2, TotalPrice: 100000, OrderStatus: "delivered"},
	}
	repo := &mockRepo{
		listCustomerOrderItemsFn: func(_ context.Context, p dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error) {
			return expected, 1, nil
		},
	}
	uc := newTestUseCase(repo)

	params := dtos.ListCustomerOrderItemsParams{CustomerID: 10, BusinessID: 1, Page: 1, PageSize: 10}
	result, total, err := uc.ListCustomerOrderItems(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].OrderStatus != "delivered" {
		t.Errorf("unexpected result: %+v", result)
	}
	if total != 1 {
		t.Errorf("expected total=1, got %d", total)
	}
}

func TestProcessOrderEvent_OrderCreated(t *testing.T) {
	var summaryUpserted, addressUpserted, productUpserted, orderItemUpserted bool
	custID := uint(10)
	prodID := "PRD_001"

	repo := &mockRepo{
		upsertCustomerSummaryFn: func(_ context.Context, s *entities.CustomerSummary) error {
			summaryUpserted = true
			if s.CustomerID != custID || s.TotalOrders != 1 {
				t.Errorf("wrong summary: customerID=%d, totalOrders=%d", s.CustomerID, s.TotalOrders)
			}
			return nil
		},
		upsertCustomerAddressFn: func(_ context.Context, a *entities.CustomerAddress) error {
			addressUpserted = true
			if a.Street != "Calle 100" || a.City != "Bogota" {
				t.Errorf("wrong address: %s, %s", a.Street, a.City)
			}
			return nil
		},
		upsertCustomerProductFn: func(_ context.Context, p *entities.CustomerProductHistory) error {
			productUpserted = true
			if p.ProductID != prodID {
				t.Errorf("wrong productID: %s", p.ProductID)
			}
			return nil
		},
		upsertCustomerOrderItemFn: func(_ context.Context, i *entities.CustomerOrderItem) error {
			orderItemUpserted = true
			if i.OrderID != "order-1" {
				t.Errorf("wrong orderID: %s", i.OrderID)
			}
			return nil
		},
	}
	uc := newTestUseCase(repo)

	event := dtos.OrderEventDTO{
		EventType:      "order.created",
		OrderID:        "order-1",
		BusinessID:     1,
		CustomerID:     &custID,
		TotalAmount:    50000,
		Platform:       "shopify",
		Status:         "pending",
		ShippingStreet: "Calle 100",
		ShippingCity:   "Bogota",
		OrderNumber:    "#80001",
		OrderedAt:      time.Now(),
		Items: []dtos.OrderEventItemDTO{
			{ProductID: &prodID, ProductName: "Creatina", Quantity: 2, UnitPrice: 25000, TotalPrice: 50000},
		},
	}

	err := uc.ProcessOrderEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !summaryUpserted {
		t.Error("summary was not upserted")
	}
	if !addressUpserted {
		t.Error("address was not upserted")
	}
	if !productUpserted {
		t.Error("product was not upserted")
	}
	if !orderItemUpserted {
		t.Error("order item was not upserted")
	}
}

func TestProcessOrderEvent_StatusChanged(t *testing.T) {
	var summaryUpdated, statusUpdated bool
	custID := uint(10)

	repo := &mockRepo{
		upsertCustomerSummaryFn: func(_ context.Context, s *entities.CustomerSummary) error {
			summaryUpdated = true
			if s.DeliveredOrders != 1 {
				t.Errorf("expected DeliveredOrders=1, got %d", s.DeliveredOrders)
			}
			if s.InProgressOrders != -1 {
				t.Errorf("expected InProgressOrders=-1, got %d", s.InProgressOrders)
			}
			return nil
		},
		updateOrderItemsStatusFn: func(_ context.Context, orderID string, status string) error {
			statusUpdated = true
			if status != "delivered" {
				t.Errorf("expected status=delivered, got %s", status)
			}
			return nil
		},
	}
	uc := newTestUseCase(repo)

	event := dtos.OrderEventDTO{
		EventType:      "order.status_changed",
		OrderID:        "order-1",
		BusinessID:     1,
		CustomerID:     &custID,
		PreviousStatus: "processing",
		CurrentStatus:  "delivered",
		TotalAmount:    50000,
		OrderedAt:      time.Now(),
	}

	err := uc.ProcessOrderEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !summaryUpdated {
		t.Error("summary was not updated")
	}
	if !statusUpdated {
		t.Error("order items status was not updated")
	}
}

func TestProcessOrderEvent_NoCustomerID_ResolvesViaPhone(t *testing.T) {
	var summaryUpserted bool
	repo := &mockRepo{
		findClientByPhoneFn: func(_ context.Context, bID uint, phone string) (*entities.Client, error) {
			if phone != "3001234567" {
				t.Errorf("wrong phone: %s", phone)
			}
			return &entities.Client{ID: 10, BusinessID: bID}, nil
		},
		upsertCustomerSummaryFn: func(_ context.Context, s *entities.CustomerSummary) error {
			summaryUpserted = true
			if s.CustomerID != 10 {
				t.Errorf("expected customerID=10, got %d", s.CustomerID)
			}
			return nil
		},
	}
	uc := newTestUseCase(repo)

	event := dtos.OrderEventDTO{
		EventType:     "order.created",
		OrderID:       "order-2",
		BusinessID:    1,
		CustomerPhone: "3001234567",
		TotalAmount:   30000,
		Status:        "pending",
		OrderedAt:     time.Now(),
	}

	err := uc.ProcessOrderEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !summaryUpserted {
		t.Error("summary should have been upserted after phone resolution")
	}
}

func TestProcessOrderEvent_NoCustomerID_NoPhone_Skips(t *testing.T) {
	repo := &mockRepo{
		upsertCustomerSummaryFn: func(_ context.Context, _ *entities.CustomerSummary) error {
			t.Fatal("should not upsert summary when customer cannot be resolved")
			return nil
		},
	}
	uc := newTestUseCase(repo)

	event := dtos.OrderEventDTO{
		EventType:  "order.created",
		OrderID:    "order-3",
		BusinessID: 1,
		OrderedAt:  time.Now(),
	}

	err := uc.ProcessOrderEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProcessOrderEvent_SkipsItemsWithoutProductID(t *testing.T) {
	custID := uint(10)
	var productUpsertCount int

	repo := &mockRepo{
		upsertCustomerProductFn: func(_ context.Context, _ *entities.CustomerProductHistory) error {
			productUpsertCount++
			return nil
		},
	}
	uc := newTestUseCase(repo)

	event := dtos.OrderEventDTO{
		EventType:  "order.created",
		OrderID:    "order-4",
		BusinessID: 1,
		CustomerID: &custID,
		Status:     "pending",
		OrderedAt:  time.Now(),
		Items: []dtos.OrderEventItemDTO{
			{ProductID: nil, ProductName: "Sin ID", Quantity: 1},
		},
	}

	err := uc.ProcessOrderEvent(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if productUpsertCount != 0 {
		t.Errorf("expected 0 product upserts for nil productID, got %d", productUpsertCount)
	}
}

func TestIsTerminalStatus(t *testing.T) {
	cases := []struct {
		status   string
		expected bool
	}{
		{"delivered", true},
		{"cancelled", true},
		{"voided", true},
		{"refunded", true},
		{"partially_refunded", true},
		{"pending", false},
		{"processing", false},
		{"shipped", false},
	}
	for _, tc := range cases {
		if got := isTerminalStatus(tc.status); got != tc.expected {
			t.Errorf("isTerminalStatus(%q) = %v, want %v", tc.status, got, tc.expected)
		}
	}
}

func TestIsCancelledStatus(t *testing.T) {
	cases := []struct {
		status   string
		expected bool
	}{
		{"cancelled", true},
		{"voided", true},
		{"refunded", true},
		{"delivered", false},
		{"pending", false},
		{"partially_refunded", false},
	}
	for _, tc := range cases {
		if got := isCancelledStatus(tc.status); got != tc.expected {
			t.Errorf("isCancelledStatus(%q) = %v, want %v", tc.status, got, tc.expected)
		}
	}
}
