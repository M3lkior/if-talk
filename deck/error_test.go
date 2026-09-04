package deck_test

import (
	"strings"
	"testing"

	"github.com/dgageot/demoit/deck"
)

func TestABrokenSlideDoesNotBreakTheDeck(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "# Un\n\n---\nlayout: pas-un-layout\n---\n# Deux\n\n---\n# Trois\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none — a slide error must not fail the whole deck", err)
	}

	if len(slides) != 3 {
		t.Fatalf("got %d slides, want 3 — a broken slide still counts", len(slides))
	}
	if !strings.Contains(string(slides[0]), "<h1>Un</h1>") {
		t.Errorf("got %q, want the first slide rendered normally", slides[0])
	}
	if !strings.Contains(string(slides[2]), "<h1>Trois</h1>") {
		t.Errorf("got %q, want the last slide rendered normally", slides[2])
	}
	if !strings.Contains(string(slides[1]), "pas-un-layout") {
		t.Errorf("got %q, want the broken slide to name the unknown layout", slides[1])
	}
	if !strings.Contains(string(slides[1]), "demoit.md:4") {
		t.Errorf("got %q, want the error to point at demoit.md line 4, the slide's first line", slides[1])
	}
	if strings.Contains(string(slides[1]), "<h1>Deux</h1>") {
		t.Errorf("got %q, want the error block to replace the slide's body, not sit beside it", slides[1])
	}
}

func TestAnUnclosedDirectiveIsReportedOnItsOwnSlide(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "# Un\n\n---\n# Deux\n\n:::split\n::term{path=sources}\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if !strings.Contains(string(slides[0]), "<h1>Un</h1>") {
		t.Errorf("got %q, want the first slide untouched by the second's failure", slides[0])
	}
	if !strings.Contains(string(slides[1]), "never closed") {
		t.Errorf("got %q, want it to report the unclosed directive", slides[1])
	}
	if !strings.Contains(string(slides[1]), "demoit.md:6") {
		t.Errorf("got %q, want the error to point at demoit.md line 6, where the split opens", slides[1])
	}
}

func TestInvalidFrontmatterIsReportedOnItsOwnSlide(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "---\nlogos: [unterminated\n---\n# Un\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if len(slides) != 1 {
		t.Fatalf("got %d slides, want 1", len(slides))
	}
	if !strings.Contains(string(slides[0]), "demoit.md:1") {
		t.Errorf("got %q, want the error to point at demoit.md line 1", slides[0])
	}
	// A line number alone is not an actionable error: the block must also carry
	// what went wrong. Without this, a scrambled or generic message still
	// passes. "yaml:" is the prefix yaml.v3 puts on its own errors.
	if !strings.Contains(string(slides[0]), "yaml:") {
		t.Errorf("got %q, want the block to say what the YAML problem was, not just where", slides[0])
	}
}
