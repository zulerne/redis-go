package commands

import (
	"context"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handlePing(_ context.Context, _ []string, _ *store.Store) string {
	return resp.EncodeString("PONG")
}

func handleEcho(_ context.Context, args []string, _ *store.Store) string {
	return resp.EncodeBulk(args[0])
}
