package app

import (
	"strings"
	"testing"
)


func TestBuildSubject(t *testing.T) {
	tests := []struct {
		eventType string
		expected  string
	}{
		{"order.created", "Notificación: order.created"},
		{"order.shipped", "Notificación: order.shipped"},
		{"", "Notificación: "},
	}

	for _, tt := range tests {
		result := buildSubject(tt.eventType)
		if result != tt.expected {
			t.Errorf("buildSubject(%q) = %q, want %q", tt.eventType, result, tt.expected)
		}
	}
}


func TestBuildHTML_ContainsEventType(t *testing.T) {
	html := buildHTML("order.created", nil)

	if !strings.Contains(html, "Evento: order.created") {
		t.Error("expected HTML to contain event type heading")
	}
}

func TestBuildHTML_ContainsHTMLStructure(t *testing.T) {
	html := buildHTML("test", nil)

	checks := []string{
		"<!DOCTYPE html>",
		"<html>",
		"</html>",
		"<body",
		"</body>",
		"charset=\"UTF-8\"",
		"Este es un mensaje automático",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("expected HTML to contain %q", check)
		}
	}
}

func TestBuildHTML_WithEventData_RendersTable(t *testing.T) {
	data := map[string]interface{}{
		"order_id": "ABC123",
	}

	html := buildHTML("order.created", data)

	if !strings.Contains(html, "<table") {
		t.Error("expected HTML to contain table when eventData is present")
	}
	if !strings.Contains(html, "order_id") {
		t.Error("expected HTML to contain key 'order_id'")
	}
	if !strings.Contains(html, "ABC123") {
		t.Error("expected HTML to contain value 'ABC123'")
	}
}

func TestBuildHTML_WithoutEventData_NoTable(t *testing.T) {
	html := buildHTML("order.created", nil)

	if strings.Contains(html, "<table") {
		t.Error("expected HTML NOT to contain table when eventData is nil")
	}
}

func TestBuildHTML_EmptyEventData_NoTable(t *testing.T) {
	html := buildHTML("order.created", map[string]interface{}{})

	if strings.Contains(html, "<table") {
		t.Error("expected HTML NOT to contain table when eventData is empty")
	}
}
