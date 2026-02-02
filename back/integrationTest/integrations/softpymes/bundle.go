package softpymes

import (
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes/internal/app/usecases"
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes/internal/domain"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

// New inicializa el módulo de SoftPymes para pruebas de integración
func New(logger log.ILogger) *SoftPymesIntegration {
	apiSimulator := usecases.NewAPISimulator(logger)

	return &SoftPymesIntegration{
		apiSimulator: apiSimulator.(*usecases.APISimulator),
		logger:       logger,
	}
}

// SoftPymesIntegration representa el módulo de integración de SoftPymes
type SoftPymesIntegration struct {
	apiSimulator *usecases.APISimulator
	logger       log.ILogger
}

// SimulateAuth simula autenticación
func (s *SoftPymesIntegration) SimulateAuth(apiKey, apiSecret, referer string) (string, error) {
	return s.apiSimulator.HandleAuth(apiKey, apiSecret, referer)
}

// SimulateInvoice simula creación de factura
func (s *SoftPymesIntegration) SimulateInvoice(token string, invoiceData map[string]interface{}) (*domain.Invoice, error) {
	return s.apiSimulator.HandleCreateInvoice(token, invoiceData)
}

// SimulateCreditNote simula creación de nota de crédito
func (s *SoftPymesIntegration) SimulateCreditNote(token string, creditNoteData map[string]interface{}) (*domain.CreditNote, error) {
	return s.apiSimulator.HandleCreateCreditNote(token, creditNoteData)
}

// ListInvoices retorna todas las facturas almacenadas
func (s *SoftPymesIntegration) ListInvoices(token string) ([]domain.Invoice, error) {
	return s.apiSimulator.HandleListDocuments(token, nil)
}

// GetRepository retorna el repositorio (para listar sin token)
func (s *SoftPymesIntegration) GetRepository() *domain.InvoiceRepository {
	return s.apiSimulator.Repository
}
