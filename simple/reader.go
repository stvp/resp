package redis

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

const (
	SIMPLE_STRING = '+'
	BULK_STRING   = '$'
	INTEGER       = ':'
	ARRAY         = '*'
	ERROR         = '-'
)

var (
	ErrInvalidSyntax = errors.New("resp: invalid syntax")
)

type RESPReader struct {
	*bufio.Reader
}

func NewReader(reader io.Reader) *RESPReader {
	return &RESPReader{
		Reader: bufio.NewReaderSize(reader, 64*1024),
	}
}

func (r *RESPReader) ReadObject() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case SIMPLE_STRING, INTEGER, ERROR:
		return line, nil
	case BULK_STRING:
		return r.readBulkString(line)
	case ARRAY:
		return r.readArray(line)
	default:
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readLine() (line []byte, err error) {
	line, err = r.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	if len(line) > 1 && line[len(line)-2] == '\r' {
		return line, nil
	} else {
		return nil, ErrInvalidSyntax
	}
}

func (r *RESPReader) readBulkString(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}
	if count == -1 {
		return line, nil
	}

	buf := make([]byte, len(line)+count+2)
	copy(buf, line)
	_, err = r.Read(buf[len(line):])
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *RESPReader) readArray(line []byte) ([]byte, error) {
	count, err := r.getCount(line)
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		buf, err := r.ReadObject()
		if err != nil {
			return nil, err
		}
		line = append(line, buf...)
	}

	return line, nil
}

func (r *RESPReader) getCount(line []byte) (int, error) {
	end := bytes.IndexByte(line, '\r')
	return strconv.Atoi(string(line[1:end]))
}
