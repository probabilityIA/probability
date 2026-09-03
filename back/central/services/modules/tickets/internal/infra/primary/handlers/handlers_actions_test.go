package handlers

import (
	"bytes"
	"context"
	stderrors "errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChangeStatus_SinEstado_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/status",
		map[string]any{"note": "sin estado"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.StatusDTO, "el estado es obligatorio")
}

func TestChangeStatus_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/abc/status",
		map[string]any{"status": "closed"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.StatusDTO)
}

func TestChangeStatus_ArmaElDTOConElUsuarioDeLaSesion(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/status",
		map[string]any{"status": "closed", "note": "duplicado"})

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, uc.StatusDTO)
	assert.Equal(t, uint(5), uc.StatusDTO.TicketID)
	assert.Equal(t, "closed", uc.StatusDTO.NewStatus)
	assert.Equal(t, "duplicado", uc.StatusDTO.Note)
	assert.Equal(t, uint(42), uc.StatusDTO.ChangedByID)
	assert.Equal(t, "closed", cuerpo(t, rec)["status"])
}

func TestChangeStatus_UsuarioDeNegocioTambienPuede(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPatch, "/tickets/5/status",
		map[string]any{"status": "closed"})

	require.Equal(t, http.StatusOK, rec.Code,
		"cambiar el estado no esta reservado a super admin como si lo estan asignar o escalar")
	require.NotNil(t, uc.StatusDTO)
	assert.Equal(t, uint(42), uc.StatusDTO.ChangedByID)
}

func TestChangeStatus_EstadoInvalido_400(t *testing.T) {
	uc := &useCaseMock{
		ChangeStatusFn: func(ctx context.Context, dto dtos.ChangeStatusDTO) (*entities.Ticket, error) {
			return nil, dom.ErrInvalidStatus
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/status",
		map[string]any{"status": "inventado"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, dom.ErrInvalidStatus.Error(), cuerpo(t, rec)["error"])
}

func TestChangeArea_SinArea_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/area",
		map[string]any{"note": "sin area"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.AreaDTO)
}

func TestChangeArea_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/0/area",
		map[string]any{"area": "desarrollo"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.AreaDTO)
}

func TestChangeArea_ArmaElDTO(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/area",
		map[string]any{"area": "desarrollo", "note": "es un bug"})

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, uc.AreaDTO)
	assert.Equal(t, uint(5), uc.AreaDTO.TicketID)
	assert.Equal(t, "desarrollo", uc.AreaDTO.NewArea)
	assert.Equal(t, "es un bug", uc.AreaDTO.Note)
	assert.Equal(t, uint(42), uc.AreaDTO.ChangedByID)
}

func TestChangeArea_AreaInvalida_400(t *testing.T) {
	uc := &useCaseMock{
		ChangeAreaFn: func(ctx context.Context, dto dtos.ChangeAreaDTO) (*entities.Ticket, error) {
			return nil, dom.ErrInvalidArea
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/area",
		map[string]any{"area": "marketing"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAssign_SoloSuperAdmin(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPatch, "/tickets/5/assign",
		map[string]any{"assigned_to_id": 7})

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, uc.AssignDTO, "un cliente no reparte trabajo del equipo")
}

func TestAssign_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/abc/assign",
		map[string]any{"assigned_to_id": 7})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.AssignDTO)
}

func TestAssign_JSONInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/assign", `{`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.AssignDTO)
}

func TestAssign_ConYSinResponsable(t *testing.T) {
	t.Run("con responsable", func(t *testing.T) {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/assign",
			map[string]any{"assigned_to_id": 7})

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotNil(t, uc.AssignDTO)
		require.NotNil(t, uc.AssignDTO.AssignedToID)
		assert.Equal(t, uint(7), *uc.AssignDTO.AssignedToID)
		assert.Equal(t, uint(42), uc.AssignDTO.ChangedByID)
	})

	t.Run("desasignar", func(t *testing.T) {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/assign",
			map[string]any{"assigned_to_id": nil})

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotNil(t, uc.AssignDTO)
		assert.Nil(t, uc.AssignDTO.AssignedToID,
			"mandar null es la forma de quitarle el responsable al ticket")
	})
}

func TestAssign_ResponsableInexistente_400(t *testing.T) {
	uc := &useCaseMock{
		AssignFn: func(ctx context.Context, dto dtos.AssignTicketDTO) (*entities.Ticket, error) {
			return nil, dom.ErrAssigneeNotFound
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/assign",
		map[string]any{"assigned_to_id": 404})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestChangeSprint_SoloSuperAdmin(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPatch, "/tickets/5/sprint",
		map[string]any{"sprint_id": 9})

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, uc.SprintDTO)
}

func TestChangeSprint_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/abc/sprint",
		map[string]any{"sprint_id": 9})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.SprintDTO)
}

func TestChangeSprint_JSONInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/sprint", `{"sprint_id":`)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.SprintDTO)
}

func TestChangeSprint_ConYSinSprint(t *testing.T) {
	t.Run("mover a un sprint", func(t *testing.T) {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/sprint",
			map[string]any{"sprint_id": 9})

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotNil(t, uc.SprintDTO)
		assert.Equal(t, uint(5), uc.SprintDTO.TicketID)
		require.NotNil(t, uc.SprintDTO.SprintID)
		assert.Equal(t, uint(9), *uc.SprintDTO.SprintID)
		assert.Equal(t, uint(42), uc.SprintDTO.ChangedByID)
	})

	t.Run("sacar del sprint", func(t *testing.T) {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/sprint",
			map[string]any{"sprint_id": nil})

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotNil(t, uc.SprintDTO)
		assert.Nil(t, uc.SprintDTO.SprintID)
	})
}

func TestChangeSprint_SprintInexistente_400(t *testing.T) {
	uc := &useCaseMock{
		ChangeSprintFn: func(ctx context.Context, dto dtos.ChangeSprintDTO) (*entities.Ticket, error) {
			return nil, dom.ErrSprintNotFound
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/sprint",
		map[string]any{"sprint_id": 404})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, dom.ErrSprintNotFound.Error(), cuerpo(t, rec)["error"])
}

func TestEscalate_SoloSuperAdmin(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPatch, "/tickets/5/escalate",
		map[string]any{"note": "urgente"})

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Nil(t, uc.EscalateDTO)
}

func TestEscalate_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/abc/escalate", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.EscalateDTO)
}

func TestEscalate_SinCuerpoOConCuerpoRoto_IgualEscala(t *testing.T) {
	for _, body := range []any{nil, `{"note":`} {
		uc := &useCaseMock{}

		rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/escalate", body)

		require.Equal(t, http.StatusOK, rec.Code,
			"la nota es opcional: un cuerpo ausente o mal formado no bloquea el escalamiento")
		require.NotNil(t, uc.EscalateDTO)
		assert.Empty(t, uc.EscalateDTO.Note)
		assert.Equal(t, uint(5), uc.EscalateDTO.TicketID)
		assert.Equal(t, uint(42), uc.EscalateDTO.ChangedByID)
	}
}

func TestEscalate_ConNota_LaTraslada(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/5/escalate",
		map[string]any{"note": "bloquea a 3 clientes"})

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "bloquea a 3 clientes", uc.EscalateDTO.Note)
	assert.Equal(t, true, cuerpo(t, rec)["escalated_to_dev"])
}

func TestEscalate_TicketInexistente_404(t *testing.T) {
	uc := &useCaseMock{
		EscalateFn: func(ctx context.Context, dto dtos.EscalateTicketDTO) (*entities.Ticket, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPatch, "/tickets/404/escalate", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListComments_SoloElEquipoVeLosInternos(t *testing.T) {
	casos := []struct {
		nombre   string
		sesion   func(*gin.Context)
		esperado bool
	}{
		{"super admin", sesionSuperAdmin, true},
		{"cliente", sesionNegocio, false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var visto bool
			uc := &useCaseMock{
				ListCommentsFn: func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
					visto = includeInternal
					return []entities.TicketComment{{ID: 1, Body: "hola"}}, nil
				},
			}

			rec := llamar(t, nuevoRouter(uc, tc.sesion), http.MethodGet, "/tickets/5/comments", nil)

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tc.esperado, visto)
		})
	}
}

func TestListComments_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/abc/comments", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListComments_SinComentarios_DevuelveListaVacia(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/comments", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []any{}, cuerpo(t, rec)["data"])
}

func TestListComments_MapeaLosComentarios(t *testing.T) {
	uc := &useCaseMock{
		ListCommentsFn: func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
			return []entities.TicketComment{
				{ID: 1, Body: "primero"},
				{ID: 2, Body: "segundo", IsInternal: true},
			}, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/comments", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	data := cuerpo(t, rec)["data"].([]any)
	require.Len(t, data, 2)
	assert.Equal(t, "primero", data[0].(map[string]any)["body"])
	assert.Equal(t, true, data[1].(map[string]any)["is_internal"])
}

func TestListComments_ErrorDelCasoDeUso_500(t *testing.T) {
	uc := &useCaseMock{
		ListCommentsFn: func(ctx context.Context, ticketID uint, includeInternal bool) ([]entities.TicketComment, error) {
			return nil, stderrors.New("db caida")
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/comments", nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAddComment_SinCuerpo_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets/5/comments",
		map[string]any{"is_internal": true})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.CommentDTO)
}

func TestAddComment_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets/abc/comments",
		map[string]any{"body": "hola"})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.CommentDTO)
}

func TestAddComment_UnClienteNoPuedeDejarComentariosInternos(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodPost, "/tickets/5/comments",
		map[string]any{"body": "hola", "is_internal": true})

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, uc.CommentDTO)
	assert.False(t, uc.CommentDTO.IsInternal,
		"un comentario interno del equipo no puede originarse en un cliente")
	assert.Equal(t, uint(5), uc.CommentDTO.TicketID)
	assert.Equal(t, uint(42), uc.CommentDTO.UserID)
	assert.Equal(t, "hola", uc.CommentDTO.Body)
}

func TestAddComment_SuperAdminSiPuedeMarcarloInterno(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets/5/comments",
		map[string]any{"body": "ojo con esto", "is_internal": true})

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, uc.CommentDTO)
	assert.True(t, uc.CommentDTO.IsInternal)
	assert.Equal(t, true, cuerpo(t, rec)["is_internal"])
}

func TestAddComment_TicketInexistente_404(t *testing.T) {
	uc := &useCaseMock{
		AddCommentFn: func(ctx context.Context, dto dtos.CreateCommentDTO) (*entities.TicketComment, error) {
			return nil, dom.ErrTicketNotFound
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodPost, "/tickets/404/comments",
		map[string]any{"body": "hola"})

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestListAttachments_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/abc/attachments", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListAttachments_MapeaLosAdjuntos(t *testing.T) {
	uc := &useCaseMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			assert.Equal(t, uint(5), ticketID)
			return []entities.TicketAttachment{{ID: 10, FileName: "a.png", Size: 99}}, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/attachments", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	data := cuerpo(t, rec)["data"].([]any)
	require.Len(t, data, 1)
	assert.Equal(t, "a.png", data[0].(map[string]any)["file_name"])
	assert.Equal(t, float64(99), data[0].(map[string]any)["size"])
}

func TestListAttachments_SinAdjuntos_DevuelveListaVacia(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/attachments", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []any{}, cuerpo(t, rec)["data"])
}

func TestListAttachments_ErrorDelCasoDeUso_500(t *testing.T) {
	uc := &useCaseMock{
		ListAttachmentsFn: func(ctx context.Context, ticketID uint) ([]entities.TicketAttachment, error) {
			return nil, stderrors.New("db caida")
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/attachments", nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func peticionConArchivo(t *testing.T, ruta, campo, nombre string, extras map[string]string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if campo != "" {
		part, err := w.CreateFormFile(campo, nombre)
		require.NoError(t, err)
		_, err = part.Write([]byte("contenido"))
		require.NoError(t, err)
	}
	for k, v := range extras {
		require.NoError(t, w.WriteField(k, v))
	}
	require.NoError(t, w.Close())

	req := httptest.NewRequest(http.MethodPost, ruta, bytes.NewReader(buf.Bytes()))
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestUploadAttachment_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}
	rec := httptest.NewRecorder()

	nuevoRouter(uc, sesionSuperAdmin).ServeHTTP(rec,
		peticionConArchivo(t, "/tickets/abc/attachments", "file", "a.png", nil))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Nil(t, uc.AdjuntoSubido)
}

func TestUploadAttachment_SinArchivo_400(t *testing.T) {
	uc := &useCaseMock{}
	rec := httptest.NewRecorder()

	nuevoRouter(uc, sesionSuperAdmin).ServeHTTP(rec,
		peticionConArchivo(t, "/tickets/5/attachments", "", "", map[string]string{"comment_id": "3"}))

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "file is required", cuerpo(t, rec)["error"])
	assert.Nil(t, uc.AdjuntoSubido)
}

func TestUploadAttachment_EntregaElArchivoYElComentario(t *testing.T) {
	uc := &useCaseMock{}
	rec := httptest.NewRecorder()

	nuevoRouter(uc, sesionSuperAdmin).ServeHTTP(rec,
		peticionConArchivo(t, "/tickets/5/attachments", "file", "captura.png",
			map[string]string{"comment_id": "3"}))

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, uc.AdjuntoSubido)
	assert.Equal(t, "captura.png", uc.AdjuntoSubido.Filename)
	require.NotNil(t, uc.ComentarioSubido)
	assert.Equal(t, uint(3), *uc.ComentarioSubido)
	assert.Equal(t, float64(42), cuerpo(t, rec)["uploaded_by_id"])
}

func TestUploadAttachment_ComentarioInvalido_SeIgnora(t *testing.T) {
	for _, valor := range []string{"abc", "0", ""} {
		uc := &useCaseMock{}
		rec := httptest.NewRecorder()

		nuevoRouter(uc, sesionSuperAdmin).ServeHTTP(rec,
			peticionConArchivo(t, "/tickets/5/attachments", "file", "captura.png",
				map[string]string{"comment_id": valor}))

		require.Equal(t, http.StatusCreated, rec.Code, "comment_id %q", valor)
		assert.Nil(t, uc.ComentarioSubido,
			"un comment_id ilegible se ignora y el adjunto queda a nivel de ticket, no rompe la subida")
	}
}

func TestUploadAttachment_TicketInexistente_404(t *testing.T) {
	uc := &useCaseMock{
		UploadAttachmentFn: func(ctx context.Context, ticketID uint, commentID *uint, uploaderID uint, file *multipart.FileHeader) (*entities.TicketAttachment, error) {
			return nil, dom.ErrTicketNotFound
		},
	}
	rec := httptest.NewRecorder()

	nuevoRouter(uc, sesionSuperAdmin).ServeHTTP(rec,
		peticionConArchivo(t, "/tickets/404/attachments", "file", "a.png", nil))

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteAttachment_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodDelete, "/tickets/attachments/abc", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Empty(t, uc.AdjuntoBorrado)
}

func TestDeleteAttachment_PasaElSolicitanteYSuRol(t *testing.T) {
	casos := []struct {
		nombre string
		sesion func(*gin.Context)
		super  bool
	}{
		{"super admin", sesionSuperAdmin, true},
		{"cliente", sesionNegocio, false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var vistoUsuario uint
			var vistoSuper bool
			uc := &useCaseMock{
				DeleteAttachmentFn: func(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error {
					vistoUsuario, vistoSuper = requesterID, isSuperAdmin
					return nil
				},
			}

			rec := llamar(t, nuevoRouter(uc, tc.sesion), http.MethodDelete, "/tickets/attachments/10", nil)

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, []uint{10}, uc.AdjuntoBorrado)
			assert.Equal(t, uint(42), vistoUsuario)
			assert.Equal(t, tc.super, vistoSuper)
			assert.Equal(t, true, cuerpo(t, rec)["success"])
		})
	}
}

func TestDeleteAttachment_DeOtroUsuario_403(t *testing.T) {
	uc := &useCaseMock{
		DeleteAttachmentFn: func(ctx context.Context, attachmentID uint, requesterID uint, isSuperAdmin bool) error {
			return dom.ErrForbidden
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionNegocio), http.MethodDelete, "/tickets/attachments/10", nil)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestListHistory_IDInvalido_400(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/abc/history", nil)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestListHistory_MapeaElHistorial(t *testing.T) {
	uc := &useCaseMock{
		ListHistoryFn: func(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
			assert.Equal(t, uint(5), ticketID)
			return []entities.TicketStatusHistory{
				{ID: 1, ChangeType: "status", FromStatus: "open", ToStatus: "closed"},
				{ID: 2, ChangeType: "area", FromArea: "soporte", ToArea: "desarrollo"},
			}, nil
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/history", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	data := cuerpo(t, rec)["data"].([]any)
	require.Len(t, data, 2)
	assert.Equal(t, "closed", data[0].(map[string]any)["to_status"])
	assert.Equal(t, "desarrollo", data[1].(map[string]any)["to_area"])
}

func TestListHistory_SinHistorial_DevuelveListaVacia(t *testing.T) {
	uc := &useCaseMock{}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/history", nil)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, []any{}, cuerpo(t, rec)["data"])
}

func TestListHistory_ErrorDelCasoDeUso_500(t *testing.T) {
	uc := &useCaseMock{
		ListHistoryFn: func(ctx context.Context, ticketID uint) ([]entities.TicketStatusHistory, error) {
			return nil, stderrors.New("db caida")
		},
	}

	rec := llamar(t, nuevoRouter(uc, sesionSuperAdmin), http.MethodGet, "/tickets/5/history", nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestSplitCSV(t *testing.T) {
	h := &Handlers{}

	assert.Nil(t, h.splitCSV(""))
	assert.Equal(t, []string{"a"}, h.splitCSV("a"))
	assert.Equal(t, []string{"a", "b"}, h.splitCSV(" a , b "))
	assert.Empty(t, h.splitCSV(" , , "), "solo separadores no produce filtros vacios")
	assert.Equal(t, []string{"a", "b"}, h.splitCSV("a,,b"))
}

func TestParseUint(t *testing.T) {
	h := &Handlers{}

	v, err := h.parseUint("42")
	require.NoError(t, err)
	assert.Equal(t, uint(42), v)

	for _, invalido := range []string{"", "abc", "-1", "4.2"} {
		_, err := h.parseUint(invalido)
		assert.Error(t, err, "valor %q", invalido)
	}
}

func TestRegisterRoutes_MontaLaTablaCompletaBajoTickets(t *testing.T) {
	middleware.Configure(nil, nil, mocks.NewSilentLogger())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	New(&useCaseMock{}, mocks.NewSilentLogger()).RegisterRoutes(r.Group("/api/v1"))

	montadas := map[string]bool{}
	for _, ruta := range r.Routes() {
		montadas[ruta.Method+" "+ruta.Path] = true
	}

	esperadas := []string{
		"GET /api/v1/tickets",
		"POST /api/v1/tickets",
		"GET /api/v1/tickets/:id",
		"PUT /api/v1/tickets/:id",
		"DELETE /api/v1/tickets/:id",
		"PATCH /api/v1/tickets/:id/status",
		"PATCH /api/v1/tickets/:id/assign",
		"PATCH /api/v1/tickets/:id/sprint",
		"PATCH /api/v1/tickets/:id/area",
		"PATCH /api/v1/tickets/:id/escalate",
		"GET /api/v1/tickets/:id/comments",
		"POST /api/v1/tickets/:id/comments",
		"GET /api/v1/tickets/:id/attachments",
		"POST /api/v1/tickets/:id/attachments",
		"DELETE /api/v1/tickets/attachments/:attachment_id",
		"GET /api/v1/tickets/:id/history",
		"GET /api/v1/tickets/categories",
	}
	for _, ruta := range esperadas {
		assert.True(t, montadas[ruta], "falta la ruta %q", ruta)
	}
	assert.Len(t, montadas, len(esperadas))
}

func TestRegisterRoutes_SinTokenNoSeLlegaAlCasoDeUso(t *testing.T) {
	middleware.Configure(nil, nil, mocks.NewSilentLogger())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	uc := &useCaseMock{}
	New(uc, mocks.NewSilentLogger()).RegisterRoutes(r.Group("/api/v1"))

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil))

	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"las rutas se montan detras del JWT, no quedan publicas")
	assert.Nil(t, uc.ListParams)
}
