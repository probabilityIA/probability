package app

import (
	"context"
	"sort"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

type channelIssue struct {
	code     string
	severity string
	title    string
	detail   string
	value    func(domain.ChannelSummary) int
}

var channelIssues = []channelIssue{
	{
		code:     domain.FindingChannelNoSKU,
		severity: domain.SeverityError,
		title:    "Publicaciones sin SKU",
		detail:   "El emparejamiento se hace por SKU. Sin el no hay forma de asociarlas ni de sincronizarles stock.",
		value:    func(s domain.ChannelSummary) int { return s.ChannelNoSKU },
	},
	{
		code:     domain.FindingSKUChanged,
		severity: domain.SeverityError,
		title:    "SKU cambiado en el canal",
		detail:   "La asociacion sigue apuntando ahi, asi que se les esta reportando el inventario del producto anterior.",
		value:    func(s domain.ChannelSummary) int { return s.SKUChanged },
	},
	{
		code:     domain.FindingSKUSpacing,
		severity: domain.SeverityWarn,
		title:    "SKU con espacios de sobra",
		detail:   "Es el mismo SKU de los dos lados: solo sobra un espacio. Quitandolo el producto queda emparejado y vuelve a recibir stock.",
		value:    func(s domain.ChannelSummary) int { return s.SKUSpacing },
	},
	{
		code:     domain.FindingSKUTypo,
		severity: domain.SeverityWarn,
		title:    "Posibles errores de digitacion",
		detail:   "SKU que se parecen mucho a un producto sin emparejar. Suele ser el mismo producto escrito distinto en cada sistema.",
		value:    func(s domain.ChannelSummary) int { return s.SKUTypo },
	},
	{
		code:     domain.FindingNotAssociated,
		severity: domain.SeverityWarn,
		title:    "Productos sin asociar",
		detail:   "Coinciden por SKU pero nadie los asocio al canal, asi que no reciben stock.",
		value:    func(s domain.ChannelSummary) int { return s.NotAssociated },
	},
}

func (uc *useCase) Findings(ctx context.Context, businessID uint) (*domain.FindingsReport, error) {
	channels, err := uc.repo.ChannelSummaries(ctx, businessID)
	if err != nil {
		return nil, err
	}
	cross, err := uc.repo.CrossChannel(ctx, businessID)
	if err != nil {
		return nil, err
	}

	report := &domain.FindingsReport{
		BusinessID: businessID,
		Channels:   channels,
		Cross:      cross,
		Findings:   []domain.Finding{},
	}

	for _, issue := range channelIssues {
		total := 0
		afectados := make([]string, 0, len(channels))
		for _, c := range channels {
			if n := issue.value(c); n > 0 {
				total += n
				afectados = append(afectados, c.IntegrationName)
			}
		}
		if total == 0 {
			continue
		}
		report.Findings = append(report.Findings, domain.Finding{
			Code:     issue.code,
			Severity: issue.severity,
			Title:    issue.title,
			Detail:   issue.detail,
			Count:    total,
			Channels: afectados,
		})
	}

	if cross.SoldNotOwned > 0 {
		report.Findings = append(report.Findings, domain.Finding{
			Code:     domain.FindingSoldNotOwned,
			Severity: domain.SeverityError,
			Title:    "Se vende algo que no existe aqui",
			Detail:   "Publicado en algun canal pero sin producto en Probability: no hay inventario que descontar y la venta se puede caer.",
			Count:    cross.SoldNotOwned,
		})
	}

	orden := map[string]int{domain.SeverityError: 0, domain.SeverityWarn: 1, domain.SeverityInfo: 2}
	sort.SliceStable(report.Findings, func(i, j int) bool {
		a, b := report.Findings[i], report.Findings[j]
		if orden[a.Severity] != orden[b.Severity] {
			return orden[a.Severity] < orden[b.Severity]
		}
		return a.Count > b.Count
	})

	for _, f := range report.Findings {
		report.Total += f.Count
	}
	return report, nil
}
