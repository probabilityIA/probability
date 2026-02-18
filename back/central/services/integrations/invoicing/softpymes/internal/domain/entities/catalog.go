package entities

// Códigos de tipo de identificación DIAN usados por Softpymes API
// Endpoint: GET /app/integration/identification_types
const (
	IdentTypeCedulaCiudadania   = "13" // Cédula de Ciudadanía
	IdentTypeTarjetaIdentidad   = "12" // Tarjeta de Identidad
	IdentTypeCedulaExtranjeria  = "22" // Cédula de Extranjería
	IdentTypeNIT                = "31" // NIT
	IdentTypeNUIP               = "91" // Número Único de Identificación Personal
)

// Códigos de tipo de documento de Softpymes API
// Endpoint: GET /app/integration/document_types
const (
	DocTypeFacturaVenta      = "FC" // Factura de Venta
	DocTypePedido            = "PE" // Pedido
	DocTypeCotizacion        = "CZ" // Cotización
	DocTypeFacturaServicios  = "FS" // Factura de Servicios Profesionales
	DocTypeRemisiones        = "RM" // Remisiones
)

// Tipos de tercero en Softpymes
const (
	ThirdTypeNatural  = "N" // Persona Natural
	ThirdTypeJuridica = "J" // Persona Jurídica
)
