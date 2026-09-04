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
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dgageot/demoit/deck/directive"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// defaultHeight is the height class the split layout uses when a slide
// declares none. Every split slide of the existing decks uses xlarge.
const defaultHeight = "xlarge"

// Load reads the deck of the presentation in folder and returns one chunk of
// HTML per slide. A Markdown deck goes through goldmark and its layouts; an
// HTML deck is split on --- and returned as it is, exactly as before.
func Load(folder, locale string) ([]template.HTML, error) {
	path, markdown, err := resolve(folder, locale)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if !markdown {
		return htmlSlides(content), nil
	}

	talk, err := LoadTalk(folder)
	if err != nil {
		return nil, err
	}

	return markdownSlides(content, talk, NewLayouts(folder), filepath.Base(path)), nil
}

// resolve finds the deck file of a presentation. A localized deck wins over
// the default one, and Markdown wins over HTML.
func resolve(folder, locale string) (string, bool, error) {
	type candidate struct {
		name     string
		markdown bool
	}

	candidates := make([]candidate, 0, 4)
	if locale != "" {
		candidates = append(candidates,
			candidate{fmt.Sprintf("demoit-%s.md", locale), true},
			candidate{fmt.Sprintf("demoit-%s.html", locale), false},
		)
	}
	candidates = append(candidates,
		candidate{"demoit.md", true},
		candidate{"demoit.html", false},
	)

	for _, option := range candidates {
		path := filepath.Join(folder, option.name)
		if _, err := os.Stat(path); err == nil {
			return path, option.markdown, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", false, err
		}
	}

	return "", false, fmt.Errorf("no deck found in %s: expected demoit.md or demoit.html", folder)
}

// htmlSlides splits a legacy HTML deck, keeping the historical behaviour.
func htmlSlides(content []byte) []template.HTML {
	parts := bytes.Split(content, []byte("---"))

	slides := make([]template.HTML, 0, len(parts))
	for _, part := range parts {
		slides = append(slides, template.HTML(part))
	}

	return slides
}

// markdownSlides renders every slide of a Markdown deck through its layout.
// It returns no error: renderSlide turns each per-slide failure into a visible
// error block so that one mistake costs one slide, never the whole deck.
func markdownSlides(content []byte, talk Talk, layouts *Layouts, file string) []template.HTML {
	raw := Split(content)

	slides := make([]template.HTML, 0, len(raw))
	for _, rawSlide := range raw {
		slides = append(slides, renderSlide(rawSlide, talk, layouts, file))
	}

	return slides
}

// renderSlide renders one slide, turning any problem into a visible error
// block so that the rest of the deck keeps working while writing.
func renderSlide(raw RawSlide, talk Talk, layouts *Layouts, file string) template.HTML {
	known, meta, err := parseFrontmatter(raw.Frontmatter)
	if err != nil {
		return errorHTML(file, raw.StartLine, err)
	}

	slide := Slide{
		Talk:      talk,
		Title:     known.Title,
		Source:    known.Source,
		Class:     known.Class,
		Height:    known.Height,
		Layout:    known.Layout,
		Meta:      meta,
		StartLine: raw.StartLine,
	}
	if slide.Layout == "" {
		slide.Layout = talk.Layout
	}
	if slide.Height == "" {
		slide.Height = defaultHeight
	}

	body, notes, line, err := renderMarkdown(raw, known.Speakernotes)
	if err != nil {
		return errorHTML(file, line, err)
	}
	slide.Content = body
	slide.Notes = notes

	rendered, err := layouts.Execute(slide.Layout, slide)
	if err != nil {
		return errorHTML(file, raw.StartLine, err)
	}

	return rendered
}

// renderMarkdown converts a slide's body and its speakernotes frontmatter key.
// It returns the deck-file line to report a failure on, so that the error block
// carries exactly one line number.
func renderMarkdown(raw RawSlide, notes string) (template.HTML, template.HTML, int, error) {
	body, line, err := convert(raw.Body, raw.BodyLine)
	if err != nil {
		return "", "", line, err
	}

	if strings.TrimSpace(notes) == "" {
		return body, "", 0, nil
	}

	// The speakernotes value sits inside the frontmatter block, so a line
	// counted from its own text points nowhere the author can use. Report the
	// slide's first line and name the key instead.
	rendered, _, err := convert([]byte(notes), raw.StartLine)
	if err != nil {
		return "", "", raw.StartLine, fmt.Errorf("in the speakernotes frontmatter key: %w", err)
	}

	return body, rendered, 0, nil
}

// convert renders Markdown to HTML. It returns the first directive problem
// with the line it sits on in the deck file, so that the caller reports one
// number rather than embedding a second one in the message.
func convert(source []byte, firstLine int) (template.HTML, int, error) {
	ctx := parser.NewContext()

	md := goldmark.New(
		goldmark.WithExtensions(directive.Extension{Context: ctx}),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			renderer.WithNodeRenderers(util.Prioritized(NewFenceRenderer(), 99)),
		),
	)

	var out bytes.Buffer
	if err := md.Convert(source, &out, parser.WithContext(ctx)); err != nil {
		return "", firstLine, err
	}

	if problems := directive.Errors(ctx); len(problems) > 0 {
		return "", firstLine + problems[0].Line - 1, errors.New(problems[0].Message)
	}

	return template.HTML(out.String()), 0, nil
}

// errorHTML renders a slide-level error so that a mistake costs one slide
// rather than the whole deck, which matters while writing with --dev.
func errorHTML(file string, line int, err error) template.HTML {
	return template.HTML(fmt.Sprintf(
		`<main class="responsive max middle-align center-align">`+
			`<article class="border round"><h5>%s:%d</h5><pre>%s</pre></article></main>`,
		template.HTMLEscapeString(file),
		line,
		template.HTMLEscapeString(err.Error()),
	))
}
