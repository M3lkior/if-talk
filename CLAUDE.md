# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A fork of [demoit](https://github.com/dgageot/demoit) (Go, Apache-2.0) used as a **speaker platform**: a local web server that renders HTML slides in the browser with live-coding web components (embedded ttys, auto-refreshing browser views, code viewer, VS Code in a container).

Two things live here at once:
- **The engine** — the Go server at the repo root (`main.go`, `handlers/`, `shell/`, `vscode/`, `livereload/`, `files/`, `flags/`).
- **The presentations** — one directory per talk (`sample/`, `impact-framework/`), each self-contained with its own slides, assets and demo sources.

## Commands

```bash
go build -o demoit .        # build (deps are vendored; no network needed)
go install                  # install to $GOPATH/bin
gofmt -l .                  # format check
golangci-lint run -c golangci.yml   # lint (config is golangci.yml, not .golangci.yml)

docker buildx bake          # cross-compile all 4 platforms into ./out (same as CI)
```

Run a presentation (arg = presentation folder, defaults to `.`):

```bash
demoit impact-framework                 # http://localhost:8888
demoit --dev impact-framework           # + live reload on file change
demoit --locale en impact-framework     # loads demoit-en.html instead of demoit.html
demoit --port 8080 --host 0.0.0.0 --shellport 9999 impact-framework
```

There are **no Go tests** in this repo. Verifying a change means running `demoit --dev <folder>` and looking at the slides.

The `impact-framework` talk needs extra local setup before the demo slides work: Docker/Colima running, and Grafana up via `docker-compose -f impact-framework/sources/grafana/grafana-compose.yml up`. See `impact-framework/README.md`.

## Architecture

**Slides are one HTML file split on `---`.** `handlers/step.go:readSteps` reads `demoit.html` (or `demoit-<locale>.html` when `--locale` is set and the file exists) from the presentation folder, splits the bytes on the literal `---`, and each part becomes a `Page` injected as raw `template.HTML` into the embedded `handlers/resources/index.tmpl.html`. Slide N is served at `/N`; slide 0 at `/`. Consequence: **`---` anywhere in the file starts a new slide**, including inside an HTML comment — that's why commented-out slides in `impact-framework/demoit.html` escape it as `-&#45;&#45;`.

**Everything the browser loads comes from `<folder>/.demoit/`.** `handlers/static.go` serves `/js/`, `/images/`, `/fonts/`, `/media/`, `/style.css`, `/favicon.ico` out of that directory. `.demoit/` is mandatory — `handlers.VerifyConfiguration` refuses to start without it. Static URLs are cache-busted with a `?hash=` sha256 prefix computed by the `hash` template func.

**The web components live in the presentation, not in the engine.** `.demoit/js/demoit.js` is per-talk (copied and diverged between `sample/` and `impact-framework/`), defines all custom elements as shadow-DOM `BaseHTMLElement` subclasses, and is loaded as an ES module by the index template. Editing a component means editing that talk's copy:

| Element | What it does | Server route it hits |
|---|---|---|
| `<web-term path="folder">` | Multi-tab tty, `+` adds a tab | `/shell/{folder}` |
| `<web-browser src="...">` | Iframed browser with URL bar + refresh | — |
| `<source-code folder= files= start-lines= end-lines= code_style=>` | IDE-like tabbed code viewer, chroma-highlighted | `/sourceCode/...` |
| `<vs-code path="folder">` | code-server in a Docker container | `/beta/vscode/{folder}` |
| `<split-view>` | auto-fit grid of the above | — |
| `<fake-window title="">` | macOS-style chrome; green dot maximizes | — |
| `<speaker-notes>` | invisible on stage; content forwarded to the notes window | — |
| `<nav-arrows>` | also binds ArrowRight/Left, PageUp/Down, Space | — |

`impact-framework/.demoit/js/demoit.js` additionally imports mermaid from a CDN and defines `<title-bar>`. Both talks style on top of [beercss](https://www.beercss.com/) (CDN) — the `s12 m6 l6`, `responsive`, `max`, `center-align` classes in the slides are beercss.

**Terminals are a reverse-proxied gotty.** `shell.ListenAndServe` runs a gotty server on `--shellport` bound to 127.0.0.1; `main.go` reverse-proxies `/tty` to it. `handlers/shell.go` doesn't render anything — it builds a `cd <folder>; source .demoit/.bashrc; HISTFILE=<copy> exec $SHELL` command line and 303-redirects to `/tty?arg=...`. So **`.demoit/.bashrc` and `.demoit/.bash_history` are how you pre-seed a demo terminal** (history is copied to a temp file first so the demo doesn't mutate the repo).

**Speaker notes are a second window synced over `BroadcastChannel("demoit_nav")`.** `/speakernotes` renders a standalone page; the main window broadcasts `{currentSlideId, stepCount, currentSlideTitle, speakerNotes}` on every load, and navigation from the notes window comes back as `{destinationSlideId}`. Both directions must keep emitting even for slides without notes, or the windows desync.

**Other routes:** `/grid` renders every step as an iframe (also warms the browser cache so slide transitions don't lag); `/pdf` drives headless Chrome via chromedp at 1920x1080@2x, 4 parallel tabs, into a gofpdf document; `/qrcode?url=` returns a PNG; `/last` redirects to the final slide.

**VS Code is opt-in and heavy.** `vscode/server.go` pulls a pinned `codercom/code-server` image and starts container `demoit-vscode` once (`sync.Once`), bind-mounting the presentation folder at `/app` on port 18080. Requires a working Docker daemon; failures are logged, not fatal.

## Conventions

- Go files carry the Apache-2.0 header from Google/David Gageot — keep it on new files in the engine.
- Dependencies are **vendored** (`vendor/`); after touching `go.mod` run `go mod vendor`.
- Adding a talk = new top-level folder with `demoit.html` + `.demoit/` (copy an existing `.demoit/js/demoit.js` and `style.css` as the starting point). Localized decks are `demoit-<locale>.html` side by side.
- Palette and per-talk theming live in `<folder>/.demoit/style.css` via `--color-main` / `--primary` custom properties.
