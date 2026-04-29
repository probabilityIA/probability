package entities

type CarrierCatalog struct {
	Code string
	Name string
}

var DefaultCarriers = []CarrierCatalog{
	{Code: "servientrega", Name: "Servientrega"},
	{Code: "interrapidisimo", Name: "Interrapidisimo"},
	{Code: "coordinadora", Name: "Coordinadora"},
	{Code: "envia", Name: "Envia"},
	{Code: "tcc", Name: "TCC"},
	{Code: "deprisa", Name: "Deprisa"},
	{Code: "99minutos", Name: "99Minutos"},
	{Code: "mipaquete", Name: "MiPaquete"},
	{Code: "enviame", Name: "Enviame"},
}
