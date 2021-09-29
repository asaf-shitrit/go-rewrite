package stream

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

const testSetHtmlTemplate = `
<html>
	<head>
	</head>
	<body>
		<h1>Example Header</h1>
		<h2>Sub Header</h2>
		<p id="meow">%s</p>
	</body>
</html>
`

//go:embed static_html/skipped-tags.html
var skippedTagsHTML string

func TestSet(t *testing.T) {

	replacedHTML := "Test Content AAA"
	testHTML := fmt.Sprintf(testSetHtmlTemplate, replacedHTML)

	valueToWrite := "Best HTML Value"
	expectedHTML := fmt.Sprintf(testSetHtmlTemplate, valueToWrite)

	buffer := &bytes.Buffer{}

	err := Set(strings.NewReader(testHTML), buffer, "id=meow", valueToWrite)
	assert.Nil(t, err)

	outputHtml := buffer.String()
	assert.NotContains(t, outputHtml, replacedHTML)
	assert.Equal(t, expectedHTML, outputHtml)

	t.Run("skipped tags check", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		setValue := "<p>Ok</p>"
		err := Set(strings.NewReader(skippedTagsHTML), buffer, "id=content", setValue)
		assert.Nil(t, err)
		output := buffer.String()
		validHTML(t, output)
		assert.Contains(t, output, setValue)
	})
}

const testAppendHtmlTemplate = `
<html>
	<head>
	</head>
	<body>
		<div id="headers">
			<h1>Example Header</h1>
		</div>
	</body>
</html>
`
const testPostAppendHtmlTemplate = `
<html>
	<head>
	</head>
	<body>
		<div id="headers">
			<h1>Example Header</h1>
			<h2>Example Sub Header</h2>
		</div>
	</body>
</html>
`

//go:embed static_html/*
var staticHtmlDir embed.FS

func TestAppend(t *testing.T) {

	t.Run("id", func(t *testing.T) {
		appendedTag := "<h2>Example Sub Header</h2>"

		buffer := &bytes.Buffer{}

		err := Append(strings.NewReader(testAppendHtmlTemplate), buffer, "id=headers", appendedTag)
		assert.Nil(t, err)

		outputHtml := buffer.String()
		validHTML(t, outputHtml)
		equalStripped(t, testPostAppendHtmlTemplate, outputHtml)
	})

	t.Run("tag", func(t *testing.T) {
		t.Run("naive case", func(t *testing.T) {
			appendedTag := "<h2>Example Sub Header</h2>"

			buffer := &bytes.Buffer{}

			err := Append(strings.NewReader(testAppendHtmlTemplate), buffer, "tag=div", appendedTag)
			assert.Nil(t, err)

			outputHtml := buffer.String()
			validHTML(t, outputHtml)
			equalStripped(t, testPostAppendHtmlTemplate, outputHtml)

		})
		t.Run("invalid input", func(t *testing.T) {
			invalidInputs := []string{"a", "3", "<invalid>", "", "meow"}
			for _, input := range invalidInputs {
				t.Run(input, func(t *testing.T) {
					err := Append(strings.NewReader(input), io.Discard, "tag=head", "<div></div>")
					assert.NotNil(t, err)
				})
			}
		})
		t.Run("general websites", func(t *testing.T) {
			sites := []string{
				"https://www.google.com/",
				"https://www.w3schools.com/",
				"https://github.com/",
				"https://cyolo.io/",
				"https://www.microsoft.com/",
				"https://www.playstation.com/",
				"https://www.ebay.com/",
				"https://www.qlik.co.il/",
			}

			injectedValue := "<script>alert(1)</script>"

			for _, site := range sites {
				t.Run(site, func(t *testing.T) {
					res, err := http.DefaultClient.Get(site)
					assert.Nil(t, err)
					output := &bytes.Buffer{}
					err = Append(res.Body, output, "tag=head", injectedValue)
					assert.Nil(t, err)
					assert.Contains(t, output.String(), injectedValue)
					validHTML(t, output.String())
				})
			}
		})

		t.Run("static html files", func(t *testing.T) {
			dirPath := "static_html"
			entries, err := staticHtmlDir.ReadDir(dirPath)
			assert.Nil(t, err)

			injectedValue := "<script>alert(1)</script>"

			for _, entry := range entries {
				t.Run(entry.Name(), func(t *testing.T) {

					content, err := staticHtmlDir.ReadFile(fmt.Sprintf("%s/%s", dirPath, entry.Name()))
					assert.Nil(t, err)

					output := &bytes.Buffer{}
					err = Append(bytes.NewReader(content), output, "tag=head", injectedValue)
					assert.Nil(t, err)
					assert.Contains(t, output.String(), injectedValue)
					validHTML(t, output.String())
				})
			}
		})
	})

}

func TestQueryPath_Match(t *testing.T) {
	path := queryPath("id=3")

	t.Run("empty value", func(t *testing.T) {
		assert.False(t, path.Match([]byte("")))
	})

	t.Run("no match", func(t *testing.T) {
		assert.False(t, path.Match([]byte("href=\"google.com\"")))
	})

	t.Run("match", func(t *testing.T) {
		assert.True(t, path.Match([]byte("id=3")))
	})

	t.Run("invalid path key", func(t *testing.T) {
		invalidPath := queryPath("bad=5")
		assert.False(t, invalidPath.Match([]byte("id=3")))
	})
}
