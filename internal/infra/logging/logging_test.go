package logging

import (
	"testing"
)

func TestNew(t *testing.T) {
	logger := New("debug")
	if logger == nil {
		t.Fatal("logger should not be nil")
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "debug level with stdout",
			cfg: Config{
				Level:  "debug",
				Format: "text",
				Output: "stdout",
			},
		},
		{
			name: "info level with stderr",
			cfg: Config{
				Level:  "info",
				Format: "json",
				Output: "stderr",
			},
		},
		{
			name: "trace level",
			cfg: Config{
				Level:  "trace",
				Format: "text",
				Output: "stdout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewWithConfig(tt.cfg)
			if logger == nil {
				t.Fatal("logger should not be nil")
			}

			// Test that logger can log messages
			logger.Logf("INFO test message for %s", tt.name)
		})
	}
}
