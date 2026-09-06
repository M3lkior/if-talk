package deck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgageot/demoit/deck"
)

func TestLayoutsRenderTheEmbeddedDefault(t *testing.T) {
	t.Parallel()

	slide := deck.Slide{
		Talk:    deck.Talk{Logos: []string{"/images/a.svg", "/images/b.jpg"}},
		Title:   "L'impact du numérique",
		Source:  "https://www.arcep.fr/",
		Content: "<h2>2,5%</h2>",
		Notes:   "<span>ADEME</span>",
	}

	got, err := deck.NewLayouts(t.TempDir()).Execute("default", slide)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	for _, want := range []string{
		`<h2 class="max center-left">L'impact du numérique</h2>`,
		`src="/images/a.svg"`,
		`src="/images/b.jpg"`,
		"<h2>2,5%</h2>",
		"Source: https://www.arcep.fr/",
		"<speaker-notes><span>ADEME</span></speaker-notes>",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("got %q, want it to contain %q", got, want)
		}
	}
}

func TestLayoutsOmitAnEmptySourceAndNotes(t *testing.T) {
	t.Parallel()

	got, err := deck.NewLayouts(t.TempDir()).Execute("default", deck.Slide{Content: "<p>rien</p>"})
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if strings.Contains(string(got), "Source:") {
		t.Errorf("got %q, want no source line when the slide declares none", got)
	}
	if strings.Contains(string(got), "<speaker-notes>") {
		t.Errorf("got %q, want no speaker-notes element when the slide declares none", got)
	}
}

func TestSplitLayoutWrapsTheContent(t *testing.T) {
	t.Parallel()

	got, err := deck.NewLayouts(t.TempDir()).Execute("split", deck.Slide{Height: "xlarge", Content: "<web-term></web-term>"})
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if want := `<split-view class="xlarge-height">`; !strings.Contains(string(got), want) {
		t.Fatalf("got %q, want it to contain %q", got, want)
	}
}

func TestEveryEmbeddedLayoutIsOverridable(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"cover", "default", "quote", "split", "content", "bare"} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			folder := t.TempDir()
			dir := filepath.Join(folder, ".demoit", "layouts")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatalf("unable to create the layouts folder: %v", err)
			}
			body := "<article>surcharge de " + name + "</article>"
			if err := os.WriteFile(filepath.Join(dir, name+".html"), []byte(body), 0o600); err != nil {
				t.Fatalf("unable to write the layout: %v", err)
			}

			got, err := deck.NewLayouts(folder).Execute(name, deck.Slide{})
			if err != nil {
				t.Fatalf("got error %v, want none", err)
			}
			if string(got) != body {
				t.Fatalf("got %q, want the talk's own layout %q", got, body)
			}
		})
	}
}

func TestLayoutsAcceptATalkSpecificName(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	dir := filepath.Join(folder, ".demoit", "layouts")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("unable to create the layouts folder: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "demo-3-panneaux.html"), []byte("<article>trois</article>"), 0o600); err != nil {
		t.Fatalf("unable to write the layout: %v", err)
	}

	got, err := deck.NewLayouts(folder).Execute("demo-3-panneaux", deck.Slide{})
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}
	if want := "<article>trois</article>"; string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLayoutsRejectAnUnknownName(t *testing.T) {
	t.Parallel()

	_, err := deck.NewLayouts(t.TempDir()).Execute("pas-un-layout", deck.Slide{})
	if err == nil {
		t.Fatal("got no error, want one naming the unknown layout")
	}
	if !strings.Contains(err.Error(), "pas-un-layout") {
		t.Fatalf("got error %v, want it to name the unknown layout", err)
	}
}

// A layout name comes from a slide's frontmatter and is joined into a file
// path, so it must not be able to walk out of the layouts folder and have an
// arbitrary file parsed as a template and rendered into the slide.
//
// The traversal target has to be a file that actually exists, otherwise the
// test passes with or without the guard: a name that resolves to nothing falls
// through to the embedded lookup and fails there as an unknown layout, which
// is an error too. So the test plants a file outside the talk folder and
// checks it stays unreachable.
func TestLayoutsRejectAPathInTheName(t *testing.T) {
	t.Parallel()

	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "secret.html"), []byte("<p>fuite</p>"), 0o600); err != nil {
		t.Fatalf("unable to write the file to protect: %v", err)
	}

	folder := t.TempDir()
	escape, err := filepath.Rel(
		filepath.Join(folder, ".demoit", "layouts"),
		filepath.Join(outside, "secret"),
	)
	if err != nil {
		t.Fatalf("unable to build the traversal path: %v", err)
	}

	for _, name := range []string{escape, "foo/bar", `foo\bar`, ".."} {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := deck.NewLayouts(folder).Execute(name, deck.Slide{})
			if err == nil {
				t.Fatalf("got %q with no error for the layout name %q, want an error", got, name)
			}
			if strings.Contains(string(got), "fuite") {
				t.Fatalf("got %q, want the file outside the talk folder to stay unreachable", got)
			}
		})
	}
}

// partials.html defines the shared header/source/notes blocks and nothing
// else, so naming it as a layout rendered an empty slide with no diagnostic.
func TestLayoutsRejectThePartialsName(t *testing.T) {
	t.Parallel()

	got, err := deck.NewLayouts(t.TempDir()).Execute("partials", deck.Slide{Content: "<p>contenu</p>"})
	if err == nil {
		t.Fatalf("got %q with no error, want one saying partials is not a layout", got)
	}
	if !strings.Contains(err.Error(), "partials") {
		t.Fatalf("got error %v, want it to name partials", err)
	}
}
