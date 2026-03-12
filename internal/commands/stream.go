package commands

import (
	"context"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handleXAdd(_ context.Context, args []string, s *store.Store) string {
	streamKey := args[0]
	entry := args[1]

	keyValues := args[2:]
	if len(keyValues)%2 != 0 {
		keyValues = keyValues[:len(keyValues)-1]
	}

	res, err := s.XAdd(streamKey, entry, keyValues)
	if err != nil {
		return resp.EncodeError(err.Error())
	}
	return resp.EncodeBulk(res)
}
