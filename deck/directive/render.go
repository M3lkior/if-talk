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
	"fmt"
	"html"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// spec says how a directive name maps to HTML.
type spec struct {
	// tag is the custom element the directive renders as.
	tag string
	// container reports whether the element wraps children.
	container bool
	// attributes renders the element's attributes, leading space included.
	attributes func(attrs map[string]string) (string, error)
}

// specs is the catalogue of directives, keyed by the name written in Markdown.
var specs = map[string]spec{
	"split":        {tag: "split-view", container: true, attributes: splitAttributes},
	"window":       {tag: "fake-window", container: true, attributes: copyAttributes("title")},
	"speakernotes": {tag: "speaker-notes", container: true, attributes: copyAttributes()},
	"term":         {tag: "web-term", attributes: copyAttributes("path")},
	"browser":      {tag: "web-browser", attributes: copyAttributes("src")},
	"vscode":       {tag: "vs-code", attributes: copyAttributes("path")},
	"code":         {tag: "source-code", attributes: codeAttributes},
	"grid":         {tag: "div", container: true, attributes: gridAttributes},
	"col":          {tag: "div", container: true, attributes: gridAttributes},
}

// nodeRenderer renders directive nodes as demoit custom elements. It carries
// the parse context so that an attribute problem is reported as a slide-level
// error instead of failing the whole conversion.
type nodeRenderer struct {
	ctx parser.Context
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *nodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.render)
}

// render writes the element a directive maps to.
func (r *nodeRenderer) render(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	directiveNode, ok := node.(*Node)
	if !ok {
		return ast.WalkContinue, nil
	}

	element, known := specs[directiveNode.Name]
	if !known {
		// Unknown names are reported by the transformer; render nothing.
		return ast.WalkSkipChildren, nil
	}

	if !entering {
		if element.container {
			_, _ = w.WriteString("</" + element.tag + ">\n")
		}

		return ast.WalkContinue, nil
	}

	attributes, err := element.attributes(directiveNode.Attrs)
	if err != nil {
		addError(r.ctx, directiveNode.Line, "%s", err.Error())

		return ast.WalkSkipChildren, nil
	}

	_, _ = w.WriteString("<" + element.tag + attributes + ">")
	if !element.container {
		_, _ = w.WriteString("</" + element.tag + ">\n")
	}

	return ast.WalkContinue, nil
}

// copyAttributes copies the named attributes through, in order, skipping the
// ones that are absent or empty. Values are HTML-escaped: %q would quote them
// the Go way, which leaves a " in a value free to end the attribute and turn
// whatever follows into markup of its own.
func copyAttributes(keys ...string) func(map[string]string) (string, error) {
	return func(attrs map[string]string) (string, error) {
		rendered := ""
		for _, key := range keys {
			if value := attrs[key]; value != "" {
				rendered += fmt.Sprintf(` %s="%s"`, key, html.EscapeString(value))
			}
		}

		return rendered, nil
	}
}

// splitAttributes renders the class of a <split-view>, built from its height.
func splitAttributes(attrs map[string]string) (string, error) {
	if height := attrs["height"]; height != "" {
		return fmt.Sprintf(` class="%s"`, html.EscapeString(height+"-height")), nil
	}

	return "", nil
}
