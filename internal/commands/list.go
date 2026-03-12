package commands

import (
	"context"
	"errors"
	"strconv"

	"github.com/zulerne/redis-go/internal/store"
	"github.com/zulerne/redis-go/pkg/resp"
)

func handleRPush(_ context.Context, args []string, s *store.Store) string {
	l, err := s.Append(args[0], args[1:]...)

	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(l)
}

func handleLPush(_ context.Context, args []string, s *store.Store) string {
	l, err := s.LeftAppend(args[0], args[1:]...)
	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeInteger(l)
}

func handleLRange(_ context.Context, args []string, s *store.Store) string {
	listName := args[0]
	start, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(LRange.WrongOptionsMessage())
	}

	stop, err := strconv.Atoi(args[2])
	if err != nil {
		return resp.EncodeError(LRange.WrongOptionsMessage())
	}

	list, err := s.GetRange(listName, start, stop)

	if err != nil {
		return resp.EncodeError(err.Error())
	}

	return resp.EncodeArray(list)
}

func handleLLen(_ context.Context, args []string, s *store.Store) string {
	listName := args[0]

	l := s.ListLen(listName)

	return resp.EncodeInteger(l)
}

func handleRPop(_ context.Context, args []string, s *store.Store) string {
	listName := args[0]

	if len(args) == 1 {
		v, err := s.RPop(listName, 1)
		if err != nil {
			return resp.EncodeNil()
		}

		return resp.EncodeBulk(v[0])
	}

	count, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(RPop.WrongOptionsMessage())
	}
	res, err := s.RPop(listName, count)
	if err != nil {
		return resp.EncodeNil()
	}

	return resp.EncodeArray(res)
}

func handleLPop(_ context.Context, args []string, s *store.Store) string {
	listName := args[0]

	if len(args) == 1 {
		res, err := s.LPop(listName, 1)
		if err != nil {
			return resp.EncodeNil()
		}
		return resp.EncodeBulk(res[0])
	}

	count, err := strconv.Atoi(args[1])
	if err != nil {
		return resp.EncodeError(LPop.WrongOptionsMessage())
	}
	res, err := s.LPop(listName, count)
	if err != nil {
		return resp.EncodeNil()
	}

	return resp.EncodeArray(res)
}

func handleBLPop(ctx context.Context, args []string, s *store.Store) string {
	listName := args[0]

	timeout, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return resp.EncodeError(BLPop.WrongOptionsMessage())
	}

	res, err := s.BLPop(ctx, listName, timeout)
	if err != nil {
		if errors.Is(err, store.ErrTimeout) {
			return resp.EncodeNilArray()
		}
		return resp.EncodeNil()
	}
	return resp.EncodeArray(res)
}
