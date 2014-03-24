package resp

import (
	"reflect"
	"testing"
)

type errorTest struct {
	given    []byte
	expected []byte
}

func TestError_Valid(t *testing.T) {
	e := Error([]byte("-oops\r\n"))
	slice := e.Slice()
	if !reflect.DeepEqual([]byte("oops"), slice) {
		t.Errorf("expected: %v\ngot: %v", []byte("oops"), slice)
	}
	bytes := e.Bytes()
	if !reflect.DeepEqual([]byte("oops"), bytes) {
		t.Errorf("expected: %v\ngot: %v", []byte("oops"), bytes)
	}
}

func TestNewError(t *testing.T) {
	e := NewError("oops")
	expected := []byte("-oops\r\n")
	if !reflect.DeepEqual(expected, []byte(e)) {
		t.Errorf("expected: %v\ngot: %v", expected, e)
	}
}
