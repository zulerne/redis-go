package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/zulerne/redis-go/internal/config"
	"github.com/zulerne/redis-go/internal/server"
	"github.com/zulerne/redis-go/internal/store"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.Default()

	redisStore := store.NewStore(logger.With("component", "store"))

	redisServer := server.NewServer(
		redisStore,
		config.DefaultConfig(),
		logger.With("component", "server"),
	)

	return redisServer.Listen(ctx)
}
