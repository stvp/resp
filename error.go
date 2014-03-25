package resp

import (
	"bytes"
)

type Error []byte

func NewError(s string) Error {
	var buf bytes.Buffer
	buf.WriteByte(ERROR_PREFIX)
	buf.WriteString(s)
	buf.Write(lineSuffix)
	return Error(buf.Bytes())
}

func (e Error) Slice() []byte {
	return e[1 : len(e)-2]
}

func (e Error) Bytes() []byte {
	slice := e.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}
