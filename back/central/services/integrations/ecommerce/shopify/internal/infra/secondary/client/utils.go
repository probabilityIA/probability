package client

import (
	"fmt"
	"strings"
)

// buildURL construye la URL para la API de Shopify.
// Si storeName es una URL completa (comienza con http:// o https://), se usa como base URL directamente.
// Esto permite modo test apuntando a un mock server (ej: http://localhost:3051).
// Si es un dominio de Shopify normal, construye la URL estándar de Shopify Admin API.
func buildURL(storeName, path string) string {
	// URL completa con protocolo -> usar directamente
	if strings.HasPrefix(storeName, "http://") || strings.HasPrefix(storeName, "https://") {
		base := strings.TrimRight(storeName, "/")
		return base + path
	}

	// Si tiene puerto (ej: "back-testing:9090") -> es un mock server, agregar http://
	if strings.Contains(storeName, ":") {
		base := strings.TrimRight(storeName, "/")
		return "http://" + base + path
	}

	// Dominio de Shopify normal
	if !strings.HasSuffix(storeName, ".myshopify.com") {
		storeName = storeName + ".myshopify.com"
	}
	return fmt.Sprintf("https://%s%s", storeName, path)
}

func parseLinkHeader(header string) string {
	if header == "" {
		return ""
	}
	links := strings.Split(header, ",")
	for _, link := range links {
		parts := strings.Split(link, ";")
		if len(parts) < 2 {
			continue
		}
		if strings.Contains(parts[1], `rel="next"`) {
			url := strings.Trim(parts[0], " <>")
			return url
		}
	}
	return ""
}
