package usecases

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/testing/integrations/siigo/internal/domain"
)

func (s *APISimulator) HandleAuth(username, accessKey, partnerID string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username is required")
	}
	if accessKey == "" {
		return "", fmt.Errorf("access_key is required")
	}

	token := fmt.Sprintf("siigo_%s", uuid.New().String())
	expiresAt := time.Now().Add(24 * time.Hour)

	s.Repository.SaveToken(&domain.AuthToken{
		Token:     token,
		ExpiresAt: expiresAt,
		Username:  username,
	})

	s.logger.Info().
		Str("username", username).
		Str("partner_id", partnerID).
		Msg("Siigo mock: token issued")

	return token, nil
}

func (s *APISimulator) HandleCreateCustomer(body map[string]interface{}) (*domain.Customer, error) {
	identification := ""
	if v, ok := body["identification"].(string); ok {
		identification = v
	}
	if identification == "" {
		if idObj, ok := body["identification"].(map[string]interface{}); ok {
			if num, ok := idObj["number"].(string); ok {
				identification = num
			}
		}
	}

	name := ""
	if names, ok := body["name"].([]interface{}); ok && len(names) > 0 {
		for i, n := range names {
			if i > 0 {
				name += " "
			}
			if str, ok := n.(string); ok {
				name += str
			}
		}
	} else if n, ok := body["name"].(string); ok {
		name = n
	}

	email := ""
	if e, ok := body["email"].(string); ok {
		email = e
	}
	phone := ""
	if contacts, ok := body["contacts"].([]interface{}); ok && len(contacts) > 0 {
		if first, ok := contacts[0].(map[string]interface{}); ok {
			if p, ok := first["phone"].(map[string]interface{}); ok {
				if num, ok := p["number"].(string); ok {
					phone = num
				}
			}
		}
	}

	customer := &domain.Customer{
		ID:             uuid.New().String(),
		Identification: identification,
		Name:           name,
		Email:          email,
		Phone:          phone,
		CreatedAt:      time.Now(),
	}

	s.Repository.SaveCustomer(customer)

	s.logger.Info().
		Str("identification", identification).
		Str("name", name).
		Msg("Siigo mock: customer created")

	return customer, nil
}

func (s *APISimulator) HandleGetCustomer(identification string) (*domain.Customer, bool) {
	return s.Repository.GetCustomerByIdentification(identification)
}

func (s *APISimulator) HandleCreateInvoice(body map[string]interface{}) (*domain.Invoice, error) {
	prefix := "FE"
	if doc, ok := body["document"].(map[string]interface{}); ok {
		if id, ok := doc["id"].(float64); ok {
			prefix = fmt.Sprintf("FE-%d", int(id))
		}
	}

	customerID := ""
	customerNIT := ""
	if cust, ok := body["customer"].(map[string]interface{}); ok {
		if v, ok := cust["identification"].(string); ok {
			customerNIT = v
		}
		if v, ok := cust["id"].(string); ok {
			customerID = v
		}
	}

	items := []domain.InvoiceItem{}
	total := 0.0
	if rawItems, ok := body["items"].([]interface{}); ok {
		for _, it := range rawItems {
			itMap, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			code, _ := itMap["code"].(string)
			desc, _ := itMap["description"].(string)
			qty := 1.0
			if v, ok := itMap["quantity"].(float64); ok {
				qty = v
			}
			price := 0.0
			if v, ok := itMap["price"].(float64); ok {
				price = v
			}
			lineTotal := qty * price
			total += lineTotal
			items = append(items, domain.InvoiceItem{
				Code:        code,
				Description: desc,
				Quantity:    qty,
				Price:       price,
				Total:       lineTotal,
			})
		}
	}

	number := s.Repository.NextInvoiceNumber()
	invoice := &domain.Invoice{
		ID:          uuid.New().String(),
		Prefix:      prefix,
		Number:      number,
		Name:        fmt.Sprintf("%s-%d", prefix, number),
		Date:        time.Now().Format("2006-01-02"),
		CustomerID:  customerID,
		CustomerNIT: customerNIT,
		Items:       items,
		Total:       total,
		Balance:     total,
		StampStatus: "Stamped",
		Status:      "active",
		CUFE:        randomCUFE(),
		CreatedAt:   time.Now(),
	}
	invoice.PublicURL = fmt.Sprintf("https://siigo-mock.local/invoices/%s", invoice.ID)

	s.Repository.SaveInvoice(invoice)

	s.logger.Info().
		Str("invoice_id", invoice.ID).
		Str("name", invoice.Name).
		Float64("total", total).
		Msg("Siigo mock: invoice created")

	return invoice, nil
}

func (s *APISimulator) HandleGetInvoice(id string) (*domain.Invoice, bool) {
	return s.Repository.GetInvoice(id)
}

func (s *APISimulator) HandleListInvoices() []*domain.Invoice {
	return s.Repository.ListInvoices()
}

func (s *APISimulator) HandleAnnulInvoice(id string) (*domain.Invoice, error) {
	inv, ok := s.Repository.GetInvoice(id)
	if !ok {
		return nil, fmt.Errorf("not_found")
	}
	if inv.Annulled {
		return nil, fmt.Errorf("annul_not_allowed")
	}
	annulled, _ := s.Repository.AnnulInvoice(id)
	s.logger.Info().
		Str("invoice_id", id).
		Str("name", annulled.Name).
		Msg("Siigo mock: invoice annulled")
	return annulled, nil
}

func (s *APISimulator) HandleGetStampErrors(id string) ([]domain.StampError, bool) {
	inv, ok := s.Repository.GetInvoice(id)
	if !ok {
		return nil, false
	}
	return inv.StampErrors, true
}

func (s *APISimulator) HandleListProducts() []*domain.Product {
	return s.Repository.ListProducts()
}

func (s *APISimulator) HandleListPaymentTypes() []*domain.PaymentType {
	return s.Repository.ListPaymentTypes()
}

func (s *APISimulator) HandleCreateVoucher(body map[string]interface{}) (*domain.Voucher, error) {
	invoiceRef := ""
	value := 0.0
	if items, ok := body["items"].([]interface{}); ok {
		for _, it := range items {
			itMap, ok := it.(map[string]interface{})
			if !ok {
				continue
			}
			if due, ok := itMap["due"].(map[string]interface{}); ok {
				prefix, _ := due["prefix"].(string)
				if cons, ok := due["consecutive"].(float64); ok {
					invoiceRef = fmt.Sprintf("%s-%d", prefix, int(cons))
				}
			}
			if v, ok := itMap["value"].(float64); ok {
				value += v
			}
		}
	}

	number := s.Repository.NextVoucherNumber()
	voucher := &domain.Voucher{
		ID:            uuid.New().String(),
		Name:          fmt.Sprintf("RC-%d", number),
		Number:        number,
		InvoiceNumber: invoiceRef,
		Value:         value,
		Date:          time.Now().Format("2006-01-02"),
		CreatedAt:     time.Now(),
	}
	s.Repository.SaveVoucher(voucher)

	s.logger.Info().
		Str("voucher_id", voucher.ID).
		Str("invoice_ref", invoiceRef).
		Float64("value", value).
		Msg("Siigo mock: voucher (cash receipt) created")

	return voucher, nil
}

func (s *APISimulator) HandleCreateJournal(body map[string]interface{}) (*domain.JournalEntry, error) {
	docID := ""
	if doc, ok := body["document"].(map[string]interface{}); ok {
		if v, ok := doc["id"].(string); ok {
			docID = v
		}
	}

	itemsRaw, _ := body["items"].([]interface{})
	items := make([]map[string]interface{}, 0, len(itemsRaw))
	total := 0.0
	for _, it := range itemsRaw {
		if m, ok := it.(map[string]interface{}); ok {
			items = append(items, m)
			if v, ok := m["value"].(float64); ok {
				total += v
			}
		}
	}

	date := time.Now().Format("2006-01-02")
	if d, ok := body["date"].(string); ok && d != "" {
		date = d
	}

	journal := &domain.JournalEntry{
		ID:         uuid.New().String(),
		Date:       date,
		DocumentID: docID,
		Items:      items,
		Total:      total,
		CreatedAt:  time.Now(),
	}

	s.Repository.SaveJournal(journal)

	s.logger.Info().
		Str("journal_id", journal.ID).
		Float64("total", total).
		Msg("Siigo mock: journal created")

	return journal, nil
}

func randomCUFE() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
