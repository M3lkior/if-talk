package directive_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dgageot/demoit/deck/directive"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// render converts src and returns the HTML plus the directive errors found.
func render(t *testing.T, src string) (string, []directive.Error) {
	t.Helper()

	ctx := parser.NewContext()

	md := goldmark.New(
		goldmark.WithExtensions(directive.Extension{Context: ctx}),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	var out bytes.Buffer
	if err := md.Convert([]byte(src), &out, parser.WithContext(ctx)); err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	return out.String(), directive.Errors(ctx)
}

func TestLeafDirectiveRendersASelfContainedElement(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::term{path=sandbox}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `<web-term path="sandbox"></web-term>`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestContainerDirectiveWrapsItsMarkdownChildren(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::speakernotes\nmontrer **if-run**\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if !strings.Contains(got, "<speaker-notes>") || !strings.Contains(got, "</speaker-notes>") {
		t.Fatalf("got %q, want it wrapped in <speaker-notes>", got)
	}
	if !strings.Contains(got, "<strong>if-run</strong>") {
		t.Fatalf("got %q, want the children parsed as Markdown", got)
	}
}

func TestContainersNestWithEqualFenceLengths(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::window{title=Terminal}\n:::split\n::term{path=sources}\n:::\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	for _, want := range []string{`<fake-window title="Terminal">`, "<split-view>", `<web-term path="sources">`, "</split-view>", "</fake-window>"} {
		if !strings.Contains(got, want) {
			t.Fatalf("got %q, want it to contain %q", got, want)
		}
	}
	if strings.Index(got, "<split-view>") < strings.Index(got, "<fake-window") {
		t.Fatal("got the split before the window, want it nested inside")
	}
}

func TestIndentedContainerContentIsNotACodeBlock(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::window{title=Demo}\n  :::split\n    ::term{path=sources}\n  :::\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if strings.Contains(got, "<pre><code>") {
		t.Fatalf("got %q, want the indented content parsed, not turned into a code block", got)
	}
	if !strings.Contains(got, `<web-term path="sources">`) {
		t.Fatalf("got %q, want the nested term rendered", got)
	}
}

func TestUnclosedContainerIsReported(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "# Titre\n\n:::split\n::term{path=sources}\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if got, want := errs[0].Line, 3; got != want {
		t.Errorf("got line %d, want %d — the line the container opened on", got, want)
	}
	if !strings.Contains(errs[0].Message, "split") {
		t.Errorf("got message %q, want it to name the split directive", errs[0].Message)
	}
}

func TestSplitCarriesItsHeightAsAClass(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::split{height=xlarge}\n::term{path=sources}\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `<split-view class="xlarge-height">`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestTwoColonsIsNotADirectiveWithoutAName(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "prose :: prose\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := "<p>prose :: prose</p>"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q — colons in prose stay prose", got, want)
	}
}

// The parser is registered under the ':' trigger, so the test above never
// reaches Open. This one does, and it is what exercises the empty-name guard.
func TestAFenceWithoutANameIsNotADirective(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::\ncontenu\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if !strings.Contains(got, "contenu") {
		t.Fatalf("got %q, want the content kept — a fence with no name is not a directive", got)
	}
}

// contentIndent is a column width and Advance moves bytes, so without a clamp
// to the line's own whitespace a less-indented continuation line loses
// characters off its front.
func TestOutdentedContinuationLineKeepsItsText(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::speakernotes\n  première ligne\ndeuxième ligne\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if !strings.Contains(got, "deuxième ligne") {
		t.Fatalf("got %q, want the outdented line kept whole", got)
	}
}

// A tab is one byte but four columns wide, so the same missing clamp ate four
// bytes off a tab-indented line and destroyed the directive on it.
func TestTabIndentedDirectiveSurvives(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::window{title=Demo}\n\t::term{path=sources}\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `<web-term path="sources">`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestAttributeValuesAreHTMLEscaped(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::window{title=\"Démo <b>&</b>\"}\ncontenu\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if strings.Contains(got, "<b>") {
		t.Fatalf("got %q, want < and & escaped inside the attribute value", got)
	}
	if !strings.Contains(got, "&lt;b&gt;") {
		t.Fatalf("got %q, want the value escaped as &lt;b&gt;", got)
	}
}
