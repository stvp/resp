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

// Int64 returns the value of the RESP integer.
func (i Integer) Int64() (int64, error) {
	return strconv.ParseInt(string(i[1:len(i)-2]), 10, 64)
}
