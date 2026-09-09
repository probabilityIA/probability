package usecasenumbers

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

const whatsAppTypeID uint = 2

const (
	StatusSinNumero        = "sin_numero"
	StatusEsperandoCodigo  = "esperando_codigo"
	StatusVerificado       = "verificado"
	StatusRegistrado       = "registrado"
	StatusNombreEnRevision = "nombre_en_revision"
)

type AddNumberInput struct {
	CountryCode  string
	PhoneNumber  string
	VerifiedName string
}

type NumberState struct {
	IntegrationID      uint
	BusinessID         uint
	PhoneNumberID      string
	Status             string
	DisplayPhoneNumber string
	VerifiedName       string
	NameStatus         string
	CodeVerification   string
	QualityRating      string
	HostedByPlatform   bool
	Active             bool
	Pin                string
}

type IUseCase interface {
	GetState(ctx context.Context, businessID uint) (*NumberState, error)
	AddNumber(ctx context.Context, businessID uint, input AddNumberInput) (*NumberState, error)
	RequestCode(ctx context.Context, businessID uint, method string) (*NumberState, error)
	VerifyCode(ctx context.Context, businessID uint, code string) (*NumberState, error)
	Register(ctx context.Context, businessID uint) (*NumberState, error)
	GetProfile(ctx context.Context, businessID uint) (*BusinessProfile, error)
	UpdateProfile(ctx context.Context, businessID uint, input BusinessProfileInput) (*BusinessProfile, error)
	UpdateProfilePhoto(ctx context.Context, businessID uint, contentType string, data []byte) (*BusinessProfile, error)
}

type IUseCaseMutable interface {
	IUseCase
	SetResolver(resolver ports.IPlatformCredentialsGetter)
	SetProfileAPI(factory ProfileAPIFactory)
}

type APIFactory func(baseURL string) ports.IPhoneNumbersAPI

type ProfileAPIFactory func(baseURL string) ports.IBusinessProfileAPI

type usecase struct {
	credentialsCache ports.ICredentialsCache
	apiFactory       APIFactory
	log              log.ILogger
	resolver         ports.IPlatformCredentialsGetter
	profileFactory   ProfileAPIFactory
}

func New(credentialsCache ports.ICredentialsCache, apiFactory APIFactory, logger log.ILogger) IUseCaseMutable {
	return &usecase{
		credentialsCache: credentialsCache,
		apiFactory:       apiFactory,
		log:              logger.WithModule("whatsapp-numbers"),
	}
}

func (u *usecase) SetResolver(resolver ports.IPlatformCredentialsGetter) {
	u.resolver = resolver
}

func (u *usecase) SetProfileAPI(factory ProfileAPIFactory) {
	u.profileFactory = factory
}
