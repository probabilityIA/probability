package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type RepositoryMock struct {
	CreateShipmentFn                  func(ctx context.Context, shipment *domain.Shipment) error
	GetShipmentByIDFn                 func(ctx context.Context, id uint) (*domain.Shipment, error)
	GetShipmentByTrackingNumberFn     func(ctx context.Context, trackingNumber string) (*domain.Shipment, error)
	GetShipmentsByOrderIDFn           func(ctx context.Context, orderID string) ([]domain.Shipment, error)
	ListShipmentsFn                   func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error)
	UpdateShipmentFn                  func(ctx context.Context, shipment *domain.Shipment) error
	DeleteShipmentFn                  func(ctx context.Context, id uint) error
	ShipmentExistsFn                  func(ctx context.Context, orderID string, trackingNumber string) (bool, error)
	GetActiveShippingCarrierFn        func(ctx context.Context, businessID uint) (*domain.CarrierInfo, error)
	GetBusinessNameFn                 func(ctx context.Context, businessID uint) (string, error)
	GetOrderBusinessIDFn              func(ctx context.Context, orderUUID string) (uint, error)
	GetShipmentBusinessIDByTrackingFn func(ctx context.Context, trackingNumber string) (uint, error)
	GetShipmentBusinessIDByIDFn       func(ctx context.Context, shipmentID uint) (uint, error)
	UpdateOrderGuideLinkFn            func(ctx context.Context, orderID string, guideLink string, trackingNumber string, carrier string) error
	UpdateOrderStatusByOrderIDFn      func(ctx context.Context, orderID string, status string) error
	ClearOrderGuideDataFn             func(ctx context.Context, orderID string) error
	EnsureAllBusinessesActiveFn       func(ctx context.Context) error
	GetOrderIntegrationIDFn           func(ctx context.Context, orderUUID string) (uint, error)
	CreateOriginAddressFn             func(ctx context.Context, address *domain.OriginAddress) error
	GetOriginAddressByIDFn            func(ctx context.Context, id uint) (*domain.OriginAddress, error)
	ListOriginAddressesByBusinessFn   func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error)
	GetDefaultOriginAddressFn         func(ctx context.Context, businessID uint) (*domain.OriginAddress, error)
	UpdateOriginAddressFn             func(ctx context.Context, address *domain.OriginAddress) error
	DeleteOriginAddressFn             func(ctx context.Context, id uint) error
	SetDefaultOriginAddressFn         func(ctx context.Context, businessID, addressID uint) error
}

func (m *RepositoryMock) CreateShipment(ctx context.Context, shipment *domain.Shipment) error {
	if m.CreateShipmentFn != nil {
		return m.CreateShipmentFn(ctx, shipment)
	}
	return nil
}

func (m *RepositoryMock) GetShipmentByID(ctx context.Context, id uint) (*domain.Shipment, error) {
	if m.GetShipmentByIDFn != nil {
		return m.GetShipmentByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	if m.GetShipmentByTrackingNumberFn != nil {
		return m.GetShipmentByTrackingNumberFn(ctx, trackingNumber)
	}
	return nil, nil
}

func (m *RepositoryMock) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	if m.GetShipmentsByOrderIDFn != nil {
		return m.GetShipmentsByOrderIDFn(ctx, orderID)
	}
	return nil, nil
}

func (m *RepositoryMock) ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]domain.Shipment, int64, error) {
	if m.ListShipmentsFn != nil {
		return m.ListShipmentsFn(ctx, page, pageSize, filters)
	}
	return nil, 0, nil
}

func (m *RepositoryMock) UpdateShipment(ctx context.Context, shipment *domain.Shipment) error {
	if m.UpdateShipmentFn != nil {
		return m.UpdateShipmentFn(ctx, shipment)
	}
	return nil
}

func (m *RepositoryMock) DeleteShipment(ctx context.Context, id uint) error {
	if m.DeleteShipmentFn != nil {
		return m.DeleteShipmentFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) ShipmentExists(ctx context.Context, orderID string, trackingNumber string) (bool, error) {
	if m.ShipmentExistsFn != nil {
		return m.ShipmentExistsFn(ctx, orderID, trackingNumber)
	}
	return false, nil
}

func (m *RepositoryMock) GetActiveShippingCarrier(ctx context.Context, businessID uint) (*domain.CarrierInfo, error) {
	if m.GetActiveShippingCarrierFn != nil {
		return m.GetActiveShippingCarrierFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetBusinessName(ctx context.Context, businessID uint) (string, error) {
	if m.GetBusinessNameFn != nil {
		return m.GetBusinessNameFn(ctx, businessID)
	}
	return "", nil
}

func (m *RepositoryMock) GetOrderBusinessID(ctx context.Context, orderUUID string) (uint, error) {
	if m.GetOrderBusinessIDFn != nil {
		return m.GetOrderBusinessIDFn(ctx, orderUUID)
	}
	return 0, nil
}

func (m *RepositoryMock) GetShipmentBusinessIDByTracking(ctx context.Context, trackingNumber string) (uint, error) {
	if m.GetShipmentBusinessIDByTrackingFn != nil {
		return m.GetShipmentBusinessIDByTrackingFn(ctx, trackingNumber)
	}
	return 0, nil
}

func (m *RepositoryMock) GetShipmentBusinessIDByID(ctx context.Context, shipmentID uint) (uint, error) {
	if m.GetShipmentBusinessIDByIDFn != nil {
		return m.GetShipmentBusinessIDByIDFn(ctx, shipmentID)
	}
	return 0, nil
}

func (m *RepositoryMock) UpdateOrderGuideLink(ctx context.Context, orderID string, guideLink string, trackingNumber string, carrier string) error {
	if m.UpdateOrderGuideLinkFn != nil {
		return m.UpdateOrderGuideLinkFn(ctx, orderID, guideLink, trackingNumber, carrier)
	}
	return nil
}

func (m *RepositoryMock) UpdateOrderStatusByOrderID(ctx context.Context, orderID string, status string) error {
	if m.UpdateOrderStatusByOrderIDFn != nil {
		return m.UpdateOrderStatusByOrderIDFn(ctx, orderID, status)
	}
	return nil
}

func (m *RepositoryMock) ClearOrderGuideData(ctx context.Context, orderID string) error {
	if m.ClearOrderGuideDataFn != nil {
		return m.ClearOrderGuideDataFn(ctx, orderID)
	}
	return nil
}

func (m *RepositoryMock) EnsureAllBusinessesActive(ctx context.Context) error {
	if m.EnsureAllBusinessesActiveFn != nil {
		return m.EnsureAllBusinessesActiveFn(ctx)
	}
	return nil
}

func (m *RepositoryMock) GetOrderIntegrationID(ctx context.Context, orderUUID string) (uint, error) {
	if m.GetOrderIntegrationIDFn != nil {
		return m.GetOrderIntegrationIDFn(ctx, orderUUID)
	}
	return 0, nil
}

func (m *RepositoryMock) CreateOriginAddress(ctx context.Context, address *domain.OriginAddress) error {
	if m.CreateOriginAddressFn != nil {
		return m.CreateOriginAddressFn(ctx, address)
	}
	return nil
}

func (m *RepositoryMock) GetOriginAddressByID(ctx context.Context, id uint) (*domain.OriginAddress, error) {
	if m.GetOriginAddressByIDFn != nil {
		return m.GetOriginAddressByIDFn(ctx, id)
	}
	return nil, nil
}

func (m *RepositoryMock) ListOriginAddressesByBusiness(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
	if m.ListOriginAddressesByBusinessFn != nil {
		return m.ListOriginAddressesByBusinessFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) GetDefaultOriginAddress(ctx context.Context, businessID uint) (*domain.OriginAddress, error) {
	if m.GetDefaultOriginAddressFn != nil {
		return m.GetDefaultOriginAddressFn(ctx, businessID)
	}
	return nil, nil
}

func (m *RepositoryMock) UpdateOriginAddress(ctx context.Context, address *domain.OriginAddress) error {
	if m.UpdateOriginAddressFn != nil {
		return m.UpdateOriginAddressFn(ctx, address)
	}
	return nil
}

func (m *RepositoryMock) DeleteOriginAddress(ctx context.Context, id uint) error {
	if m.DeleteOriginAddressFn != nil {
		return m.DeleteOriginAddressFn(ctx, id)
	}
	return nil
}

func (m *RepositoryMock) SetDefaultOriginAddress(ctx context.Context, businessID, addressID uint) error {
	if m.SetDefaultOriginAddressFn != nil {
		return m.SetDefaultOriginAddressFn(ctx, businessID, addressID)
	}
	return nil
}
