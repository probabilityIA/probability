package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/secamc93/probability/back/integrationTest/integrations/shopify"
	"github.com/secamc93/probability/back/integrationTest/integrations/softpymes"
	"github.com/secamc93/probability/back/integrationTest/integrations/whatsapp"
	"github.com/secamc93/probability/back/integrationTest/shared/env"
	"github.com/secamc93/probability/back/integrationTest/shared/log"
)

func main() {
	logger := log.New()
	config := env.New()

	webhookBaseURL := config.Get("WEBHOOK_BASE_URL")
	if webhookBaseURL == "" {
		logger.Fatal().Msg("WEBHOOK_BASE_URL no configurado en .env")
		os.Exit(1)
	}

	shopDomain := config.Get("SHOPIFY_SHOP_DOMAIN")
	if shopDomain == "" {
		logger.Fatal().Msg("SHOPIFY_SHOP_DOMAIN no configurado en .env")
		os.Exit(1)
	}

	logger.Info().
		Str("webhook_base_url", webhookBaseURL).
		Str("shop_domain", shopDomain).
		Msg("Inicializando simulador de webhooks")

	shopifyIntegration := shopify.New(config, logger)
	whatsappIntegration := whatsapp.New(config, logger)
	softpymesIntegration := softpymes.New(logger)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Simulador de Webhooks ===")
		fmt.Println("\n📦 SHOPIFY:")
		fmt.Println("1. orders/create (crear nueva orden aleatoria)")
		fmt.Println("2. orders/paid (marcar orden como pagada)")
		fmt.Println("3. orders/updated (actualizar orden existente)")
		fmt.Println("4. orders/cancelled (cancelar orden existente)")
		fmt.Println("5. orders/fulfilled (marcar orden como cumplida)")
		fmt.Println("6. orders/partially_fulfilled (marcar orden como parcialmente cumplida)")
		fmt.Println("7. Listar órdenes almacenadas")
		fmt.Println("\n💬 WHATSAPP:")
		fmt.Println("8. Simular respuesta de usuario (manual)")
		fmt.Println("9. Simular respuesta automática (por template)")
		fmt.Println("10. Listar conversaciones almacenadas")
		fmt.Println("\n📄 SOFTPYMES (Facturación):")
		fmt.Println("11. Simular autenticación")
		fmt.Println("12. Simular creación de factura")
		fmt.Println("13. Simular nota de crédito")
		fmt.Println("14. Listar facturas almacenadas")
		fmt.Println("\n0. Salir")
		fmt.Print("\nOpción: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1", "2", "3", "4", "5", "6":
			var topic string
			switch input {
			case "1":
				topic = "orders/create"
			case "2":
				topic = "orders/paid"
			case "3":
				topic = "orders/updated"
			case "4":
				topic = "orders/cancelled"
			case "5":
				topic = "orders/fulfilled"
			case "6":
				topic = "orders/partially_fulfilled"
			}

			logger.Info().Str("topic", topic).Msg("Simulando webhook")

			if err := shopifyIntegration.SimulateOrder(topic); err != nil {
				logger.Error().Err(err).Str("topic", topic).Msg("Error al simular webhook")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().Str("topic", topic).Msg("Webhook simulado exitosamente")
				fmt.Printf("✅ Webhook '%s' enviado exitosamente\n", topic)
			}
		case "7":
			orders := shopifyIntegration.GetAllOrders()
			if len(orders) == 0 {
				fmt.Println("📭 No hay órdenes almacenadas")
			} else {
				fmt.Printf("\n📦 Órdenes almacenadas (%d):\n", len(orders))
				for i, order := range orders {
					status := order.FinancialStatus
					if order.FulfillmentStatus != nil {
						status += " / " + *order.FulfillmentStatus
					}
					fmt.Printf("  %d. %s - %s - %s %s - Estado: %s\n",
						i+1, order.Name, order.Email, order.Currency, order.TotalPrice, status)
				}
			}
		case "8":
			fmt.Print("Número de teléfono (ej: +573001234567): ")
			phoneInput, _ := reader.ReadString('\n')
			phoneNumber := strings.TrimSpace(phoneInput)

			fmt.Print("Respuesta del usuario (ej: Confirmar pedido): ")
			responseInput, _ := reader.ReadString('\n')
			response := strings.TrimSpace(responseInput)

			logger.Info().
				Str("phone_number", phoneNumber).
				Str("response", response).
				Msg("Simulando respuesta de usuario")

			if err := whatsappIntegration.SimulateUserResponse(phoneNumber, response); err != nil {
				logger.Error().Err(err).Msg("Error al simular respuesta de usuario")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().Msg("Respuesta de usuario simulada exitosamente")
				fmt.Printf("✅ Respuesta '%s' enviada al sistema para %s\n", response, phoneNumber)
			}
		case "9":
			fmt.Print("Número de teléfono (ej: +573001234567): ")
			phoneInput, _ := reader.ReadString('\n')
			phoneNumber := strings.TrimSpace(phoneInput)

			fmt.Print("Nombre del template (ej: confirmacion_pedido_contraentrega): ")
			templateInput, _ := reader.ReadString('\n')
			templateName := strings.TrimSpace(templateInput)

			logger.Info().
				Str("phone_number", phoneNumber).
				Str("template", templateName).
				Msg("Simulando respuesta automática")

			if err := whatsappIntegration.SimulateAutoResponse(phoneNumber, templateName); err != nil {
				logger.Error().Err(err).Msg("Error al simular respuesta automática")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().Msg("Respuesta automática simulada exitosamente")
				fmt.Printf("✅ Respuesta automática enviada al sistema para %s\n", phoneNumber)
			}
		case "10":
			conversations := whatsappIntegration.GetAllConversations()
			if len(conversations) == 0 {
				fmt.Println("📭 No hay conversaciones almacenadas")
			} else {
				fmt.Printf("\n💬 Conversaciones almacenadas (%d):\n", len(conversations))
				for i, conv := range conversations {
					messages := whatsappIntegration.GetMessages(conv.ID)
					fmt.Printf("  %d. %s - Estado: %s - Orden: %s - Mensajes: %d\n",
						i+1, conv.PhoneNumber, conv.CurrentState, conv.OrderNumber, len(messages))
				}
			}
		case "11":
			fmt.Print("API Key: ")
			apiKeyInput, _ := reader.ReadString('\n')
			apiKey := strings.TrimSpace(apiKeyInput)

			fmt.Print("API Secret: ")
			apiSecretInput, _ := reader.ReadString('\n')
			apiSecret := strings.TrimSpace(apiSecretInput)

			fmt.Print("Referer (ej: https://tutienda.com): ")
			refererInput, _ := reader.ReadString('\n')
			referer := strings.TrimSpace(refererInput)

			logger.Info().Msg("Simulando autenticación de SoftPymes")

			token, err := softpymesIntegration.SimulateAuth(apiKey, apiSecret, referer)
			if err != nil {
				logger.Error().Err(err).Msg("Error al autenticar")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().Str("token", token).Msg("Autenticación exitosa")
				fmt.Printf("✅ Token generado: %s\n", token)
				fmt.Println("💡 Guarda este token para crear facturas")
			}
		case "12":
			fmt.Print("Token (obtenido en opción 11): ")
			tokenInput, _ := reader.ReadString('\n')
			token := strings.TrimSpace(tokenInput)

			fmt.Print("Order ID (ej: ORD-001): ")
			orderIDInput, _ := reader.ReadString('\n')
			orderID := strings.TrimSpace(orderIDInput)

			fmt.Print("Nombre cliente: ")
			customerNameInput, _ := reader.ReadString('\n')
			customerName := strings.TrimSpace(customerNameInput)

			fmt.Print("Email cliente: ")
			customerEmailInput, _ := reader.ReadString('\n')
			customerEmail := strings.TrimSpace(customerEmailInput)

			fmt.Print("NIT cliente: ")
			customerNITInput, _ := reader.ReadString('\n')
			customerNIT := strings.TrimSpace(customerNITInput)

			fmt.Print("Total (ej: 100000): ")
			totalInput, _ := reader.ReadString('\n')
			totalStr := strings.TrimSpace(totalInput)
			var total float64
			fmt.Sscanf(totalStr, "%f", &total)

			invoiceData := map[string]interface{}{
				"order_id": orderID,
				"customer": map[string]interface{}{
					"name":  customerName,
					"email": customerEmail,
					"nit":   customerNIT,
				},
				"items": []interface{}{
					map[string]interface{}{
						"description": "Producto Test",
						"quantity":    1.0,
						"unit_price":  total,
						"tax":         total * 0.19,
						"total":       total * 1.19,
					},
				},
				"total": total,
			}

			logger.Info().Msg("Simulando creación de factura")

			invoice, err := softpymesIntegration.SimulateInvoice(token, invoiceData)
			if err != nil {
				logger.Error().Err(err).Msg("Error al crear factura")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().
					Str("invoice_number", invoice.InvoiceNumber).
					Str("cufe", invoice.CUFE).
					Msg("Factura creada exitosamente")
				fmt.Printf("✅ Factura creada:\n")
				fmt.Printf("  Número: %s\n", invoice.InvoiceNumber)
				fmt.Printf("  CUFE: %s\n", invoice.CUFE)
				fmt.Printf("  Total: $%.2f %s\n", invoice.Total, invoice.Currency)
				fmt.Printf("  PDF: %s\n", invoice.PDFURL)
			}
		case "13":
			fmt.Print("Token: ")
			tokenInput, _ := reader.ReadString('\n')
			token := strings.TrimSpace(tokenInput)

			fmt.Print("Invoice ID (external_id de la factura): ")
			invoiceIDInput, _ := reader.ReadString('\n')
			invoiceID := strings.TrimSpace(invoiceIDInput)

			fmt.Print("Monto a acreditar: ")
			amountInput, _ := reader.ReadString('\n')
			amountStr := strings.TrimSpace(amountInput)
			var amount float64
			fmt.Sscanf(amountStr, "%f", &amount)

			fmt.Print("Razón (ej: Devolución de producto): ")
			reasonInput, _ := reader.ReadString('\n')
			reason := strings.TrimSpace(reasonInput)

			fmt.Print("Tipo (total/partial): ")
			noteTypeInput, _ := reader.ReadString('\n')
			noteType := strings.TrimSpace(noteTypeInput)

			creditNoteData := map[string]interface{}{
				"invoice_id": invoiceID,
				"amount":     amount,
				"reason":     reason,
				"note_type":  noteType,
			}

			logger.Info().Msg("Simulando creación de nota de crédito")

			creditNote, err := softpymesIntegration.SimulateCreditNote(token, creditNoteData)
			if err != nil {
				logger.Error().Err(err).Msg("Error al crear nota de crédito")
				fmt.Printf("❌ Error: %v\n", err)
			} else {
				logger.Info().
					Str("note_number", creditNote.CreditNoteNumber).
					Str("cufe", creditNote.CUFE).
					Msg("Nota de crédito creada exitosamente")
				fmt.Printf("✅ Nota de crédito creada:\n")
				fmt.Printf("  Número: %s\n", creditNote.CreditNoteNumber)
				fmt.Printf("  CUFE: %s\n", creditNote.CUFE)
				fmt.Printf("  Monto: $%.2f\n", creditNote.Amount)
				fmt.Printf("  Tipo: %s\n", creditNote.NoteType)
				fmt.Printf("  PDF: %s\n", creditNote.PDFURL)
			}
		case "14":
			repo := softpymesIntegration.GetRepository()
			invoices := repo.GetAllInvoices()
			creditNotes := repo.GetAllCreditNotes()

			if len(invoices) == 0 && len(creditNotes) == 0 {
				fmt.Println("📭 No hay documentos almacenados")
			} else {
				if len(invoices) > 0 {
					fmt.Printf("\n📄 Facturas almacenadas (%d):\n", len(invoices))
					for i, invoice := range invoices {
						fmt.Printf("  %d. %s - %s - $%.2f %s - Cliente: %s\n",
							i+1, invoice.InvoiceNumber, invoice.OrderID, invoice.Total, invoice.Currency, invoice.CustomerName)
					}
				}
				if len(creditNotes) > 0 {
					fmt.Printf("\n💳 Notas de crédito almacenadas (%d):\n", len(creditNotes))
					for i, note := range creditNotes {
						fmt.Printf("  %d. %s - Factura: %s - $%.2f - Tipo: %s\n",
							i+1, note.CreditNoteNumber, note.InvoiceID, note.Amount, note.NoteType)
					}
				}
			}
		case "0":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("❌ Opción inválida")
		}
	}
}













