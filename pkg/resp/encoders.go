package resp

import (
	"fmt"
	"strings"
)

// RespByte represents RESP protocol type markers.
type RespByte = byte

const (
	StringChar  RespByte = '+' // Simple string
	ErrorChar   RespByte = '-' // Error
	IntegerChar RespByte = ':' // Integer
	BulkChar    RespByte = '$' // Bulk string
	ArrayChar   RespByte = '*' // Array
)

const (
	OkMessage  = "OK"
	ErrMessage = "ERR"
)

func EncodeBulk(s string) string {
	return fmt.Sprintf("%c%d\r\n%s\r\n", BulkChar, len(s), s)
}

func EncodeString(s string) string {
	return fmt.Sprintf("%c%s\r\n", StringChar, s)
}

func EncodeNil() string {
	return fmt.Sprintf("%c-1\r\n", BulkChar)
}

func EncodeError(e string) string {
	return fmt.Sprintf("%c%s %s\r\n", ErrorChar, ErrMessage, e)
}

func EncodeInteger(i int) string {
	return fmt.Sprintf("%c%d\r\n", IntegerChar, i)
}

func EncodeArray(a []string) string {
	var r strings.Builder
	r.WriteString(fmt.Sprintf("%c%d\r\n", ArrayChar, len(a)))

	for _, s := range a {
		r.WriteString(EncodeBulk(s))
	}

	return r.String()
}

func EncodeNilArray() string {
	return fmt.Sprintf("%c-1\r\n", ArrayChar)
}
