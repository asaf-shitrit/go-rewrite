package html_overwrite

import (
	"github.com/html-overwrite/model"
	"github.com/html-overwrite/std"
	"github.com/html-overwrite/stream"
	"io"
)

// Load allows inputting html docs in string format
// into the stdLibWriter so they could be modified.
func Load(r io.Reader) (w model.Writer, err error) {
	if w, err = std.NewWriter(r); err != nil {
		return
	}

	return
}

func Append(r io.Reader, w io.Writer, path, value string) error {
	return stream.Append(r, w, path, value)
}

func Set(r io.Reader, w io.Writer, path, value string) error {
	return stream.Set(r, w, path, value)
}
