// Package logging provides a standardized and configurable logging setup,
// built on top of the standard library's `log/slog` package.
//
// The primary goal is to abstract away the boilerplate of configuring slog for a
// command-line application. It provides user-friendly verbosity levels, command-line
// flag integration, and simultaneous console (stderr) and file output.
//
// # Key Features
//
//   - Flag Integration: Wires logging config (-verbosity, -logfile, -logformat) to flags.
//   - Dual Output: Logs to stderr and optionally to a file.
//   - Formats: Supports human-readable "text" and machine-parseable "json" output.
//   - Testability: Uses a filesystem abstraction for easy testing.

package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/AdlerCov/internal/filesystem"
)

// Verbosity enum & helpers
type VerbosityLevel int

const (
	Verbose VerbosityLevel = iota
	Info
	Warning
	Error
	Off
)

var verbosityToSlog = map[VerbosityLevel]slog.Level{
	Verbose: slog.LevelDebug,
	Info:    slog.LevelInfo,
	Warning: slog.LevelWarn,
	Error:   slog.LevelError,
	Off:     slog.Level(slog.LevelError + 128), // silence
}

func (v VerbosityLevel) SlogLevel() slog.Level { return verbosityToSlog[v] }

func ParseVerbosity(s string) (VerbosityLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "verbose":
		return Verbose, nil
	case "info":
		return Info, nil
	case "warning", "warn":
		return Warning, nil
	case "error":
		return Error, nil
	case "off", "silent":
		return Off, nil
	default:
		return Info, fmt.Errorf("invalid verbosity level %q", s)
	}
}

// Config + Flag wiring

type Config struct {
	Verbosity VerbosityLevel        // default Info
	File      string                // "" = console only
	Format    string                // "text" (default) | "json"
	FS        filesystem.Filesystem // if nil, uses filesystem.DefaultFS{}
}

// --- NEW: Multi-destination slog.Handler ---
// MultiDestHandler sends log records to multiple underlying handlers.
type MultiDestHandler struct {
	handlers []slog.Handler
}

// NewMultiDestHandler creates a new handler that delegates to the provided handlers.
func NewMultiDestHandler(handlers ...slog.Handler) *MultiDestHandler {
	return &MultiDestHandler{
		handlers: handlers,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler is enabled if any of its sub-handlers is enabled.
func (h *MultiDestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle handles the log record by passing it to all its sub-handlers.
func (h *MultiDestHandler) Handle(ctx context.Context, r slog.Record) error {
	// In a production system, you might want to aggregate errors.
	// For this use case, handling the record on each handler is sufficient.
	for _, handler := range h.handlers {
		// We ignore the error from sub-handlers for simplicity.
		_ = handler.Handle(ctx, r)
	}
	return nil
}

// WithAttrs returns a new MultiDestHandler whose sub-handlers have the given attributes.
func (h *MultiDestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiDestHandler(newHandlers...)
}

// WithGroup returns a new MultiDestHandler whose sub-handlers have the given group.
func (h *MultiDestHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return NewMultiDestHandler(newHandlers...)
}

// --- MODIFIED: Init function ---
// Init initializes the global logger with separate configurations for console and file.
func Init(cfg *Config) (io.Closer, error) {
	if cfg == nil {
		cfg = &Config{Verbosity: Info, Format: "text"}
	}

	fs := cfg.FS
	if fs == nil {
		fs = filesystem.DefaultFS{}
	}

	var handlers []slog.Handler
	var closer io.Closer // This will hold the file handle to be closed on exit

	// 1. Configure the "clean" console handler for os.Stderr.
	// We use ReplaceAttr to remove the timestamp.
	consoleOpts := &slog.HandlerOptions{
		Level: cfg.Verbosity.SlogLevel(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the timestamp from console logs
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	// The console output is always text for better readability.
	handlers = append(handlers, slog.NewTextHandler(os.Stderr, consoleOpts))

	// 2. Configure the "complete" file handler if a log file is specified.
	if cfg.File != "" {
		// Create directory if it doesn't exist
		if err := fs.MkdirAll(filepath.Dir(cfg.File), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// Use filesystem abstraction to create the file
		f, err := fs.Create(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %w", err)
		}
		closer = f // Assign the file to the closer

		// This handler does NOT have ReplaceAttr, so it will be complete.
		fileOpts := &slog.HandlerOptions{Level: cfg.Verbosity.SlogLevel()}
		var fileHandler slog.Handler
		if strings.ToLower(cfg.Format) == "json" {
			fileHandler = slog.NewJSONHandler(f, fileOpts)
		} else {
			fileHandler = slog.NewTextHandler(f, fileOpts)
		}
		handlers = append(handlers, fileHandler)
	}

	// 3. Create the multi-destination handler and set it as the default.
	combinedHandler := NewMultiDestHandler(handlers...)
	slog.SetDefault(slog.New(combinedHandler))

	return closer, nil
}

func InitWithFS(fs filesystem.Filesystem, verbosity VerbosityLevel, file, format string) (io.Closer, error) {
	cfg := &Config{
		Verbosity: verbosity,
		File:      file,
		Format:    format,
		FS:        fs,
	}
	return Init(cfg)
}

func Nop() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 128}))
}
