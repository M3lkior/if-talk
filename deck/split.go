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

// Package deck reads a presentation's deck file and renders its slides.
package deck

import "strings"

// RawSlide is one slide of a deck file, before rendering.
type RawSlide struct {
	// Frontmatter is the slide's YAML block, without its --- fences. It is nil
	// when the slide has none.
	Frontmatter []byte
	// Body is everything after the frontmatter.
	Body []byte
	// StartLine is the 1-based line of the slide's first line in the deck file.
	StartLine int
	// BodyLine is the 1-based line of Body's first line in the deck file.
	BodyLine int
}

// Split cuts a deck file into slides. A slide ends at a line of three or more
// dashes that sits outside any code fence. A --- inside a ``` or ~~~ fence
// never cuts, which is what lets a mermaid or YAML block contain one. That
// same separator is the opening fence of the next slide's frontmatter, so a
// deck never needs two --- in a row.
func Split(src []byte) []RawSlide {
	lines := splitLines(src)

	var slides []RawSlide
	for i, fileStart := 0, true; ; fileStart = false {
		slide, next, more := readSlide(lines, i, fileStart)
		slides = append(slides, slide)
		if !more {
			return slides
		}
		i = next
	}
}

// readSlide reads one slide starting at lines[start]. It returns the slide, the
// index of the next slide's first line, and whether a separator ended this one.
func readSlide(lines []string, start int, fileStart bool) (RawSlide, int, bool) {
	frontmatter, bodyStart := readFrontmatter(lines, start, fileStart)
	end, more := findSlideEnd(lines, bodyStart)

	return RawSlide{
		Frontmatter: frontmatter,
		Body:        []byte(strings.Join(lines[bodyStart:end], "")),
		StartLine:   start + 1,
		BodyLine:    bodyStart + 1,
	}, end + 1, more
}

// readFrontmatter reads the YAML block at the top of a slide and returns it
// with the index of the slide's first body line. The block's opening fence is
// the --- that separated this slide from the previous one, already consumed by
// the caller; for the first slide of the file it is a leading --- line. A
// block is frontmatter when its first line is a lowercase `key:` pair, which
// Markdown prose never is, and it ends at the next --- line. An unterminated
// block is not treated as frontmatter.
func readFrontmatter(lines []string, start int, fileStart bool) ([]byte, int) {
	i := skipBlank(lines, start)

	// At the start of the file a leading --- line opens the frontmatter. It is
	// consumed whatever follows: when the block turns out not to be
	// frontmatter, that line was a redundant leading separator, and handing it
	// back to findSlideEnd would cut an empty slide in front of the deck and
	// shift every slide number by one.
	fallback := start
	if fileStart && i < len(lines) && isSlideSeparator(lines[i]) {
		fallback = i + 1
		i = skipBlank(lines, i+1)
	}

	if i >= len(lines) || !isFrontmatterStart(lines[i]) {
		return nil, fallback
	}

	for j := i + 1; j < len(lines); j++ {
		if isSlideSeparator(lines[j]) {
			return []byte(strings.Join(lines[i:j], "")), j + 1
		}
	}

	return nil, fallback
}

// skipBlank returns the index of the first non-blank line at or after start.
func skipBlank(lines []string, start int) int {
	i := start
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	return i
}

// isFrontmatterStart reports whether the line opens a frontmatter block: a
// lowercase key in column 0, followed by a colon. Prose after a slide break
// starts with #, *, <, a capital or a blank line, so the --- that separates
// two slides stays unambiguous without probing the block as YAML.
func isFrontmatterStart(line string) bool {
	trimmed := strings.TrimRight(line, "\r\n")

	colon := strings.IndexByte(trimmed, ':')
	if colon < 1 {
		return false
	}

	for i := 0; i < colon; i++ {
		c := trimmed[i]
		lowercase := c >= 'a' && c <= 'z'
		if i == 0 && !lowercase {
			return false
		}
		if !lowercase && !(c >= '0' && c <= '9') && c != '_' && c != '-' {
			return false
		}
	}

	return true
}

// findSlideEnd returns the index of the separator that ends the slide, or
// len(lines) at end of file, and whether a separator was found.
func findSlideEnd(lines []string, start int) (int, bool) {
	fence := ""

	for i := start; i < len(lines); i++ {
		line := strings.TrimRight(lines[i], "\r\n")

		open := fenceOpen(line)

		switch {
		case fence != "":
			if isFenceClose(line, fence) {
				fence = ""
			}
		case open != "":
			fence = open
		case isSlideSeparator(line):
			return i, true
		}
	}

	return len(lines), false
}

// splitLines splits src into lines, each keeping its trailing newline.
func splitLines(src []byte) []string {
	rest := string(src)

	var lines []string
	for len(rest) > 0 {
		i := strings.IndexByte(rest, '\n')
		if i < 0 {
			lines = append(lines, rest)
			break
		}
		lines = append(lines, rest[:i+1])
		rest = rest[i+1:]
	}

	return lines
}

// fenceOpen returns the code fence the line opens, or an empty string when it
// opens none. A line indented by four spaces or more opens nothing.
func fenceOpen(line string) string {
	trimmed := strings.TrimLeft(line, " ")
	if len(line)-len(trimmed) > 3 {
		return ""
	}

	for _, char := range []byte{'`', '~'} {
		if run := runLength(trimmed, char); run >= 3 {
			return trimmed[:run]
		}
	}

	return ""
}

// isFenceClose reports whether the line closes the given code fence.
func isFenceClose(line, fence string) bool {
	trimmed := strings.TrimLeft(line, " ")
	run := runLength(trimmed, fence[0])

	return run >= len(fence) && strings.TrimSpace(trimmed[run:]) == ""
}

// runLength returns the number of leading char bytes of s.
func runLength(s string, char byte) int {
	i := 0
	for i < len(s) && s[i] == char {
		i++
	}

	return i
}

// isSlideSeparator reports whether the line separates two slides.
func isSlideSeparator(line string) bool {
	trimmed := strings.TrimSpace(line)

	return len(trimmed) >= 3 && strings.Trim(trimmed, "-") == ""
}
