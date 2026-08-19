package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

func (uc *invoicingUseCase) ListCatalogs(ctx context.Context, integrationID uint) (*dtos.Catalogs, error) {
	if integrationID == 0 {
		return nil, fmt.Errorf("integration_id es requerido")
	}

	creds, err := uc.resolveWebhookCredentials(ctx, fmt.Sprintf("%d", integrationID))
	if err != nil {
		return nil, err
	}
	credentials := &creds

	catalogs := &dtos.Catalogs{}
	var mu sync.Mutex
	var wg sync.WaitGroup

	registrar := func(nombre string, err error) {
		if err == nil {
			return
		}
		mu.Lock()
		catalogs.Errors = append(catalogs.Errors, fmt.Sprintf("%s: %v", nombre, err))
		mu.Unlock()
		uc.log.Warn(ctx).Err(err).Str("catalogo", nombre).Msg("No se pudo consultar un catalogo de Siigo")
	}

	tareas := []struct {
		nombre  string
		ejecuta func() error
	}{
		{"tipos de documento FV", func() error {
			items, err := uc.siigoClient.ListDocumentTypes(ctx, *credentials, "FV")
			catalogs.DocumentTypesFV = items
			return err
		}},
		{"tipos de documento NC", func() error {
			items, err := uc.siigoClient.ListDocumentTypes(ctx, *credentials, "NC")
			catalogs.DocumentTypesNC = items
			return err
		}},
		{"tipos de documento RC", func() error {
			items, err := uc.siigoClient.ListDocumentTypes(ctx, *credentials, "RC")
			catalogs.DocumentTypesRC = items
			return err
		}},
		{"tipos de documento CC", func() error {
			items, err := uc.siigoClient.ListDocumentTypes(ctx, *credentials, "CC")
			catalogs.DocumentTypesCC = items
			return err
		}},
		{"medios de pago FV", func() error {
			items, err := uc.siigoClient.ListPaymentTypes(ctx, *credentials, "FV")
			catalogs.PaymentTypesFV = mapPaymentTypes(items)
			return err
		}},
		{"medios de pago RC", func() error {
			items, err := uc.siigoClient.ListPaymentTypes(ctx, *credentials, "RC")
			catalogs.PaymentTypesRC = mapPaymentTypes(items)
			return err
		}},
		{"vendedores", func() error {
			items, err := uc.siigoClient.ListSellers(ctx, *credentials)
			catalogs.Sellers = items
			return err
		}},
		{"impuestos", func() error {
			items, err := uc.siigoClient.ListTaxes(ctx, *credentials)
			catalogs.Taxes = items
			return err
		}},
		{"centros de costo", func() error {
			items, err := uc.siigoClient.ListCostCenters(ctx, *credentials)
			catalogs.CostCenters = items
			return err
		}},
		{"bodegas", func() error {
			items, err := uc.siigoClient.ListWarehouses(ctx, *credentials)
			catalogs.Warehouses = mapWarehouses(items)
			return err
		}},
	}

	for _, tarea := range tareas {
		wg.Add(1)
		go func(nombre string, ejecuta func() error) {
			defer wg.Done()
			registrar(nombre, ejecuta())
		}(tarea.nombre, tarea.ejecuta)
	}

	wg.Wait()

	uc.log.Info(ctx).
		Uint("integration_id", integrationID).
		Int("con_error", len(catalogs.Errors)).
		Msg("Catalogos de Siigo consultados")

	return catalogs, nil
}

func mapPaymentTypes(items []dtos.PaymentTypeItem) []dtos.CatalogItem {
	out := make([]dtos.CatalogItem, 0, len(items))
	for _, i := range items {
		out = append(out, dtos.CatalogItem{ID: i.ID, Name: i.Name, Detail: i.Type})
	}
	return out
}

func mapWarehouses(items []dtos.WarehouseItem) []dtos.CatalogItem {
	out := make([]dtos.CatalogItem, 0, len(items))
	for _, i := range items {
		out = append(out, dtos.CatalogItem{ID: i.ID, Name: i.Name})
	}
	return out
}
