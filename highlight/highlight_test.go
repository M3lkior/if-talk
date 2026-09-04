package highlight_test

import (
	"strings"
	"testing"

	"github.com/dgageot/demoit/highlight"
)

func TestForLanguageKnowsYAML(t *testing.T) {
	t.Parallel()

	if got := highlight.ForLanguage("yaml"); got == nil {
		t.Fatal("got no lexer for yaml, want one")
	}
}

func TestForLanguageFallsBackOnUnknownLanguage(t *testing.T) {
	t.Parallel()

	if got := highlight.ForLanguage("not-a-language"); got == nil {
		t.Fatal("got no lexer for an unknown language, want the fallback lexer")
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
