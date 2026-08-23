package request

type ApplyRequest struct {
	IntegrationID uint     `json:"integration_id" binding:"required"`
	ExternalIDs   []string `json:"external_ids" binding:"required"`
}
