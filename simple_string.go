package resp

type String []byte

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

	// Assume simple string
	return s[1 : len(s)-2]
}

func (s String) Bytes() []byte {
	slice := s.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}
