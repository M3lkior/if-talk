/*
Copyright 2018 Google LLC
Copyright 2022 David Gageot

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package directive adds block directives to goldmark: `:::name{attrs}` for a
// container and `::name{attrs}` for a leaf. Each one renders as the custom
// element the demoit web components define.
package directive

import "github.com/yuin/goldmark/ast"

// Kind is the NodeKind of Node.
var Kind = ast.NewNodeKind("Directive")

// Node is a block directive.
type Node struct {
	ast.BaseBlock

	// Name is the directive name, such as "split", "term" or "code".
	Name string
	// Attrs are the attributes written between braces.
	Attrs map[string]string
	// FenceLength is the number of colons of the opening fence.
	FenceLength int
	// Container reports whether the directive wraps children.
	Container bool
	// Line is the 1-based line the directive opened on, for error messages.
	Line int

	// closed records that a fence closed this container, as opposed to the end
	// of the file.
	closed bool
}

// Kind implements ast.Node.
func (n *Node) Kind() ast.NodeKind {
	return Kind
}

// Dump implements ast.Node.
func (n *Node) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Name": n.Name}, nil)
}

// NewNode returns a directive node.
func NewNode(name string, attrs map[string]string, fenceLength, line int) *Node {
	return &Node{
		Name:        name,
		Attrs:       attrs,
		FenceLength: fenceLength,
		Container:   fenceLength >= 3,
		Line:        line,
	}
}
