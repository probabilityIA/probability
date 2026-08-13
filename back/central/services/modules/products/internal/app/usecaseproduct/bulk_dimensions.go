package usecaseproduct

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCaseProduct) ExportDimensions(ctx context.Context, businessID uint) ([]domain.Product, error) {
	all := make([]domain.Product, 0)
	page := 1
	const pageSize = 100

	for {
		batch, total, err := uc.repo.ListProducts(ctx, businessID, page, pageSize, nil)
		if err != nil {
			return nil, fmt.Errorf("error listing products for export: %w", err)
		}
		all = append(all, batch...)
		if int64(len(all)) >= total || len(batch) == 0 {
			break
		}
		page++
	}

	return all, nil
}

func (uc *UseCaseProduct) BulkUpdateDimensions(ctx context.Context, businessID uint, rows []domain.DimensionRow) (*domain.BulkDimensionsResult, error) {
	result := &domain.BulkDimensionsResult{
		TotalRows: len(rows),
		Errors:    []string{},
	}

	for i, row := range rows {
		rowNum := i + 2

		product, err := uc.repo.GetProductBySKU(ctx, businessID, row.SKU)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d (SKU %s): no existe en este negocio", rowNum, row.SKU))
			continue
		}

		req := &domain.UpdateProductRequest{
			Weight: row.Weight,
			Length: row.Length,
			Width:  row.Width,
			Height: row.Height,
		}

		if _, err := uc.UpdateProduct(ctx, businessID, product.ID, req); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("fila %d (SKU %s): %s", rowNum, row.SKU, err.Error()))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}
