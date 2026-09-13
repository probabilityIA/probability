package entities

import (
	"sort"
	"time"
)

const campaignDateLayout = "2006-01-02"

type CampaignOccurrence struct {
	Round     uint
	SendDay   bool
	Last      bool
	Exhausted bool
}

func (c *Campaign) Location() *time.Location {
	location, err := time.LoadLocation(c.Timezone)
	if err != nil {
		return time.UTC
	}
	return location
}

func (c *Campaign) TotalOccurrences() uint {
	if c.ScheduleMode == CampaignScheduleDates {
		return uint(len(NormalizeCampaignDates(c.SendDates)))
	}
	return c.Occurrences
}

func (c *Campaign) OccurrenceAt(now time.Time) CampaignOccurrence {
	location := c.Location()
	today := now.In(location).Format(campaignDateLayout)

	if c.ScheduleMode == CampaignScheduleDates {
		return c.dateOccurrence(today)
	}

	step := uint(1)
	if c.ScheduleMode == CampaignScheduleInterval && c.IntervalDays > 1 {
		step = c.IntervalDays
	}

	anchor := now
	if c.ScheduledAt != nil {
		anchor = *c.ScheduledAt
	} else if c.StartedAt != nil {
		anchor = *c.StartedAt
	}

	days := daysBetween(anchor.In(location).Format(campaignDateLayout), today)
	if days < 0 {
		return CampaignOccurrence{}
	}

	index := uint(days) / step
	sendDay := uint(days)%step == 0
	total := c.Occurrences

	if total > 0 && (index+1 > total || (index+1 == total && !sendDay)) {
		return CampaignOccurrence{Exhausted: true}
	}

	if !sendDay {
		return CampaignOccurrence{}
	}

	return CampaignOccurrence{
		Round:   index + 1,
		SendDay: true,
		Last:    total > 0 && index+1 == total,
	}
}

func (c *Campaign) dateOccurrence(today string) CampaignOccurrence {
	dates := NormalizeCampaignDates(c.SendDates)
	if len(dates) == 0 || today > dates[len(dates)-1] {
		return CampaignOccurrence{Exhausted: true}
	}

	for i, date := range dates {
		if date == today {
			return CampaignOccurrence{
				Round:   uint(i + 1),
				SendDay: true,
				Last:    i == len(dates)-1,
			}
		}
	}

	return CampaignOccurrence{}
}

func NormalizeCampaignDates(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))

	for _, value := range values {
		parsed, err := time.Parse(campaignDateLayout, value)
		if err != nil {
			continue
		}
		date := parsed.Format(campaignDateLayout)
		if _, ok := seen[date]; ok {
			continue
		}
		seen[date] = struct{}{}
		out = append(out, date)
	}

	sort.Strings(out)
	return out
}

func IsValidCampaignDate(value string) bool {
	_, err := time.Parse(campaignDateLayout, value)
	return err == nil
}

func daysBetween(from, to string) int {
	start, err := time.Parse(campaignDateLayout, from)
	if err != nil {
		return 0
	}
	end, err := time.Parse(campaignDateLayout, to)
	if err != nil {
		return 0
	}
	return int(end.Sub(start).Hours() / 24)
}
