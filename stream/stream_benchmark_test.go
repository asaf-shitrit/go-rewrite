package stream

import (
	"io/ioutil"
	"strings"
	"testing"
)

func BenchmarkNewParseCtx(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		newParseCtx(nil, nil)
	}
}

func BenchmarkSeekTagEnd(b *testing.B) {
	pc := &parseContext{
		r:             strings.NewReader("<html><body><div id=\"meow\"></div></body></html>"),
		w:             ioutil.Discard,
		runeBuffer:    [3]byte{0, '<', 'h'},
		generalBuffer: make([]byte, 1024),
		end:           false,
		skipWrite:     false,
		i:             1,
	}
	b.ResetTimer()
	seekMatchingTagEnd(pc, "id=meow")
	b.ReportAllocs()
}

func BenchmarkQueryPath(b *testing.B) {

	b.Run("Match", func(b *testing.B) {
		queryPath := queryPath("id=meow")
		value := []byte("id=meow")

		b.ResetTimer()
		queryPath.Match(value)
		b.ReportAllocs()
	})

	b.Run("Type", func(b *testing.B) {
		queryPath := queryPath("id=meow")
		b.ResetTimer()
		queryPath.Type()
		b.ReportAllocs()
	})

	b.Run("kv", func(b *testing.B) {
		queryPath := queryPath("id=meow")
		b.ResetTimer()
		_, _ = queryPath.kv()
		b.ReportAllocs()
	})
}

func BenchmarkSplitStringOnEqual(b *testing.B) {
	testStr := "a=b"
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		splitStringOnEqual(testStr)
	}
	b.ReportAllocs()
}

func BenchmarkSplitBytesOnEqual(b *testing.B) {
	testStr := []byte("a=b")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		splitBytesOnEqual(testStr)
	}
	b.ReportAllocs()

}

func Benchmark_withCtx(b *testing.B) {
	r := strings.NewReader("value")
	_ = withCtx(r, ioutil.Discard, func(pc *parseContext) error {
		return nil
	})
	b.ReportAllocs()
}
