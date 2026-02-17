package domain

import (
	"fmt"
	"sync"
)

// InvoiceRepository almacena las facturas simuladas en memoria
type InvoiceRepository struct {
	mu          sync.RWMutex
	invoices    map[string]*Invoice
	creditNotes map[string]*CreditNote
	customers   map[string]*Customer
	tokens      map[string]*AuthToken
	invoiceSeq  int
	creditSeq   int
}

// NewInvoiceRepository crea una nueva instancia del repositorio
func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{
		invoices:    make(map[string]*Invoice),
		creditNotes: make(map[string]*CreditNote),
		customers:   make(map[string]*Customer),
		tokens:      make(map[string]*AuthToken),
		invoiceSeq:  1000,
		creditSeq:   2000,
	}
}

// SaveInvoice guarda una factura
func (r *InvoiceRepository) SaveInvoice(invoice *Invoice) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoices[invoice.ID] = invoice
}

// GetInvoice obtiene una factura por ID
func (r *InvoiceRepository) GetInvoice(id string) (*Invoice, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	invoice, exists := r.invoices[id]
	return invoice, exists
}

// GetInvoiceByNumber obtiene una factura por su número de factura
func (r *InvoiceRepository) GetInvoiceByNumber(invoiceNumber string) (*Invoice, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, invoice := range r.invoices {
		if invoice.InvoiceNumber == invoiceNumber {
			return invoice, true
		}
	}
	return nil, false
}

// GetAllInvoices retorna todas las facturas
func (r *InvoiceRepository) GetAllInvoices() []*Invoice {
	r.mu.RLock()
	defer r.mu.RUnlock()
	invoices := make([]*Invoice, 0, len(r.invoices))
	for _, invoice := range r.invoices {
		invoices = append(invoices, invoice)
	}
	return invoices
}

// SaveCreditNote guarda una nota de crédito
func (r *InvoiceRepository) SaveCreditNote(note *CreditNote) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creditNotes[note.ID] = note
}

// GetCreditNote obtiene una nota de crédito por ID
func (r *InvoiceRepository) GetCreditNote(id string) (*CreditNote, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	note, exists := r.creditNotes[id]
	return note, exists
}

// GetAllCreditNotes retorna todas las notas de crédito
func (r *InvoiceRepository) GetAllCreditNotes() []*CreditNote {
	r.mu.RLock()
	defer r.mu.RUnlock()
	notes := make([]*CreditNote, 0, len(r.creditNotes))
	for _, note := range r.creditNotes {
		notes = append(notes, note)
	}
	return notes
}

// SaveCustomer guarda un cliente
func (r *InvoiceRepository) SaveCustomer(customer *Customer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.customers[customer.Identification] = customer
}

// GetCustomer obtiene un cliente por identificación
func (r *InvoiceRepository) GetCustomer(identification string) (*Customer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, exists := r.customers[identification]
	return c, exists
}

// SaveToken guarda un token de autenticación
func (r *InvoiceRepository) SaveToken(token *AuthToken) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.Token] = token
}

// GetToken obtiene un token
func (r *InvoiceRepository) GetToken(token string) (*AuthToken, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, exists := r.tokens[token]
	return t, exists
}

// GenerateInvoiceNumber genera un número de factura secuencial
func (r *InvoiceRepository) GenerateInvoiceNumber() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoiceSeq++
	return fmt.Sprintf("SPY-%04d", r.invoiceSeq)
}

// GenerateCreditNoteNumber genera un número de nota de crédito secuencial
func (r *InvoiceRepository) GenerateCreditNoteNumber() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.creditSeq++
	return fmt.Sprintf("NC-%04d", r.creditSeq)
}
