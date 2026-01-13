package commands

import (
	"fmt"
	"strings"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

type Command string

const (
	// Connection

	Ping Command = "PING"
	Echo Command = "ECHO"

	// String

	Set Command = "SET"
	Get Command = "GET"

	// List

	RPush  Command = "RPUSH"
	LPush  Command = "LPUSH"
	LRange Command = "LRANGE"
	LLen   Command = "LLEN"
	RPop   Command = "RPOP"
	LPop   Command = "LPOP"
	BLPop  Command = "BLPOP"

	// Generic

	Type Command = "TYPE"

	// Stream

	XAdd Command = "XADD"
)

func (cn Command) WrongArgumentsMessage() string {
	return fmt.Sprintf("wrong number of arguments for '%s' command", strings.ToLower(string(cn)))
}

func (cn Command) WrongOptionsMessage() string {
	return fmt.Sprintf("wrong options for '%s' command", strings.ToLower(string(cn)))
}

type HandlerFunc func(args []string, s *store.Store) string

type CommandHandler struct {
	Action  HandlerFunc
	MinArgs int
}

var handlers = map[Command]CommandHandler{
	Ping: {Action: handlePing, MinArgs: 0},
	Echo: {Action: handleEcho, MinArgs: 1},

	Set: {Action: handleSet, MinArgs: 2},
	Get: {Action: handleGet, MinArgs: 1},

	RPush:  {Action: handleRPush, MinArgs: 2},
	LPush:  {Action: handleLPush, MinArgs: 2},
	LRange: {Action: handleLRange, MinArgs: 3},
	LLen:   {Action: handleLLen, MinArgs: 1},
	RPop:   {Action: handleRPop, MinArgs: 1},
	LPop:   {Action: handleLPop, MinArgs: 1},
	BLPop:  {Action: handleBLPop, MinArgs: 2},

	Type: {Action: handleType, MinArgs: 1},

	XAdd: {Action: handleXAdd, MinArgs: 4},
}

func Handle(cmdName string, args []string, store *store.Store) string {
	cmd := Command(cmdName)

	handler, ok := handlers[cmd]

	if !ok {
		return resp.EncodeError(fmt.Sprintf("unknown command '%s'", cmd))
	}
	if len(args) < handler.MinArgs {
		return resp.EncodeError(cmd.WrongArgumentsMessage())
	}

	return handler.Action(args, store)
}
