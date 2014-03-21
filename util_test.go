package resp

import (
	"testing"
)

type respLenTest struct {
	given         []byte
	expected      int
	errorExpected bool
}

func TestParseLenLine(t *testing.T) {
	tests := []respLenTest{
		// Invalid lines
		{[]byte{}, 0, true},
		{[]byte(""), 0, true},
		{[]byte("-\r\n"), 0, true},
		{[]byte("-OK\r\n"), 0, true},
		{[]byte("*0x2\r\n"), 0, true},
		{[]byte("*-19\r\n"), 0, true},
		// Valid lines
		{[]byte("*1\r\n"), 1, false},
		{[]byte("*100\r\n"), 100, false},
		{[]byte("$9876\r\n"), 9876, false},
	}

	for i, test := range tests {
		size, err := parseLenLine(test.given)
		if test.errorExpected {
			if err == nil {
				t.Errorf("tests[%d]: expected an error but didn't get one", i)
			}
		} else {
			if err != nil {
				t.Errorf("tests[%d]: %s", i, err.Error())
			} else if test.expected != size {
				t.Errorf("tests[%d]: expected: %v, got: %v", i, test.expected, size)
			}
		}
	}
}
