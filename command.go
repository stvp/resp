package resp

import (
	"bytes"
	"fmt"
)

// Command points to the bytes for a valid RESP command (an array of bulk
// strings) and provides methods from extracting the raw arguments.
type Command []byte

// NewCommand takes any number of command arguments and returns a Command slice
// pointing to the raw RESP representation of that command.
func NewCommand(args ...string) Command {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "*%d\r\n", len(args))

	for _, arg := range args {
		fmt.Fprintf(&buf, "$%d\r\n%s\r\n", len(arg), arg)
	}

	return Command(buf.Bytes())
}

// Slice returns a slice of byte slices that point to the raw command arguments
// of this Command. If the contents of Command change, the returned byte slices
// will be invalid.
func (c Command) Slice() ([][]byte, error) {
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

// Bytes is the same as Slice except that it returns copies of the byte slices
// for each individual command argument.
func (c Command) Bytes() ([][]byte, error) {
	slices, err := c.Slice()
	if slices == nil {
		return nil, err
	}

	bytes := make([][]byte, len(slices))
	for i, slice := range slices {
		bytes[i] = make([]byte, len(slice))
		copy(bytes[i], slice)
	}

	return bytes, err
}

// Strings returns the Command's arguments as strings. If the underlying RESP
// format is invalid, an error will be returned.
func (c Command) Strings() ([]string, error) {
	slices, err := c.Slice()
	if slices == nil {
		return nil, err
	}

	strings := make([]string, len(slices))
	for i, slice := range slices {
		strings[i] = string(slice)
	}

	return strings, err
}
