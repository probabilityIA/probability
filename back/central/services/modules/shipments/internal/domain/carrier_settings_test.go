package domain

import "testing"

func settings() []CarrierSetting {
	return []CarrierSetting{
		{Code: "TCC", Enabled: false, AllowCOD: true, AllowPrepaid: true},
		{Code: "SERVIENTREGA", Enabled: true, AllowCOD: false, AllowPrepaid: true},
		{Code: "ENVIA", Enabled: true, AllowCOD: true, AllowPrepaid: true},
	}
}

func TestCarrierAllowed_ApagadaNoPasa(t *testing.T) {
	if CarrierAllowed(settings(), "TCC", false) {
		t.Fatal("una transportadora apagada no debe pasar el filtro")
	}
	if CarrierAllowed(settings(), "tcc", true) {
		t.Fatal("el filtro debe ignorar mayusculas y minusculas")
	}
}

func TestCarrierAllowed_ModoDePago(t *testing.T) {
	if CarrierAllowed(settings(), "Servientrega", true) {
		t.Fatal("servientrega no permite contra entrega")
	}
	if !CarrierAllowed(settings(), "Servientrega", false) {
		t.Fatal("servientrega si permite prepago")
	}
}

func TestCarrierAllowed_SinConfigNoFiltra(t *testing.T) {
	if !CarrierAllowed(nil, "TCC", false) {
		t.Fatal("sin configuracion no se filtra nada")
	}
	if !CarrierAllowed(settings(), "MELONN", false) {
		t.Fatal("una transportadora ausente de la configuracion pasa")
	}
	if !CarrierAllowed(settings(), "", true) {
		t.Fatal("una tarifa sin carrier no se descarta")
	}
}
