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
	"net/http"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/dgageot/demoit/files"
	"github.com/dgageot/demoit/highlight"
)

// Code returns the content of a source file.
func Code(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/sourceCode/")

	var contents []byte

	if !files.Exists(filename) {
		http.NotFound(w, r)
		return
	}
	var err error
	contents, err = files.Read(filename)
	if err != nil {
		http.Error(w, "Unable to read "+filename, http.StatusInternalServerError)
		return
	}

	lexer := highlight.ForFile(filename)
	style := highlight.Style(r.FormValue("style"))
	lines := highligtedLines(r)
	formatter := html.New(html.Standalone(true), html.WithLineNumbers(true), html.HighlightLines(lines), html.WithClasses(true))

	iterator, err := lexer.Tokenise(nil, string(contents))
	if err != nil {
		http.Error(w, "Unable to tokenize "+filename, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := formatter.Format(w, style, iterator); err != nil {
		http.Error(w, "Unable to format source code", http.StatusInternalServerError)
		return
	}
}

func highligtedLines(r *http.Request) [][2]int {
	lines := [][2]int{}

	startLines := strings.Split(r.FormValue("startLine"), ",")
	endLines := strings.Split(r.FormValue("endLine"), ",")

	for i := range startLines {
		startLine, _ := strconv.Atoi(startLines[i])
		endLine, _ := strconv.Atoi(endLines[i])

		lines = append(lines, [2]int{startLine, endLine})
	}

	return lines
}
