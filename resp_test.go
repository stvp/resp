package resp

import (
	"testing"
)

func TestParse(t *testing.T) {
	// Simple string
	obj, err := Parse([]byte("+OK\r\n"))
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Bulk string
	obj, err = Parse([]byte("$4\r\ncool\r\n"))
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Error
	obj, err = Parse([]byte("-oops\r\n"))
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := obj.(Error); !ok {
		t.Errorf("expected Error, got %#v", obj)
	}

	// Integer
	obj, err = Parse([]byte(":123\r\n"))
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(Integer); !ok {
		t.Errorf("expected Integer, got %#v", obj)
	}

	// Array
	obj, err = Parse([]byte("*1\r\n+OK\r\n"))
	if err != nil {
		t.Error(err)
	}
	if _, ok := obj.(Array); !ok {
		t.Errorf("expected Array, got %#v", obj)
	}
}
