package invoicing

import (
	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/router"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes"
	worldoffice "github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

// New inicializa todos los proveedores de facturación electrónica y el router de colas.
// Registra cada provider en integrationCore bajo su type_id correspondiente.
// Debe llamarse después de inicializar integrationCore y antes de que el servidor empiece a recibir tráfico.
func New(
	config env.IConfig,
	logger log.ILogger,
	rabbitMQ rabbitmq.IQueue,
	integrationCore core.IIntegrationCore,
) {
	// Softpymes (type_id=5)
	softpymesBundle := softpymes.New(config, logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeInvoicing, softpymesBundle)

	// Factus (type_id=7)
	factusBundle := factus.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeFactus, factusBundle)

	// Siigo (type_id=8)
	siigoBundle := siigo.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeSiigo, siigoBundle)

	// Alegra (type_id=9)
	alegraBundle := alegra.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeAlegra, alegraBundle)

	// World Office (type_id=10)
	worldOfficeBundle := worldoffice.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeWorldOffice, worldOfficeBundle)

	// Helisa (type_id=11)
	helisaBundle := helisa.New(logger, rabbitMQ, integrationCore)
	integrationCore.RegisterIntegration(core.IntegrationTypeHelisa, helisaBundle)

	// Router: consume invoicing.requests y enruta al proveedor correcto.
	// Se inicializa al final para que las colas de proveedores ya estén declaradas.
	router.New(logger, rabbitMQ)
}
