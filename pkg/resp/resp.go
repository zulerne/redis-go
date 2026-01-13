package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Type string

const (
	ArrayType Type = "array"
	BulkType  Type = "bulk"
)

type Value struct {
	Type  Type
	Value string
	Array []Value
}
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) readLine() ([]byte, error) {
	var line []byte
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, err
		}

		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], nil
}

func (r *Resp) readInteger() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i64), nil
}

func (r *Resp) Read() (Value, error) {
	dataType, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch dataType {
	case ArrayChar:
		return r.readArray()
	case BulkChar:
		return r.readBulk()
	default:
		return Value{}, fmt.Errorf("unknown type %q: %w", string(dataType), ErrUnknownProtocol)
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.Type = ArrayType

	l, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.Array = make([]Value, 0, l)
	for i := 0; i < l; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.Type = BulkType

	l, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, l)
	_, err = io.ReadFull(r.reader, bulk)
	if err != nil {
		return v, err
	}
	v.Value = string(bulk)

	// Read the trailing CRLF
	if _, err = r.readLine(); err != nil {
		return v, err
	}

	return v, nil
}
