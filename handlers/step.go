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

package handlers

import (
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dgageot/demoit/deck"
	"github.com/dgageot/demoit/files"
	"github.com/dgageot/demoit/flags"
	"github.com/gorilla/mux"
)

//go:embed resources/index.tmpl.html
var indexHTML string
var indexTemplate = template.Must(template.New("index").Funcs(templateFuncs).Parse(indexHTML))

// templateFuncs is shared by every page template. `hash` busts the cache of a
// file served out of the presentation folder; `engineHash` does the same for
// the stylesheet embedded in the binary, which `hash` cannot see.
var templateFuncs = template.FuncMap{
	"hash":       hash,
	"engineHash": EngineCSSHash,
}

// Page describes a page of the demo.
type Page struct {
	WorkingDir  string
	HTML        template.HTML
	URL         string
	PrevURL     string
	NextURL     string
	CurrentStep int
	StepCount   int
	DevMode     bool
	// Stage and Dark come from the talk's talk.yml theme block. They
	// gate the fixed stage and the theme switcher, so a talk that declares
	// neither renders exactly as it did before the chassis existed.
	Stage bool
	Dark  bool
}

// Step renders a given page.
func Step(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	steps, err := readSteps(files.Root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read steps: %v", err), http.StatusInternalServerError)
		return
	}

	id := 0
	if vars["id"] != "" {
		id, err = strconv.Atoi(vars["id"])
		if err != nil || id >= len(steps) {
			http.NotFound(w, r)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html")
	if err := indexTemplate.Execute(w, steps[id]); err != nil {
		http.Error(w, "Unable to render page", http.StatusInternalServerError)
		return
	}
}

// LastStep redirects to the latest page.
func LastStep(w http.ResponseWriter, r *http.Request) {
	steps, err := readSteps(files.Root)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read steps: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%d", len(steps)-1), http.StatusSeeOther)
}

func readSteps(folder string) ([]Page, error) {
	locale := ""
	if flags.Locale != nil {
		locale = *flags.Locale
	}

	rendered, err := deck.Load(folder, locale)
	if err != nil {
		return nil, err
	}

	// The talk is read again here rather than threaded out of deck.Load, which
	// returns rendered slides only. It is the same file deck.Load already read,
	// and a missing talk.yml is not an error.
	talk, err := deck.LoadTalk(folder)
	if err != nil {
		return nil, err
	}

	steps := make([]Page, 0, len(rendered))
	for i, html := range rendered {
		url := "/"
		if i > 0 {
			url = fmt.Sprintf("/%d", i)
		}

		steps = append(steps, Page{
			WorkingDir:  folder,
			HTML:        html,
			DevMode:     *flags.DevMode,
			CurrentStep: i,
			URL:         url,
			Stage:       talk.Theme.Stage,
			Dark:        talk.Theme.Dark,
		})
	}

	for i := range steps {
		steps[i].StepCount = len(steps) - 1
		if i > 0 {
			steps[i].PrevURL = steps[i-1].URL
		}
		if i < len(steps)-1 {
			steps[i].NextURL = steps[i+1].URL
		}
	}

	return steps, nil
}

// VerifyConfiguration runs a couple of verifications on the configuration.
func VerifyConfiguration() error {
	if _, err := readSteps(files.Root); err != nil {
		return err
	}

	info, err := os.Stat(filepath.Join(files.Root, ".demoit"))
	if os.IsNotExist(err) {
		return errors.New(`mandatory resource folder ".demoit" doesn't exist`)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return errors.New(`mandatory resource folder ".demoit" is not a folder`)
	}

	return nil
}

// Ignore errors and return empty string if an error occurs.
func hash(path string) string {
	h, err := files.Sha256(path)
	if err != nil {
		return ""
	}

	return h[:10]
}
