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
		`<div class="grid xlarge-height">`,
		`<div class="s4">`,
		`<div class="s8">`,
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
	s4 := strings.Index(got, `<div class="s4">`)
	s8 := strings.Index(got, `<div class="s8">`)
	term := strings.Index(got, "<web-term")
	code := strings.Index(got, "<vs-code")

	if s4 > term || term > s8 || s8 > code {
		t.Errorf("got %q, want the s4 column then the term then the s8 column then vs-code", got)
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
