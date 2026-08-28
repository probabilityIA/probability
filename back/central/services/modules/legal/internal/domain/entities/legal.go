package entities

import "time"

type LegalDocument struct {
	ID            uint
	Code          string
	Version       string
	Title         string
	FileURL       string
	SHA256        string
	EffectiveFrom time.Time
	IsActive      bool
}

type LegalAcceptance struct {
	ID              uint
	UserID          uint
	LegalDocumentID uint
	BusinessID      *uint
	DocumentCode    string
	DocumentVersion string
	DocumentSHA256  string
	AcceptedAt      time.Time
	IPAddress       string
	UserAgent       string
	Method          string
}

const AcceptanceMethodWebModal = "web_modal"
