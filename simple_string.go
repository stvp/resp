package resp

import (
	"bytes"
)

type SimpleString []byte

func NewSimpleString(s string) SimpleString {
	var buf bytes.Buffer
	buf.WriteByte(SIMPLE_STRING_PREFIX)
	buf.WriteString(s)
	buf.Write(LineEnding)
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

func (s SimpleString) IsOk() bool {
	return len(s) == 5 && s[1] == 'O' && s[2] == 'K'
}

func (s SimpleString) IsPong() bool {
	return len(s) == 7 && s[1] == 'P' && s[2] == 'O' && s[3] == 'N' && s[4] == 'G'
}
