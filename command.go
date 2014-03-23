package resp

// Command points to the bytes for a valid RESP command (an array of bulk
// strings) and provides methods from extracting the raw arguments.
type Command []byte

// Args returns a slice of byte slices that point to the raw command arguments
// of this Command. If the contents of Command change, the returned byte slices
// will be invalid.
func (c Command) Args() ([][]byte, error) {
	if len(c) < MIN_COMMAND_LENGTH || c[0] != '*' {
		return nil, ErrSyntaxError
	}

	// Find the number of args
	argCount, cursor, err := parseLenLine(c)
	if err != nil {
		return nil, err
	}

	args := make([][]byte, argCount)
	var end, length int
	for i, _ := range args {
		cursor += 1
		if cursor >= len(c) {
			return nil, ErrSyntaxError
		}

		length, end, err = parseLenLine(c[cursor:])
		cursor += end + 1
		if err != nil {
			return nil, err
		}
		// Null bulk strings are invalid in RESP commands
		if length < 0 {
			return nil, ErrSyntaxError
		}

		if cursor+length+2 > len(c) {
			return nil, ErrSyntaxError
		}

		args[i] = c[cursor : cursor+length]

		// Move cursor to final character of current line
		cursor += length + 1
	}

	return args, nil
}
