package domain

import "time"

type GetOrdersParams struct {
	Statuses []string
	After    *time.Time
	Before   *time.Time
	Page     int
	PerPage  int
}

type GetOrdersResult struct {
	Orders []JumpsellerOrder
	Count  int
}

type Credential struct {
	APIKey      string
	APISecret   string
	AccessToken string
	BaseURL     string
}

func (c Credential) PacerKey() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	return c.AccessToken
}
