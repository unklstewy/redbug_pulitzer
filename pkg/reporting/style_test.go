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
				HeaderColor:   "#3498db", // Blue
				BorderColor:   "#2980b9",
				FooterColor:   "#2980b9",
				AccentColor:   "#1abc9c",
				Mode:          "Read Mode",
				ModeEmoji:     "📥",
				CommandLabel:  "Read Command",
				ResponseLabel: "Radio Response",
			},
		},
		{
			name: "Write Mode",
			mode: WriteMode,
			expected: StyleConfig{
				HeaderColor:   "#e74c3c", // Red
				BorderColor:   "#c0392b",
				FooterColor:   "#c0392b",
				AccentColor:   "#f39c12",
				Mode:          "Write Mode",
				ModeEmoji:     "📤",
				CommandLabel:  "Write Command",
				ResponseLabel: "Radio Response",
			},
		},
		{
			name: "Default Mode (should use neutral styling)",
			mode: Mode("something else"),
			expected: StyleConfig{
				HeaderColor:   "#2c3e50", // Dark blue/gray
				BorderColor:   "#34495e",
				FooterColor:   "#34495e",
				AccentColor:   "#9b59b6",
				Mode:          "Analysis Mode",
				ModeEmoji:     "🔍",
				CommandLabel:  "Command",
				ResponseLabel: "Response",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := GetStyleConfig(tt.mode)

			if style.HeaderColor != tt.expected.HeaderColor {
				t.Errorf("HeaderColor = %v, want %v", style.HeaderColor, tt.expected.HeaderColor)
			}

			if style.Mode != tt.expected.Mode {
				t.Errorf("Mode = %v, want %v", style.Mode, tt.expected.Mode)
			}

			if style.CommandLabel != tt.expected.CommandLabel {
				t.Errorf("CommandLabel = %v, want %v", style.CommandLabel, tt.expected.CommandLabel)
			}
		})
	}
}
