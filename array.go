package resp

// An Array is a RESP array, including all of the array's contained RESP
// objects.
type Array []byte

func (a Array) Raw() []byte { return a }

// TODO
// func NewArray(objects interface{}...) Array {}
// func (a Array) Objects() ([]interface{}, error) {}
