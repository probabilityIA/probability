package domain

import "testing"

func TestMapEnvioClickStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusStep string
		incidence  bool
		want       string
	}{
		{
			name:       "pendiente de recoleccion",
			statusStep: "Pendiente de Recolecci\u00f3n",
			incidence:  false,
			want:       "pending",
		},
		{
			name:       "pendiente",
			statusStep: "Pendiente",
			incidence:  false,
			want:       "pending",
		},
		{
			name:       "en transito con acento",
			statusStep: "En tr\u00e1nsito",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "en transito mayuscula con acento",
			statusStep: "En Tr\u00e1nsito",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "en transito sin acento",
			statusStep: "En Transito",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "en transito minuscula sin acento",
			statusStep: "En transito",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "envio recolectado con acento",
			statusStep: "Env\u00edo Recolectado",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "envio recolectado sin acento",
			statusStep: "Envio Recolectado",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "en distribucion con acento",
			statusStep: "En Distribuci\u00f3n",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "en distribucion sin acento",
			statusStep: "En distribucion",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "entregado masculino",
			statusStep: "Entregado",
			incidence:  false,
			want:       "delivered",
		},
		{
			name:       "entregada femenino",
			statusStep: "Entregada",
			incidence:  false,
			want:       "delivered",
		},
		{
			name:       "incidencia por flag con status entregado",
			statusStep: "Entregado",
			incidence:  true,
			want:       "failed",
		},
		{
			name:       "incidencia por flag con status pendiente",
			statusStep: "Pendiente",
			incidence:  true,
			want:       "failed",
		},
		{
			name:       "incidencia por status step",
			statusStep: "Incidencia",
			incidence:  false,
			want:       "failed",
		},
		{
			name:       "novedad",
			statusStep: "Novedad",
			incidence:  false,
			want:       "failed",
		},
		{
			name:       "no entregado minuscula",
			statusStep: "No entregado",
			incidence:  false,
			want:       "failed",
		},
		{
			name:       "no entregado mayuscula",
			statusStep: "No Entregado",
			incidence:  false,
			want:       "failed",
		},
		{
			name:       "status desconocido retorna in_transit por defecto",
			statusStep: "Estado Desconocido XYZ",
			incidence:  false,
			want:       "in_transit",
		},
		{
			name:       "status vacio retorna in_transit por defecto",
			statusStep: "",
			incidence:  false,
			want:       "in_transit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapEnvioClickStatus(tt.statusStep, tt.incidence)
			if got != tt.want {
				t.Errorf("MapEnvioClickStatus(%q, %v) = %q, want %q", tt.statusStep, tt.incidence, got, tt.want)
			}
		})
	}
}
