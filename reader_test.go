package resp

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

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

func TestReadObjectSlice_Invalid(t *testing.T) {
	tests := [][]byte{
		// empty
		[]byte{},
		// too small
		[]byte("\r\n"),
		[]byte("-\r\n"),
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

func TestReadObjectSlice_BufferErrors(t *testing.T) {
	reply := []byte("-OK\r\n")
	reader := NewReaderSize(bytes.NewReader(reply), len(reply)-1)
	object, err := reader.ReadObjectSlice()
	if err != ErrBufferFull {
		t.Errorf("expected ErrBufferFull but got %#v", err)
	}
	if !reflect.DeepEqual(object, reply[0:len(reply)-1]) {
		t.Errorf("expected: %v\ngot: %v", reply[0:len(reply)-1], object)
	}
}

type multipleReadTest struct {
	reads    [][]byte
	expected []byte
}

func TestReadObjectSlice_MultipleReads_Valid(t *testing.T) {
	tests := []multipleReadTest{
		{
			[][]byte{[]byte("-O"), []byte("K\r"), []byte("\n")},
			[]byte("-OK\r\n"),
		},
		{
			[][]byte{[]byte("$3\r"), []byte("\nfo"), []byte("o\r\n")},
			[]byte("$3\r\nfoo\r\n"),
		},
		{
			[][]byte{[]byte("*2\r\n"), []byte("-OK\r"), []byte("\n-O"), []byte("K\r\n")},
			[]byte("*2\r\n-OK\r\n-OK\r\n"),
		},
		{
			[][]byte{[]byte("*2\r\n*"), []byte("1\r\n-OK\r"), []byte("\n-O"), []byte("K\r\n")},
			[]byte("*2\r\n*1\r\n-OK\r\n-OK\r\n"),
		},
	}

	for i, test := range tests {
		readers := []io.Reader{}
		for _, piece := range test.reads {
			readers = append(readers, bytes.NewReader(piece))
		}
		reader := NewReader(io.MultiReader(readers...))
		object, err := reader.ReadObjectSlice()
		if err != nil {
			t.Errorf("tests[%d]: %s", i, err.Error())
		} else if !reflect.DeepEqual(test.expected, object) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, object)
		}
	}
}

func TestReadObjectSlice_MultipleReads_Invalid(t *testing.T) {
	tests := []multipleReadTest{
		{
			[][]byte{[]byte("-O"), []byte("K\r")},
			[]byte("-OK\r"),
		},
		{
			[][]byte{[]byte("$3\r")},
			[]byte("$3\r"),
		},
	}

	for i, test := range tests {
		readers := []io.Reader{}
		for _, piece := range test.reads {
			readers = append(readers, bytes.NewReader(piece))
		}
		reader := NewReader(io.MultiReader(readers...))
		object, err := reader.ReadObjectSlice()
		if err == nil {
			t.Errorf("tests[%d]: expected error but didn't get one", i)
		}
		if !reflect.DeepEqual(test.expected, object) {
			t.Errorf("tests[%d]:\nexpected: %v\ngot: %v", i, test.expected, object)
		}
	}
}

type LoopReader struct {
	bytes []byte
	i     int
}

func (r *LoopReader) Read(p []byte) (n int, err error) {
	if len(p) >= len(r.bytes) {
		copy(p, r.bytes)
		return len(r.bytes), nil
	} else {
		panic("TODO handle len(p) < len(r.bytes)")
	}
}

func BenchmarkReaderReadObjectSliceSmall(b *testing.B) {
	// 23 bytes
	resp := []byte("*2\r\n$4\r\nINFO\r\n$3\r\nALL\r\n")
	reader := NewReader(&LoopReader{resp, 0})
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = reader.ReadObjectSlice()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderReadObjectSliceMedium(b *testing.B) {
	// 95 bytes
	resp := []byte("*10\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n")
	reader := NewReader(&LoopReader{resp, 0})
	var err error
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err = reader.ReadObjectSlice()
		if err != nil {
			b.Fatal(err)
		}
	}
}
