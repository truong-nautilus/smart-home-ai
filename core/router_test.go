package core

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Valid light on",
			input:   `{"action":"light.on","device":"phong_khach"}`,
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
			_, err := ParseCommand(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecurityManager(t *testing.T) {
	security := NewSecurityManager()
	cmd := &Command{
		Action: "light.on",
		Device: "test",
	}
	if err := security.ValidateCommand(cmd); err != nil {
		t.Errorf("ValidateCommand() error = %v", err)
	}
}
