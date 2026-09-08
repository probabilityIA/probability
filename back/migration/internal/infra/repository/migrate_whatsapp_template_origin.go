package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

var systemTemplateNames = []string{
	"confirmacion_pedido",
	"confirmacion_pedido_contraentrega",
	"confirmacion_pedido_contraentrega_sin_valor",
	"confirmar_cancelacion_pedido",
	"guia_envio_generada",
	"guia_envio_generada_cod",
	"handoff_asesor",
	"menu_no_confirmacion",
	"motivo_cancelacion_pedido",
	"novedad_cambio_direccion",
	"novedad_cambio_medio_pago",
	"novedad_cambio_productos",
	"pedido_cancelado",
	"pedido_confirmado_v2",
	"pedido_en_reparto",
	"pedido_en_reparto_cod",
	"pedido_entregado",
	"pedido_entregado_cod",
	"tipo_novedad_pedido",
}

var internalTemplateNames = []string{
	"alerta_servidor",
	"recuperacion_codigo",
	"reporte_saldo_billetera",
	"reporte_saldo_billetera_v2",
	"resumen_pago_suscripcion",
	"prueba_conexion",
	"hello_world",
}

func (r *Repository) migrateWhatsappTemplateOrigin(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	if err := conn.AutoMigrate(&models.WhatsappTemplate{}); err != nil {
		return fmt.Errorf("failed to auto-migrate whatsapp_templates origin: %w", err)
	}

	if err := conn.Exec(
		`ALTER TABLE whatsapp_templates ALTER COLUMN business_id DROP NOT NULL`,
	).Error; err != nil {
		return fmt.Errorf("failed to make whatsapp_templates.business_id nullable: %w", err)
	}

	seed := func(names []string, origin, scope string) error {
		for _, name := range names {
			if err := conn.Exec(`
				INSERT INTO whatsapp_templates
					(created_at, updated_at, business_id, origin, scope, name, language, category,
					 body_text, header_text, footer_text, status)
				SELECT NOW(), NOW(), NULL, ?, ?, ?, 'es', 'UTILITY', '', '', '', 'approved'
				WHERE NOT EXISTS (
					SELECT 1 FROM whatsapp_templates
					WHERE name = ? AND business_id IS NULL AND deleted_at IS NULL
				)`, origin, scope, name, name).Error; err != nil {
				return fmt.Errorf("failed to seed template %s: %w", name, err)
			}
		}
		return nil
	}

	if err := seed(systemTemplateNames, models.WhatsappTemplateOriginSystem, models.WhatsappTemplateScopeOrderEvent); err != nil {
		return err
	}

	if err := seed(internalTemplateNames, models.WhatsappTemplateOriginInternal, models.WhatsappTemplateScopeInternal); err != nil {
		return err
	}

	if err := conn.Exec(
		`UPDATE whatsapp_templates SET scope = ? WHERE origin = ? AND scope <> ?`,
		models.WhatsappTemplateScopeOrderEvent,
		models.WhatsappTemplateOriginSystem,
		models.WhatsappTemplateScopeOrderEvent,
	).Error; err != nil {
		return fmt.Errorf("failed to backfill system scope: %w", err)
	}

	if err := conn.Exec(
		`UPDATE whatsapp_templates SET scope = ? WHERE origin = ? AND scope <> ?`,
		models.WhatsappTemplateScopeInternal,
		models.WhatsappTemplateOriginInternal,
		models.WhatsappTemplateScopeInternal,
	).Error; err != nil {
		return fmt.Errorf("failed to backfill internal scope: %w", err)
	}

	return nil
}
