package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

// VTEXOrderDetailResponse es la respuesta JSON de GET /api/oms/pvt/orders/{orderId}.
type VTEXOrderDetailResponse struct {
	OrderID            string                      `json:"orderId"`
	Sequence           string                      `json:"sequence"`
	MarketplaceOrderID string                      `json:"marketplaceOrderId"`
	Status             string                      `json:"status"`
	StatusDescription  string                      `json:"statusDescription"`
	Value              int                         `json:"value"`
	TotalItems         int                         `json:"totalItems"`
	TotalDiscount      int                         `json:"totalDiscount"`
	TotalFreight       int                         `json:"totalFreight"`
	CreationDate       string                      `json:"creationDate"`
	LastChange         string                      `json:"lastChange"`
	Items              []VTEXItemResponse          `json:"items"`
	ShippingData       *VTEXShippingDataResponse   `json:"shippingData"`
	PaymentData        *VTEXPaymentDataResponse    `json:"paymentData"`
	ClientProfileData  *VTEXClientProfileResponse  `json:"clientProfileData"`
	RatesAndBenefits   []VTEXRateAndBenefitResp    `json:"ratesAndBenefitsData"`
	Sellers            []VTEXSellerResponse        `json:"sellers"`
	Totals             []VTEXTotalResponse         `json:"totals"`
	PackageAttachment  *VTEXPackageAttachmentResp  `json:"packageAttachment"`
}

type VTEXItemResponse struct {
	UniqueID        string  `json:"uniqueId"`
	ID              string  `json:"id"`
	ProductID       string  `json:"productId"`
	EANID           string  `json:"ean"`
	RefID           string  `json:"refId"`
	Name            string  `json:"name"`
	SKUName         string  `json:"skuName"`
	ImageURL        string  `json:"imageUrl"`
	DetailURL       string  `json:"detailUrl"`
	Quantity        int     `json:"quantity"`
	Price           int     `json:"price"`
	ListPrice       int     `json:"listPrice"`
	SellingPrice    int     `json:"sellingPrice"`
	Tax             int     `json:"tax"`
	MeasurementUnit string  `json:"measurementUnit"`
	UnitMultiplier  float64 `json:"unitMultiplier"`
}

type VTEXShippingDataResponse struct {
	Address       *VTEXAddressResponse       `json:"address"`
	LogisticsInfo []VTEXLogisticsInfoResp    `json:"logisticsInfo"`
	SelectedSLA   string                     `json:"selectedSla"`
	TrackingHints []VTEXTrackingHintResp     `json:"trackingHints"`
}

type VTEXAddressResponse struct {
	AddressType    string    `json:"addressType"`
	ReceiverName   string    `json:"receiverName"`
	Street         string    `json:"street"`
	Number         string    `json:"number"`
	Complement     string    `json:"complement"`
	Neighborhood   string    `json:"neighborhood"`
	City           string    `json:"city"`
	State          string    `json:"state"`
	Country        string    `json:"country"`
	PostalCode     string    `json:"postalCode"`
	Reference      string    `json:"reference"`
	GeoCoordinates []float64 `json:"geoCoordinates"`
}

type VTEXLogisticsInfoResp struct {
	ItemIndex            int                        `json:"itemIndex"`
	SelectedSLA          string                     `json:"selectedSla"`
	LockTTL              string                     `json:"lockTTL"`
	Price                int                        `json:"price"`
	ListPrice            int                        `json:"listPrice"`
	SellingPrice         int                        `json:"sellingPrice"`
	DeliveryWindow       *VTEXDeliveryWindowResp    `json:"deliveryWindow"`
	ShippingEstimate     string                     `json:"shippingEstimate"`
	ShippingEstimateDate *string                    `json:"shippingEstimateDate"`
	DeliveryCompany      string                     `json:"deliveryCompany"`
	DeliveryIDs          []VTEXDeliveryIDResp       `json:"deliveryIds"`
}

type VTEXDeliveryWindowResp struct {
	StartDateUTC string `json:"startDateUtc"`
	EndDateUTC   string `json:"endDateUtc"`
	Price        int    `json:"price"`
}

type VTEXDeliveryIDResp struct {
	CourierID   string `json:"courierId"`
	CourierName string `json:"courierName"`
	DockID      string `json:"dockId"`
	Quantity    int    `json:"quantity"`
	WarehouseID string `json:"warehouseId"`
}

type VTEXTrackingHintResp struct {
	CourierName   string `json:"courierName"`
	TrackingID    string `json:"trackingId"`
	TrackingURL   string `json:"trackingUrl"`
	TrackingLabel string `json:"trackingLabel"`
}

type VTEXPaymentDataResponse struct {
	Transactions []VTEXTransactionResp `json:"transactions"`
}

type VTEXTransactionResp struct {
	IsActive      bool                `json:"isActive"`
	TransactionID string              `json:"transactionId"`
	MerchantName  string              `json:"merchantName"`
	Payments      []VTEXPaymentResp   `json:"payments"`
}

type VTEXPaymentResp struct {
	ID                 string            `json:"id"`
	PaymentSystem      string            `json:"paymentSystem"`
	PaymentSystemName  string            `json:"paymentSystemName"`
	Value              int               `json:"value"`
	ReferenceValue     int               `json:"referenceValue"`
	Group              string            `json:"group"`
	ConnectorResponses map[string]string `json:"connectorResponses"`
	InstallmentCount   int               `json:"installments"`
	CardHolder         string            `json:"cardHolder"`
	FirstDigits        string            `json:"firstDigits"`
	LastDigits         string            `json:"lastDigits"`
	URL                string            `json:"url"`
	TID                string            `json:"tid"`
}

type VTEXClientProfileResponse struct {
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	DocumentType  string `json:"documentType"`
	Document      string `json:"document"`
	Phone         string `json:"phone"`
	CorporateName string `json:"corporateName"`
	IsCorporate   bool   `json:"isCorporate"`
}

type VTEXRateAndBenefitResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VTEXSellerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SubSellerID string `json:"subSellerId"`
}

type VTEXTotalResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type VTEXPackageAttachmentResp struct {
	Packages []VTEXPackageResp `json:"packages"`
}

type VTEXPackageResp struct {
	Items          []VTEXPackageItemResp  `json:"items"`
	CourierStatus  *VTEXCourierStatusResp `json:"courierStatus"`
	TrackingNumber string                 `json:"trackingNumber"`
	TrackingURL    string                 `json:"trackingUrl"`
	InvoiceNumber  string                 `json:"invoiceNumber"`
	InvoiceValue   int                    `json:"invoiceValue"`
	InvoiceURL     string                 `json:"invoiceUrl"`
	InvoiceKey     string                 `json:"invoiceKey"`
	Courier        string                 `json:"courier"`
	Type           string                 `json:"type"`
}

type VTEXPackageItemResp struct {
	ItemIndex   int    `json:"itemIndex"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type VTEXCourierStatusResp struct {
	Status   string                  `json:"status"`
	Finished bool                    `json:"finished"`
	Data     []VTEXCourierEventResp  `json:"data"`
}

type VTEXCourierEventResp struct {
	Description string `json:"description"`
	Date        string `json:"date"`
	City        string `json:"city"`
	State       string `json:"state"`
}

// VTEXOrderListAPIResponse es la respuesta JSON de GET /api/oms/pvt/orders (lista).
type VTEXOrderListAPIResponse struct {
	List   []VTEXOrderSummaryResp `json:"list"`
	Paging VTEXPagingResp         `json:"paging"`
}

type VTEXOrderSummaryResp struct {
	OrderID      string `json:"orderId"`
	Sequence     string `json:"sequence"`
	Status       string `json:"status"`
	CreationDate string `json:"creationDate"`
	LastChange   string `json:"lastChange"`
	TotalValue   int    `json:"totalValue"`
	CurrencyCode string `json:"currencyCode"`
	Origin       string `json:"origin"`
}

type VTEXPagingResp struct {
	Total       int `json:"total"`
	Pages       int `json:"pages"`
	CurrentPage int `json:"currentPage"`
	PerPage     int `json:"perPage"`
}

// VTEXWebhookBody es el body JSON que VTEX envía al webhook (Hook v1).
type VTEXWebhookBody struct {
	Domain        string              `json:"Domain"`
	OrderID       string              `json:"OrderId"`
	State         string              `json:"State"`
	LastState     string              `json:"LastState"`
	LastChange    string              `json:"LastChange"`
	CurrentChange string              `json:"CurrentChange"`
	Origin        *VTEXOriginResp     `json:"Origin"`
}

type VTEXOriginResp struct {
	Account string `json:"Account"`
	Key     string `json:"Key"`
}

// ==============================================================
// ToDomain() conversions
// ==============================================================

func (r *VTEXOrderDetailResponse) ToDomain() domain.VTEXOrder {
	order := domain.VTEXOrder{
		OrderID:            r.OrderID,
		Sequence:           r.Sequence,
		MarketplaceOrderID: r.MarketplaceOrderID,
		Status:             r.Status,
		StatusDescription:  r.StatusDescription,
		Value:              r.Value,
		TotalItems:         r.TotalItems,
		TotalDiscount:      r.TotalDiscount,
		TotalFreight:       r.TotalFreight,
		CreationDate:       parseVTEXDate(r.CreationDate),
		LastChange:         parseVTEXDate(r.LastChange),
	}

	// Items
	order.Items = make([]domain.VTEXOrderItem, len(r.Items))
	for i, item := range r.Items {
		order.Items[i] = domain.VTEXOrderItem{
			UniqueID:        item.UniqueID,
			ID:              item.ID,
			ProductID:       item.ProductID,
			EANID:           item.EANID,
			RefID:           item.RefID,
			Name:            item.Name,
			SKUName:         item.SKUName,
			ImageURL:        item.ImageURL,
			DetailURL:       item.DetailURL,
			Quantity:        item.Quantity,
			Price:           item.Price,
			ListPrice:       item.ListPrice,
			SellingPrice:    item.SellingPrice,
			Tax:             item.Tax,
			MeasurementUnit: item.MeasurementUnit,
			UnitMultiplier:  item.UnitMultiplier,
		}
	}

	// ShippingData
	if r.ShippingData != nil {
		sd := &domain.VTEXShippingData{
			SelectedSLA: r.ShippingData.SelectedSLA,
		}
		if r.ShippingData.Address != nil {
			a := r.ShippingData.Address
			sd.Address = &domain.VTEXAddress{
				AddressType:    a.AddressType,
				ReceiverName:   a.ReceiverName,
				Street:         a.Street,
				Number:         a.Number,
				Complement:     a.Complement,
				Neighborhood:   a.Neighborhood,
				City:           a.City,
				State:          a.State,
				Country:        a.Country,
				PostalCode:     a.PostalCode,
				Reference:      a.Reference,
				GeoCoordinates: a.GeoCoordinates,
			}
		}
		sd.LogisticsInfo = make([]domain.VTEXLogisticsInfo, len(r.ShippingData.LogisticsInfo))
		for i, li := range r.ShippingData.LogisticsInfo {
			info := domain.VTEXLogisticsInfo{
				ItemIndex:            li.ItemIndex,
				SelectedSLA:          li.SelectedSLA,
				LockTTL:              li.LockTTL,
				Price:                li.Price,
				ListPrice:            li.ListPrice,
				SellingPrice:         li.SellingPrice,
				ShippingEstimate:     li.ShippingEstimate,
				ShippingEstimateDate: parseVTEXDatePtr(li.ShippingEstimateDate),
				DeliveryCompany:      li.DeliveryCompany,
			}
			if li.DeliveryWindow != nil {
				info.DeliveryWindow = &domain.VTEXDeliveryWindow{
					StartDateUTC: parseVTEXDate(li.DeliveryWindow.StartDateUTC),
					EndDateUTC:   parseVTEXDate(li.DeliveryWindow.EndDateUTC),
					Price:        li.DeliveryWindow.Price,
				}
			}
			info.DeliveryIDs = make([]domain.VTEXDeliveryID, len(li.DeliveryIDs))
			for j, did := range li.DeliveryIDs {
				info.DeliveryIDs[j] = domain.VTEXDeliveryID{
					CourierID:   did.CourierID,
					CourierName: did.CourierName,
					DockID:      did.DockID,
					Quantity:    did.Quantity,
					WarehouseID: did.WarehouseID,
				}
			}
			sd.LogisticsInfo[i] = info
		}
		sd.TrackingHints = make([]domain.VTEXTrackingHint, len(r.ShippingData.TrackingHints))
		for i, th := range r.ShippingData.TrackingHints {
			sd.TrackingHints[i] = domain.VTEXTrackingHint{
				CourierName:   th.CourierName,
				TrackingID:    th.TrackingID,
				TrackingURL:   th.TrackingURL,
				TrackingLabel: th.TrackingLabel,
			}
		}
		order.ShippingData = sd
	}

	// PaymentData
	if r.PaymentData != nil {
		pd := &domain.VTEXPaymentData{}
		pd.Transactions = make([]domain.VTEXTransaction, len(r.PaymentData.Transactions))
		for i, tx := range r.PaymentData.Transactions {
			t := domain.VTEXTransaction{
				IsActive:      tx.IsActive,
				TransactionID: tx.TransactionID,
				MerchantName:  tx.MerchantName,
			}
			t.Payments = make([]domain.VTEXPayment, len(tx.Payments))
			for j, p := range tx.Payments {
				t.Payments[j] = domain.VTEXPayment{
					ID:                 p.ID,
					PaymentSystem:      p.PaymentSystem,
					PaymentSystemName:  p.PaymentSystemName,
					Value:              p.Value,
					ReferenceValue:     p.ReferenceValue,
					Group:              p.Group,
					ConnectorResponses: p.ConnectorResponses,
					InstallmentCount:   p.InstallmentCount,
					CardHolder:         p.CardHolder,
					FirstDigits:        p.FirstDigits,
					LastDigits:         p.LastDigits,
					URL:                p.URL,
					TID:                p.TID,
				}
			}
			pd.Transactions[i] = t
		}
		order.PaymentData = pd
	}

	// ClientProfileData
	if r.ClientProfileData != nil {
		order.ClientProfileData = &domain.VTEXClientProfile{
			Email:         r.ClientProfileData.Email,
			FirstName:     r.ClientProfileData.FirstName,
			LastName:      r.ClientProfileData.LastName,
			DocumentType:  r.ClientProfileData.DocumentType,
			Document:      r.ClientProfileData.Document,
			Phone:         r.ClientProfileData.Phone,
			CorporateName: r.ClientProfileData.CorporateName,
			IsCorporate:   r.ClientProfileData.IsCorporate,
		}
	}

	// Sellers
	order.Sellers = make([]domain.VTEXSeller, len(r.Sellers))
	for i, s := range r.Sellers {
		order.Sellers[i] = domain.VTEXSeller{
			ID:          s.ID,
			Name:        s.Name,
			SubSellerID: s.SubSellerID,
		}
	}

	// Totals
	order.Totals = make([]domain.VTEXTotal, len(r.Totals))
	for i, t := range r.Totals {
		order.Totals[i] = domain.VTEXTotal{
			ID:    t.ID,
			Name:  t.Name,
			Value: t.Value,
		}
	}

	// PackageAttachment
	if r.PackageAttachment != nil {
		pa := &domain.VTEXPackageAttachment{}
		pa.Packages = make([]domain.VTEXPackage, len(r.PackageAttachment.Packages))
		for i, pkg := range r.PackageAttachment.Packages {
			p := domain.VTEXPackage{
				TrackingNumber: pkg.TrackingNumber,
				TrackingURL:    pkg.TrackingURL,
				InvoiceNumber:  pkg.InvoiceNumber,
				InvoiceValue:   pkg.InvoiceValue,
				InvoiceURL:     pkg.InvoiceURL,
				InvoiceKey:     pkg.InvoiceKey,
				Courier:        pkg.Courier,
				Type:           pkg.Type,
			}
			p.Items = make([]domain.VTEXPackageItem, len(pkg.Items))
			for j, pi := range pkg.Items {
				p.Items[j] = domain.VTEXPackageItem{
					ItemIndex:   pi.ItemIndex,
					Quantity:    pi.Quantity,
					Price:       pi.Price,
					Description: pi.Description,
				}
			}
			if pkg.CourierStatus != nil {
				cs := &domain.VTEXCourierStatus{
					Status:   pkg.CourierStatus.Status,
					Finished: pkg.CourierStatus.Finished,
				}
				cs.Data = make([]domain.VTEXCourierEvent, len(pkg.CourierStatus.Data))
				for j, ev := range pkg.CourierStatus.Data {
					cs.Data[j] = domain.VTEXCourierEvent{
						Description: ev.Description,
						Date:        ev.Date,
						City:        ev.City,
						State:       ev.State,
					}
				}
				p.CourierStatus = cs
			}
			pa.Packages[i] = p
		}
		order.PackageAttachment = pa
	}

	return order
}

func (r *VTEXOrderListAPIResponse) ToDomain() domain.VTEXOrderListResponse {
	result := domain.VTEXOrderListResponse{
		Paging: domain.VTEXPaging{
			Total:       r.Paging.Total,
			Pages:       r.Paging.Pages,
			CurrentPage: r.Paging.CurrentPage,
			PerPage:     r.Paging.PerPage,
		},
	}
	result.List = make([]domain.VTEXOrderSummary, len(r.List))
	for i, s := range r.List {
		result.List[i] = domain.VTEXOrderSummary{
			OrderID:      s.OrderID,
			Sequence:     s.Sequence,
			Status:       s.Status,
			CreationDate: parseVTEXDate(s.CreationDate),
			LastChange:   parseVTEXDate(s.LastChange),
			TotalValue:   s.TotalValue,
			CurrencyCode: s.CurrencyCode,
			Origin:       s.Origin,
		}
	}
	return result
}

func (r *VTEXWebhookBody) ToDomain() domain.VTEXWebhookPayload {
	payload := domain.VTEXWebhookPayload{
		Domain:        r.Domain,
		OrderID:       r.OrderID,
		State:         r.State,
		LastState:     r.LastState,
		LastChange:    r.LastChange,
		CurrentChange: r.CurrentChange,
	}
	if r.Origin != nil {
		payload.Origin = &domain.VTEXWebhookOrigin{
			Account: r.Origin.Account,
			Key:     r.Origin.Key,
		}
	}
	return payload
}

// parseVTEXDate parsea un string de fecha de la API de VTEX.
// VTEX usa formato ISO 8601: "2026-02-24T10:30:00.0000000+00:00".
func parseVTEXDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02T15:04:05-07:00", s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return t
	}
	t, err = time.Parse("2006-01-02T15:04:05.0000000+00:00", s)
	if err == nil {
		return t
	}
	return time.Time{}
}

func parseVTEXDatePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t := parseVTEXDate(*s)
	if t.IsZero() {
		return nil
	}
	return &t
}
