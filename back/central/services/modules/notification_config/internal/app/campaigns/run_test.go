package campaigns

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func launchableCampaign() *entities.Campaign {
	return &entities.Campaign{
		ID:                 3,
		BusinessID:         26,
		WhatsappTemplateID: templateID(7),
		Name:               "Ruta 30",
		SenderName:         "Isabel Rojas",
		AudienceType:       entities.CampaignAudienceFiltered,
		Status:             entities.CampaignStatusDraft,
		Timezone:           defaultTimezone,
		SendWindowStart:    "00:00",
		SendWindowEnd:      "23:59",
		DailySendCap:       100,
		BatchSize:          10,
	}
}

func TestLaunchRequiresOwnNumber(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.sender.integrationID = 0

	_, err := h.uc.Launch(context.Background(), 3, 26)
	if err == nil {
		t.Fatal("se esperaba error: marketing solo sale desde el numero propio")
	}
}

func TestLaunchPropagatesSenderError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.sender.fail = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del querier de numero propio")
	}
}

func TestLaunchRequiresApprovedTemplate(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.templates.template.Status = entities.TemplateStatusPending

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la plantilla no esta aprobada")
	}
}

func TestLaunchRequiresTemplateID(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.WhatsappTemplateID = nil
	h.campaigns.byID[3] = campaign

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la campana no tiene plantilla")
	}
}

func TestLaunchRejectsMissingTemplate(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.templates.template = nil

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la plantilla ya no existe")
	}
}

func TestLaunchRejectsForeignTemplate(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	other := uint(99)
	h.templates.template.BusinessID = &other

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la plantilla es de otro negocio")
	}
}

func TestLaunchRejectsEmptyAudience(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.candidates = nil

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la audiencia quedo vacia")
	}
}

func TestLaunchQueuesAudienceAndSchedules(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.candidates = []entities.CampaignCandidate{
		{ClientID: 1, Phone: "573001112233", Name: "Ana Perez"},
		{ClientID: 2, Phone: "573004445566", Name: "Luis Gomez"},
	}

	campaign, err := h.uc.Launch(context.Background(), 3, 26)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if campaign.Status != entities.CampaignStatusScheduled {
		t.Fatalf("status esperado scheduled, llego %s", campaign.Status)
	}
	if campaign.AudienceCount != 2 {
		t.Fatalf("audiencia esperada 2, llego %d", campaign.AudienceCount)
	}
	if campaign.StartedAt == nil || campaign.ScheduledAt == nil {
		t.Fatal("no se marcaron las fechas de arranque")
	}
	if len(h.sends.bulkCreated) != 2 {
		t.Fatalf("se esperaban 2 envios creados, llegaron %d", len(h.sends.bulkCreated))
	}
}

func TestLaunchRejectsAlreadyLaunched(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusRunning
	h.campaigns.byID[3] = campaign

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la campana ya fue lanzada")
	}
}

func TestLaunchFailsWhenCampaignMissing(t *testing.T) {
	h := newHarness()

	if _, err := h.uc.Launch(context.Background(), 404, 26); err == nil {
		t.Fatal("se esperaba error porque la campana no existe")
	}
}

func TestLaunchPropagatesAudienceError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.failFind = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del querier de audiencia")
	}
}

func TestLaunchPropagatesBulkError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.candidates = []entities.CampaignCandidate{{ClientID: 1, Phone: "573001112233"}}
	h.sends.failBulk = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error al crear los envios")
	}
}

func TestLaunchPropagatesUpdateError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.candidates = []entities.CampaignCandidate{{ClientID: 1, Phone: "573001112233"}}
	h.campaigns.failUpdate = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestPauseAndResume(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusRunning
	h.campaigns.byID[3] = campaign

	paused, err := h.uc.Pause(context.Background(), 3, 26)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if paused.Status != entities.CampaignStatusPaused {
		t.Fatalf("status esperado paused, llego %s", paused.Status)
	}

	resumed, err := h.uc.Resume(context.Background(), 3, 26)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if resumed.Status != entities.CampaignStatusScheduled {
		t.Fatalf("status esperado scheduled, llego %s", resumed.Status)
	}
}

func TestPauseRejectsDraft(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()

	if _, err := h.uc.Pause(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: solo se pausa lo programado o en curso")
	}
}

func TestResumeRejectsNonPaused(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()

	if _, err := h.uc.Resume(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: solo se reanuda lo pausado")
	}
}

func TestCancelSetsFinishedAt(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusRunning
	h.campaigns.byID[3] = campaign

	cancelled, err := h.uc.Cancel(context.Background(), 3, 26)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if cancelled.Status != entities.CampaignStatusCancelled || cancelled.FinishedAt == nil {
		t.Fatal("la cancelacion no quedo registrada")
	}
}

func TestCancelRejectsFinished(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusCompleted
	h.campaigns.byID[3] = campaign

	if _, err := h.uc.Cancel(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba error: la campana ya termino")
	}
}

func TestTransitionFailsWhenMissing(t *testing.T) {
	h := newHarness()

	if _, err := h.uc.Cancel(context.Background(), 404, 26); err == nil {
		t.Fatal("se esperaba error porque la campana no existe")
	}
}

func TestTransitionPropagatesGetError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	if _, err := h.uc.Cancel(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestTransitionPropagatesUpdateError(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusRunning
	h.campaigns.byID[3] = campaign
	h.campaigns.failUpdate = true

	if _, err := h.uc.Pause(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func openWindowNow() (string, string) {
	location, err := time.LoadLocation(defaultTimezone)
	if err != nil {
		location = time.UTC
	}
	now := time.Now().In(location)
	start := now.Add(-30 * time.Minute)
	end := now.Add(30 * time.Minute)

	if start.Day() != now.Day() || end.Day() != now.Day() {
		return "", ""
	}

	return start.Format("15:04"), end.Format("15:04")
}

func runnableCampaign() entities.Campaign {
	campaign := launchableCampaign()
	campaign.Status = entities.CampaignStatusScheduled
	campaign.SendWindowStart, campaign.SendWindowEnd = openWindowNow()
	return *campaign
}

func TestRunDueCampaignsPublishesBatch(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pending = []entities.CampaignSend{
		{ID: 11, ClientID: 1, Phone: "573001112233", ClientName: "Ana Perez"},
		{ID: 12, ClientID: 2, Phone: "573004445566", ClientName: "Luis Gomez"},
	}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.publisher.published) != 2 {
		t.Fatalf("se esperaban 2 mensajes publicados, llegaron %d", len(h.publisher.published))
	}

	first := h.publisher.published[0]
	if first.TemplateName != "ruta_30_dian" || first.Language != "es" {
		t.Fatalf("la plantilla no se resolvio: %+v", first)
	}
	if len(first.Parameters) != 2 || first.Parameters[0] != "Ana" || first.Parameters[1] != "Isabel Rojas" {
		t.Fatalf("los parametros no se armaron bien: %+v", first.Parameters)
	}
	if len(h.sends.queued) != 2 {
		t.Fatal("los envios no quedaron marcados como encolados")
	}
}

func TestRunDueCampaignsRespectsDailyCap(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.DailySendCap = 5
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.sentSince = 5
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "573001112233"}}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.publisher.published) != 0 {
		t.Fatal("no debio publicarse nada: el tope diario ya se alcanzo")
	}
}

func TestRunDueCampaignsLimitsBatchToRemaining(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.DailySendCap = 3
	campaign.BatchSize = 10
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.sentSince = 2
	h.sends.pending = []entities.CampaignSend{
		{ID: 11, Phone: "1"}, {ID: 12, Phone: "2"}, {ID: 13, Phone: "3"},
	}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.publisher.published) != 1 {
		t.Fatalf("solo quedaba 1 mensaje de cupo, se publicaron %d", len(h.publisher.published))
	}
}

func TestRunDueCampaignsCompletesWhenNothingPending(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pending = nil

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.campaigns.updated) == 0 {
		t.Fatal("la campana no se cerro")
	}
	last := h.campaigns.updated[len(h.campaigns.updated)-1]
	if last.Status != entities.CampaignStatusCompleted || last.FinishedAt == nil {
		t.Fatalf("la campana debio quedar completada, llego %s", last.Status)
	}
}

func TestRunDueCampaignsSkipsOutsideWindow(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.SendWindowStart = "03:00"
	campaign.SendWindowEnd = "03:01"
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "1"}}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	now := time.Now().In(mustLoad(defaultTimezone))
	insideWindow := now.Hour() == 3 && now.Minute() == 0
	if !insideWindow && len(h.publisher.published) != 0 {
		t.Fatal("no debio publicarse nada fuera de la ventana")
	}
}

func TestRunDueCampaignsPausesOnBadTemplate(t *testing.T) {
	h := newHarness()
	h.campaigns.due = []entities.Campaign{runnableCampaign()}
	h.templates.template.Status = entities.TemplateStatusRejected

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.campaigns.updated) == 0 {
		t.Fatal("la campana no se actualizo")
	}
	last := h.campaigns.updated[len(h.campaigns.updated)-1]
	if last.Status != entities.CampaignStatusPaused || last.ErrorMessage == "" {
		t.Fatalf("la campana debio pausarse con motivo, llego %s", last.Status)
	}
}

func TestRunDueCampaignsPropagatesListError(t *testing.T) {
	h := newHarness()
	h.campaigns.failDue = true

	if err := h.uc.RunDueCampaigns(context.Background()); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestDispatchBatchPropagatesCountError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.failCount = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error del conteo")
	}
}

func TestDispatchBatchPropagatesPendingError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.failPending = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error al listar pendientes")
	}
}

func TestDispatchBatchMarksFailedWhenPublishFails(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "573001112233"}}
	h.publisher.fail = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.sends.results) != 1 || h.sends.results[0] != entities.CampaignSendStatusFailed {
		t.Fatal("el envio debio quedar como fallido")
	}
}

func TestDispatchBatchSkipsWhenQueueMarkFails(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "573001112233"}}
	h.sends.queuedErrFor = 11

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.publisher.published) != 0 {
		t.Fatal("no debio publicarse si no se pudo reservar el envio")
	}
}

func TestDispatchBatchWithoutPublisher(t *testing.T) {
	h := newHarness()
	h.impl.publisher = nil
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "573001112233"}}

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if len(h.sends.queued) != 0 {
		t.Fatal("sin publicador no se debe reservar nada")
	}
}

func TestDispatchBatchPropagatesBatchMarkError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "1"}}
	h.campaigns.failBatch = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error al marcar la tanda")
	}
}

func TestDispatchBatchPropagatesCountersError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "1"}}
	h.campaigns.failCounters = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error al recalcular contadores")
	}
}

func TestCompleteCampaignPropagatesErrors(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.campaigns.failUpdate = true

	if err := h.impl.completeCampaign(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestFailCampaignPropagatesUpdateError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.campaigns.failUpdate = true

	if err := h.impl.failCampaign(context.Background(), &campaign, errBoom); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestMaterializeAudienceReturnsZeroWhenEmpty(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()

	count, err := h.impl.materializeAudience(context.Background(), &campaign)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if count != 0 {
		t.Fatalf("se esperaba 0, llego %d", count)
	}
}

func TestVariableValuesOverrideResolution(t *testing.T) {
	campaign := runnableCampaign()
	campaign.VariableValues = map[string]string{"customer.first_name": "Doctora"}

	value := resolveVariable(
		entities.TemplateVariable{Position: 1, Source: "customer.first_name"},
		&campaign,
		entities.CampaignSend{ClientName: "Ana Perez"},
	)
	if value != "Doctora" {
		t.Fatalf("el valor fijo de la campana debe ganar, llego %q", value)
	}
}

func TestResolveVariableSources(t *testing.T) {
	campaign := runnableCampaign()
	send := entities.CampaignSend{ClientName: "Ana Maria Perez"}

	cases := map[string]string{
		"customer.first_name": "Ana",
		"customer.full_name":  "Ana Maria Perez",
		senderNameVariable:    "Isabel Rojas",
		campaignNameVariable:  "Ruta 30",
		"algo.inventado":      "",
	}

	for source, expected := range cases {
		got := resolveVariable(entities.TemplateVariable{Source: source}, &campaign, send)
		if got != expected {
			t.Fatalf("para %s se esperaba %q, llego %q", source, expected, got)
		}
	}
}

func TestBuildParametersUsesFallbackAndDash(t *testing.T) {
	campaign := runnableCampaign()
	template := &entities.WhatsappTemplate{
		Variables: []entities.TemplateVariable{
			{Position: 1, Source: "algo.inventado", Fallback: "cliente"},
			{Position: 2, Source: "otra.cosa"},
		},
	}

	params := buildParameters(template, &campaign, entities.CampaignSend{})
	if len(params) != 2 || params[0] != "cliente" || params[1] != "-" {
		t.Fatalf("los parametros de respaldo no se aplicaron: %+v", params)
	}
}

func TestFirstNameEmpty(t *testing.T) {
	if firstName("   ") != "" {
		t.Fatal("un nombre vacio debe devolver cadena vacia")
	}
}

func TestIsInsideWindowFallsBackWhenClockInvalid(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.SendWindowStart = "no-es-hora"

	if !h.impl.isInsideWindow(&campaign, time.Now()) {
		t.Fatal("con una hora invalida no se debe bloquear el envio")
	}

	campaign = runnableCampaign()
	campaign.SendWindowEnd = "tampoco"
	if !h.impl.isInsideWindow(&campaign, time.Now()) {
		t.Fatal("con una hora de fin invalida no se debe bloquear el envio")
	}
}

func TestIsInsideWindowFallsBackToUTC(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.Timezone = "Marte/Olimpo"
	campaign.SendWindowStart = "00:00"
	campaign.SendWindowEnd = "23:59"

	if !h.impl.isInsideWindow(&campaign, time.Now()) {
		t.Fatal("con zona horaria invalida se cae a UTC y la ventana completa aplica")
	}
}

func TestStartOfDayFallsBackToUTC(t *testing.T) {
	campaign := runnableCampaign()
	campaign.Timezone = "Marte/Olimpo"

	start := startOfDay(&campaign, time.Now())
	if start.Hour() != 0 || start.Minute() != 0 {
		t.Fatalf("el inicio del dia debe ser medianoche, llego %s", start)
	}
}

func mustLoad(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return location
}

func TestUpdatePropagatesGetError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	_, err := h.uc.Update(context.Background(), updateDTOFor(3))
	if err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestLaunchPropagatesTemplateRepoError(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.templates.fail = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio de plantillas")
	}
}

func TestLaunchPropagatesGetError(t *testing.T) {
	h := newHarness()
	h.campaigns.failGet = true

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("se esperaba el error del repositorio")
	}
}

func TestRunDueCampaignsKeepsGoingAfterBatchError(t *testing.T) {
	h := newHarness()
	first := runnableCampaign()
	first.ID = 3
	second := runnableCampaign()
	second.ID = 4
	h.campaigns.due = []entities.Campaign{first, second}
	h.sends.failCount = true

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatalf("un fallo de una campana no debe cortar el resto: %v", err)
	}
}

func TestDispatchBatchPropagatesUpdateError(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	h.sends.pending = []entities.CampaignSend{{ID: 11, Phone: "1"}}
	h.campaigns.failUpdate = true

	if err := h.impl.dispatchBatch(context.Background(), &campaign); err == nil {
		t.Fatal("se esperaba el error al guardar la campana")
	}
}
