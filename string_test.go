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

func TestNewSimpleString(t *testing.T) {
	s := NewSimpleString("hi")
	expected := []byte("+hi\r\n")
	if !reflect.DeepEqual(expected, []byte(s)) {
		t.Errorf("expected: %v\ngot: %v", expected, s)
	}
}
