// Path: internal/reporter/context.go
package reporter

import (
	"io"
	"log/slog"

	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/reportconfig"
	"github.com/IgorBayerl/ReportGenerator/go_report_generator/internal/settings"
)

type IBuilderContext interface {
	ReportConfiguration() *reportconfig.ReportConfiguration
	Settings() *settings.Settings
	Logger() *slog.Logger
}

type BuilderContext struct {
	Cfg   *reportconfig.ReportConfiguration
	Stngs *settings.Settings
	L     *slog.Logger
}

func (bc *BuilderContext) ReportConfiguration() *reportconfig.ReportConfiguration { return bc.Cfg }

func (bc *BuilderContext) Settings() *settings.Settings { return bc.Stngs }

func (bc *BuilderContext) Logger() *slog.Logger { return bc.L }

func NewBuilderContext(config *reportconfig.ReportConfiguration, settings *settings.Settings, logger *slog.Logger) *BuilderContext {
	if logger == nil {
		// Default to a discarded logger if none is provided to prevent nil pointer panics.
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return &BuilderContext{
		Cfg:   config,
		Stngs: settings,
		L:     logger,
	}
}
