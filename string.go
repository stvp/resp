package resp

var (
	// Common strings
	OK   = String("+OK\r\n")
	PONG = String("+PONG\r\n")
)

// String points to the bytes for a RESP simple string or bulk string.
type String []byte

// Slice returns a slice pointing to the bytes of the string contents.
func (s String) Slice() []byte {
	if s[0] == BULK_STRING_PREFIX {
		length, lengthEndIndex, err := parseLenLine(s)
		if err != nil {
			return nil
		} else if length == -1 {
			return []byte{}
		} else {
			return s[lengthEndIndex+1 : len(s)-2]
		}
	}

	// Otherwise, assume simple string
	return s[1 : len(s)-2]
}

// Bytes is the same as Slice except that it returns a copied slice.
func (s String) Bytes() []byte {
	slice := s.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}
