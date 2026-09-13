package entities

import (
	"testing"
	"time"
)

func bogota(t *testing.T, value string) time.Time {
	t.Helper()
	location, err := time.LoadLocation("America/Bogota")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04", value, location)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func scheduleCampaign(t *testing.T, mode string, start string) *Campaign {
	anchor := bogota(t, start)
	return &Campaign{Timezone: "America/Bogota", ScheduleMode: mode, ScheduledAt: &anchor}
}

func TestDailyOccurrenceEveryDay(t *testing.T) {
	c := scheduleCampaign(t, CampaignScheduleDaily, "2026-09-14 08:00")

	got := c.OccurrenceAt(bogota(t, "2026-09-16 10:00"))
	if !got.SendDay || got.Round != 3 || got.Last || got.Exhausted {
		t.Fatalf("unexpected occurrence: %+v", got)
	}
}

func TestOccurrenceBeforeStartIsIdle(t *testing.T) {
	c := scheduleCampaign(t, CampaignScheduleInterval, "2026-09-14 08:00")
	c.IntervalDays = 8

	got := c.OccurrenceAt(bogota(t, "2026-09-13 10:00"))
	if got.SendDay || got.Exhausted {
		t.Fatalf("unexpected occurrence: %+v", got)
	}
}

func TestIntervalOccurrenceOnlyOnStepDays(t *testing.T) {
	c := scheduleCampaign(t, CampaignScheduleInterval, "2026-09-14 08:00")
	c.IntervalDays = 8

	cases := map[string]CampaignOccurrence{
		"2026-09-14 09:00": {Round: 1, SendDay: true},
		"2026-09-15 09:00": {},
		"2026-09-21 09:00": {},
		"2026-09-22 09:00": {Round: 2, SendDay: true},
		"2026-09-30 09:00": {Round: 3, SendDay: true},
	}

	for when, want := range cases {
		if got := c.OccurrenceAt(bogota(t, when)); got != want {
			t.Fatalf("%s: got %+v want %+v", when, got, want)
		}
	}
}

func TestIntervalWithOccurrencesEnds(t *testing.T) {
	c := scheduleCampaign(t, CampaignScheduleInterval, "2026-09-14 08:00")
	c.IntervalDays = 7
	c.Occurrences = 2

	if got := c.OccurrenceAt(bogota(t, "2026-09-21 09:00")); !got.SendDay || !got.Last || got.Round != 2 {
		t.Fatalf("second week should be the last send: %+v", got)
	}
	if got := c.OccurrenceAt(bogota(t, "2026-09-22 09:00")); !got.Exhausted {
		t.Fatalf("day after the last send should be exhausted: %+v", got)
	}
	if got := c.OccurrenceAt(bogota(t, "2026-09-18 09:00")); got.Exhausted || got.SendDay {
		t.Fatalf("between sends should wait: %+v", got)
	}
}

func TestDatesOccurrence(t *testing.T) {
	c := &Campaign{
		Timezone:     "America/Bogota",
		ScheduleMode: CampaignScheduleDates,
		SendDates:    []string{"2026-09-22", "2026-09-15", "2026-09-15", "basura"},
	}

	if got := c.OccurrenceAt(bogota(t, "2026-09-15 23:30")); !got.SendDay || got.Round != 1 || got.Last {
		t.Fatalf("first date: %+v", got)
	}
	if got := c.OccurrenceAt(bogota(t, "2026-09-18 10:00")); got.SendDay || got.Exhausted {
		t.Fatalf("between dates: %+v", got)
	}
	if got := c.OccurrenceAt(bogota(t, "2026-09-22 10:00")); !got.SendDay || got.Round != 2 || !got.Last {
		t.Fatalf("last date: %+v", got)
	}
	if got := c.OccurrenceAt(bogota(t, "2026-09-23 10:00")); !got.Exhausted {
		t.Fatalf("after last date: %+v", got)
	}
	if c.TotalOccurrences() != 2 {
		t.Fatalf("expected two valid dates, got %d", c.TotalOccurrences())
	}
}

func TestDatesUseCampaignTimezone(t *testing.T) {
	c := &Campaign{
		Timezone:     "America/Bogota",
		ScheduleMode: CampaignScheduleDates,
		SendDates:    []string{"2026-09-15"},
	}

	utcLateNight := time.Date(2026, 9, 16, 2, 0, 0, 0, time.UTC)
	if got := c.OccurrenceAt(utcLateNight); !got.SendDay {
		t.Fatalf("02:00 UTC is still the 15th in Bogota: %+v", got)
	}
}

func TestEmptyDatesAreExhausted(t *testing.T) {
	c := &Campaign{ScheduleMode: CampaignScheduleDates}
	if got := c.OccurrenceAt(time.Now()); !got.Exhausted {
		t.Fatalf("no dates should be exhausted: %+v", got)
	}
}

func TestNormalizeCampaignDates(t *testing.T) {
	got := NormalizeCampaignDates([]string{"2026-10-01", "2026-09-30", "2026-10-01", "x"})
	if len(got) != 2 || got[0] != "2026-09-30" || got[1] != "2026-10-01" {
		t.Fatalf("unexpected dates: %v", got)
	}
}
