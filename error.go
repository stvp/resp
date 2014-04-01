package resp

import (
	"bytes"
)

// An Error is a RESP error byte slice.
type Error []byte

// NewError returns a RESP error with the given error message.
func NewError(msg string) Error {
	var buf bytes.Buffer
	buf.WriteByte(ERROR_PREFIX)
	buf.WriteString(msg)
	buf.Write(lineSuffix)
	return Error(buf.Bytes())
}

// Slice returns a slice pointing to this Error's message bytes.
func (e Error) Slice() []byte {
	return e[1 : len(e)-2]
}

// Bytes is the same as Slice except that it returns a copied slice.
func (e Error) Bytes() []byte {
	slice := e.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}

// Error returns the error message. This allows Error to satisfy the error
// interface.
func (e Error) Error() string {
	return string(e.Slice())
}
