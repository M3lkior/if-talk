package deck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgageot/demoit/deck"
)

// writeTalk writes a .demoit/talk.yml into a fresh folder and returns it.
func writeTalk(t *testing.T, content string) string {
	t.Helper()

	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, ".demoit"), 0o755); err != nil {
		t.Fatalf("unable to create .demoit: %v", err)
	}
	if content != "" {
		path := filepath.Join(folder, ".demoit", "talk.yml")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("unable to write talk.yml: %v", err)
		}
	}

	return folder
}

func TestLoadTalkReadsEveryKey(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "title: Impact Framework\nlayout: content\nlogos:\n  - /images/a.svg\n  - /images/b.jpg\nfooter: adapted from demoit\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := string(talk.Title), "Impact Framework"; got != want {
		t.Errorf("got title %q, want %q", got, want)
	}
	if got, want := talk.Layout, "content"; got != want {
		t.Errorf("got layout %q, want %q", got, want)
	}
	if got, want := len(talk.Logos), 2; got != want {
		t.Errorf("got %d logos, want %d", got, want)
	}
	if got, want := talk.Footer, "adapted from demoit"; got != want {
		t.Errorf("got footer %q, want %q", got, want)
	}
}

// A talk's title and subtitle go through the same Markdown pipeline as a
// slide's title: key, and for the same reason: a cover writes
// `*Impact Framework*` so the word takes the accent treatment, and an
// apostrophe must not come back as an HTML entity the way html/template would
// escape a plain string.
func TestLoadTalkRendersTheTitleAndSubtitleAsMarkdown(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "title: A la découverte d'*Impact Framework*.\nsubtitle: 10 avril 2025 — **Ludovic Dussart**\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := string(talk.Title), "A la découverte d'<em>Impact Framework</em>."; got != want {
		t.Errorf("got title %q, want %q", got, want)
	}
	if got, want := string(talk.Subtitle), "10 avril 2025 — <strong>Ludovic Dussart</strong>"; got != want {
		t.Errorf("got subtitle %q, want %q", got, want)
	}
}

// A talk with neither key renders two empty strings, not a <p></p>: that is
// what lets a layout test them with {{ with }}.
func TestLoadTalkLeavesAnAbsentTitleEmpty(t *testing.T) {
	t.Parallel()

	talk, err := deck.LoadTalk(writeTalk(t, "layout: content\n"))
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if talk.Title != "" {
		t.Errorf("got title %q, want it empty", talk.Title)
	}
	if talk.Subtitle != "" {
		t.Errorf("got subtitle %q, want it empty", talk.Subtitle)
	}
}

func TestLoadTalkAcceptsAMissingFile(t *testing.T) {
	t.Parallel()

	talk, err := deck.LoadTalk(writeTalk(t, ""))
	if err != nil {
		t.Fatalf("got error %v, want none — a talk without talk.yml must work", err)
	}

	if got, want := talk.Layout, "default"; got != want {
		t.Errorf("got layout %q, want %q", got, want)
	}
	if len(talk.Logos) != 0 {
		t.Errorf("got %d logos, want none", len(talk.Logos))
	}
}

func TestLoadTalkRejectsInvalidYAML(t *testing.T) {
	t.Parallel()

	if _, err := deck.LoadTalk(writeTalk(t, "logos: [unterminated\n")); err == nil {
		t.Fatal("got no error, want one for invalid YAML")
	}
}

func TestLoadTalkReadsTheThemeBlock(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "layout: content\ntheme:\n  stage: true\n  dark: true\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !talk.Theme.Stage {
		t.Error("got theme.stage false, want true")
	}
	if !talk.Theme.Dark {
		t.Error("got theme.dark false, want true")
	}
}

// Both flags default to false, and that default is the point: it is what keeps
// sample/ off the new chassis. sample declares no theme block, so it gets
// neither the fixed stage nor a theme switcher -- and its own demoit.js carries
// a light window chrome in hard-coded colours that would not follow a dark
// page.
func TestLoadTalkDefaultsTheThemeBlockToOff(t *testing.T) {
	t.Parallel()

	talk, err := deck.LoadTalk(writeTalk(t, "layout: content\n"))
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if talk.Theme.Stage {
		t.Error("got theme.stage true, want false for a talk that declares no theme block")
	}
	if talk.Theme.Dark {
		t.Error("got theme.dark true, want false for a talk that declares no theme block")
	}
}

// One flag without the other must work: a talk may want the fixed stage and no
// dark mode, or the reverse.
func TestLoadTalkReadsOneThemeFlagAlone(t *testing.T) {
	t.Parallel()

	talk, err := deck.LoadTalk(writeTalk(t, "theme:\n  stage: true\n"))
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if !talk.Theme.Stage {
		t.Error("got theme.stage false, want true")
	}
	if talk.Theme.Dark {
		t.Error("got theme.dark true, want false")
	}
}
