/*

	HPSCHD Main - v2

*/

package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func init() {
	// Init data locations
	// store ::: ephemeral mesostic cache
	localDirs([]string{"store"})

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
	// Start OpenTelemetry Trace Provider
	tp, err := NewTraceProviderOTEL()
	if err != nil {
		slog.Warn("Could not start trace provider, continuing...", slog.Any("error", err))
	}
	defer tp.Shutdown(context.Background())

	// Runtime Flags
	nofetch := flag.Bool("nofetch", false, "Do not start NASA APOD cronjob")
	port := flag.String("port", "9876", "Server port")
	flag.Parse()

	sp := &ServePoems{}

	// Start NASA APOD fetching in background (unless disabled)
	if !*nofetch {
		// Configure ticker interval for NASA APOD fetching
		t := envVar("HPSCHD_APOD_FREQUENCY", "88")
		ti, err := strconv.Atoi(t)
		if err != nil {
			slog.Error("unreadable frequency")
		}

		tid := time.Duration(ti) * time.Second
		ctx := context.Background()
		sp.Ticker = time.NewTicker(tid)
		defer sp.Ticker.Stop()

		go sp.TickerAPOD(ctx)
	}

	// Create API server with OTEL wrapper around Mux
	addr := ":" + *port
	sp.Server = &http.Server{
		Addr: addr,
		Handler: otelhttp.NewHandler(sp.SetupMux(), "MesosticAPIV2",
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			})),
	}

	// Even if APOD fetching is disabled, the homepage can still serve the contents of /store
	slog.Info("Starting Mesostic Generation Engine", slog.String("addr", addr))
	if err = sp.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
