/*

	HPSCHD Main - v2

*/

package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	// Set up slog with JSON handler for structured logging
	// Default to Info level, can be overridden with HPSCHD_LOG_LEVEL env var
	logLevel := slog.LevelInfo
	if level := os.Getenv("HPSCHD_LOG_LEVEL"); level != "" {
		switch level {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		}
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// Runtime Flags
	nofetch := flag.Bool("nofetch", false, "Do not start NASA APOD cronjob")
	port := flag.String("port", "9876", "Server port")

	// Parse Flags
	flag.Parse()

	// Init data locations
	// store ::: ephemeral mesostic cache
	localDirs([]string{"store"})

	// Configure ticker interval for NASA APOD fetching
	t := envVar("HPSCHD_APOD_FREQUENCY", "88")
	ti, err := strconv.Atoi(t)
	if err != nil {
		slog.Error("unreadable frequency")
	}
	tid := time.Duration(ti) * time.Second
	if *nofetch {
		slog.Info("Running with integrated NASA APOD fetch disabled")
		tid = 0 // Disable ticker
	}

	// Initialize v2 server with ticker
	sp := &ServePoems{Ticker: time.NewTicker(tid)}
	defer sp.Ticker.Stop()

	// Start NASA APOD fetching in background (unless disabled)
	if !*nofetch {
		go sp.TickerAPOD()
	}

	// Start v2 API server (blocking)
	addr := ":" + *port
	sp.Server = &http.Server{
		Addr:    addr,
		Handler: sp.SetupMux(),
	}

	slog.Info("Starting v2 server", slog.String("addr", addr))
	if err := sp.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
