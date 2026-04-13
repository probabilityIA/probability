package entities

import "testing"

func TestDisplayTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		dt       DisplayType
		expected string
	}{
		{"modal_image", DisplayTypeModalImage, "modal_image"},
		{"modal_text", DisplayTypeModalText, "modal_text"},
		{"ticker", DisplayTypeTicker, "ticker"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.dt) != tt.expected {
				t.Errorf("got %s, want %s", tt.dt, tt.expected)
			}
		})
	}
}

func TestFrequencyTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		ft       FrequencyType
		expected string
	}{
		{"once", FrequencyOnce, "once"},
		{"daily", FrequencyDaily, "daily"},
		{"always", FrequencyAlways, "always"},
		{"requires_acceptance", FrequencyRequiresAcceptance, "requires_acceptance"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.ft) != tt.expected {
				t.Errorf("got %s, want %s", tt.ft, tt.expected)
			}
		})
	}
}

func TestAnnouncementStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		status   AnnouncementStatus
		expected string
	}{
		{"draft", StatusDraft, "draft"},
		{"scheduled", StatusScheduled, "scheduled"},
		{"active", StatusActive, "active"},
		{"inactive", StatusInactive, "inactive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("got %s, want %s", tt.status, tt.expected)
			}
		})
	}
}

func TestViewActionConstants(t *testing.T) {
	tests := []struct {
		name     string
		action   ViewAction
		expected string
	}{
		{"viewed", ViewActionViewed, "viewed"},
		{"closed", ViewActionClosed, "closed"},
		{"clicked_link", ViewActionClickedLink, "clicked_link"},
		{"accepted", ViewActionAccepted, "accepted"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.action) != tt.expected {
				t.Errorf("got %s, want %s", tt.action, tt.expected)
			}
		})
	}
}
