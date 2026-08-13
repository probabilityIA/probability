package inventorycompare

import (
	"sort"
	"time"
)

const (
	PageSizeDefault = 100
	PageSizeMax     = 200
)

type Page struct {
	Rows       []Row     `json:"rows"`
	Totals     Totals    `json:"totals"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	Total      int       `json:"total"`
	TotalPages int       `json:"total_pages"`
	CheckedAt  time.Time `json:"checked_at"`
}

func NormalizePaging(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = PageSizeDefault
	}
	if pageSize > PageSizeMax {
		pageSize = PageSizeMax
	}
	return page, pageSize
}

func Slice[T any](items []T, page, pageSize int) []T {
	inicio := (page - 1) * pageSize
	if inicio >= len(items) {
		return nil
	}
	fin := inicio + pageSize
	if fin > len(items) {
		fin = len(items)
	}
	return items[inicio:fin]
}

func TotalPages(total, pageSize int) int {
	if pageSize <= 0 || total <= 0 {
		return 1
	}
	paginas := total / pageSize
	if total%pageSize != 0 {
		paginas++
	}
	return paginas
}

const (
	ActionUpdate    = "update"
	ActionUnchanged = "unchanged"
	ActionSkip      = "skip"
)

const (
	ReasonSinInventario  = "Probability no lleva inventario de este producto, se enviaria 0"
	ReasonSinPublicacion = "no se encontro la publicacion en el canal"
	ReasonSinManejoStock = "la publicacion no maneja stock en el canal"
)

type Item struct {
	ProductID         string
	SKU               string
	Name              string
	ExternalItemID    string
	ExternalVariantID string
	ProbabilityQty    int
	HasStock          bool
	ChannelQty        int
	ChannelFound      bool
	ChannelManaged    bool
	SkipReason        string
}

type Row struct {
	ProductID         string `json:"product_id"`
	SKU               string `json:"sku"`
	Name              string `json:"name"`
	ExternalItemID    string `json:"external_item_id"`
	ExternalVariantID string `json:"external_variant_id,omitempty"`
	ProbabilityQty    *int   `json:"probability_qty"`
	ChannelQty        *int   `json:"channel_qty"`
	Delta             *int   `json:"delta"`
	Action            string `json:"action"`
	Reason            string `json:"reason,omitempty"`
}

type Totals struct {
	Total     int `json:"total"`
	ToUpdate  int `json:"to_update"`
	Unchanged int `json:"unchanged"`
	Skipped   int `json:"skipped"`
}

func BuildRow(item Item) Row {
	fila := Row{
		ProductID:         item.ProductID,
		SKU:               item.SKU,
		Name:              item.Name,
		ExternalItemID:    item.ExternalItemID,
		ExternalVariantID: item.ExternalVariantID,
	}
	if !item.HasStock {
		item.ProbabilityQty = 0
	}
	cantidadPropia := item.ProbabilityQty
	fila.ProbabilityQty = &cantidadPropia
	if item.ChannelFound {
		cantidad := item.ChannelQty
		fila.ChannelQty = &cantidad
	}

	switch {
	case item.SkipReason != "":
		fila.Action = ActionSkip
		fila.Reason = item.SkipReason
	case !item.ChannelFound:
		fila.Action = ActionSkip
		fila.Reason = ReasonSinPublicacion
	case !item.ChannelManaged:
		fila.Action = ActionSkip
		fila.Reason = ReasonSinManejoStock
	case item.ChannelQty == item.ProbabilityQty:
		fila.Action = ActionUnchanged
	default:
		fila.Action = ActionUpdate
	}

	if fila.Action == ActionUpdate || fila.Action == ActionUnchanged {
		delta := item.ProbabilityQty - item.ChannelQty
		fila.Delta = &delta
	}
	if !item.HasStock && fila.Action == ActionUpdate {
		fila.Reason = ReasonSinInventario
	}

	return fila
}

func Build(items []Item) []Row {
	filas := make([]Row, 0, len(items))
	for _, item := range items {
		filas = append(filas, BuildRow(item))
	}
	sort.SliceStable(filas, func(i, j int) bool { return filas[i].SKU < filas[j].SKU })
	return filas
}

func Summarize(filas []Row) Totals {
	totales := Totals{Total: len(filas)}
	for _, fila := range filas {
		switch fila.Action {
		case ActionUpdate:
			totales.ToUpdate++
		case ActionUnchanged:
			totales.Unchanged++
		default:
			totales.Skipped++
		}
	}
	return totales
}
