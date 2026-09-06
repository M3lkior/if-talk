package highlight_test

import (
	"strings"
	"testing"

	"github.com/dgageot/demoit/highlight"
)

// The lexer's identity is what these tests pin, not its existence: both
// functions end in `return lexers.Fallback`, so a nil check can never fail and
// leaves the YAML special-casing — and every code fence's highlighting —
// unpinned. Both were confirmed to fail when the function body is replaced
// with a bare `return lexers.Fallback`.
func TestForLanguageKnowsYAML(t *testing.T) {
	t.Parallel()

	if got := highlight.ForLanguage("yaml").Config().Name; got != "YAML" {
		t.Fatalf("got the %q lexer for yaml, want YAML", got)
	}
}

func TestForLanguageFallsBackOnUnknownLanguage(t *testing.T) {
	t.Parallel()

	if got := highlight.ForLanguage("not-a-language").Config().Name; got != "fallback" {
		t.Fatalf("got the %q lexer for an unknown language, want the fallback lexer", got)
	}
}

// ForFile is what /sourceCode calls on every request (handlers/code.go), and
// it had no test at all.
func TestForFileKnowsYAML(t *testing.T) {
	t.Parallel()

	if got := highlight.ForFile("x.yml").Config().Name; got != "YAML" {
		t.Fatalf("got the %q lexer for x.yml, want YAML", got)
	}
}

func TestForFileFallsBackOnAnUnknownExtension(t *testing.T) {
	t.Parallel()

	if got := highlight.ForFile("x.zzz").Config().Name; got != "fallback" {
		t.Fatalf("got the %q lexer for x.zzz, want the fallback lexer", got)
	}
}

func TestStyleDefaultsToGitHub(t *testing.T) {
	t.Parallel()

	if got := highlight.Style("").Name; got != "github" {
		t.Fatalf("got style %q, want github", got)
	}
}

func TestWriteHighlightsInline(t *testing.T) {
	t.Parallel()

	var out strings.Builder
	if err := highlight.Write(&out, "name: demoit\n", "yaml"); err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	got := out.String()
	if !strings.Contains(got, "demoit") {
		t.Fatalf("got %q, want it to contain the highlighted code", got)
	}
	if strings.Contains(got, "<html") {
		t.Fatalf("got %q, want inline HTML with no standalone document", got)
	}
	if strings.Contains(got, "class=") {
		t.Fatalf("got %q, want inline style attributes rather than CSS classes", got)
	}
}
