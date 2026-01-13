package commands

import (
	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handleType(args []string, s *store.Store) string {
	key := args[0]
	res := s.Type(key)
	return resp.EncodeString(res)
}
