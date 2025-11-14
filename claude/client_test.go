package claude

import (
	"encoding/json"
	"testing"

	"github.com/truong-nautilus/smart-home-ai/core"
)

func TestCommandParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected core.Command
		wantErr  bool
	}{
		{
			name:  "Valid light on command",
			input: `{"action":"light.on","device":"phong_khach","value":100}`,
			expected: core.Command{
				Action: "light.on",
				Device: "phong_khach",
				Value:  float64(100),
			},
			wantErr: false,
		},
		{
			name:  "Valid AC command",
			input: `{"action":"ac.set_temp","device":"dieu_hoa_phong_khach","value":26}`,
			expected: core.Command{
				Action: "ac.set_temp",
				Device: "dieu_hoa_phong_khach",
				Value:  float64(26),
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `{"action":"light.on"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cmd core.Command
			err := json.Unmarshal([]byte(tt.input), &cmd)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if cmd.Action != tt.expected.Action {
				t.Errorf("Action = %v, want %v", cmd.Action, tt.expected.Action)
			}

			if cmd.Device != tt.expected.Device {
				t.Errorf("Device = %v, want %v", cmd.Device, tt.expected.Device)
			}
		})
	}
}

func TestClaudeConfig(t *testing.T) {
	config := ClaudeConfig{
		APIKey:       "test-key",
		Model:        "claude-3-5-sonnet-20241022",
		WebSocketURL: "wss://api.anthropic.com/v1/messages",
		MaxTokens:    1024,
		Temperature:  0.7,
	}

	if config.Model != "claude-3-5-sonnet-20241022" {
		t.Errorf("Expected model claude-3-5-sonnet-20241022, got %s", config.Model)
	}

	if config.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", config.Temperature)
	}
}
