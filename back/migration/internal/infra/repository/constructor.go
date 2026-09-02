package repository

import (
	"context"

	"github.com/secamc93/probability/back/migration/shared/db"
	"github.com/secamc93/probability/back/migration/shared/env"
)

type Repository struct {
	db  db.IDatabase
	cfg env.IConfig
}

func New(db db.IDatabase, cfg env.IConfig) *Repository {
	return &Repository{
		db:  db,
		cfg: cfg,
	}
}

func (r *Repository) Migrate(ctx context.Context) error {
	return nil
}

func (r *Repository) migrateHistorico(ctx context.Context) error {
	if err := r.migrateUserTourProgress(ctx); err != nil {
		return err
	}
	if err := r.migrateWalletTxBusinessID(ctx); err != nil {
		return err
	}
	if err := r.migrateWalletTxConcept(ctx); err != nil {
		return err
	}
	if err := r.migrateWalletKPISelection(ctx); err != nil {
		return err
	}
	if err := r.migrateCatalogPricing(ctx); err != nil {
		return err
	}
	if err := r.migrateCodReport(ctx); err != nil {
		return err
	}
	if err := r.migrateShippingMarginCOD(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentCodMargin(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentCodRefactor(ctx); err != nil {
		return err
	}
	if err := r.migrateShipmentProbabilityGuide(ctx); err != nil {
		return err
	}
	if err := r.migrateGuideFormats(ctx); err != nil {
		return err
	}
	if err := r.migrateShippingQuotes(ctx); err != nil {
		return err
	}
	if err := r.migrateInvoicePartialUniqueIndex(ctx); err != nil {
		return err
	}
	if err := r.migrateIntegrationSyncRuns(ctx); err != nil {
		return err
	}
	if err := r.migrateProductChannelCategories(ctx); err != nil {
		return err
	}
	if err := r.migrateWarehouseLayout(ctx); err != nil {
		return err
	}
	if err := r.migrateWarehouseDimensions(ctx); err != nil {
		return err
	}
	if err := r.migrateRackSide(ctx); err != nil {
		return err
	}
	if err := r.migrateDemoAutoregistro(ctx); err != nil {
		return err
	}
	if err := r.migrateBusinessIsDemo(ctx); err != nil {
		return err
	}
	if err := r.migratePasswordResetTokens(ctx); err != nil {
		return err
	}
	if err := r.migrateWooShippingTokens(ctx); err != nil {
		return err
	}
	if err := r.migrateWooCommerceTestURL(ctx); err != nil {
		return err
	}
	if err := r.migrateChannelRawDataNullable(ctx); err != nil {
		return err
	}
	if err := r.migrateJumpsellerIntegrationType(ctx); err != nil {
		return err
	}
	if err := r.migrateJumpsellerOrderStatuses(ctx); err != nil {
		return err
	}
	if err := r.migrateTiendanubeOrderStatuses(ctx); err != nil {
		return err
	}
	if err := r.migrateVtexIntegrationType(ctx); err != nil {
		return err
	}
	if err := r.migrateShipitIntegrationType(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderGeoConfidence(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderIntegrationExternalUnique(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderIsCod(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderCodIncludesShipping(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderFreeShipping(ctx); err != nil {
		return err
	}
	if err := r.migrateWooCommerceStatusMappings(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderStatusSource(ctx); err != nil {
		return err
	}
	if err := r.backfillGeocodePendingOrders(ctx); err != nil {
		return err
	}
	if err := r.backfillOrdersGeozoneByPoint(ctx); err != nil {
		return err
	}
	if err := r.backfillOrdersGeozone(ctx); err != nil {
		return err
	}
	if err := r.migrateSubscriptionPlans(ctx); err != nil {
		return err
	}
	if err := r.migrateFreeTrialSubscriptionPlans(ctx); err != nil {
		return err
	}
	if err := r.migrateIntegrationStats(ctx); err != nil {
		return err
	}
	if err := r.migrateIntegrationsStoreIDUnique(ctx); err != nil {
		return err
	}
	if err := r.migrateAccounting(ctx); err != nil {
		return err
	}
	if err := r.migrateAccountingInvoices(ctx); err != nil {
		return err
	}
	if err := r.migrateFiscalProfiles(ctx); err != nil {
		return err
	}
	if err := r.migrateAccountingServices(ctx); err != nil {
		return err
	}
	if err := r.migrateNotificationConfigCODOnly(ctx); err != nil {
		return err
	}
	if err := r.migrateOrderStatusChangedAt(ctx); err != nil {
		return err
	}
	if err := r.migrateBusinessBarColors(ctx); err != nil {
		return err
	}
	if err := r.migrateBackfillBusinessSubscriptions(ctx); err != nil {
		return err
	}
	if err := r.backfillMysticMessageLogs(ctx); err != nil {
		return err
	}
	if err := r.fixVig0010Cod(ctx); err != nil {
		return err
	}
	if err := r.fixVig0095CotizacionFallida(ctx); err != nil {
		return err
	}
	if err := r.migrateTrazabilidadUsuario(ctx); err != nil {
		return err
	}
	if err := r.fixVigaCodRealPayout(ctx); err != nil {
		return err
	}
	if err := r.seedVigaCodMarginAmount(ctx); err != nil {
		return err
	}
	if err := r.fixVigaCodEnCurso(ctx); err != nil {
		return err
	}
	if err := r.fixVigaCodCalibracionFallida(ctx); err != nil {
		return err
	}
	if err := r.FixVigaCodPromesaCorte52(ctx); err != nil {
		return err
	}
	if err := r.FixVigaCodPromesaResto(ctx); err != nil {
		return err
	}
	if err := r.FixVigaCodRevertConfirmadas(ctx); err != nil {
		return err
	}
	if err := r.seedCodMarginAmount(ctx); err != nil {
		return err
	}
	if err := r.pruneOrderErrors(ctx); err != nil {
		return err
	}
	if err := r.migrateProductMatchRules(ctx); err != nil {
		return err
	}
	if err := r.migratePublicCheckout(ctx); err != nil {
		return err
	}
	if err := r.migrateWebsiteSectionsOrder(ctx); err != nil {
		return err
	}
	if err := r.migrateClientMultiBusiness(ctx); err != nil {
		return err
	}
	if err := r.migrateTiktokIntegrationType(ctx); err != nil {
		return err
	}
	if err := r.migrateMarketingLeads(ctx); err != nil {
		return err
	}
	if err := r.migratePaymentWebhookEvents(ctx); err != nil {
		return err
	}
	if err := r.migrateBancolombiaQR(ctx); err != nil {
		return err
	}
	if err := r.migrateNavbarContent(ctx); err != nil {
		return err
	}
	if err := r.migrateWhatsappConversationType(ctx); err != nil {
		return err
	}
	if err := r.migrateCommercialProspects(ctx); err != nil {
		return err
	}
	if err := r.migrateSiigoReferrals(ctx); err != nil {
		return err
	}
	if err := r.migrateSubscriptionAuditLogs(ctx); err != nil {
		return err
	}
	if err := r.migrateBusinessModuleOverrideExpiry(ctx); err != nil {
		return err
	}
	if err := r.migrateTiendanubeURLs(ctx); err != nil {
		return err
	}
	if err := r.migrateSubscriptionCutoffDay(ctx); err != nil {
		return err
	}
	if err := r.migrateSubscriptionCourtesyUntil(ctx); err != nil {
		return err
	}
	return r.seedCommercialProspects(ctx)
}
