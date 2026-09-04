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

// Package highlight wraps chroma so that source files served by /sourceCode
// and Markdown code fences are highlighted by the same code.
package highlight

import (
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// nonDefaultYAMLLexer turns bare YAML scalars into single-quoted string
// tokens, so that values stand out from keys.
type nonDefaultYAMLLexer struct {
	chroma.Lexer
}

// Tokenise implements chroma.Lexer.
func (n *nonDefaultYAMLLexer) Tokenise(options *chroma.TokeniseOptions, text string) (chroma.Iterator, error) {
	iterator, err := n.Lexer.Tokenise(nil, text)
	if err != nil {
		return nil, err
	}

	updated := iterator.Tokens()

	for i, token := range updated {
		if token.Type == chroma.Text {
			if token.Value == "-" {
				continue
			}

			if i+1 >= len(updated) {
				continue
			}

			next := updated[i+1]
			if next.Type == chroma.Punctuation && next.Value == ":" {
				continue
			}

			token.Type = chroma.LiteralStringSingle
			updated[i] = token
		}
	}

	return chroma.Literator(updated...), nil
}

// ForFile returns the lexer to use for the given file name.
func ForFile(file string) chroma.Lexer {
	if lexer := lexers.Match(file); lexer != nil {
		if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
			fmt.Println("Using non default YAML Lexer")
			return &nonDefaultYAMLLexer{lexers.Get(".yaml")}
		}
		return lexer
	}

	return lexers.Fallback
}

// ForLanguage returns the lexer to use for the language of a Markdown code
// fence. An unknown language yields the fallback lexer.
func ForLanguage(language string) chroma.Lexer {
	if language == "yaml" || language == "yml" {
		if lexer := lexers.Get("yaml"); lexer != nil {
			return &nonDefaultYAMLLexer{lexer}
		}
	}

	if lexer := lexers.Get(language); lexer != nil {
		return lexer
	}

	return lexers.Fallback
}

// Style returns the chroma style with the given name, GitHub by default.
func Style(name string) *chroma.Style {
	if name != "" {
		style := styles.Get(name)
		if style != nil {
			return style
		}
	}

	return styles.GitHub
}

// Write writes code as HTML highlighted for the given language, with inline
// style attributes so that a slide needs no extra stylesheet.
func Write(w io.Writer, code, language string) error {
	iterator, err := ForLanguage(language).Tokenise(nil, code)
	if err != nil {
		return err
	}

	formatter := chromahtml.New(chromahtml.WithClasses(false))

	return formatter.Format(w, Style(""), iterator)
}
