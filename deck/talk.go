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

// Talk is the identity of a presentation, shared by all its localized decks.
// It is read once from <folder>/.demoit/talk.yml and handed to every layout.
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
	// Logos are the images the header and cover layouts display, in order.
	Logos []string `yaml:"logos"`
	// LogosDark are the dark-theme counterparts of Logos, in the same order.
	// A talk that declares none renders Logos in both themes, which is what
	// every talk did before this field existed -- the layouts emit the
	// conditional class only when this list is non-empty.
	LogosDark []string `yaml:"logosDark"`
	// Footer is the text layouts may display at the bottom of a slide.
	Footer string `yaml:"footer"`
	// Theme is the rendering chassis the talk opts into.
	Theme Theme `yaml:"theme"`
}

// LoadTalk reads <folder>/.demoit/talk.yml. A missing file is not an error: it
// yields a Talk with the default layout and no logos, so that a talk can
// consist of a single Markdown file.
func LoadTalk(folder string) (Talk, error) {
	talk := Talk{Layout: "default"}

	content, err := os.ReadFile(filepath.Join(folder, ".demoit", "talk.yml"))
	if errors.Is(err, fs.ErrNotExist) {
		return talk, nil
	}
	if err != nil {
		return talk, err
	}

	if err := yaml.Unmarshal(content, &talk); err != nil {
		return talk, fmt.Errorf("unable to parse .demoit/talk.yml: %w", err)
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
