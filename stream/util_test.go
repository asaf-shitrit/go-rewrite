package stream

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"io/ioutil"
	"strings"
	"testing"
)

func stripNewLineAndTabs(s string) string {
	return strings.NewReplacer("\t", "", "\n", "").Replace(s)
}

func equalStripped(t *testing.T, a, b string) {
	assert.Equal(t, stripNewLineAndTabs(a), stripNewLineAndTabs(b))
}

func validHTML(t *testing.T, s string) {
	nodes, err := html.ParseFragment(strings.NewReader(s), nil)
	if err != nil {
		t.Fatalf("non valid html: %v", err)
	}
	for _, n := range nodes {
		if err := html.Render(ioutil.Discard, n); err != nil {
			t.Fatalf("non valid html: %v", err)
		}
	}
}

func Test_splitBytesOnEqual(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		b := []byte("a=b")
		left, right := splitBytesOnEqual(b)
		assert.Equal(t, left, []byte("a"))
		assert.Equal(t, right, []byte("b"))
	})

	t.Run("should panic", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()
		splitBytesOnEqual([]byte("aaaaa"))
	})
}

func Test_splitStringsOnEqual(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		b := "a=b"
		left, right := splitStringOnEqual(b)
		assert.Equal(t, left, "a")
		assert.Equal(t, right, "b")
	})

	t.Run("should panic", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.NotNil(t, r)
		}()
		splitStringOnEqual("aaaaa")
	})
}

func Test_stripValueParentheses(t *testing.T) {
	t.Run("too short value", func(t *testing.T) {
		a := []byte("a")
		mutatedA := stripValueParentheses(a)
		assert.Equal(t, a, mutatedA)
	})

	t.Run("valid value", func(t *testing.T) {
		a := []byte("\"test value\"")
		a = stripValueParentheses(a)
		assert.Equal(t, a, []byte("test value"))
	})
}
