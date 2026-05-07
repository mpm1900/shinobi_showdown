package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shinobi_showdown/internal/db"
	"shinobi_showdown/internal/server"
)

func main() {
	logger := slog.Default()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx)
	if err != nil {
		logger.Error("Error connecting to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()
	queries := db.New(pool)

	s := server.NewServer(ctx, queries)
	go run(s, logger)
	logger.Info("Server is running...")

	<-ctx.Done()
	shutdown(s, logger)
}

func run(s *server.Server, logger *slog.Logger) {
	err := s.Run()
	if err != nil {
		logger.Error("Error running server", "err", err)
		os.Exit(1)
	}
}

func shutdown(s *server.Server, logger *slog.Logger) {
	logger.Info("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down server", "err", err)
		os.Exit(1)
	}
}
