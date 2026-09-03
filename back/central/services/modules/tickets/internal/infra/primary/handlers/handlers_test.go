package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/app"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type useCaseMock struct {
	CreateFn           func(ctx context.Context, dto dtos.CreateTicketDTO) (*entities.Ticket, error)
	GetFn              func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error)
	ListFn             func(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error)
	ListCategoriesFn   func(ctx context.Context) ([]string, error)
	UpdateFn           func(ctx context.Context, dto dtos.UpdateTicketDTO) (*entities.Ticket, error)
	DeleteFn           func(ctx context.Context, id uint) error
	ChangeStatusFn     func(ctx context.Context, dto dtos.ChangeStatusDTO) (*entities.Ticket, error)
	ChangeAreaFn       func(ctx context.Context, dto dtos.ChangeAreaDTO) (*entities.Ticket, error)
	AssignFn           func(ctx context.Context, dto dtos.AssignTicketDTO) (*entities.Ticket, error)
	ChangeSprintFn     func(ctx context.Context, dto dtos.ChangeSprintDTO) (*entities.Ticket, error)
	EscalateFn         func(ctx context.Context, dto dtos.EscalateTicketDTO) (*entities.Ticket, error)
	AddCommentFn       func(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error)
	ListCommentsFn     func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error)
	UploadAttachmentFn func(ctx context.Context, ticketID uint, commentID *uint, uploaderID uint, file *multipart.FileHeader) (*entities.TicketAttachment, error)
	ListAttachmentsFn  func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error)
	DeleteAttachmentFn func(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error
	ListHistoryFn      func(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error)

	CreateDTO        *dtos.CreateTicketDTO
	UpdateDTO        *dtos.UpdateTicketDTO
	ListParams       *dtos.ListTicketsParams
	StatusDTO        *dtos.ChangeStatusDTO
	AreaDTO          *dtos.ChangeAreaDTO
	AssignDTO        *dtos.AssignTicketDTO
	SprintDTO        *dtos.ChangeSprintDTO
	EscalateDTO      *dtos.EscalateTicketDTO
	CommentDTO       *dtos.CreateCommentDTO
	BorradoID        *uint
	AdjuntoBorrado   []uint
	AdjuntoSubido    *multipart.FileHeader
	ComentarioSubido *uint
}

var _ app.IUseCase = (*useCaseMock)(nil)

func (m *useCaseMock) Create(ctx context.Context, dto dtos.CreateTicketDTO) (*entities.Ticket, error) {
	m.CreateDTO = &dto
	if m.CreateFn != nil {
		return m.CreateFn(ctx, dto)
	}
	return &entities.Ticket{ID: 1, Code: "TKT-000001", Title: dto.Title}, nil
}

func (m *useCaseMock) Get(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
	if m.GetFn != nil {
		return m.GetFn(ctx, id, userID, businessID, isSuperAdmin)
	}
	return &entities.Ticket{ID: id}, nil
}

func (m *useCaseMock) List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
	m.ListParams = &params
	if m.ListFn != nil {
		return m.ListFn(ctx, params)
	}
	return nil, 0, nil
}

func (m *useCaseMock) Update(ctx context.Context, dto dtos.UpdateTicketDTO) (*entities.Ticket, error) {
	m.UpdateDTO = &dto
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.ID}, nil
}

func (m *useCaseMock) ListCategories(ctx context.Context) ([]string, error) {
	if m.ListCategoriesFn != nil {
		return m.ListCategoriesFn(ctx)
	}
	return []string{}, nil
}

func (m *useCaseMock) Delete(ctx context.Context, id uint) error {
	m.BorradoID = &id
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	return nil
}

func (m *useCaseMock) ChangeStatus(ctx context.Context, dto dtos.ChangeStatusDTO) (*entities.Ticket, error) {
	m.StatusDTO = &dto
	if m.ChangeStatusFn != nil {
		return m.ChangeStatusFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.TicketID, Status: dto.NewStatus}, nil
}

func (m *useCaseMock) ChangeArea(ctx context.Context, dto dtos.ChangeAreaDTO) (*entities.Ticket, error) {
	m.AreaDTO = &dto
	if m.ChangeAreaFn != nil {
		return m.ChangeAreaFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.TicketID, Area: dto.NewArea}, nil
}

func (m *useCaseMock) Assign(ctx context.Context, dto dtos.AssignTicketDTO) (*entities.Ticket, error) {
	m.AssignDTO = &dto
	if m.AssignFn != nil {
		return m.AssignFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.TicketID, AssignedToID: dto.AssignedToID}, nil
}

func (m *useCaseMock) ChangeSprint(ctx context.Context, dto dtos.ChangeSprintDTO) (*entities.Ticket, error) {
	m.SprintDTO = &dto
	if m.ChangeSprintFn != nil {
		return m.ChangeSprintFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.TicketID, SprintID: dto.SprintID}, nil
}

func (m *useCaseMock) Escalate(ctx context.Context, dto dtos.EscalateTicketDTO) (*entities.Ticket, error) {
	m.EscalateDTO = &dto
	if m.EscalateFn != nil {
		return m.EscalateFn(ctx, dto)
	}
	return &entities.Ticket{ID: dto.TicketID, EscalatedToDev: true}, nil
}

func (m *useCaseMock) AddComment(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error) {
	m.CommentDTO = &dto
	if m.AddCommentFn != nil {
		return m.AddCommentFn(ctx, dto)
	}
	return &entities.TicketComment{ID: 1, TicketID: dto.TicketID, Body: dto.Body, IsInternal: dto.IsInternal}, nil
}

func (m *useCaseMock) ListComments(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
	if m.ListCommentsFn != nil {
		return m.ListCommentsFn(ctx, ticketID, includeInternal)
	}
	return nil, nil
}

func (m *useCaseMock) UploadAttachment(ctx context.Context, ticketID uint, commentID *uint, uploaderID uint, file *multipart.FileHeader) (*entities.TicketAttachment, error) {
	m.AdjuntoSubido = file
	m.ComentarioSubido = commentID
	if m.UploadAttachmentFn != nil {
		return m.UploadAttachmentFn(ctx, ticketID, commentID, uploaderID, file)
	}
	return &entities.TicketAttachment{ID: 1, TicketID: ticketID, CommentID: commentID, UploadedByID: uploaderID}, nil
}

func (m *useCaseMock) ListAttachments(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
	if m.ListAttachmentsFn != nil {
		return m.ListAttachmentsFn(ctx, ticketID)
	}
	return nil, nil
}

func (m *useCaseMock) DeleteAttachment(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error {
	m.AdjuntoBorrado = append(m.AdjuntoBorrado, attachmentID)
	if m.DeleteAttachmentFn != nil {
		return m.DeleteAttachmentFn(ctx, attachmentID, requesterID, isSuperAdmin)
	}
	return nil
}

func (m *useCaseMock) ListHistory(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
	if m.ListHistoryFn != nil {
		return m.ListHistoryFn(ctx, ticketID)
	}
	return nil, nil
}

func sesionSuperAdmin(c *gin.Context) {
	c.Set("user_id", uint(42))
	c.Set("business_id", uint(0))
	c.Set("auth_info", &middleware.AuthInfo{UserID: 42, BusinessID: 0})
}

func sesionNegocio(c *gin.Context) {
	c.Set("user_id", uint(42))
	c.Set("business_id", uint(26))
	c.Set("auth_info", &middleware.AuthInfo{UserID: 42, BusinessID: 26})
}

func sesionAnonima(c *gin.Context) {}

func nuevoRouter(uc app.IUseCase, sesion func(*gin.Context)) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		sesion(c)
		c.Next()
	})
	h := &Handlers{uc: uc, log: mocks.NewSilentLogger()}
	g := r.Group("/tickets")
	g.GET("", h.List)
	g.POST("", h.Create)
	g.GET(":id", h.Get)
	g.PUT(":id", h.Update)
	g.DELETE(":id", h.Delete)
	g.PATCH(":id/status", h.ChangeStatus)
	g.PATCH(":id/assign", h.Assign)
	g.PATCH(":id/sprint", h.ChangeSprint)
	g.PATCH(":id/area", h.ChangeArea)
	g.PATCH(":id/escalate", h.Escalate)
	g.GET(":id/comments", h.ListComments)
	g.POST(":id/comments", h.AddComment)
	g.GET(":id/attachments", h.ListAttachments)
	g.POST(":id/attachments", h.UploadAttachment)
	g.DELETE("attachments/:attachment_id", h.DeleteAttachment)
	g.GET(":id/history", h.ListHistory)
	return r
}

func llamar(t *testing.T, r *gin.Engine, metodo, ruta string, cuerpo any) *httptest.ResponseRecorder {
	t.Helper()
	var lector *bytes.Reader
	if cuerpo != nil {
		var b []byte
		switch v := cuerpo.(type) {
		case string:
			b = []byte(v)
		default:
			var err error
			b, err = json.Marshal(v)
			require.NoError(t, err)
		}
		lector = bytes.NewReader(b)
	} else {
		lector = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(metodo, ruta, lector)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func cuerpo(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out),
		"la respuesta no es JSON: %q", rec.Body.String())
	return out
}

func TestCreate_SinCamposObligatorios_400(t *testing.T) {
	casos := []struct {
		nombre string
		body   any
	}{
		{"sin titulo", map[string]any{"description": "d"}},
		{"sin descripcion", map[string]any{"title": "t"}},
		{"json roto", `{"title":`},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			uc := &useCaseMock{}
			rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets", tc.body)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Nil(t, uc.CreateDTO, "no se llama al caso de uso con un payload invalido")
			assert.Contains(t, cuerpo(t, rec), "error")
		})
	}
}

func TestCreate_SinSesion_401(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionAnonima), http.MethodPost, "/tickets",
		map[string]any{"title": "t", "description": "d"})

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Nil(t, uc.CreateDTO)
}

func TestCreate_SuperAdmin_RespetaOrigenResponsableYSprint(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets", map[string]any{
		"title": "t", "description": "d", "source": "business",
		"business_id": 26, "assigned_to_id": 7, "sprint_id": 9,
	})

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, uc.CreateDTO)
	assert.Equal(t, "business", uc.CreateDTO.Source)
	require.NotNil(t, uc.CreateDTO.BusinessID)
	assert.Equal(t, uint(26), *uc.CreateDTO.BusinessID)
	require.NotNil(t, uc.CreateDTO.AssignedToID)
	assert.Equal(t, uint(7), *uc.CreateDTO.AssignedToID)
	require.NotNil(t, uc.CreateDTO.SprintID)
	assert.Equal(t, uint(9), *uc.CreateDTO.SprintID)
	assert.Equal(t, uint(42), uc.CreateDTO.CreatedByID)
}

func TestCreate_SuperAdminSinOrigen_CaeAInterno(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets",
		map[string]any{"title": "t", "description": "d"})

	require.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "internal", uc.CreateDTO.Source)
}

func TestCreate_UsuarioDeNegocio_IgnoraLoQueNoLeCorresponde(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPost, "/tickets", map[string]any{
		"title": "t", "description": "d", "source": "internal",
		"business_id": 99, "assigned_to_id": 7, "sprint_id": 9,
	})

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, uc.CreateDTO)
	assert.Equal(t, "business", uc.CreateDTO.Source,
		"un cliente no puede hacer pasar su ticket por interno")
	require.NotNil(t, uc.CreateDTO.BusinessID)
	assert.Equal(t, uint(26), *uc.CreateDTO.BusinessID,
		"el negocio sale del token, no del body")
	assert.Nil(t, uc.CreateDTO.AssignedToID, "un cliente no asigna responsables")
	assert.Nil(t, uc.CreateDTO.SprintID, "un cliente no mete tickets a un sprint")
}

func TestCreate_ErrorDeValidacionDelCasoDeUso_400(t *testing.T) {
	uc := &useCaseMock{
		CreateFn: func(ctx context.Context, dto dtos.CreateTicketDTO) (*entities.Ticket, error) {
			return nil, dom.ErrInvalidType
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets",
		map[string]any{"title": "t", "description": "d"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, dom.ErrInvalidType.Error(), cuerpo(t, rec)["error"])
}

func TestGet_IDInvalido_400(t *testing.T) {
	for _, ruta := range []string{"/tickets/abc", "/tickets/0"} {
		t.Run(ruta, func(t *testing.T) {
			uc := &useCaseMock{
				GetFn: func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
					t.Fatal("no se debe consultar con un id invalido")
					return nil, nil
				},
			}

			rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, ruta, nil)

			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Equal(t, "invalid id", cuerpo(t, rec)["error"])
		})
	}
}

func TestGet_PasaElContextoDelSolicitante(t *testing.T) {
	var vistoID, vistoUsuario uint
	var vistoNegocio *uint
	var vistoSuper bool
	uc := &useCaseMock{
		GetFn: func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
			vistoID, vistoUsuario, vistoNegocio, vistoSuper = id, userID, businessID, isSuperAdmin
			return &entities.Ticket{ID: id, Code: "TKT-000005"}, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodGet, "/tickets/5", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, uint(5), vistoID)
	assert.Equal(t, uint(42), vistoUsuario)
	require.NotNil(t, vistoNegocio)
	assert.Equal(t, uint(26), *vistoNegocio)
	assert.False(t, vistoSuper)
	assert.Equal(t, "TKT-000005", cuerpo(t, rec)["code"])
}

func TestGet_SuperAdminPuedeMirarElNegocioQuePideEnLaQuery(t *testing.T) {
	var vistoNegocio *uint
	var vistoSuper bool
	uc := &useCaseMock{
		GetFn: func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
			vistoNegocio, vistoSuper = businessID, isSuperAdmin
			return &entities.Ticket{ID: id}, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5?business_id=26", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.True(t, vistoSuper)
	require.NotNil(t, vistoNegocio)
	assert.Equal(t, uint(26), *vistoNegocio)
}

func TestGet_SuperAdminConBusinessIDBasura_LoIgnora(t *testing.T) {
	for _, query := range []string{"?business_id=abc", "?business_id=0", ""} {
		var vistoNegocio *uint
		uc := &useCaseMock{
			GetFn: func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
				vistoNegocio = businessID
				return &entities.Ticket{ID: id}, nil
			},
		}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5"+query, nil)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Nil(t, vistoNegocio, "query %q", query)
	}
}

func TestGet_ErroresDelCasoDeUso_SeTraducenAlEstadoHTTP(t *testing.T) {
	casos := []struct {
		err  error
		code int
	}{
		{dom.ErrTicketNotFound, http.StatusNotFound},
		{dom.ErrCommentNotFound, http.StatusNotFound},
		{dom.ErrAttachmentNotFound, http.StatusNotFound},
		{dom.ErrForbidden, http.StatusForbidden},
		{dom.ErrInvalidStatus, http.StatusBadRequest},
		{dom.ErrInvalidPriority, http.StatusBadRequest},
		{dom.ErrInvalidType, http.StatusBadRequest},
		{dom.ErrInvalidSeverity, http.StatusBadRequest},
		{dom.ErrInvalidArea, http.StatusBadRequest},
		{dom.ErrTitleRequired, http.StatusBadRequest},
		{dom.ErrDescriptionRequired, http.StatusBadRequest},
		{dom.ErrAssigneeNotFound, http.StatusBadRequest},
		{dom.ErrSprintNotFound, http.StatusBadRequest},
		{stderrors.New("db caida"), http.StatusInternalServerError},
	}

	for _, tc := range casos {
		t.Run(tc.err.Error(), func(t *testing.T) {
			uc := &useCaseMock{
				GetFn: func(ctx context.Context, id uint, userID uint, businessID *uint, isSuperAdmin bool) (*entities.Ticket, error) {
					return nil, tc.err
				},
			}

			rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5", nil)

			assert.Equal(t, tc.code, rec.Code)
			assert.Equal(t, tc.err.Error(), cuerpo(t, rec)["error"])
		})
	}
}

func TestList_ValoresPorDefectoDePaginacion(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, uc.ListParams)
	assert.Equal(t, 1, uc.ListParams.Page)
	assert.Equal(t, 10, uc.ListParams.PageSize)

	body := cuerpo(t, rec)
	assert.Equal(t, float64(0), body["total"])
	assert.Equal(t, float64(0), body["total_pages"])
	assert.Equal(t, []any{}, body["data"], "sin resultados el front recibe [] y no null")
}

func TestList_CalculaLasPaginasRedondeandoHaciaArriba(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize string
		esperado float64
	}{
		{0, "10", 0},
		{10, "10", 1},
		{11, "10", 2},
		{99, "25", 4},
	}

	for _, tc := range casos {
		uc := &useCaseMock{
			ListFn: func(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
				return []entities.Ticket{{ID: 1}}, tc.total, nil
			},
		}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet,
			"/tickets?page_size="+tc.pageSize, nil)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tc.esperado, cuerpo(t, rec)["total_pages"], "total %d", tc.total)
	}
}

func TestList_TraduceLosFiltrosDeLaQuery(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet,
		"/tickets?page=3&page_size=25&status=open,%20closed&priority=high&type=bug&area=soporte,desarrollo"+
			"&source=business&escalated=true&search=factura&sort_by=code&sort_order=asc&only_mine=true", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	p := uc.ListParams
	require.NotNil(t, p)
	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 25, p.PageSize)
	assert.Equal(t, []string{"open", "closed"}, p.Status, "el CSV se separa y se recorta")
	assert.Equal(t, []string{"high"}, p.Priority)
	assert.Equal(t, []string{"bug"}, p.Type)
	assert.Equal(t, []string{"soporte", "desarrollo"}, p.Area)
	assert.Equal(t, "business", p.Source)
	assert.True(t, p.EscalatedOnly)
	assert.Equal(t, "factura", p.Search)
	assert.Equal(t, "code", p.SortBy)
	assert.Equal(t, "asc", p.SortOrder)
	assert.True(t, p.OnlyMine)
	assert.Equal(t, uint(42), p.UserID)
	assert.True(t, p.IsSuperAdmin)
}

func TestList_FiltrosVaciosNoGeneranListasFantasma(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet,
		"/tickets?status=&priority=,,&escalated=false&only_mine=no", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Nil(t, uc.ListParams.Status)
	assert.Empty(t, uc.ListParams.Priority)
	assert.False(t, uc.ListParams.EscalatedOnly)
	assert.False(t, uc.ListParams.OnlyMine)
}

func TestList_FiltroDeSprint(t *testing.T) {
	casos := []struct {
		query      string
		wantNone   bool
		wantSprint *uint
	}{
		{"", false, nil},
		{"&sprint_id=none", true, nil},
		{"&sprint_id=NULL", true, nil},
		{"&sprint_id=%20none%20", true, nil},
		{"&sprint_id=9", false, func() *uint { v := uint(9); return &v }()},
		{"&sprint_id=abc", false, nil},
		{"&sprint_id=0", false, nil},
	}

	for _, tc := range casos {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets?page=1"+tc.query, nil)

		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, tc.wantNone, uc.ListParams.SprintNone, "query %q", tc.query)
		if tc.wantSprint == nil {
			assert.Nil(t, uc.ListParams.SprintID, "query %q", tc.query)
		} else {
			require.NotNil(t, uc.ListParams.SprintID, "query %q", tc.query)
			assert.Equal(t, *tc.wantSprint, *uc.ListParams.SprintID)
		}
	}
}

func TestList_SuperAdminFiltraPorNegocioCreadorYResponsable(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet,
		"/tickets?business_id=26&created_by_id=42&assigned_to_id=7", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	p := uc.ListParams
	require.NotNil(t, p.BusinessID)
	assert.Equal(t, uint(26), *p.BusinessID)
	require.NotNil(t, p.CreatedByID)
	assert.Equal(t, uint(42), *p.CreatedByID)
	require.NotNil(t, p.AssignedToID)
	assert.Equal(t, uint(7), *p.AssignedToID)
}

func TestList_SuperAdminConFiltrosBasura_LosIgnora(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet,
		"/tickets?business_id=abc&created_by_id=xx&assigned_to_id=yy", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Nil(t, uc.ListParams.BusinessID)
	assert.Nil(t, uc.ListParams.CreatedByID)
	assert.Nil(t, uc.ListParams.AssignedToID)
}

func TestList_UsuarioDeNegocio_QuedaEncerradoEnSuNegocio(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodGet,
		"/tickets?business_id=99&created_by_id=1&assigned_to_id=2", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	p := uc.ListParams
	require.NotNil(t, p.BusinessID)
	assert.Equal(t, uint(26), *p.BusinessID,
		"el negocio del token gana sobre el de la query")
	assert.Nil(t, p.CreatedByID, "los filtros de super admin no aplican a un cliente")
	assert.Nil(t, p.AssignedToID)
	assert.False(t, p.IsSuperAdmin)
}

func TestList_UsuarioSinNegocioNiSuperAdmin_NoRecibeNegocio(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionAnonima), http.MethodGet, "/tickets", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Nil(t, uc.ListParams.BusinessID,
		"sin negocio el repositorio no debe recibir un filtro que lo abra todo")
	assert.False(t, uc.ListParams.IsSuperAdmin)
}

func TestList_DevuelveLosTicketsMapeados(t *testing.T) {
	uc := &useCaseMock{
		ListFn: func(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
			return []entities.Ticket{
				{ID: 1, Code: "TKT-000001", Title: "uno"},
				{ID: 2, Code: "TKT-000002", Title: "dos"},
			}, 2, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	data, ok := cuerpo(t, rec)["data"].([]any)
	require.True(t, ok)
	require.Len(t, data, 2)
	assert.Equal(t, "TKT-000001", data[0].(map[string]any)["code"])
	assert.Equal(t, "dos", data[1].(map[string]any)["title"])
}

func TestList_ErrorDelCasoDeUso_500(t *testing.T) {
	uc := &useCaseMock{
		ListFn: func(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
			return nil, 0, stderrors.New("db caida")
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets", nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestUpdate_SoloSuperAdmin(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPut, "/tickets/5",
		map[string]any{"title": "nuevo"})

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, uc.UpdateDTO)
}

func TestUpdate_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPut, "/tickets/abc",
		map[string]any{"title": "nuevo"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.UpdateDTO)
}

func TestUpdate_JSONInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPut, "/tickets/5", `{"title":`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.UpdateDTO)
}

func TestUpdate_TrasladaTodosLosCamposAlCasoDeUso(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPut, "/tickets/5", map[string]any{
		"title": "nuevo", "description": "otra", "type": "bug", "category": "ordenes",
		"priority": "high", "severity": "low", "area": "desarrollo",
		"assigned_to_id": 7, "sprint_id": 9, "clear_sprint": true,
		"due_date": "2026-12-31", "clear_due_date": true,
	})

	require.Equal(t, http.StatusOK, rec.Code)
	d := uc.UpdateDTO
	require.NotNil(t, d)
	assert.Equal(t, uint(5), d.ID)
	assert.Equal(t, "nuevo", *d.Title)
	assert.Equal(t, "otra", *d.Description)
	assert.Equal(t, "bug", *d.Type)
	assert.Equal(t, "ordenes", *d.Category)
	assert.Equal(t, "high", *d.Priority)
	assert.Equal(t, "low", *d.Severity)
	assert.Equal(t, "desarrollo", *d.Area)
	assert.Equal(t, uint(7), *d.AssignedToID)
	assert.Equal(t, uint(9), *d.SprintID)
	assert.True(t, d.ClearSprint)
	assert.Equal(t, "2026-12-31", *d.DueDate)
	assert.True(t, d.ClearDueDate)
}

func TestUpdate_CamposAusentes_ViajanComoNil(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPut, "/tickets/5", map[string]any{})

	require.Equal(t, http.StatusOK, rec.Code)
	d := uc.UpdateDTO
	require.NotNil(t, d)
	assert.Nil(t, d.Title, "un campo ausente no es un campo vacio: no debe borrar nada")
	assert.Nil(t, d.Description)
	assert.Nil(t, d.SprintID)
	assert.False(t, d.ClearSprint)
	assert.False(t, d.ClearDueDate)
}

func TestUpdate_TicketInexistente_404(t *testing.T) {
	uc := &useCaseMock{
		UpdateFn: func(ctx context.Context, dto dtos.UpdateTicketDTO) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPut, "/tickets/404",
		map[string]any{"title": "nuevo"})

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDelete_SoloSuperAdmin(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodDelete, "/tickets/5", nil)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, uc.BorradoID)
}

func TestDelete_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodDelete, "/tickets/abc", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.BorradoID)
}

func TestDelete_Correcto_200(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodDelete, "/tickets/5", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, uc.BorradoID)
	assert.Equal(t, uint(5), *uc.BorradoID)
	assert.Equal(t, true, cuerpo(t, rec)["success"])
}

func TestDelete_ErrorDelCasoDeUso_500(t *testing.T) {
	uc := &useCaseMock{
		DeleteFn: func(ctx context.Context, id uint) error { return stderrors.New("fk violation") },
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodDelete, "/tickets/5", nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
