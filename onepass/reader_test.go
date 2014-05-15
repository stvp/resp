package onepass

import (
	"bytes"
	"reflect"
	"testing"
)

type RespTest struct {
	Given    []byte
	Expected []byte
}

func TestReader_ValidRESP(t *testing.T) {
	tests := []RespTest{
		// simple string
		{[]byte("+OK\r\n"), []byte("+OK\r\n")},
		// integer
		{[]byte(":1234\r\n"), []byte(":1234\r\n")},
		// ignore trailing junk
		{[]byte("+OK\r\n..."), []byte("+OK\r\n")},
		// read only one full response
		{[]byte("+OK\r\n+ERR\r\n"), []byte("+OK\r\n")},
		// array
		{[]byte("*2\r\n+OK\r\n+OK\r\n"), []byte("*2\r\n+OK\r\n+OK\r\n")},
		// null array
		{[]byte("*-1\r\n"), []byte("*-1\r\n")},
		// empty array
		{[]byte("*0\r\n"), []byte("*0\r\n")},
		// bulk string
		{[]byte("$4\r\ncool\r\n"), []byte("$4\r\ncool\r\n")},
		// null bulk string
		{[]byte("$-1\r\n"), []byte("$-1\r\n")},
		// bulk string with \r in string
		{[]byte("$3\r\na\rb\r\n"), []byte("$3\r\na\rb\r\n")},
		// bulk string with \n in string
		{[]byte("$3\r\na\nb\r\n"), []byte("$3\r\na\nb\r\n")},
		// bulk string with line ending in string
		{[]byte("$4\r\na\r\nb\r\n"), []byte("$4\r\na\r\nb\r\n")},
		// empty bulk string
		{[]byte("$0\r\n\r\n"), []byte("$0\r\n\r\n")},
		// array of arrays
		{[]byte("*2\r\n*1\r\n+OK\r\n*1\r\n+OK\r\n"), []byte("*2\r\n*1\r\n+OK\r\n*1\r\n+OK\r\n")},
		// array with null bulk string
		{[]byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n"), []byte("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n")},
	}

	for i, test := range tests {
		reader := NewReader(bytes.NewReader(test.Given))
		object, err := reader.ReadObject()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.Expected, object) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.Expected, object)
		}
	}
}
