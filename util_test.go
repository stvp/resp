package resp

import (
	"testing"
)

type respLenTest struct {
	given         []byte
	length        int
	endIndex      int
	errorExpected bool
}

func TestParseLenLine(t *testing.T) {
	tests := []respLenTest{
		// Invalid lines
		{[]byte{}, 0, -1, true},
		{[]byte(""), 0, -1, true},
		{[]byte("-\r\n"), 0, -1, true},
		{[]byte("-OK\r\n"), 0, -1, true},
		{[]byte("*0x2\r\n"), 0, -1, true},
		{[]byte("*-19\r\n"), 0, -1, true},
		{[]byte("*1"), 0, -1, true},
		{[]byte("*1\r"), 0, -1, true},
		// Valid lines
		{[]byte("*-1\r\n"), -1, 4, false},
		{[]byte("*1\r\n"), 1, 3, false},
		{[]byte("*100\r\n"), 100, 5, false},
		{[]byte("$9876\r\n"), 9876, 6, false},
		{[]byte("$1\r\n$10"), 1, 3, false},
		{[]byte("$1\r\n$100\r\n"), 1, 3, false},
	}

	for i, test := range tests {
		size, endIndex, err := parseLenLine(test.given)
		if test.errorExpected {
			if err == nil {
				t.Errorf("tests[%d]: expected an error but didn't get one", i)
			}
		} else {
			if err != nil {
				t.Errorf("tests[%d]: %s", i, err.Error())
			} else {
				if test.length != size {
					t.Errorf("tests[%d]: expected: %v, got: %v", i, test.length, size)
				}
				if test.endIndex != endIndex {
					t.Errorf("tests[%d]: expected: %v, got: %v", i, test.endIndex, endIndex)
				}
			}
		}
	}
}

func BenchmarkParseLenLine(b *testing.B) {
	line := []byte("*250000\r\n")
	for i := 0; i < b.N; i++ {
		parseLenLine(line)
	}
}
