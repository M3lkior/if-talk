package deck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dgageot/demoit/deck"
)

// writeDeck writes files into a fresh presentation folder and returns it.
func writeDeck(t *testing.T, files map[string]string) string {
	t.Helper()

	folder := t.TempDir()
	for name, content := range files {
		path := filepath.Join(folder, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("unable to create %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("unable to write %s: %v", name, err)
		}
	}

	return folder
}

func TestLoadRendersAMarkdownDeck(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		".demoit/talk.yml": "logos: [/images/a.svg]\n",
		"demoit.md":        "---\nlayout: cover\n---\n# Impact Framework\n\n---\ntitle: Les chiffres\nsource: https://arcep.fr\n---\n**2,5%** des émissions.\n\n:::speakernotes\nADEME\n:::\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if want := "<h1>Impact Framework</h1>"; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want it to contain %q", slides[0], want)
	}
	if want := `src="/images/a.svg"`; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want the cover to show the talk's logo", slides[0])
	}
	for _, want := range []string{
		`<h2 class="max center-left">Les chiffres</h2>`,
		"<strong>2,5%</strong>",
		"Source: https://arcep.fr",
		"<speaker-notes>",
	} {
		if !strings.Contains(string(slides[1]), want) {
			t.Errorf("got %q, want it to contain %q", slides[1], want)
		}
	}
}

func TestLoadRendersTheTitleAsMarkdown(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "---\nlayout: default\ntitle: L'IA et le **carbone**\n---\ncontenu\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	// The exact string below already proves there is no <p> wrapper around the
	// title (goldmark would otherwise have produced
	// `<h2 ...><p>L'IA et le <strong>carbone</strong></p>\n</h2>`) and that the
	// apostrophe is bare rather than HTML-entity-escaped as `&#39;`.
	if want := `<h2 class="max center-left">L'IA et le <strong>carbone</strong></h2>`; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want it to contain %q", slides[0], want)
	}
}

func TestLoadDefaultsAMainClassFromTheLayout(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "---\nlayout: split\n---\n::term{path=sources}\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if want := `<main class="responsive max">`; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want it to contain %q — a slide with no class: key should get its layout's default from deck.defaultClasses", slides[0], want)
	}
}

func TestLoadClassKeyReplacesTheLayoutDefault(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md": "---\nlayout: default\nclass: main responsive xlarge-height\n---\ncontenu\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if want := `<main class="main responsive xlarge-height">`; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want it to contain %q — a class: key must replace the default layout's classes wholesale, not add to them", slides[0], want)
	}
	if strings.Contains(string(slides[0]), "center-align") {
		t.Errorf("got %q, want no center-align — class: replaces the default layout's hardcoded classes instead of appending to them, so a slide that declares xlarge-height must not also carry the default's own large-height/center-align", slides[0])
	}
}

func TestLoadPrefersMarkdownOverHTML(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md":   "# markdown\n",
		"demoit.html": "<h1>html</h1>\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !strings.Contains(string(slides[0]), "<h1>markdown</h1>") {
		t.Fatalf("got %q, want the Markdown deck to win", slides[0])
	}
}

func TestLoadPrefersTheLocalizedDeck(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md":    "# français\n",
		"demoit-en.md": "# english\n",
	})

	slides, err := deck.Load(folder, "en")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !strings.Contains(string(slides[0]), "english") {
		t.Fatalf("got %q, want the English deck", slides[0])
	}
}

func TestLoadFallsBackToTheLocalizedHTMLDeck(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.md":      "# français\n",
		"demoit-en.html": "<h1>english</h1>\n",
	})

	slides, err := deck.Load(folder, "en")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !strings.Contains(string(slides[0]), "<h1>english</h1>") {
		t.Fatalf("got %q, want the localized HTML deck to win over the default Markdown one", slides[0])
	}
}

func TestLoadKeepsHTMLDecksUntouched(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		"demoit.html": "<main>un</main>\n---\n<main>deux</main>\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if len(slides) != 2 {
		t.Fatalf("got %d slides, want 2", len(slides))
	}
	if got, want := string(slides[0]), "<main>un</main>\n"; got != want {
		t.Fatalf("got %q, want %q — an HTML deck must pass through unchanged", got, want)
	}
}

func TestLoadReportsAMissingDeck(t *testing.T) {
	t.Parallel()

	if _, err := deck.Load(t.TempDir(), ""); err == nil {
		t.Fatal("got no error, want one for a folder with no deck file")
	}
}

func TestLoadExposesUnknownFrontmatterKeys(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		".demoit/layouts/mine.html": `<p>{{ .Meta.badge }}</p>`,
		"demoit.md":                 "---\nlayout: mine\nbadge: nouveau\n---\n# Titre\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if want := "<p>nouveau</p>"; !strings.Contains(string(slides[0]), want) {
		t.Fatalf("got %q, want it to contain %q", slides[0], want)
	}
}

// Every layout is overridable, and that has to include its default classes:
// the engine used to fill .Class from a map keyed by layout name before the
// template ran, so a talk that replaced default.html found its own fallback
// unreachable. The fallback now lives in the template, which is the file the
// talk overrides.
func TestLoadLetsATalkOverrideALayoutsDefaultClasses(t *testing.T) {
	t.Parallel()

	folder := writeDeck(t, map[string]string{
		".demoit/layouts/default.html": `<main class="{{ if .Class }}{{ .Class }}{{ else }}classes-du-talk{{ end }}">{{ .Content }}</main>`,
		"demoit.md":                    "---\nlayout: default\n---\ncontenu\n",
	})

	slides, err := deck.Load(folder, "")
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if want := `<main class="classes-du-talk">`; !strings.Contains(string(slides[0]), want) {
		t.Errorf("got %q, want it to contain %q — a talk's own layout must supply its own default classes", slides[0], want)
	}
	if strings.Contains(string(slides[0]), "large-height") {
		t.Errorf("got %q, want none of the embedded default.html's classes — the talk's layout replaced it wholesale", slides[0])
	}
}
