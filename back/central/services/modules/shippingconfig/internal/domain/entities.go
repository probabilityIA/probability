package domain

import (
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/shared/shippingpkg"
)

const (
	StrategyProductDimensions = shippingpkg.StrategyProductDimensions
	StrategyStandardBox       = shippingpkg.StrategyStandardBox
)

type Box = shippingpkg.Box

type DirectIntegration struct {
	Enabled       bool   `json:"enabled"`
	IntegrationID *uint  `json:"integration_id"`
	Status        string `json:"status"`
}

type CarrierSetting struct {
	Code         string            `json:"code"`
	Enabled      bool              `json:"enabled"`
	AllowCOD     bool              `json:"allow_cod"`
	AllowPrepaid bool              `json:"allow_prepaid"`
	Direct       DirectIntegration `json:"direct"`
}

type ShippingConfig struct {
	ID              uint             `json:"id"`
	BusinessID      uint             `json:"business_id"`
	WarehouseID     *uint            `json:"warehouse_id"`
	PackageStrategy string           `json:"package_strategy"`
	Boxes           []Box            `json:"boxes"`
	Carriers        []CarrierSetting `json:"carriers"`
	AlwaysInsure    bool             `json:"always_insure"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

func (c *ShippingConfig) CarrierFor(code string) *CarrierSetting {
	if c == nil {
		return nil
	}
	target := NormalizeCarrier(code)
	for i := range c.Carriers {
		if NormalizeCarrier(c.Carriers[i].Code) == target {
			return &c.Carriers[i]
		}
	}
	return nil
}

func (c *ShippingConfig) EnabledCarrierCodes(cod bool) []string {
	if c == nil {
		return nil
	}
	out := make([]string, 0, len(c.Carriers))
	for i := range c.Carriers {
		cs := c.Carriers[i]
		if !cs.Enabled {
			continue
		}
		if cod && !cs.AllowCOD {
			continue
		}
		if !cod && !cs.AllowPrepaid {
			continue
		}
		out = append(out, NormalizeCarrier(cs.Code))
	}
	return out
}

type WarehouseOrigin struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	City      string `json:"city"`
	State     string `json:"state"`
	Phone     string `json:"phone"`
	IsDefault bool   `json:"is_default"`
	IsActive  bool   `json:"is_active"`
	HasConfig bool   `json:"has_config"`
}

type Carrier struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	DirectAvailable bool   `json:"direct_available"`
}

func AvailableCarriers() []Carrier {
	return []Carrier{
		{Code: "ENVIA", Name: "Envia"},
		{Code: "COORDINADORA", Name: "Coordinadora"},
		{Code: "INTERRAPIDISIMO", Name: "Interrapidisimo"},
		{Code: "SERVIENTREGA", Name: "Servientrega"},
		{Code: "TCC", Name: "TCC"},
		{Code: "DEPRISA", Name: "Deprisa"},
		{Code: "99MINUTOS", Name: "99 Minutos"},
		{Code: "PIBOX", Name: "Pibox"},
		{Code: "MELONN", Name: "Melonn"},
	}
}

func NormalizeCarrier(name string) string {
	return strings.ToUpper(strings.TrimSpace(name))
}

func (c *ShippingConfig) SelectBox(totalQuantity int, itemLength, itemWidth, itemHeight float64) *Box {
	if c == nil {
		return nil
	}
	return shippingpkg.SelectBox(c.PackageStrategy, c.Boxes, totalQuantity, itemLength, itemWidth, itemHeight)
}

const (
	DirectStatusUnavailable = "unavailable"
	DirectStatusPending     = "pending"
	DirectStatusActive      = "active"
)

func DirectIntegrationAvailable(code string) bool {
	return false
}

func IsKnownCarrier(code string) bool {
	target := NormalizeCarrier(code)
	for _, c := range AvailableCarriers() {
		if c.Code == target {
			return true
		}
	}
	return false
}
