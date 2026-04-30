package domain

import "testing"

func TestMapEnvioClickEventRealCases(t *testing.T) {
	tests := []struct {
		name         string
		carrier      string
		status       string
		statusStep   string
		statusDetail string
		incidence    bool
		want         ProbabilityShipmentStatus
		wantTable    bool
	}{
		{"envia pendiente captura", "envia", "Pendiente de recolección", "Pendiente de recolección", "Fecha de captura en el sistema", false, StatusPending, true},
		{"envia recolectado", "envia", "En tránsito", "En tránsito", "Fecha de recolección del envio", false, StatusPickedUp, true},
		{"envia salida a ruta", "envia", "En tránsito", "En tránsito", "Salida a ruta", false, StatusInTransit, true},
		{"envia ruta entrega final", "envia", "En tránsito", "En ruta de entrega final", "En ruta de entrega", false, StatusOutForDelivery, true},
		{"envia entregado", "envia", "Entregado", "Entregado", "Envio entregado", false, StatusDelivered, true},
		{"envia cancelado", "envia", "Cancelado", "Cancelado", "Cancelado", false, StatusCancelled, true},
		{"envia incidencia espera ruta", "envia", "Incidencia de Entrega", "En tránsito", "EN ESPERA DE RUTA POBLACION ALEDAÑA", true, StatusOnHold, false},
		{"envia direccion incorrecta", "envia", "Dirección incorrecta/ insuficiente", "En tránsito", "Dirección incorrecta/ insuficiente", true, StatusOnHold, false},
		{"envia no conocen destinatario", "envia", "Direccion incorrecta/ insuficiente", "En transito", "NO CONOCEN DESTINATARIO EN DIRECCION DESTINO", true, StatusOnHold, false},

		{"interrapidisimo envio admitido", "interrapidisimo", "En tránsito", "En tránsito", "Envío Admitido", false, StatusPickedUp, true},
		{"interrapidisimo despachado bodega", "interrapidisimo", "En tránsito", "En tránsito", "Despachado Para Bodega", false, StatusInTransit, true},
		{"interrapidisimo ingresado bodega", "interrapidisimo", "En tránsito", "En tránsito", "Ingresado A Bodega", false, StatusInTransit, true},
		{"interrapidisimo viajando regional", "interrapidisimo", "En tránsito", "En tránsito", "Viajando En Ruta Regional", false, StatusInTransit, true},
		{"interrapidisimo viajando nacional", "interrapidisimo", "En tránsito", "En tránsito", "Viajando En Ruta Nacional", false, StatusInTransit, true},

		{"deprisa guia generada", "deprisa", "Pendiente de Recolección", "Pendiente de Recolección", "GUÍA GENERADA EN ESPERA DE TRANSPORTADORA", false, StatusPending, true},
		{"deprisa cancelado", "deprisa", "Cancelado", "Cancelado", "Cancelado", false, StatusCancelled, true},

		{"servientrega alistamiento cliente", "servientrega", "EN ALISTAMIENTO DEL CLIENTE", "En tránsito", "EN ALISTAMIENTO DEL CLIENTE", false, StatusPending, true},
		{"servientrega cancelado", "servientrega", "Cancelado", "Cancelado", "Cancelado", false, StatusCancelled, true},

		{"carrier desconocido fallback heuristica", "rapidoseguro", "Entregado", "Entregado", "Envio entregado", false, StatusDelivered, false},
		{"detail desconocido fallback heuristica", "envia", "En tránsito", "En tránsito", "Frase rara nueva", false, StatusInTransit, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, fromTable, _ := MapEnvioClickEvent(tt.carrier, tt.status, tt.statusStep, tt.statusDetail, tt.incidence)
			if got != tt.want {
				t.Errorf("MapEnvioClickEvent(%q,%q,%q,%q,%v) = %q, want %q", tt.carrier, tt.status, tt.statusStep, tt.statusDetail, tt.incidence, got, tt.want)
			}
			if fromTable != tt.wantTable {
				t.Errorf("mappedFromTable = %v, want %v", fromTable, tt.wantTable)
			}
		})
	}
}
