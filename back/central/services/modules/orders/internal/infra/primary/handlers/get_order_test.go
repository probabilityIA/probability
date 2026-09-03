package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
)

type orderUseCaseStub struct {
	orden *dtos.OrderResponse
	err   error
}

func (s *orderUseCaseStub) GetOrderByID(_ context.Context, _ string) (*dtos.OrderResponse, error) {
	return s.orden, s.err
}
func (s *orderUseCaseStub) GetOrderRaw(context.Context, string) (*dtos.OrderRawResponse, error) {
	return nil, nil
}
func (s *orderUseCaseStub) GetOrderHistory(context.Context, string) ([]dtos.OrderHistoryResponse, error) {
	return nil, nil
}
func (s *orderUseCaseStub) GetOrderNotifications(context.Context, string, uint) (*dtos.OrderNotificationsResponse, error) {
	return nil, nil
}
func (s *orderUseCaseStub) ListOrders(context.Context, int, int, map[string]interface{}) (*dtos.OrdersListResponse, error) {
	return nil, nil
}
func (s *orderUseCaseStub) UpdateOrder(context.Context, string, *dtos.UpdateOrderRequest) (*dtos.OrderResponse, error) {
	return nil, nil
}
func (s *orderUseCaseStub) DeleteOrder(context.Context, string) error { return nil }

func pedirOrden(t *testing.T, stub *orderUseCaseStub, businessIDToken interface{}) (int, map[string]interface{}) {
	t.Helper()
	gin.SetMode(gin.TestMode)

	h := &Handlers{orderCRUD: stub}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	if businessIDToken != nil {
		c.Set("business_id", businessIDToken)
	}

	h.GetOrderByID(c)

	var body map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &body)
	return w.Code, body
}

func ordenDe(businessID uint) *dtos.OrderResponse {
	return &dtos.OrderResponse{ID: "abc", BusinessID: &businessID}
}

func TestGetOrderByID_OtroNegocioDaNoAutorizado(t *testing.T) {
	code, body := pedirOrden(t, &orderUseCaseStub{orden: ordenDe(46)}, uint(26))

	if code != http.StatusForbidden {
		t.Fatalf("se esperaba 403, llego %d", code)
	}
	if body["message"] != "No autorizado" {
		t.Fatalf("mensaje inesperado: %v", body["message"])
	}
}

func TestGetOrderByID_MismoNegocioPasa(t *testing.T) {
	code, body := pedirOrden(t, &orderUseCaseStub{orden: ordenDe(26)}, uint(26))

	if code != http.StatusOK || body["success"] != true {
		t.Fatalf("se esperaba 200 exitoso, llego %d %v", code, body)
	}
}

func TestGetOrderByID_SuperAdminVeCualquierOrden(t *testing.T) {
	code, _ := pedirOrden(t, &orderUseCaseStub{orden: ordenDe(46)}, uint(0))
	if code != http.StatusOK {
		t.Fatalf("el super admin debe poder ver cualquier orden, llego %d", code)
	}

	code, _ = pedirOrden(t, &orderUseCaseStub{orden: ordenDe(46)}, nil)
	if code != http.StatusOK {
		t.Fatalf("sin business_id en el contexto se trata como super admin, llego %d", code)
	}
}

func TestGetOrderByID_SinNegocioEnLaOrdenNoSeEntrega(t *testing.T) {
	code, _ := pedirOrden(t, &orderUseCaseStub{orden: &dtos.OrderResponse{ID: "abc"}}, uint(26))
	if code != http.StatusForbidden {
		t.Fatalf("una orden sin business_id no debe entregarse a un negocio, llego %d", code)
	}
}

func TestGetOrderByID_InexistenteDa404AunqueElErrorVengaEnvuelto(t *testing.T) {
	err := fmt.Errorf("error getting order: %w", errors.New("order not found"))
	code, body := pedirOrden(t, &orderUseCaseStub{err: err}, uint(26))

	if code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, llego %d", code)
	}
	if body["message"] != "Orden no encontrada" {
		t.Fatalf("mensaje inesperado: %v", body["message"])
	}
}
