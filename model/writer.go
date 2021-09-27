package model

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
	// loaded into the stdLibWriter in a string format.
	String() string
}
