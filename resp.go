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
)

var (
	ErrSyntaxError = fmt.Errorf("resp: syntax error")
	ErrBufferFull  = fmt.Errorf("resp: object is larger than buffer")
)
