package resp

// Error points to the bytes for a valid RESP error and provides methods for
// extracting the error message.
type Error RESP

func NewError(resp []byte) (Error, error) {
	if !validRESPLine(ERROR_PREFIX, resp) {
		return nil, ErrSyntaxError
	}
	return Error(resp), nil
}

// NewErrorString takes an error message and returns an Error pointing to the
// RESP representation of that error message.
func NewErrorString(s string) Error {
	bytes := make([]byte, 1+len(s)+2)
	bytes[0] = '-'
	copy(bytes[1:], []byte(s))
	bytes[len(bytes)-2] = '\r'
	bytes[len(bytes)-1] = '\n'
	return Error(bytes)
}

func (e Error) Slice() []byte {
	return e[1 : len(e)-2]
}

func (e Error) Bytes() []byte {
	slice := e.Slice()
	bytes := make([]byte, len(slice))
	copy(bytes, slice)
	return bytes
}
