package resp

import (
	"bytes"
	"strconv"
)

// A String is a RESP simple string or bulk string.
type String []byte

// NewBulkString returns a String (as a bulk string) with the given contents.
func NewBulkString(s string) String {
	var buf bytes.Buffer
	buf.WriteByte(BULK_STRING_PREFIX)
	buf.WriteString(strconv.Itoa(len(s)))
	buf.Write(lineSuffix)
	buf.WriteString(s)
	buf.Write(lineSuffix)
	return String(buf.Bytes())
}

// NewSimpleString returns a String (as a simple string) with the given contents.
func NewSimpleString(s string) String {
	var buf bytes.Buffer
	buf.WriteByte(SIMPLE_STRING_PREFIX)
	buf.WriteString(s)
	buf.Write(lineSuffix)
	return String(buf.Bytes())
}

func (s String) Raw() []byte { return s }

// Slice returns a slice pointing to the string contained in this RESP simple
// string or bulk string.
func (s String) Slice() []byte {
	if s[0] == BULK_STRING_PREFIX {
		length, lengthEndIndex, err := parseLenLine(s)
		if err != nil || length == -1 {
			return nil
		} else {
			return s[lengthEndIndex+1 : len(s)-2]
		}
	} else {
		// Assume simple string
		return s[1 : len(s)-2]
	}
}

// Bytes is the same as Slice except that it returns a copied slice.
func (s String) Bytes() []byte {
	slice := s.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}

func (s String) String() string {
	return string(s.Slice())
}
