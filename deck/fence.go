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

package deck

import (
	"bytes"
	"html/template"

	"github.com/dgageot/demoit/highlight"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// mermaidLanguage is the fence language that renders as a mermaid diagram
// rather than as highlighted code.
const mermaidLanguage = "mermaid"

// fenceRenderer renders fenced code blocks: a mermaid fence becomes the
// <pre class="mermaid"> the talk's demoit.js looks for, and every other fence
// is highlighted with chroma.
type fenceRenderer struct{}

// NewFenceRenderer returns the node renderer for fenced code blocks.
func NewFenceRenderer() renderer.NodeRenderer {
	return &fenceRenderer{}
}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *fenceRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
}

// render writes a fenced code block.
func (r *fenceRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	block, ok := node.(*ast.FencedCodeBlock)
	if !ok {
		return ast.WalkContinue, nil
	}

	language := string(block.Language(source))
	code := blockText(block, source)

	if language == mermaidLanguage {
		_, _ = w.WriteString(`<pre class="mermaid">`)
		_, _ = w.WriteString(template.HTMLEscapeString(code))
		_, _ = w.WriteString("</pre>\n")

		return ast.WalkSkipChildren, nil
	}

	if err := highlight.Write(w, code, language); err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkSkipChildren, nil
}

// blockText returns the raw text of a fenced code block.
func blockText(block *ast.FencedCodeBlock, source []byte) string {
	var text bytes.Buffer

	lines := block.Lines()
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		text.Write(segment.Value(source))
	}

	return text.String()
}
