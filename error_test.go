package resp

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError("oops")
	expected := []byte("-oops\r\n")
	if !reflect.DeepEqual(expected, []byte(err)) {
		t.Errorf("expected: %v\ngot: %v", expected, err)
	}
}

func TestErrorBytes_Invalid(t *testing.T) {
	tests := [][]byte{
		// empty
		[]byte(""),
		// wrong type
		[]byte("+oops\r\n"),
		// bad line ending
		[]byte("+oops\r"),
	}

	for i, test := range tests {
		e := Error(test)
		_, err := e.Bytes()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
	}
}

type errorTest struct {
	given    []byte
	expected []byte
}

func TestErrorBytes_Valid(t *testing.T) {
	tests := []errorTest{
		{[]byte("-\r\n"), []byte("")},
		{[]byte("-oops\r\n"), []byte("oops")},
	}

	for i, test := range tests {
		e := Error(test.given)
		bytes, err := e.Bytes()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, bytes) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, bytes)
		}
	}
}
