package domain

import (
	"fmt"
	"sync"
)

type Repository struct {
	mu           sync.RWMutex
	customers    map[string]*Customer
	invoices     map[string]*Invoice
	journals     map[string]*JournalEntry
	tokens       map[string]*AuthToken
	vouchers     map[string]*Voucher
	creditNotes  map[string]*CreditNote
	products     []*Product
	paymentTypes []*PaymentType
	invoiceSeq   int
	voucherSeq   int
	creditSeq    int
}

func NewRepository() *Repository {
	r := &Repository{
		customers:  make(map[string]*Customer),
		invoices:   make(map[string]*Invoice),
		journals:   make(map[string]*JournalEntry),
		tokens:      make(map[string]*AuthToken),
		vouchers:    make(map[string]*Voucher),
		creditNotes: make(map[string]*CreditNote),
		invoiceSeq:  1000,
		voucherSeq:  500,
		creditSeq:   700,
	}
	r.seedCatalogs()
	return r
}

func (r *Repository) seedCatalogs() {
	r.products = []*Product{
		{ID: "prod-1", Code: "ITEM-1", Name: "Producto demo 1", Description: "Producto de prueba 1", Price: 50000},
		{ID: "prod-2", Code: "ITEM-2", Name: "Producto demo 2", Description: "Producto de prueba 2", Price: 75000},
		{ID: "prod-3", Code: "SERV-1", Name: "Servicio demo", Description: "Servicio de prueba", Price: 120000},
	}
	r.paymentTypes = []*PaymentType{
		{ID: 5636, Name: "Efectivo", Type: "Cash"},
		{ID: 5637, Name: "Transferencia", Type: "Transfer"},
		{ID: 5638, Name: "Tarjeta de credito", Type: "CreditCard"},
	}
}

func (r *Repository) SaveCustomer(c *Customer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.customers[c.ID] = c
}

func (r *Repository) GetCustomerByIdentification(id string) (*Customer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.customers {
		if c.Identification == id {
			return c, true
		}
	}
	return nil, false
}

func (r *Repository) SaveInvoice(inv *Invoice) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoices[inv.ID] = inv
}

func (r *Repository) GetInvoice(id string) (*Invoice, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inv, ok := r.invoices[id]
	return inv, ok
}

func (r *Repository) GetInvoiceByName(name string) (*Invoice, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, inv := range r.invoices {
		if inv.Name == name {
			return inv, true
		}
	}
	return nil, false
}

func (r *Repository) ListInvoices() []*Invoice {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*Invoice, 0, len(r.invoices))
	for _, inv := range r.invoices {
		out = append(out, inv)
	}
	return out
}

func (r *Repository) AnnulInvoice(id string) (*Invoice, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	inv, ok := r.invoices[id]
	if !ok {
		return nil, false
	}
	inv.Annulled = true
	inv.Status = "annulled"
	inv.Balance = 0
	return inv, true
}

func (r *Repository) SaveJournal(j *JournalEntry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.journals[j.ID] = j
}

func (r *Repository) SaveToken(t *AuthToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[t.Token] = t
}

func (r *Repository) SaveVoucher(v *Voucher) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.vouchers[v.ID] = v
	if inv, ok := r.invoices[v.InvoiceNumber]; ok {
		inv.Balance -= v.Value
		if inv.Balance < 0 {
			inv.Balance = 0
		}
	}
	for _, inv := range r.invoices {
		if inv.Name == v.InvoiceNumber {
			inv.Balance -= v.Value
			if inv.Balance < 0 {
				inv.Balance = 0
			}
		}
	}
}

func (r *Repository) SaveCreditNote(n *CreditNote) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creditNotes[n.ID] = n
}

func (r *Repository) NextCreditNoteNumber() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creditSeq++
	return r.creditSeq
}

func (r *Repository) ListProducts() []*Product {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.products
}

func (r *Repository) ListPaymentTypes() []*PaymentType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.paymentTypes
}

func (r *Repository) NextInvoiceNumber() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoiceSeq++
	return r.invoiceSeq
}

func (r *Repository) NextVoucherNumber() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.voucherSeq++
	return r.voucherSeq
}

func (r *Repository) NextInvoiceName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, r.NextInvoiceNumber())
}
