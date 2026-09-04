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

package directive

import (
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// transformer rewrites a split directive that carries column weights into a
// beercss grid, and reports directives whose name is not in the catalogue.
type transformer struct {
	ctx parser.Context
}

// NewTransformer returns the AST transformer of the directive extension.
func NewTransformer(ctx parser.Context) parser.ASTTransformer {
	return &transformer{ctx: ctx}
}

// Transform implements parser.ASTTransformer.
func (t *transformer) Transform(node *ast.Document, _ text.Reader, _ parser.Context) {
	var directives []*Node

	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if found, ok := n.(*Node); ok {
			directives = append(directives, found)
		}

		return ast.WalkContinue, nil
	})

	for _, found := range directives {
		if _, known := specs[found.Name]; !known {
			addError(t.ctx, found.Line, "unknown directive %q", found.Name)

			continue
		}
		if found.Name == "split" && found.Attrs["cols"] != "" {
			t.toGrid(found)
		}
	}
}

// toGrid turns a split node with column weights into a grid node whose direct
// children are each wrapped in a column node.
func (t *transformer) toGrid(node *Node) {
	cols := splitList(node.Attrs["cols"])

	children := make([]ast.Node, 0, node.ChildCount())
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		children = append(children, child)
	}

	if len(cols) != len(children) {
		addError(t.ctx, node.Line, "the split directive declares %d columns for %d children: write one column weight per child", len(cols), len(children))

		return
	}

	class := "grid"
	if height := node.Attrs["height"]; height != "" {
		class += " " + height + "-height"
	}

	node.Name = "grid"
	node.Attrs = map[string]string{"class": class}

	for i, child := range children {
		column := NewNode("col", map[string]string{"class": "s" + cols[i]}, node.FenceLength, node.Line)
		node.ReplaceChild(node, child, column)
		column.AppendChild(column, child)
	}
}

// gridAttributes renders the class of a grid or column node.
func gridAttributes(attrs map[string]string) (string, error) {
	if class := strings.TrimSpace(attrs["class"]); class != "" {
		return " class=\"" + class + "\"", nil
	}

	return "", nil
}
