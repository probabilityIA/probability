package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/app"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/jwt"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/storage"
)

type IHandlers interface {
	GetBusinessPage(c *gin.Context)
	ListCatalog(c *gin.Context)
	GetProduct(c *gin.Context)
	SubmitContact(c *gin.Context)
	CreateCheckoutSession(c *gin.Context)
	GetCheckoutStatus(c *gin.Context)
	GetSession(c *gin.Context)
	GetMyPrices(c *gin.Context)
	GetMyOrders(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UploadAvatar(c *gin.Context)
	ListAddresses(c *gin.Context)
	AddAddress(c *gin.Context)
	SetPrimaryAddress(c *gin.Context)
	DeleteAddress(c *gin.Context)
	GetCategories(c *gin.Context)
	RegisterRoutes(router *gin.RouterGroup)
}

type Handlers struct {
	uc     app.IUseCase
	logger log.ILogger
	env    env.IConfig
	jwt    jwt.IJWTService
	s3     storage.IS3Service
}

func New(uc app.IUseCase, logger log.ILogger, environment env.IConfig, s3 storage.IS3Service) IHandlers {
	return &Handlers{uc: uc, logger: logger, env: environment, jwt: jwt.New(environment.Get("JWT_SECRET")), s3: s3}
}

func (h *Handlers) optionalUserID(c *gin.Context) uint {
	header := c.GetHeader("Authorization")
	if header == "" {
		return 0
	}
	token := strings.TrimPrefix(header, "Bearer ")
	claims, err := h.jwt.ValidateToken(token)
	if err != nil || claims == nil {
		return 0
	}
	return claims.UserID
}

func (h *Handlers) getImageURLBase() string {
	return h.env.Get("URL_BASE_DOMAIN_S3")
}
