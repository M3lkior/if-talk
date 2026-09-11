package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgageot/demoit/handlers"
)

// serveEngineCSS returns the body the handler writes.
func serveEngineCSS(t *testing.T) (*httptest.ResponseRecorder, string) {
	t.Helper()

	w := httptest.NewRecorder()
	handlers.EngineCSS(w, httptest.NewRequest(http.MethodGet, "/demoit.css", nil))

	return w, w.Body.String()
}

// The engine stylesheet is a build artefact: hack/css.sh generates it and it is
// committed. Nothing in `go build` regenerates it, so a stale copy -- or one
// rebuilt from a broken source -- would ship silently. These tests are the
// guard.
func TestEngineCSSIsServed(t *testing.T) {
	t.Parallel()

	w, css := serveEngineCSS(t)

	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
	if got, want := w.Header().Get("Content-Type"), "text/css"; !strings.HasPrefix(got, want) {
		t.Errorf("got content type %q, want it to start with %q", got, want)
	}
	if css == "" {
		t.Fatal("got an empty stylesheet, want the generated CSS")
	}
}

// The classes deck/directive assembles are built in Go, so they exist in no
// source file and Tailwind's scanner can never meet them. They reach the
// stylesheet only through the `@source inline(...)` lines of
// styles/demoit.src.css. Dropping those leaves every split{cols=} pane
// unstyled on stage with nothing to say why -- the same silent failure
// transform.go already documents for a weight outside 1..12.
func TestEngineCSSCarriesTheClassesBuiltInGo(t *testing.T) {
	t.Parallel()

	_, css := serveEngineCSS(t)

	for _, want := range []string{
		".col-span-1", ".col-span-4", ".col-span-8", ".col-span-12",
		".h-stage-small", ".h-stage-medium", ".h-stage-large", ".h-stage-xlarge",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("got a stylesheet without %q: is the matching @source inline() line missing from styles/demoit.src.css?", want)
		}
	}
}

// An @utility is emitted only where a scanned file uses it, so the chassis is
// declared inline in styles/demoit.src.css rather than left to the layouts.
// This pins that: a talk overriding a layout must not have to re-emit the
// chassis from its own stylesheet.
func TestEngineCSSCarriesTheChassis(t *testing.T) {
	t.Parallel()

	_, css := serveEngineCSS(t)

	for _, want := range []string{
		".slide-main", ".slide-header", ".slide-footer", ".slide-prose",
		".contact-row", ".card", ".stage", ".no-stage",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("got a stylesheet without %q", want)
		}
	}
}

// The theme and the display profiles are what the switcher acts on: a build
// that lost them would leave `t` and `d` pressing on nothing.
func TestEngineCSSCarriesTheThemeAndProfiles(t *testing.T) {
	t.Parallel()

	_, css := serveEngineCSS(t)

	for _, want := range []string{
		".dark",
		`[data-display="tv"]`,
		`[data-display="projector"]`,
		"--dm-window-bg",
		"--stage-scale",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("got a stylesheet without %q", want)
		}
	}
}
