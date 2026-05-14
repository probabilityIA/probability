package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/infra/secondary/client/request"
)

func BuildCreateJournalRequest(req *dtos.CreateJournalRequest) request.SiigoJournal {
	config := req.Config

	documentID := getIntFromConfig(config, "journal_document_id", 0)
	defaultWarehouseID := getIntFromConfig(config, "default_warehouse_id", 0)
	defaultAccountCode := getStringFromConfig(config, "default_account_code", "")
	defaultMovement := getStringFromConfig(config, "default_movement", "Debit")
	defaultCostCenter := getIntFromConfig(config, "default_cost_center", 0)
	defaultTaxID := getIntFromConfig(config, "journal_tax_id", 0)
	defaultTaxName := getStringFromConfig(config, "journal_tax_name", "")
	defaultTaxType := getStringFromConfig(config, "journal_tax_type", "")

	currency := req.Currency
	if currency == "" {
		currency = "COP"
	}

	date := req.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	items := make([]request.SiigoJournalItem, 0, len(req.Items))
	for _, item := range req.Items {
		accountCode := item.AccountCode
		if accountCode == "" {
			accountCode = defaultAccountCode
		}

		movement := item.Movement
		if movement == "" {
			movement = defaultMovement
		}

		warehouseID := item.WarehouseID
		if warehouseID == 0 {
			warehouseID = defaultWarehouseID
		}

		costCenter := item.CostCenter
		if costCenter == 0 {
			costCenter = defaultCostCenter
		}

		taxID := item.TaxID
		if taxID == 0 {
			taxID = defaultTaxID
		}

		journalItem := request.SiigoJournalItem{
			Account: request.SiigoJournalAccount{
				Code:     accountCode,
				Movement: movement,
			},
			Description: item.Name,
			Value:       item.TotalPrice,
			CostCenter:  costCenter,
		}

		if item.CustomerDNI != "" {
			journalItem.Customer = &request.SiigoJournalCustomer{
				Identification: item.CustomerDNI,
			}
		}

		if item.SKU != "" {
			journalItem.Product = &request.SiigoJournalProduct{
				Code:      item.SKU,
				Quantity:  item.Quantity,
				Warehouse: warehouseID,
			}
		}

		if taxID > 0 {
			tax := &request.SiigoJournalTax{ID: taxID}
			if item.TaxName != "" {
				tax.Name = item.TaxName
			} else if defaultTaxName != "" {
				tax.Name = defaultTaxName
			}
			if item.TaxType != "" {
				tax.Type = item.TaxType
			} else if defaultTaxType != "" {
				tax.Type = defaultTaxType
			}
			if item.TaxPercentage > 0 {
				tax.Percentage = item.TaxPercentage
			}
			journalItem.Tax = tax
		}

		items = append(items, journalItem)
	}

	journal := request.SiigoJournal{
		Document:     request.SiigoDocument{ID: documentID},
		Date:         date,
		Items:        items,
		Observations: req.Observations,
	}

	if currency != "COP" {
		journal.Currency = &request.SiigoCurrency{Code: currency}
	}

	return journal
}
