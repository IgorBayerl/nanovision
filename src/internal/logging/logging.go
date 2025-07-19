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
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/filesystem"
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

func Init(cfg *Config) (io.Closer, error) {
	if cfg == nil {
		cfg = &Config{Verbosity: Info, Format: "text"}
	}

	fs := cfg.FS
	if fs == nil {
		fs = filesystem.DefaultFS{}
	}

	// Collect writers
	var writers []io.Writer
	writers = append(writers, os.Stderr)

	var closer io.Closer
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
		writers = append(writers, f)
		closer = f
	}

	var out io.Writer
	if len(writers) == 1 {
		out = writers[0]
	} else {
		out = io.MultiWriter(writers...)
	}

	opts := &slog.HandlerOptions{Level: cfg.Verbosity.SlogLevel()}
	var handler slog.Handler
	if strings.ToLower(cfg.Format) == "json" {
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = slog.NewTextHandler(out, opts)
	}

	slog.SetDefault(slog.New(handler))
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
