package usecases

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
)

// Carrier info for Colombian transport companies
type carrierInfo struct {
	ID            int64
	Name          string
	TrackPrefix   string
	Products      []productInfo
	BaseDelivery  int // base delivery days
}

type productInfo struct {
	ID   int64
	Name string
}

var carriers = []carrierInfo{
	{
		ID:           1,
		Name:         "Servientrega",
		TrackPrefix:  "SRV",
		BaseDelivery: 2,
		Products: []productInfo{
			{ID: 101, Name: "Mercancia Premier"},
			{ID: 102, Name: "Mercancia Estandar"},
		},
	},
	{
		ID:           2,
		Name:         "Coordinadora",
		TrackPrefix:  "CRD",
		BaseDelivery: 3,
		Products: []productInfo{
			{ID: 201, Name: "Paqueteria Express"},
			{ID: 202, Name: "Paqueteria Estandar"},
		},
	},
	{
		ID:           3,
		Name:         "Envia",
		TrackPrefix:  "ENV",
		BaseDelivery: 2,
		Products: []productInfo{
			{ID: 301, Name: "Express Nacional"},
			{ID: 302, Name: "Economico"},
		},
	},
	{
		ID:           4,
		Name:         "InterRapidisimo",
		TrackPrefix:  "IRP",
		BaseDelivery: 1,
		Products: []productInfo{
			{ID: 401, Name: "Hoy Mismo"},
			{ID: 402, Name: "Dia Siguiente"},
		},
	},
	{
		ID:           5,
		Name:         "TCC",
		TrackPrefix:  "TCC",
		BaseDelivery: 3,
		Products: []productInfo{
			{ID: 501, Name: "Mensajeria Express"},
			{ID: 502, Name: "Carga Ligera"},
		},
	},
}

// Valid DANE codes - supports both 5-char (legacy) and 8-char (EnvioClick format) codes.
// EnvioClick API uses 8-digit codes; the 8-char format is the canonical format.
var validDaneCodes = map[string]string{
	// 8-char codes (EnvioClick format) - departmental capitals
	"11001000": "Bogota D.C.",
	"05001000": "Medellin",
	"76001000": "Cali",
	"08001000": "Barranquilla",
	"13001000": "Cartagena",
	"68001000": "Bucaramanga",
	"54001000": "Cucuta",
	"73001000": "Ibague",
	"17001000": "Manizales",
	"66001000": "Pereira",
	"41001000": "Neiva",
	"15001000": "Tunja",
	"52001000": "Pasto",
	"23001000": "Monteria",
	"47001000": "Santa Marta",
	"50001000": "Villavicencio",
	"19001000": "Popayan",
	"63001000": "Armenia",
	"44001000": "Riohacha",
	"27001000": "Quibdo",
	"70001000": "Sincelejo",
	"18001000": "Florencia",
	"85001000": "Yopal",
	"20001000": "Valledupar",
	"86001000": "Mocoa",
	"88001000": "San Andres",
	"91001000": "Leticia",
	"94001000": "Inirida",
	"95001000": "San Jose del Guaviare",
	"97001000": "Mitu",
	"99001000": "Puerto Carreno",
	// Common municipalities
	"05002000": "Abejorral",
	"05004000": "Abriaqui",
	"05021000": "Alejandria",
	"05030000": "Amaga",
	"05031000": "Amalfi",
	"05034000": "Andes",
	"05036000": "Angelopolis",
	"05091000": "Bello",
	"05093000": "Betania",
	"05101000": "Briceno",
	"05107000": "Buritica",
	"05113000": "Caceres",
	"05120000": "Caldas",
	"05125000": "Campamento",
	"05129000": "Canamomo",
	"05134000": "Caracoli",
	"05138000": "Carepa",
	"05142000": "Carmen de Viboral",
	"05145000": "Carolina",
	"05147000": "Caucasia",
	"05150000": "Cisneros",
	"05154000": "Ciudad Bolivar",
	"05172000": "Concordia",
	"05190000": "Copacabana",
	"05197000": "Dabeiba",
	"05206000": "Don Matias",
	"05209000": "Ebejico",
	"05212000": "El Bagre",
	"05234000": "Entrerrios",
	"05237000": "Envigado",
	"05240000": "Fredonia",
	"05250000": "Frontino",
	"05264000": "Giraldo",
	"05266000": "Girardota",
	"05282000": "Gomez Plata",
	"05284000": "Granada",
	"05306000": "Guarne",
	"05308000": "Guatape",
	"05310000": "Heliconia",
	"05313000": "Hispania",
	"05315000": "Itagui",
	"05318000": "Ituango",
	"05321000": "Jardin",
	"05347000": "La Ceja",
	"05353000": "La Estrella",
	"05360000": "La Pintada",
	"05364000": "La Union",
	"05368000": "Liborina",
	"05376000": "Maceo",
	"05380000": "Marinilla",
	"05390000": "Montebello",
	"05400000": "Murindo",
	"05411000": "Mutata",
	"05425000": "Narino",
	"05440000": "Nechi",
	"05467000": "Olaya",
	"05475000": "Penol",
	"05480000": "Peque",
	"05490000": "Pueblorrico",
	"05495000": "Puerto Berrio",
	"05501000": "Puerto Nare",
	"05541000": "Remedios",
	"05543000": "Retiro",
	"05576000": "Sabanalarga",
	"05579000": "Sabaneta",
	"05585000": "Salgar",
	"05591000": "San Andres de Cuerquia",
	"05604000": "San Carlos",
	"05607000": "San Francisco",
	"05615000": "San Luis",
	"05628000": "San Pedro de los Milagros",
	"05631000": "San Rafael",
	"05642000": "Santa Barbara",
	"05647000": "Santa Rosa de Osos",
	"05649000": "Santo Domingo",
	"05652000": "El Santuario",
	"05658000": "Segovia",
	"05664000": "Sonson",
	"05665000": "Sopetran",
	"05670000": "Tamesis",
	"05674000": "Taraza",
	"05679000": "Tarso",
	"05686000": "Titiribi",
	"05690000": "Toledo",
	"05697000": "Turbo",
	"05736000": "Urrao",
	"05745000": "Valdivia",
	"05756000": "Valparaiso",
	"05761000": "Vegachi",
	"05764000": "Venecia",
	"05789000": "Yali",
	"05790000": "Yanacona",
	"05792000": "Yolombo",
	"05809000": "Zaragoza",
	// Cundinamarca municipalities
	"25001000": "Agua de Dios",
	"25019000": "Alban",
	"25035000": "Anapoima",
	"25040000": "Anolaima",
	"25053000": "Arbelaez",
	"25086000": "Beltran",
	"25095000": "Bituima",
	"25099000": "Bojaca",
	"25120000": "Cabrera",
	"25123000": "Cachipay",
	"25126000": "Cajica",
	"25148000": "Caparrapi",
	"25151000": "Caqueza",
	"25154000": "Carmen de Carupa",
	"25168000": "Chaguani",
	"25175000": "Chia",
	"25178000": "Chipaque",
	"25181000": "Choachi",
	"25183000": "Choconta",
	"25200000": "Cogua",
	"25214000": "Cota",
	"25224000": "Cucunuba",
	"25245000": "El Colegio",
	"25258000": "El Penon",
	"25260000": "El Rosal",
	"25269000": "Facatativa",
	"25279000": "Fomeque",
	"25281000": "Fosca",
	"25286000": "Funza",
	"25288000": "Fusagasuga",
	"25290000": "Gachala",
	"25293000": "Gachancipa",
	"25295000": "Gacheta",
	"25297000": "Gama",
	"25307000": "Girardot",
	"25312000": "Granada",
	"25317000": "Guacheta",
	"25320000": "Guaduas",
	"25322000": "Guasca",
	"25324000": "Guataqui",
	"25326000": "Guatavita",
	"25328000": "Guayabal de Siquima",
	"25335000": "Guayabetal",
	"25339000": "Gutierrez",
	"25368000": "Jerusalen",
	"25372000": "Junin",
	"25377000": "La Calera",
	"25386000": "La Mesa",
	"25394000": "La Palma",
	"25398000": "La Pena",
	"25402000": "La Vega",
	"25407000": "Lenguazaque",
	"25426000": "Macheta",
	"25430000": "Madrid",
	"25436000": "Manta",
	"25438000": "Medina",
	"25473000": "Mosquera",
	"25483000": "Nemocon",
	"25486000": "Nilo",
	"25488000": "Nimaima",
	"25489000": "Nocaima",
	"25491000": "Venecia",
	"25506000": "Pacho",
	"25513000": "Paime",
	"25518000": "Pandi",
	"25524000": "Paratebueno",
	"25530000": "Pasca",
	"25535000": "Puerto Salgar",
	"25572000": "Puli",
	"25580000": "Quebradanegra",
	"25592000": "Quetame",
	"25594000": "Quipile",
	"25596000": "Apulo",
	"25612000": "Ricaurte",
	"25645000": "San Antonio del Tequendama",
	"25649000": "San Bernardo",
	"25653000": "San Cayetano",
	"25658000": "San Francisco",
	"25662000": "San Juan de Rio Seco",
	"25718000": "Sasaima",
	"25736000": "Sesquile",
	"25740000": "Sibate",
	"25743000": "Silvania",
	"25745000": "Simijaca",
	"25754000": "Soacha",
	"25758000": "Sopo",
	"25769000": "Subachoque",
	"25772000": "Suesca",
	"25777000": "Supata",
	"25779000": "Susa",
	"25781000": "Sutatausa",
	"25785000": "Tabio",
	"25793000": "Tausa",
	"25797000": "Tena",
	"25799000": "Tenjo",
	"25805000": "Tibacuy",
	"25807000": "Tibirita",
	"25815000": "Tocaima",
	"25817000": "Tocancipa",
	"25823000": "Topaipi",
	"25839000": "Ubala",
	"25841000": "Ubaque",
	"25843000": "Villa de San Diego de Ubate",
	"25845000": "Une",
	"25851000": "Utica",
	"25862000": "Vergara",
	"25867000": "Viani",
	"25871000": "Villagomez",
	"25873000": "Villapinzon",
	"25875000": "Villeta",
	"25878000": "Viota",
	"25885000": "Yacopi",
	"25898000": "Zipacon",
	"25899000": "Zipaquira",
	// Valle del Cauca municipalities
	"76020000": "Alcala",
	"76036000": "Andalucia",
	"76041000": "Ansermanuevo",
	"76054000": "Argelia",
	"76100000": "Bolivar",
	"76109000": "Buenaventura",
	"76111000": "Guadalajara de Buga",
	"76113000": "Bugalagrande",
	"76122000": "Caicedonia",
	"76126000": "Calima",
	"76130000": "Candelaria",
	"76147000": "Cartago",
	"76233000": "Dagua",
	"76243000": "El Aguila",
	"76246000": "El Cairo",
	"76248000": "El Cerrito",
	"76250000": "El Dovio",
	"76275000": "Florida",
	"76306000": "Ginebra",
	"76318000": "Guacari",
	"76364000": "Jamundi",
	"76377000": "La Cumbre",
	"76400000": "La Union",
	"76403000": "La Victoria",
	"76497000": "Obando",
	"76520000": "Palmira",
	"76563000": "Pradera",
	"76606000": "Restrepo",
	"76616000": "Riofrio",
	"76622000": "Roldanillo",
	"76670000": "San Pedro",
	"76736000": "Sevilla",
	"76823000": "Toro",
	"76828000": "Trujillo",
	"76834000": "Tulua",
	"76845000": "Ulloa",
	"76863000": "Versalles",
	"76869000": "Vijes",
	"76890000": "Yotoco",
	"76892000": "Yumbo",
	"76895000": "Zarzal",
	// 5-char codes (legacy fallback)
	"11001": "Bogota D.C.",
	"05001": "Medellin",
	"76001": "Cali",
	"08001": "Barranquilla",
	"13001": "Cartagena",
	"68001": "Bucaramanga",
	"54001": "Cucuta",
	"73001": "Ibague",
	"17001": "Manizales",
	"66001": "Pereira",
	"41001": "Neiva",
	"15001": "Tunja",
	"52001": "Pasto",
	"23001": "Monteria",
	"47001": "Santa Marta",
	"50001": "Villavicencio",
	"19001": "Popayan",
	"63001": "Armenia",
	"44001": "Riohacha",
	"27001": "Quibdo",
	"70001": "Sincelejo",
	"18001": "Florencia",
	"85001": "Yopal",
	"20001": "Valledupar",
	"86001": "Mocoa",
	"88001": "San Andres",
	"91001": "Leticia",
	"94001": "Inirida",
	"95001": "San Jose del Guaviare",
	"97001": "Mitu",
	"99001": "Puerto Carreno",
}

// IsValidDaneCode validates a DANE code.
// Accepts 8-char codes (EnvioClick format) and 5-char legacy codes.
// For 8-char codes, any numeric code with a recognized 5-char prefix is accepted.
func IsValidDaneCode(code string) bool {
	if code == "" {
		return false
	}
	// Check exact match first
	if _, ok := validDaneCodes[code]; ok {
		return true
	}
	// For 8-char numeric codes, check if the 5-char prefix is a known department
	if len(code) == 8 {
		prefix := code[:5]
		if _, ok := validDaneCodes[prefix]; ok {
			return true
		}
		// Accept any 8-char numeric code (EnvioClick format) - all are valid for mock
		allDigits := true
		for _, c := range code {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		return allDigits
	}
	return false
}

// GetCityName returns the city name for a DANE code
func GetCityName(code string) string {
	if name, ok := validDaneCodes[code]; ok {
		return name
	}
	// For 8-char codes not in the map, try the 5-char prefix
	if len(code) == 8 {
		if name, ok := validDaneCodes[code[:5]]; ok {
			return name
		}
	}
	return "Desconocido"
}

// GenerateRates generates 3-5 mock rates with realistic COP prices
func GenerateRates(req domain.QuoteRequest) []domain.Rate {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	totalWeight := 0.0
	totalVolWeight := 0.0
	for _, pkg := range req.Packages {
		totalWeight += pkg.Weight
		totalVolWeight += (pkg.Height * pkg.Width * pkg.Length) / 5000 // volumetric
	}

	billableWeight := math.Max(totalWeight, totalVolWeight)

	// Pick 3-5 random carrier+product combos
	numRates := 3 + rng.Intn(3) // 3 to 5
	used := make(map[string]bool)
	rates := make([]domain.Rate, 0, numRates)
	idCounter := int64(rng.Intn(9000) + 1000)

	for i := 0; i < numRates && i < len(carriers)*2; i++ {
		c := carriers[rng.Intn(len(carriers))]
		p := c.Products[rng.Intn(len(c.Products))]
		key := fmt.Sprintf("%d-%d", c.ID, p.ID)
		if used[key] {
			continue
		}
		used[key] = true

		flete := calculateFlete(billableWeight, req.ContentValue, rng)
		deliveryDays := c.BaseDelivery + rng.Intn(3) // base + 0-2 extra days

		idCounter++
		rates = append(rates, domain.Rate{
			IDRate:        idCounter,
			IDProduct:     p.ID,
			Product:       p.Name,
			IDCarrier:     c.ID,
			Carrier:       c.Name,
			Flete:         flete,
			DeliveryDays:  deliveryDays,
			QuotationType: "standard",
		})
	}

	return rates
}

// calculateFlete calculates a realistic shipping cost in COP
func calculateFlete(billableWeight, contentValue float64, rng *rand.Rand) float64 {
	// Base: 8000 COP
	base := 8000.0
	// Weight component: 2500 COP per kg
	weightCost := billableWeight * 2500
	// Insurance component: 1.5% of declared value
	insuranceCost := contentValue * 0.015
	// Random variation +/- 15%
	variation := 1.0 + (rng.Float64()*0.3 - 0.15)

	total := (base + weightCost + insuranceCost) * variation
	// Round to nearest 100 COP
	return math.Round(total/100) * 100
}

// GenerateTrackingNumber generates a carrier-specific tracking number
func GenerateTrackingNumber(carrier carrierInfo, rng *rand.Rand) string {
	num := rng.Intn(900000000) + 100000000 // 9 digit number
	return fmt.Sprintf("%s-%d", carrier.TrackPrefix, num)
}

// GenerateTrackingHistory generates 2-4 realistic tracking events
func GenerateTrackingHistory(carrier string, createdAt time.Time) []domain.TrackHistory {
	cities := []string{"Bogota D.C.", "Medellin", "Cali", "Barranquilla", "Bucaramanga"}
	rng := rand.New(rand.NewSource(createdAt.UnixNano()))

	events := []domain.TrackHistory{
		{
			Date:        createdAt.Format("2006-01-02 15:04:05"),
			Status:      "Recolectado",
			Description: "Paquete recolectado en punto de origen",
			Location:    cities[rng.Intn(len(cities))],
		},
		{
			Date:        createdAt.Add(2 * time.Hour).Format("2006-01-02 15:04:05"),
			Status:      "En centro de distribucion",
			Description: fmt.Sprintf("Paquete recibido en centro de distribucion %s", carrier),
			Location:    cities[rng.Intn(len(cities))],
		},
	}

	// Optionally add more events
	if rng.Intn(2) == 0 {
		events = append(events, domain.TrackHistory{
			Date:        createdAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Status:      "En transito",
			Description: "Paquete en ruta hacia destino",
			Location:    cities[rng.Intn(len(cities))],
		})
	}
	if rng.Intn(3) == 0 {
		events = append(events, domain.TrackHistory{
			Date:        createdAt.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
			Status:      "En reparto",
			Description: "Paquete en vehiculo de reparto",
			Location:    cities[rng.Intn(len(cities))],
		})
	}

	return events
}

// FindCarrierByID finds a carrier by its ID
func FindCarrierByID(id int64) *carrierInfo {
	for _, c := range carriers {
		if c.ID == id {
			return &c
		}
	}
	return nil
}

// FindCarrierByRateID finds a carrier by a rate's carrier ID
func FindCarrierByName(name string) *carrierInfo {
	for _, c := range carriers {
		if c.Name == name {
			return &c
		}
	}
	return nil
}
