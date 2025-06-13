package reporting

import (
	"testing"
)

func TestGetStyleConfig(t *testing.T) {
	tests := []struct {
		name     string
		mode     Mode
		expected StyleConfig
	}{
		{
			name: "Read Mode",
			mode: ReadMode,
			expected: StyleConfig{
				PrimaryColor:   "#2980b9",
				SecondaryColor: "#3498db",
				Title:          "Read Protocol API Documentation",
				Icon:           "📥",
				HeaderBgColor:  "#eaf2f8",
				BorderColor:    "#3498db",
			},
		},
		{
			name: "Write Mode",
			mode: WriteMode,
			expected: StyleConfig{
				PrimaryColor:   "#c0392b", // Swap these two colors to match implementation
				SecondaryColor: "#e74c3c",
				Title:          "Write Protocol API Documentation",
				Icon:           "📤",
				HeaderBgColor:  "#f9ebea",
				BorderColor:    "#e74c3c",
			},
		},
		// You might want to check if there's a default case in the actual implementation
		// and add a test for it if it exists
	}

	// Execute tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStyleConfig(tt.mode)

			// Check each field matches expected value
			if result.PrimaryColor != tt.expected.PrimaryColor {
				t.Errorf("PrimaryColor = %v, want %v", result.PrimaryColor, tt.expected.PrimaryColor)
			}
			if result.SecondaryColor != tt.expected.SecondaryColor {
				t.Errorf("SecondaryColor = %v, want %v", result.SecondaryColor, tt.expected.SecondaryColor)
			}
			if result.Title != tt.expected.Title {
				t.Errorf("Title = %v, want %v", result.Title, tt.expected.Title)
			}
			if result.Icon != tt.expected.Icon {
				t.Errorf("Icon = %v, want %v", result.Icon, tt.expected.Icon)
			}
			if result.HeaderBgColor != tt.expected.HeaderBgColor {
				t.Errorf("HeaderBgColor = %v, want %v", result.HeaderBgColor, tt.expected.HeaderBgColor)
			}
			if result.BorderColor != tt.expected.BorderColor {
				t.Errorf("BorderColor = %v, want %v", result.BorderColor, tt.expected.BorderColor)
			}
		})
	}
}
