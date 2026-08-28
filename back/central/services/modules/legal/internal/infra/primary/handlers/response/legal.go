package response

type LegalDocument struct {
	ID            uint   `json:"id"`
	Code          string `json:"code"`
	Version       string `json:"version"`
	Title         string `json:"title"`
	FileURL       string `json:"file_url"`
	SHA256        string `json:"sha256"`
	ContentHTML   string `json:"content_html"`
	EffectiveFrom string `json:"effective_from"`
}

type PendingDocuments struct {
	RequiresAcceptance bool            `json:"requires_acceptance"`
	Documents          []LegalDocument `json:"documents"`
}

type AcceptResult struct {
	AcceptedAt  string `json:"accepted_at"`
	DocumentIDs []uint `json:"document_ids"`
}
