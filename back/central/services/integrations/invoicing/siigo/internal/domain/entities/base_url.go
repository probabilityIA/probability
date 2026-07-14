package entities

import "strings"

// ResolveSiigoBaseURL decide contra que URL de Siigo se habla.
// Precedencia (la misma que usa Probar Conexion):
//  1. base_url_test del integration_type, si la integracion esta en modo pruebas.
//  2. api_url: override opcional guardado como credencial de la integracion.
//  3. base_url del integration_type (produccion).
//
// Si devuelve "", el llamador debe cortar con un error claro: sin base URL el
// cliente HTTP arma rutas relativas ("/auth") y falla con
// "unsupported protocol scheme", que no le dice nada al usuario.
func ResolveSiigoBaseURL(isTesting bool, baseURLTest, apiURLOverride, baseURL string) string {
	test := strings.TrimSpace(baseURLTest)
	override := strings.TrimSpace(apiURLOverride)
	prod := strings.TrimSpace(baseURL)

	if isTesting && test != "" {
		return test
	}
	if override != "" {
		return override
	}
	return prod
}
