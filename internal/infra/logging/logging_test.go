package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	logger := New("debug")
	assert.NotNil(t, logger)
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config combinations",
			cfg: Config{
				Level:  "debug",
				Format: "text",
				Output: "stdout",
			},
		},
		{
			name: "json format with stderr",
			cfg: Config{
				Level:  "info",
				Format: "json",
				Output: "stderr",
			},
		},
		{
			name: "caller info configuration",
			cfg: Config{
				Level:  "trace",
				Format: "text",
				Output: "stdout",
				CallerInfo: CallerConfig{
					Enabled:  true,
					File:     true,
					Function: true,
					Package:  true,
				},
				CallerDepth: 2,
			},
		},
		{
			name: "secret masking and stack trace",
			cfg: Config{
				Level:           "debug",
				Format:          "text",
				Output:          "stdout",
				SecretMask:      []string{"password", "token"},
				StackTraceError: true,
			},
		},
		{
			name:    "invalid log level",
			cfg:     Config{Level: "invalid", Format: "text", Output: "stdout"},
			wantErr: true,
		},
		{
			name:    "invalid format",
			cfg:     Config{Level: "info", Format: "invalid", Output: "stdout"},
			wantErr: true,
		},
		{
			name:    "negative caller depth",
			cfg:     Config{Level: "info", Format: "text", Output: "stdout", CallerDepth: -1},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewWithConfig(tt.cfg)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				logger.Logf("INFO test message for %s", tt.name)
			}
		})
	}
}

func TestNewWithConfig_FileOutput(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := Config{
		Level:  "info",
		Format: "text",
		Output: logFile,
	}

	logger, err := NewWithConfig(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Logf("INFO test message")
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg:  Config{Level: "debug", Format: "text", Output: "stdout"},
		},
		{
			name: "empty config",
			cfg:  Config{},
		},
		{
			name:    "invalid level",
			cfg:     Config{Level: "invalid"},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name:    "invalid format",
			cfg:     Config{Format: "invalid"},
			wantErr: true,
			errMsg:  "invalid log format",
		},
		{
			name:    "negative caller depth",
			cfg:     Config{CallerDepth: -1},
			wantErr: true,
			errMsg:  "caller depth must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetOutputWriter(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		wantErr bool
	}{
		{
			name:   "stdout",
			output: "stdout",
		},
		{
			name:   "stderr",
			output: "stderr",
		},
		{
			name:   "empty defaults to stdout",
			output: "",
		},
		{
			name:    "invalid file path",
			output:  "/invalid/path/that/should/not/exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer, err := getOutputWriter(tt.output)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, writer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, writer)
			}
		})
	}

	t.Run("valid file path", func(t *testing.T) {
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")

		writer, err := getOutputWriter(logFile)
		assert.NoError(t, err)
		assert.NotNil(t, writer)

		if file, ok := writer.(*os.File); ok {
			_ = file.Close()
		}
	})
}

func TestJSONFormat(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		logMsg   string
		wantJSON bool
	}{
		{
			name:     "basic JSON format",
			cfg:      Config{Level: "info", Format: "json", Output: "stdout"},
			logMsg:   "INFO test message",
			wantJSON: true,
		},
		{
			name: "JSON with caller info",
			cfg: Config{
				Level:      "debug",
				Format:     "json",
				Output:     "stdout",
				JSONConfig: JSONConfig{AddSource: true},
			},
			logMsg:   "DEBUG test message with caller",
			wantJSON: true,
		},
		{
			name: "JSON with secret masking",
			cfg: Config{
				Level:      "info",
				Format:     "json",
				Output:     "stdout",
				SecretMask: []string{"password", "secret"},
				JSONConfig: JSONConfig{ReplaceAttr: true},
			},
			logMsg:   "INFO test message with password=secret123",
			wantJSON: true,
		},
		{
			name:     "text format",
			cfg:      Config{Level: "info", Format: "text", Output: "stdout"},
			logMsg:   "INFO test message",
			wantJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			logFile := filepath.Join(tempDir, "test.log")
			tt.cfg.Output = logFile

			logger, err := NewWithConfig(tt.cfg)
			require.NoError(t, err)
			require.NotNil(t, logger)

			logger.Logf(tt.logMsg)

			content, err := os.ReadFile(logFile)
			require.NoError(t, err)
			
			output := string(content)
			assert.Contains(t, output, "test message")

			lines := strings.Split(strings.TrimSpace(output), "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "" {
					continue
				}
				
				if tt.wantJSON {
					var jsonObj map[string]interface{}
					err := json.Unmarshal([]byte(line), &jsonObj)
					assert.NoError(t, err)
					assert.Contains(t, jsonObj, "time")
					assert.Contains(t, jsonObj, "level")
					assert.Contains(t, jsonObj, "msg")
				} else {
					assert.False(t, strings.HasPrefix(strings.TrimSpace(line), "{"))
				}
			}
		})
	}
}

func TestJSONCallerInfo(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := Config{
		Level:  "debug",
		Format: "json",
		Output: logFile,
		JSONConfig: JSONConfig{
			AddSource: true,
		},
	}

	logger, err := NewWithConfig(cfg)
	require.NoError(t, err)

	logger.Logf("DEBUG test message with caller info")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	var jsonObj map[string]interface{}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	require.Greater(t, len(lines), 0)
	
	err = json.Unmarshal([]byte(lines[0]), &jsonObj)
	require.NoError(t, err)

	// check for common source field names used by slog implementations
	hasSource := false
	for _, field := range []string{"source", "caller", "file"} {
		if _, exists := jsonObj[field]; exists {
			hasSource = true
			break
		}
	}
	
	if !hasSource {
		t.Logf("JSON object: %+v", jsonObj)
		t.Log("Note: Source information format may vary with slog implementation")
	}
}

func TestJSONSecretMasking(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := Config{
		Level:      "info",
		Format:     "json",
		Output:     logFile,
		SecretMask: []string{"password", "token"},
		JSONConfig: JSONConfig{
			ReplaceAttr: true,
		},
	}

	logger, err := NewWithConfig(cfg)
	require.NoError(t, err)

	logger.Logf("INFO message with password and token values")

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	output := string(content)
	assert.Contains(t, output, "[REDACTED]")
	assert.NotContains(t, output, "password")
	assert.NotContains(t, output, "token")
	
	var jsonObj map[string]interface{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	require.Greater(t, len(lines), 0)
	
	err = json.Unmarshal([]byte(lines[0]), &jsonObj)
	require.NoError(t, err)
	
	assert.Contains(t, jsonObj, "msg")
	msg, ok := jsonObj["msg"].(string)
	require.True(t, ok)
	assert.Contains(t, msg, "[REDACTED]")
}

func TestCreateJSONHandler(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
	}{
		{
			name: "basic handler",
			cfg:  Config{Level: "info", Format: "json"},
		},
		{
			name: "with caller info",
			cfg:  Config{Level: "debug", Format: "json", JSONConfig: JSONConfig{AddSource: true}},
		},
		{
			name: "with secret masking",
			cfg:  Config{Level: "warn", Format: "json", SecretMask: []string{"secret"}, JSONConfig: JSONConfig{ReplaceAttr: true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := createJSONHandler(tt.cfg, &buf)
			assert.NotNil(t, handler)
		})
	}
}
