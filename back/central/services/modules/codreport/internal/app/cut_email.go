package app

import (
	"context"
	"fmt"
	"html"
	"sort"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

var bogota = func() *time.Location {
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return time.FixedZone("COT", -5*3600)
	}
	return loc
}()

var monthsES = []string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}

func (uc *UseCase) CutEmailPreview(ctx context.Context, businessID uint, cutID uint) (*dtos.CutEmailPreview, error) {
	cut, err := uc.repo.CutByID(ctx, businessID, cutID)
	if err != nil {
		return nil, err
	}
	if cut.Status != "confirmed" {
		return nil, fmt.Errorf("solo se puede enviar por correo un corte confirmado")
	}

	orders, err := uc.CutOrders(ctx, businessID, cutID)
	if err != nil {
		return nil, err
	}

	businessName := uc.repo.BusinessName(ctx, businessID)
	return &dtos.CutEmailPreview{
		Subject: fmt.Sprintf("Corte de pago contra entrega %s - %s", periodLabel(cut.PeriodStart, cut.PeriodEnd), businessName),
		HTML:    renderCutEmailHTML(cut, orders, businessName),
	}, nil
}

func (uc *UseCase) SendCutEmail(ctx context.Context, d dtos.SendCutEmailDTO) error {
	if uc.email == nil {
		return fmt.Errorf("servicio de correo no configurado")
	}
	if len(d.Recipients) == 0 {
		return fmt.Errorf("sin destinatarios")
	}

	preview, err := uc.CutEmailPreview(ctx, d.BusinessID, d.CutID)
	if err != nil {
		return err
	}
	subject, body := preview.Subject, preview.HTML

	var failed []string
	for _, to := range d.Recipients {
		if err := uc.email.SendHTML(ctx, to, subject, body); err != nil {
			uc.log.Error(ctx).Err(err).Uint("cut_id", d.CutID).Str("to", to).Msg("Error enviando corte COD por correo")
			failed = append(failed, to)
		}
	}
	if len(failed) == len(d.Recipients) {
		return fmt.Errorf("no se pudo enviar el correo a ningun destinatario")
	}
	if len(failed) > 0 {
		return fmt.Errorf("no se pudo enviar a: %s", strings.Join(failed, ", "))
	}
	return nil
}

func dateES(t time.Time) string {
	t = t.In(bogota)
	return fmt.Sprintf("%02d %s %d", t.Day(), monthsES[t.Month()-1], t.Year())
}

func dateTimeES(t time.Time) string {
	t = t.In(bogota)
	return fmt.Sprintf("%02d %s %d, %02d:%02d", t.Day(), monthsES[t.Month()-1], t.Year(), t.Hour(), t.Minute())
}

func dateOnlyES(t time.Time) string {
	t = t.UTC()
	return fmt.Sprintf("%02d %s %d", t.Day(), monthsES[t.Month()-1], t.Year())
}

func periodLabel(start, end time.Time) string {
	return dateOnlyES(start) + " - " + dateOnlyES(end)
}

func formatCOP(v float64) string {
	neg := v < 0
	if neg {
		v = -v
	}
	whole := int64(v + 0.5)
	s := fmt.Sprintf("%d", whole)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := "$" + strings.Join(parts, ".")
	if neg {
		return "-" + out
	}
	return out
}

func carrierTitle(c string) string {
	c = strings.TrimSpace(c)
	if c == "" || c == "SIN TRANSPORTADORA" {
		return "Sin transportadora"
	}
	lower := strings.ToLower(c)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

const logoURL = "https://probabilityia.com.co/logo2recortado.png"

const cellStyle = `padding:7px 8px;border-bottom:1px solid #f3f4f6;font-size:12px`
const headStyle = `padding:8px;border-bottom:1px solid #e5e7eb;font-size:11px;text-transform:uppercase;color:#6b7280`

func renderCutEmailHTML(cut *entities.PaymentCut, orders []entities.CodOrder, businessName string) string {
	esc := html.EscapeString
	var b strings.Builder

	b.WriteString(`<div style="font-family:Arial,Helvetica,sans-serif;max-width:760px;margin:0 auto;color:#1f2937">`)
	b.WriteString(`<div style="padding:18px 24px;border:1px solid #e5e7eb;border-bottom:none;border-radius:8px 8px 0 0;background:#ffffff">`)
	b.WriteString(`<img src="` + logoURL + `" alt="ProbabilityIA" width="170" style="display:block;height:auto;max-width:170px;border:0">`)
	b.WriteString(`</div>`)
	b.WriteString(`<div style="background:#111827;color:#ffffff;padding:22px 24px">`)
	b.WriteString(`<h1 style="margin:0;font-size:19px">Corte de pago contra entrega</h1>`)
	b.WriteString(`<p style="margin:6px 0 0;font-size:13px;color:#d1d5db">` + esc(businessName) + `</p>`)
	b.WriteString(`</div>`)
	b.WriteString(`<div style="border:1px solid #e5e7eb;border-top:none;padding:22px 24px;border-radius:0 0 8px 8px">`)

	b.WriteString(`<table style="width:100%;font-size:13px;margin-bottom:14px;border-collapse:collapse">`)
	b.WriteString(`<tr><td style="padding:3px 0"><strong>Periodo:</strong> ` + esc(periodLabel(cut.PeriodStart, cut.PeriodEnd)) + `</td>`)
	b.WriteString(`<td style="padding:3px 0;text-align:right"><strong>Corte No.</strong> ` + fmt.Sprintf("%d", cut.ID) + `</td></tr>`)
	if cut.ConfirmedAt != nil {
		b.WriteString(`<tr><td style="padding:3px 0" colspan="2"><strong>Confirmado:</strong> ` + esc(dateTimeES(*cut.ConfirmedAt)))
		if cut.ConfirmedByName != "" {
			b.WriteString(` por ` + esc(cut.ConfirmedByName))
		}
		b.WriteString(`</td></tr>`)
	}
	b.WriteString(`</table>`)

	b.WriteString(`<table style="width:100%;border-collapse:collapse;margin-bottom:18px"><tr>`)
	b.WriteString(summaryBox("Ordenes consignadas", fmt.Sprintf("%d", cut.OrdersCount), "#1f2937"))
	b.WriteString(summaryBox("Total consignado", formatCOP(cut.TotalCollected), "#059669"))
	b.WriteString(`</tr></table>`)

	if len(cut.ByCarrier) > 0 {
		aggs := append([]entities.CarrierAggregate{}, cut.ByCarrier...)
		sort.Slice(aggs, func(i, j int) bool { return aggs[i].TotalCollected > aggs[j].TotalCollected })
		b.WriteString(`<h3 style="font-size:13px;margin:0 0 6px;color:#374151">Resumen por transportadora</h3>`)
		b.WriteString(`<table style="width:100%;border-collapse:collapse;margin-bottom:18px"><thead><tr style="background:#f3f4f6">`)
		b.WriteString(`<th style="text-align:left;` + headStyle + `">Transportadora</th>`)
		b.WriteString(`<th style="text-align:right;` + headStyle + `">Ordenes</th>`)
		b.WriteString(`<th style="text-align:right;` + headStyle + `">Consignado</th>`)
		b.WriteString(`</tr></thead><tbody>`)
		for _, a := range aggs {
			b.WriteString(`<tr><td style="` + cellStyle + `">` + esc(carrierTitle(a.Carrier)) + `</td>`)
			b.WriteString(`<td style="text-align:right;` + cellStyle + `">` + fmt.Sprintf("%d", a.OrdersCount) + `</td>`)
			b.WriteString(`<td style="text-align:right;` + cellStyle + `">` + formatCOP(a.TotalCollected) + `</td></tr>`)
		}
		b.WriteString(`</tbody></table>`)
	}

	b.WriteString(`<h3 style="font-size:13px;margin:0 0 6px;color:#374151">Relacion de ordenes (` + fmt.Sprintf("%d", len(orders)) + `)</h3>`)
	b.WriteString(`<table style="width:100%;border-collapse:collapse"><thead><tr style="background:#f3f4f6">`)
	b.WriteString(`<th style="text-align:left;` + headStyle + `">Orden</th>`)
	b.WriteString(`<th style="text-align:left;` + headStyle + `">Cliente</th>`)
	b.WriteString(`<th style="text-align:left;` + headStyle + `">Transportadora</th>`)
	b.WriteString(`<th style="text-align:left;` + headStyle + `">Guia</th>`)
	b.WriteString(`<th style="text-align:left;` + headStyle + `">Entregado</th>`)
	b.WriteString(`<th style="text-align:right;` + headStyle + `">Recaudado</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	var total float64
	for _, o := range orders {
		total += o.CodTotal
		guide := o.GuideNumber
		if guide == "" {
			guide = "-"
		}
		delivered := "-"
		if o.DeliveredAt != nil {
			delivered = dateES(*o.DeliveredAt)
		}
		b.WriteString(`<tr>`)
		b.WriteString(`<td style="` + cellStyle + `;font-weight:bold">` + esc(o.OrderNumber) + `</td>`)
		b.WriteString(`<td style="` + cellStyle + `">` + esc(o.CustomerName) + `</td>`)
		b.WriteString(`<td style="` + cellStyle + `">` + esc(carrierTitle(o.Carrier)) + `</td>`)
		b.WriteString(`<td style="` + cellStyle + `;font-family:monospace">` + esc(guide) + `</td>`)
		b.WriteString(`<td style="` + cellStyle + `">` + esc(delivered) + `</td>`)
		b.WriteString(`<td style="text-align:right;` + cellStyle + `">` + formatCOP(o.CodTotal) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody><tfoot><tr>`)
	b.WriteString(`<td colspan="5" style="padding:10px 8px;text-align:right;font-size:13px"><strong>Total consignado</strong></td>`)
	b.WriteString(`<td style="padding:10px 8px;text-align:right;font-size:14px;color:#059669"><strong>` + formatCOP(total) + `</strong></td>`)
	b.WriteString(`</tr></tfoot></table>`)

	b.WriteString(`<p style="font-size:11px;color:#9ca3af;margin-top:24px">Este correo fue generado por la plataforma Probability.</p>`)
	b.WriteString(`</div></div>`)
	return b.String()
}

func summaryBox(label, value, color string) string {
	return `<td style="width:50%;padding:0 6px 0 0"><div style="border:1px solid #e5e7eb;border-radius:8px;padding:12px 14px">` +
		`<div style="font-size:11px;text-transform:uppercase;color:#6b7280">` + html.EscapeString(label) + `</div>` +
		`<div style="font-size:20px;font-weight:bold;color:` + color + `;margin-top:4px">` + html.EscapeString(value) + `</div>` +
		`</div></td>`
}
