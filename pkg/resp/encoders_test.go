package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBulk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple string",
			input: "hello",
			want:  "$5\r\nhello\r\n",
		},
		{
			name:  "empty string",
			input: "",
			want:  "$0\r\n\r\n",
		},
		{
			name:  "string with spaces",
			input: "hello world",
			want:  "$11\r\nhello world\r\n",
		},
		{
			name:  "string with newlines",
			input: "line1\nline2",
			want:  "$11\r\nline1\nline2\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeBulk(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "PONG response",
			input: "PONG",
			want:  "+PONG\r\n",
		},
		{
			name:  "OK response",
			input: "OK",
			want:  "+OK\r\n",
		},
		{
			name:  "empty string",
			input: "",
			want:  "+\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeNil(t *testing.T) {
	got := EncodeNil()
	want := "$-1\r\n"

	assert.Equal(t, want, got)
}

func TestEncodeNilArray(t *testing.T) {
	got := EncodeNilArray()
	want := "*-1\r\n"

	assert.Equal(t, want, got)
}

func TestEncodeError(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple error",
			input: "something went wrong",
			want:  "-ERR something went wrong\r\n",
		},
		{
			name:  "empty error",
			input: "",
			want:  "-ERR \r\n",
		},
		{
			name:  "wrong number of arguments",
			input: "wrong number of arguments for 'get' command",
			want:  "-ERR wrong number of arguments for 'get' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeError(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeInteger(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  string
	}{
		{
			name:  "positive number",
			input: 42,
			want:  ":42\r\n",
		},
		{
			name:  "zero",
			input: 0,
			want:  ":0\r\n",
		},
		{
			name:  "negative number",
			input: -10,
			want:  ":-10\r\n",
		},
		{
			name:  "large number",
			input: 1000000,
			want:  ":1000000\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeInteger(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeArray(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "empty array",
			input: []string{},
			want:  "*0\r\n",
		},
		{
			name:  "single element",
			input: []string{"hello"},
			want:  "*1\r\n$5\r\nhello\r\n",
		},
		{
			name:  "two elements",
			input: []string{"list_key", "foo"},
			want:  "*2\r\n$8\r\nlist_key\r\n$3\r\nfoo\r\n",
		},
		{
			name:  "multiple elements",
			input: []string{"a", "b", "c"},
			want:  "*3\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n",
		},
		{
			name:  "elements with different lengths",
			input: []string{"short", "medium text", "x"},
			want:  "*3\r\n$5\r\nshort\r\n$11\r\nmedium text\r\n$1\r\nx\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeArray(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEncodeArray_BLPOP(t *testing.T) {
	input := []string{"mylist", "element1"}
	got := EncodeArray(input)
	want := "*2\r\n$6\r\nmylist\r\n$8\r\nelement1\r\n"

	assert.Equal(t, want, got)
}
