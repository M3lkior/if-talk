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
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// priority puts the directive parser and renderer ahead of every default one.
// goldmark's own block parsers start at 100 for setext headings, so 99 runs
// before all of them.
const priority = 99

// Extension adds `:::name{attrs}` containers and `::name{attrs}` leaves to a
// goldmark Markdown. Errors found while rendering are recorded into Context,
// which must be the context passed to Convert.
type Extension struct {
	// Context collects the errors found in directives. When nil, a fresh
	// context is used and errors are dropped.
	Context parser.Context
}

// Extend implements goldmark.Extender.
func (e Extension) Extend(md goldmark.Markdown) {
	ctx := e.Context
	if ctx == nil {
		ctx = parser.NewContext()
	}

	md.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewParser(), priority),
		),
		parser.WithASTTransformers(
			util.Prioritized(NewTransformer(ctx), priority),
		),
	)
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&nodeRenderer{ctx: ctx}, priority),
	))
}
