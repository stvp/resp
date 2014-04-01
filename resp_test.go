package resp

import (
	"errors"
	"testing"
)

func TestParse(t *testing.T) {
	// Return error
	e := errors.New("oops")
	obj, err := Parse([]byte("+OK\r\n"), e)
	if err != e {
		t.Error(err)
	}

	// Simple string
	obj, err = Parse([]byte("+OK\r\n"), nil)
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Bulk string
	obj, err = Parse([]byte("$4\r\ncool\r\n"), nil)
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Error
	obj, err = Parse([]byte("-oops\r\n"), nil)
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(Error); !ok {
		t.Errorf("expected Error as first value, got %#v", obj)
	}

	// Integer
	obj, err = Parse([]byte(":123\r\n"), nil)
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(Integer); !ok {
		t.Errorf("expected Integer, got %#v", obj)
	}

	// Array
	obj, err = Parse([]byte("*1\r\n+OK\r\n"), nil)
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(Array); !ok {
		t.Errorf("expected Array, got %#v", obj)
	}
}
