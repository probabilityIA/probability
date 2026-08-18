package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
)

func newActorContext(setUser bool, userID uint, email string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("POST", "/", nil)
	if setUser {
		c.Set("user_id", userID)
		c.Set("user_email", email)
	}
	return c
}

func newActorHandlers(repo *mocks.RepositoryMock) *Handlers {
	return &Handlers{uc: usecases.New(repo, nil, nil)}
}

func TestResolveActor_SinUsuarioEnElToken(t *testing.T) {
	h := newActorHandlers(&mocks.RepositoryMock{})

	a := h.resolveActor(newActorContext(false, 0, ""))

	if a.ID != nil {
		t.Errorf("sin usuario en el token no debe haber actor, se obtuvo %v", *a.ID)
	}
	if a.Name != "" {
		t.Errorf("sin usuario en el token no debe haber nombre, se obtuvo %q", a.Name)
	}
}

func TestResolveActor_UsaElNombreDelUsuario(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetUserDisplayNameFn: func(ctx context.Context, userID uint) string {
			if userID != 7 {
				t.Fatalf("se esperaba el user_id 7 del token, llego %d", userID)
			}
			return "Sebastian Camacho"
		},
	}
	h := newActorHandlers(repo)

	a := h.resolveActor(newActorContext(true, 7, "sebas@probability.com"))

	if a.ID == nil || *a.ID != 7 {
		t.Fatalf("se esperaba el actor 7, se obtuvo %v", a.ID)
	}
	if a.Name != "Sebastian Camacho" {
		t.Errorf("nombre esperado del usuario, obtenido: %q", a.Name)
	}
}

func TestResolveActor_CaeAlEmailCuandoNoHayNombre(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetUserDisplayNameFn: func(ctx context.Context, userID uint) string { return "" },
	}
	h := newActorHandlers(repo)

	a := h.resolveActor(newActorContext(true, 9, "sebas@probability.com"))

	if a.Name != "sebas@probability.com" {
		t.Errorf("sin nombre debe quedar el email como identificacion, obtenido: %q", a.Name)
	}
}
