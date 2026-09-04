package directive_test

import (
	"strings"
	"testing"
)

func TestCodeDirectiveTranslatesSeparators(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::code{folder=sources files=a.yml,b.yml lines=11-20,4-9}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	for _, want := range []string{
		`folder="sources"`,
		`files="a.yml b.yml"`,
		`start-lines="11;4"`,
		`end-lines="20;9"`,
		`code_style="vs"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
}

func TestCodeDirectiveAlwaysEmitsLineAttributes(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::code{folder=sources files=pipelines.yml}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	for _, want := range []string{`start-lines=""`, `end-lines=""`} {
		if !strings.Contains(got, want) {
			t.Errorf("got %q, want it to contain %q — demoit.js splits these attributes unconditionally", got, want)
		}
	}
}

func TestCodeDirectiveRejectsMismatchedRanges(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "::code{folder=sources files=a.yml,b.yml lines=11-20}\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if !strings.Contains(errs[0].Message, "one range per file") {
		t.Errorf("got message %q, want it to explain the one-range-per-file rule", errs[0].Message)
	}
}

func TestCodeDirectiveNeedsFolderAndFiles(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "::code{files=a.yml}\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if !strings.Contains(errs[0].Message, "folder") {
		t.Errorf("got message %q, want it to name the missing folder attribute", errs[0].Message)
	}
}

func TestCodeDirectiveHonoursAnExplicitStyle(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::code{folder=sources files=a.yml style=monokai}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `code_style="monokai"`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestMalformedAttributeBlockIsReported(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "::term{path=sources\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1 — an unclosed brace must not pass silently", len(errs))
	}
	if !strings.Contains(errs[0].Message, "term") {
		t.Errorf("got message %q, want it to name the directive", errs[0].Message)
	}
}

func TestADirectiveWithoutAttributesIsNotReported(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::term\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none — no braces at all is legitimate", errs)
	}
	if want := "<web-term></web-term>"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

// A URL is the value goldmark's own attribute parser could not read: it stops
// at the first slash. Five of the deck's slides carry one.
func TestUnquotedURLValueSurvives(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::browser{src=https://if.greensoftware.foundation/users/quick-start}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `<web-browser src="https://if.greensoftware.foundation/users/quick-start">`; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestUnterminatedQuoteIsReported(t *testing.T) {
	t.Parallel()

	_, errs := render(t, "::term{path=\"sans fin}\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1 — an unterminated quote must not pass silently", len(errs))
	}
}

// A malformed block must render no attributes at all. Keeping the pairs read
// before the failure produced `<web-term path="sources&#10;">` next to the
// error: a component that looks configured and is not.
func TestMalformedAttributeBlockRendersNoAttributes(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::term{path=sources\n")

	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if want := "<web-term></web-term>"; !strings.Contains(got, want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestCodeDirectiveEscapesItsAttributes(t *testing.T) {
	t.Parallel()

	got, errs := render(t, "::code{folder=sources files=\"a&b.yml\"}\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if strings.Contains(got, "a&b.yml") {
		t.Fatalf("got %q, want the ampersand escaped", got)
	}
	if !strings.Contains(got, "a&amp;b.yml") {
		t.Fatalf("got %q, want it to contain a&amp;b.yml", got)
	}
}
