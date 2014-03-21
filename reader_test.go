package resp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadObjectSlice_Invalid(t *testing.T) {
	tests := [][]byte{
		// empty
		[]byte{},
		// no delimiter
		[]byte("-OK"),
		// invalid delimiter
		[]byte("-OK\r"),
		// invalid prefix
		[]byte("OK\r\n"),
		// array with invalid length
		[]byte("*5\r\n-OK\r\n"),
	}

	for i, test := range tests {
		reader := NewReader(bytes.NewReader(test))
		_, err := reader.ReadObjectSlice()
		if err == nil {
			t.Errorf("tests[%d]: expected an error but didn't get one", i)
		}
	}
}

type respTest struct {
	given    []byte
	expected []byte
}

func TestReadObjectSlice_Valid(t *testing.T) {
	tests := []respTest{
		// simple string
		{[]byte("-OK\r\n"), []byte("-OK\r\n")},
		// ignore trailing junk
		{[]byte("-OK\r\n..."), []byte("-OK\r\n")},
		// read only one full response
		{[]byte("-OK\r\n-ERR\r\n"), []byte("-OK\r\n")},
		// array
		{[]byte("*2\r\n-OK\r\n-OK\r\n"), []byte("*2\r\n-OK\r\n-OK\r\n")},
		// empty array
		{[]byte("*0\r\n"), []byte("*0\r\n")},
		// bulk string
		{[]byte("$4\r\ncool\r\n"), []byte("$4\r\ncool\r\n")},
		// null bulk string
		{[]byte("$-1\r\n"), []byte("$-1\r\n")},
		// empty bulk string
		{[]byte("$0\r\n\r\n"), []byte("$0\r\n\r\n")},
		// array of arrays
		{[]byte("*2\r\n*1\r\n-OK\r\n*1\r\n-OK\r\n"), []byte("*2\r\n*1\r\n-OK\r\n*1\r\n-OK\r\n")},
		// array with null bulk string
		{[]byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n"), []byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")},
	}

	for i, test := range tests {
		reader := NewReader(bytes.NewReader(test.given))
		object, err := reader.ReadObjectSlice()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, object) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, object)
		}
	}
}
