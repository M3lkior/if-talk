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
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed layouts/*.html
var embeddedLayouts embed.FS

// partialsName is the layout file that only defines shared fragments.
const partialsName = "partials"

// Slide is one slide of a deck, ready to be handed to a layout.
type Slide struct {
	// Talk is the identity of the presentation the slide belongs to.
	Talk Talk
	// Title is the slide's title, displayed by the header of most layouts. It
	// is rendered from the frontmatter title through the same Markdown pipeline
	// as the slide body, so an author can write **bold** or other inline
	// Markdown in a title and an apostrophe is not HTML-entity-escaped.
	Title template.HTML
	// Source is the reference displayed at the bottom of the slide.
	Source string
	// Class holds the CSS classes for the slide's main element. A slide that
	// declares none gets its layout's default (deck.defaultClasses); a slide
	// that declares a class: key replaces the default wholesale rather than
	// adding to it, because real slides need to drop a class as often as add
	// one — several carry no center-align at all.
	Class string
	// Height is the height class of the split layout's column container.
	Height string
	// Layout is the name of the layout the slide is rendered with.
	Layout string
	// Content is the slide's rendered Markdown body.
	Content template.HTML
	// Notes is the rendered speakernotes key of the slide's frontmatter.
	Notes template.HTML
	// Meta holds every frontmatter key the engine does not know about.
	Meta map[string]any
	// StartLine is the 1-based line of the slide in the deck file.
	StartLine int
}

// Layouts resolves and executes slide layouts. A layout the presentation
// declares in .demoit/layouts always wins over the one embedded in the binary,
// and a presentation may add layouts the engine knows nothing about.
type Layouts struct {
	folder string
}

// NewLayouts returns the layouts of the presentation in folder.
func NewLayouts(folder string) *Layouts {
	return &Layouts{folder: folder}
}

// Execute renders the slide through the layout with the given name.
func (l *Layouts) Execute(name string, slide Slide) (template.HTML, error) {
	partials, err := l.read(partialsName)
	if err != nil {
		return "", err
	}

	body, err := l.read(name)
	if err != nil {
		return "", err
	}

	parsed, err := template.New(name).Parse(string(partials))
	if err != nil {
		return "", fmt.Errorf("unable to parse the layout partials: %w", err)
	}

	parsed, err = parsed.Parse(string(body))
	if err != nil {
		return "", fmt.Errorf("unable to parse the %q layout: %w", name, err)
	}

	var out bytes.Buffer
	if err := parsed.Execute(&out, slide); err != nil {
		return "", fmt.Errorf("unable to render the %q layout: %w", name, err)
	}

	return template.HTML(out.String()), nil
}

// read returns the source of a layout, preferring the presentation's own copy.
func (l *Layouts) read(name string) ([]byte, error) {
	if !isLayoutName(name) {
		return nil, fmt.Errorf("invalid layout name %q: a layout name is a bare word, with no path separator", name)
	}

	own, err := os.ReadFile(filepath.Join(l.folder, ".demoit", "layouts", name+".html"))
	if err == nil {
		return own, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	embedded, err := embeddedLayouts.ReadFile("layouts/" + name + ".html")
	if errors.Is(err, fs.ErrNotExist) {
		if name == partialsName {
			return nil, nil
		}

		return nil, fmt.Errorf("unknown layout %q: it is neither in .demoit/layouts nor embedded in demoit", name)
	}

	return embedded, err
}

// isLayoutName reports whether name is a bare layout name. The name arrives
// from a slide's frontmatter and is joined into a file path, so it must not be
// able to walk out of the layouts folder: without this, a deck could name
// `../../../../etc/passwd` as its layout and have the file parsed as a
// template and rendered into a slide. Custom layout names are unaffected —
// they are bare words like `demo-3-panneaux`.
func isLayoutName(name string) bool {
	if name == "" {
		return false
	}

	return !strings.ContainsAny(name, `/\`) && !strings.Contains(name, "..")
}
