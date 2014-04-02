package resp

import (
	"errors"
)

const (
	// A large INFO ALL response can be over 4kb, so we set the default to 8kb.
	DEFAULT_BUFFER = 8192

	// Smallest valid RESP object is ":0\r\n".
	MIN_OBJECT_LENGTH = 4

	// The minimum valid command is "*1\r\n$4\r\nPING\r\n"
	MIN_COMMAND_LENGTH = 14

	// RESP object prefixes
	SIMPLE_STRING_PREFIX = '+'
	ERROR_PREFIX         = '-'
	INTEGER_PREFIX       = ':'
	BULK_STRING_PREFIX   = '$'
	ARRAY_PREFIX         = '*'
)

var (
	// Common responses
	OK   = NewSimpleString("OK")
	PONG = NewSimpleString("PONG")

	// Errors
	ErrSyntaxError = errors.New("resp: syntax error")
	ErrBufferFull  = errors.New("resp: object is larger than buffer")

	lineSuffix = []byte("\r\n")
)

type Object interface {
	Raw() []byte
}

type InvalidObject []byte

func (o InvalidObject) Raw() []byte { return o }

// Parse takes a slice pointing to valid a valid RESP object and returns the
// RESP as the corresponding type.
func Parse(resp []byte) Object {
	switch resp[0] {
	case SIMPLE_STRING_PREFIX:
		return String(resp)
	case ERROR_PREFIX:
		return Error(resp)
	case INTEGER_PREFIX:
		return Integer(resp)
	case BULK_STRING_PREFIX:
		return String(resp)
	case ARRAY_PREFIX:
		return Array(resp)
	default:
		// This will never happen when being used with Reader
		return InvalidObject(resp)
	}
}
