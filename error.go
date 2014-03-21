package resp

import (
	"fmt"
)

var (
	ErrSyntaxError = fmt.Errorf("resp: syntax error")
	ErrBufferFull  = fmt.Errorf("resp: object is larger than buffer")
)
