package directive_test

import (
	"strings"
	"testing"
)

func TestSplitWithColsRendersAWeightedGrid(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::split{cols=4,8 height=xlarge}\n::term{path=sources}\n\n::vscode{path=sources}\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	for _, want := range []string{
		`<div class="grid grid-cols-12 gap-4 h-stage-xlarge">`,
		`<div class="col-span-4">`,
		`<div class="col-span-8">`,
		`<web-term path="sources">`,
		`<vs-code path="sources">`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
	if strings.Contains(got, "<split-view") {
		t.Errorf("got %q, want a weighted grid rather than a split-view", got)
	}

	// Order matters and substring presence does not prove it. ReplaceChild
	// rewires the sibling links, so collecting the children in the wrong order
	// silently swaps the panes: the terminal would get the wide column and
	// VS Code the narrow one, with every assertion above still passing.
	narrow := strings.Index(got, `<div class="col-span-4">`)
	wide := strings.Index(got, `<div class="col-span-8">`)
	term := strings.Index(got, "<web-term")
	code := strings.Index(got, "<vs-code")

	if narrow > term || term > wide || wide > code {
		t.Errorf("got %q, want the col-span-4 column then the term then the col-span-8 column then vs-code", got)
	}
}

func TestSplitWithColsRejectsAChildCountMismatch(t *testing.T) {
	t.Parallel()

	_, errs := render(t, ":::split{cols=4,8}\n::term{path=sources}\n:::\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if !strings.Contains(errs[0].Message, "2 columns") {
		t.Errorf("got message %q, want it to compare columns and children", errs[0].Message)
	}
}

func TestUnknownDirectiveIsReported(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "::teminal{path=sources}\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if !strings.Contains(errs[0].Message, "teminal") {
		t.Errorf("got message %q, want it to name the unknown directive", errs[0].Message)
	}
}

// The class of a hand-written :::grid / :::col is escaped the same way every
// other directive attribute is: an unescaped " would end the attribute early
// and turn the rest of the class into markup of its own.
func TestGridClassIsHTMLEscaped(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::grid{class=\"a&b\"}\ncontenu\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if strings.Contains(got, `class="a&b"`) {
		t.Fatalf("got %q, want the ampersand escaped", got)
	}
	if want := `<div class="a&amp;b">`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

// A weight becomes a "col-span-N" utility verbatim, so anything that is not a
// whole number from 1 to 12 is a class no stylesheet defines.
func TestSplitWithColsRejectsANonNumericWeight(t *testing.T) {
	t.Parallel()

	for _, cols := range []string{"a,b", "13,4"} {
		cols := cols
		t.Run(cols, func(t *testing.T) {
			t.Parallel()

			got, errs := render(t, ":::split{cols="+cols+"}\n::term{path=a}\n\n::term{path=b}\n:::\n")

			if len(errs) != 1 {
				t.Fatalf("got %d errors for cols=%s, want 1", len(errs), cols)
			}
			if !strings.Contains(errs[0].Message, "12") {
				t.Errorf("got message %q, want it to name the twelve-column grid", errs[0].Message)
			}
			if strings.Contains(got, `class="col-span-`) {
				t.Errorf("got %q, want no col-span class built from a weight that is not one", got)
			}
		})
	}
}
