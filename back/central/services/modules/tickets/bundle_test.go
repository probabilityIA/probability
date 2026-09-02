package tickets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_MontaElModuloCompletoSobreElRouter(t *testing.T) {
	middleware.Configure(nil, nil, mocks.NewSilentLogger())
	gin.SetMode(gin.TestMode)

	r := gin.New()
	New(r.Group("/api/v1"), testkit.NewDB(t), mocks.NewSilentLogger(), nil)

	montadas := map[string]bool{}
	for _, ruta := range r.Routes() {
		montadas[ruta.Method+" "+ruta.Path] = true
	}

	require.Len(t, montadas, 16, "el modulo expone 16 rutas")
	assert.True(t, montadas["GET /api/v1/tickets"])
	assert.True(t, montadas["POST /api/v1/tickets"])
	assert.True(t, montadas["PATCH /api/v1/tickets/:id/sprint"])
	assert.True(t, montadas["DELETE /api/v1/tickets/attachments/:attachment_id"])
}

func TestNew_LasRutasQuedanDetrasDelJWT(t *testing.T) {
	middleware.Configure(nil, nil, mocks.NewSilentLogger())
	gin.SetMode(gin.TestMode)

	r := gin.New()
	New(r.Group("/api/v1"), testkit.NewDB(t), mocks.NewSilentLogger(), nil)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil))

	assert.Equal(t, http.StatusUnauthorized, rec.Code,
		"sin token no se toca la base ni el caso de uso")
}
