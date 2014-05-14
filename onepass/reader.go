package onepass

import (
	"fmt"
	"io"
)

const (
	SIMPLE_STRING_PREFIX = '+'
	ERROR_PREFIX         = '-'
	INTEGER_PREFIX       = ':'
	BULK_STRING_PREFIX   = '$'
	ARRAY_PREFIX         = '*'
)

type Reader struct {
	reader io.Reader
	buf    []byte
	state  stateFn
	err    error

	// start position of RESP object currently being read
	r int
	// current position in RESP buffer
	i int
	// index where new data should be written to the buffer
	w int
}

func NewReader(reader io.Reader) *Reader {
	r := &Reader{
		reader: reader,
		buf:    make([]byte, 8192),
		r:      0,
		i:      -1,
		w:      0,
	}
	return r
}

// ReadObject reads a single RESP object from the io.Reader wrapped by this
// Reader. It returns an error if the underlying io.Reader returns an error or
// if the RESP is invalid.
func (r *Reader) ReadObject() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	for r.state = handlePrefix; r.state != nil; {
		r.state = r.state(r)
	}

	return r.cut(), r.err
}

// -- Buffer methods

// fill reads more data in to the buffer, growing the buffer as needed.
func (r *Reader) fill() error {
	// If the active buffer window starts past the halfway point of the whole
	// buffer, slide the window over to free up space for writing.
	if r.r > len(r.buf)/2 {
		copy(r.buf, r.buf[r.r:r.w])
		r.i -= r.r
		r.w -= r.r
		r.r = 0
	}

	// If the active buffer window is over half of the whole buffer size, double
	// the buffer size.
	if r.w > len(r.buf)/2 {
		buf := make([]byte, len(r.buf)*2)
		copy(buf, r.buf)
		r.buf = buf
	}

	n, err := r.reader.Read(r.buf[r.w:])
	if n < 0 {
		return fmt.Errorf("read a negative number of bytes")
	}
	r.w += n

	return err
}

// forward moves the cursor forward in the buffer n number of bytes, reading
// new data from the wrapped io.Reader as needed. It sets an error on the
// Reader if it is unable to fill the buffer with at least n bytes.
func (r *Reader) forward(n int) {
	if n > 0 {
		for r.i+n >= r.w {
			r.err = r.fill()
			if r.err != nil {
				return
			}
		}
		r.i = r.i + n
	} else if n < 0 {
		panic("tried to read backwards")
	}
}

// forwardLine moves the cursor forward to the end of the next RESP line
// separator, reading new data as needed. It sets an error on the Reader if it
// is unable to find a line ending.
func (r *Reader) forwardLine() {
	for {
		r.forward(1)
		if r.err != nil || (r.buf[r.i-1] == '\r' && r.buf[r.i] == '\n') {
			return
		}
	}
}

// forwardLineEnding moves the cursor forward over a RESP line ending. It sets
// an error on the Reader if the next bytes are not a RESP line ending.
func (r *Reader) forwardLineEnding() {
	r.forward(2)
	if r.err == nil {
		if r.buf[r.i-1] != '\r' || r.buf[r.i] != '\n' {
			r.err = fmt.Errorf("invalid line ending: %s", string(r.buf[r.i-1:r.i]))
		}
	}
}

// forwardLength moves the cursor over a bulk string or array length
// specification, returning the length. It returns -1 for null lengths. If the
// length contains invalid characters, it sets an error on the Reader.
func (r *Reader) forwardLength() (length int) {
	var b byte

	r.forward(1)
	if r.err != nil {
		return
	}

	// Check initial byte so we can escape early for null length specifications.
	b = r.buf[r.i]
	if b == '-' {
		r.forward(1)
		if r.err != nil {
			return
		}
		if r.buf[r.i] == '1' {
			return -1
		}
	} else if b >= '0' && b <= '9' {
		length = int(b - '0')
	} else {
		r.err = fmt.Errorf("invalid character in length: %s", string(b))
		return
	}

	// Get the rest of the length.
	for {
		r.forward(1)
		if r.err != nil {
			return
		}

		b = r.buf[r.i]
		if b >= '0' && b <= '9' {
			length = (length * 10) + int(b-'0')
		} else if b == '\r' {
			r.forward(1)
			if r.err != nil {
				return
			} else if r.buf[r.i] == '\n' {
				return length
			} else {
				r.err = fmt.Errorf("invalid line ending: %s", string([]byte{b, r.buf[r.i]}))
				return
			}
		} else {
			r.err = fmt.Errorf("invalid character in length: %s", string(b))
			return
		}
	}
}

// cut returns a copy of the current active buffer window and resets the
// window.
func (r *Reader) cut() []byte {
	slice := make([]byte, r.i+1)
	copy(slice, r.buf[r.r:r.i+1])
	r.r = r.i + 1
	return slice
}

// -- State machine

// stateFn takes a Reader and returns the next state. If the next state is nil,
// we've either fully read a single RESP object, or we've encountered an error.
type stateFn func(*Reader) stateFn

func handlePrefix(r *Reader) stateFn {
	r.forward(1)
	if r.err != nil {
		return nil
	}

	switch r.buf[r.i] {
	case SIMPLE_STRING_PREFIX, ERROR_PREFIX:
		return handleLine(r)
	case BULK_STRING_PREFIX:
		return handleBulkString(r)
	case ARRAY_PREFIX:
		return handleArray(r)
	default:
		r.err = fmt.Errorf("invalid prefix: %s", string(r.buf[r.i]))
		return nil
	}
}

func handleLine(r *Reader) stateFn {
	r.forwardLine()
	return nil
}

func handleBulkString(r *Reader) stateFn {
	length := r.forwardLength()
	if r.err != nil {
		return nil
	}
	if length > 0 {
		r.forward(length)
	}
	r.forwardLineEnding()
	return nil
}

func handleArray(r *Reader) stateFn {
	length := r.forwardLength()
	if r.err != nil {
		return nil
	}
	if length == -1 {
		r.forwardLineEnding()
		return nil
	}
	for i := 0; i < length; i++ {
		for r.state = handlePrefix; r.state != nil; {
			r.state = r.state(r)
		}
		if r.err != nil {
			break
		}
	}
	return nil
}
