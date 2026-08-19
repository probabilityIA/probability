package client

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const siigoNameMaxLen = 100

var siigoNameTransformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func sanitizeSiigoName(raw string) string {
	folded, _, err := transform.String(siigoNameTransformer, raw)
	if err != nil {
		folded = raw
	}

	var b strings.Builder
	for _, r := range folded {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ':
			b.WriteRune(r)
		default:
			b.WriteRune(' ')
		}
	}

	cleaned := strings.Join(strings.Fields(b.String()), " ")
	if len(cleaned) > siigoNameMaxLen {
		cleaned = strings.TrimSpace(cleaned[:siigoNameMaxLen])
	}

	return cleaned
}

func buildSiigoCustomerName(raw, personType string) []string {
	parts := strings.Fields(sanitizeSiigoName(raw))

	if strings.EqualFold(personType, "company") {
		joined := strings.Join(parts, " ")
		if joined == "" {
			joined = "Sin Nombre"
		}
		if len(joined) > siigoNameMaxLen {
			joined = strings.TrimSpace(joined[:siigoNameMaxLen])
		}
		return []string{joined}
	}

	switch len(parts) {
	case 0:
		return []string{"Sin", "Nombre"}
	case 1:
		return []string{parts[0], parts[0]}
	case 2:
		return []string{parts[0], parts[1]}
	case 3:
		return []string{parts[0], strings.Join(parts[1:], " ")}
	default:
		mid := len(parts) / 2
		return []string{strings.Join(parts[:mid], " "), strings.Join(parts[mid:], " ")}
	}
}

func normalizeSiigoPersonType(raw string) string {
	if strings.EqualFold(raw, "company") {
		return "Company"
	}
	return "Person"
}
