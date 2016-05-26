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
	buf.WriteByte(bulkStringPrefix)
	buf.WriteString(strconv.Itoa(len(s)))
	buf.Write(lineSuffix)
	buf.WriteString(s)
	buf.Write(lineSuffix)
	return String(buf.Bytes())
}

// NewSimpleString returns a String (as a simple string) with the given contents.
func NewSimpleString(s string) String {
	var buf bytes.Buffer
	buf.WriteByte(simpleStringPrefix)
	buf.WriteString(s)
	buf.Write(lineSuffix)
	return String(buf.Bytes())
}

// Raw returns the underlying bytes of this RESP object.
func (s String) Raw() []byte { return s }

// Slice returns a slice pointing to the string contained in this RESP simple
// string or bulk string.
func (s String) Slice() []byte {
	if s[0] == bulkStringPrefix {
		length, lengthEndIndex, err := parseLenLine(s)
		if err != nil || length == -1 {
			return nil
		}
		return s[lengthEndIndex+1 : len(s)-2]
	}

	// Assume simple string
	return s[1 : len(s)-2]
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
