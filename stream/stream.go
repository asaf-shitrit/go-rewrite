package stream

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type parseContext struct {
	r             io.Reader // reader to read from
	w             io.Writer // writer to write to
	runeBuffer    [3]byte   // buffer that holds current seeked runes
	generalBuffer []byte    // buffer to hold attributes to match to
	end           bool
	skipWrite     bool
	i             int
}

func (pc *parseContext) resetGeneralBuffer() {
	pc.generalBuffer = pc.generalBuffer[:0]
}

func (pc *parseContext) attribute() []byte {
	return pc.generalBuffer
}

func (pc *parseContext) initialRead() {
	now := pc.now()

	_, err := pc.r.Read(pc.runeBuffer[1:])
	if err != nil {
		panic(err)
	}

	// assign the newly read values from the read buffer
	// to the runes buffer
	pc.runeBuffer[0] = now
}

func (pc *parseContext) read() {
	// read & shift buffer positions
	now := pc.now()
	following := pc.following()

	_, err := pc.r.Read(pc.runeBuffer[2:])
	if err != nil {
		if err.Error() == "EOF" {
			pc.end = true
			return
		}
		panic(err)
	}

	// assign the newly read values from the read buffer
	// to the runes buffer
	pc.runeBuffer[0], pc.runeBuffer[1] = now, following
}

func (pc *parseContext) writeOutput() []byte {

	if pc.end {
		return pc.runeBuffer[2:]
	}

	return pc.runeBuffer[1:2]
}

func (pc *parseContext) next() {

	if pc.end {
		panic(errors.New("cannot continue when reader is at EOF"))
	}

	if pc.i == 0 {
		// on first iteration we need to read twice
		pc.initialRead()
	} else {
		pc.read()
	}

	if !pc.skipWrite {
		// we need to write
		if _, err := pc.w.Write(pc.writeOutput()); err != nil {
			panic(fmt.Errorf("failed to write output: %v", err))
		}
	}

	pc.i++
}

func (pc *parseContext) now() uint8 {
	return pc.runeBuffer[1]
}

func (pc *parseContext) following() uint8 {
	return pc.runeBuffer[2]
}

func (pc *parseContext) before() uint8 {
	return pc.runeBuffer[0]
}

func (pc *parseContext) reset(r io.Reader, w io.Writer) {
	pc.r = r
	pc.w = w
	pc.generalBuffer = pc.generalBuffer[:0]
	pc.skipWrite = false
	pc.end = false
	pc.i = 0
}

type queryPath string

func (q queryPath) kv() (string, string) {
	return splitStringOnEqual(string(q))
}

var nonClosingHeaderTags = [][]byte{[]byte("meta"), []byte("link")}

func isNonClosingTag(pc *parseContext) bool {
	for _, tag := range nonClosingHeaderTags {
		if bytes.Equal(pc.generalBuffer, tag) {
			return true
		}
	}
	return false
}

// given a ctx where 'now' is the comment tag opener
// it will continue until the end of the comment tag
func untilCommentEnd(pc *parseContext) {
	for ; !pc.end; pc.next() {
		if pc.now() == '>' && pc.before() == '-' {
			break
		}
	}
}

func untilCurrentTagCloseTagStart(pc *parseContext) {
	depth := 0
	for ; !pc.end; pc.next() {
		untilNextOpen(pc)

		// skip over comment tags
		if pc.following() == '!' {
			untilCommentEnd(pc)
			continue
		}

		// make sure were in an opening tag
		if pc.following() != '/' {
			pc.resetGeneralBuffer()
			for ; !pc.end; pc.next() {
				if pc.now() == ' ' || pc.now() == '>' {
					break
				}
				if pc.now() == '<' {
					continue
				}
				pc.generalBuffer = append(pc.generalBuffer, pc.now())
			}

			if isNonClosingTag(pc) {
				// skip over non-closing tags
				continue
			}

			if shouldTagContentBeSkipped(pc) {
				// skip over content first
				untilNextOpen(pc)
				untilNextEnd(pc)
				continue
			}
		}

		if pc.following() == '/' {
			if depth == 0 {
				// found closing tag
				return
			}
			depth--
			continue
		}

		// skip over to end
		for ; !pc.end; pc.next() {
			if pc.now() != '>' {
				continue
			}
			break
		}

		if pc.before() == '/' {
			continue
		}

		depth++
	}
}

func (q queryPath) Type() string {
	key, _ := splitStringOnEqual(string(q))
	return key
}

func (q queryPath) Match(value []byte) bool {

	if len(value) == 0 {
		return false
	}

	qk, qv := q.kv()

	switch qk {
	case "tag":
		return qv == string(value)
	case "id", "class":
		_, val := splitBytesOnEqual(value)
		return qv == string(stripValueParentheses(val))
	default:
		return false
	}
}

// untilNextOpen skips until the next open tag
// and returns if its ok to continue
func untilNextOpen(pc *parseContext) {
	for ; !pc.end; pc.next() {
		if pc.now() != '<' {
			continue
		}
		return
	}
}

func untilHtmlTagOpen(pc *parseContext) {
	for ; !pc.end; pc.next() {
		if pc.now() == '<' && pc.following() == 'h' {
			break
		}
	}
}

func seekToEnd(pc *parseContext) {
	for ; !pc.end; pc.next() {
	}
}

func untilNextEnd(pc *parseContext) {
	for ; !pc.end; pc.next() {
		if pc.now() != '>' {
			continue
		}
		return
	}
}

var skippableTags = [][]byte{[]byte("script"), []byte("style"), []byte("noscript")}

func shouldTagContentBeSkipped(pc *parseContext) bool {
	for _, tag := range skippableTags {
		if bytes.Equal(pc.generalBuffer, tag) {
			return true
		}
	}
	return false
}

func seekMatchingTagEnd(pc *parseContext, path queryPath) {
	// start of tag
	if pc.now() != '<' {
		panic(errors.New("invalid element start"))
	}

	// skip over closing tag to next element
	if pc.following() == '/' {
		pc.next()
		untilNextOpen(pc)
		seekMatchingTagEnd(pc, path)
		return
	}

	// opening tag

	// iterate over tag name

	// clean out buffer
	pc.resetGeneralBuffer()

	for ; !pc.end; pc.next() {
		// tag closed with no info
		if pc.now() == '>' {

			// in case of no attributes we try to match the buffer here
			if path.Type() == "tag" && path.Match(pc.generalBuffer) {
				untilNextEnd(pc)
				return
			}

			untilNextOpen(pc)
			seekMatchingTagEnd(pc, path)
			return
		}
		// if not empty space skip
		if pc.now() != ' ' {
			// in case of a tag where we aren't currently
			// on the tag opener we copy it over to the general buffer
			if pc.now() != '<' && path.Type() == "tag" {
				pc.generalBuffer = append(pc.generalBuffer, pc.now())
			}
			continue
		}
		// break on first empty space
		break
	}

	// in case of attributes we try to match the buffer first
	if path.Type() == "tag" && path.Match(pc.generalBuffer) {
		untilNextEnd(pc)
		return
	}

	//TODO: investigate why this is not used

	// check if tag content should be skipped
	//if shouldTagContentBeSkipped(pc) {
	//	untilNextEnd(pc)
	//	untilNextOpen(pc)
	//	seekMatchingTagEnd(pc, path)
	//	return
	//}

	// iterate over tag attributes
	for ; !pc.end; pc.next() {

		// arrived at tag end
		if pc.now() == '>' {
			untilNextOpen(pc)
			seekMatchingTagEnd(pc, path)
			return
		}

		// space - skip ahead
		if pc.now() == ' ' {
			continue
		}

		// attr found
		pc.resetGeneralBuffer()
		for ; !pc.end; pc.next() {

			pc.generalBuffer = append(pc.generalBuffer, pc.now())

			// in case the following rune is a space or a tag
			// closer we try to match the attr stored in the
			// general buffer.
			if pc.following() == ' ' || pc.following() == '>' {

				if path.Match(pc.attribute()) {
					untilNextEnd(pc)
					return
				}

				break
			}
		}
	}
}

func newParseCtx(r io.Reader, w io.Writer) *parseContext {
	return &parseContext{
		r:             r,
		w:             w,
		runeBuffer:    [3]byte{},
		generalBuffer: make([]byte, 0, 2048),
		end:           false,
		skipWrite:     false,
		i:             0,
	}
}

var closingTag = []byte("<")

func withCtx(r io.Reader, w io.Writer, f func(pc *parseContext) error) (err error) {
	pc := defaultPool.Get(r, w)
	defer defaultPool.Put(pc)

	defer func() {
		if r := recover(); r != nil {

			err = r.(error)

			// finish r/w operations on panic
			// in case were not at end of file
			// or context of read ended
			if !pc.end && err.Error() != "EOF" {
				seekToEnd(pc)
			}
		}
	}()

	return f(pc)
}

func Append(r io.Reader, w io.Writer, path, value string) error {

	if len(value) == 0 || value[0] != '<' {
		return errors.New("value must start with an html open tag '<'")
	}

	return withCtx(r, w, func(pc *parseContext) (err error) {
		untilHtmlTagOpen(pc)
		seekMatchingTagEnd(pc, queryPath(path))
		untilCurrentTagCloseTagStart(pc)

		// write the value while omitting the first tag opener
		//TODO: fix issue with tag opener written to 'w' because of how
		// untilCurrentTagCloseTagStart exit condition.
		if _, err = pc.w.Write(unsafeGetBytes(value[1:])); err != nil {
			return
		}
		if _, err = pc.w.Write(closingTag); err != nil {
			return
		}

		seekToEnd(pc)
		return
	})
}

func Set(r io.Reader, w io.Writer, path string, value string) error {
	return withCtx(r, w, func(pc *parseContext) (err error) {
		untilHtmlTagOpen(pc)
		seekMatchingTagEnd(pc, queryPath(path))

		if _, err = w.Write(unsafeGetBytes(value)); err != nil {
			return
		}

		pc.skipWrite = true
		untilCurrentTagCloseTagStart(pc)

		// we need to add the closing tag ourselves
		if _, err = pc.w.Write(closingTag); err != nil {
			return
		}

		pc.skipWrite = false
		seekToEnd(pc)

		return
	})
}
