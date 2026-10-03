package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/IgorBayerl/nanovision/internal/logging"
	"github.com/IgorBayerl/nanovision/internal/reporter/htmlreact"
	"github.com/IgorBayerl/nanovision/internal/server"
	"github.com/IgorBayerl/nanovision/internal/store"
)

// runServe is `nanovision serve`: the team server. CI uploads runs to it,
// and it shows them in the browser.
func runServe(args []string) int {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: nanovision serve [flags]")
		fmt.Fprintln(os.Stderr, "\nRuns the team server: CI uploads runs to it with the upload token, and")
		fmt.Fprintln(os.Stderr, "everyone opens them in the browser. Local runs use the static report.")
		fmt.Fprintln(os.Stderr)
		fs.PrintDefaults()
	}
	storeDir := fs.String("store", "", "Folder that holds the store of the server; it is made when missing")
	addr := fs.String("addr", "localhost:7070", "Address to listen on; ':7070' serves the whole network")
	tokenFile := fs.String("token-file", "", "File that holds the upload token; "+tokenEnv+" works too")
	publicURL := fs.String("public-url", "", "The URL people use to reach this server, for the links it hands out")
	reviewDays := fs.Int("keep-review-days", 30, "Days to keep review runs; 0 keeps them all")
	maxGB := fs.Int("max-store-gb", 20, "Size limit of the store in GB: above it the oldest runs are deleted; 0 means no limit")
	verbosity := fs.String("verbosity", "Info", "Logging level: Verbose, Info, Warning, Error, Off")
	fs.Parse(args)

	level, err := logging.ParseVerbosity(*verbosity)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if _, err := logging.Init(&logging.Config{Verbosity: level, Format: "text"}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	logger := slog.Default()

	if *storeDir == "" {
		logger.Error("serve needs -store: the folder that holds the store of the server")
		return 1
	}
	dir, err := filepath.Abs(*storeDir)
	if err != nil {
		logger.Error("Store folder", "error", err)
		return 1
	}
	token, err := uploadToken(*tokenFile)
	if err == nil && token == "" {
		err = errors.New("serve needs the upload token that CI sends: set " + tokenEnv + " or pass -token-file")
	}
	if err != nil {
		logger.Error("Upload token", "error", err)
		return 1
	}

	s, err := store.Open(dir)
	if err != nil {
		logger.Error("Could not open the store", "error", err)
		return 1
	}
	defer s.Close()

	ui, err := htmlreact.DistFS()
	if err != nil {
		logger.Error("Web UI is missing from this build", "error", err)
		return 1
	}
	// the build keys the browser cache of the views; a development build has no
	// commit, so each start gets its own
	build := version + "-" + commit
	if buildInfo().Dev() {
		build += "-" + strconv.FormatInt(time.Now().Unix(), 36)
	}
	srv, err := server.New(server.Options{
		Store: s, Token: token, Version: version, Build: build,
		PublicURL: *publicURL, UI: ui, Logger: logger,
	})
	if err != nil {
		logger.Error("Could not start the server", "error", err)
		return 1
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Error("Could not listen", "addr", *addr, "error", err)
		return 1
	}
	httpServer := &http.Server{Handler: srv, ReadHeaderTimeout: 10 * time.Second}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go srv.RunMaintenance(ctx, server.Upkeep{ReviewDays: *reviewDays, MaxBytes: int64(*maxGB) << 30})

	limit := "no size limit"
	if *maxGB > 0 {
		limit = fmt.Sprintf("limit %d GB", *maxGB)
	}
	fmt.Printf("\nnanovision serve  |  store %s (%s)\n\n  Open: %s\n\n", dir, limit, displayURL(ln.Addr(), *publicURL))

	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		httpServer.Shutdown(shutdown)
	}()
	if err := httpServer.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("Server stopped", "error", err)
		return 1
	}
	return 0
}

func uploadToken(file string) (string, error) {
	if file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}
	return strings.TrimSpace(os.Getenv(tokenEnv)), nil
}

func displayURL(addr net.Addr, public string) string {
	if public != "" {
		return public
	}
	tcp, ok := addr.(*net.TCPAddr)
	if !ok {
		return "http://" + addr.String()
	}
	host := tcp.IP.String()
	if tcp.IP.IsUnspecified() {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s", net.JoinHostPort(host, fmt.Sprint(tcp.Port)))
}
