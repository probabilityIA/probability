package domain

import (
	"fmt"
	"sync"
)

type Repository struct {
	mu         sync.RWMutex
	customers  map[string]*Customer
	invoices   map[string]*Invoice
	journals   map[string]*JournalEntry
	tokens     map[string]*AuthToken
	invoiceSeq int
}

func NewRepository() *Repository {
	return &Repository{
		customers:  make(map[string]*Customer),
		invoices:   make(map[string]*Invoice),
		journals:   make(map[string]*JournalEntry),
		tokens:     make(map[string]*AuthToken),
		invoiceSeq: 1000,
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

func (r *Repository) NextInvoiceNumber() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.invoiceSeq++
	return r.invoiceSeq
}

func (r *Repository) NextInvoiceName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, r.NextInvoiceNumber())
}
