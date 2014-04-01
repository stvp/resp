package resp

import (
	"bytes"
	"fmt"
)

// Command points to the bytes for a RESP command (an array of bulk strings).
type Command []byte

// NewCommand takes a Redis command and arguments and returns a Command byte
// slice pointing to the RESP for the command.
func NewCommand(args ...string) Command {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "*%d\r\n", len(args))

	for _, arg := range args {
		fmt.Fprintf(&buf, "$%d\r\n%s\r\n", len(arg), arg)
	}

	return buf.Bytes()
}

// Slices returns a slice of byte slices that point to each argument in this
// Command. It returns a ErrSyntaxError error if the command RESP bytes are
// invalid.
func (c Command) Slices() ([][]byte, error) {
	// Check for basic validity
	if len(c) < MIN_COMMAND_LENGTH || c[0] != ARRAY_PREFIX || c[len(c)-2] != '\r' || c[len(c)-1] != '\n' {
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

// Bytes is the same as Slices except that it returns slices that point to
// copies of the bytes.
func (c Command) Bytes() ([][]byte, error) {
	slices, err := c.Slices()
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

// Strings is the same as Slices except that it returns strings for each
// argument.
func (c Command) Strings() ([]string, error) {
	slices, err := c.Slices()
	if slices == nil {
		return nil, err
	}

	strings := make([]string, len(slices))
	for i, slice := range slices {
		strings[i] = string(slice)
	}

	return strings, err
}
