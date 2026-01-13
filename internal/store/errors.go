package store

import "errors"

var (
	ErrWrongType             = errors.New("wrong type")
	ErrKeyNotFound           = errors.New("key not found")
	ErrListNotExist          = errors.New("list doesn't exist")
	ErrListEmpty             = errors.New("list is empty")
	ErrTimeout               = errors.New("timeout")
	ErrInternalTypeAssertion = errors.New("internal: unexpected type assertion failure")
	ErrInvalidEntryId        = errors.New("invalid entry id")
	ErrSmallerEntryId        = errors.New("The ID specified in XADD is equal or smaller than the target stream top item")
	ErrZeroEntryId           = errors.New("The ID specified in XADD must be greater than 0-0")
)
