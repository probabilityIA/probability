package usecasemessaging

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidatePhoneNumber valida que el número de teléfono tenga el formato correcto
// Acepta formatos: +[código_país][número], 00[código_país][número], [código_país][número]
// Valida que tenga un código de país válido y formato internacional
func ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return fmt.Errorf("número de teléfono vacío")
	}

	// Limpiar espacios y caracteres especiales
	cleanPhone := strings.TrimSpace(phone)

	// Remover caracteres no numéricos excepto + al inicio
	if cleanPhone[0] == '+' {
		cleanPhone = "+" + regexp.MustCompile(`\D`).ReplaceAllString(cleanPhone[1:], "")
	} else {
		cleanPhone = regexp.MustCompile(`\D`).ReplaceAllString(cleanPhone, "")
	}

	// Extraer código de país y número
	var countryCode, phoneDigits string
	var minDigits, maxDigits int

	if strings.HasPrefix(cleanPhone, "+") {
		// Formato: +[código_país][número]
		cleanPhone = cleanPhone[1:] // Remover el +
		countryCode, phoneDigits = extractCountryCodeAndNumber(cleanPhone)
	} else if strings.HasPrefix(cleanPhone, "00") {
		// Formato: 00[código_país][número]
		cleanPhone = cleanPhone[2:] // Remover "00"
		countryCode, phoneDigits = extractCountryCodeAndNumber(cleanPhone)
	} else {
		// Formato: [código_país][número] (sin prefijo)
		countryCode, phoneDigits = extractCountryCodeAndNumber(cleanPhone)
	}

	if countryCode == "" {
		return fmt.Errorf("número debe incluir un código de país válido (ej: +1, +57, +34)")
	}

	// Obtener longitud esperada para el código de país
	minDigits, maxDigits = getCountryCodeLength(countryCode)
	if minDigits == 0 {
		return fmt.Errorf("código de país '%s' no reconocido", countryCode)
	}

	// Validar longitud del número según el país
	if len(phoneDigits) < minDigits || len(phoneDigits) > maxDigits {
		return fmt.Errorf("número debe tener entre %d y %d dígitos para el código de país %s, tiene %d",
			minDigits, maxDigits, countryCode, len(phoneDigits))
	}

	// Validar que todos los caracteres sean dígitos
	if !regexp.MustCompile(`^\d+$`).MatchString(phoneDigits) {
		return fmt.Errorf("número debe contener solo dígitos después del código de país")
	}

	return nil
}

// extractCountryCodeAndNumber extrae el código de país y el número del teléfono
func extractCountryCodeAndNumber(phone string) (string, string) {
	// Lista de códigos de país comunes ordenados por longitud (de mayor a menor)
	countryCodes := []string{
		// Códigos de 3 dígitos
		"1242", "1246", "1264", "1268", "1284", "1340", "1345", "1441", "1473", "1649",
		"1664", "1670", "1671", "1684", "1721", "1758", "1767", "1784", "1787", "1809",
		"1829", "1849", "1868", "1869", "1876", "1939",

		// Códigos de 2 dígitos
		"20", "27", "30", "31", "32", "33", "34", "36", "39", "40", "41", "43", "44", "45", "46", "47", "48", "49",
		"51", "52", "53", "54", "55", "56", "57", "58", "60", "61", "62", "63", "64", "65", "66", "81", "82", "84", "86", "90", "91", "92", "93", "94", "95", "98",

		// Códigos de 1 dígito
		"1", "7",
	}

	for _, code := range countryCodes {
		if strings.HasPrefix(phone, code) {
			return code, phone[len(code):]
		}
	}

	return "", phone
}

// getCountryCodeLength retorna la longitud mínima y máxima de números para códigos de país comunes
func getCountryCodeLength(countryCode string) (int, int) {
	// Mapeo de códigos de país a longitud de número
	lengths := map[string][]int{
		// América del Norte
		"1": {10, 10}, // USA, Canadá

		// América del Sur
		"51": {8, 9},   // Perú
		"52": {10, 10}, // México
		"53": {8, 8},   // Cuba
		"54": {10, 11}, // Argentina
		"55": {10, 11}, // Brasil
		"56": {8, 9},   // Chile
		"57": {10, 10}, // Colombia
		"58": {10, 11}, // Venezuela

		// Europa
		"30": {10, 10}, // Grecia
		"31": {9, 9},   // Países Bajos
		"32": {8, 9},   // Bélgica
		"33": {9, 9},   // Francia
		"34": {9, 9},   // España
		"36": {8, 9},   // Hungría
		"39": {9, 11},  // Italia
		"40": {9, 9},   // Rumania
		"41": {9, 9},   // Suiza
		"43": {10, 13}, // Austria
		"44": {10, 10}, // Reino Unido
		"45": {8, 8},   // Dinamarca
		"46": {8, 9},   // Suecia
		"47": {8, 8},   // Noruega
		"48": {9, 9},   // Polonia
		"49": {10, 12}, // Alemania

		// Asia
		"60": {8, 10},  // Malasia
		"61": {9, 9},   // Australia
		"62": {8, 12},  // Indonesia
		"63": {10, 10}, // Filipinas
		"64": {8, 9},   // Nueva Zelanda
		"65": {8, 8},   // Singapur
		"66": {8, 9},   // Tailandia
		"81": {10, 11}, // Japón
		"82": {9, 10},  // Corea del Sur
		"84": {9, 10},  // Vietnam
		"86": {11, 11}, // China
		"90": {10, 10}, // Turquía
		"91": {10, 10}, // India
		"92": {10, 10}, // Pakistán
		"93": {8, 9},   // Afganistán
		"94": {9, 9},   // Sri Lanka
		"95": {8, 10},  // Myanmar
		"98": {10, 10}, // Irán

		// África
		"20": {9, 10}, // Egipto
		"27": {9, 9},  // Sudáfrica

		// Rusia y países vecinos
		"7": {10, 10}, // Rusia, Kazajistán
	}

	if lengths, exists := lengths[countryCode]; exists {
		return lengths[0], lengths[1]
	}

	// Valores por defecto para códigos no especificados
	return 7, 15
}
