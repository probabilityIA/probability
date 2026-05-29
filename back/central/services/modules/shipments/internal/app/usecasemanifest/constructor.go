package usecasemanifest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type UseCaseManifest struct {
	repo domain.IRepository
}

func New(repo domain.IRepository) *UseCaseManifest {
	return &UseCaseManifest{repo: repo}
}

type PendingFilter struct {
	BusinessID      uint
	IncludeChildren bool
	Carrier         string
	Page            int
	PageSize        int
}

type PendingShipmentDTO struct {
	ShipmentID         uint       `json:"shipment_id"`
	OrderID            *string    `json:"order_id,omitempty"`
	OrderNumber        string     `json:"order_number"`
	TrackingNumber     string     `json:"tracking_number"`
	Carrier            string     `json:"carrier"`
	CarrierCode        string     `json:"carrier_code"`
	CustomerName       string     `json:"customer_name"`
	CustomerDocument   string     `json:"customer_document"`
	DestinationAddress string     `json:"destination_address"`
	DestinationCity    string     `json:"destination_city"`
	DestinationState   string     `json:"destination_state"`
	Weight             float64    `json:"weight"`
	DeclaredValue      float64    `json:"declared_value"`
	CodTotal           float64    `json:"cod_total"`
	BusinessID         uint       `json:"business_id"`
	BusinessName       string     `json:"business_name"`
	ShipmentCreatedAt  *time.Time `json:"shipment_created_at,omitempty"`
	OrderCreatedAt     *time.Time `json:"order_created_at,omitempty"`
	ShipmentStatus     string     `json:"shipment_status"`
	OrderStatus        string     `json:"order_status"`
}

type PendingPage struct {
	Carrier    string               `json:"carrier"`
	Data       []PendingShipmentDTO `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

type CarrierOption struct {
	Carrier string `json:"carrier"`
	Count   int64  `json:"count"`
}

func (uc *UseCaseManifest) ListCarriers(ctx context.Context, businessID uint, includeChildren bool) ([]CarrierOption, error) {
	rows, err := uc.repo.ListPendingCarriers(ctx, businessID, includeChildren)
	if err != nil {
		return nil, err
	}
	out := make([]CarrierOption, 0, len(rows))
	for _, r := range rows {
		out = append(out, CarrierOption{Carrier: r.Carrier, Count: r.Count})
	}
	return out, nil
}

func (uc *UseCaseManifest) ListPending(ctx context.Context, f PendingFilter) (*PendingPage, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 25
	}
	if f.PageSize > 200 {
		f.PageSize = 200
	}

	rows, total, err := uc.repo.ListPendingForManifest(ctx, domain.ManifestFilter{
		BusinessID:      f.BusinessID,
		IncludeChildren: f.IncludeChildren,
		Carrier:         f.Carrier,
		Page:            f.Page,
		PageSize:        f.PageSize,
	})
	if err != nil {
		return nil, err
	}

	data := make([]PendingShipmentDTO, 0, len(rows))
	for _, r := range rows {
		data = append(data, PendingShipmentDTO{
			ShipmentID:         r.ShipmentID,
			OrderID:            r.OrderID,
			OrderNumber:        r.OrderNumber,
			TrackingNumber:     r.TrackingNumber,
			Carrier:            r.Carrier,
			CarrierCode:        r.CarrierCode,
			CustomerName:       r.CustomerName,
			CustomerDocument:   r.CustomerDocument,
			DestinationAddress: r.DestinationAddress,
			DestinationCity:    r.DestinationCity,
			DestinationState:   r.DestinationState,
			Weight:             r.Weight,
			DeclaredValue:      r.DeclaredValue,
			CodTotal:           r.CodTotal,
			BusinessID:         r.BusinessID,
			BusinessName:       r.BusinessName,
			ShipmentCreatedAt:  r.ShipmentCreatedAt,
			OrderCreatedAt:     r.OrderCreatedAt,
			ShipmentStatus:     r.ShipmentStatus,
			OrderStatus:        r.OrderStatus,
		})
	}

	totalPages := int((total + int64(f.PageSize) - 1) / int64(f.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	return &PendingPage{
		Carrier:    f.Carrier,
		Data:       data,
		Total:      total,
		Page:       f.Page,
		PageSize:   f.PageSize,
		TotalPages: totalPages,
	}, nil
}

type GeneratePDFInput struct {
	BusinessID  uint
	ShipmentIDs []uint
	Carrier     string
	UserName    string
}

type GeneratedManifest struct {
	Carrier  string
	Filename string
	PDF      []byte
}

func (uc *UseCaseManifest) GeneratePDF(ctx context.Context, in GeneratePDFInput) ([]GeneratedManifest, error) {
	if in.BusinessID == 0 {
		return nil, fmt.Errorf("business_id requerido")
	}
	if len(in.ShipmentIDs) == 0 {
		return nil, fmt.Errorf("no hay envios seleccionados")
	}

	rows, _, err := uc.repo.ListPendingForManifest(ctx, domain.ManifestFilter{
		BusinessID:      in.BusinessID,
		IncludeChildren: true,
		Carrier:         in.Carrier,
	})
	if err != nil {
		return nil, err
	}

	idset := map[uint]bool{}
	for _, id := range in.ShipmentIDs {
		idset[id] = true
	}

	selected := make([]domain.ManifestShipmentRow, 0, len(in.ShipmentIDs))
	for _, r := range rows {
		if idset[r.ShipmentID] {
			selected = append(selected, r)
		}
	}
	if len(selected) == 0 {
		return nil, fmt.Errorf("ninguno de los envios seleccionados esta pendiente")
	}

	byCarrier := map[string][]domain.ManifestShipmentRow{}
	carrierOrder := []string{}
	for _, r := range selected {
		k := strings.TrimSpace(r.Carrier)
		if k == "" {
			k = "Sin asignar"
		}
		if _, ok := byCarrier[k]; !ok {
			carrierOrder = append(carrierOrder, k)
		}
		byCarrier[k] = append(byCarrier[k], r)
	}

	biz, err := uc.repo.GetBusinessForManifest(ctx, in.BusinessID)
	if err != nil {
		return nil, err
	}
	if biz == nil {
		return nil, fmt.Errorf("negocio no encontrado")
	}

	out := make([]GeneratedManifest, 0, len(carrierOrder))
	now := time.Now()
	for _, c := range carrierOrder {
		input := domain.ManifestPDFInput{
			BusinessID:   biz.ID,
			BusinessName: biz.Name,
			BusinessCode: biz.Code,
			OriginCity:   biz.City,
			GeneratedAt:  now,
			GeneratedBy:  in.UserName,
			ManifestNo:   fmt.Sprintf("M-%d-%d", biz.ID, now.Unix()),
			Carrier:      c,
			Rows:         byCarrier[c],
		}
		pdf, err := buildManifestPDF(input)
		if err != nil {
			return nil, err
		}
		safeCarrier := strings.ReplaceAll(strings.ToLower(c), " ", "-")
		out = append(out, GeneratedManifest{
			Carrier:  c,
			Filename: fmt.Sprintf("manifiesto-%s-%s.pdf", safeCarrier, now.Format("20060102-1504")),
			PDF:      pdf,
		})
	}
	return out, nil
}
