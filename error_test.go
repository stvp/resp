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
	e, err := NewError([]byte("-oops\r\n"))
	if err != nil {
		t.Error(err)
	} else {
		slice := e.Slice()
		if !reflect.DeepEqual([]byte("oops"), slice) {
			t.Errorf("expected: %v\ngot: %v", []byte("oops"), slice)
		}
		bytes := e.Bytes()
		if !reflect.DeepEqual([]byte("oops"), bytes) {
			t.Errorf("expected: %v\ngot: %v", []byte("oops"), bytes)
		}
	}
}

func TestError_Invalid(t *testing.T) {
	tests := [][]byte{
		// empty
		[]byte(""),
		[]byte("-\r\n"),
		// wrong type
		[]byte("+oops\r\n"),
		// bad line ending
		[]byte("+oops\r"),
	}

	for i, test := range tests {
		_, err := NewError(test)
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
	}
}

func TestNewErrorString(t *testing.T) {
	e := NewErrorString("oops")
	expected := []byte("-oops\r\n")
	if !reflect.DeepEqual(expected, []byte(e)) {
		t.Errorf("expected: %v\ngot: %v", expected, e)
	}
}
