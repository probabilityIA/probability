package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/shippingconfig/internal/domain"
)

type UseCase struct {
	repo domain.IRepository
}

func New(repo domain.IRepository) domain.IUseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetOverview(ctx context.Context, businessID uint) (*domain.Overview, error) {
	configs, err := uc.repo.ListConfigs(ctx, businessID)
	if err != nil {
		return nil, err
	}

	warehouses, err := uc.repo.ListWarehouseOrigins(ctx, businessID)
	if err != nil {
		return nil, err
	}

	out := &domain.Overview{
		Warehouses: warehouses,
		Overrides:  make([]domain.ShippingConfig, 0, len(configs)),
		Carriers:   domain.AvailableCarriers(),
	}

	withConfig := make(map[uint]bool)
	for i := range configs {
		cfg := configs[i]
		if cfg.WarehouseID == nil {
			c := cfg
			out.Business = &c
			continue
		}
		withConfig[*cfg.WarehouseID] = true
		out.Overrides = append(out.Overrides, cfg)
	}

	for i := range out.Warehouses {
		out.Warehouses[i].HasConfig = withConfig[out.Warehouses[i].ID]
	}

	if out.Business == nil {
		out.Business = &domain.ShippingConfig{
			BusinessID:      businessID,
			PackageStrategy: domain.StrategyProductDimensions,
			Boxes:           []domain.Box{},
			Carriers:        defaultCarrierSettings(),
		}
	}

	return out, nil
}

func (uc *UseCase) SaveBusinessConfig(ctx context.Context, businessID uint, req domain.SaveConfigRequest) (*domain.ShippingConfig, error) {
	return uc.save(ctx, businessID, nil, req)
}

func (uc *UseCase) SaveWarehouseConfig(ctx context.Context, businessID, warehouseID uint, req domain.SaveConfigRequest) (*domain.ShippingConfig, error) {
	if err := uc.assertWarehouseBelongs(ctx, businessID, warehouseID); err != nil {
		return nil, err
	}
	return uc.save(ctx, businessID, &warehouseID, req)
}

func (uc *UseCase) save(ctx context.Context, businessID uint, warehouseID *uint, req domain.SaveConfigRequest) (*domain.ShippingConfig, error) {
	strategy := req.PackageStrategy
	if strategy != domain.StrategyStandardBox && strategy != domain.StrategyProductDimensions {
		return nil, fmt.Errorf("estrategia de empaque invalida: %s", req.PackageStrategy)
	}

	boxes := req.Boxes
	if boxes == nil {
		boxes = []domain.Box{}
	}
	if strategy == domain.StrategyStandardBox && len(boxes) == 0 {
		return nil, fmt.Errorf("la estrategia de caja estandar requiere al menos una caja")
	}
	for i := range boxes {
		if boxes[i].Name == "" {
			return nil, fmt.Errorf("cada caja debe tener nombre")
		}
	}

	carriers, err := normalizeCarrierSettings(req.Carriers)
	if err != nil {
		return nil, err
	}

	cfg := &domain.ShippingConfig{
		BusinessID:      businessID,
		WarehouseID:     warehouseID,
		PackageStrategy: strategy,
		Boxes:           boxes,
		Carriers:        carriers,
	}

	if err := uc.repo.UpsertConfig(ctx, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (uc *UseCase) RemoveWarehouseConfig(ctx context.Context, businessID, warehouseID uint) error {
	if err := uc.assertWarehouseBelongs(ctx, businessID, warehouseID); err != nil {
		return err
	}
	return uc.repo.DeleteConfig(ctx, businessID, warehouseID)
}

func (uc *UseCase) SetDefaultWarehouse(ctx context.Context, businessID, warehouseID uint) error {
	if err := uc.assertWarehouseBelongs(ctx, businessID, warehouseID); err != nil {
		return err
	}
	return uc.repo.SetDefaultWarehouse(ctx, businessID, warehouseID)
}

func (uc *UseCase) Resolve(ctx context.Context, businessID uint, warehouseID *uint) (*domain.ShippingConfig, error) {
	if warehouseID != nil && *warehouseID > 0 {
		cfg, err := uc.repo.GetConfig(ctx, businessID, warehouseID)
		if err != nil {
			return nil, err
		}
		if cfg != nil {
			return cfg, nil
		}
	}
	return uc.repo.GetConfig(ctx, businessID, nil)
}

func (uc *UseCase) assertWarehouseBelongs(ctx context.Context, businessID, warehouseID uint) error {
	warehouses, err := uc.repo.ListWarehouseOrigins(ctx, businessID)
	if err != nil {
		return err
	}
	for _, w := range warehouses {
		if w.ID == warehouseID {
			return nil
		}
	}
	return fmt.Errorf("la bodega no pertenece a este negocio")
}

func defaultCarrierSettings() []domain.CarrierSetting {
	catalog := domain.AvailableCarriers()
	out := make([]domain.CarrierSetting, 0, len(catalog))
	for _, c := range catalog {
		out = append(out, domain.CarrierSetting{
			Code:         c.Code,
			Enabled:      true,
			AllowCOD:     true,
			AllowPrepaid: true,
			Direct:       domain.DirectIntegration{Status: domain.DirectStatusUnavailable},
		})
	}
	return out
}

func normalizeCarrierSettings(in []domain.CarrierSetting) ([]domain.CarrierSetting, error) {
	if len(in) == 0 {
		return defaultCarrierSettings(), nil
	}

	seen := make(map[string]bool, len(in))
	out := make([]domain.CarrierSetting, 0, len(in))
	for _, cs := range in {
		code := domain.NormalizeCarrier(cs.Code)
		if code == "" || seen[code] {
			continue
		}
		if !domain.IsKnownCarrier(code) {
			return nil, fmt.Errorf("transportadora desconocida: %s", cs.Code)
		}
		seen[code] = true

		direct := cs.Direct
		direct.Status = domain.DirectStatusUnavailable
		if direct.Enabled {
			if !domain.DirectIntegrationAvailable(code) {
				direct.Status = domain.DirectStatusPending
			} else if direct.IntegrationID != nil && *direct.IntegrationID > 0 {
				direct.Status = domain.DirectStatusActive
			} else {
				direct.Status = domain.DirectStatusPending
			}
		} else {
			direct.IntegrationID = nil
		}

		out = append(out, domain.CarrierSetting{
			Code:         code,
			Enabled:      cs.Enabled,
			AllowCOD:     cs.AllowCOD,
			AllowPrepaid: cs.AllowPrepaid,
			Direct:       direct,
		})
	}
	return out, nil
}
