package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

// MeliOrderResponse es la respuesta JSON de GET /orders/{id} de MercadoLibre.
type MeliOrderResponse struct {
	ID           int64                   `json:"id"`
	Status       string                  `json:"status"`
	StatusDetail *MeliStatusDetailResp   `json:"status_detail"`
	DateCreated  string                  `json:"date_created"`
	DateClosed   *string                 `json:"date_closed"`
	LastUpdated  string                  `json:"last_updated"`
	TotalAmount  float64                 `json:"total_amount"`
	CurrencyID   string                  `json:"currency_id"`
	Buyer        MeliBuyerResponse       `json:"buyer"`
	Seller       MeliSellerResponse      `json:"seller"`
	OrderItems   []MeliOrderItemResponse `json:"order_items"`
	Payments     []MeliPaymentResponse   `json:"payments"`
	Shipping     *MeliShippingRefResp    `json:"shipping"`
	Tags         []string                `json:"tags"`
	PackID       *int64                  `json:"pack_id"`
	CouponAmount float64                 `json:"coupon_amount"`
	CouponID     *string                 `json:"coupon_id"`
}

type MeliStatusDetailResp struct {
	Description string `json:"description"`
	Code        string `json:"code"`
}

type MeliBuyerResponse struct {
	ID          int64                    `json:"id"`
	Nickname    string                   `json:"nickname"`
	FirstName   string                   `json:"first_name"`
	LastName    string                   `json:"last_name"`
	Email       string                   `json:"email"`
	Phone       *MeliPhoneResponse       `json:"phone"`
	BillingInfo *MeliBillingInfoResponse `json:"billing_info"`
}

type MeliPhoneResponse struct {
	AreaCode  string `json:"area_code"`
	Number    string `json:"number"`
	Extension string `json:"extension"`
}

type MeliBillingInfoResponse struct {
	DocType   string `json:"doc_type"`
	DocNumber string `json:"doc_number"`
}

type MeliSellerResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type MeliOrderItemResponse struct {
	Item          MeliItemResponse `json:"item"`
	Quantity      int              `json:"quantity"`
	UnitPrice     float64          `json:"unit_price"`
	FullUnitPrice float64          `json:"full_unit_price"`
	CurrencyID    string           `json:"currency_id"`
	SaleFee       float64          `json:"sale_fee"`
}

type MeliItemResponse struct {
	ID                  string                        `json:"id"`
	Title               string                        `json:"title"`
	CategoryID          string                        `json:"category_id"`
	VariationID         *int64                        `json:"variation_id"`
	SellerCustomField   *string                       `json:"seller_custom_field"`
	SellerSKU           *string                       `json:"seller_sku"`
	Condition           string                        `json:"condition"`
	VariationAttributes []MeliVariationAttrResponse   `json:"variation_attributes"`
}

type MeliVariationAttrResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ValueName string `json:"value_name"`
}

type MeliPaymentResponse struct {
	ID                 int64    `json:"id"`
	OrderID            int64    `json:"order_id"`
	PayerID            int64    `json:"payer_id"`
	Status             string   `json:"status"`
	StatusDetail       string   `json:"status_detail"`
	TransactionAmount  float64  `json:"transaction_amount"`
	CurrencyID         string   `json:"currency_id"`
	DateCreated        string   `json:"date_created"`
	DateApproved       *string  `json:"date_approved"`
	DateLastModified   *string  `json:"date_last_modified"`
	PaymentMethodID    string   `json:"payment_method_id"`
	PaymentType        string   `json:"payment_type"`
	OperationType      string   `json:"operation_type"`
	InstallmentAmount  *float64 `json:"installment_amount"`
	Installments       int      `json:"installments"`
	TransactionOrderID *string  `json:"transaction_order_id"`
}

// MeliShippingRefResp es la referencia al envío dentro de la orden (solo ID).
type MeliShippingRefResp struct {
	ID int64 `json:"id"`
}

// MeliOrdersSearchResponse es la respuesta de GET /orders/search.
type MeliOrdersSearchResponse struct {
	Query   string              `json:"query"`
	Results []MeliOrderResponse `json:"results"`
	Paging  MeliPagingResponse  `json:"paging"`
}

type MeliPagingResponse struct {
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// MeliNotificationBody es el body del POST que MeLi envía como notificación IPN.
type MeliNotificationBody struct {
	Resource      string `json:"resource"`
	UserID        int64  `json:"user_id"`
	Topic         string `json:"topic"`
	ApplicationID int64  `json:"application_id"`
	Attempts      int    `json:"attempts"`
	Sent          string `json:"sent"`
	Received      string `json:"received"`
}

// MeliTokenResponse es la respuesta de POST /oauth/token.
type MeliTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// MeliUserResponse es la respuesta de GET /users/me.
type MeliUserResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

// ==============================================================
// ToDomain() conversions
// ==============================================================

func (r *MeliOrderResponse) ToDomain() domain.MeliOrder {
	order := domain.MeliOrder{
		ID:           r.ID,
		Status:       r.Status,
		DateCreated:  parseMeliDate(r.DateCreated),
		DateClosed:   parseMeliDatePtr(r.DateClosed),
		LastUpdated:  parseMeliDate(r.LastUpdated),
		TotalAmount:  r.TotalAmount,
		CurrencyID:   r.CurrencyID,
		Tags:         r.Tags,
		PackID:       r.PackID,
		CouponAmount: r.CouponAmount,
		CouponID:     r.CouponID,
	}

	// Status detail
	if r.StatusDetail != nil {
		order.StatusDetail = &domain.MeliStatusDetail{
			Description: r.StatusDetail.Description,
			Code:        r.StatusDetail.Code,
		}
	}

	// Buyer
	order.Buyer = domain.MeliBuyer{
		ID:        r.Buyer.ID,
		Nickname:  r.Buyer.Nickname,
		FirstName: r.Buyer.FirstName,
		LastName:  r.Buyer.LastName,
		Email:     r.Buyer.Email,
	}
	if r.Buyer.Phone != nil {
		order.Buyer.Phone = domain.MeliPhone{
			AreaCode:  r.Buyer.Phone.AreaCode,
			Number:    r.Buyer.Phone.Number,
			Extension: r.Buyer.Phone.Extension,
		}
	}
	if r.Buyer.BillingInfo != nil {
		order.Buyer.BillingInfo = &domain.MeliBillingInfo{
			DocType:   r.Buyer.BillingInfo.DocType,
			DocNumber: r.Buyer.BillingInfo.DocNumber,
		}
	}

	// Seller
	order.Seller = domain.MeliSeller{
		ID:       r.Seller.ID,
		Nickname: r.Seller.Nickname,
	}

	// Order items
	order.OrderItems = make([]domain.MeliOrderItem, len(r.OrderItems))
	for i, item := range r.OrderItems {
		meliItem := domain.MeliItem{
			ID:                item.Item.ID,
			Title:             item.Item.Title,
			CategoryID:        item.Item.CategoryID,
			VariationID:       item.Item.VariationID,
			SellerCustomField: item.Item.SellerCustomField,
			SellerSKU:         item.Item.SellerSKU,
			Condition:         item.Item.Condition,
		}
		meliItem.VariationAttributes = make([]domain.MeliVariationAttribute, len(item.Item.VariationAttributes))
		for j, va := range item.Item.VariationAttributes {
			meliItem.VariationAttributes[j] = domain.MeliVariationAttribute{
				ID:    va.ID,
				Name:  va.Name,
				Value: va.ValueName,
			}
		}
		order.OrderItems[i] = domain.MeliOrderItem{
			Item:          meliItem,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			FullUnitPrice: item.FullUnitPrice,
			Currency:      item.CurrencyID,
			SaleFee:       item.SaleFee,
		}
	}

	// Payments
	order.Payments = make([]domain.MeliPayment, len(r.Payments))
	for i, p := range r.Payments {
		order.Payments[i] = domain.MeliPayment{
			ID:                 p.ID,
			OrderID:            p.OrderID,
			PayerID:            p.PayerID,
			Status:             p.Status,
			StatusDetail:       p.StatusDetail,
			TransactionAmount:  p.TransactionAmount,
			CurrencyID:         p.CurrencyID,
			DateCreated:        parseMeliDate(p.DateCreated),
			DateApproved:       parseMeliDatePtr(p.DateApproved),
			DateLastModified:   parseMeliDatePtr(p.DateLastModified),
			PaymentMethodID:    p.PaymentMethodID,
			PaymentType:        p.PaymentType,
			OperationType:      p.OperationType,
			InstallmentAmount:  p.InstallmentAmount,
			Installments:       p.Installments,
			TransactionOrderID: p.TransactionOrderID,
		}
	}

	// Shipping ref
	if r.Shipping != nil {
		order.Shipping = &domain.MeliShippingRef{
			ID: r.Shipping.ID,
		}
	}

	return order
}

func (r *MeliTokenResponse) ToDomain() domain.TokenResponse {
	return domain.TokenResponse{
		AccessToken:  r.AccessToken,
		TokenType:    r.TokenType,
		ExpiresIn:    r.ExpiresIn,
		Scope:        r.Scope,
		UserID:       r.UserID,
		RefreshToken: r.RefreshToken,
	}
}

func (r *MeliUserResponse) ToDomain() domain.MeliSeller {
	return domain.MeliSeller{
		ID:       r.ID,
		Nickname: r.Nickname,
	}
}

func (r *MeliNotificationBody) ToDomain() domain.MeliNotification {
	return domain.MeliNotification{
		Resource:      r.Resource,
		UserID:        r.UserID,
		Topic:         r.Topic,
		ApplicationID: r.ApplicationID,
		Attempts:      r.Attempts,
		Sent:          parseMeliDate(r.Sent),
		Received:      parseMeliDate(r.Received),
	}
}

// parseMeliDate parsea un string de fecha de la API de MercadoLibre.
// MeLi usa formato RFC3339 con timezone: "2024-01-15T10:30:00.000-04:00".
func parseMeliDate(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		return t
	}
	// Fallback: intentar sin milliseconds
	t, err = time.Parse("2006-01-02T15:04:05-07:00", s)
	if err == nil {
		return t
	}
	// Fallback: sin timezone
	t, err = time.Parse("2006-01-02T15:04:05", s)
	if err == nil {
		return t
	}
	return time.Time{}
}

func parseMeliDatePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t := parseMeliDate(*s)
	if t.IsZero() {
		return nil
	}
	return &t
}
