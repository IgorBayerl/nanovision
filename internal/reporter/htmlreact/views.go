package htmlreact

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/model"
)

// ViewOptions shape the JSON a server builds for a stored run. The shapes are
// the ones the static report embeds, so the UI renders both the same way.
type ViewOptions struct {
	// FileURL links a file row to its details page
	FileURL func(path string) string
	// replaces the information block
	Metadata []MetadataItem
	// when the run was made; a stored run's views never change
	GeneratedAt time.Time
}

func viewBuilder(cfg *config.AppConfig, opts ViewOptions) *HtmlReactReportBuilder {
	return &HtmlReactReportBuilder{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		config: cfg,
		view:   &opts,
	}
}

// SummaryJSON is the summaryV1 document of an annotated tree.
func SummaryJSON(tree *model.SummaryTree, cfg *config.AppConfig, opts ViewOptions) ([]byte, error) {
	summary, err := viewBuilder(cfg, opts).transformTree(tree)
	if err != nil {
		return nil, err
	}
	return encodeView(summary)
}

// DetailsJSON is the detailsV1 document of one file of an annotated tree.
// source holds the file's lines; nil gives empty lines.
func DetailsJSON(tree *model.SummaryTree, file *model.FileNode, source []string, cfg *config.AppConfig, opts ViewOptions) ([]byte, error) {
	details, err := viewBuilder(cfg, opts).detailsWithSource(tree, file, source)
	if err != nil {
		return nil, err
	}
	return encodeView(details)
}

func encodeView(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}
