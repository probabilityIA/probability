package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/primary/consumer/consumerevent/request"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// orderQueries implementa consultas de órdenes para el consumer
type orderQueries struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewOrderQueries crea una nueva instancia de consultas de órdenes
func NewOrderQueries(database db.IDatabase, logger log.ILogger) request.OrderRepository {
	return &orderQueries{
		db:     database,
		logger: logger.WithModule("order_queries"),
	}
}

// GetByID obtiene una orden por su ID
func (a *orderQueries) GetByID(ctx context.Context, id string) (*request.OrderData, error) {
	var order struct {
		ID              string
		OrderNumber     string
		Status          string
		PaymentMethodID uint
		CustomerPhone   string
		TotalAmount     float64
		Currency        string
		BusinessID      *uint
	}

	err := a.db.Conn(ctx).
		Table("orders").
		Select("id, order_number, status, payment_method_id, customer_phone, total_amount, currency, business_id").
		Where("id = ?", id).
		First(&order).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		a.logger.Error().Err(err).Str("order_id", id).Msg("Error getting order")
		return nil, err
	}

	return &request.OrderData{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		Status:          order.Status,
		PaymentMethodID: order.PaymentMethodID,
		CustomerPhone:   order.CustomerPhone,
		TotalAmount:     order.TotalAmount,
		Currency:        order.Currency,
		BusinessID:      order.BusinessID,
	}, nil
}
