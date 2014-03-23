package resp

import (
	"reflect"
	"testing"
)

type newCommandTest struct {
	args  []string
	bytes []byte
}

func TestNewCommand(t *testing.T) {
	tests := []newCommandTest{
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

		_, err := command.Slice()
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

		args, err := command.Slice()
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
