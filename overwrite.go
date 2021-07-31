package html_overwrite

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// writer is the underlying basic
// implementation of the Writer
// interface.
type writer struct {
	root *html.Node
}

// Writer allows mutating HTML nodes
// at ease using a simple query language
// and raw html string values.
type Writer interface {
	// Set will query for nodes matching the
	// given path and set their content to be the
	// given value.
	Set(path string, value string) error
	// Append will query for nodes matching the
	// given path and append a new child node
	// as the given value.
	Append(path string, value string) error
	// String will return the active HTML node
	// loaded into the writer in a string format.
	String() string
}

// a query will return a list of found nodes
// based on a given path
// PATH FORMAT:
// variables consist of passing either
// id/class an equals sign and their
// respective value split up by a comma.
// Example:
// class=name-of-class,id=3
func (w *writer) query(path string) (res []*html.Node) {
	res = make([]*html.Node, 0)
	items := strings.Split(path, ",")
	for _, item := range items {
		kv := strings.Split(item, "=")
		k, v := kv[0], kv[1]
		switch k {
		case "id":
			if n, err := id(w.root, v); err == nil {
				res = append(res, n)
				return
			}

		case "class":
			nodes := filter(w.root, func(n *html.Node) bool {
				for _, attr := range n.Attr {
					if attr.Key == "class" {
						return strings.Contains(attr.Val, v)
					}
				}
				return false
			})
			res = append(res, nodes...)
		}

	}
	return
}

// find the first node in the root tree matching the
// id criteria.
func id(root *html.Node, id string) (*html.Node, error) {
	var found *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		for _, attr := range node.Attr {
			// check if attr is "id" and if its value
			// matches our given id value
			if attr.Key == "id" && attr.Val == id {
				found = node
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(root)

	if found == nil {
		return nil, fmt.Errorf("failed to find id")
	}

	return found, nil
}

// filters a given root node & all of its children using a
// given function that returns a bool result.
func filter(root *html.Node, run func(n *html.Node) bool) []*html.Node {
	found := make([]*html.Node, 0)
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if run(node) {
			found = append(found, node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(root)
	return found
}

// parsePartial parses a given value as a partial HTML
// format string.
func parsePartial(value string) (*html.Node, error) {
	node, err := html.Parse(strings.NewReader(value))
	if err != nil {
		return nil, err
	}

	child := node.LastChild.LastChild.LastChild
	child.Parent = nil
	child.PrevSibling = nil
	child.NextSibling = nil

	return child, nil
}

// removes all of a given node children.
func removeNodeChildren(node *html.Node) {
	children := make([]*html.Node, 0)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, child)
	}

	for _, child := range children {
		node.RemoveChild(child)
	}
}

// Set will query for nodes matching the
// given path and set their content to be the
// given value.
func (w *writer) Set(path, value string) (err error) {
	nodes := w.query(path)

	// parse value as html node
	var newNode *html.Node
	if newNode, err = parsePartial(value); err != nil {
		return err
	}

	for _, node := range nodes {
		removeNodeChildren(node)
		node.AppendChild(newNode)
	}

	return
}

// Append will query for nodes matching the
// given path and append a new child node
// as the given value.
func (w *writer) Append(path, value string) (err error) {
	nodes := w.query(path)

	// parse value as html node
	var newNode *html.Node
	if newNode, err = parsePartial(value); err != nil {
		return err
	}

	for _, node := range nodes {

		// create new node
		node.AppendChild(newNode)
	}

	return
}

// renderNode will convert the given HTML node
// to a string format.
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return buf.String()
}

// String will return the active HTML node
// loaded into the writer in a string format.
func (w *writer) String() string {
	return renderNode(w.root)
}

// Load allows inputting html docs in string format
// into the writer so they could be modified.
func Load(rawHTML string) (w Writer, err error) {

	var doc *html.Node
	doc, err = html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return
	}

	w = &writer{doc}
	return
}
