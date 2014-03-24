package resp

import (
	"bytes"
	"fmt"
)

type BulkString RESP

func NewBulkString(resp []byte) (BulkString, error) {
	if len(resp) < MIN_OBJECT_LENGTH || resp[0] != BULK_STRING_PREFIX || !bytes.HasSuffix(resp, LINE_ENDING) {
		return nil, ErrSyntaxError
	}
	return BulkString(resp), nil
}

func NewBulkStringString(s string) BulkString {
	var buf bytes.Buffer
	buf.WriteByte(BULK_STRING_PREFIX)
	fmt.Fprintf(&buf, "%d", len(s))
	buf.Write(LINE_ENDING)
	buf.WriteString(s)
	buf.Write(LINE_ENDING)
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
