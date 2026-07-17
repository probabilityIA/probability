package response

import (
	"strings"
	"time"
)

type JSTime struct {
	time.Time
}

const jumpsellerTimeLayout = "2006-01-02 15:04:05 MST"

func (t *JSTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		t.Time = time.Time{}
		return nil
	}
	parsed, err := time.Parse(jumpsellerTimeLayout, s)
	if err != nil {
		parsed, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	t.Time = parsed
	return nil
}
