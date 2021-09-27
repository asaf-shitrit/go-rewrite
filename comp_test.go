package html_overwrite

import (
	"bytes"
	"strings"
	"testing"
)

const testHTML = `
<html>
	<head>
	</head>
	<body id="body">
		<h1>Example Header</h1>
		<h2>Sub Header</h2>
		<p id="meow">Asaf Test Content</p>
	</body>
</html>
`

const testValue = "Great Test Value"

func BenchmarkAppend(b *testing.B) {

	appendedValue := "<div></div>"

	b.Run("std lib", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			writer, _ := Load(testHTML)
			_ = writer.Append("id=body", appendedValue)
		}
		b.ReportAllocs()
	})

	b.Run("stream", func(b *testing.B) {
		readers := make([]*strings.Reader, b.N)
		buffers := make([]*bytes.Buffer, b.N)

		for n := 0; n < b.N; n++ {
			readers[n] = strings.NewReader(testHTML)
			buffers[n] = bytes.NewBuffer(make([]byte, 0, len(testHTML)*2))
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = Append(readers[n], buffers[n], "id=body", appendedValue)
		}
		b.ReportAllocs()
	})
}

func BenchmarkSet(b *testing.B) {
	b.Run("std lib", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			writer, _ := Load(testHTML)
			_ = writer.Set("id=meow", testValue)
		}
		b.ReportAllocs()
	})

	b.Run("stream", func(b *testing.B) {
		readers := make([]*strings.Reader, b.N)
		buffers := make([]*bytes.Buffer, b.N)

		for n := 0; n < b.N; n++ {
			readers[n] = strings.NewReader(testHTML)
			buffers[n] = bytes.NewBuffer(make([]byte, 0, len(testHTML)*2))
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			_ = Set(readers[n], buffers[n], "id=meow", testValue)
		}
		b.ReportAllocs()
	})
}
