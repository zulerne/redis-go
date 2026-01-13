package store

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newStore() *Store {
	return NewStore(slog.Default())
}

func TestStore_SetGet(t *testing.T) {
	s := newStore()
	key := "foo"
	val := "bar"

	s.Set(key, val, 0)

	got, err := s.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, val, got)
}

func TestStore_Expiration(t *testing.T) {
	s := newStore()
	key := "temp"
	val := "data"
	ttl := 100 * time.Millisecond

	s.Set(key, val, ttl)

	_, err := s.Get(key)
	assert.Nil(t, err)

	time.Sleep(200 * time.Millisecond)

	_, err = s.Get(key)
	assert.NotNil(t, err)
}

func TestStore_Concurrency(t *testing.T) {
	s := newStore()

	for i := 0; i < 100; i++ {
		go func() {
			s.Set("key", "value", 10*time.Millisecond)
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, _ = s.Get("key")
		}()
	}

	time.Sleep(500 * time.Millisecond)
}

func TestStore_Append(t *testing.T) {
	s := newStore()
	key := "mylist"

	len1, _ := s.Append(key, "val1")
	assert.Equal(t, 1, len1)

	len2, _ := s.Append(key, "val2")
	assert.Equal(t, 2, len2)

	len3, _ := s.Append(key, "val3", "val4")
	assert.Equal(t, 4, len3)
}

func TestStore_GetRange(t *testing.T) {
	s := newStore()
	key := "letters"
	// Setup: List ["a", "b", "c", "d", "e"]
	_, _ = s.Append(key, "a", "b", "c", "d", "e")

	tests := []struct {
		name     string
		start    int
		stop     int
		expected int // expected length of result
	}{
		{"Full range", 0, 4, 5},
		{"Partial range", 0, 2, 3},    // a, b, c
		{"Middle range", 1, 3, 3},     // b, c, d
		{"Oversized stop", 0, 100, 5}, // Should return all 5
		{"Start > Stop", 2, 1, 0},     // Empty
		{"Start out of bounds", 10, 12, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := s.GetRange(key, tt.start, tt.stop)
			assert.Equal(t, tt.expected, len(got))
		})
	}
}
func TestStore_GetRange_Negative(t *testing.T) {
	s := newStore()
	key := "nums"
	// List: ["1", "2", "3", "4", "5"]
	_, _ = s.Append(key, "1", "2", "3", "4", "5")

	tests := []struct {
		name     string
		start    int
		stop     int
		expected int
		firstVal string
	}{
		{"Last element", -1, -1, 1, "5"},
		{"Last two", -2, -1, 2, "4"},
		{"Negative range overshoot", -100, -1, 5, "1"},
		{"Mixed indices", 0, -1, 5, "1"},              // 0 .. last
		{"Mixed overshoot", -100, 2, 3, "1"},          // 0 .. 2
		{"Stop before start negative", -4, -5, 0, ""}, // start(1) > stop(0)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := s.GetRange(key, tt.start, tt.stop)
			assert.Equal(t, tt.expected, len(got))
			if tt.expected > 0 {
				assert.Equal(t, tt.firstVal, got[0])
			}
		})
	}
}

func TestStore_LeftAppend(t *testing.T) {
	s := newStore()
	key := "l_list"

	// LPUSH key "a" "b" "c"
	len1, _ := s.LeftAppend(key, "a", "b", "c")
	assert.Equal(t, 3, len1)

	rangeRes, _ := s.GetRange(key, 0, -1)
	expectedOrder := []string{"c", "b", "a"}

	for i, val := range rangeRes {
		assert.Equal(t, expectedOrder[i], val)
	}
}

func TestStore_ListLen(t *testing.T) {
	s := newStore()
	key := "my_len_list"

	assert.Equal(t, 0, s.ListLen(key))

	_, _ = s.Append(key, "a", "b", "c")

	assert.Equal(t, 3, s.ListLen(key))
}

func TestStore_RPop(t *testing.T) {
	s := newStore()
	key := "list"

	_, _ = s.Append(key, "one", "two", "three")

	val, err := s.RPop(key, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	assert.Equal(t, "three", val[0])

	assert.Equal(t, 2, s.ListLen(key))

	val, _ = s.RPop(key, 1)
	assert.Equal(t, "two", val[0])
}

func TestStore_RPopMultiple(t *testing.T) {
	s := newStore()
	key := "multilist"
	_, _ = s.Append(key, "1", "2", "3", "4", "5")

	res, err := s.RPop(key, 2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "4", res[0])
	assert.Equal(t, "5", res[1])

	assert.Equal(t, 3, s.ListLen(key))

	res, _ = s.RPop(key, 10)
	assert.Equal(t, 3, len(res))

	assert.Equal(t, 0, s.ListLen(key))
}

func TestStore_LPop(t *testing.T) {
	s := newStore()
	key := "queue"

	_, _ = s.Append(key, "one", "two", "three")

	val, err := s.LPop(key, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	assert.Equal(t, "one", val[0])

	assert.Equal(t, 2, s.ListLen(key))

	val, _ = s.LPop(key, 1)
	assert.Equal(t, "two", val[0])
}

func TestStore_LPopMultiple(t *testing.T) {
	s := newStore()
	key := "multilist"
	_, _ = s.Append(key, "1", "2", "3", "4", "5")

	res, err := s.LPop(key, 2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "1", res[0])
	assert.Equal(t, "2", res[1])

	assert.Equal(t, 3, s.ListLen(key))

	res, _ = s.LPop(key, 10)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, 0, s.ListLen(key))
}

func TestStore_BlpopTimeout(t *testing.T) {
	s := newStore()
	listName := "test_list"

	start := time.Now()
	res, err := s.BLPop(listName, 0.1)
	elapsed := time.Since(start)

	assert.True(t, errors.Is(err, ErrTimeout))
	assert.Nil(t, res)
	assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
}

func TestStore_BlpopSuccess(t *testing.T) {
	s := newStore()
	listName := "test_list"

	go func() {
		time.Sleep(50 * time.Millisecond)
		_, _ = s.Append(listName, "value")
	}()

	res, err := s.BLPop(listName, 1.0)

	assert.Nil(t, err)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, res[0], listName)
	assert.Equal(t, res[1], "value")
}

func TestStore_Type(t *testing.T) {
	s := newStore()
	key := "foo"
	s.Set(key, "bar", 0)

	listKey := "list"
	_, _ = s.Append(listKey, "a", "b", "c")

	streamKey := "stream"
	_, _ = s.XAdd(streamKey, "1-1", []string{"field", "value"})

	assert.Equal(t, "string", s.Type(key))
	assert.Equal(t, "list", s.Type(listKey))
	assert.Equal(t, "stream", s.Type(streamKey))
	assert.Equal(t, "none", s.Type("nonexistent"))
}

func TestStore_XAdd_ValidIds(t *testing.T) {
	tests := []struct {
		name    string
		entries []string
	}{
		{
			name:    "simple sequential",
			entries: []string{"1-1", "1-2", "1-3"},
		},
		{
			name:    "increasing ms resets seq",
			entries: []string{"1-10", "2-1"},
		},
		{
			name:    "large jump in ms",
			entries: []string{"100-1", "200-1", "300-1"},
		},
		{
			name:    "first entry can be 0-1",
			entries: []string{"0-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore()

			for _, id := range tt.entries {
				result, err := s.XAdd("stream", id, []string{"k", "v"})
				assert.Nil(t, err, "ID %s should be valid", id)
				assert.Equal(t, id, result)
			}
		})
	}
}

func TestStore_XAdd_InvalidIds(t *testing.T) {
	tests := []struct {
		name        string
		setupId     string
		invalidId   string
		expectedErr error
	}{
		{
			name:        "0-0 is always invalid",
			setupId:     "",
			invalidId:   "0-0",
			expectedErr: ErrZeroEntryId,
		},
		{
			name:        "same ID as last",
			setupId:     "1-1",
			invalidId:   "1-1",
			expectedErr: ErrSmallerEntryId,
		},
		{
			name:        "smaller seq with same ms",
			setupId:     "1-5",
			invalidId:   "1-3",
			expectedErr: ErrSmallerEntryId,
		},
		{
			name:        "smaller ms even with larger seq",
			setupId:     "5-1",
			invalidId:   "3-100",
			expectedErr: ErrSmallerEntryId,
		},
		{
			name:        "0-0 after valid entry",
			setupId:     "1-1",
			invalidId:   "0-0",
			expectedErr: ErrZeroEntryId,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newStore()

			if tt.setupId != "" {
				_, err := s.XAdd("stream", tt.setupId, []string{"k", "v"})
				assert.Nil(t, err)
			}

			_, err := s.XAdd("stream", tt.invalidId, []string{"k", "v"})
			assert.True(t, errors.Is(err, tt.expectedErr),
				"expected %v, got %v", tt.expectedErr, err)
		})
	}
}
