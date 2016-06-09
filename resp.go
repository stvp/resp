package resp

import (
	"errors"
)

const (
	// A large INFO ALL response can be over 4kb, so we set the default to 8kb.
	defaultBufferLen = 8192

	// Smallest valid RESP object is ":0\r\n".
	minObjectLen = 4

	// The minimum valid command is a single character bulk string: "*1\r\n$1\r\nX\r\n"
	minCommandLen = 11

	// RESP object prefixes
	simpleStringPrefix = '+'
	errorPrefix        = '-'
	integerPrefix      = ':'
	bulkStringPrefix   = '$'
	arrayPrefix        = '*'
)

var (
	// Common responses
	OK   = NewSimpleString("OK")
	PONG = NewSimpleString("PONG")

	// ErrSyntaxError is returned when invalid RESP is encountered.
	ErrSyntaxError = errors.New("resp: syntax error")

	// ErrBufferFull is returned when a RESP object is larger than the
	// buffer can accommodate.
	ErrBufferFull = errors.New("resp: object is larger than buffer")

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
	case simpleStringPrefix:
		return String(resp)
	case errorPrefix:
		return Error(resp)
	case integerPrefix:
		return Integer(resp)
	case bulkStringPrefix:
		return String(resp)
	case arrayPrefix:
		return Array(resp)
	default:
		// This will never happen when being used with Reader
		return InvalidObject(resp)
	}
}
