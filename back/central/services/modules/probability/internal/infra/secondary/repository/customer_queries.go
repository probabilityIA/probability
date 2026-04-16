package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/probability/internal/domain/entities"
)

// ============================================
// MÉTODOS DE CONSULTA A TABLAS DE CLIENTES, ENVÍOS Y PAGOS
// (Replicados localmente - no compartir repos entre módulos)
// ============================================

// GetCustomerOrderHistory obtiene el historial de órdenes de un cliente,
// excluyendo la orden actual que se está evaluando.
//
// Tablas consultadas: orders, payments (gestionadas por módulos orders/payments)
func (r *Repository) GetCustomerOrderHistory(ctx context.Context, customerID uint, excludeOrderID string) (*entities.CustomerHistory, error) {
	var result struct {
		TotalOrders       int64
		TotalSpent        float64
		AvgOrderValue     float64
		FirstOrderDate    *time.Time
		LastOrderDate     *time.Time
		NoveltyCount      int64
		CODOrderCount     int64
		DistinctAddresses int64
	}

	err := r.db.Conn(ctx).Raw(`
		SELECT
			COUNT(*) AS total_orders,
			COALESCE(SUM(total_amount), 0) AS total_spent,
			COALESCE(AVG(total_amount), 0) AS avg_order_value,
			MIN(created_at) AS first_order_date,
			MAX(created_at) AS last_order_date,
			COUNT(CASE WHEN novelty IS NOT NULL AND novelty != '' THEN 1 END) AS novelty_count,
			COUNT(CASE WHEN cod_total IS NOT NULL AND cod_total > 0 THEN 1 END) AS cod_order_count,
			COUNT(DISTINCT shipping_street) AS distinct_addresses
		FROM orders
		WHERE customer_id = ? AND id != ? AND deleted_at IS NULL
	`, customerID, excludeOrderID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Failed payments: separate query joining payments
	var failedPayments int64
	r.db.Conn(ctx).Raw(`
		SELECT COUNT(*)
		FROM payments p
		JOIN orders o ON o.id = p.order_id
		WHERE o.customer_id = ? AND p.status = 'failed' AND p.deleted_at IS NULL AND o.deleted_at IS NULL
	`, customerID).Scan(&failedPayments)

	return &entities.CustomerHistory{
		TotalOrders:       int(result.TotalOrders),
		TotalSpent:        result.TotalSpent,
		AvgOrderValue:     result.AvgOrderValue,
		FirstOrderDate:    result.FirstOrderDate,
		LastOrderDate:     result.LastOrderDate,
		NoveltyCount:      int(result.NoveltyCount),
		CODOrderCount:     int(result.CODOrderCount),
		DistinctAddresses: int(result.DistinctAddresses),
		FailedPayments:    int(failedPayments),
	}, nil
}

// GetCustomerDeliveryHistory obtiene el historial de entregas de un cliente.
//
// Tablas consultadas: shipments, orders (gestionadas por módulos shipments/orders)
func (r *Repository) GetCustomerDeliveryHistory(ctx context.Context, customerID uint) (*entities.DeliveryHistory, error) {
	var result struct {
		TotalShipments  int64
		FailedShipments int64
	}

	err := r.db.Conn(ctx).Raw(`
		SELECT
			COUNT(*) AS total_shipments,
			COUNT(CASE WHEN s.status = 'failed' THEN 1 END) AS failed_shipments
		FROM shipments s
		JOIN orders o ON o.id = s.order_id
		WHERE o.customer_id = ? AND s.deleted_at IS NULL AND o.deleted_at IS NULL
	`, customerID).Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &entities.DeliveryHistory{
		TotalShipments:  int(result.TotalShipments),
		FailedShipments: int(result.FailedShipments),
	}, nil
}

// GetOrderItemCount obtiene la cantidad total de items de una orden.
//
// Tabla consultada: order_items (gestionada por módulo orders)
func (r *Repository) GetOrderItemCount(ctx context.Context, orderID string) (int, error) {
	var count int64
	err := r.db.Conn(ctx).Raw(`
		SELECT COALESCE(SUM(quantity), 0)
		FROM order_items
		WHERE order_id = ? AND deleted_at IS NULL
	`, orderID).Scan(&count).Error
	return int(count), err
}

// GetPaymentMethodCategory obtiene la categoría de un método de pago por su ID.
//
// Tabla consultada: payment_methods (gestionada por módulo payments)
func (r *Repository) GetPaymentMethodCategory(ctx context.Context, paymentMethodID uint) (string, error) {
	if paymentMethodID == 0 {
		return "", nil
	}
	var category string
	err := r.db.Conn(ctx).Raw(`
		SELECT COALESCE(category, '') FROM payment_methods WHERE id = ?
	`, paymentMethodID).Scan(&category).Error
	return category, err
}
