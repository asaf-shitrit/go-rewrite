package html_overwrite

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
)

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

const BaseHTMLTemplate = `
<html>
	<head>
	</head>
	<body>
		<div id="content">
			%s
		</div>
	</body>
</html>
`

const TestNode = `<p>Very Cool</p>`

func stdLibIdBasedTests(t *testing.T) {
	t.Run("Set", func(t *testing.T) {

		createdPNode := `<p id="time">Maybe this time.</p>`
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, TestNode)

		w, err := Load(strings.NewReader(initialHTML))
		assert.Nil(t, err)
		assert.NotNil(t, w)

		err = w.Set("id=content", createdPNode)
		assert.Nil(t, err)

		newHtml := w.String()
		assert.NotContains(t, newHtml, TestNode)
		assert.Contains(t, newHtml, createdPNode)
	})

	t.Run("Append", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, "")
		w, err := Load(strings.NewReader(initialHTML))
		assert.Nil(t, err)
		assert.NotNil(t, w)

		amount := 3

		for i := 0; i < amount; i++ {
			err = w.Append("id=content", TestNode)
			assert.Nil(t, err)
		}

		newHTML := w.String()

		re := regexp.MustCompile(TestNode)
		matches := re.FindAllString(newHTML, -1)

		assert.Equal(t, len(matches), amount)
	})
}

const DivsWithClasses = `
<div class="great-name"></div>
<div class="great-name second-great-name"></div>
<div id="good stuff" class="great-name"></div>
`

const DivsWithClassesAndContent = `
<div class="great-name">
	<h3>the best place</h3>
</div>
<div class="great-name second-great-name">
	<h3>to go</h3>
</div>
<div id="good stuff" class="great-name">
	<h3>shopping</h3>
</div>
`

func stdLibClassBasedTests(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, DivsWithClasses)
		w, err := Load(strings.NewReader(initialHTML))
		assert.Nil(t, err)
		assert.NotNil(t, w)
		injectedPNode := `<p>Great Content Also!</p>`
		err = w.Set("class=great-name", injectedPNode)
		assert.Nil(t, err)
		newHtmlString := w.String()
		assert.NotNil(t, newHtmlString)
		re := regexp.MustCompile(injectedPNode)
		matches := re.FindAllString(newHtmlString, -1)
		assert.Equal(t, len(matches), 3)
	})

	t.Run("Append", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, DivsWithClassesAndContent)
		w, err := Load(strings.NewReader(initialHTML))
		assert.Nil(t, err)
		assert.NotNil(t, w)
		injectedPNode := `<p>Come visit Kenyon Lev Hadera</p>`
		err = w.Append("class=great-name", injectedPNode)
		assert.Nil(t, err)
		newHtmlString := w.String()
		assert.NotNil(t, newHtmlString)
		re := regexp.MustCompile(injectedPNode)
		matches := re.FindAllString(newHtmlString, -1)
		assert.Equal(t, len(matches), 3)
		assert.Contains(t, newHtmlString, "the best place")
		assert.Contains(t, newHtmlString, "to go")
		assert.Contains(t, newHtmlString, "shopping")
	})
}

func streamIdBasedTests(t *testing.T) {

	t.Run("Set", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, `<p id="time">Maybe this time.</p>`)
		outputHTML := &bytes.Buffer{}
		err := Set(strings.NewReader(initialHTML), outputHTML, "id=time", "This time for sure.")
		assert.Nil(t, err)

		output := outputHTML.String()
		validHTML(t, output)
		assert.Contains(t, output, "This time for sure")
	})

	t.Run("Append", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, `<p>Maybe this time.</p>`)
		outputHTML := &bytes.Buffer{}
		appendedValue := "<p>Another one</p>"
		err := Set(strings.NewReader(initialHTML), outputHTML, "id=content", appendedValue)
		assert.Nil(t, err)

		output := outputHTML.String()
		validHTML(t, output)
		assert.Contains(t, output, appendedValue)
	})
}

func streamTagBasedTests(t *testing.T) {
	t.Run("Append", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, "")
		outputHTML := &bytes.Buffer{}
		injectedValue := "<script>alert(1)</script>"
		err := Append(strings.NewReader(initialHTML), outputHTML, "tag=head", injectedValue)
		assert.Nil(t, err)

		output := outputHTML.String()
		validHTML(t, output)
		assert.Contains(t, output, injectedValue)
	})
}

type ErrReader struct{ Error error }

func (e *ErrReader) Read([]byte) (int, error) {
	return 0, e.Error
}

func TestHtmlMutations(t *testing.T) {
	t.Run("std lib based", func(t *testing.T) {
		t.Run("bad inputs", func(t *testing.T) {
			t.Run("err reader", func(t *testing.T) {
				r := &ErrReader{Error: errors.New("random failure")}
				_, err := Load(r)
				assert.NotNil(t, err)
			})
		})
		t.Run("id based tests", stdLibIdBasedTests)
		t.Run("class based tests", stdLibClassBasedTests)
	})
	t.Run("stream based", func(t *testing.T) {
		t.Run("id based tests", streamIdBasedTests)
		t.Run("tag based tests", streamTagBasedTests)
	})
}
