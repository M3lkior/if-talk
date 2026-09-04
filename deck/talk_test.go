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

	if got, want := talk.Title, "Impact Framework"; got != want {
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
