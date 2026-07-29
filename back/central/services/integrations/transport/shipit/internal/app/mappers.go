package app

import (
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

const maxReferenceLength = 15

func buildReference(req domain.GuideRequest) string {
	ref := req.MyShipmentReference
	if ref == "" {
		ref = req.ExternalOrderID
	}
	if ref == "" && req.OrderUUID != "" {
		ref = req.OrderUUID
	}
	if len(ref) > maxReferenceLength {
		ref = ref[len(ref)-maxReferenceLength:]
	}
	return ref
}

func buildSizes(packages []domain.Package) domain.Sizes {
	if len(packages) == 0 {
		return domain.Sizes{Width: 10, Height: 10, Length: 10, Weight: 1}
	}
	sizes := domain.Sizes{
		Width:  packages[0].Width,
		Height: packages[0].Height,
		Length: packages[0].Length,
	}
	for _, pkg := range packages {
		sizes.Weight += pkg.Weight
		if pkg.Width > sizes.Width {
			sizes.Width = pkg.Width
		}
		if pkg.Height > sizes.Height {
			sizes.Height = pkg.Height
		}
		if pkg.Length > sizes.Length {
			sizes.Length = pkg.Length
		}
	}
	if sizes.Weight <= 0 {
		sizes.Weight = 1
	}
	return sizes
}

func splitStreetNumber(address string) (string, string) {
	trimmed := strings.TrimSpace(address)
	if trimmed == "" {
		return "", "S/N"
	}
	fields := strings.Fields(trimmed)
	if len(fields) < 2 {
		return trimmed, "S/N"
	}
	last := fields[len(fields)-1]
	hasDigit := false
	for _, r := range last {
		if r >= '0' && r <= '9' {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return trimmed, "S/N"
	}
	return strings.Join(fields[:len(fields)-1], " "), last
}

func buildShipmentBody(req domain.GuideRequest, communeID uint, communeName string, sandbox bool) domain.ShipmentBody {
	street, number := splitStreetNumber(req.Destination.Address)
	fullName := strings.TrimSpace(req.Destination.FirstName + " " + req.Destination.LastName)

	items := len(req.Packages)
	if items == 0 {
		items = 1
	}

	courier := domain.CourierSpec{
		Client:    strings.ToLower(strings.TrimSpace(req.Carrier)),
		Selected:  strings.TrimSpace(req.Carrier) != "",
		Payable:   req.CODValue > 0,
		Algorithm: 1,
	}

	reference := buildReference(req)
	detail := req.Description
	if detail == "" {
		detail = "Mercancia orden " + reference
	}

	return domain.ShipmentBody{
		Kind:      0,
		Platform:  2,
		Reference: reference,
		Items:     items,
		Sandbox:   sandbox,
		Sizes:     buildSizes(req.Packages),
		Destiny: domain.Destiny{
			Street:      street,
			Number:      number,
			Complement:  req.Destination.Reference,
			CommuneID:   communeID,
			CommuneName: communeName,
			FullName:    fullName,
			Email:       req.Destination.Email,
			Phone:       req.Destination.Phone,
			Kind:        "home_delivery",
		},
		Courier: courier,
		Insurance: domain.Insurance{
			TicketNumber: reference,
			TicketAmount: req.ContentValue,
			Detail:       detail,
			Extra:        false,
		},
	}
}

func mapGenerateResponse(req domain.GuideRequest, resp *domain.ShipmentResponse) *domain.GenerateResponse {
	tracker := ""
	if resp.TrackingNumber != nil && *resp.TrackingNumber != "" {
		tracker = *resp.TrackingNumber
	} else if resp.ShipitCode != "" {
		tracker = resp.ShipitCode
	} else {
		tracker = fmt.Sprintf("SHIPIT-%d", resp.ID)
	}

	labelURL := ""
	switch {
	case resp.TicketPDFURL != nil && *resp.TicketPDFURL != "":
		labelURL = *resp.TicketPDFURL
	case resp.TicketURL != nil && *resp.TicketURL != "":
		labelURL = *resp.TicketURL
	case resp.CourierURL != "":
		labelURL = resp.CourierURL
	case resp.TrackingLink != "" && resp.TrackingLink != "Sin tracking":
		labelURL = resp.TrackingLink
	}

	carrier := resp.CourierForClient
	if carrier == "" {
		carrier = req.Carrier
	}
	if carrier == "" {
		carrier = "shipit"
	}

	return &domain.GenerateResponse{
		Status: "success",
		Data: domain.GenerateData{
			TrackingNumber:   tracker,
			LabelURL:         labelURL,
			MyGuideReference: resp.Reference,
			Carrier:          carrier,
			IDOrder:          resp.ID,
		},
	}
}

func mapRatesResponse(resp *domain.RatesResponse) *domain.QuoteResponse {
	rates := make([]domain.Rate, 0, len(resp.Prices))
	for i, price := range resp.Prices {
		if !price.AvailableToShipping {
			continue
		}
		rates = append(rates, domain.Rate{
			IDRate:       int64(i + 1),
			Product:      price.Name,
			Carrier:      price.Courier.Name,
			Flete:        price.Price,
			DeliveryDays: price.Days,
		})
	}
	return &domain.QuoteResponse{
		Status: "success",
		Data:   domain.QuoteData{Rates: rates},
	}
}

func mapTrackingResponse(resp *domain.TrackingResponse) *domain.TrackOutput {
	status, rawDetail, _ := domain.MapTrackingData(&resp.Data)

	history := make([]domain.TrackEvent, 0, len(resp.Data.Statuses))
	for _, st := range resp.Data.Statuses {
		date := st.CourierUpdatedAt
		if date == "" {
			date = st.CreatedAt
		}
		mapped, _ := domain.MapShipitStatus(st.CourierStatus)
		history = append(history, domain.TrackEvent{
			Date:        date,
			Status:      mapped.String(),
			Description: st.CourierStatus,
			Location:    "",
		})
	}

	return &domain.TrackOutput{
		Status:         status.String(),
		Carrier:        resp.Data.Courier,
		TrackingNumber: resp.Data.Number,
		StatusDetail:   rawDetail,
		IsDelivered:    resp.Data.IsDelivered,
		History:        history,
	}
}
