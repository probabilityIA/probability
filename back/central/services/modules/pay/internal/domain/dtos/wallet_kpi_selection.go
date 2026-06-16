package dtos

type WalletKPISelectionResponse struct {
	ID                  uint   `json:"id"`
	SelectedBusinessIDs []uint `json:"selected_business_ids"`
}

type UpdateWalletKPISelectionRequest struct {
	SelectedBusinessIDs []uint `json:"selected_business_ids"`
}
