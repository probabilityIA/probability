package app

import (
	"fmt"
	"html"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

func formatCOP(v float64) string {
	neg := v < 0
	if neg {
		v = -v
	}
	whole := int64(v + 0.005)
	cents := int64((v-float64(whole))*100 + 0.5)
	s := fmt.Sprintf("%d", whole)
	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	out := "$ " + strings.Join(parts, ".")
	if cents > 0 {
		out += fmt.Sprintf(",%02d", cents)
	}
	if neg {
		out = "-" + out
	}
	return out
}

func renderInvoiceHTML(inv *entities.Invoice) string {
	var b strings.Builder
	esc := html.EscapeString
	b.WriteString(`<div style="font-family:Arial,Helvetica,sans-serif;max-width:640px;margin:0 auto;color:#1f2937">`)
	b.WriteString(`<div style="background:#1f2937;color:#ffffff;padding:24px;border-radius:8px 8px 0 0">`)
	b.WriteString(`<h1 style="margin:0;font-size:20px">Probability</h1>`)
	b.WriteString(`<p style="margin:4px 0 0;font-size:13px;color:#d1d5db">Factura ` + esc(inv.Number) + `</p>`)
	b.WriteString(`</div>`)
	b.WriteString(`<div style="border:1px solid #e5e7eb;border-top:none;padding:24px;border-radius:0 0 8px 8px">`)
	b.WriteString(`<table style="width:100%;font-size:13px;margin-bottom:16px"><tr>`)
	b.WriteString(`<td><strong>Cliente:</strong> ` + esc(inv.BusinessName) + `</td>`)
	b.WriteString(`<td style="text-align:right"><strong>Fecha:</strong> ` + inv.IssueDate.Format("2006-01-02") + `</td>`)
	b.WriteString(`</tr>`)
	if inv.DueDate != nil {
		b.WriteString(`<tr><td></td><td style="text-align:right"><strong>Vence:</strong> ` + inv.DueDate.Format("2006-01-02") + `</td></tr>`)
	}
	b.WriteString(`</table>`)
	b.WriteString(`<table style="width:100%;border-collapse:collapse;font-size:13px">`)
	b.WriteString(`<thead><tr style="background:#f3f4f6">`)
	b.WriteString(`<th style="text-align:left;padding:8px;border-bottom:1px solid #e5e7eb">Descripcion</th>`)
	b.WriteString(`<th style="text-align:right;padding:8px;border-bottom:1px solid #e5e7eb">Cant.</th>`)
	b.WriteString(`<th style="text-align:right;padding:8px;border-bottom:1px solid #e5e7eb">Precio</th>`)
	b.WriteString(`<th style="text-align:right;padding:8px;border-bottom:1px solid #e5e7eb">Total</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	for _, item := range inv.Items {
		b.WriteString(`<tr>`)
		b.WriteString(`<td style="padding:8px;border-bottom:1px solid #f3f4f6">` + esc(item.Description) + `</td>`)
		b.WriteString(fmt.Sprintf(`<td style="text-align:right;padding:8px;border-bottom:1px solid #f3f4f6">%g</td>`, item.Quantity))
		b.WriteString(`<td style="text-align:right;padding:8px;border-bottom:1px solid #f3f4f6">` + formatCOP(item.UnitPrice) + `</td>`)
		b.WriteString(`<td style="text-align:right;padding:8px;border-bottom:1px solid #f3f4f6">` + formatCOP(item.Amount) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	b.WriteString(`<table style="width:100%;font-size:13px;margin-top:16px">`)
	b.WriteString(`<tr><td style="text-align:right;padding:4px"><strong>Subtotal:</strong></td><td style="text-align:right;padding:4px;width:120px">` + formatCOP(inv.Subtotal) + `</td></tr>`)
	for _, tax := range inv.TaxDetail {
		b.WriteString(fmt.Sprintf(`<tr><td style="text-align:right;padding:4px">%s (%g%%):</td><td style="text-align:right;padding:4px">%s</td></tr>`, esc(tax.TaxName), tax.RatePercent, formatCOP(tax.Amount)))
	}
	b.WriteString(`<tr><td style="text-align:right;padding:8px 4px;font-size:16px"><strong>Total:</strong></td><td style="text-align:right;padding:8px 4px;font-size:16px"><strong>` + formatCOP(inv.Total) + `</strong></td></tr>`)
	for _, ret := range inv.WithholdingDetail {
		b.WriteString(fmt.Sprintf(`<tr><td style="text-align:right;padding:4px;color:#6b7280">%s (%g%%) informativa:</td><td style="text-align:right;padding:4px;color:#6b7280">-%s</td></tr>`, esc(ret.TaxName), ret.RatePercent, formatCOP(ret.Amount)))
	}
	if inv.WithholdingTotal > 0 {
		b.WriteString(`<tr><td style="text-align:right;padding:4px"><strong>Neto a recibir:</strong></td><td style="text-align:right;padding:4px"><strong>` + formatCOP(inv.Total-inv.WithholdingTotal) + `</strong></td></tr>`)
	}
	b.WriteString(`</table>`)
	if inv.Notes != "" {
		b.WriteString(`<p style="font-size:12px;color:#6b7280;margin-top:16px;border-top:1px solid #e5e7eb;padding-top:12px">` + esc(inv.Notes) + `</p>`)
	}
	b.WriteString(`<p style="font-size:11px;color:#9ca3af;margin-top:24px">Este correo fue generado por la plataforma Probability.</p>`)
	b.WriteString(`</div></div>`)
	return b.String()
}
