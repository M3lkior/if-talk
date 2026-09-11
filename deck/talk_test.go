package deck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgageot/demoit/deck"
)

// writeTalk writes a talk.yml into a fresh folder and returns it.
func writeTalk(t *testing.T, content string) string {
	t.Helper()

	folder := t.TempDir()
	if content != "" {
		path := filepath.Join(folder, "talk.yml")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("unable to write talk.yml: %v", err)
		}
	}

	return folder
}

// The file sits at the folder's root. A copy left at the old .demoit/talk.yml
// is not read, and silently keeping a talk on stale settings is worse than
// showing it the defaults, so this pins the one location that counts.
func TestLoadTalkIgnoresTheOldDemoitLocation(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	if err := os.MkdirAll(filepath.Join(folder, ".demoit"), 0o755); err != nil {
		t.Fatalf("unable to create .demoit: %v", err)
	}
	if err := os.WriteFile(filepath.Join(folder, ".demoit", "talk.yml"), []byte("layout: content\n"), 0o600); err != nil {
		t.Fatalf("unable to write talk.yml: %v", err)
	}

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := talk.Layout, "default"; got != want {
		t.Errorf("got layout %q, want %q — .demoit/talk.yml must not be read", got, want)
	}
}

func TestLoadTalkReadsEveryKey(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "title: Impact Framework\nlayout: content\nlogo:\n  white: /images/a.svg\n  dark: /images/a-dark.svg\neventLogo:\n  white: /images/b.jpg\nfooter: adapted from demoit\n")

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
	if got, want := talk.Logo.White, "/images/a.svg"; got != want {
		t.Errorf("got logo %q, want %q", got, want)
	}
	if got, want := talk.Logo.Dark, "/images/a-dark.svg"; got != want {
		t.Errorf("got dark logo %q, want %q", got, want)
	}
	if got, want := talk.EventLogo.White, "/images/b.jpg"; got != want {
		t.Errorf("got event logo %q, want %q", got, want)
	}
	if got, want := talk.Footer, "adapted from demoit"; got != want {
		t.Errorf("got footer %q, want %q", got, want)
	}
}

// The order is the one the slides show: the speaker's logo, then the event's.
func TestLogosListsTheSpeakerThenTheEvent(t *testing.T) {
	t.Parallel()

	talk := deck.Talk{
		Logo:      deck.Logo{White: "/images/me.svg"},
		EventLogo: deck.Logo{White: "/images/event.svg"},
	}

	logos := talk.Logos()
	if got, want := len(logos), 2; got != want {
		t.Fatalf("got %d logos, want %d", got, want)
	}
	if got, want := logos[0].White, "/images/me.svg"; got != want {
		t.Errorf("got first logo %q, want %q", got, want)
	}
	if got, want := logos[1].White, "/images/event.svg"; got != want {
		t.Errorf("got second logo %q, want %q", got, want)
	}
}

// An undeclared eventLogo: block yields a zero Logo, which must be dropped
// rather than rendered as an <img> with an empty src.
func TestLogosDropsAnUndeclaredEventLogo(t *testing.T) {
	t.Parallel()

	talk := deck.Talk{Logo: deck.Logo{White: "/images/me.svg"}}

	if got, want := len(talk.Logos()), 1; got != want {
		t.Fatalf("got %d logos, want %d", got, want)
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

// A logo declaring only `white` is legal and means "this one image works on
// either ground": the layouts then emit it once, with no conditional class.
func TestLoadTalkAcceptsALogoWithNoDarkVariant(t *testing.T) {
	t.Parallel()

	folder := writeTalk(t, "eventLogo:\n  white: /images/event.jpg\n")

	talk, err := deck.LoadTalk(folder)
	if err != nil {
		t.Fatalf("got error %v, want none", err)
	}

	if got, want := talk.EventLogo.White, "/images/event.jpg"; got != want {
		t.Errorf("got event logo %q, want %q", got, want)
	}
	if talk.EventLogo.Dark != "" {
		t.Errorf("got dark variant %q, want none", talk.EventLogo.Dark)
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
	if len(talk.Logos()) != 0 {
		t.Errorf("got %d logos, want none", len(talk.Logos()))
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
