package resp

import (
	"bytes"
	"fmt"
)

type BulkString []byte

func NewBulkString(s string) BulkString {
	var buf bytes.Buffer
	buf.WriteByte(BULK_STRING_PREFIX)
	fmt.Fprintf(&buf, "%d", len(s))
	buf.Write(LineEnding)
	buf.WriteString(s)
	buf.Write(LineEnding)
	return BulkString(buf.Bytes())
}

func (b BulkString) Slice() ([]byte, error) {
	length, lengthEndIndex, err := parseLenLine(b)
	if err != nil {
		return nil, err
	}
	if length == -1 {
		return []byte{}, nil
	}
	return b[lengthEndIndex+1 : len(b)-2], nil
}

func (b BulkString) Bytes() ([]byte, error) {
	slice, err := b.Slice()
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes, nil
}
