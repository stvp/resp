package resp

import (
	"bytes"
	"io"
)

// Reader implements a buffered RESP object reader for an io.Reader object.
type Reader struct {
	rd   io.Reader
	buf  []byte
	r, w int
	err  error
}

// NewReader returns a new Reader with the default buffer size.
func NewReader(r io.Reader) *Reader {
	return NewReaderSize(r, -1)
}

// NewReaderSize returns a new Reader with the given buffer size. If the buffer
// size is less than 1, the default buffer size will be used.
func NewReaderSize(r io.Reader, size int) *Reader {
	if size < 1 {
		size = DEFAULT_BUFFER
	}

	return &Reader{
		rd:  r,
		buf: make([]byte, size),
	}
}

func (r *Reader) ReadObject() (interface{}, error) {
	return Parse(r.ReadObjectBytes())
}

// ReadObjectSlice reads until the buffer contains one full valid RESP object
// and returns a slice pointing at the slice of the buffer that contains the
// object. The byte slice stops being valid after the next read on this Reader.
// If ReadObjectSlice encounters an error before finding a valid RESP object,
// it returns all data in the buffer and the error itself. A ErrBufferFull
// error typically indicates that the RESP object is larger than the buffer. In
// general. Errors returned by ReadObjectSlice should be considered fatal
// because there's no easy way to recover from them when processing a stream of
// RESP objects.
func (r *Reader) ReadObjectSlice() ([]byte, error) {
	i := r.indexObjectEnd(r.r)
	if i > r.r {
		object := r.buf[r.r : i+1]
		r.r = i + 1
		return object, nil
	}

	for {
		r.fill()

		i := r.indexObjectEnd(r.r)
		if i > r.r {
			object := r.buf[r.r : i+1]
			r.r = i + 1
			return object, nil
		}

		if r.err != nil {
			brokenObject := r.buf[r.r:r.w]
			r.r = 0
			r.w = 0
			return brokenObject, r.readErr()
		}
	}
}

// ReadObjectBytes behaves similarly to ReadObjectSlice except that it returns
// a copied slice of bytes that remains valid after the next read.
func (r *Reader) ReadObjectBytes() ([]byte, error) {
	bytes, err := r.ReadObjectSlice()
	copied := make([]byte, len(bytes))
	copy(copied, bytes)
	return copied, err
}

// Buffered returns the number of bytes currently buffered.
func (r *Reader) Buffered() int {
	return r.w - r.r
}

// indexObjectEnd returns the buffer index of the final character of the object
// beginning at the given position. It returns -1 if a valid object can't be
// found.
func (r *Reader) indexObjectEnd(start int) int {
	if r.Buffered()-start < MIN_OBJECT_LENGTH {
		return -1
	}

	switch r.buf[start] {
	case '+', '-', ':':
		lineEnd := bytes.Index(r.buf[start:], lineSuffix)
		if lineEnd < 0 {
			return -1
		}
		if lineEnd+1 < MIN_OBJECT_LENGTH {
			r.err = ErrSyntaxError
			return -1
		}
		return start + lineEnd + 1
	case '$':
		length, lineEnd, err := parseLenLine(r.buf[start:])
		if err != nil {
			r.err = err
			return -1
		}

		lineEnd += start
		if length == -1 {
			return lineEnd
		} else {
			bulkStringEnd := lineEnd + length + 2
			if bulkStringEnd >= r.w {
				return -1
			}
			return bulkStringEnd
		}
	case '*':
		length, lineEnd, err := parseLenLine(r.buf[start:])
		if err != nil {
			r.err = err
			return -1
		}

		lineEnd += start
		if length == -1 {
			return lineEnd
		} else {
			end := lineEnd
			for i := 0; i < length; i++ {
				end = r.indexObjectEnd(end + 1)
				if end < 0 {
					return -1
				}
			}
			return end
		}
	default:
		r.err = ErrSyntaxError
		return -1
	}
}

// fill reads new data into the buffer, if possible. If the io.Reader returns
// an error, it is set on this Reader for future returning.
func (r *Reader) fill() {
	if r.Buffered() >= len(r.buf)-1 {
		r.err = ErrBufferFull
		return
	}

	if r.r > 0 {
		copy(r.buf, r.buf[r.r:r.w])
		r.w -= r.r
		r.r = 0
	}

	// Add new data
	n, err := r.rd.Read(r.buf[r.w:])
	if n < 0 {
		panic("read negative bytes")
	}
	r.w += n
	if err != nil {
		r.err = err
	}
}

func (r *Reader) readErr() error {
	err := r.err
	r.err = nil
	return err
}
