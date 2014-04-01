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
	OK   = String("+OK\r\n")
	PONG = String("+PONG\r\n")

	// Errors
	ErrSyntaxError = errors.New("resp: syntax error")
	ErrBufferFull  = errors.New("resp: object is larger than buffer")

	lineSuffix = []byte("\r\n")
)

// Parse takes a bytes slice for a single RESP object and returns the bytes
// wrapped with the correct type (String, Error, Integer, or Array). If an
// error is passed in, Parse will simply reply with the given arguments.  If
// the given RESP is invalid, ErrSyntaxError will be returned.
func Parse(resp []byte, err error) (interface{}, error) {
	if err != nil {
		return resp, err
	}
	if len(resp) < MIN_OBJECT_LENGTH || resp[len(resp)-2] != '\r' || resp[len(resp)-1] != '\n' {
		return resp, ErrSyntaxError
	}

	switch resp[0] {
	case SIMPLE_STRING_PREFIX:
		return String(resp), nil
	case ERROR_PREFIX:
		return Error(resp), nil
	case INTEGER_PREFIX:
		return Integer(resp), nil
	case BULK_STRING_PREFIX:
		return String(resp), nil
	case ARRAY_PREFIX:
		return Array(resp), nil
	default:
		return resp, ErrSyntaxError
	}
}
