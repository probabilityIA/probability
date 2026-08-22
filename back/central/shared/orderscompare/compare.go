package orderscompare

import (
	"sort"
	"strings"
	"time"
)

const (
	ActionCreate = "create"
	ActionInSync = "in_sync"
	ActionOnlyIn = "only_in_probability"
)

type ChannelOrder struct {
	ExternalID        string
	Number            string
	CustomerName      string
	Status            string
	RawStatus         string
	FulfillmentStatus string
	Total             float64
	Currency          string
	Items             int
	CreatedAt         time.Time
	URL               string
}

type ChannelFilters struct {
	From  *time.Time
	To    *time.Time
	Limit int
}

type ImportResult struct {
	Queued   []string
	Failed   map[string]string
	NotFound []string
}

type LocalOrder struct {
	OrderID      string
	OrderNumber  string
	ExternalID   string
	Status       string
	Total        float64
	Currency     string
	CustomerName string
	CreatedAt    time.Time
}

type Row struct {
	ExternalID     string    `json:"external_id"`
	Number         string    `json:"number"`
	CustomerName   string    `json:"customer_name"`
	ChannelStatus  string    `json:"channel_status"`
	RawStatus      string    `json:"raw_status"`
	Fulfillment    string    `json:"fulfillment_status,omitempty"`
	LocalStatus    string    `json:"local_status,omitempty"`
	OrderID        string    `json:"order_id,omitempty"`
	OrderNumber    string    `json:"order_number,omitempty"`
	Total          float64   `json:"total"`
	LocalTotal     float64   `json:"local_total,omitempty"`
	Currency       string    `json:"currency"`
	Items          int       `json:"items"`
	CreatedAt      time.Time `json:"created_at"`
	URL            string    `json:"url,omitempty"`
	Action         string    `json:"action"`
	MovesInventory bool      `json:"moves_inventory"`
	InventoryNote  string    `json:"inventory_note,omitempty"`
	StatusMismatch bool      `json:"status_mismatch"`
	TotalMismatch  bool      `json:"total_mismatch"`
}

type Totals struct {
	Total              int `json:"total"`
	ToCreate           int `json:"to_create"`
	InSync             int `json:"in_sync"`
	OnlyInProbability  int `json:"only_in_probability"`
	WithoutInventory   int `json:"without_inventory"`
	WithStatusMismatch int `json:"with_status_mismatch"`
}

type Result struct {
	Rows   []Row  `json:"rows"`
	Totals Totals `json:"totals"`
}

func Build(channelOrders []ChannelOrder, locals []LocalOrder) Result {
	localByExternal := make(map[string]LocalOrder, len(locals))
	for _, l := range locals {
		key := normalize(l.ExternalID)
		if key == "" {
			continue
		}
		localByExternal[key] = l
	}

	seen := make(map[string]bool, len(channelOrders))
	rows := make([]Row, 0, len(channelOrders)+len(locals))
	totals := Totals{}

	for _, co := range channelOrders {
		key := normalize(co.ExternalID)
		seen[key] = true

		row := Row{
			ExternalID:    co.ExternalID,
			Number:        co.Number,
			CustomerName:  co.CustomerName,
			ChannelStatus: co.Status,
			RawStatus:     co.RawStatus,
			Fulfillment:   co.FulfillmentStatus,
			Total:         co.Total,
			Currency:      co.Currency,
			Items:         co.Items,
			CreatedAt:     co.CreatedAt,
			URL:           co.URL,
		}

		skips, reason := SkipsInventoryFor(co.Status, co.FulfillmentStatus)
		row.MovesInventory = !skips
		row.InventoryNote = reason

		if local, ok := localByExternal[key]; ok {
			row.Action = ActionInSync
			row.OrderID = local.OrderID
			row.OrderNumber = local.OrderNumber
			row.LocalStatus = local.Status
			row.LocalTotal = local.Total
			row.StatusMismatch = !sameStatus(local.Status, co.Status)
			row.TotalMismatch = !sameTotal(local.Total, co.Total)
			totals.InSync++
			if row.StatusMismatch {
				totals.WithStatusMismatch++
			}
		} else {
			row.Action = ActionCreate
			totals.ToCreate++
			if skips {
				totals.WithoutInventory++
			}
		}

		rows = append(rows, row)
	}

	for _, l := range locals {
		key := normalize(l.ExternalID)
		if key != "" && seen[key] {
			continue
		}
		rows = append(rows, Row{
			ExternalID:     l.ExternalID,
			Number:         l.OrderNumber,
			CustomerName:   l.CustomerName,
			LocalStatus:    l.Status,
			OrderID:        l.OrderID,
			OrderNumber:    l.OrderNumber,
			LocalTotal:     l.Total,
			Currency:       l.Currency,
			CreatedAt:      l.CreatedAt,
			Action:         ActionOnlyIn,
			MovesInventory: false,
		})
		totals.OnlyInProbability++
	}

	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].CreatedAt.After(rows[j].CreatedAt)
	})

	totals.Total = len(rows)
	return Result{Rows: rows, Totals: totals}
}

func Paginate(rows []Row, page, pageSize int) ([]Row, int) {
	if pageSize <= 0 {
		pageSize = 50
	}
	if page <= 0 {
		page = 1
	}
	total := len(rows)
	start := (page - 1) * pageSize
	if start >= total {
		return []Row{}, total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return rows[start:end], total
}

func normalize(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

func sameStatus(a, b string) bool {
	return normalize(a) == normalize(b)
}

func sameTotal(a, b float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < 0.01
}
