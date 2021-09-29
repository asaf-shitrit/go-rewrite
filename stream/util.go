package stream

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

func unsafeGetBytes(s string) []byte {
	return (*[0x7fff0000]byte)(unsafe.Pointer(
		(*reflect.StringHeader)(unsafe.Pointer(&s)).Data),
	)[:len(s):len(s)]
}

func splitStringOnEqual(s string) (string, string) {
	i := strings.Index(s, "=")
	if i < 0 {
		panic(errors.New("failed to split key value"))
	}
	return s[:i], s[i+1:]
}

func splitBytesOnEqual(s []byte) ([]byte, []byte) {
	i := bytes.IndexRune(s, '=')
	if i < 0 {
		panic(errors.New("failed to split key value"))
	}
	return s[:i], s[i+1:]
}

func stripValueParentheses(b []byte) []byte {
	if len(b) < 3 {
		return b
	}

	return b[1 : len(b)-1]
}
