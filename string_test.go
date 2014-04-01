package resp

import (
	"reflect"
	"testing"
)

func TestNewBulkString(t *testing.T) {
	s := NewBulkString("hi")
	expected := []byte("$2\r\nhi\r\n")
	if !reflect.DeepEqual(expected, []byte(s)) {
		t.Errorf("expected: %v\ngot: %v", expected, s)
	}
}
