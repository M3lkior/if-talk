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
	"errors"
	"fmt"
	"html"
	"strconv"
	"strings"
)

// codeAttributes turns the directive's single comma-separated form into the
// three separators <source-code> expects: spaces between file names,
// semicolons between per-file line ranges, and the code_style attribute the
// component reads. start-lines and end-lines are always emitted, even empty,
// because demoit.js splits them without checking that they exist.
func codeAttributes(attrs map[string]string) (string, error) {
	folder := attrs["folder"]
	if folder == "" {
		return "", errors.New(`the code directive needs a "folder" attribute`)
	}

	files := splitList(attrs["files"])
	if len(files) == 0 {
		return "", errors.New(`the code directive needs a "files" attribute`)
	}

	starts, ends, err := lineRanges(attrs["lines"], len(files))
	if err != nil {
		return "", err
	}

	style := attrs["style"]
	if style == "" {
		style = "vs"
	}

	return fmt.Sprintf(` folder="%s" files="%s" start-lines="%s" end-lines="%s" code_style="%s"`,
		html.EscapeString(folder),
		html.EscapeString(strings.Join(files, " ")),
		html.EscapeString(starts),
		html.EscapeString(ends),
		html.EscapeString(style)), nil
}

// lineRanges turns "11-20,4-9" into the "11;4" and "20;9" that <source-code>
// expects. There is one range per file, in the order the files are listed.
func lineRanges(spec string, fileCount int) (string, string, error) {
	if strings.TrimSpace(spec) == "" {
		return "", "", nil
	}

	ranges := splitList(spec)
	if len(ranges) != fileCount {
		return "", "", fmt.Errorf("the code directive got %d line ranges for %d files: write one range per file", len(ranges), fileCount)
	}

	first := make([]string, 0, len(ranges))
	last := make([]string, 0, len(ranges))

	for _, lineRange := range ranges {
		start, end, found := strings.Cut(lineRange, "-")
		if !found {
			return "", "", fmt.Errorf("the line range %q is not of the form start-end", lineRange)
		}

		// Both halves have to be numbers. <source-code> reads them straight
		// out of the attribute and highlights nothing when they are not, so
		// `lines=11-` and `lines=a-b` used to travel all the way to the slide
		// as dead markup — exactly the failure ::code exists to prevent by
		// owning the three-separator contract on the author's behalf.
		startLine, err := lineNumber(start, lineRange)
		if err != nil {
			return "", "", err
		}
		endLine, err := lineNumber(end, lineRange)
		if err != nil {
			return "", "", err
		}

		first = append(first, startLine)
		last = append(last, endLine)
	}

	return strings.Join(first, ";"), strings.Join(last, ";"), nil
}

// splitList splits a comma-separated attribute, dropping empty entries.
func splitList(value string) []string {
	var items []string

	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			items = append(items, trimmed)
		}
	}

	return items
}

// lineNumber trims one half of a line range and checks it really is a number,
// naming the whole range in the error so the author can find it on the line.
func lineNumber(half, lineRange string) (string, error) {
	trimmed := strings.TrimSpace(half)
	if _, err := strconv.Atoi(trimmed); err != nil {
		return "", fmt.Errorf("the line range %q is not of the form start-end: %q is not a line number", lineRange, trimmed)
	}

	return trimmed, nil
}
