package resp

// An Array is a RESP array, including all of the array's contained RESP
// objects.
type Array []byte

// Raw returns the underlying bytes of this RESP object.
func (a Array) Raw() []byte { return a }

// TODO
// func NewArray(objects interface{}...) Array {}
// func (a Array) Objects() ([]interface{}, error) {}
