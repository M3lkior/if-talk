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
go test ./...               # run the Go test suite (deck/, deck/directive/, highlight/)
gofmt -l .                  # format check
golangci-lint run -c golangci.yml   # lint (config is golangci.yml, not .golangci.yml)

docker buildx bake          # cross-compile all 4 platforms into ./out (same as CI)
```

Run a presentation (arg = presentation folder, defaults to `.`):

```bash
demoit impact-framework                 # http://localhost:8888
demoit --dev impact-framework           # + live reload on file change
demoit --locale en impact-framework     # loads the localized deck (demoit-en.md or demoit-en.html) if present
demoit --port 8080 --host 0.0.0.0 --shellport 9999 impact-framework
```

This repo has a Go test suite — `deck/`, `deck/directive/` and `highlight/` all carry tests, and `go test ./...` is how you verify a change to that code. It doesn't assert what a browser paints, though: a rendering or layout change (a `.demoit/layouts/*.html`, CSS, the index template) still needs `demoit --dev <folder>` and a look at the slides — that's exactly how this plan's own migration of `impact-framework` to Markdown was checked, slide-by-slide against the old HTML render, not against the suite.

The `impact-framework` talk needs extra local setup before the demo slides work: Docker/Colima running, and Grafana up via `docker-compose -f impact-framework/sources/grafana/grafana-compose.yml up`. See `impact-framework/README.md`.

## Architecture

**Slides come from a deck file resolved by `deck.Load`.** `handlers/step.go:readSteps` calls `deck.Load` (`deck/deck.go`), which tries, in this order, `demoit-<locale>.md`, `demoit-<locale>.html`, `demoit.md`, `demoit.html` inside the presentation folder and renders whichever it finds first (`--locale` only changes which of the first two names is tried). Each slide becomes a `Page` whose `HTML` is injected as raw `template.HTML` into the embedded `handlers/resources/index.tmpl.html`. Slide N is served at `/N`; slide 0 at `/`.

An **HTML deck** (`sample/demoit.html`, `impact-framework/demoit-en.html`) goes through the original, unchanged path: `bytes.Split` on the literal `---`. `---` anywhere in the file starts a new slide, including inside an HTML comment, so a commented-out slide there still needs its `---` escaped (e.g. `-&#45;&#45;`) to survive.

A **Markdown deck** (`demoit.md` / `demoit-<locale>.md`, e.g. `impact-framework/demoit.md`) is split instead by `deck/split.go`, which only cuts on a line of three-or-more dashes that sits outside a fenced code block (triple-backtick or `~~~`). Two consequences: `---` is no longer usable as Markdown's thematic-break syntax (write `***`), and it no longer needs escaping *inside a code fence* — but the splitter is not HTML-comment-aware, so a `---` left inside a commented-out block still cuts a slide. Each slide is then parsed by goldmark with a hand-written directive extension (`deck/directive/`) and rendered through a layout template (`deck/layouts/`, see below). `deck/frontmatter.go` reads the slide's `key: value` header. `highlight/` wraps chroma and is shared by plain fenced code blocks in the Markdown body and by the `::code` directive's viewer route (`handlers/code.go`). A slide that fails to parse renders as an error card naming the file and line instead of crashing the server (`deck.errorHTML`).

**Frontmatter is per-slide.** The `---` line that separates two slides doubles as the opening fence of the next slide's frontmatter: the block right after a separator (or at file start) is read as frontmatter when its first line looks like a lowercase `key:` pair (`deck/split.go:isFrontmatterStart`), and the scan continues to the next `---` line. A body paragraph that happens to start the same way (`note: ceci est du texte`) followed later by a `---` is swallowed as frontmatter the same way — a capitalized first word (`Note:`) reads as prose instead. Known keys (`deck/frontmatter.go`):

| Key | Effect |
|---|---|
| `layout` | Name of the layout template; defaults to the talk's own `layout:` (from `talk.yml`, itself defaulting to `default`). |
| `title` | Rendered through the same Markdown pipeline as the body — `**bold**` works, an apostrophe is not entity-escaped — and shown in the layout's header. |
| `source` | Shown as a "Source: …" line under the slide. |
| `class` | CSS classes for the slide's `<main>`. **Replaces** the layout's default class list (`deck.defaultClasses`) wholesale rather than adding to it. |
| `height` | Height class for a `split` layout's column container (`xlarge` by default). |
| `speakernotes` | Rendered through Markdown too — a `-` bullet list becomes `<ul><li>`, invisible on stage — and forwarded to the speaker-notes window. |

Any other key lands untouched in `.Meta`, for a custom layout that wants to read it.

**Directives are the Markdown deck's web-component syntax** (`deck/directive/`). `::name{key=value key2="a value with spaces"}` on its own line is self-closing (2 colons); `:::name{...}` opens a container that needs a matching `:::` line to close (3 colons) and can nest. Attributes are always `key=value`, space-separated — there is no goldmark `#id`/`.class` shorthand. Authoritative list is `deck/directive/render.go`'s `specs` map:

| Directive | Renders as | Notes |
|---|---|---|
| `::term{path=}` | `<web-term path=...>` | |
| `::browser{src=}` | `<web-browser src=...>` | |
| `::vscode{path=}` | `<vs-code path=...>` | |
| `::code{folder= files= lines= style=}` | `<source-code ...>` | see below |
| `:::window{title=}` … `:::` | `<fake-window title=...>` | |
| `:::speakernotes` … `:::` | `<speaker-notes>` | the `speakernotes:` frontmatter key covers the common case; this exists for reuse inside a directive tree |
| `:::split{height=}` … `:::` | `<split-view class="{height}-height">` | same tag `layout: split` produces; children are auto-arranged by the web component |
| `:::split{cols=4,8 height=}` … `:::` | a weighted `<div class="grid ...">`, one `<div class="sN">` per direct child | each pane needs a **blank line** before it: the wrapper counts direct block children, and two panes glued together without a blank line become one child, so the column count won't match |

`grid`/`col` can also be authored directly, for a layout `split{cols=}` can't produce: `split{cols=}` always emits a single row of panes, one hardcoded `class="sN"` div per direct child and nothing more. `impact-framework/demoit.md:349-380` builds a two-row layout this way — a full-width `s12` heading, then an `s4 xlarge-height` / `s8 xlarge-height` row below it — with `:::grid{class="grid xlarge-height"}` and three `:::col{...}` panes, each free to carry a compound class (`s4 xlarge-height`) that `split{cols=}`'s generated columns can never have. For a single row of plain weighted panes, `split{cols=}` stays the shorter way to write it.

`::code` translates Markdown's comma-separated lists into `source-code`'s attributes: `files=a.yml,b.yml` becomes `files="a.yml b.yml"`; `lines=11-20,4-9` (one range per file) becomes `start-lines="11;4" end-lines="20;9"`. These ranges **highlight** lines in the viewer (`html.HighlightLines` in `handlers/code.go`) — the whole file is always displayed, never truncated. `style=` maps to `code_style=` and defaults to `"vs"` when omitted, matching `demoit.js`'s own default.

**Layouts.** Six ship embedded in `deck/layouts/`: `bare`, `content`, `cover`, `default`, `quote`, `split` (plus `partials.html`, which only defines the shared `header`/`source`/`notes` template blocks — not a layout of its own). Every one, `partials` included, is overridden by a same-named file in `<talk>/.demoit/layouts/` (`TestEveryEmbeddedLayoutIsOverridable`), and a talk may add layout names the engine has never heard of — `impact-framework/.demoit/layouts/default-h3.html` (an `<h3>` header instead of `<h2>`) is a live example. A talk-added layout gets no entry in `deck.defaultClasses`, so `.Class` arrives empty for a slide naming no `class:` key; `default-h3.html` hardcodes the same fallback `default.html` itself uses for that case. Note for future cleanup: `default.html` and `content.html` are, as of this writing, byte-for-byte identical — they differ only by which `defaultClasses` entry maps to them — so the catalogue currently gives no visible reason to pick one over the other.

**`.demoit/talk.yml`** (`deck/talk.go`) declares `layout` (the deck-wide default a slide falls back to when its own frontmatter names none) and `logos` (image paths rendered by `cover.html` and the header partial). It also parses `title` and `footer`, but no embedded or `impact-framework` layout reads either field — the on-screen footer (progress bar + "Slides framework adapted from demoit") is static markup in `handlers/resources/index.tmpl.html`, unrelated to `Talk.Footer`.

**Known, accepted gaps in `impact-framework`'s Markdown migration** (harmless — don't "fix" these as bugs): a slide's `Content` has exactly one insertion point inside `<main>`, so a paragraph the original HTML placed as a sibling *after* `</main>` now renders inside it (the "Et l'IA dans tout ça ?" slide); a disabled/commented block that used to sit *before* a slide's header in the raw file now renders *after* it, because the layout always emits the header first (same slide, plus two others).

**Everything the browser loads comes from `<folder>/.demoit/`.** `handlers/static.go` serves `/js/`, `/images/`, `/fonts/`, `/media/`, `/style.css`, `/favicon.ico` out of that directory. `.demoit/` is mandatory — `handlers.VerifyConfiguration` refuses to start without it. Static URLs are cache-busted with a `?hash=` sha256 prefix computed by the `hash` template func.

**The web components live in the presentation, not in the engine.** `.demoit/js/demoit.js` is per-talk (copied and diverged between `sample/` and `impact-framework/`), defines all custom elements as shadow-DOM `BaseHTMLElement` subclasses, and is loaded as an ES module by the index template. Editing a component means editing that talk's copy:

| Element | What it does | Server route it hits |
|---|---|---|
| `<web-term path="folder">` | Multi-tab tty, `+` adds a tab | `/shell/{folder}` |
| `<web-browser src="...">` | Iframed browser with URL bar + refresh | — |
| `<source-code folder= files= start-lines= end-lines= code_style=>` | IDE-like tabbed code viewer, chroma-highlighted; `start-lines`/`end-lines` **highlight** those lines, they don't truncate the file | `/sourceCode/...` |
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
- Adding a talk = new top-level folder with `demoit.md` (or `demoit.html`) + `.demoit/` (copy an existing `.demoit/js/demoit.js` and `style.css` as the starting point). Localized decks are `demoit-<locale>.md`/`demoit-<locale>.html` side by side.
- Palette and per-talk theming live in `<folder>/.demoit/style.css` via `--color-main` / `--primary` custom properties.
