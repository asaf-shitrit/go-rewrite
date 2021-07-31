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

func idBasedTests(t *testing.T){
	t.Run("simple node set", func(t *testing.T) {

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

	t.Run("simple node append", func(t *testing.T) {
		initialHTML := fmt.Sprintf(BaseHTMLTemplate, "")
		w, err := Load(initialHTML)
		assert.Nil(t, err)
		assert.NotNil(t, w)

		amount := 3

		for i:=0;i<amount;i++ {
			err = w.Append("id=content", TestNode)
			assert.Nil(t, err)
		}

		newHTML := w.String()

		re := regexp.MustCompile(TestNode)
		matches := re.FindAllString(newHTML, -1)

		assert.Equal(t, len(matches), amount)
	})
}

func TestLoad(t *testing.T) {
	t.Run("id based tests", idBasedTests)
}