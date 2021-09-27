package std

import (
	"fmt"
	"github.com/html-overwrite/model"
	"golang.org/x/net/html"
	"strings"
)

// writer is the underlying basic
// implementation of the Writer
// interface.
type writer struct {
	root *html.Node
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
		appendChild(node, newNode)
	}

	return
}

// appendChild is a replacement call for the html.Node
// AppendChild func due to it being pretty annoying in
// case of same instances of a child appended to different
// nodes (which causes parent collision) so this is an actual
// useful implementation of it that clones the child beforehand.
func appendChild(node *html.Node, child *html.Node) {
	// we have to copy html nodes to
	// prevent a collision of parent pointer
	// in the internal node logic
	clonedChild := cloneNode(child)
	node.AppendChild(clonedChild)
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
		appendChild(node, newNode)
	}

	return
}

// String will return the active HTML node
// loaded into the writer in a string format.
func (w *writer) String() string {
	return renderNode(w.root)
}

func NewWriter(rawHTML string) (model.Writer, error) {
	var doc *html.Node
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return nil, err
	}

	return &writer{doc}, nil
}
