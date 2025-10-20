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

	"github.com/IgorBayerl/AdlerCov/filesystem"
)

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
	Off:     slog.Level(slog.LevelError + 128),
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

type Config struct {
	Verbosity VerbosityLevel        // default Info
	File      string                // "" = console only
	Format    string                // "text" (default) | "json"
	FS        filesystem.Filesystem // if nil, uses filesystem.DefaultFS{}
}

type MultiDestHandler struct {
	handlers []slog.Handler
}

func NewMultiDestHandler(handlers ...slog.Handler) *MultiDestHandler {
	return &MultiDestHandler{
		handlers: handlers,
	}
}

func (h *MultiDestHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiDestHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		_ = handler.Handle(ctx, r)
	}
	return nil
}

func (h *MultiDestHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiDestHandler(newHandlers...)
}

func (h *MultiDestHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return NewMultiDestHandler(newHandlers...)
}

func Init(cfg *Config) (io.Closer, error) {
	if cfg == nil {
		cfg = &Config{Verbosity: Info, Format: "text"}
	}

	fs := cfg.FS
	if fs == nil {
		fs = filesystem.DefaultFS{}
	}

	var handlers []slog.Handler
	var closer io.Closer

	consoleOpts := &slog.HandlerOptions{
		Level: cfg.Verbosity.SlogLevel(),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	handlers = append(handlers, slog.NewTextHandler(os.Stderr, consoleOpts))

	if cfg.File != "" {
		if err := fs.MkdirAll(filepath.Dir(cfg.File), 0o755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		f, err := fs.Create(cfg.File)
		if err != nil {
			return nil, fmt.Errorf("failed to create log file: %w", err)
		}
		closer = f

		fileOpts := &slog.HandlerOptions{Level: cfg.Verbosity.SlogLevel()}
		var fileHandler slog.Handler
		if strings.ToLower(cfg.Format) == "json" {
			fileHandler = slog.NewJSONHandler(f, fileOpts)
		} else {
			fileHandler = slog.NewTextHandler(f, fileOpts)
		}
		handlers = append(handlers, fileHandler)
	}

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
