package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/request"
)

// BuildCreateJournalRequest construye el body de la request de Siigo para crear un comprobante contable
func BuildCreateJournalRequest(req *dtos.CreateJournalRequest) request.SiigoJournal {
	config := req.Config

	// Document ID (tipo de documento CC desde la config de la integración)
	documentID := getIntFromConfig(config, "journal_document_id", 0)

	// Defaults desde config
	defaultWarehouseID := getIntFromConfig(config, "default_warehouse_id", 0)
	defaultAccountCode := getStringFromConfig(config, "default_account_code", "")
	defaultMovement := getStringFromConfig(config, "default_movement", "Debit")
	defaultCostCenter := getIntFromConfig(config, "default_cost_center", 0)
	defaultTaxID := getIntFromConfig(config, "journal_tax_id", 0)

	// Moneda
	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	// Fecha
	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Construir items
	items := make([]request.SiigoJournalItem, 0, len(req.Items))
	for _, item := range req.Items {
		// Account code: item override -> config default
		accountCode := item.AccountCode
		if accountCode == "" {
			accountCode = defaultAccountCode
		}

		// Movement: item override -> config default
		movement := item.Movement
		if movement == "" {
			movement = defaultMovement
		}

		// Warehouse: item override -> config default
		warehouseID := item.WarehouseID
		if warehouseID == 0 {
			warehouseID = defaultWarehouseID
		}

		// Cost center: item override -> config default
		costCenter := item.CostCenter
		if costCenter == 0 {
			costCenter = defaultCostCenter
		}

		// Tax: item override -> config default
		taxID := item.TaxID
		if taxID == 0 {
			taxID = defaultTaxID
		}

		journalItem := request.SiigoJournalItem{
			Account:     request.SiigoJournalAccount{Code: accountCode},
			Description: item.Name,
			Movement:    movement,
			Value:       item.TotalPrice,
			CostCenter:  costCenter,
		}

		// Customer (si tiene DNI)
		if item.CustomerDNI != "" {
			journalItem.Customer = &request.SiigoJournalCustomer{
				Identification: item.CustomerDNI,
			}
		}

		// Product (si tiene SKU)
		if item.SKU != "" {
			journalItem.Product = &request.SiigoJournalProduct{
				Code:      item.SKU,
				Quantity:  item.Quantity,
				Warehouse: warehouseID,
			}
		}

		// Tax
		if taxID > 0 {
			journalItem.Taxes = []request.SiigoTax{{ID: taxID}}
		}

		items = append(items, journalItem)
	}

	return request.SiigoJournal{
		Document:     request.SiigoDocument{ID: documentID},
		Date:         date,
		Currency:     request.SiigoCurrency{Code: currency},
		Items:        items,
		Observations: req.Observations,
	}
}
