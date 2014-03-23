package resp

// Error points to the bytes for a valid RESP error and provides methods for
// extracting the error message.
type Error []byte

// NewError accepts an error message and returns the RESP representation of the
// error as an Error.
func NewError(str string) Error {
	buf := make([]byte, 1+len(str)+2)
	buf[0] = '-'
	copy(buf[1:], []byte(str))
	buf[len(buf)-2] = '\r'
	buf[len(buf)-1] = '\n'
	return Error(buf)
}

// Slice returns a slice of bytes that points to the raw error message. If the
// contents of this Error change, the returned byte slice will be invalid. If
// the RESP format is invalid, a Go error will be returned.
func (e Error) Slice() ([]byte, error) {
	if len(e) < 3 || e[0] != '-' || e[len(e)-2] != '\r' || e[len(e)-1] != '\n' {
		return nil, ErrSyntaxError
	}

	return e[1 : len(e)-2], nil
}

// Bytes is the same as Slice except that it returns a copy of the raw error
// message bytes.
func (e Error) Bytes() ([]byte, error) {
	slice, err := e.Slice()
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes, err
}

// String returns the raw error message. If the RESP format is invalid, a Go
// error will be returned.
func (e Error) String() (string, error) {
	slice, err := e.Slice()
	if err != nil {
		return "", err
	}
	return string(slice), err
}
