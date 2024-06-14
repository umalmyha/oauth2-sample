package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/umalmyha/oauth2-sample/internal/config"
	"github.com/umalmyha/oauth2-sample/internal/handler"
	"github.com/umalmyha/oauth2-sample/internal/oauth2"
)

const shutdownTimeout = 30 * time.Second

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo},
		),
	)

	if err := start(logger); err != nil {
		logger.Error("error occurred on application routine", slog.Any("error", err))
		os.Exit(1)
	}
}

func start(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	h := handler.NewHandler(
		&handler.OAuthCredentials{
			ClientID:     cfg.Credentials.ClientID,
			ClientSecret: cfg.Credentials.ClientSecret,
		},
		oauth2.NewTokenIssuer(
			cfg.JWT.PrivateKey,
			oauth2.WithIssuer(cfg.JWT.Issuer),
			oauth2.WithTTL(cfg.JWT.TTL),
		),
		logger,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /token", h.Handle)

	addr := fmt.Sprintf(":%d", cfg.HTTPServer.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, os.Kill)

	go func() {
		logger.Info("starting HTTP server at " + addr)
		if srvErr := srv.ListenAndServe(); err != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			errCh <- srvErr
		}
	}()

	select {
	case err = <-errCh:
		return fmt.Errorf("error occurred while trying to start HTTP server: %w", err)
	case <-stopCh:
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err = srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}
	}

	return nil
}
