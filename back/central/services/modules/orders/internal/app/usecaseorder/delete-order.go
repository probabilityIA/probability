package usecaseorder

import (
	"context"
	"errors"
	"fmt"
)

// DeleteOrder elimina (soft delete) una orden
func (uc *UseCaseOrder) DeleteOrder(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("order ID is required")
	}

	// Verificar que la orden existe
	_, err := uc.repo.GetOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting order: %w", err)
	}

	// Eliminar la orden
	if err := uc.repo.DeleteOrder(ctx, id); err != nil {
		return fmt.Errorf("error deleting order: %w", err)
	}

	return nil
}
