package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	models "github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Wallet

func (r *Repository) GetWalletByBusinessID(ctx context.Context, businessID uint) (*entities.Wallet, error) {
	var m models.Wallet
	err := r.db.Conn(ctx).Where("business_id = ?", businessID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return walletToDomain(&m), nil
}

func (r *Repository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*entities.Wallet, error) {
	var m models.Wallet
	err := r.db.Conn(ctx).Where("id = ?", walletID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet not found")
		}
		return nil, err
	}
	return walletToDomain(&m), nil
}

func (r *Repository) CreateWallet(ctx context.Context, wallet *entities.Wallet) error {
	m := walletToModel(wallet)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	wallet.ID = m.ID
	wallet.CreatedAt = m.CreatedAt
	wallet.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *Repository) UpdateWallet(ctx context.Context, wallet *entities.Wallet) error {
	m := walletToModel(wallet)
	if err := r.db.Conn(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}
	return nil
}

func (r *Repository) GetAllWallets(ctx context.Context) ([]*entities.Wallet, error) {
	var list []models.Wallet
	if err := r.db.Conn(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	result := make([]*entities.Wallet, len(list))
	for i, m := range list {
		result[i] = walletToDomain(&m)
	}
	return result, nil
}

// WalletTransactions

func (r *Repository) CreateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	m := walletTxToModel(tx)
	if err := r.db.Conn(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("failed to create wallet transaction: %w", err)
	}
	tx.ID = m.ID
	tx.CreatedAt = m.CreatedAt
	return nil
}

func (r *Repository) GetWalletTransactionByID(ctx context.Context, id uuid.UUID) (*entities.WalletTransaction, error) {
	var m models.WalletTransaction
	err := r.db.Conn(ctx).Where("id = ?", id).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("wallet transaction not found")
		}
		return nil, err
	}
	return walletTxToDomain(&m), nil
}

func (r *Repository) GetWalletTransactionByReference(ctx context.Context, reference string) (*entities.WalletTransaction, error) {
	var m models.WalletTransaction
	err := r.db.Conn(ctx).Where("reference = ?", reference).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return walletTxToDomain(&m), nil
}

func (r *Repository) SaveWalletTransactionGatewayResponse(ctx context.Context, id uuid.UUID, response []byte) error {
	return r.db.Conn(ctx).
		Table("transaction").
		Where("id = ?", id).
		Update("gateway_response", datatypes.JSON(response)).Error
}

func (r *Repository) UpdateWalletTransaction(ctx context.Context, tx *entities.WalletTransaction) error {
	m := walletTxToModel(tx)
	if err := r.db.Conn(ctx).Save(m).Error; err != nil {
		return fmt.Errorf("failed to update wallet transaction: %w", err)
	}
	return nil
}

func (r *Repository) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) ([]*entities.WalletTransaction, error) {
	type row struct {
		models.WalletTransaction
		IntegrationImageURL string `gorm:"column:integration_image_url"`
	}
	var rows []row
	err := r.db.Conn(ctx).
		Table("transaction AS t").
		Select("t.*, it.image_url AS integration_image_url").
		Joins("LEFT JOIN integration_types it ON it.id = t.integration_type_id").
		Where("t.wallet_id = ?", walletID).
		Order("t.created_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	imageBase := strings.TrimRight(os.Getenv("URL_BASE_DOMAIN_S3"), "/")
	out := make([]*entities.WalletTransaction, len(rows))
	for i := range rows {
		e := walletTxToDomain(&rows[i].WalletTransaction)
		raw := rows[i].IntegrationImageURL
		if raw != "" && imageBase != "" && !strings.HasPrefix(raw, "http") {
			e.IntegrationImageURL = imageBase + "/" + strings.TrimLeft(raw, "/")
		} else {
			e.IntegrationImageURL = raw
		}
		out[i] = e
	}
	return out, nil
}

func (r *Repository) GetPendingRechargeTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	var list []models.WalletTransaction
	err := r.db.Conn(ctx).
		Where("status = ? AND type = ?", entities.WalletTxStatusPending, entities.WalletTxTypeRecharge).
		Order("created_at ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return walletTxListToDomain(list), nil
}

func (r *Repository) GetProcessedTransactions(ctx context.Context) ([]*entities.WalletTransaction, error) {
	var list []models.WalletTransaction
	err := r.db.Conn(ctx).
		Where("status IN ?", []string{entities.WalletTxStatusCompleted, entities.WalletTxStatusFailed}).
		Order("created_at DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return walletTxListToDomain(list), nil
}

func (r *Repository) DeleteTransactionsByWalletIDAndType(ctx context.Context, walletID uuid.UUID, txType string) error {
	return r.db.Conn(ctx).
		Where("wallet_id = ? AND type = ?", walletID, txType).
		Delete(&models.WalletTransaction{}).Error
}

func (r *Repository) DeleteAllTransactionsByWalletID(ctx context.Context, walletID uuid.UUID) error {
	return r.db.Conn(ctx).
		Where("wallet_id = ?", walletID).
		Delete(&models.WalletTransaction{}).Error
}

// Wallet Mappers

func walletToModel(e *entities.Wallet) *models.Wallet {
	return &models.Wallet{
		ID:         e.ID,
		BusinessID: e.BusinessID,
		Balance:    e.Balance,
		CreatedAt:  e.CreatedAt,
		UpdatedAt:  e.UpdatedAt,
	}
}

func walletToDomain(m *models.Wallet) *entities.Wallet {
	return &entities.Wallet{
		ID:         m.ID,
		BusinessID: m.BusinessID,
		Balance:    m.Balance,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func walletTxToModel(e *entities.WalletTransaction) *models.WalletTransaction {
	return &models.WalletTransaction{
		ID:                   e.ID,
		WalletID:             e.WalletID,
		Amount:               e.Amount,
		Type:                 e.Type,
		Status:               e.Status,
		Reference:            e.Reference,
		QrCode:               e.QrCode,
		PaymentTransactionID: e.PaymentTransactionID,
		UserID:               e.UserID,
		IntegrationTypeID:    e.IntegrationTypeID,
		IntegrationID:        e.IntegrationID,
		GatewayRequest:       datatypes.JSON(e.GatewayRequest),
		GatewayResponse:      datatypes.JSON(e.GatewayResponse),
		CreatedAt:            e.CreatedAt,
	}
}

func walletTxToDomain(m *models.WalletTransaction) *entities.WalletTransaction {
	return &entities.WalletTransaction{
		ID:                   m.ID,
		WalletID:             m.WalletID,
		Amount:               m.Amount,
		Type:                 m.Type,
		Status:               m.Status,
		Reference:            m.Reference,
		QrCode:               m.QrCode,
		PaymentTransactionID: m.PaymentTransactionID,
		UserID:               m.UserID,
		IntegrationTypeID:    m.IntegrationTypeID,
		IntegrationID:        m.IntegrationID,
		GatewayRequest:       []byte(m.GatewayRequest),
		GatewayResponse:      []byte(m.GatewayResponse),
		CreatedAt:            m.CreatedAt,
	}
}

func walletTxListToDomain(list []models.WalletTransaction) []*entities.WalletTransaction {
	result := make([]*entities.WalletTransaction, len(list))
	for i, m := range list {
		result[i] = walletTxToDomain(&m)
	}
	return result
}

// Financial Stats

func (r *Repository) GetFinancialStats(ctx context.Context, dto *dtos.FinancialStatsDTO) (*dtos.FinancialStatsResponse, error) {
	// Parsear fechas
	startDate, err := time.Parse("2006-01-02", dto.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", dto.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format: %w", err)
	}
	endDate = endDate.Add(24 * time.Hour) // Incluir todo el día final

	// Obtener suscripciones pagadas
	type SubscriptionStats struct {
		BusinessID uint
		Amount     float64
	}
	var subscriptions []SubscriptionStats
	subQuery := r.db.Conn(ctx).
		Table("business_subscriptions").
		Select("business_id, SUM(amount) as amount").
		Where("status = ? AND created_at BETWEEN ? AND ?", "paid", startDate, endDate)
	if dto.BusinessID != nil {
		subQuery = subQuery.Where("business_id = ?", *dto.BusinessID)
	}
	err = subQuery.Group("business_id").Scan(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription stats: %w", err)
	}

	// Obtener cantidad de guías (shipments no canceladas)
	// Usar DISTINCT ON(order_id) para contar solo el shipment más reciente por orden (igual que en dashboard)
	type GuideStats struct {
		BusinessID uint
		GuideCount int
	}
	var guides []GuideStats

	baseSQLGuides := `
		SELECT o.business_id, COUNT(*) as guide_count
		FROM (
			SELECT DISTINCT ON (order_id) id, order_id, status
			FROM shipments
			WHERE status != ? AND created_at BETWEEN ? AND ?
			ORDER BY order_id, created_at DESC
		) s
		JOIN orders o ON o.id = s.order_id
		WHERE o.deleted_at IS NULL
	`

	guideParams := []interface{}{
		"cancelled", startDate, endDate,
	}

	if dto.BusinessID != nil {
		baseSQLGuides += " AND o.business_id = ?"
		guideParams = append(guideParams, *dto.BusinessID)
	}

	baseSQLGuides += " GROUP BY o.business_id"

	err = r.db.Conn(ctx).Raw(baseSQLGuides, guideParams...).Scan(&guides).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get guide stats: %w", err)
	}

	// Combinar datos y obtener nombres de negocios
	statsMap := make(map[uint]*dtos.BusinessFinancialStats)

	for _, sub := range subscriptions {
		if _, exists := statsMap[sub.BusinessID]; !exists {
			statsMap[sub.BusinessID] = &dtos.BusinessFinancialStats{
				BusinessID: sub.BusinessID,
			}
		}
		statsMap[sub.BusinessID].SubscriptionIncome = sub.Amount
	}

	for _, guide := range guides {
		if _, exists := statsMap[guide.BusinessID]; !exists {
			statsMap[guide.BusinessID] = &dtos.BusinessFinancialStats{
				BusinessID: guide.BusinessID,
			}
		}
		statsMap[guide.BusinessID].GuideCount = guide.GuideCount
		statsMap[guide.BusinessID].GuideIncome = float64(guide.GuideCount) * 2290 // Precio fijo
	}

	// Obtener nombres de negocios (solo para los negocios en statsMap)
	type BusinessName struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var businesses []BusinessName

	// Extraer IDs de negocios del statsMap
	businessIDs := make([]uint, 0)
	for id := range statsMap {
		businessIDs = append(businessIDs, id)
	}

	businessNameQuery := r.db.Conn(ctx).
		Table("business").
		Select("id, name")

	if len(businessIDs) > 0 {
		businessNameQuery = businessNameQuery.Where("id IN ?", businessIDs)
	}

	err = businessNameQuery.Scan(&businesses).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get business names: %w", err)
	}

	businessNameMap := make(map[uint]string)
	for _, b := range businesses {
		businessNameMap[b.ID] = b.Name
	}

	// Construir respuesta
	var businessList []dtos.BusinessFinancialStats
	var totalIncome float64
	for _, stats := range statsMap {
		stats.BusinessName = businessNameMap[stats.BusinessID]
		stats.TotalIncome = stats.SubscriptionIncome + stats.GuideIncome
		totalIncome += stats.TotalIncome
		businessList = append(businessList, *stats)
	}

	return &dtos.FinancialStatsResponse{
		Period: dtos.PeriodInfo{
			Start: dto.StartDate,
			End:   dto.EndDate,
		},
		TotalIncome: totalIncome,
		Businesses:  businessList,
	}, nil
}
