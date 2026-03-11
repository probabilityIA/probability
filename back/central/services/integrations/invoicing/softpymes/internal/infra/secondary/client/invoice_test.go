package client

import (
	"testing"
)

func TestResolveItemCode(t *testing.T) {
	productID := "PROD-123"

	tests := []struct {
		name         string
		itemSKU      string
		itemName     string
		productID    *string
		itemMappings map[string]interface{}
		expected     string
	}{
		{
			name:         "nil mappings returns SKU",
			itemSKU:      "ABC-001",
			itemName:     "Blue T-Shirt",
			productID:    nil,
			itemMappings: nil,
			expected:     "ABC-001",
		},
		{
			name:         "nil mappings with empty SKU falls back to productID",
			itemSKU:      "",
			itemName:     "Product",
			productID:    &productID,
			itemMappings: nil,
			expected:     "PROD-123",
		},
		{
			name:         "nil mappings with empty SKU and nil productID returns empty",
			itemSKU:      "",
			itemName:     "",
			productID:    nil,
			itemMappings: nil,
			expected:     "",
		},
		{
			name:     "regular product uses SKU even with mappings",
			itemSKU:  "PT01001",
			itemName: "Blue T-Shirt",
			itemMappings: map[string]interface{}{
				"tip":        "SE03001",
				"membership": "SE01001",
			},
			expected: "PT01001",
		},
		// Tip matching by name
		{
			name:     "item named 'Tip' uses tip code",
			itemSKU:  "VAR-0",
			itemName: "Tip",
			itemMappings: map[string]interface{}{
				"tip": "SE03001",
			},
			expected: "SE03001",
		},
		{
			name:     "item named 'tip' (lowercase) uses tip code",
			itemSKU:  "VAR-0",
			itemName: "tip",
			itemMappings: map[string]interface{}{
				"tip": "SE03001",
			},
			expected: "SE03001",
		},
		{
			name:     "item named 'Propina' uses tip code",
			itemSKU:  "PROP-1",
			itemName: "Propina",
			itemMappings: map[string]interface{}{
				"tip": "SE03001",
			},
			expected: "SE03001",
		},
		{
			name:     "tip with empty code falls back to SKU",
			itemSKU:  "VAR-0",
			itemName: "Tip",
			itemMappings: map[string]interface{}{
				"tip": "",
			},
			expected: "VAR-0",
		},
		// Membership matching by name
		{
			name:     "item named 'Membresía Anual' uses membership code",
			itemSKU:  "SE01001",
			itemName: "Membresía Anual",
			itemMappings: map[string]interface{}{
				"membership": "MEM-SOFT",
			},
			expected: "MEM-SOFT",
		},
		{
			name:     "item named 'Annual Membership' uses membership code",
			itemSKU:  "MEM-123",
			itemName: "Annual Membership",
			itemMappings: map[string]interface{}{
				"membership": "SE01001",
			},
			expected: "SE01001",
		},
		{
			name:     "item named 'Membresia Premium' (no accent) uses membership code",
			itemSKU:  "MEM-P",
			itemName: "Membresia Premium",
			itemMappings: map[string]interface{}{
				"membership": "SE01001",
			},
			expected: "SE01001",
		},
		{
			name:     "membership with empty code falls back to SKU",
			itemSKU:  "SE01001",
			itemName: "Membresía Anual",
			itemMappings: map[string]interface{}{
				"membership": "",
			},
			expected: "SE01001",
		},
		// Non-matching names
		{
			name:     "item named 'Makeup Tips' does NOT match tip (not exact)",
			itemSKU:  "MK-001",
			itemName: "Makeup Tips",
			itemMappings: map[string]interface{}{
				"tip": "SE03001",
			},
			expected: "MK-001",
		},
		{
			name:      "empty name with empty SKU falls back to productID",
			itemSKU:   "",
			itemName:  "",
			productID: &productID,
			itemMappings: map[string]interface{}{
				"tip": "SE03001",
			},
			expected: "PROD-123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveItemCode(tt.itemSKU, tt.itemName, tt.productID, tt.itemMappings)
			if result != tt.expected {
				t.Errorf("resolveItemCode(%q, %q, %v, %v) = %q, want %q",
					tt.itemSKU, tt.itemName, tt.productID, tt.itemMappings, result, tt.expected)
			}
		})
	}
}

func TestResolveShippingItemCode(t *testing.T) {
	tests := []struct {
		name         string
		itemMappings map[string]interface{}
		expected     string
	}{
		{
			name:         "nil mappings returns SHIPPING",
			itemMappings: nil,
			expected:     "SHIPPING",
		},
		{
			name:         "empty mappings returns SHIPPING",
			itemMappings: map[string]interface{}{},
			expected:     "SHIPPING",
		},
		{
			name: "shipping code configured",
			itemMappings: map[string]interface{}{
				"shipping": "SE02001",
			},
			expected: "SE02001",
		},
		{
			name: "empty shipping code falls back to SHIPPING",
			itemMappings: map[string]interface{}{
				"shipping": "",
			},
			expected: "SHIPPING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolveShippingItemCode(tt.itemMappings)
			if result != tt.expected {
				t.Errorf("resolveShippingItemCode(%v) = %q, want %q",
					tt.itemMappings, result, tt.expected)
			}
		})
	}
}
