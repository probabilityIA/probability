package client

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

type failedShipmentsBody struct {
	FailedShipments []struct {
		ShipmentID string `json:"shipment_id"`
		Message    string `json:"message"`
		ErrorCode  string `json:"error_code"`
	} `json:"failed_shipments"`
}

func classifyLabelBody(body []byte) error {
	var parsed failedShipmentsBody
	if err := json.Unmarshal(body, &parsed); err != nil || len(parsed.FailedShipments) == 0 {
		return nil
	}

	for _, fallo := range parsed.FailedShipments {
		motivo := strings.ToLower(fallo.Message + " " + fallo.ErrorCode)
		if strings.Contains(motivo, "shipped") || strings.Contains(motivo, "delivered") ||
			strings.Contains(motivo, "not_ready_to_print") {
			return domain.ErrLabelAlreadyShipped
		}
	}

	return domain.ErrLabelNotAvailable
}

func looksLikePDF(body []byte) bool {
	return bytes.HasPrefix(bytes.TrimLeft(body, " \r\n\t"), []byte("%PDF-"))
}

func looksLikeJSON(body []byte) bool {
	trimmed := bytes.TrimLeft(body, " \r\n\t")
	return bytes.HasPrefix(trimmed, []byte("{")) || bytes.HasPrefix(trimmed, []byte("["))
}
