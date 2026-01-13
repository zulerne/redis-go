package commands

import (
	"errors"
	"strconv"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handleRPush(args []string, store *store.Store) string {
	l, err := store.Append(args[0], args[1:]...)

	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(l)
}

func handleLPush(args []string, store *store.Store) string {
	l, err := store.LeftAppend(args[0], args[1:]...)
	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(l)
}

func handleLRange(args []string, store *store.Store) string {
	listName := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(LRange.WrongOptionsMessage())
	}

	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return resp.EncodeError(LRange.WrongOptionsMessage())
	}

	list, err := store.GetRange(listName, start, stop)

	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeArray(list)
}

func handleLLen(args []string, store *store.Store) string {
	listName := args[0]

	l := store.ListLen(listName)

	return resp.EncodeInteger(l)
}

func handleRPop(args []string, store *store.Store) string {
	listName := args[0]

	if len(args) == 1 {
		v, err := store.RPop(listName, 1)
		if err != nil {
			return resp.EncodeNil()
		}

		return resp.EncodeBulk(v[0])
	}

	count, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(LPop.WrongOptionsMessage())
	}
	res, err := store.RPop(listName, count)
	if err != nil {
		return resp.EncodeNil()
	}

	return resp.EncodeArray(res)
}

func handleLPop(args []string, store *store.Store) string {
	listName := args[0]

	if len(args) == 1 {
		res, err := store.LPop(listName, 1)
		if err != nil {
			return resp.EncodeNil()
		}
		return resp.EncodeBulk(res[0])
	}

	count, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(LPop.WrongOptionsMessage())
	}
	res, err := store.LPop(listName, count)
	if err != nil {
		return resp.EncodeNil()
	}

	return resp.EncodeArray(res)
}

func handleBLPop(args []string, s *store.Store) string {
	listName := args[0]

	timeout, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return resp.EncodeError(BLPop.WrongOptionsMessage())
	}

	res, err := s.BLPop(listName, timeout)
	if err != nil {
		if errors.Is(err, store.ErrTimeout) {
			return resp.EncodeNilArray()
		}
		return resp.EncodeNil()
	}
	return resp.EncodeArray(res)
}
