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

// Bytes returns a slice of bytes that points to the raw error message. If the
// contents of this Error change, the returned byte slice will be invalid.
func (e Error) Bytes() ([]byte, error) {
	if len(e) < 3 || e[0] != '-' || e[len(e)-2] != '\r' || e[len(e)-1] != '\n' {
		return nil, ErrSyntaxError
	}

	return e[1 : len(e)-2], nil
}
