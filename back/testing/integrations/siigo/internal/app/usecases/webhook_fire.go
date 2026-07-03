package usecases

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func (s *APISimulator) FireWebhooks(topic, code string) (int, []string) {
	webhooks := s.Repository.ListWebhooks()

	payload := map[string]interface{}{
		"topic":       topic,
		"company_key": "MOCKCOMPANYSAS",
		"resource": map[string]string{
			"code": code,
		},
		"fired_at": time.Now().UTC().Format(time.RFC3339),
	}
	body, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 15 * time.Second}
	fired := 0
	errs := []string{}

	for _, w := range webhooks {
		if topic != "" && w.Topic != topic {
			continue
		}
		req, err := http.NewRequest(http.MethodPost, w.URL, bytes.NewReader(body))
		if err != nil {
			errs = append(errs, w.URL+": "+err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			errs = append(errs, w.URL+": "+err.Error())
			continue
		}
		resp.Body.Close()
		fired++
	}

	return fired, errs
}
