package response

type BusinessSimpleResponse struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Code            string `json:"code,omitempty"`
	LogoURL         string `json:"logo_url,omitempty"`
	PrimaryColor    string `json:"primary_color,omitempty"`
	SecondaryColor  string `json:"secondary_color,omitempty"`
	TertiaryColor   string `json:"tertiary_color,omitempty"`
	QuaternaryColor string `json:"quaternary_color,omitempty"`
}

type GetBusinessesSimpleResponse struct {
	Success    bool                     `json:"success"`
	Message    string                   `json:"message"`
	Data       []BusinessSimpleResponse `json:"data"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	PageSize   int                      `json:"page_size"`
	TotalPages int                      `json:"total_pages"`
}
