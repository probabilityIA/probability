package usecasemessaging

import (
	"context"
	"sort"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/mocks"
)

func TestSystemFlowCatalogCoincideConLaMaquinaDeEstados(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationCacheMock{},
		&mocks.PersistencePublisherMock{},
		&mocks.CredentialsCacheMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	var walk func(node *entities.SystemFlowNode, path string)
	walk = func(node *entities.SystemFlowNode, path string) {
		if node == nil {
			return
		}
		if node.Kind == entities.SystemFlowNodeTemplate {
			if _, ok := entities.Templates[node.Template]; !ok {
				t.Errorf("%s: la plantilla %q no existe en el catalogo de plantillas", path, node.Template)
			}
		}
		if node.State != "" {
			labels := make([]string, 0, len(node.Branches))
			for _, branch := range node.Branches {
				labels = append(labels, branch.Label)
			}
			available := append([]string{}, uc.GetAvailableResponses(node.State)...)
			sort.Strings(labels)
			sort.Strings(available)
			if len(labels) != len(available) {
				t.Errorf("%s: el flujo muestra %v pero el estado %s acepta %v", path, labels, node.State, available)
			}
			for i := range labels {
				if i < len(available) && labels[i] != available[i] {
					t.Errorf("%s: el flujo muestra %v pero el estado %s acepta %v", path, labels, node.State, available)
					break
				}
			}

			for _, branch := range node.Branches {
				step := path + " > " + branch.Label
				transition, err := uc.TransitionState(context.Background(), newActiveConversation(node.State), branch.Label)
				if err != nil {
					t.Errorf("%s: la maquina de estados no acepta el boton: %v", step, err)
					continue
				}
				switch {
				case branch.BackTo != "":
					if transition.TemplateName != branch.BackTo {
						t.Errorf("%s: el flujo dice que vuelve a %q y la maquina envia %q", step, branch.BackTo, transition.TemplateName)
					}
				case branch.Node == nil:
					t.Errorf("%s: rama sin destino", step)
				case branch.Node.Kind == entities.SystemFlowNodeDecision:
					if transition.TemplateName != "" || transition.EventType != branch.Node.Event {
						t.Errorf("%s: se esperaba decision con evento %q, la maquina envia %q con evento %q", step, branch.Node.Event, transition.TemplateName, transition.EventType)
					}
				default:
					if transition.TemplateName != branch.Node.Template {
						t.Errorf("%s: el flujo muestra %q y la maquina envia %q", step, branch.Node.Template, transition.TemplateName)
					}
					if transition.EventType != branch.Node.Event {
						t.Errorf("%s: el flujo dice evento %q y la maquina publica %q", step, branch.Node.Event, transition.EventType)
					}
				}
			}
		}
		for _, branch := range node.Branches {
			walk(branch.Node, path+" > "+branch.Label)
		}
	}

	for _, flow := range SystemFlowCatalog() {
		walk(flow.Root, flow.Key)
	}
}
