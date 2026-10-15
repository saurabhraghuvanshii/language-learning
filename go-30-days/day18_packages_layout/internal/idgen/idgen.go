// Package idgen lives under internal/. The Go toolchain refuses to compile any
// package OUTSIDE the tree rooted at go30/day18_packages_layout that imports it.
// So it is a private helper of this project, free to change without breaking users.
package idgen

import "fmt"

// Generator hands out sequential ids like "SKU-0001".
// The field is unexported (lower case): only code inside this package can touch it.
type Generator struct {
	prefix string
	n      int
}

// New is an exported constructor (upper case first letter = visible to importers).
func New(prefix string) *Generator {
	return &Generator{prefix: prefix}
}

// Next returns the next id. Not safe for concurrent use (add a mutex if needed).
func (g *Generator) Next() string {
	g.n++
	return fmt.Sprintf("%s-%04d", g.prefix, g.n)
}
