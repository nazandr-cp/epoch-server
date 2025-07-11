package logging

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/go-pkgz/lgr"
)

const (
	levelTrace = "trace"
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "warn"
	levelError = "error"
	
	formatJSON = "json"
	formatText = "text"
	
	outputStdout = "stdout"
	outputStderr = "stderr"
)

type Config struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
	Output string `yaml:"output" json:"output"`

	CallerInfo       CallerConfig `yaml:"caller" json:"caller"`
	SecretMask       []string     `yaml:"secrets" json:"secrets"`
	StackTraceError  bool         `yaml:"stack_trace_error" json:"stack_trace_error"`
	CustomTemplate   string       `yaml:"template" json:"template"`
	CallerDepth      int          `yaml:"caller_depth" json:"caller_depth"`

	// JSON format specific settings for slog integration
	JSONConfig JSONConfig `yaml:"json" json:"json"`
}

// JSONConfig controls JSON logging behavior using slog handlers
type JSONConfig struct {
	AddSource   bool `yaml:"add_source" json:"add_source"`
	ReplaceAttr bool `yaml:"replace_attr" json:"replace_attr"`
}

// CallerConfig controls caller information in logs
type CallerConfig struct {
	Enabled  bool `yaml:"enabled" json:"enabled"`
	File     bool `yaml:"file" json:"file"`
	Function bool `yaml:"function" json:"function"`
	Package  bool `yaml:"package" json:"package"`
}

func New(level string) lgr.L {
	cfg := Config{
		Level:  level,
		Format: formatText,
		Output: outputStdout,
	}
	logger, err := NewWithConfig(cfg)
	if err != nil {
		return lgr.New(lgr.Debug, lgr.Msec, lgr.LevelBraces)
	}
	return logger
}

func NewWithConfig(cfg Config) (lgr.L, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	var options []lgr.Option
	options = append(options, lgr.Msec)

	switch strings.ToLower(cfg.Level) {
	case levelTrace:
		options = append(options, lgr.Trace)
	case levelDebug:
		options = append(options, lgr.Debug)
	}

	output, err := getOutputWriter(cfg.Output)
	if err != nil {
		return nil, err
	}

	// JSON format uses slog handler for structured logging
	switch strings.ToLower(cfg.Format) {
	case formatJSON:
		jsonHandler := createJSONHandler(cfg, output)
		options = append(options, lgr.SlogHandler(jsonHandler))
	default:
		options = append(options, lgr.LevelBraces, lgr.Out(output))
	}

	// caller information for text format only (JSON uses slog source)
	if strings.ToLower(cfg.Format) != formatJSON {
		if cfg.CallerInfo.Enabled {
			if cfg.CallerInfo.File {
				options = append(options, lgr.CallerFile)
			}
			if cfg.CallerInfo.Function {
				options = append(options, lgr.CallerFunc)
			}
			if cfg.CallerInfo.Package {
				options = append(options, lgr.CallerPkg)
			}
			if cfg.CallerDepth > 0 {
				options = append(options, lgr.CallerDepth(cfg.CallerDepth))
			}
		} else {
			level := strings.ToLower(cfg.Level)
			if level == levelTrace || level == levelDebug {
				options = append(options, lgr.CallerFile, lgr.CallerFunc)
			}
		}
	}

	// text format options (JSON handles these through slog attributes)
	isJSON := strings.ToLower(cfg.Format) == formatJSON
	if !isJSON {
		if len(cfg.SecretMask) > 0 {
			options = append(options, lgr.Secret(cfg.SecretMask...))
		}
		if cfg.StackTraceError {
			options = append(options, lgr.StackTraceOnError)
		}
		if cfg.CustomTemplate != "" {
			options = append(options, lgr.Format(cfg.CustomTemplate))
		}
		if strings.ToLower(cfg.Output) != outputStderr {
			options = append(options, lgr.Err(os.Stderr))
		}
	}

	return lgr.New(options...), nil
}

// createJSONHandler creates a slog JSON handler with mapped levels and custom attributes
func createJSONHandler(cfg Config, output io.Writer) *slog.JSONHandler {
	var slogLevel slog.Level
	switch strings.ToLower(cfg.Level) {
	case levelTrace, levelDebug:
		slogLevel = slog.LevelDebug
	case levelInfo:
		slogLevel = slog.LevelInfo
	case levelWarn:
		slogLevel = slog.LevelWarn
	case levelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: cfg.JSONConfig.AddSource,
	}

	if cfg.JSONConfig.ReplaceAttr {
		handlerOptions.ReplaceAttr = createReplaceAttrFunc(cfg)
	}

	return slog.NewJSONHandler(output, handlerOptions)
}

// createReplaceAttrFunc handles secret masking and time formatting for JSON logs
func createReplaceAttrFunc(cfg Config) func(groups []string, a slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		// mask secrets in messages for security
		if len(cfg.SecretMask) > 0 && a.Key == slog.MessageKey {
			value := a.Value.String()
			for _, secret := range cfg.SecretMask {
				if strings.Contains(value, secret) {
					value = strings.ReplaceAll(value, secret, "[REDACTED]")
				}
			}
			return slog.Attr{
				Key:   a.Key,
				Value: slog.StringValue(value),
			}
		}

		// standardize timestamp format for consistency
		if a.Key == slog.TimeKey {
			return slog.Attr{
				Key:   a.Key,
				Value: slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z07:00")),
			}
		}

		return a
	}
}

func validateConfig(cfg Config) error {
	level := strings.ToLower(cfg.Level)
	validLevels := []string{levelTrace, levelDebug, levelInfo, levelWarn, levelError}
	if level != "" && !contains(validLevels, level) {
		return errors.New("invalid log level: " + cfg.Level + ", must be one of: trace, debug, info, warn, error")
	}

	format := strings.ToLower(cfg.Format)
	validFormats := []string{formatText, formatJSON}
	if format != "" && !contains(validFormats, format) {
		return errors.New("invalid log format: " + cfg.Format + ", must be one of: text, json")
	}

	if cfg.CallerDepth < 0 {
		return errors.New("caller depth must be non-negative, got: " + strconv.Itoa(cfg.CallerDepth))
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getOutputWriter(output string) (io.Writer, error) {
	switch strings.ToLower(output) {
	case "", outputStdout:
		return os.Stdout, nil
	case outputStderr:
		return os.Stderr, nil
	default:
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return nil, errors.New("failed to open log file " + output + ": " + err.Error())
		}
		return file, nil
	}
}
