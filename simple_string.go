package resp

import (
	"bytes"
)

type SimpleString RESP

func NewSimpleString(resp []byte) (SimpleString, error) {
	if !validRESPLine(SIMPLE_STRING_PREFIX, resp) {
		return nil, ErrSyntaxError
	}
	return SimpleString(resp), nil
}

func NewSimpleStringString(s string) SimpleString {
	var buf bytes.Buffer
	buf.WriteByte(SIMPLE_STRING_PREFIX)
	buf.WriteString(s)
	buf.Write(LINE_ENDING)
	return SimpleString(buf.Bytes())
}

func (s SimpleString) Slice() []byte {
	return s[1 : len(s)-2]
}

func (s SimpleString) Bytes() []byte {
	slice := s.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}
