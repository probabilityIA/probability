package usecasemessaging

import "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"

func (u *usecases) SystemFlows() []entities.SystemFlow {
	return SystemFlowCatalog()
}

func SystemFlowCatalog() []entities.SystemFlow {
	novelty := func(label, template, detail string) entities.SystemFlowBranch {
		return entities.SystemFlowBranch{
			Label: label,
			Node: &entities.SystemFlowNode{
				Kind:     entities.SystemFlowNodeTemplate,
				Template: template,
				Event:    "novelty",
				Effects:  []string{detail, "La orden queda sin confirmar"},
			},
		}
	}

	return []entities.SystemFlow{
		{
			Key:         "order_confirmation",
			Name:        "Confirmación de pedido",
			Description: "Pide al cliente que confirme su pedido y resuelve cancelaciones, novedades y paso a un asesor sin intervención del negocio.",
			Trigger:     "Se envía cuando entra un pedido y la regla “Confirmación de pedido” de WhatsApp está activa.",
			Root: &entities.SystemFlowNode{
				Kind:     entities.SystemFlowNodeTemplate,
				Template: "confirmacion_pedido_contraentrega",
				Detail:   "Se usa la variante que corresponda a la orden: contra entrega, sin valor, sin contra entrega o con mapa.",
				State:    entities.StateAwaitingConfirmation,
				Branches: []entities.SystemFlowBranch{
					{
						Label: "Confirmar pedido",
						Node: &entities.SystemFlowNode{
							Kind:     entities.SystemFlowNodeTemplate,
							Template: "pedido_confirmado_v2",
							Event:    "confirmed",
							Effects:  []string{"La orden queda confirmada"},
						},
					},
					{
						Label: "No confirmar",
						Node: &entities.SystemFlowNode{
							Kind:     entities.SystemFlowNodeTemplate,
							Template: "menu_no_confirmacion",
							State:    entities.StateAwaitingMenuSelection,
							Branches: []entities.SystemFlowBranch{
								{
									Label: "Presentar novedad",
									Node: &entities.SystemFlowNode{
										Kind:     entities.SystemFlowNodeTemplate,
										Template: "tipo_novedad_pedido",
										State:    entities.StateAwaitingNoveltyType,
										Branches: []entities.SystemFlowBranch{
											novelty("Cambio de dirección", "novedad_cambio_direccion", "Se registra la novedad de cambio de dirección"),
											novelty("Cambio de productos", "novedad_cambio_productos", "Se registra la novedad de cambio de productos"),
											novelty("Cambio medio de pago", "novedad_cambio_medio_pago", "Se registra la novedad de cambio de medio de pago"),
										},
									},
								},
								{
									Label: "Cancelar pedido",
									Node: &entities.SystemFlowNode{
										Kind:     entities.SystemFlowNodeTemplate,
										Template: "confirmar_cancelacion_pedido",
										State:    entities.StateAwaitingCancelConfirm,
										Branches: []entities.SystemFlowBranch{
											{
												Label: "Sí, cancelar",
												Node: &entities.SystemFlowNode{
													Kind:   entities.SystemFlowNodeDecision,
													Title:  "¿La orden ya tiene guía?",
													Detail: "Se revisa el último envío: guía generada y no anulada.",
													Event:  "cancelled",
													Branches: []entities.SystemFlowBranch{
														{
															Label: "No tiene guía",
															Node: &entities.SystemFlowNode{
																Kind:     entities.SystemFlowNodeTemplate,
																Template: "pedido_cancelado",
																Effects:  []string{"La orden pasa a Cancelada", "Vía avisa: Orden cancelada"},
															},
														},
														{
															Label: "Ya tiene guía",
															Node: &entities.SystemFlowNode{
																Kind:     entities.SystemFlowNodeTemplate,
																Template: "solicitud_cancelacion_recibida",
																Effects:  []string{"La orden pasa a Cliente solicita cancelar", "Vía avisa: El cliente solicita cancelar"},
															},
														},
													},
												},
											},
											{
												Label:  "No, volver",
												BackTo: "menu_no_confirmacion",
											},
										},
									},
								},
								{
									Label: "Asesor",
									Node: &entities.SystemFlowNode{
										Kind:     entities.SystemFlowNodeTemplate,
										Template: "handoff_asesor",
										Event:    "handoff",
										Effects:  []string{"La conversación pasa a una persona del negocio"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
