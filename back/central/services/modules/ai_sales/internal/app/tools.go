package app

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
)

// GetToolDefinitions retorna las definiciones de herramientas para Bedrock
func GetToolDefinitions() []domain.ToolDefinition {
	return []domain.ToolDefinition{
		{
			Name:        "SearchProducts",
			Description: "Busca productos en el catalogo de la tienda por nombre, descripcion o categoria. Usa esta herramienta SIEMPRE antes de recomendar un producto al cliente.",
			InputSchema: `{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Termino de busqueda para encontrar productos (nombre, categoria, descripcion)"
					},
					"limit": {
						"type": "integer",
						"description": "Numero maximo de resultados (default 5, max 20)"
					}
				},
				"required": ["query"]
			}`,
		},
		{
			Name:        "CreateOrder",
			Description: "Crea un pedido para el cliente. Solo usar despues de que el cliente confirme los productos y cantidades que desea comprar.",
			InputSchema: `{
				"type": "object",
				"properties": {
					"customer_name": {
						"type": "string",
						"description": "Nombre del cliente que hace el pedido"
					},
					"customer_phone": {
						"type": "string",
						"description": "Numero de telefono del cliente"
					},
					"items": {
						"type": "array",
						"description": "Lista de productos a incluir en el pedido",
						"items": {
							"type": "object",
							"properties": {
								"product_sku": {
									"type": "string",
									"description": "SKU del producto"
								},
								"quantity": {
									"type": "integer",
									"description": "Cantidad deseada"
								}
							},
							"required": ["product_sku", "quantity"]
						}
					}
				},
				"required": ["customer_name", "customer_phone", "items"]
			}`,
		},
		{
			Name:        "SearchCustomer",
			Description: "Busca un cliente existente por DNI, email, telefono o nombre. Usa esta herramienta para verificar si el cliente ya existe antes de crear un pedido.",
			InputSchema: `{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Termino de busqueda: DNI, email, telefono o nombre del cliente"
					}
				},
				"required": ["query"]
			}`,
		},
		{
			Name:        "GetCustomerLastAddress",
			Description: "Obtiene la ultima direccion de envio de un cliente a partir de sus pedidos anteriores. Usa esta herramienta despues de identificar al cliente con SearchCustomer para sugerir su direccion.",
			InputSchema: `{
				"type": "object",
				"properties": {
					"customer_id": {
						"type": "integer",
						"description": "ID del cliente obtenido de SearchCustomer"
					}
				},
				"required": ["customer_id"]
			}`,
		},
	}
}

// toolDeps agrupa las dependencias necesarias para ejecutar tools
type toolDeps struct {
	productRepo    domain.IProductRepository
	customerRepo   domain.ICustomerRepository
	orderPublisher domain.IAIOrderPublisher
	businessID     uint
}

// DispatchTool ejecuta la herramienta solicitada por Bedrock y retorna el resultado como string JSON
func DispatchTool(ctx context.Context, toolName string, inputJSON string, deps *toolDeps) (string, error) {
	switch toolName {
	case "SearchProducts":
		return executeSearchProducts(ctx, inputJSON, deps)
	case "CreateOrder":
		return executeCreateOrder(ctx, inputJSON, deps)
	case "SearchCustomer":
		return executeSearchCustomer(ctx, inputJSON, deps)
	case "GetCustomerLastAddress":
		return executeGetCustomerLastAddress(ctx, inputJSON, deps)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// parseToolInput deserializa el input JSON de una herramienta
func parseToolInput(inputJSON string, target any) error {
	if err := json.Unmarshal([]byte(inputJSON), target); err != nil {
		return fmt.Errorf("error parsing tool input: %w", err)
	}
	return nil
}
