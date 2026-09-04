package deck_test

import (
	"testing"

	"github.com/dgageot/demoit/deck"
)

func TestSplitReadsFrontmatterAndBody(t *testing.T) {
	t.Parallel()

	src := []byte("---\nlayout: cover\n---\n# Un\n\n---\n# Deux\n")

	slides := deck.Split(src)

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if got, want := string(slides[0].Frontmatter), "layout: cover\n"; got != want {
		t.Errorf("got frontmatter %q, want %q", got, want)
	}
	if got, want := string(slides[0].Body), "# Un\n\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
	if got, want := slides[0].BodyLine, 4; got != want {
		t.Errorf("got body line %d, want %d", got, want)
	}
	if got, want := slides[1].StartLine, 7; got != want {
		t.Errorf("got start line %d, want %d", got, want)
	}
	if slides[1].Frontmatter != nil {
		t.Errorf("got frontmatter %q on a slide without one, want none", slides[1].Frontmatter)
	}
}

// The separator that ends a slide is also the opening fence of the next
// slide's frontmatter. This is the case the naive rule gets wrong, so it is
// the case that matters most.
func TestSplitReadsFrontmatterOnALaterSlide(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n---\nlayout: split\ntitle: Deux\n---\n::term{path=sources}\n\n---\n# Trois\n")

	slides := deck.Split(src)

	if len(slides) != 3 {
		t.Fatalf("got %d slides, want 3", len(slides))
	}
	if slides[0].Frontmatter != nil {
		t.Errorf("got frontmatter %q on the first slide, want none", slides[0].Frontmatter)
	}
	if got, want := string(slides[1].Frontmatter), "layout: split\ntitle: Deux\n"; got != want {
		t.Errorf("got frontmatter %q, want %q", got, want)
	}
	if got, want := string(slides[1].Body), "::term{path=sources}\n\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
	if got, want := slides[1].StartLine, 4; got != want {
		t.Errorf("got start line %d, want %d", got, want)
	}
	if got, want := slides[1].BodyLine, 7; got != want {
		t.Errorf("got body line %d, want %d", got, want)
	}
	if got, want := string(slides[2].Body), "# Trois\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
}

func TestSplitDoesNotMistakeProseForFrontmatter(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n---\nNote: attention, ceci est de la prose.\n---\n# Trois\n")

	slides := deck.Split(src)

	if len(slides) != 3 {
		t.Fatalf("got %d slides, want 3 — a capitalized key is prose, not frontmatter", len(slides))
	}
	if slides[1].Frontmatter != nil {
		t.Errorf("got frontmatter %q, want the line kept as body", slides[1].Frontmatter)
	}
	if got, want := string(slides[1].Body), "Note: attention, ceci est de la prose.\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
}

// Documented limitation: a body paragraph that starts with a lowercase key at
// column 0 and is followed by a --- is read as frontmatter. Pinned here so the
// behaviour is a known trade-off rather than a surprise.
func TestSplitTakesALowercaseKeyAsFrontmatter(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n---\nnote: ceci sera pris pour du frontmatter\n---\n# Trois\n")

	slides := deck.Split(src)

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if got, want := string(slides[1].Frontmatter), "note: ceci sera pris pour du frontmatter\n"; got != want {
		t.Errorf("got frontmatter %q, want %q", got, want)
	}
}

func TestSplitIgnoresSeparatorsInsideAFence(t *testing.T) {
	t.Parallel()

	src := []byte("```mermaid\n---\ntitle: diagramme\n---\nflowchart LR\n```\n")

	slides := deck.Split(src)

	if len(slides) != 1 {
		t.Fatalf("got %d slides, want 1 — a --- inside a code fence must not cut", len(slides))
	}
}

func TestSplitIgnoresSeparatorsInsideAYAMLFence(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n```yaml\nfoo: 1\n---\nbar: 2\n```\n\n---\n# Deux\n")

	slides := deck.Split(src)

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
}

func TestSplitAcceptsLongerSeparators(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n-----\n# Deux\n")

	slides := deck.Split(src)

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
}

func TestSplitIgnoresSeparatorsInsideATildeFence(t *testing.T) {
	t.Parallel()

	src := []byte("# Un\n\n~~~\nfoo\n---\nbar\n~~~\n\n---\n# Deux\n")

	slides := deck.Split(src)

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2 — a --- inside a ~~~ fence must not cut", len(slides))
	}
}

// A deck whose first line is --- but whose frontmatter never closes must not
// gain a phantom empty first slide: the leading --- is consumed either way,
// because handing it back to the body scan makes it cut a slide of its own and
// shifts every following slide number by one.
func TestSplitHandlesUnterminatedFrontmatterAtFileStart(t *testing.T) {
	t.Parallel()

	src := []byte("---\nlayout: split\n# reste\n")

	slides := deck.Split(src)

	if len(slides) != 1 {
		t.Fatalf("got %d slides, want 1 — a leading --- must never produce an empty slide", len(slides))
	}
	if slides[0].Frontmatter != nil {
		t.Errorf("got frontmatter %q, want none — the block never closes", slides[0].Frontmatter)
	}
	if got, want := string(slides[0].Body), "layout: split\n# reste\n"; got != want {
		t.Errorf("got body %q, want %q", got, want)
	}
	if got, want := slides[0].BodyLine, 2; got != want {
		t.Errorf("got body line %d, want %d", got, want)
	}
}
