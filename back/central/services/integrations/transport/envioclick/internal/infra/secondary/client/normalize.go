package client

import (
	"errors"
	"strings"
	"unicode"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

const (
	maxAddressLen     = 50
	maxSuburbLen      = 30
	maxDescriptionLen = 25
	maxFirstNameLen   = 14
	maxLastNameLen    = 14
	maxCompanyLen     = 30
	maxReferenceLen   = 28
	minTextLen        = 2
	maxPhoneDigits    = 10
)

var (
	ErrDestinationAddressMissing = errors.New("la dirección de destino está vacía: completá la dirección del cliente en la orden")
	ErrOriginAddressMissing      = errors.New("la dirección de la bodega de origen está vacía: completá los datos de la bodega")
)

func collapseSpaces(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func clampText(text string, limit int) string {
	text = collapseSpaces(text)
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}

	recortado := string(runes[:limit])
	if corte := strings.LastIndex(recortado, " "); corte >= limit/2 {
		recortado = recortado[:corte]
	}
	return strings.TrimRight(recortado, " ,.-")
}

func clampPhone(phone string) string {
	var digitos strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			digitos.WriteRune(r)
		}
	}

	limpio := digitos.String()
	if len(limpio) > maxPhoneDigits {
		return limpio[len(limpio)-maxPhoneDigits:]
	}
	return limpio
}

func normalizeAddress(addr *domain.Address) {
	addr.Address = clampText(addr.Address, maxAddressLen)
	addr.Suburb = clampText(addr.Suburb, maxSuburbLen)
	addr.FirstName = clampText(addr.FirstName, maxFirstNameLen)
	addr.LastName = clampText(addr.LastName, maxLastNameLen)
	addr.Company = clampText(addr.Company, maxCompanyLen)
	addr.CrossStreet = clampText(addr.CrossStreet, maxAddressLen)
	addr.Reference = clampText(addr.Reference, maxAddressLen)
	addr.Email = collapseSpaces(addr.Email)
	addr.Phone = clampPhone(addr.Phone)

	if len([]rune(addr.Suburb)) < minTextLen {
		addr.Suburb = addr.Address
		if len([]rune(addr.Suburb)) > maxSuburbLen {
			addr.Suburb = clampText(addr.Suburb, maxSuburbLen)
		}
	}
}

func normalizeQuoteRequest(req *domain.QuoteRequest) error {
	normalizeAddress(&req.Origin)
	normalizeAddress(&req.Destination)

	req.Description = clampText(req.Description, maxDescriptionLen)
	if len([]rune(req.Description)) < minTextLen {
		req.Description = "Compra en linea"
	}

	req.MyShipmentReference = clampText(req.MyShipmentReference, maxReferenceLen)
	req.ExternalOrderID = clampText(req.ExternalOrderID, maxReferenceLen)

	if len([]rune(req.Destination.Address)) < minTextLen {
		return ErrDestinationAddressMissing
	}
	if len([]rune(req.Origin.Address)) < minTextLen {
		return ErrOriginAddressMissing
	}

	return nil
}
