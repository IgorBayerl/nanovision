package reporter

import (
	"io"
	"log/slog"

	"github.com/IgorBayerl/AdlerCov/internal/config"
)

// IBuilderContext defines the contract for the context passed to report builders.
type IBuilderContext interface {
	Config() *config.AppConfig
	Logger() *slog.Logger
}

// BuilderContext is the concrete implementation of IBuilderContext.
type BuilderContext struct {
	AppCfg *config.AppConfig
	L      *slog.Logger
}

func (bc *BuilderContext) Config() *config.AppConfig { return bc.AppCfg }
func (bc *BuilderContext) Logger() *slog.Logger      { return bc.L }

// NewBuilderContext creates a new context for report builders.
func NewBuilderContext(appConfig *config.AppConfig, logger *slog.Logger) IBuilderContext {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	return &BuilderContext{
		AppCfg: appConfig,
		L:      logger,
	}
}
