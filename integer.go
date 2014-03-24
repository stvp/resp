package resp

import (
	"strconv"
)

type Integer []byte

func NewInteger(i int64) Integer {
	buf := []byte{INTEGER_PREFIX}
	strconv.AppendInt(buf, i, 10)
	buf = append(buf, '\r', '\n')
	return Integer(buf)
}

func (i Integer) Int64() (int64, error) {
	return strconv.ParseInt(string(i[1:len(i)-2]), 10, 64)
}
