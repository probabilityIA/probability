package response

// CreatedBillDoc tipo de documento de la factura creada
type CreatedBillDoc struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreatedBillPaymentForm forma de pago de la factura creada
type CreatedBillPaymentForm struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreatedBill datos de la factura retornada por POST /v1/bills/validate
type CreatedBill struct {
	ID          int                    `json:"id"`
	Number      string                 `json:"number"`
	CUFE        string                 `json:"cufe"`
	QR          string                 `json:"qr"`
	QRImage     string                 `json:"qr_image"`
	Total       string                 `json:"total"`
	Status      int                    `json:"status"`
	Validated   string                 `json:"validated"`
	CreatedAt   string                 `json:"created_at"`
	Document    CreatedBillDoc         `json:"document"`
	PaymentForm CreatedBillPaymentForm `json:"payment_form"`
}

// CreateBillData data anidada en la respuesta
type CreateBillData struct {
	Bill CreatedBill `json:"bill"`
}

// CreateBill respuesta completa de POST /v1/bills/validate
type CreateBill struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    CreateBillData `json:"data"`
}
