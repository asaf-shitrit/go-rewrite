package html_overwrite

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

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

func idBasedTests(t *testing.T) {
	t.Run("Set", func(t *testing.T) {

		createdPNode := `<p id="time">Maybe this time.</p>`
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, TestNode)

		w, err := Load(initialHTML)
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
		w, err := Load(initialHTML)
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

func classBasedTests(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, DivsWithClasses)
		w, err := Load(initialHTML)
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
		w, err := Load(initialHTML)
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

func TestLoad(t *testing.T) {
	t.Run("id based tests", idBasedTests)
	t.Run("class based tests", classBasedTests)
}
