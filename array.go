package resp

type Array RESP

func NewArray(resp []byte) (Array, error) {
	if !validRESPLine(ARRAY_PREFIX, resp) {
		return nil, ErrSyntaxError
	}
	return Array(resp), nil
}

// TODO
