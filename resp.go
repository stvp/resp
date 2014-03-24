package resp

import (
	"bytes"
	"fmt"
)

const (
	// A large INFO ALL response can be over 4kb, so we set the default to 8kb.
	DEFAULT_BUFFER = 8192

	// Smallest valid RESP object is ":0\r\n".
	MIN_OBJECT_LENGTH = 4

	// The minimum valid command is "*1\r\n$4\r\nPING\r\n"
	MIN_COMMAND_LENGTH = 14

	// RESP type prefixes
	SIMPLE_STRING_PREFIX = '+'
	ERROR_PREFIX         = '-'
	INTEGER_PREFIX       = ':'
	BULK_STRING_PREFIX   = '$'
	ARRAY_PREFIX         = '*'

	// Common replies
	OK   = "OK"
	PONG = "PONG"
)

var (
	LineEnding = []byte{'\r', '\n'}

	// Errors
	ErrSyntaxError = fmt.Errorf("resp: syntax error")
	ErrBufferFull  = fmt.Errorf("resp: object is larger than buffer")
)

func Load(line []byte) (interface{}, error) {
	if len(line) < MIN_OBJECT_LENGTH || !bytes.HasSuffix(line, LineEnding) {
		return nil, ErrSyntaxError
	}

	switch line[0] {
	case SIMPLE_STRING_PREFIX:
		if len(line) == 5 && line[1] == 'O' && line[2] == 'K' {
			return OK, nil
		} else if len(line) == 7 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G' {
			return PONG, nil
		} else {
			return SimpleString(line), nil
		}
	case ERROR_PREFIX:
		return Error(line), nil
	case INTEGER_PREFIX:
		return Integer(line), nil
	case BULK_STRING_PREFIX:
		return BulkString(line), nil
	case ARRAY_PREFIX:
		return Array(line), nil
	default:
		return nil, ErrSyntaxError
	}
}
