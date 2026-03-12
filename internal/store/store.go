package store

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DataType string

const (
	TypeString DataType = "string"
	TypeList   DataType = "list"
	TypeStream DataType = "stream"
	TypeNone   DataType = "none"
)

type Item struct {
	Value     any
	Type      DataType
	ExpiresAt time.Time
}
type StreamEntry struct {
	ms     int
	seq    int
	Values map[string]string
}

func (se StreamEntry) Id() string {
	return fmt.Sprintf("%d-%d", se.ms, se.seq)
}

type Store struct {
	data           map[string]*Item
	logger         *slog.Logger
	blockingClient map[string][]chan string
	mu             sync.RWMutex
}

func NewStore(logger *slog.Logger) *Store {
	return &Store{
		data:           make(map[string]*Item),
		logger:         logger,
		blockingClient: make(map[string][]chan string),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	s.data[key] = &Item{Type: TypeString, Value: value, ExpiresAt: expiresAt}
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		delete(s.data, key)
		return "", ErrKeyNotFound
	}

	if item.Type != TypeString {
		return "", ErrWrongType
	}

	val, ok := item.Value.(string)
	if !ok {
		return "", ErrInternalTypeAssertion
	}

	return val, nil
}

func (s *Store) Append(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]

	if !ok {
		item = &Item{
			Value: []string{},
			Type:  TypeList,
		}
		s.data[key] = item
	}

	if item.Type != TypeList {
		return 0, ErrWrongType
	}

	values = s.notifyWaiters(key, values)

	list, ok := item.Value.([]string)
	if !ok {
		return 0, ErrInternalTypeAssertion
	}
	if len(values) > 0 {
		list = append(list, values...)
	}
	item.Value = list

	return len(list), nil
}

func (s *Store) LeftAppend(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	slices.Reverse(values)

	item, ok := s.data[key]

	if !ok {
		item = &Item{
			Value: []string{},
			Type:  TypeList,
		}
		s.data[key] = item
	}

	if item.Type != TypeList {
		return 0, ErrWrongType
	}

	values = s.notifyWaiters(key, values)

	list, ok := item.Value.([]string)
	if !ok {
		return 0, ErrInternalTypeAssertion
	}
	if len(values) > 0 {
		list = append(values, list...)
		item.Value = list
	}

	return len(list), nil
}

func (s *Store) GetRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.data[key]

	if !ok {
		return []string{}, ErrListNotExist
	}

	if item.Type != TypeList {
		return []string{}, ErrWrongType
	}

	list, ok := item.Value.([]string)
	if !ok {
		return []string{}, ErrInternalTypeAssertion
	}

	length := len(list)

	if start < 0 {
		start = max(length+start, 0)
	}
	if stop < 0 {
		stop = max(length+stop, 0)
	}

	if start > stop {
		return []string{}, nil
	}
	if start >= length {
		return []string{}, nil
	}

	if stop >= length {
		stop = length - 1
	}

	res := make([]string, stop+1-start)
	copy(res, list[start:stop+1])

	return res, nil
}

func (s *Store) ListLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.data[key]

	if !ok {
		return 0
	}

	if item.Type != TypeList {
		return 0
	}

	list, ok := item.Value.([]string)
	if !ok {
		return 0
	}

	return len(list)
}

func (s *Store) RPop(key string, count int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]

	if !ok {
		return []string{}, ErrListNotExist
	}

	if item.Type != TypeList {
		return []string{}, ErrWrongType
	}

	list, ok := item.Value.([]string)
	if !ok {
		return []string{}, ErrInternalTypeAssertion
	}

	length := len(list)
	if length == 0 {
		return []string{}, ErrListEmpty
	}
	if count > length {
		count = length
	}

	r := make([]string, count)
	copy(r, list[length-count:])

	list = list[:length-count]
	item.Value = list

	if len(list) == 0 {
		delete(s.data, key)
	}

	return r, nil
}

func (s *Store) LPop(key string, count int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]

	if !ok {
		return []string{}, ErrListNotExist
	}

	if item.Type != TypeList {
		return []string{}, ErrWrongType
	}

	list, ok := item.Value.([]string)
	if !ok {
		return []string{}, ErrInternalTypeAssertion
	}

	length := len(list)
	if length == 0 {
		return []string{}, ErrListEmpty
	}
	if count > length {
		count = length
	}

	r := make([]string, count)
	copy(r, list[:count])

	list = list[count:]
	item.Value = list

	if len(list) == 0 {
		delete(s.data, key)
	}

	return r, nil
}

func (s *Store) BLPop(ctx context.Context, key string, timeout float64) ([]string, error) {
	s.mu.Lock()

	item, ok := s.data[key]
	if ok {
		if item.Type != TypeList {
			s.mu.Unlock()
			return nil, ErrWrongType
		}

		list, ok := item.Value.([]string)
		if !ok {
			s.mu.Unlock()
			return nil, ErrInternalTypeAssertion
		}
		if len(list) > 0 {
			val := list[0]
			list = list[1:]
			item.Value = list

			if len(list) == 0 {
				delete(s.data, key)
			}
			s.mu.Unlock()
			return []string{key, val}, nil
		}
	}

	ch := make(chan string, 1)
	s.blockingClient[key] = append(s.blockingClient[key], ch)

	s.mu.Unlock()

	var timer <-chan time.Time
	if timeout > 0 {
		timer = time.After(time.Duration(timeout * float64(time.Second)))
	}

	select {
	case v := <-ch:
		return []string{key, v}, nil
	case <-timer:
		s.cleanupBlockingClients(key, ch)
		return nil, ErrTimeout
	case <-ctx.Done():
		s.cleanupBlockingClients(key, ch)
		return nil, ctx.Err()
	}
}

func (s *Store) Type(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, ok := s.data[key]; ok {
		if !val.ExpiresAt.IsZero() && time.Now().After(val.ExpiresAt) {
			delete(s.data, key)
			return string(TypeNone)
		}
		return string(val.Type)
	}

	return string(TypeNone)
}

func (s *Store) XAdd(key, entryId string, keyValues []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.data[key]

	if !ok {
		item = &Item{
			Value: []StreamEntry{},
			Type:  TypeStream,
		}
		s.data[key] = item
	}

	if item.Type != TypeStream {
		return "", ErrWrongType
	}

	stream, ok := item.Value.([]StreamEntry)
	if !ok {
		return "", ErrInternalTypeAssertion
	}

	var last StreamEntry
	if len(stream) > 0 {
		last = stream[len(stream)-1]
	}
	ms, seq, err := s.parseEntryId(entryId, last)
	if err != nil {
		return "", err
	}

	values := make(map[string]string, len(keyValues)/2)
	for i := 0; i+1 < len(keyValues); i += 2 {
		values[keyValues[i]] = keyValues[i+1]
	}

	entry := StreamEntry{
		ms:     ms,
		seq:    seq,
		Values: values,
	}

	stream = append(stream, entry)
	item.Value = stream

	return entry.Id(), nil
}

func (s *Store) notifyWaiters(listName string, values []string) []string {
	for len(values) > 0 {
		waiters, ok := s.blockingClient[listName]
		if !ok || len(waiters) == 0 {
			break
		}

		ch := waiters[0]
		s.blockingClient[listName] = waiters[1:]
		if len(s.blockingClient[listName]) == 0 {
			delete(s.blockingClient, listName)
		}

		select {
		case ch <- values[0]:
			values = values[1:]
		default:
		}
	}
	return values
}

func (s *Store) cleanupBlockingClients(listName string, ch chan string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cl, ok := s.blockingClient[listName]
	if !ok {
		return
	}

	for i, channel := range s.blockingClient[listName] {
		if ch == channel {
			s.blockingClient[listName] = append(cl[:i], cl[i+1:]...)
			break
		}
	}
}

func (s *Store) parseEntryId(entryId string, prev StreamEntry) (ms, seq int, err error) {
	msStr, seqStr, ok := strings.Cut(entryId, "-")
	if !ok {
		return 0, 0, ErrInvalidEntryId
	}
	ms, err = strconv.Atoi(msStr)
	if err != nil {
		return 0, 0, ErrInvalidEntryId
	}
	seq, err = strconv.Atoi(seqStr)
	if err != nil {
		return 0, 0, ErrInvalidEntryId
	}

	if ms == 0 && seq == 0 {
		return 0, 0, ErrZeroEntryId
	}
	if ms < prev.ms {
		return 0, 0, ErrSmallerEntryId
	}
	if ms == prev.ms && seq <= prev.seq {
		return 0, 0, ErrSmallerEntryId
	}

	return ms, seq, nil
}
