package resp

import (
	"strconv"
)

// Error points to the bytes for a RESP integer.
type Integer []byte

// NewInteger takes an integer and returns an Integer slice containing a valid
// RESP integer.
func NewInteger(i int64) Integer {
	buf := []byte{INTEGER_PREFIX}
	strconv.AppendInt(buf, i, 10)
	buf = append(buf, '\r', '\n')
	return Integer(buf)
}

// Int returns the value of the RESP integer as an int.
func (i Integer) Int() (int, error) {
	n, err := strconv.Atoi(string(i[1 : len(i)-2]))
	if err != nil {
		return 0, ErrSyntaxError
	}
	return n, nil
}

// Int64 returns the value of the RESP integer as in int64.
func (i Integer) Int64() (int64, error) {
	n, err := strconv.ParseInt(string(i[1:len(i)-2]), 10, 64)
	if err != nil {
		return 0, ErrSyntaxError
	}
	return n, nil
}
