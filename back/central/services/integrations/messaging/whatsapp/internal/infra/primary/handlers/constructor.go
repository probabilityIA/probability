package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasemessaging"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/app/usecasetemplates"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type IHandler interface {
	RegisterRoutes(router *gin.RouterGroup)

	SendTemplate(c *gin.Context)

	SendManualReply(c *gin.Context)

	PauseAI(c *gin.Context)
	ResumeAI(c *gin.Context)

	VerifyWebhook(c *gin.Context)
	ReceiveWebhook(c *gin.Context)

	GetTemplatesStatus(c *gin.Context)
	ProvisionTemplates(c *gin.Context)

	SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter)
}

type handler struct {
	useCase             usecasemessaging.IUseCase
	templatesUseCase    usecasetemplates.IUseCase
	log                 log.ILogger
	config              env.IConfig
	platformCredsGetter ports.IPlatformCredentialsGetter
	rabbit              rabbitmq.IQueue
}

func New(
	useCase usecasemessaging.IUseCase,
	templatesUseCase usecasetemplates.IUseCase,
	logger log.ILogger,
	config env.IConfig,
	rabbit rabbitmq.IQueue,
) IHandler {
	return &handler{
		useCase:          useCase,
		templatesUseCase: templatesUseCase,
		log:              logger,
		config:           config,
		rabbit:           rabbit,
	}
}

func (h *handler) SetPlatformCredsGetter(getter ports.IPlatformCredentialsGetter) {
	h.platformCredsGetter = getter
}
