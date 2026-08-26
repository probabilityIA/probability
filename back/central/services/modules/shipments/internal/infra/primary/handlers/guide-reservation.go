package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

type guideConflict struct {
	status int
	body   gin.H
}

func activeGuideConflict(s *domain.Shipment) *guideConflict {
	tracking := ""
	if s.TrackingNumber != nil {
		tracking = *s.TrackingNumber
	}
	return &guideConflict{
		status: http.StatusConflict,
		body: gin.H{
			"error":           "La orden ya tiene una guia activa. Cancela la guia existente antes de generar una nueva.",
			"shipment_id":     s.ID,
			"tracking_number": tracking,
		},
	}
}

func inFlightConflict(shipmentID uint) *guideConflict {
	return &guideConflict{
		status: http.StatusConflict,
		body: gin.H{
			"error":       "La guia de esta orden ya se esta generando. Espera el resultado antes de intentar de nuevo.",
			"shipment_id": shipmentID,
			"in_flight":   true,
		},
	}
}

func needsVerificationConflict(shipmentID uint) *guideConflict {
	return &guideConflict{
		status: http.StatusConflict,
		body: gin.H{
			"error":              "No pudimos confirmar si la transportadora creo esta guia. Verificala con la transportadora antes de generar otra para no duplicarla.",
			"shipment_id":        shipmentID,
			"needs_verification": true,
		},
	}
}

func (h *Handlers) reserveShipmentForGuide(ctx context.Context, req *domain.CreateShipmentRequest) (uint, *guideConflict, error) {
	var shipmentID uint
	var conflict *guideConflict

	reserve := func() error {
		if req.OrderID != nil && *req.OrderID != "" {
			existing, _ := h.uc.Repo().GetShipmentsByOrderID(ctx, *req.OrderID)
			now := time.Now()
			for i := range existing {
				s := &existing[i]
				if s.HasActiveGuide() {
					conflict = activeGuideConflict(s)
					return nil
				}
				if s.GuideInFlight(now) {
					conflict = inFlightConflict(s.ID)
					return nil
				}
				if s.NeedsVerification() {
					conflict = needsVerificationConflict(s.ID)
					return nil
				}
			}
			for i := range existing {
				if existing[i].ReusableForGuide(now) {
					shipmentID = existing[i].ID
					break
				}
			}
		}

		if shipmentID == 0 {
			resp, err := h.uc.CreateShipment(ctx, req)
			if err != nil {
				return err
			}
			shipmentID = resp.ID
		} else {
			updateReq := &domain.UpdateShipmentRequest{
				UpdatedBy:     req.CreatedBy,
				UpdatedByName: req.CreatedByName,
				TotalCost:     req.TotalCost,
				CodCarrierFee: req.CodCarrierFee,
				Carrier:       req.Carrier,
				CarrierCode:   req.CarrierCode,
				Weight:        req.Weight,
				Height:        req.Height,
				Width:         req.Width,
				Length:        req.Length,
			}
			if _, err := h.uc.UpdateShipment(ctx, shipmentID, updateReq); err != nil {
				return err
			}
		}

		locked, err := h.uc.Repo().MarkShipmentGenerating(ctx, shipmentID, time.Now().Add(-domain.GuideGenerationStaleAfter))
		if err != nil {
			return err
		}
		if !locked {
			conflict = inFlightConflict(shipmentID)
		}
		return nil
	}

	var err error
	if req.OrderID != nil && *req.OrderID != "" {
		err = h.uc.Repo().WithOrderGuideLock(ctx, *req.OrderID, reserve)
	} else {
		err = reserve()
	}
	if err != nil {
		return 0, nil, err
	}
	return shipmentID, conflict, nil
}
