package html_overwrite

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"strings"
)

// removeNodeChildren removes all of a given
// node children.
func removeNodeChildren(node *html.Node) {
	children := make([]*html.Node, 0)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		children = append(children, child)
	}

	for _, child := range children {
		node.RemoveChild(child)
	}
}

// renderNode will convert the given HTML node
// to a string format.
func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return buf.String()
}

// isChildTextNode helps in cases
// in cases of text elements like
// <p>,<a>,<span>check to check if a given node
// child is a singular text node.
func isChildTextNode(n *html.Node) bool {

	// filter out nodes with no children
	if n.LastChild == nil || n.FirstChild == nil {
		return false
	}

	return n.FirstChild.Type == html.TextNode && n.LastChild.Type == html.TextNode
}

// cloneNode is a direct copy of html.Node internal
// copy implementation which is a way to export it
// into this codebase.
func cloneNode(n *html.Node) *html.Node {
	m := &html.Node{
		Type:     n.Type,
		DataAtom: n.DataAtom,
		Data:     n.Data,
		Attr:     make([]html.Attribute, len(n.Attr)),
	}
	copy(m.Attr, n.Attr)

	// in case the child of the node is a text
	// node we clone its children with it
	if isChildTextNode(n) {
		clonedTextNode := cloneNode(n.FirstChild)
		m.FirstChild = clonedTextNode
		m.LastChild = clonedTextNode
	}

	return m
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
