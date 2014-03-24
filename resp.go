package resp

import (
	"fmt"
)

const (
	// A large INFO ALL response can be over 4kb, so we set the default to 8kb.
	DEFAULT_BUFFER = 8192

	// Smallest valid RESP object is ":0\r\n".
	MIN_OBJECT_LENGTH = 4

	// The minimum valid command object is "*1\r\n$4\r\nPING\r\n"
	MIN_COMMAND_LENGTH = 14

	// RESP type prefixes
	SIMPLE_STRING_PREFIX = '+'
	ERROR_PREFIX         = '-'
	INTEGER_PREFIX       = ':'
	BULK_STRING_PREFIX   = '$'
	ARRAY_PREFIX         = '*'
)

var (
	// Not really a constant, but...
	LINE_ENDING = []byte("\r\n")

	// Errors
	ErrSyntaxError = fmt.Errorf("resp: syntax error")
	ErrBufferFull  = fmt.Errorf("resp: object is larger than buffer")
)

// RESP points to the bytes of a RESP object.
type RESP []byte

func NewRESP(bytes []byte) (interface{}, error) {
	if len(bytes) == 0 {
		return nil, ErrSyntaxError
	}

	switch bytes[0] {
	case SIMPLE_STRING_PREFIX:
		return NewSimpleString(bytes)
	case ERROR_PREFIX:
		return NewError(bytes)
	case INTEGER_PREFIX:
		return NewInteger(bytes)
	case BULK_STRING_PREFIX:
		return NewBulkString(bytes)
	case ARRAY_PREFIX:
		return NewArray(bytes)
	default:
		return nil, ErrSyntaxError
	}
}
