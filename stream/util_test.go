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
