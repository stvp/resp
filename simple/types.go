package redis

import (
	"bytes"
	"errors"
	"fmt"
)

func String(resp []byte, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch resp[0] {
	case SIMPLE_STRING:
		return string(resp[1 : len(resp)-2]), nil
	case BULK_STRING:
		start := bytes.IndexByte(resp, '\n')
		return string(resp[start+1 : len(resp)-2]), nil
	case ERROR:
		return "", errors.New(string(resp[1 : len(resp)-2]))
	default:
		return "", incorrectTypeError("string", resp[0])
	}
}

func Ok(resp []byte, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	switch resp[0] {
	case SIMPLE_STRING:
		return resp[1] == 'O' && resp[2] == 'K', nil
	case ERROR:
		return false, errors.New(string(resp[1 : len(resp)-2]))
	default:
		return false, incorrectTypeError("OK", resp[0])
	}
}

// func Int(resp []byte, err error) (int, error) {
// if err != nil {
// return 0, err
// }

// }

func incorrectTypeError(expected string, prefix byte) error {
	switch prefix {
	case SIMPLE_STRING, BULK_STRING:
		return fmt.Errorf("%s expected, but received string")
	case ERROR:
		return fmt.Errorf("%s expected, but received error")
	case INTEGER:
		return fmt.Errorf("%s expected, but received integer")
	case ARRAY:
		return fmt.Errorf("%s expected, but received array")
	default:
		return fmt.Errorf("%s expected, but received invalid prefix %q", prefix)
	}
}
