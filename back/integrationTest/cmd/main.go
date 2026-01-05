package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/secamc93/probability/back/integrationTest/integrations/shopify"
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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== Simulador de Webhooks Shopify ===")
		fmt.Println("Selecciona una opción:")
		fmt.Println("1. orders/create (crear nueva orden aleatoria)")
		fmt.Println("2. orders/paid (marcar orden como pagada)")
		fmt.Println("3. orders/updated (actualizar orden existente)")
		fmt.Println("4. orders/cancelled (cancelar orden existente)")
		fmt.Println("5. orders/fulfilled (marcar orden como cumplida)")
		fmt.Println("6. orders/partially_fulfilled (marcar orden como parcialmente cumplida)")
		fmt.Println("7. Listar órdenes almacenadas")
		fmt.Println("0. Salir")
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
		case "0":
			fmt.Println("Saliendo...")
			os.Exit(0)
		default:
			fmt.Println("❌ Opción inválida")
		}
	}
}













