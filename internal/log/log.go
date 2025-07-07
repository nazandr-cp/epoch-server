package log

import (
	"io"
	"os"
	"strings"

	"github.com/go-pkgz/lgr"
)

type Config struct {
	Level  string
	Format string
	Output string
}

func New(level string) lgr.L {
	cfg := Config{
		Level:  level,
		Format: "text", // default format
		Output: "stdout", // default output
	}
	return NewWithConfig(cfg)
}

func NewWithConfig(cfg Config) lgr.L {
	var options []lgr.Option
	
	// Add basic timestamp with milliseconds for transaction tracking
	options = append(options, lgr.Msec)
	
	// Set log level
	switch strings.ToLower(cfg.Level) {
	case "trace":
		options = append(options, lgr.Trace)
	case "debug":
		options = append(options, lgr.Debug)
	case "info", "warn", "error":
		// INFO is the default level in lgr
	}
	
	// Configure format-specific options
	switch strings.ToLower(cfg.Format) {
	case "json":
		// For JSON format, don't use level braces as they interfere with JSON structure
		// The application manually prefixes levels like "INFO", "ERROR" in log messages
	case "text":
		// For text format, use level braces for better readability
		options = append(options, lgr.LevelBraces)
	default:
		// Default to text format with level braces
		options = append(options, lgr.LevelBraces)
	}
	
	// Set output destination
	var output io.Writer = os.Stdout
	switch strings.ToLower(cfg.Output) {
	case "stderr":
		output = os.Stderr
	case "stdout":
		output = os.Stdout
	default:
		// Could be a file path - for now default to stdout
		output = os.Stdout
	}
	options = append(options, lgr.Out(output))
	
	// Add caller information for debug and trace levels (helpful for blockchain debugging)
	level := strings.ToLower(cfg.Level)
	if level == "trace" || level == "debug" {
		options = append(options, lgr.CallerFile, lgr.CallerFunc)
	}
	
	// For error output, also send to stderr in addition to main output
	// This is useful for container environments where errors should go to stderr
	if strings.ToLower(cfg.Output) != "stderr" {
		options = append(options, lgr.Err(os.Stderr))
	}
	
	return lgr.New(options...)
}
