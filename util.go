package resp

// parseLenLine takes a slice that points to the start of a RESP array or bulk
// string length specification line and returns the array size or bulk string
// length (respectively) and the end index of the length specification line in
// the given slice. If the line is invalid, an error will be returned. All
// bytes after the end of the length specification line are ignored.
func parseLenLine(line []byte) (length int, endIndex int, err error) {
	if len(line) < minObjectLen {
		// Bad line length
		return 0, 0, ErrSyntaxError
	}
	if line[0] != arrayPrefix && line[0] != bulkStringPrefix {
		// Bad line prefix
		return 0, 0, ErrSyntaxError
	}
	if len(line) >= 5 && line[1] == '-' && line[2] == '1' && line[3] == '\r' && line[4] == '\n' {
		// Null length shortcut
		return -1, 4, nil
	}

	var n int
	var b byte
	var i int
	for i, b = range line[1:] {
		if b < '0' || b > '9' {
			if b == '\r' && len(line) > i+2 && line[i+2] == '\n' {
				return n, i + 2, nil
			}
			return 0, 0, ErrSyntaxError
		}
		n = (n * 10) + int(b-'0')
	}

	// Missing line ending
	return 0, 0, ErrSyntaxError
}
