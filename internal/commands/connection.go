package commands

import (
	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handlePing(_ []string, _ *store.Store) string {
	return resp.EncodeString("PONG")
}

func handleEcho(args []string, _ *store.Store) string {
	return resp.EncodeBulk(args[0])
}
