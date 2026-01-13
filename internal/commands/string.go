package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

// SET command options
const (
	optEX = "EX"
	optPX = "PX"
)

func handleSet(args []string, s *store.Store) string {
	var ttl time.Duration
	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			opt := strings.ToUpper(args[i])
			switch opt {
			case optEX, optPX:
				if i+1 >= len(args) {
					return resp.EncodeError(Set.WrongOptionsMessage())
				}
				amount, err := strconv.Atoi(args[i+1])
				if err != nil {
					return resp.EncodeError(Set.WrongOptionsMessage())
				}
				if opt == optPX {
					ttl = time.Duration(amount) * time.Millisecond
				} else {
					ttl = time.Duration(amount) * time.Second
				}
				i++
			}
		}
	}
	s.Set(args[0], args[1], ttl)

	return resp.EncodeString(resp.OkMessage)
}

func handleGet(args []string, s *store.Store) string {
	val, err := s.Get(args[0])

	if err != nil {
		return resp.EncodeNil()
	}

	return resp.EncodeBulk(val)
}
