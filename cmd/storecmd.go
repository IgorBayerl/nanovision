package main

import (
	"path/filepath"
	"strings"

	"context"
	"flag"
	"fmt"
	"go.yaml.in/yaml/v3"
	"os"

	"github.com/IgorBayerl/nanovision/internal/config"
	"github.com/IgorBayerl/nanovision/internal/store"
)

const storeUsage = `Usage: nanovision store <command> [flags] [file]

Commands:
  gc             delete the blobs no run uses
  backup <file>  write a consistent copy of the store; safe while the server runs

Flags:
  -store <dir>   store folder (default: history.store of nanovision.yaml, else .nanovision)
  -config <file> nanovision.yaml to read history.store from
`

// runStoreCommand is `nanovision store ...`: the upkeep of a store.
func runStoreCommand(args []string) int {
	if len(args) == 0 || args[0] == "-h" || args[0] == "-help" || args[0] == "--help" {
		fmt.Fprint(os.Stderr, storeUsage)
		return 2
	}
	cmd := args[0]
	fs := flag.NewFlagSet("store "+cmd, flag.ExitOnError)
	fs.Usage = func() { fmt.Fprint(os.Stderr, storeUsage) }
	storeDir := fs.String("store", "", "")
	configPath := fs.String("config", "", "")
	fs.Parse(args[1:])

	dir, err := resolveStoreDir(*storeDir, *configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Fprintf(os.Stderr, "no store at %s\n", dir)
		return 1
	}
	s, err := store.Open(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer s.Close()
	ctx := context.Background()

	switch cmd {
	case "gc":
		res, err := s.GC(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Printf("Deleted %d unused blobs (%d bytes).\n", res.Blobs, res.Bytes)
	case "backup":
		if fs.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "usage: nanovision store backup <file>")
			return 2
		}
		if err := s.Backup(ctx, fs.Arg(0)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Printf("Backup written to %s\n", fs.Arg(0))
	default:
		fmt.Fprintf(os.Stderr, "unknown store command %q\n\n%s", cmd, storeUsage)
		return 2
	}
	return 0
}

// resolveStoreDir picks the store folder: the flag, then history.store of the
// config, then .nanovision next to it.
func resolveStoreDir(flagValue, configPath string) (string, error) {
	if flagValue != "" {
		return filepath.Abs(flagValue)
	}
	if configPath == "" {
		// a missing file is fine: the store then has its default folder
		if configPath, _ = config.FindFile("."); configPath == "" {
			configPath = config.FileNames[0]
		}
	}
	base := filepath.Dir(configPath)
	if data, err := os.ReadFile(configPath); err == nil {
		var cfg struct {
			History struct {
				Store string `yaml:"store"`
			} `yaml:"history"`
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return "", fmt.Errorf("%s: %w", configPath, err)
		}
		if st := cfg.History.Store; st != "" {
			lower := strings.ToLower(st)
			if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
				return "", fmt.Errorf("history.store in %s is the team server %s; open it in the browser, or pass -store with a local folder", configPath, st)
			}
			if !filepath.IsAbs(st) {
				st = filepath.Join(base, st)
			}
			return filepath.Abs(st)
		}
	}
	return filepath.Abs(filepath.Join(base, ".nanovision"))
}
