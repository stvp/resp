package resp

import (
	"testing"
)

func TestParse(t *testing.T) {
	// Simple string
	obj := Parse([]byte("+OK\r\n"))
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Bulk string
	obj = Parse([]byte("$4\r\ncool\r\n"))
	if _, ok := obj.(String); !ok {
		t.Errorf("expected String, got %#v", obj)
	}

	// Error
	obj = Parse([]byte("-oops\r\n"))
	if _, ok := obj.(Error); !ok {
		t.Errorf("expected Error as first value, got %#v", obj)
	}

	// Integer
	obj = Parse([]byte(":123\r\n"))
	if _, ok := obj.(Integer); !ok {
		t.Errorf("expected Integer, got %#v", obj)
	}

	// Array
	obj = Parse([]byte("*1\r\n+OK\r\n"))
	if _, ok := obj.(Array); !ok {
		t.Errorf("expected Array, got %#v", obj)
	}
}
