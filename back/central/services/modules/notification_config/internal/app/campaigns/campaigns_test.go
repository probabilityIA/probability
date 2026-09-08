package campaigns

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func templateID(id uint) *uint {
	return &id
}

func baseDTO() dtos.CreateCampaignDTO {
	return dtos.CreateCampaignDTO{
		BusinessID:         26,
		WhatsappTemplateID: templateID(7),
		Name:               "Ruta 30 septiembre",
		SenderName:         "Isabel Rojas",
		AudienceType:       entities.CampaignAudienceFiltered,
	}
}

func TestCreateAppliesDefaults(t *testing.T) {
	h := newHarness()

	campaign, err := h.uc.Create(context.Background(), baseDTO())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if campaign.Status != entities.CampaignStatusDraft {
		t.Fatalf("status esperado draft, llego %s", campaign.Status)
	}
	if campaign.Timezone != defaultTimezone {
		t.Fatalf("timezone esperado %s, llego %s", defaultTimezone, campaign.Timezone)
	}
	if campaign.SendWindowStart != defaultSendWindowStart || campaign.SendWindowEnd != defaultSendWindowEnd {
		t.Fatalf("ventana por defecto incorrecta: %s-%s", campaign.SendWindowStart, campaign.SendWindowEnd)
	}
	if campaign.DailySendCap != defaultDailySendCap {
		t.Fatalf("tope diario esperado %d, llego %d", defaultDailySendCap, campaign.DailySendCap)
	}
	if campaign.BatchSize != defaultBatchSize {
		t.Fatalf("tanda esperada %d, llego %d", defaultBatchSize, campaign.BatchSize)
	}
}

func TestCreateRejectsEmptyName(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.Name = "   "

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por nombre vacio")
	}
}

func TestCreateRejectsMissingBusiness(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.BusinessID = 0

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por business_id")
	}
}

func TestCreateRejectsUnsupportedAudience(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.AudienceType = "clientes_de_la_competencia"

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por audiencia no soportada")
	}
}

func TestCreateDefaultsAudienceWhenEmpty(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.AudienceType = ""

	campaign, err := h.uc.Create(context.Background(), dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.AudienceType != entities.CampaignAudienceFiltered {
		t.Fatalf("audiencia esperada filtered, llego %s", campaign.AudienceType)
	}
}

func TestCreateRejectsMissingTemplate(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.WhatsappTemplateID = nil

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por plantilla faltante")
	}
}

func TestCreateRejectsTemplateFromAnotherBusiness(t *testing.T) {
	h := newHarness()
	other := uint(99)
	h.templates.template.BusinessID = &other

	if _, err := h.uc.Create(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba error por plantilla ajena")
	}
}

func TestCreateRejectsUnknownTemplate(t *testing.T) {
	h := newHarness()
	h.templates.template = nil

	if _, err := h.uc.Create(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba error porque la plantilla no existe")
	}
}

func TestCreatePropagatesTemplateRepoError(t *testing.T) {
	h := newHarness()
	h.templates.fail = true

	if _, err := h.uc.Create(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba el error del repositorio de plantillas")
	}
}

func TestCreateRejectsInvalidTimezone(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.Timezone = "Marte/Olimpo"

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por zona horaria invalida")
	}
}

func TestCreateRejectsInvertedWindow(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.SendWindowStart = "20:00"
	dto.SendWindowEnd = "08:00"

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por ventana invertida")
	}
}

func TestCreateRejectsMalformedWindow(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.SendWindowStart = "manana"

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por hora invalida")
	}

	dto = baseDTO()
	dto.SendWindowEnd = "tarde"

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por hora de fin invalida")
	}
}

func TestCreateRejectsDailyCapOverLimit(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.DailySendCap = entities.CampaignMaxDailySendCap + 1

	if _, err := h.uc.Create(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por tope diario excesivo")
	}
}

func TestCreateClampsBatchSize(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.BatchSize = entities.CampaignMaxBatchSize + 500
	dto.DailySendCap = 900

	campaign, err := h.uc.Create(context.Background(), dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.BatchSize != entities.CampaignMaxBatchSize {
		t.Fatalf("la tanda debio recortarse a %d, llego %d", entities.CampaignMaxBatchSize, campaign.BatchSize)
	}
}

func TestCreateBatchNeverExceedsDailyCap(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.DailySendCap = 10
	dto.BatchSize = 100

	campaign, err := h.uc.Create(context.Background(), dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.BatchSize != 10 {
		t.Fatalf("la tanda no puede superar el tope diario, llego %d", campaign.BatchSize)
	}
}

func TestCreatePropagatesRepoError(t *testing.T) {
	h := newHarness()
	h.campaigns.failCreate = true

	if _, err := h.uc.Create(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestUpdateRejectsRunningCampaign(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26, Status: entities.CampaignStatusRunning}

	_, err := h.uc.Update(context.Background(), dtos.UpdateCampaignDTO{ID: 3, CreateCampaignDTO: baseDTO()})
	if err == nil {
		t.Fatal("se esperaba error porque la campana esta en curso")
	}
}

func TestUpdateKeepsStatusAndCreator(t *testing.T) {
	h := newHarness()
	creator := uint(5)
	h.campaigns.byID[3] = &entities.Campaign{
		ID:          3,
		BusinessID:  26,
		Status:      entities.CampaignStatusScheduled,
		CreatedByID: &creator,
	}

	campaign, err := h.uc.Update(context.Background(), dtos.UpdateCampaignDTO{ID: 3, CreateCampaignDTO: baseDTO()})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.Status != entities.CampaignStatusScheduled {
		t.Fatalf("el status no debio cambiar, llego %s", campaign.Status)
	}
	if campaign.CreatedByID == nil || *campaign.CreatedByID != creator {
		t.Fatal("el creador original se perdio")
	}
}

func TestUpdateFailsWhenCampaignMissing(t *testing.T) {
	h := newHarness()

	if _, err := h.uc.Update(context.Background(), dtos.UpdateCampaignDTO{ID: 404, CreateCampaignDTO: baseDTO()}); err == nil {
		t.Fatal("se esperaba error porque la campana no existe")
	}
}

func TestUpdatePropagatesValidationError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26, Status: entities.CampaignStatusDraft}

	dto := baseDTO()
	dto.Name = ""

	if _, err := h.uc.Update(context.Background(), dtos.UpdateCampaignDTO{ID: 3, CreateCampaignDTO: dto}); err == nil {
		t.Fatal("se esperaba error de validacion")
	}
}

func TestUpdatePropagatesRepoError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26, Status: entities.CampaignStatusDraft}
	h.campaigns.failUpdate = true

	if _, err := h.uc.Update(context.Background(), dtos.UpdateCampaignDTO{ID: 3, CreateCampaignDTO: baseDTO()}); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestGetByIDRejectsForeignBusiness(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 99}

	if _, err := h.uc.GetByID(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error de aislamiento por negocio")
	}
}

func TestGetByIDReturnsNilWhenMissing(t *testing.T) {
	h := newHarness()

	campaign, err := h.uc.GetByID(context.Background(), 404, 26)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign != nil {
		t.Fatal("se esperaba nil")
	}
}

func TestGetByIDPropagatesRepoError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	if _, err := h.uc.GetByID(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestListNormalizesPaging(t *testing.T) {
	h := newHarness()
	h.campaigns.list = []entities.Campaign{{ID: 1}}
	h.campaigns.listTotal = 1

	items, total, err := h.uc.List(context.Background(), 26, "", -3, 0)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(items) != 1 || total != 1 {
		t.Fatal("no se devolvieron los datos del repositorio")
	}
}

func TestListPropagatesRepoError(t *testing.T) {
	h := newHarness()
	h.campaigns.failList = true

	if _, _, err := h.uc.List(context.Background(), 26, "", 1, 10); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestDeleteBlocksRunningCampaign(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26, Status: entities.CampaignStatusRunning}

	if err := h.uc.Delete(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: hay que pausar antes de borrar")
	}
}

func TestDeleteRemovesDraft(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26, Status: entities.CampaignStatusDraft}

	if err := h.uc.Delete(context.Background(), 3, 26); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if h.campaigns.deletedID != 3 {
		t.Fatal("no se borro la campana")
	}
}

func TestDeleteFailsWhenMissing(t *testing.T) {
	h := newHarness()

	if err := h.uc.Delete(context.Background(), 404, 26); err == nil {
		t.Fatal("se esperaba error porque la campana no existe")
	}
}

func TestDeletePropagatesGetError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	if err := h.uc.Delete(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestPreviewAudienceReturnsCounts(t *testing.T) {
	h := newHarness()
	h.audience.total = 120
	h.audience.optedOut = 8
	h.audience.noPhone = 12
	h.audience.reachable = 95
	h.audience.candidates = []entities.CampaignCandidate{
		{ClientID: 1, Name: "Ana Perez"},
		{ClientID: 2, Name: "Luis Gomez"},
	}

	preview, err := h.uc.PreviewAudience(context.Background(), baseDTO())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if preview.Total != 120 || preview.Reachable != 95 {
		t.Fatalf("conteos incorrectos: %+v", preview)
	}
	if len(preview.SampleName) != 2 {
		t.Fatalf("se esperaban 2 nombres de muestra, llegaron %d", len(preview.SampleName))
	}
}

func TestPreviewAudienceRequiresBusiness(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.BusinessID = 0

	if _, err := h.uc.PreviewAudience(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por business_id")
	}
}

func TestPreviewAudienceRejectsUnsupportedAudience(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.AudienceType = "inventada"

	if _, err := h.uc.PreviewAudience(context.Background(), dto); err == nil {
		t.Fatal("se esperaba error por audiencia no soportada")
	}
}

func TestPreviewAudiencePropagatesCountError(t *testing.T) {
	h := newHarness()
	h.audience.failCount = true

	if _, err := h.uc.PreviewAudience(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba el error del querier")
	}
}

func TestPreviewAudiencePropagatesFindError(t *testing.T) {
	h := newHarness()
	h.audience.failFind = true

	if _, err := h.uc.PreviewAudience(context.Background(), baseDTO()); err == nil {
		t.Fatal("se esperaba el error del querier")
	}
}

func TestListSendsRequiresCampaign(t *testing.T) {
	h := newHarness()

	if _, _, err := h.uc.ListSends(context.Background(), 404, 26, "", 1, 10); err == nil {
		t.Fatal("se esperaba error porque la campana no existe")
	}
}

func TestListSendsReturnsRows(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = &entities.Campaign{ID: 3, BusinessID: 26}
	h.sends.list = []entities.CampaignSend{{ID: 1}}
	h.sends.listTotal = 1

	items, total, err := h.uc.ListSends(context.Background(), 3, 26, "sent", 1, 200)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(items) != 1 || total != 1 {
		t.Fatal("no se devolvieron los envios")
	}
}

func TestListSendsPropagatesGetError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	if _, _, err := h.uc.ListSends(context.Background(), 3, 26, "", 1, 10); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestMarkSendResultDelegates(t *testing.T) {
	h := newHarness()

	err := h.uc.MarkSendResult(context.Background(), dtos.CampaignSendResult{SendID: 1, Status: "sent"})
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.sends.results) != 1 || h.sends.results[0] != "sent" {
		t.Fatal("no se guardo el resultado")
	}
}

func TestMarkSendResultPropagatesError(t *testing.T) {
	h := newHarness()
	h.sends.failResult = true

	if err := h.uc.MarkSendResult(context.Background(), dtos.CampaignSendResult{SendID: 1}); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestNormalizePagingClampsUpperBound(t *testing.T) {
	page, size := normalizePaging(2, 5000)
	if page != 2 || size != 100 {
		t.Fatalf("paginacion incorrecta: %d %d", page, size)
	}
}

func TestCampaignHelpers(t *testing.T) {
	draft := entities.Campaign{Status: entities.CampaignStatusDraft}
	if !draft.IsEditable() {
		t.Fatal("un borrador debe ser editable")
	}

	done := entities.Campaign{Status: entities.CampaignStatusCompleted}
	if !done.IsFinished() || done.IsEditable() {
		t.Fatal("una campana completada no se edita y esta terminada")
	}

	cancelled := entities.Campaign{Status: entities.CampaignStatusCancelled}
	if !cancelled.IsFinished() {
		t.Fatal("una campana cancelada esta terminada")
	}

	if entities.IsSupportedAudience("otra_cosa") {
		t.Fatal("audiencia no soportada aceptada")
	}
	if !entities.IsSupportedAudience(entities.CampaignAudienceAllClients) {
		t.Fatal("all_clients debe ser soportada")
	}
}

func TestBuildAudienceParamsTrimsCity(t *testing.T) {
	dto := baseDTO()
	dto.City = "  Bogota  "
	dto.ClientIDs = []uint{1, 2}
	dto.CreatedFromDays = 30
	dto.OnlyWithoutOrder = true

	params := buildAudienceParams(dto)
	if params.City != "Bogota" {
		t.Fatalf("la ciudad no se limpio: %q", params.City)
	}
	if len(params.ClientIDs) != 2 || params.CreatedFromDays != 30 || !params.OnlyWithoutOrder {
		t.Fatal("los filtros no se copiaron")
	}
}

func TestScheduledAtIsKept(t *testing.T) {
	h := newHarness()
	when := time.Now().Add(2 * time.Hour)
	dto := baseDTO()
	dto.ScheduledAt = &when

	campaign, err := h.uc.Create(context.Background(), dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.ScheduledAt == nil || !campaign.ScheduledAt.Equal(when) {
		t.Fatal("la fecha programada no se conservo")
	}
}

func updateDTOFor(id uint) dtos.UpdateCampaignDTO {
	return dtos.UpdateCampaignDTO{ID: id, CreateCampaignDTO: baseDTO()}
}
