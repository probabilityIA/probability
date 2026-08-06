package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/marketingleads/internal/mocks"
)

func nuevoCaso(repo *mocks.RepositoryMock, wa *mocks.WhatsAppSenderMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if wa == nil {
		return New(repo, nil, mocks.NewSilentLogger())
	}
	return New(repo, wa, mocks.NewSilentLogger())
}

func leadValido() dtos.CreateLeadDTO {
	return dtos.CreateLeadDTO{
		Name:       "Sebastian Camacho",
		Email:      "sebas@tienda.com",
		Phone:      "573001112233",
		SurveySlug: "diagnostico-logistico",
		ScoreTotal: 34,
		ScoreMax:   50,
		Level:      "intermedio",
	}
}

func TestCreateLead_GuardaElLeadConLosDatosRecortados(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := leadValido()
	dto.Name = "  Sebastian Camacho  "
	dto.Email = "\tsebas@tienda.com\n"
	dto.Phone = " 573001112233 "

	err := nuevoCaso(repo, nil).CreateLead(context.Background(), dto)

	require.NoError(t, err)
	require.Len(t, repo.Created, 1)
	assert.Equal(t, "Sebastian Camacho", repo.Created[0].Name,
		"el formulario publico llega con espacios pegados; se normalizan antes de persistir")
	assert.Equal(t, "sebas@tienda.com", repo.Created[0].Email)
	assert.Equal(t, "573001112233", repo.Created[0].Phone)
}

func TestCreateLead_RechazaCuandoFaltaNombreCorreoOTelefono(t *testing.T) {
	casos := map[string]func(*dtos.CreateLeadDTO){
		"sin nombre":           func(d *dtos.CreateLeadDTO) { d.Name = "" },
		"sin correo":           func(d *dtos.CreateLeadDTO) { d.Email = "" },
		"sin telefono":         func(d *dtos.CreateLeadDTO) { d.Phone = "" },
		"nombre solo espacios": func(d *dtos.CreateLeadDTO) { d.Name = "   " },
		"correo solo espacios": func(d *dtos.CreateLeadDTO) { d.Email = "  " },
		"telefono solo tabs":   func(d *dtos.CreateLeadDTO) { d.Phone = "\t\t" },
	}

	for nombre, romper := range casos {
		t.Run(nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			dto := leadValido()
			romper(&dto)

			err := nuevoCaso(repo, nil).CreateLead(context.Background(), dto)

			assert.ErrorIs(t, err, domainerrors.ErrInvalidLead)
			assert.Empty(t, repo.Created, "un lead incompleto no llega a la base")
		})
	}
}

func TestCreateLead_ConservaPuntajeNivelRespuestasYRecomendaciones(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	dto := leadValido()
	dto.Answers = []entities.LeadAnswer{{Question: "Cuantas ordenes al mes?", Answer: "500"}}
	dto.Recommendations = []entities.LeadRecommendation{{Title: "Automatiza guias", Desc: "d", Type: "logistica"}}

	require.NoError(t, nuevoCaso(repo, nil).CreateLead(context.Background(), dto))

	guardado := repo.Created[0]
	assert.Equal(t, 34, guardado.ScoreTotal)
	assert.Equal(t, 50, guardado.ScoreMax)
	assert.Equal(t, "intermedio", guardado.Level)
	assert.Equal(t, "diagnostico-logistico", guardado.SurveySlug)
	assert.Len(t, guardado.Answers, 1)
	assert.Len(t, guardado.Recommendations, 1)
}

func TestCreateLead_SiFallaElRepositorioNoSeEnviaWhatsApp(t *testing.T) {
	fallo := errors.New("conexion perdida")
	repo := &mocks.RepositoryMock{
		CreateLeadFn: func(ctx context.Context, lead *entities.MarketingLead) error { return fallo },
	}
	wa := &mocks.WhatsAppSenderMock{}

	err := nuevoCaso(repo, wa).CreateLead(context.Background(), leadValido())

	assert.ErrorIs(t, err, fallo)
	assert.Empty(t, wa.Calls,
		"sin lead guardado no se manda el mensaje: evita prometerle resultados a alguien que no quedo registrado")
}

func TestCreateLead_EnviaElWhatsAppConElTituloDeLaEncuesta(t *testing.T) {
	casos := []struct {
		slug     string
		esperado string
	}{
		{"diagnostico-logistico", "Diagn\u00f3stico Log\u00edstico"},
		{"escanea-tu-logistica", "Diagn\u00f3stico de Operaci\u00f3n E-commerce"},
		{"encuesta-que-no-existe", "Diagn\u00f3stico"},
		{"", "Diagn\u00f3stico"},
	}

	for _, caso := range casos {
		t.Run(caso.slug, func(t *testing.T) {
			wa := &mocks.WhatsAppSenderMock{}
			dto := leadValido()
			dto.SurveySlug = caso.slug

			require.NoError(t, nuevoCaso(nil, wa).CreateLead(context.Background(), dto))

			require.Len(t, wa.Calls, 1)
			assert.Equal(t, caso.esperado, wa.Calls[0].SurveyTitle,
				"un slug desconocido cae al titulo generico en vez de mandar el mensaje vacio")
		})
	}
}

func TestCreateLead_ElWhatsAppUsaElTelefonoYNombreYaNormalizados(t *testing.T) {
	wa := &mocks.WhatsAppSenderMock{}
	dto := leadValido()
	dto.Name = "  Sebastian  "
	dto.Phone = "  573001112233  "

	require.NoError(t, nuevoCaso(nil, wa).CreateLead(context.Background(), dto))

	require.Len(t, wa.Calls, 1)
	assert.Equal(t, "573001112233", wa.Calls[0].Phone,
		"un telefono con espacios rompe el envio en la API de Meta")
	assert.Equal(t, "Sebastian", wa.Calls[0].Name)
	assert.Equal(t, "intermedio", wa.Calls[0].Level)
}

func TestCreateLead_GuardaElMessageIDDevueltoPorWhatsApp(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	wa := &mocks.WhatsAppSenderMock{
		SendFn: func(ctx context.Context, phone, name, surveyTitle, level string) (string, error) {
			return "wamid.HBgMNTcz", nil
		},
	}

	require.NoError(t, nuevoCaso(repo, wa).CreateLead(context.Background(), leadValido()))

	require.Len(t, repo.MessageIDSet, 1)
	assert.Equal(t, "wamid.HBgMNTcz", repo.MessageIDSet[0].MessageID)
	assert.Equal(t, repo.Created[0].ID, repo.MessageIDSet[0].LeadID,
		"el message id se ata al lead recien creado para poder rastrear la entrega despues")
}

func TestCreateLead_SiWhatsAppNoDevuelveMessageIDNoSeIntentaGuardarlo(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	wa := &mocks.WhatsAppSenderMock{
		SendFn: func(ctx context.Context, phone, name, surveyTitle, level string) (string, error) {
			return "", nil
		},
	}

	require.NoError(t, nuevoCaso(repo, wa).CreateLead(context.Background(), leadValido()))

	assert.Empty(t, repo.MessageIDSet)
}

func TestCreateLead_SiFallaElEnvioDeWhatsAppElLeadSigueGuardado(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	wa := &mocks.WhatsAppSenderMock{
		SendFn: func(ctx context.Context, phone, name, surveyTitle, level string) (string, error) {
			return "", errors.New("Meta devolvio 400: plantilla no aprobada")
		},
	}

	err := nuevoCaso(repo, wa).CreateLead(context.Background(), leadValido())

	assert.NoError(t, err,
		"el lead es el activo comercial: un fallo de WhatsApp no puede devolverle error al formulario ni perder el contacto")
	assert.Len(t, repo.Created, 1)
	assert.Empty(t, repo.MessageIDSet)
}

func TestCreateLead_SiFallaGuardarElMessageIDElLeadIgualQuedaCreado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		SetWhatsAppMessageIDFn: func(ctx context.Context, leadID uint, messageID string) error {
			return errors.New("update fallido")
		},
	}

	err := nuevoCaso(repo, &mocks.WhatsAppSenderMock{}).CreateLead(context.Background(), leadValido())

	assert.NoError(t, err, "solo se pierde la trazabilidad del mensaje, no el lead")
	assert.Len(t, repo.Created, 1)
}

func TestCreateLead_SinIntegracionDeWhatsAppSoloGuardaElLead(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	require.NoError(t, nuevoCaso(repo, nil).CreateLead(context.Background(), leadValido()))

	assert.Len(t, repo.Created, 1)
	assert.Empty(t, repo.MessageIDSet)
}

func TestListLeads_NormalizaPaginaYTamanoFueraDeRango(t *testing.T) {
	casos := []struct {
		nombre           string
		page, pageSize   int
		pageEsp, sizeEsp int
	}{
		{"pagina cero", 0, 20, 1, 20},
		{"pagina negativa", -5, 20, 1, 20},
		{"tamano cero", 1, 0, 1, 20},
		{"tamano negativo", 1, -1, 1, 20},
		{"tamano por encima del maximo", 1, 101, 1, 20},
		{"tamano en el maximo permitido", 1, 100, 1, 100},
		{"valores validos se respetan", 3, 15, 3, 15},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}

			_, res, err := nuevoCaso(repo, nil).ListLeads(context.Background(), caso.page, caso.pageSize)

			require.NoError(t, err)
			require.Len(t, repo.ListCalls, 1)
			assert.Equal(t, caso.pageEsp, repo.ListCalls[0].Page,
				"la paginacion se normaliza antes del query para no mandar un OFFSET negativo a postgres")
			assert.Equal(t, caso.sizeEsp, repo.ListCalls[0].PageSize)
			assert.Equal(t, caso.pageEsp, res.Page)
			assert.Equal(t, caso.sizeEsp, res.PageSize)
		})
	}
}

func TestListLeads_CalculaElTotalDePaginas(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize int
		esperado int
	}{
		{0, 20, 1},
		{1, 20, 1},
		{20, 20, 1},
		{21, 20, 2},
		{100, 20, 5},
		{101, 20, 6},
	}

	for _, caso := range casos {
		repo := &mocks.RepositoryMock{
			ListLeadsFn: func(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error) {
				return nil, caso.total, nil
			},
		}

		_, res, err := nuevoCaso(repo, nil).ListLeads(context.Background(), 1, caso.pageSize)

		require.NoError(t, err)
		assert.Equal(t, caso.esperado, res.TotalPages,
			"con cero resultados el front igual espera una pagina, no cero")
		assert.Equal(t, caso.total, res.Total)
	}
}

func TestListLeads_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListLeadsFn: func(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error) {
			return nil, 0, fallo
		},
	}

	leads, res, err := nuevoCaso(repo, nil).ListLeads(context.Background(), 1, 20)

	assert.ErrorIs(t, err, fallo)
	assert.Nil(t, leads)
	assert.Zero(t, res.Total)
}

func TestListLeads_DevuelveLosLeadsDelRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListLeadsFn: func(ctx context.Context, page, pageSize int) ([]entities.MarketingLead, int64, error) {
			return []entities.MarketingLead{{ID: 7, Name: "Ana"}, {ID: 8, Name: "Luis"}}, 2, nil
		},
	}

	leads, res, err := nuevoCaso(repo, nil).ListLeads(context.Background(), 1, 20)

	require.NoError(t, err)
	require.Len(t, leads, 2)
	assert.Equal(t, "Ana", leads[0].Name)
	assert.EqualValues(t, 2, res.Total)
}
