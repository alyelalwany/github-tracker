package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alyelalwany/github-tracker/cmd/app/handler"
	"github.com/alyelalwany/github-tracker/internal/echo"
	"github.com/alyelalwany/github-tracker/internal/health"
	"github.com/alyelalwany/github-tracker/internal/util"
)

func main() {
	cfg := util.LoadConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	slog.SetDefault(logger)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := handler.NewServer() // fork/exec gh once, build githubv4.Client

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Healthz)
	mux.HandleFunc("/readyz", health.Readyz)
	mux.HandleFunc("POST /echo", echo.Handler)
	mux.HandleFunc("GET /repos/{kind}/{login}", s.GetRepos)
	mux.HandleFunc("GET /user/{login}", s.GetUserDetails)

	srv := &http.Server{
		Addr:              cfg.BindAddr + ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
		slog.Info("Server is up at"+srv.Addr, "addr", srv.Addr)
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "err", err)
	}
	slog.Info("Shutdown..")

}
