package usecases

import (
	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/log"
	"github.com/secamc93/probability/back/testing/shared/storage"
)

// APISimulator simula el API de EnvioClick
type APISimulator struct {
	logger     log.ILogger
	Repository *domain.ShipmentRepository // Exportado para acceso desde bundle
	s3         storage.IS3Service         // nil si S3 no está configurado
	urlBase    string                     // URL_BASE_DOMAIN_S3 para construir URLs públicas
}

// NewAPISimulator crea una nueva instancia del simulador de API
// s3 puede ser nil si no se desea subir PDFs a S3
func NewAPISimulator(logger log.ILogger, s3 storage.IS3Service, urlBase string) *APISimulator {
	return &APISimulator{
		logger:     logger,
		Repository: domain.NewShipmentRepository(),
		s3:         s3,
		urlBase:    urlBase,
	}
}
