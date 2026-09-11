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
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Theme says which parts of the rendering chassis a talk opts into. Both
// fields default to false, and that default matters: a talk that declares no
// theme block keeps the full-window rendering it has today and gets no theme
// switcher, which is what leaves an older talk untouched by the chassis.
type Theme struct {
	// Stage puts the slides on the fixed 1920x1080 stage, scaled to fit.
	Stage bool `yaml:"stage"`
	// Dark renders the theme switcher and lets the dark palette apply.
	Dark bool `yaml:"dark"`
}

// Logo is one logo and the two themes it has to survive. White is the image a
// light theme shows, Dark the one a dark theme shows.
//
// The pair is per logo rather than two parallel lists, and that is the point:
// a list of light logos beside a list of dark ones only works while both have
// the same length, so a partner with no dark variant had to be repeated in the
// dark list to keep its place. Here it declares `white` alone and the layouts
// show it in both themes -- which is also the old behaviour of a talk that
// declared no dark logos at all, now available one logo at a time.
type Logo struct {
	// White is the image a light theme shows. A logo with nothing else is
	// shown in both themes.
	White string `yaml:"white"`
	// Dark is the image a dark theme shows instead of White.
	Dark string `yaml:"dark"`
}

// declared reports whether the logo names any image at all. An undeclared
// eventLogo: block yields a zero Logo, which must render nothing rather than
// an <img> with an empty src.
func (l Logo) declared() bool {
	return l.White != "" || l.Dark != ""
}

// Talk is the identity of a presentation, shared by all its localized decks.
// It is read once from <folder>/talk.yml and handed to every layout.
type Talk struct {
	// Title is the talk's title, rendered from talk.yml through the same
	// Markdown pipeline as a slide's title: key -- so a cover can write
	// `*Impact Framework*` and have the word take the accent treatment, and an
	// apostrophe is not HTML-entity-escaped the way html/template escapes a
	// plain string.
	Title template.HTML `yaml:"title"`
	// Subtitle is the line a cover displays under the title, typically the
	// date and the speaker. Rendered as Markdown too.
	Subtitle template.HTML `yaml:"subtitle"`
	// Layout is the layout slides use when they declare none.
	Layout string `yaml:"layout"`
	// Logo is the speaker's own logo, the one every slide of the talk carries.
	Logo Logo `yaml:"logo"`
	// EventLogo is the logo of the event the talk is given at, shown beside
	// Logo on the cover and in every header. Optional, and the only reason it
	// is a field of its own rather than a second entry in a list: an event
	// changes from one conference to the next while the speaker's logo does
	// not, so a talk is re-pointed at a new event by filling one block.
	EventLogo Logo `yaml:"eventLogo"`
	// Footer is the text layouts may display at the bottom of a slide.
	Footer string `yaml:"footer"`
	// Theme is the rendering chassis the talk opts into.
	Theme Theme `yaml:"theme"`
}

// Logos returns the logos a layout displays, in order: the speaker's own
// first, then the event's when one is declared.
//
// Layouts range over this rather than testing the two fields themselves. They
// are the one part of the rendering a talk is free to replace, so every logo
// rule kept in a template is a rule a talk's override silently loses -- which
// is exactly how impact-framework's default-h3 layout ended up showing no dark
// logo at all. Here a layout writes one range and the decision stays in Go.
func (t Talk) Logos() []Logo {
	logos := make([]Logo, 0, 2)

	for _, logo := range []Logo{t.Logo, t.EventLogo} {
		if logo.declared() {
			logos = append(logos, logo)
		}
	}

	return logos
}

// LoadTalk reads <folder>/talk.yml. A missing file is not an error: it yields
// a Talk with the default layout and no logos, so that a talk can consist of a
// single Markdown file.
//
// At the folder's root rather than inside .demoit: it is the one file a
// speaker edits by hand for reasons that have nothing to do with the browser,
// and .demoit is otherwise the folder the server serves assets from.
func LoadTalk(folder string) (Talk, error) {
	talk := Talk{Layout: "default"}

	content, err := os.ReadFile(filepath.Join(folder, "talk.yml"))
	if errors.Is(err, fs.ErrNotExist) {
		return talk, nil
	}
	if err != nil {
		return talk, err
	}

	if err := yaml.Unmarshal(content, &talk); err != nil {
		return talk, fmt.Errorf("unable to parse talk.yml: %w", err)
	}

	// Rendered after the parse rather than during it: yaml puts the raw
	// Markdown in the field -- template.HTML has string as its underlying type
	// -- and we replace it with its rendering. renderTitle returns the empty
	// string untouched, so a talk that declares neither key keeps two empty
	// fields a layout can test with {{ with }}.
	title, _, err := renderTitle(string(talk.Title), 0)
	if err != nil {
		return talk, fmt.Errorf("unable to render the talk title: %w", err)
	}
	talk.Title = title

	subtitle, _, err := renderTitle(string(talk.Subtitle), 0)
	if err != nil {
		return talk, fmt.Errorf("unable to render the talk subtitle: %w", err)
	}
	talk.Subtitle = subtitle

	if talk.Layout == "" {
		talk.Layout = "default"
	}

	return talk, nil
}
