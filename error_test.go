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

func TestError_Invalid(t *testing.T) {
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
		_, err := e.Slice()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
		_, err = e.Bytes()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
		_, err = e.String()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
	}
}

type errorTest struct {
	given    []byte
	expected []byte
}

func TestError_Valid(t *testing.T) {
	tests := []errorTest{
		{[]byte("-\r\n"), []byte("")},
		{[]byte("-oops\r\n"), []byte("oops")},
	}

	for i, test := range tests {
		e := Error(test.given)
		slice, err := e.Slice()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, slice) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, slice)
		}
		bytes, err := e.Bytes()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, bytes) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, bytes)
		}
		str, err := e.String()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if string(test.expected) != str {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, string(test.expected), str)
		}
	}
}
