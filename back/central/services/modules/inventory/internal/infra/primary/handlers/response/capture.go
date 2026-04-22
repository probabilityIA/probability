package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
)

type LicensePlateLineResponse struct {
	ID         uint      `json:"id"`
	LpnID      uint      `json:"lpn_id"`
	BusinessID uint      `json:"business_id"`
	ProductID  string    `json:"product_id"`
	LotID      *uint     `json:"lot_id"`
	SerialID   *uint     `json:"serial_id"`
	Qty        int       `json:"qty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LicensePlateResponse struct {
	ID                uint                        `json:"id"`
	BusinessID        uint                        `json:"business_id"`
	Code              string                      `json:"code"`
	LpnType           string                      `json:"lpn_type"`
	CurrentLocationID *uint                       `json:"current_location_id"`
	Status            string                      `json:"status"`
	CreatedAt         time.Time                   `json:"created_at"`
	UpdatedAt         time.Time                   `json:"updated_at"`
	Lines             []LicensePlateLineResponse  `json:"Lines,omitempty"`
}

type LicensePlateListResponse struct {
	Data       []LicensePlateResponse `json:"data"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalPages int                    `json:"total_pages"`
}

type ScanResolutionResponse struct {
	Code       string         `json:"code"`
	CodeType   string         `json:"code_type"`
	MatchedID  *uint          `json:"matched_id,omitempty"`
	ProductID  string         `json:"product_id,omitempty"`
	LocationID *uint          `json:"location_id,omitempty"`
	LotID      *uint          `json:"lot_id,omitempty"`
	SerialID   *uint          `json:"serial_id,omitempty"`
	LpnID      *uint          `json:"lpn_id,omitempty"`
	Suggested  string         `json:"suggested,omitempty"`
	Data       map[string]any `json:"data,omitempty"`
}

type ScanEventResponse struct {
	ID          uint      `json:"id"`
	BusinessID  uint      `json:"business_id"`
	UserID      *uint     `json:"user_id"`
	DeviceID    string    `json:"device_id"`
	ScannedCode string    `json:"scanned_code"`
	CodeType    string    `json:"code_type"`
	Action      string    `json:"action"`
	ScannedAt   time.Time `json:"scanned_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ScanResponse struct {
	Resolved   bool                    `json:"resolved"`
	Resolution *ScanResolutionResponse `json:"resolution,omitempty"`
	Event      *ScanEventResponse      `json:"event,omitempty"`
}

type InventorySyncLogResponse struct {
	ID            uint       `json:"id"`
	BusinessID    uint       `json:"business_id"`
	IntegrationID *uint      `json:"integration_id"`
	Direction     string     `json:"direction"`
	PayloadHash   string     `json:"payload_hash"`
	Status        string     `json:"status"`
	Error         string     `json:"error"`
	SyncedAt      *time.Time `json:"synced_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

type InventorySyncLogListResponse struct {
	Data       []InventorySyncLogResponse `json:"data"`
	Total      int64                      `json:"total"`
	Page       int                        `json:"page"`
	PageSize   int                        `json:"page_size"`
	TotalPages int                        `json:"total_pages"`
}

type InboundSyncResultResponse struct {
	Log       *InventorySyncLogResponse `json:"log"`
	Duplicate bool                      `json:"duplicate"`
}

func LicensePlateLineFromEntity(l entities.LicensePlateLine) LicensePlateLineResponse {
	return LicensePlateLineResponse{
		ID:         l.ID,
		LpnID:      l.LpnID,
		BusinessID: l.BusinessID,
		ProductID:  l.ProductID,
		LotID:      l.LotID,
		SerialID:   l.SerialID,
		Qty:        l.Qty,
		CreatedAt:  l.CreatedAt,
		UpdatedAt:  l.UpdatedAt,
	}
}

func LicensePlateFromEntity(lpn *entities.LicensePlate) LicensePlateResponse {
	lines := make([]LicensePlateLineResponse, len(lpn.Lines))
	for i, l := range lpn.Lines {
		lines[i] = LicensePlateLineFromEntity(l)
	}
	return LicensePlateResponse{
		ID:                lpn.ID,
		BusinessID:        lpn.BusinessID,
		Code:              lpn.Code,
		LpnType:           lpn.LpnType,
		CurrentLocationID: lpn.CurrentLocationID,
		Status:            lpn.Status,
		CreatedAt:         lpn.CreatedAt,
		UpdatedAt:         lpn.UpdatedAt,
		Lines:             lines,
	}
}

func ScanResolutionFromEntity(r *entities.ScanResolution) *ScanResolutionResponse {
	if r == nil {
		return nil
	}
	return &ScanResolutionResponse{
		Code:       r.Code,
		CodeType:   r.CodeType,
		MatchedID:  r.MatchedID,
		ProductID:  r.ProductID,
		LocationID: r.LocationID,
		LotID:      r.LotID,
		SerialID:   r.SerialID,
		LpnID:      r.LpnID,
		Suggested:  r.Suggested,
		Data:       r.Data,
	}
}

func ScanEventFromEntity(e *entities.ScanEvent) *ScanEventResponse {
	if e == nil {
		return nil
	}
	return &ScanEventResponse{
		ID:          e.ID,
		BusinessID:  e.BusinessID,
		UserID:      e.UserID,
		DeviceID:    e.DeviceID,
		ScannedCode: e.ScannedCode,
		CodeType:    e.CodeType,
		Action:      e.Action,
		ScannedAt:   e.ScannedAt,
		CreatedAt:   e.CreatedAt,
	}
}

func SyncLogFromEntity(l *entities.InventorySyncLog) *InventorySyncLogResponse {
	if l == nil {
		return nil
	}
	return &InventorySyncLogResponse{
		ID:            l.ID,
		BusinessID:    l.BusinessID,
		IntegrationID: l.IntegrationID,
		Direction:     l.Direction,
		PayloadHash:   l.PayloadHash,
		Status:        l.Status,
		Error:         l.Error,
		SyncedAt:      l.SyncedAt,
		CreatedAt:     l.CreatedAt,
	}
}
