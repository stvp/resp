package resp

import (
	"bytes"
)

// parseLenLine takes a RESP array or bulk string length specification and
// returns the specified length as well as the index of the final character of
// the length specification line. It returns a length of -1 for null bulk
// strings. If the line is invalid, an error will be returned. Trailing
// characters are ignored.
func parseLenLine(line []byte) (length int, endIndex int, err error) {
	if len(line) < MIN_OBJECT_LENGTH {
		return -1, -1, ErrSyntaxError
	}

	if line[0] != ARRAY_PREFIX && line[0] != BULK_STRING_PREFIX {
		return -1, -1, ErrSyntaxError
	}

	// Shortcut for null bulk strings and null arrays
	if bytes.HasPrefix(line[1:], nullLength) {
		return -1, 4, nil
	}

	var n int
	var b byte
	var i int
	for i, b = range line[1 : len(line)-2] {
		if b < '0' || b > '9' {
			if b == '\r' {
				return n, i + 2, nil
			} else {
				return -1, i + 3, ErrSyntaxError
			}
		}
		n *= 10
		n += int(b - '0')
	}

	return n, i + 3, nil
}

// indexLineEnd returns the index of the final character of the first line in
// the given RESP byte slice. If no valid line ending can be found, it returns
// -1.
func indexLineEnd(slice []byte) int {
	i := bytes.IndexByte(slice, '\n')
	if i > 0 && slice[i-1] == '\r' {
		return i
	}
	return -1
}
