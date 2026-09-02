package app

import (
	"strings"
	"time"
)

const dateOnlyLayout = "2006-01-02"

var acceptedDateTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
}

func parseDateValue(value string) (time.Time, bool, bool) {
	v := strings.TrimSpace(value)
	if v == "" {
		return time.Time{}, false, false
	}
	if parsed, err := time.Parse(dateOnlyLayout, v); err == nil {
		return parsed.UTC(), true, true
	}
	for _, layout := range acceptedDateTimeLayouts {
		if parsed, err := time.Parse(layout, v); err == nil {
			return parsed.UTC(), false, true
		}
	}
	return time.Time{}, false, false
}

func parseStartDate(value string) (time.Time, bool) {
	parsed, _, ok := parseDateValue(value)
	return parsed, ok
}

func parseEndDate(value string) (time.Time, bool) {
	parsed, dateOnly, ok := parseDateValue(value)
	if !ok {
		return time.Time{}, false
	}
	if dateOnly {
		parsed = parsed.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}
	return parsed, true
}
