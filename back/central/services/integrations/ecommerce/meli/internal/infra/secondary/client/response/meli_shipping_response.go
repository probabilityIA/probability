package response

import (
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// MeliShippingDetailResponse es la respuesta de GET /shipments/{id}.
type MeliShippingDetailResponse struct {
	ID              int64                       `json:"id"`
	Status          string                      `json:"status"`
	SubStatus       string                      `json:"substatus"`
	ShipmentType    string                      `json:"shipment_type"`
	DateCreated     *string                     `json:"date_created"`
	ReceiverAddress *MeliReceiverAddressResp    `json:"receiver_address"`
	SenderAddress   *MeliSenderAddressResp      `json:"sender_address"`
	ShippingOption  *MeliShippingOptionResp     `json:"shipping_option"`
}

type MeliReceiverAddressResp struct {
	ID           int64              `json:"id"`
	AddressLine  string             `json:"address_line"`
	StreetName   string             `json:"street_name"`
	StreetNumber string             `json:"street_number"`
	ZipCode      string             `json:"zip_code"`
	City         MeliLocationResp   `json:"city"`
	State        MeliLocationResp   `json:"state"`
	Country      MeliLocationResp   `json:"country"`
	Neighborhood *MeliLocationResp  `json:"neighborhood"`
	Latitude     *float64           `json:"latitude"`
	Longitude    *float64           `json:"longitude"`
	Comment      string             `json:"comment"`
}

type MeliSenderAddressResp struct {
	ID           int64            `json:"id"`
	AddressLine  string           `json:"address_line"`
	StreetName   string           `json:"street_name"`
	StreetNumber string           `json:"street_number"`
	ZipCode      string           `json:"zip_code"`
	City         MeliLocationResp `json:"city"`
	State        MeliLocationResp `json:"state"`
	Country      MeliLocationResp `json:"country"`
}

type MeliLocationResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MeliShippingOptionResp struct {
	ID                    int64                      `json:"id"`
	Name                  string                     `json:"name"`
	CurrencyID            string                     `json:"currency_id"`
	Cost                  float64                    `json:"cost"`
	ListCost              float64                    `json:"list_cost"`
	ShippingMethodID      int64                      `json:"shipping_method_id"`
	DeliveryType          string                     `json:"delivery_type"`
	EstimatedDeliveryTime *MeliEstimatedDeliveryResp `json:"estimated_delivery_time"`
}

type MeliEstimatedDeliveryResp struct {
	Type string  `json:"type"`
	Date *string `json:"date"`
}

func (r *MeliShippingDetailResponse) ToDomain() domain.MeliShippingDetail {
	detail := domain.MeliShippingDetail{
		ID:           r.ID,
		Status:       r.Status,
		SubStatus:    r.SubStatus,
		ShipmentType: r.ShipmentType,
		DateCreated:  parseMeliDatePtr(r.DateCreated),
	}

	if r.ReceiverAddress != nil {
		ra := &domain.MeliReceiverAddress{
			ID:           r.ReceiverAddress.ID,
			AddressLine:  r.ReceiverAddress.AddressLine,
			StreetName:   r.ReceiverAddress.StreetName,
			StreetNumber: r.ReceiverAddress.StreetNumber,
			ZipCode:      r.ReceiverAddress.ZipCode,
			City:         mapLocation(r.ReceiverAddress.City),
			State:        mapLocation(r.ReceiverAddress.State),
			Country:      mapLocation(r.ReceiverAddress.Country),
			Latitude:     r.ReceiverAddress.Latitude,
			Longitude:    r.ReceiverAddress.Longitude,
			Comment:      r.ReceiverAddress.Comment,
		}
		if r.ReceiverAddress.Neighborhood != nil {
			n := mapLocation(*r.ReceiverAddress.Neighborhood)
			ra.Neighborhood = &n
		}
		detail.ReceiverAddress = ra
	}

	if r.SenderAddress != nil {
		detail.SenderAddress = &domain.MeliSenderAddress{
			ID:           r.SenderAddress.ID,
			AddressLine:  r.SenderAddress.AddressLine,
			StreetName:   r.SenderAddress.StreetName,
			StreetNumber: r.SenderAddress.StreetNumber,
			ZipCode:      r.SenderAddress.ZipCode,
			City:         mapLocation(r.SenderAddress.City),
			State:        mapLocation(r.SenderAddress.State),
			Country:      mapLocation(r.SenderAddress.Country),
		}
	}

	if r.ShippingOption != nil {
		opt := &domain.MeliShippingOption{
			ID:               r.ShippingOption.ID,
			Name:             r.ShippingOption.Name,
			CurrencyID:       r.ShippingOption.CurrencyID,
			Cost:             r.ShippingOption.Cost,
			ListCost:         r.ShippingOption.ListCost,
			ShippingMethodID: r.ShippingOption.ShippingMethodID,
			DeliveryType:     r.ShippingOption.DeliveryType,
		}
		if r.ShippingOption.EstimatedDeliveryTime != nil {
			est := &domain.MeliEstimatedDelivery{
				Type: r.ShippingOption.EstimatedDeliveryTime.Type,
				Date: parseMeliDatePtr(r.ShippingOption.EstimatedDeliveryTime.Date),
			}
			opt.EstimatedDeliveryTime = est
		}
		detail.ShippingOption = opt
	}

	return detail
}

func mapLocation(loc MeliLocationResp) domain.MeliLocation {
	return domain.MeliLocation{
		ID:   loc.ID,
		Name: loc.Name,
	}
}
