package worker

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

// OrderScheduler programa la generación automática de órdenes cada 5 minutos
type OrderScheduler struct {
	useCases      *usecases.UseCases
	logger        log.ILogger
	config        *SchedulerConfig
	ticker        *time.Ticker
	stopChan      chan bool
	isRunning     bool
	platformIndex int // Índice para rotar entre plataformas
}

// SchedulerConfig contiene la configuración para la generación automática
type SchedulerConfig struct {
	Interval        time.Duration // Intervalo entre generaciones (default: 5 minutos)
	OrdersPerBatch  int           // Cantidad de órdenes a generar por batch (default: 1)
	IntegrationID   uint          // ID de la integración
	BusinessID      *uint         // ID del negocio (opcional)
	Platform        string        // Plataforma (default: "test")
	Status          string        // Estado inicial (default: "pending")
	IncludePayment  bool          // Si incluir información de pago
	IncludeShipment bool          // Si incluir información de envío
}

// NewOrderScheduler crea una nueva instancia del scheduler
func NewOrderScheduler(uc *usecases.UseCases, logger log.ILogger, config *SchedulerConfig) *OrderScheduler {
	// Valores por defecto
	if config.Interval == 0 {
		config.Interval = 5 * time.Minute
	}
	if config.OrdersPerBatch == 0 {
		config.OrdersPerBatch = 1
	}
	if config.Platform == "" {
		config.Platform = "test"
	}
	if config.Status == "" {
		config.Status = "pending"
	}

	return &OrderScheduler{
		useCases:      uc,
		logger:        logger,
		config:        config,
		stopChan:      make(chan bool),
		isRunning:     false,
		platformIndex: 0,
	}
}

// Start inicia el scheduler en una goroutine
func (s *OrderScheduler) Start(ctx context.Context) {
	if s.isRunning {
		s.logger.Warn().Msg("Order scheduler is already running")
		return
	}

	s.isRunning = true
	s.ticker = time.NewTicker(s.config.Interval)

	s.logger.Info().
		Dur("interval", s.config.Interval).
		Int("orders_per_batch", s.config.OrdersPerBatch).
		Msg("Order scheduler started - will generate orders automatically")

	// Ejecutar inmediatamente al iniciar
	go s.generateOrders(ctx)

	// Ejecutar periódicamente
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.generateOrders(ctx)
			case <-s.stopChan:
				s.ticker.Stop()
				s.isRunning = false
				s.logger.Info().Msg("Order scheduler stopped")
				return
			case <-ctx.Done():
				s.Stop()
				return
			}
		}
	}()
}

// Stop detiene el scheduler
func (s *OrderScheduler) Stop() {
	if !s.isRunning {
		return
	}
	s.stopChan <- true
}

// generateOrders genera las órdenes según la configuración
// Genera una orden de cada plataforma (Shopify, Meli, WooCommerce) en cada ciclo
func (s *OrderScheduler) generateOrders(ctx context.Context) {
	// Definir las tres plataformas con sus respectivos integration IDs
	platforms := []struct {
		Platform      string
		IntegrationID uint
	}{
		{"shopify", 1},
		{"meli", 5},
		{"woocommerce", 6},
	}

	// Generar una orden de cada plataforma
	for _, platform := range platforms {
		req := &domain.GenerateOrderRequest{
			Count:           1,
			IntegrationID:   platform.IntegrationID,
			BusinessID:      s.config.BusinessID,
			Platform:        platform.Platform,
			Status:          s.config.Status,
			IncludePayment:  s.config.IncludePayment,
			IncludeShipment: s.config.IncludeShipment,
		}

		response, err := s.useCases.GenerateAndPublishOrders(ctx, req)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("platform", platform.Platform).
				Uint("integration_id", platform.IntegrationID).
				Int("generated", response.Generated).
				Int("published", response.Published).
				Int("failed", response.Failed).
				Msg("Error generating orders in scheduler")
			continue
		}

		s.logger.Info().
			Str("platform", platform.Platform).
			Uint("integration_id", platform.IntegrationID).
			Int("generated", response.Generated).
			Int("published", response.Published).
			Int("failed", response.Failed).
			Msg("Orders generated automatically by scheduler")
	}
}
