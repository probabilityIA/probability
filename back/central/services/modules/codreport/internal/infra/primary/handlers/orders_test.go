package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/secamc93/probability/back/central/services/auth/middleware"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ucMock struct {
	listOrdersFn func(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error)
}

func (m *ucMock) Summary(_ context.Context, _ dtos.ReportFilter) (*entities.CodSummary, error) {
	return &entities.CodSummary{}, nil
}

func (m *ucMock) ListOrders(ctx context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
	return m.listOrdersFn(ctx, f)
}

func (m *ucMock) ListCuts(_ context.Context, _ uint, _ bool) ([]entities.PaymentCut, error) {
	return nil, nil
}

func (m *ucMock) ConfirmCut(_ context.Context, _ dtos.ConfirmCutDTO) (*entities.PaymentCut, error) {
	return &entities.PaymentCut{}, nil
}

func (m *ucMock) CarrierConfigs(_ context.Context, _ uint) ([]entities.CarrierConfig, error) {
	return nil, nil
}

func (m *ucMock) SaveCarrierConfig(_ context.Context, _ dtos.SaveCarrierConfigDTO) (*entities.CarrierConfig, error) {
	return &entities.CarrierConfig{}, nil
}

func newOrdersRequest(uc *ucMock, query string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	h := handlers.New(uc, log.New())

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/cod-report/orders"+query, nil)
	c.Set("auth_info", &middleware.AuthInfo{BusinessID: 10})

	h.ListOrders(c)
	return w
}

func TestListOrders_HasGuideFilter(t *testing.T) {
	cases := []struct {
		name      string
		query     string
		wantSet   bool
		wantValue bool
	}{
		{name: "con guia", query: "?has_guide=true", wantSet: true, wantValue: true},
		{name: "sin guia", query: "?has_guide=false", wantSet: true, wantValue: false},
		{name: "sin filtro", query: "", wantSet: false},
		{name: "valor invalido", query: "?has_guide=xyz", wantSet: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var captured dtos.OrdersFilter
			uc := &ucMock{
				listOrdersFn: func(_ context.Context, f dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
					captured = f
					return []entities.CodOrder{}, 0, nil
				},
			}

			w := newOrdersRequest(uc, tc.query)

			if w.Code != http.StatusOK {
				t.Fatalf("status: esperado 200, obtenido %d", w.Code)
			}
			if !tc.wantSet {
				if captured.HasGuide != nil {
					t.Fatalf("HasGuide: esperado nil, obtenido %v", *captured.HasGuide)
				}
				return
			}
			if captured.HasGuide == nil {
				t.Fatalf("HasGuide: esperado %v, obtenido nil", tc.wantValue)
			}
			if *captured.HasGuide != tc.wantValue {
				t.Errorf("HasGuide: esperado %v, obtenido %v", tc.wantValue, *captured.HasGuide)
			}
		})
	}
}

func TestListOrders_ResponseIncludesHasGuide(t *testing.T) {
	uc := &ucMock{
		listOrdersFn: func(_ context.Context, _ dtos.OrdersFilter) ([]entities.CodOrder, int64, error) {
			return []entities.CodOrder{
				{OrderID: "a", OrderNumber: "1001", Carrier: "SERVIENTREGA", HasGuide: true},
				{OrderID: "b", OrderNumber: "1002", Carrier: "SIN TRANSPORTADORA", HasGuide: false},
			}, 2, nil
		},
	}

	w := newOrdersRequest(uc, "")

	if w.Code != http.StatusOK {
		t.Fatalf("status: esperado 200, obtenido %d", w.Code)
	}

	var resp struct {
		Data []struct {
			OrderID  string `json:"order_id"`
			HasGuide bool   `json:"has_guide"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error al deserializar respuesta: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("data: esperado 2 filas, obtenido %d", len(resp.Data))
	}
	if !resp.Data[0].HasGuide {
		t.Errorf("fila 0: esperado has_guide true")
	}
	if resp.Data[1].HasGuide {
		t.Errorf("fila 1: esperado has_guide false")
	}
}
