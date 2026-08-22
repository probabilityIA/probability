package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orderscompare/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/orderscompare"
)

const (
	defaultLimit    = 200
	maxLimit        = 1000
	defaultPageSize = 50
	maxPageSize     = 200
)

func (uc *useCase) Compare(ctx context.Context, query dtos.CompareQuery) (*dtos.ComparePage, error) {
	integration, err := uc.resolveIntegration(ctx, query.IntegrationID, query.BusinessID)
	if err != nil {
		return nil, err
	}

	channel := dtos.ChannelInfo{
		IntegrationID:   integration.ID,
		IntegrationName: integration.Name,
		IntegrationType: integration.IntegrationType,
		Supported:       uc.channels.Supports(integration.IntegrationType),
	}

	if !channel.Supported {
		return nil, fmt.Errorf("el canal %s todavia no expone sus ordenes para comparar", integration.Name)
	}

	limit := query.Limit
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	channelOrders, err := uc.channels.ListOrders(ctx, integration.IntegrationType, fmt.Sprintf("%d", integration.ID), orderscompare.ChannelFilters{
		From:  query.From,
		To:    query.To,
		Limit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("leyendo las ordenes del canal: %w", err)
	}

	locals, err := uc.repo.ListLocalOrders(ctx, integration.ID, query.BusinessID, query.From, query.To, limit)
	if err != nil {
		return nil, fmt.Errorf("leyendo las ordenes de Probability: %w", err)
	}

	result := orderscompare.Build(channelOrders, locals)
	rows := filterRows(result.Rows, query)

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	page := query.Page
	if page <= 0 {
		page = 1
	}

	pageRows, total := orderscompare.Paginate(rows, page, pageSize)
	totalPages := (total + pageSize - 1) / pageSize

	return &dtos.ComparePage{
		Rows:       pageRows,
		Totals:     result.Totals,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		CheckedAt:  time.Now(),
		Channel:    channel,
	}, nil
}

func (uc *useCase) resolveIntegration(ctx context.Context, integrationID, businessID uint) (*ports.Integration, error) {
	if integrationID == 0 {
		return nil, fmt.Errorf("integration_id es requerido")
	}
	integration, err := uc.repo.GetIntegration(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	if integration == nil {
		return nil, fmt.Errorf("la integracion no existe")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}
	return integration, nil
}

func filterRows(rows []orderscompare.Row, query dtos.CompareQuery) []orderscompare.Row {
	search := strings.ToLower(strings.TrimSpace(query.Search))
	if !query.OnlyDiff && search == "" {
		return rows
	}

	filtered := make([]orderscompare.Row, 0, len(rows))
	for _, row := range rows {
		if query.OnlyDiff && row.Action == orderscompare.ActionInSync && !row.StatusMismatch && !row.TotalMismatch {
			continue
		}
		if search != "" && !matches(row, search) {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func matches(row orderscompare.Row, search string) bool {
	return strings.Contains(strings.ToLower(row.ExternalID), search) ||
		strings.Contains(strings.ToLower(row.Number), search) ||
		strings.Contains(strings.ToLower(row.OrderNumber), search) ||
		strings.Contains(strings.ToLower(row.CustomerName), search)
}
