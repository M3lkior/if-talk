package deck_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dgageot/demoit/deck"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// renderFences converts src with the fence renderer installed.
func renderFences(t *testing.T, src string) string {
	t.Helper()

	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
			renderer.WithNodeRenderers(util.Prioritized(deck.NewFenceRenderer(), 99)),
		),
	)

	var out bytes.Buffer
	if err := md.Convert([]byte(src), &out); err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	return out.String()
}

func TestMermaidFenceBecomesAPreBlock(t *testing.T) {
	t.Parallel()

	got := renderFences(t, "```mermaid\nblock-beta\n    plugin1 -- \"output\" --> cpu_out\n```\n")

	if want := `<pre class="mermaid">`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
	if !strings.Contains(got, "block-beta") {
		t.Fatalf("got %q, want the diagram source kept", got)
	}
	if strings.Contains(got, "<code") {
		t.Fatalf("got %q, want no <code> element around a mermaid diagram", got)
	}
}

func TestMermaidFenceEscapesItsArrows(t *testing.T) {
	t.Parallel()

	got := renderFences(t, "```mermaid\nflowchart LR\n    a --> b\n```\n")

	if strings.Contains(got, "-->") {
		t.Fatalf("got %q, want the arrow escaped so it cannot close an HTML comment", got)
	}
	if !strings.Contains(got, "--&gt;") {
		t.Fatalf("got %q, want the arrow escaped as --&gt;", got)
	}
}

func TestOtherFencesAreHighlighted(t *testing.T) {
	t.Parallel()

	got := renderFences(t, "```yaml\nname: demoit\n```\n")

	if !strings.Contains(got, "demoit") {
		t.Fatalf("got %q, want the code kept", got)
	}
	if !strings.Contains(got, "style=") {
		t.Fatalf("got %q, want chroma's inline styles", got)
	}
}
