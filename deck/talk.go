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
	// Title is the talk's title, available to layouts that want a fallback
	// when a slide declares none.
	Title string `yaml:"title"`
	// Layout is the layout slides use when they declare none.
	Layout string `yaml:"layout"`
	// Logos are the images the header and cover layouts display, in order.
	Logos []string `yaml:"logos"`
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

	if talk.Layout == "" {
		talk.Layout = "default"
	}

	return talk, nil
}
