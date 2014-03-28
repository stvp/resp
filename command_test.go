package resp

import (
	"errors"
	"reflect"
	"testing"
)

type commandTest struct {
	given    []byte
	expected [][]byte
}

func TestCommandSlice_Valid(t *testing.T) {
	tests := []commandTest{
		{[]byte("*1\r\n$4\r\nPING\r\n"), [][]byte{[]byte("PING")}},
		{[]byte("*2\r\n$4\r\nINFO\r\n$3\r\nALL\r\n"), [][]byte{[]byte("INFO"), []byte("ALL")}},
	}

	for i, test := range tests {
		command := Command(test.given)

		args, err := command.Slices()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, args) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, args)
		}

		args, err = command.Bytes()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, args) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, args)
		}

		expectedStrings := make([]string, len(test.expected))
		for i, bytes := range test.expected {
			expectedStrings[i] = string(bytes)
		}
		stringArgs, err := command.Strings()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(expectedStrings, stringArgs) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, expectedStrings, stringArgs)
		}
	}
}

func TestCommandSlice_Invalid(t *testing.T) {
	tests := [][]byte{
		// empty
		[]byte(""),
		// wrong type
		[]byte("$3\r\nfoo\r\n"),
		// missing array elements
		[]byte("*2\r\n$1\r\nX\r\n"),
		// bad bulk string
		[]byte("*1\r\n$100\r\noops"),
		// nil bulk string
		[]byte("*1\r\n$-1\r\n"),
		// too short
		[]byte("*1\r\n$3\r\nLOL\r\n"),
	}

	for i, test := range tests {
		command := Command(test)

		_, err := command.Slices()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}

		_, err = command.Bytes()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}

		_, err = command.Strings()
		if err == nil {
			t.Errorf("test[%d]: expected error but got none", i)
		}
	}
}

type newCommandStringsTest struct {
	args  []string
	bytes []byte
}

func TestNewCommandStrings(t *testing.T) {
	tests := []newCommandStringsTest{
		{[]string{}, []byte("*0\r\n")},
		{[]string{"PING"}, []byte("*1\r\n$4\r\nPING\r\n")},
		{[]string{"INFO", "ALL"}, []byte("*2\r\n$4\r\nINFO\r\n$3\r\nALL\r\n")},
		{[]string{"INFO", ""}, []byte("*2\r\n$4\r\nINFO\r\n$0\r\n\r\n")},
	}

	for i, test := range tests {
		command := NewCommand(test.args...)
		if !reflect.DeepEqual(test.bytes, []byte(command)) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.bytes, command)
		}
	}
}

func TestParseCommand(t *testing.T) {
	good := [][]byte{
		[]byte("*1\r\n$4\r\nPING\r\n"),
		[]byte("*2\r\n$4\r\nINFO\r\n$3\r\nALL\r\n"),
	}

	bad := [][]byte{
		[]byte(""),
		[]byte("$4\r\noops\r\n"),
		[]byte("*1\r\n$4\r\nPING"),
		[]byte("*1\r\n$4\r\nPING\r"),
	}

	for i, goodBytes := range good {
		_, err := ParseCommand(goodBytes, nil)
		if err != nil {
			t.Errorf("good[%d]: %s", i, err.Error())
		}
	}

	for i, badBytes := range bad {
		_, err := ParseCommand(badBytes, nil)
		if err != ErrSyntaxError {
			t.Errorf("bad[%d]: expected a ErrSyntaxError, got: %#v", i, err)
		}
	}

	expectedError := errors.New("oops")
	_, err := ParseCommand(good[0], expectedError)
	if err != expectedError {
		t.Errorf("passing an error to ParseCommand returned %#v instead of the given error", err)
	}
}
