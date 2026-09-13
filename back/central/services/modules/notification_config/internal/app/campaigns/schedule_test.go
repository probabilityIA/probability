package campaigns

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func todayIn(campaign *entities.Campaign, offsetDays int) string {
	return time.Now().In(campaign.Location()).AddDate(0, 0, offsetDays).Format("2006-01-02")
}

func TestCreateDefaultsToDailyDistribute(t *testing.T) {
	h := newHarness()

	campaign, err := h.uc.Create(context.Background(), baseDTO())
	if err != nil {
		t.Fatal(err)
	}
	if campaign.ScheduleMode != entities.CampaignScheduleDaily || campaign.DeliveryMode != entities.CampaignDeliveryDistribute {
		t.Fatalf("unexpected defaults: %s %s", campaign.ScheduleMode, campaign.DeliveryMode)
	}
}

func TestCreateValidatesSchedule(t *testing.T) {
	cases := map[string]func(*dtos.CreateCampaignDTO){
		"modo desconocido":      func(d *dtos.CreateCampaignDTO) { d.ScheduleMode = "weekly" },
		"entrega desconocida":   func(d *dtos.CreateCampaignDTO) { d.DeliveryMode = "blast" },
		"intervalo en cero":     func(d *dtos.CreateCampaignDTO) { d.ScheduleMode = entities.CampaignScheduleInterval },
		"intervalo gigante":     func(d *dtos.CreateCampaignDTO) { d.ScheduleMode = entities.CampaignScheduleInterval; d.IntervalDays = 400 },
		"fechas vacias":         func(d *dtos.CreateCampaignDTO) { d.ScheduleMode = entities.CampaignScheduleDates },
		"fecha mal escrita":     func(d *dtos.CreateCampaignDTO) { d.ScheduleMode = entities.CampaignScheduleDates; d.SendDates = []string{"15/09/2026"} },
		"repetir sin fin":       func(d *dtos.CreateCampaignDTO) { d.DeliveryMode = entities.CampaignDeliveryRepeat },
		"demasiadas repeticion": func(d *dtos.CreateCampaignDTO) { d.Occurrences = 1000 },
	}

	for name, mutate := range cases {
		h := newHarness()
		dto := baseDTO()
		mutate(&dto)
		if _, err := h.uc.Create(context.Background(), dto); err == nil {
			t.Fatalf("%s: expected error", name)
		}
	}
}

func TestCreateNormalizesDates(t *testing.T) {
	h := newHarness()
	dto := baseDTO()
	dto.ScheduleMode = entities.CampaignScheduleDates
	dto.DeliveryMode = entities.CampaignDeliveryRepeat
	dto.SendDates = []string{"2026-10-02", "2026-10-01", "2026-10-02"}
	dto.Occurrences = 9

	campaign, err := h.uc.Create(context.Background(), dto)
	if err != nil {
		t.Fatal(err)
	}
	if len(campaign.SendDates) != 2 || campaign.SendDates[0] != "2026-10-01" {
		t.Fatalf("dates not normalized: %v", campaign.SendDates)
	}
	if campaign.Occurrences != 0 {
		t.Fatalf("dates mode must ignore occurrences, got %d", campaign.Occurrences)
	}
}

func TestLaunchRejectsPastDates(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(campaign, -2)}
	h.campaigns.byID[3] = campaign
	h.audience.candidates = []entities.CampaignCandidate{{ClientID: 1, Phone: "573001"}}

	if _, err := h.uc.Launch(context.Background(), 3, 26); err == nil {
		t.Fatal("expected error for past dates")
	}
}

func TestLaunchStartsRoundOne(t *testing.T) {
	h := newHarness()
	h.campaigns.byID[3] = launchableCampaign()
	h.audience.candidates = []entities.CampaignCandidate{{ClientID: 1, Phone: "573001"}}

	campaign, err := h.uc.Launch(context.Background(), 3, 26)
	if err != nil {
		t.Fatal(err)
	}
	if campaign.CurrentRound != 1 || h.sends.bulkCreated[0].Round != 1 {
		t.Fatalf("launch must materialize round 1, got %d / %d", campaign.CurrentRound, h.sends.bulkCreated[0].Round)
	}
}

func TestRunDueCampaignsSkipsDayOutsideSchedule(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(&campaign, 3)}
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pending = []entities.CampaignSend{{ID: 1, Phone: "573001"}}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(h.publisher.published) != 0 {
		t.Fatalf("no messages should go out on a day that is not scheduled, got %d", len(h.publisher.published))
	}
}

func TestRunDueCampaignsSendsOnScheduledDate(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(&campaign, 0), todayIn(&campaign, 5)}
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pending = []entities.CampaignSend{{ID: 1, Phone: "573001"}}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(h.publisher.published) != 1 {
		t.Fatalf("expected one message, got %d", len(h.publisher.published))
	}
}

func TestRunDueCampaignsPausesWhenDatesEndWithPending(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(&campaign, -1)}
	h.campaigns.due = []entities.Campaign{campaign}
	h.sends.pendingCount = 12

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	last := h.campaigns.updated[len(h.campaigns.updated)-1]
	if last.Status != entities.CampaignStatusPaused || last.ErrorMessage == "" {
		t.Fatalf("expected paused with reason, got %s %q", last.Status, last.ErrorMessage)
	}
}

func TestRunDueCampaignsCompletesWhenDatesEndWithoutPending(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(&campaign, -1)}
	h.campaigns.due = []entities.Campaign{campaign}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	if last := h.campaigns.updated[len(h.campaigns.updated)-1]; last.Status != entities.CampaignStatusCompleted {
		t.Fatalf("expected completed, got %s", last.Status)
	}
}

func TestRunDueCampaignsRepeatStartsNewRound(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.DeliveryMode = entities.CampaignDeliveryRepeat
	campaign.SendDates = []string{todayIn(&campaign, -7), todayIn(&campaign, 0), todayIn(&campaign, 7)}
	campaign.CurrentRound = 1
	h.campaigns.due = []entities.Campaign{campaign}
	h.audience.candidates = []entities.CampaignCandidate{{ClientID: 1, Phone: "573001"}, {ClientID: 2, Phone: "573002"}}
	h.sends.pending = []entities.CampaignSend{{ID: 10, Phone: "573001"}, {ID: 11, Phone: "573002"}}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(h.sends.skippedFrom) != 1 || h.sends.skippedFrom[0] != 2 {
		t.Fatalf("leftovers of round 1 must be skipped, got %v", h.sends.skippedFrom)
	}
	if len(h.sends.bulkCreated) != 2 || h.sends.bulkCreated[0].Round != 2 {
		t.Fatalf("round 2 must be materialized again for the same audience, got %+v", h.sends.bulkCreated)
	}
	if len(h.publisher.published) != 2 {
		t.Fatalf("expected the whole audience again, got %d", len(h.publisher.published))
	}
}

func TestRunDueCampaignsRepeatWaitsBetweenRounds(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.DeliveryMode = entities.CampaignDeliveryRepeat
	campaign.SendDates = []string{todayIn(&campaign, 0), todayIn(&campaign, 7)}
	campaign.CurrentRound = 1
	h.campaigns.due = []entities.Campaign{campaign}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, updated := range h.campaigns.updated {
		if updated.Status == entities.CampaignStatusCompleted {
			t.Fatal("a repeat campaign must not complete before its last date")
		}
	}
}

func TestRunDueCampaignsRepeatCompletesAfterLastDate(t *testing.T) {
	h := newHarness()
	campaign := runnableCampaign()
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.DeliveryMode = entities.CampaignDeliveryRepeat
	campaign.SendDates = []string{todayIn(&campaign, -1)}
	h.campaigns.due = []entities.Campaign{campaign}

	if err := h.uc.RunDueCampaigns(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(h.sends.skippedFrom) != 1 || h.sends.skippedFrom[0] != 0 {
		t.Fatalf("all leftovers must be skipped, got %v", h.sends.skippedFrom)
	}
	if last := h.campaigns.updated[len(h.campaigns.updated)-1]; last.Status != entities.CampaignStatusCompleted {
		t.Fatalf("expected completed, got %s", last.Status)
	}
}

func TestResumeRejectsExhaustedSchedule(t *testing.T) {
	h := newHarness()
	campaign := launchableCampaign()
	started := time.Now().AddDate(0, 0, -3)
	campaign.Status = entities.CampaignStatusPaused
	campaign.StartedAt = &started
	campaign.ScheduleMode = entities.CampaignScheduleDates
	campaign.SendDates = []string{todayIn(campaign, -1)}
	h.campaigns.byID[3] = campaign

	if _, err := h.uc.Resume(context.Background(), 3, 26); err == nil {
		t.Fatal("expected error resuming a campaign without future dates")
	}
}

func TestUpdateKeepsRoundAndAudience(t *testing.T) {
	h := newHarness()
	current := launchableCampaign()
	current.Status = entities.CampaignStatusPaused
	current.CurrentRound = 2
	current.AudienceCount = 40
	h.campaigns.byID[3] = current

	updated, err := h.uc.Update(context.Background(), updateDTOFor(3))
	if err != nil {
		t.Fatal(err)
	}
	if updated.CurrentRound != 2 || updated.AudienceCount != 40 {
		t.Fatalf("update lost progress: round %d audience %d", updated.CurrentRound, updated.AudienceCount)
	}
}
