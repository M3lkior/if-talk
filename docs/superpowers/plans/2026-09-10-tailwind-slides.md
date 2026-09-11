# Migration beercss → Tailwind Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remplacer beercss par Tailwind CSS v4, poser les slides sur une scène fixe 1920×1080 mise à l'échelle, et ajouter un thème sombre et des profils de projection commutables que les runtimes embarqués suivent.

**Architecture:** Deux feuilles Tailwind générées et commitées — une embarquée dans le binaire pour le châssis du moteur, une par talk pour les utilitaires de ses slides — construites par un binaire Tailwind standalone épinglé, sans `node_modules` et sans réseau à l'exécution. Le châssis met les slides sur une scène de 1920×1080 px mise à l'échelle par `transform: scale()`, ce qui rend la géométrie identique sur tout écran 16:9. Thème et profil d'affichage sont des jetons CSS commutés par une classe et un attribut sur `<html>`, propagés aux shadow DOM des composants par custom properties.

**Tech Stack:** Go 1.x (deps vendorées), Tailwind CSS v4.3.3 (CLI standalone), goldmark, chroma, gotty, beercss (retiré).

**Spec:** `docs/superpowers/specs/2026-09-10-tailwind-slides-design.md`

## Global Constraints

- **Tailwind CLI standalone v4.3.3**, exactement cette version, épinglée dans `hack/css.sh`. Binaires disponibles et vérifiés : `macos-arm64`, `macos-x64`, `linux-x64`, `linux-arm64`.
- **Aucun `node_modules`, aucun `package.json`.** Le réseau ne sert qu'au téléchargement du CLI et au build ; jamais à l'exécution du serveur.
- **Le CSS généré est commité.** `handlers/resources/demoit.css` et `<talk>/.demoit/tailwind.css` sont des artefacts versionnés, pas des fichiers ignorés.
- **`.tools/` est gitignoré** et contient le CLI téléchargé.
- **Les fichiers Go du moteur portent l'en-tête Apache-2.0** de Google/David Gageot. **Les `_test.go` n'en portent pas** — convention établie du repo, et `goheader` n'a pas de template dans `golangci.yml`.
- **Un fichier importé avec `theme(reference)` ne doit contenir que des blocs `@theme`.** Vérifié : toute autre règle produit `Error: Files imported with '@import "…" theme(reference)' must only contain '@theme' blocks.` D'où la séparation `tokens.css` / `components.css`.
- **La feuille d'un talk doit référencer le thème par défaut de Tailwind** (`@import "tailwindcss/theme.css" theme(reference);`), sans quoi `gap-4`, `text-white` et tout utilitaire adossé à `--spacing` ou `--color-*` ne sont **pas générés** — silencieusement. Vérifié.
- **`@source "../"` scanne les dossiers cachés** : `.demoit/layouts/*.html` est bien pris. Vérifié.
- **Après toute modification de `go.mod`** : `go mod vendor`.
- **Vérification à chaque tâche** : `go build -o /dev/null .`, `go test ./...`, `gofmt -l .` (vide hors `vendor/`), `golangci-lint run -c golangci.yml`.
- **`rem` vaut 16 px**, aujourd'hui comme après. Mesuré, pas supposé : le `font-size: 3vw` de `style.css:28` est écrasé par le `html{font-size:var(--size)}` de beercss, chargé après lui dans `index.tmpl.html`. Aucune taille en `rem` du deck n'est à convertir, et supposer un `rem` à 57.6 px triplerait tout.
- **La référence visuelle est hors dépôt**, dans le scratchpad de session : voir Tâche 0.

---

## Tâche 0 : La référence visuelle (prérequis, déjà exécutée)

Cette tâche est **déjà faite** et n'est décrite que pour que l'exécutant sache où se trouve la référence et comment la reproduire.

L'arbre a été capturé **avant toute modification**, depuis le worktree à `d4af1bd` :

- `<scratchpad>/before/deck-before.pdf` — export `/pdf` du deck complet, 18 slides, rendu par le pipeline chromedp de l'app à 1920×1080@2x. C'est la référence canonique.
- `<scratchpad>/before/png-1920/slide-NN.png` — une capture par slide à 1920×1080.
- `<scratchpad>/shoot1.sh`, `<scratchpad>/shootall.sh` — les scripts de capture.

**Pourquoi un plafond de temps dans les scripts :** Chrome ne sort pas de lui-même sur les slides qui tiennent une socket vivante (les iframes de tty). `--virtual-time-budget` bloque au lieu d'aider. Les scripts lancent donc un Chrome par slide et le tuent au bout de N secondes ; la capture est déjà écrite.

Pour recapturer l'après à un moment donné :

```bash
go build -o /tmp/demoit-after .
/tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shootall.sh <scratchpad>/after/png-1920 http://127.0.0.1:8899 17 18
```

---

## File Structure

**Créés — le socle CSS du moteur**

| Fichier | Responsabilité |
|---|---|
| `styles/tokens.css` | **Uniquement** des blocs `@theme`. Couleurs, familles de police, rayons, hauteurs de scène. Importable en `theme(reference)`. |
| `styles/components.css` | `@custom-variant dark`, et les `@utility` du châssis : `slide-main`, `slide-header`, `slide-footer`, `slide-prose`, `contact-row`, `card`. |
| `styles/demoit.src.css` | L'entrée du build moteur : couches, preflight, imports, `@source`, et le CSS non-utilitaire du châssis (`.stage`, `#progression`, la fenêtre de notes, les profils d'affichage, les custom properties de chrome). |
| `handlers/resources/demoit.css` | **Généré, commité, embarqué.** Servi sur `/demoit.css`. |
| `handlers/resources/css.go` | Le `//go:embed` de `demoit.css` et son handler HTTP. |
| `hack/css.sh` | Télécharge le CLI épinglé et construit les deux cibles. |

**Créés — côté talk**

| Fichier | Responsabilité |
|---|---|
| `impact-framework/.demoit/tokens.css` | Copie des jetons, surchargeable par le talk. |
| `impact-framework/.demoit/components.css` | Copie des composants du châssis. |
| `impact-framework/.demoit/tailwind.src.css` | L'entrée du build du talk. |
| `impact-framework/.demoit/tailwind.css` | **Généré, commité.** Servi sur `/tailwind.css`. |
| `impact-framework/.demoit/fonts/*.woff2` | Poppins 100, 400, 500, 700, vendorés. |

**Modifiés**

| Fichier | Changement |
|---|---|
| `handlers/resources/index.tmpl.html` | La scène, le script de thème pré-paint, les nouvelles feuilles ; beercss, `beer.min.js`, `material-dynamic-colors`, `<dialog id="maximized">` et `class="light"` retirés. |
| `handlers/resources/grid.tmpl.html` | Thème et jetons. |
| `handlers/resources/speakernotes.tmpl.html` | Thème, jetons, écoute du canal de thème. |
| `deck/layouts/{bare,content,cover,default,quote,split,partials}.html` | Classes beercss → châssis Tailwind. |
| `deck/directive/transform.go` | Émet `grid grid-cols-12 gap-4` et `col-span-N`. |
| `deck/directive/render.go` | `splitAttributes` émet `h-stage-*`. |
| `deck/talk.go` | Bloc `theme:` (`stage`, `dark`). |
| `handlers/step.go` | `Page.Stage`, `Page.Dark`. |
| `main.go` | Routes `/demoit.css` et `/tailwind.css`. |
| `impact-framework/demoit.md` | Classes des 18 slides. |
| `impact-framework/demoit-en.html` | Idem, deck HTML. |
| `impact-framework/.demoit/layouts/{cover,default-h3}.html` | Idem, layouts d'override. |
| `impact-framework/.demoit/style.css` | Nettoyé de ce qui devient jeton ; Poppins local. |
| `impact-framework/.demoit/js/demoit.js` | Chrome par custom properties, thème de `source-code` et mermaid, échelle de scène, classes mortes retirées. |
| `impact-framework/.demoit/talk.yml` | Bloc `theme:`. |
| `impact-framework/demoit.html` | Une ligne d'en-tête. |
| `.gitignore` | `.tools/`. |
| `CLAUDE.md` | Cinq sections qui décrivent beercss. |

**Supprimés :** `impact-framework/.demoit/poppins.css`.

**Tests modifiés ou créés :** `deck/directive/transform_test.go`, `deck/directive/directive_test.go`, `deck/talk_test.go`, `handlers/resources_test.go` (créé).

---

## Tâche 1 : La chaîne de build et la feuille du moteur

Objectif : `hack/css.sh` produit une feuille Tailwind embarquée et servie, **sans encore la lier** dans les pages. Aucun changement visuel.

**Files:**
- Create: `hack/css.sh`, `styles/tokens.css`, `styles/components.css`, `styles/demoit.src.css`, `handlers/resources/css.go`
- Create (généré): `handlers/resources/demoit.css`
- Modify: `.gitignore`, `main.go:59-64`
- Test: `handlers/resources_test.go`

**Interfaces:**
- Produces: la route `/demoit.css` ; `handlers.EngineCSS` (`func(http.ResponseWriter, *http.Request)`) ; les jetons `--color-main`, `--color-fg`, `--color-fg-muted`, `--color-surface`, `--color-void`, `--color-outline`, `--height-stage-{small,medium,large,xlarge}`, `--radius-card` ; les utilitaires `slide-main`, `slide-header`, `slide-footer`, `slide-prose`, `contact-row`, `card` ; la variante `dark:`.
- Consumes: rien.

- [ ] **Step 1: Écrire `hack/css.sh`**

```bash
#!/usr/bin/env bash
# Builds the Tailwind stylesheets of demoit.
#
#   hack/css.sh engine                     -> handlers/resources/demoit.css
#   hack/css.sh <talk-folder>              -> <talk-folder>/.demoit/tailwind.css
#   hack/css.sh <talk-folder> --watch      -> rebuilds on change
#
# The CLI is a pinned standalone binary: no node_modules, and the network is
# only ever touched here, never while a talk is running.
set -euo pipefail

TAILWIND_VERSION=4.3.3
TOOLS_DIR=${TOOLS_DIR:-.tools}
CLI="$TOOLS_DIR/tailwindcss-$TAILWIND_VERSION"

usage() {
    echo "usage: hack/css.sh engine|<talk-folder> [--watch]" >&2
    exit 2
}

platform() {
    case "$(uname -s)/$(uname -m)" in
        Darwin/arm64)   echo macos-arm64 ;;
        Darwin/x86_64)  echo macos-x64 ;;
        Linux/aarch64)  echo linux-arm64 ;;
        Linux/arm64)    echo linux-arm64 ;;
        Linux/x86_64)   echo linux-x64 ;;
        *)
            echo "no pinned Tailwind binary for $(uname -s)/$(uname -m)" >&2
            exit 1
            ;;
    esac
}

ensure_cli() {
    if [ -x "$CLI" ]; then
        return
    fi

    mkdir -p "$TOOLS_DIR"
    url="https://github.com/tailwindlabs/tailwindcss/releases/download/v$TAILWIND_VERSION/tailwindcss-$(platform)"
    echo "downloading Tailwind CLI $TAILWIND_VERSION" >&2
    curl -sSfL --retry 3 -o "$CLI.tmp" "$url"
    chmod +x "$CLI.tmp"
    mv "$CLI.tmp" "$CLI"
}

[ $# -ge 1 ] || usage

target=$1
shift

watch=""
if [ "${1:-}" = "--watch" ]; then
    watch="--watch"
    shift
fi
[ $# -eq 0 ] || usage

ensure_cli

if [ "$target" = "engine" ]; then
    input=styles/demoit.src.css
    output=handlers/resources/demoit.css
else
    input="$target/.demoit/tailwind.src.css"
    output="$target/.demoit/tailwind.css"
    if [ ! -f "$input" ]; then
        echo "$input not found: is $target a talk folder?" >&2
        exit 1
    fi
fi

echo "building $output" >&2
# --optimize, not --minify: the generated CSS is committed, so it stays
# readable in diffs and reviewable, while still being deduplicated.
exec "$CLI" --input "$input" --output "$output" --optimize $watch
```

- [ ] **Step 2: Rendre le script exécutable et ignorer le CLI**

```bash
chmod +x hack/css.sh
printf '.tools/\n' >> .gitignore
```

- [ ] **Step 3: Écrire `styles/tokens.css` — uniquement des `@theme`**

Rappel de la contrainte vérifiée : une seule règle non-`@theme` dans ce fichier casse le build du talk.

```css
/*
Design tokens of the demoit chassis.

This file must contain @theme blocks and nothing else: a talk's stylesheet
imports it with `theme(reference)` to get the vocabulary without emitting the
variables twice, and Tailwind rejects a referenced file that carries any other
rule.

Sizes are in rem on purpose. The root font size stays at 16px -- the value
that is already in effect, and Tailwind's own default -- so every rem in the
deck keeps the length it has today. The `font-size: 3vw` in the talk's
style.css never applied: beercss loads after it and its own `html` rule wins
by source order.
*/
@theme {
    --font-sans: "Poppins", ui-sans-serif, system-ui, sans-serif;
    --font-mono: "Roboto Mono", ui-monospace, monospace;

    /* A talk overrides --color-main in its own .demoit/style.css. */
    --color-main: #0f15fd;
    --color-main-dark: #7c81ff;

    --color-fg: #3d4043;
    --color-fg-muted: #6b7075;
    --color-surface: #ffffff;
    --color-outline: #d5d7da;

    /* The letterbox around the stage, on screens whose aspect ratio differs. */
    --color-void: #000000;

    /*
    Pane heights, as a share of the slide's main area rather than absolute
    lengths. Two reasons, neither of them "the old values were absurd" -- at
    16px/rem they were sane. First, .xlarge-height was 64rem, 1024px, for a
    1080px stage: add a header and a footer and the pane overflows. Second,
    its fallback came from an `@media (max-height:1024px)` query -- a question
    about the window, which a fixed stage makes meaningless. A share of the
    main area cannot overflow, and asks no media query.

    Names match what the deck writes in frontmatter (height: xlarge), because
    split.html interpolates that value verbatim.
    */
    --height-stage-small: 40%;
    --height-stage-medium: 65%;
    --height-stage-large: 85%;
    --height-stage-xlarge: 100%;

    --radius-card: 2rem;
}
```

- [ ] **Step 4: Écrire `styles/components.css`**

```css
/*
The chassis: the dark variant and the handful of named utilities that would
otherwise cost a dozen classes each. Imported normally (not by reference) by
both the engine and each talk, so a utility is emitted wherever it is used.
*/

/*
Class-driven, not prefers-color-scheme: the speaker switches the theme, the
machine does not decide it.
*/
@custom-variant dark (&:where(.dark, .dark *));

/*
The slide's main area. Replaces beercss's `main.responsive`, which was
`flex:1; padding:.5rem; overflow-x:hidden` plus `max-inline-size:75rem;
margin:0 auto` from a second rule, narrowed again to 100% by `.max`. min-h-0
is what lets a pane inside it take a percentage height without the flex item
refusing to shrink.
*/
@utility slide-main {
    flex: 1 1 0%;
    min-block-size: 0;
    inline-size: 100%;
    padding-inline: 0.5rem;
    overflow-x: hidden;
}

@utility slide-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding-inline: 0.3rem;
    color: var(--color-main);
    border-block-end: var(--rule) solid var(--color-main);
}

@utility slide-footer {
    position: absolute;
    inset-block-end: 0;
    inset-inline: 0;
}

/*
beercss gave every slide its vertical rhythm through one global rule:

  *+:is(address,article,blockquote,code,h1..h6,nav,ol,p,pre,section,table,ul,
        .grid,.row,.field,.tabs,aside,fieldset,form) { margin-block-start:1rem }

It is written in no slide, so nothing signals it -- and every element of every
slide collapses together the moment beercss leaves. This is that rule, made
explicit and scoped to the slide body.
*/
@utility slide-prose {
    & > * + * {
        margin-block-start: 1rem;
    }
}

/*
The "round avatar + handle" line, repeated 12 times across two slides of
impact-framework (demoit.md:127-143 and 481-514).
*/
@utility contact-row {
    grid-column: span 4 / span 4;
    display: flex;
    align-items: center;
    gap: 1rem;
    text-align: start;

    /* 3rem is the size the avatars render at today; measured, not guessed. */
    & img {
        block-size: 3rem;
        inline-size: 3rem;
        border-radius: 9999px;
        object-fit: cover;
    }
}

@utility card {
    border-radius: var(--radius-card);
    padding: 0.5rem;
    background-color: var(--color-surface);
}
```

- [ ] **Step 5: Écrire `styles/demoit.src.css`**

```css
/*
Build entry of the engine stylesheet. Output: handlers/resources/demoit.css,
committed and embedded in the binary.

Scanned sources are the engine's own embedded templates. A talk's slides are
scanned by its own build (hack/css.sh <talk>), because a released binary has
neither deck/layouts nor handlers/resources on disk.
*/
@layer theme, base, components, utilities;
@import "tailwindcss/theme.css" layer(theme);
@import "tailwindcss/preflight.css" layer(base);
@import "tailwindcss/utilities.css" layer(utilities);

@import "./tokens.css";
@import "./components.css";

@source "../handlers/resources/";
@source "../deck/layouts/";

/*
deck/directive builds these class names in Go (transform.go toGrid, render.go
splitAttributes), so no scanner will ever meet them in a source file. Without
these two lines the panes of every split{cols=} slide silently take the whole
row on stage, with nothing saying why -- the same failure mode transform.go
already documents for a weight outside 1..12.
*/
@source inline("col-span-{1,2,3,4,5,6,7,8,9,10,11,12}");
@source inline("h-stage-{small,medium,large,xlarge}");

/*
Legibility tokens. A display profile redefines these and only these: never a
geometry value, so a profile cannot break a layout.
*/
:root {
    --rule: 1px;
    --rule-opacity: 0.4;
    --weight-body: 400;
    --weight-thin: 100;

    /* Component chrome. Custom properties cross the shadow DOM; classes do not. */
    --dm-window-bg: #ffffff;
    --dm-chrome-bg: linear-gradient(to bottom, #edeaed 0%, #dddfdd 100%);
    --dm-chrome-border: #cbcbcb;
    --dm-tabs-bg: #f3f3f3;
    --dm-tab-bg: #ececec;
    --dm-tab-fg: #000000;
    --dm-code-fg: #212121;
    --dm-code-selection: #bfd6ff;
    /* Chroma style names, passed to /sourceCode by the source-code component. */
    --dm-code-style: "vs";
}

:root[data-display="tv"] {
    --rule: 1.5px;
    --rule-opacity: 0.7;
    --weight-body: 500;
    --weight-thin: 300;
    --color-fg: #2a2d30;
}

:root[data-display="projector"] {
    --rule: 2px;
    --rule-opacity: 1;
    --weight-body: 500;
    --weight-thin: 400;
    --color-fg: #000000;
    --color-fg-muted: #2a2d30;
}

.dark {
    --color-fg: #e6e1e6;
    --color-fg-muted: #b4b0b6;
    --color-surface: #16171a;
    --color-outline: #3a3d42;
    --color-main: var(--color-main-dark);

    --dm-window-bg: #1c1d21;
    --dm-chrome-bg: linear-gradient(to bottom, #2b2d33 0%, #23252a 100%);
    --dm-chrome-border: #3a3d42;
    --dm-tabs-bg: #23252a;
    --dm-tab-bg: #2b2d33;
    --dm-tab-fg: #e6e1e6;
    --dm-code-fg: #e6e1e6;
    --dm-code-selection: #2f4a7a;
    --dm-code-style: "github-dark";
}

/*
The stage. 1920x1080 is not a new convention: /pdf already renders at
1920x1080@2x (handlers/pdf.go) and /grid already iframes at 1920x1080
(grid.tmpl.html), while style.css mixed rem, px and vw and beercss switched
columns on width thresholds. Pinning the stage makes the geometry the same on
every 16:9 screen and letterboxes the rest.

--stage-scale is set by demoit.js on resize: dividing a length by a length in
calc() is specified but not yet implemented in browsers.

The stage also retires the vw units style.css was full of: a
`calc(100% - 14vw)` on split-view used to measure the window and now measures
the stage, which is a constant.
*/
html {
    /*
    No font-size here on purpose: 16px is both the value already in effect and
    Tailwind's preflight default, so every rem in the deck keeps its length.
    */
    background-color: var(--color-void);
}

body {
    margin: 0;
    min-block-size: 100vh;
    background-color: var(--color-void);
    color: var(--color-fg);
    font-family: var(--font-sans);
    font-weight: var(--weight-body);
}

.stage {
    position: absolute;
    inset-block-start: 50%;
    inset-inline-start: 50%;
    inline-size: 1920px;
    block-size: 1080px;
    transform: translate(-50%, -50%) scale(var(--stage-scale, 1));
    transform-origin: center;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    background-color: var(--color-surface);
}

/*
A talk that does not opt into the stage keeps the full window, as before.
*/
.no-stage {
    inline-size: 100%;
    min-block-size: 100vh;
    background-color: var(--color-surface);
}

#progression {
    position: absolute;
    inset-block-end: 0;
    inset-inline-start: 0;
    block-size: 6px;
    margin: 0;
    border-radius: 0 3px 3px 0;
    background-color: var(--color-fg);
    opacity: var(--rule-opacity);
}
```

- [ ] **Step 6: Construire la feuille et regarder ce qui en sort**

Run:
```bash
hack/css.sh engine
wc -l handlers/resources/demoit.css
grep -c "col-span-12" handlers/resources/demoit.css
grep -c "h-stage-xlarge" handlers/resources/demoit.css
```
Expected: le fichier existe, `col-span-12` et `h-stage-xlarge` présents au moins une fois chacun. Si l'un vaut 0, la ligne `@source inline(...)` correspondante est absente ou mal écrite.

- [ ] **Step 7: Écrire le test d'embarquement (il doit échouer)**

`handlers/resources_test.go` — pas d'en-tête de licence, c'est la convention des `_test.go` du repo.

```go
package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgageot/demoit/handlers"
)

// The engine stylesheet is a build artefact: it is generated by hack/css.sh
// and committed. Nothing in `go build` regenerates it, so a stale or
// regenerated-without-@source-inline copy would ship silently. These tests are
// the guard.
func TestEngineCSSIsServed(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	handlers.EngineCSS(w, httptest.NewRequest(http.MethodGet, "/demoit.css", nil))

	if got, want := w.Code, http.StatusOK; got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
	if got, want := w.Header().Get("Content-Type"), "text/css"; !strings.HasPrefix(got, want) {
		t.Errorf("got content type %q, want it to start with %q", got, want)
	}
	if w.Body.Len() == 0 {
		t.Fatal("got an empty stylesheet, want the generated CSS")
	}
}

// The classes deck/directive builds in Go are invisible to Tailwind's scanner:
// they exist in no source file. styles/demoit.src.css declares them with
// `@source inline(...)`, and dropping that declaration would leave every
// split{cols=} pane unstyled on stage with nothing to say why.
func TestEngineCSSCarriesTheClassesBuiltInGo(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	handlers.EngineCSS(w, httptest.NewRequest(http.MethodGet, "/demoit.css", nil))
	css := w.Body.String()

	for _, want := range []string{
		".col-span-1", ".col-span-4", ".col-span-8", ".col-span-12",
		".h-stage-small", ".h-stage-medium", ".h-stage-large", ".h-stage-xlarge",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("got a stylesheet without %q: is the matching @source inline() line missing from styles/demoit.src.css?", want)
		}
	}
}

// The chassis utilities are what the layouts render; a build that dropped
// styles/components.css would leave every slide unstyled.
func TestEngineCSSCarriesTheChassis(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	handlers.EngineCSS(w, httptest.NewRequest(http.MethodGet, "/demoit.css", nil))
	css := w.Body.String()

	for _, want := range []string{".slide-main", ".slide-header", ".slide-prose", ".stage"} {
		if !strings.Contains(css, want) {
			t.Errorf("got a stylesheet without %q", want)
		}
	}
}
```

- [ ] **Step 8: Lancer le test pour le voir échouer**

Run: `go test ./handlers/ -run TestEngineCSS -v`
Expected: FAIL — `undefined: handlers.EngineCSS`

- [ ] **Step 9: Écrire `handlers/resources/css.go`**

En-tête Apache-2.0 obligatoire (fichier du moteur, pas un test).

```go
/*
Copyright 2018 Google LLC
Copyright 2022 David Gageot

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handlers

import (
	_ "embed"
	"net/http"
	"strings"
	"time"
)

// engineCSS is the chassis stylesheet: tokens, the stage, the slide-* classes
// and the display profiles. It is generated by `hack/css.sh engine` and
// committed, because a talk folder cannot scan the engine's embedded templates
// and a released binary carries no styles/ directory.
//
//go:embed demoit.css
var engineCSS string

// engineCSSModTime is fixed so that the served stylesheet is cacheable and
// its ETag stable: the content only ever changes when the binary does.
var engineCSSModTime = time.Now()

// EngineCSS serves the embedded chassis stylesheet.
func EngineCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeContent(w, r, "demoit.css", engineCSSModTime, strings.NewReader(engineCSS))
}
```

Note : le fichier vit dans `handlers/resources/` mais appartient au paquet `handlers`. Vérifier que le paquet compile — si Go refuse un fichier de paquet dans un sous-dossier, déplacer `css.go` en `handlers/css.go` et changer la directive en `//go:embed resources/demoit.css`, comme `handlers/step.go:36` le fait déjà pour `index.tmpl.html`.

- [ ] **Step 10: Lancer le test pour le voir passer**

Run: `go test ./handlers/ -run TestEngineCSS -v`
Expected: PASS, 3 tests.

- [ ] **Step 11: Brancher les deux routes**

`main.go`, à côté de `/style.css` (ligne 63) :

```go
	r.HandleFunc("/demoit.css", handlers.EngineCSS).Methods("GET")
	r.HandleFunc("/tailwind.css", handlers.Static).Methods("GET")
```

- [ ] **Step 12: Vérifier que rien n'a bougé visuellement**

Run:
```bash
go build -o /tmp/demoit . && /tmp/demoit --port 8899 --shellport 9998 impact-framework &
curl -s -o /dev/null -w "demoit.css %{http_code} %{size_download}\n" http://127.0.0.1:8899/demoit.css
curl -s -o /dev/null -w "slide0 %{http_code}\n" http://127.0.0.1:8899/
```
Expected: `demoit.css 200 <taille non nulle>`, `slide0 200`. Les pages ne lient pas encore la feuille, donc le rendu est **inchangé** — c'est le but de cette tâche.

- [ ] **Step 13: Vérifications et commit**

```bash
go test ./... && gofmt -l . | grep -v '^vendor/' ; golangci-lint run -c golangci.yml
git add hack/css.sh styles handlers/resources/demoit.css handlers/resources/css.go handlers/resources_test.go main.go .gitignore
git commit -m "$(cat <<'MSG'
feat(css): build and serve a Tailwind chassis stylesheet

Adds the pinned standalone Tailwind CLI build chain and the engine's own
stylesheet: tokens, the stage, the slide-* chassis and the display
profiles. Generated by hack/css.sh and committed, because a talk folder
cannot scan the engine's embedded templates and a released binary has no
styles/ directory on disk.

Nothing links the stylesheet yet, so this changes no rendering.

The tests guard the two things a build can silently get wrong: the
classes deck/directive assembles in Go exist in no source file, so they
only reach the stylesheet through `@source inline(...)`, and the chassis
utilities only reach it through styles/components.css.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 2 : Le bloc `theme:` de `talk.yml`

Objectif : un talk déclare s'il veut la scène et le commutateur de thème. Pur Go, TDD, aucun changement visuel.

**Files:**
- Modify: `deck/talk.go:11-18`, `handlers/step.go:41-51` et `:89-115`
- Test: `deck/talk_test.go`

**Interfaces:**
- Produces: `deck.Talk.Theme` de type `deck.Theme` avec les champs `Stage bool` et `Dark bool` ; `handlers.Page.Stage bool` et `handlers.Page.Dark bool`.
- Consumes: `deck.LoadTalk(folder string) (Talk, error)` — déjà exporté (`deck/talk.go:20`).

- [ ] **Step 1: Écrire les tests (ils doivent échouer)**

À ajouter à `deck/talk_test.go`. Le helper `writeTalk` existe déjà en haut du fichier.

```go
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

// Both flags default to false, and that default is what keeps sample/ off the
// new chassis: it declares no theme block, so it gets neither the fixed stage
// nor a theme switcher its own demoit.js chrome would not follow.
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
```

- [ ] **Step 2: Lancer les tests pour les voir échouer**

Run: `go test ./deck/ -run TestLoadTalk -v`
Expected: FAIL — `talk.Theme undefined (type deck.Talk has no field or method Theme)`

- [ ] **Step 3: Ajouter le type et le champ dans `deck/talk.go`**

Au-dessus de `type Talk struct` :

```go
// Theme says which parts of the rendering chassis a talk opts into. Both
// default to false: a talk that declares no theme block keeps the full-window
// rendering and gets no theme switcher, which is what leaves sample/ alone.
type Theme struct {
	// Stage puts the slides on the fixed 1920x1080 stage.
	Stage bool `yaml:"stage"`
	// Dark renders the theme switcher and the dark palette.
	Dark bool `yaml:"dark"`
}
```

Et dans `Talk` :

```go
	Theme Theme `yaml:"theme"`
```

- [ ] **Step 4: Lancer les tests pour les voir passer**

Run: `go test ./deck/ -run TestLoadTalk -v`
Expected: PASS

- [ ] **Step 5: Exposer les deux drapeaux au template**

`handlers/step.go`, dans `Page` :

```go
	Stage       bool
	Dark        bool
```

Dans `readSteps`, avant la boucle qui construit les `steps` :

```go
	talk, err := deck.LoadTalk(folder)
	if err != nil {
		return nil, err
	}
```

Et dans le littéral `Page{...}` :

```go
			Stage:       talk.Theme.Stage,
			Dark:        talk.Theme.Dark,
```

- [ ] **Step 6: Vérifications et commit**

```bash
go build -o /dev/null . && go test ./... && gofmt -l . | grep -v '^vendor/'
git add deck/talk.go deck/talk_test.go handlers/step.go
git commit -m "$(cat <<'MSG'
feat(deck): let a talk opt into the stage and the theme switcher

talk.yml gains a theme block with two flags, both defaulting to false.
The default is the point: a talk that declares nothing keeps the
full-window rendering it has today and gets no switcher, which is how
sample/ stays off the new chassis -- its own demoit.js carries a light
window chrome in hard-coded colours that would not follow a dark page.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 3 : Les classes que le Go assemble

Objectif : `split{cols=}` et `height=` émettent des classes Tailwind. **L'API d'auteur ne change pas.**

État attendu après cette tâche : les slides `split` du deck sont **temporairement dégradées** — elles portent des classes Tailwind alors que la page ne lie pas encore la feuille. C'est un état intermédiaire assumé ; la Tâche 5 le résout.

**Files:**
- Modify: `deck/directive/transform.go:26-33`, `:84-108`, `deck/directive/render.go:119-126`
- Test: `deck/directive/transform_test.go`, `deck/directive/directive_test.go`

**Interfaces:**
- Produces: `:::split{cols=4,8 height=xlarge}` rend `<div class="grid grid-cols-12 gap-4 h-stage-xlarge">` avec `<div class="col-span-4">` et `<div class="col-span-8">` ; `layout: split` + `height: xlarge` rend `<split-view class="h-stage-xlarge">`.
- Consumes: les jetons `--height-stage-*` et les `col-span-*` de la Tâche 1.

- [ ] **Step 1: Mettre à jour le test existant (il doit échouer)**

`deck/directive/transform_test.go`, `TestSplitWithColsRendersAWeightedGrid` — remplacer les attentes et les index :

```go
	for _, want := range []string{
		`<div class="grid grid-cols-12 gap-4 h-stage-xlarge">`,
		`<div class="col-span-4">`,
		`<div class="col-span-8">`,
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
	narrow := strings.Index(got, `<div class="col-span-4">`)
	wide := strings.Index(got, `<div class="col-span-8">`)
	term := strings.Index(got, "<web-term")
	code := strings.Index(got, "<vs-code")

	if narrow > term || term > wide || wide > code {
		t.Errorf("got %q, want the col-span-4 column then the term then the col-span-8 column then vs-code", got)
	}
```

Et le message d'erreur attendu, dans `TestSplitWithColsRejectsAnOutOfRangeWeight` s'il existe (voir `transform_test.go:90` et son commentaire), doit parler de la grille sans nommer beercss. Ajouter aussi ce test :

```go
// A pane height travels as a Tailwind height utility backed by a --height-stage-*
// token, not as the old {name}-height class. The author-facing API is
// unchanged: the deck still writes height=xlarge.
func TestSplitHeightRendersAStageHeightClass(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::split{height=lg}\n::term{path=sources}\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := `<split-view class="h-stage-lg">`; !strings.Contains(got, want) {
		t.Errorf("got %q, want it to contain %q", got, want)
	}
}

// A split with no height at all must carry no class, as before: the layout's
// own default applies.
func TestSplitWithoutHeightCarriesNoClass(t *testing.T) {
	t.Parallel()

	got, errs := render(t, ":::split\n::term{path=sources}\n:::\n")

	if len(errs) != 0 {
		t.Fatalf("got errors %v, want none", errs)
	}
	if want := "<split-view>"; !strings.Contains(got, want) {
		t.Errorf("got %q, want a bare %q", got, want)
	}
}
```

- [ ] **Step 2: Lancer les tests pour les voir échouer**

Run: `go test ./deck/directive/ -v`
Expected: FAIL — les attentes portent sur `grid grid-cols-12 gap-4` et `col-span-4` alors que le code émet encore `grid` et `s4`.

- [ ] **Step 3: Réécrire `splitAttributes` dans `render.go`**

Remplacer les lignes 119-126 :

```go
// splitAttributes renders the class of a <split-view>, built from its height.
// A height name travels as a Tailwind height utility, backed by the
// --height-stage-* tokens of styles/tokens.css. Those tokens are a share of
// the slide's main area rather than an absolute length, so a pane can no
// longer be taller than the stage that holds it.
func splitAttributes(attrs map[string]string) (string, error) {
	if height := attrs["height"]; height != "" {
		return fmt.Sprintf(` class="%s"`, html.EscapeString(stageHeightClass(height))), nil
	}

	return "", nil
}

// stageHeightClass is the height utility a height name maps to. It is shared
// with the transformer, which appends the same class to a weighted grid.
func stageHeightClass(height string) string {
	return "h-stage-" + height
}
```

- [ ] **Step 4: Réécrire `toGrid` dans `transform.go`**

Remplacer le bloc des lignes 26-33 :

```go
// columnCount is the width of the grid a split{cols=} lays panes out on: a
// weight outside 1..12 has no matching "col-span-N" utility.
const columnCount = 12

// transformer rewrites a split directive that carries column weights into a
// twelve-column grid, and reports directives whose name is not in the
// catalogue.
```

Puis, dans `toGrid`, la validation et la construction (lignes 84-108) :

```go
	// A weight becomes a "col-span-N" utility verbatim, and the grid has twelve
	// columns. Anything else -- `cols=a,b`, `cols=13,4` -- is a class no
	// stylesheet defines, so the pane silently takes the whole row on stage
	// and nothing says why.
	for _, weight := range cols {
		if width, err := strconv.Atoi(weight); err != nil || width < 1 || width > columnCount {
			addError(t.ctx, node.Line, "the split directive has the column weight %q: the grid has %d columns, so every weight is a whole number from 1 to %d", weight, columnCount, columnCount)

			return
		}
	}

	class := "grid grid-cols-12 gap-4"
	if height := node.Attrs["height"]; height != "" {
		class += " " + stageHeightClass(height)
	}

	node.Name = "grid"
	node.Attrs = map[string]string{"class": class}

	for i, child := range children {
		column := NewNode("col", map[string]string{"class": "col-span-" + cols[i]}, node.FenceLength, node.Line)
		node.ReplaceChild(node, child, column)
		column.AppendChild(column, child)
	}
```

- [ ] **Step 5: Lancer les tests pour les voir passer**

Run: `go test ./deck/directive/ -v`
Expected: PASS

- [ ] **Step 6: Vérifier qu'aucun test ne mentionne encore beercss**

Run: `grep -rn beercss deck/ handlers/ main.go`
Expected: aucune sortie.

- [ ] **Step 7: Vérifications et commit**

```bash
go build -o /dev/null . && go test ./... && gofmt -l . | grep -v '^vendor/' ; golangci-lint run -c golangci.yml
git add deck/directive/
git commit -m "$(cat <<'MSG'
feat(deck): emit Tailwind grid and height classes

split{cols=} now builds a twelve-column Tailwind grid and col-span-N
panes; a height name becomes an h-stage-* utility. The author-facing API
is untouched -- a deck still writes cols=4,8 and height=xlarge.

The height tokens behind that utility are a share of the slide's main
area rather than an absolute length. The old values could not fit:
.xlarge-height was 1024px for a 1080px stage, so a header and a footer
were enough to make a pane overflow, and its fallback came from a
max-height media query -- a question about the window, which a fixed
stage makes meaningless.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 4 : Les layouts

Objectif : les 7 fichiers de `deck/layouts/` et les 2 overrides du talk rendent le châssis Tailwind. Les classes sans effet disparaissent.

État attendu : rendu intermédiaire dégradé, résolu par la Tâche 5.

**Files:**
- Modify: `deck/layouts/partials.html`, `bare.html`, `content.html`, `cover.html`, `default.html`, `quote.html`, `split.html`
- Modify: `impact-framework/.demoit/layouts/cover.html`, `default-h3.html`
- Test: `deck/layout_test.go` (vérifier que `TestEveryEmbeddedLayoutIsOverridable` passe toujours)

**Interfaces:**
- Consumes: `slide-header`, `slide-main`, `slide-prose`, `contact-row`, `card` (Tâche 1) ; `.Class`, `.Title`, `.Source`, `.Notes`, `.Talk.Logos`, `.Height`, `.Content` (inchangés).
- Produces: les classes par défaut de chaque layout, que la Tâche 6 doit connaître pour réécrire les `class:` du deck.

- [ ] **Step 1: Réécrire `deck/layouts/partials.html`**

`center-left` et `center-right` n'existent pas dans beercss — aucune règle, aucun effet. Elles ne sont pas traduites, elles sont supprimées. `max` dans un `<nav>` valait `flex:1`, d'où `flex-1`. `<img class="responsive">` valait `block-size:3rem; inline-size:3rem; object-fit:cover`, d'où `size-12 object-cover` — 3rem valent 48 px, la taille à laquelle les logos rendent effectivement sur la référence.

```html
{{ define "header" }}<header class="slide-header">
    <h2 class="flex-1 m-0 text-left">{{ .Title }}</h2>
    <div class="flex items-center gap-4">
        {{ range .Talk.Logos }}<img class="size-12 rounded-full object-cover" src="{{ . }}">
        {{ end }}
    </div>
</header>
{{ end }}

{{ define "source" }}{{ with .Source }}<div class="px-2 text-center text-fg-muted">
    <span>Source: {{ . }}</span>
</div>
{{ end }}{{ end }}

{{ define "notes" }}{{ with .Notes }}<speaker-notes>{{ . }}</speaker-notes>
{{ end }}{{ end }}
```

Note : le `<nav>` disparaît. beercss lui donnait `display:flex; align-items:center; gap:1rem`, ce que `slide-header` porte maintenant directement ; garder un `<nav>` vide de sens ajouterait un niveau pour rien. Le `<div class="grid center-right">` autour des logos devient un simple `flex` : il n'a jamais été une grille utile, `center-right` n'existant pas.

- [ ] **Step 2: Réécrire les six layouts embarqués**

Chaque layout garde son propre repli de classes — c'est la règle du repo (`deck/layout.go`, et `CLAUDE.md` : « A layout owns its default `<main>` classes as well as its markup »), et c'est ce qui rend chaque layout surchargeable par un talk.

`bare.html` :
```html
<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-prose{{ end }}">
{{ .Content }}
</main>
{{ template "notes" . }}
```

`content.html` :
```html
{{ template "header" . }}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-prose{{ end }}">
{{ .Content }}
</main>
{{ template "source" . }}{{ template "notes" . }}
```

`default.html` :
```html
{{ template "header" . }}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-prose flex flex-col justify-center text-center{{ end }}">
{{ .Content }}
</main>
{{ template "source" . }}{{ template "notes" . }}
```

`quote.html` :
```html
{{ template "header" . }}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main flex flex-col items-center justify-center text-center{{ end }}">
<blockquote class="border-l-4 border-main px-4 text-left">
{{ .Content }}
</blockquote>
</main>
{{ template "source" . }}{{ template "notes" . }}
```

`split.html` :
```html
{{ template "header" . }}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main{{ end }}">
<split-view class="{{ .Height }}">
{{ .Content }}
</split-view>
</main>
{{ template "source" . }}{{ template "notes" . }}
```

**Pourquoi les jetons s'appellent `small`/`medium`/`large`/`xlarge`** et non `sm`/`md`/`lg`/`xl` : `split.html` interpole `.Height` verbatim, et `.Height` vaut `xlarge` par défaut (`deck/frontmatter.go`). Nommer les jetons autrement aurait exigé une table de correspondance en Go pour rien. C'est déjà ce que la Tâche 1 a écrit dans `styles/tokens.css` et dans l'`@source inline`.

`cover.html` :
```html
<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main flex flex-col items-center justify-center text-center{{ end }}">
{{ .Content }}
<div class="mt-12 flex items-center justify-center gap-12">
    {{ range .Talk.Logos }}<img class="h-[25rem] w-[25rem] object-contain" src="{{ . }}" loading="lazy">
    {{ end }}
</div>
</main>
{{ template "notes" . }}
```

Note : les `<div class="large-space">` empilés et la grille `s12 m6 l6` deviennent un `flex gap-12`. `xsmall-height`/`xsmall-width` valaient `25rem` (surcharge du talk, `style.css:42-47`), d'où `h-[25rem] w-[25rem]` ; `object-contain` plutôt que `cover` parce qu'il s'agit de logos, et que le `object-fit:cover` que beercss imposait les recadrait.

- [ ] **Step 3: Réécrire les deux overrides du talk**

`impact-framework/.demoit/layouts/cover.html` :
```html
{{/* Same fallback the embedded cover layout carries: a talk that overrides a
     layout owns its default classes too, since the engine no longer fills
     .Class in before the template runs. */}}<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main flex flex-col items-center justify-center text-center{{ end }}">
    <div class="card px-8 py-4">
{{ .Content }}
    </div>
    <div class="mt-12 flex items-center justify-center gap-12">
        <img class="h-[25rem] w-[25rem] object-contain" src="/images/logo.jpg" loading="lazy">
        <img class="h-[25rem] w-[25rem] object-contain" src="/images/zatsit_logo_noir.svg" loading="lazy"/>
    </div>
</main>
{{ template "notes" . }}
```

`impact-framework/.demoit/layouts/default-h3.html` — même contenu que `default.html` avec un `<h3>` :
```html
<header class="slide-header">
    <h3 class="flex-1 m-0 text-left">{{ .Title }}</h3>
    <div class="flex items-center gap-4">
        {{ range .Talk.Logos }}<img class="size-12 rounded-full object-cover" src="{{ . }}">
        {{ end }}
    </div>
</header>
{{/* slide.Class arrives empty when the slide names no class: key, so every
     layout -- embedded or talk-added -- carries its own fallback for that case.
     This one repeats the fallback the embedded "default" layout uses. A slide
     that does declare a class: key still gets it verbatim, replacing this
     fallback wholesale, the same replace semantics every other layout
     follows. */}}
<main class="{{ if .Class }}{{ .Class }}{{ else }}slide-main slide-prose flex flex-col justify-center text-center{{ end }}">
{{ .Content }}
</main>
{{ template "source" . }}{{ template "notes" . }}
```

- [ ] **Step 4: Reconstruire la feuille du moteur**

Les layouts ont changé, donc les classes que le scanner voit ont changé.

Run:
```bash
hack/css.sh engine
grep -cE "\.slide-header|\.slide-main|\.slide-prose" handlers/resources/demoit.css
go test ./... 2>&1 | tail -20
```
Expected: les trois classes du châssis présentes, tous les tests verts.

- [ ] **Step 5: Vérifier que la surchargeabilité tient toujours**

Run: `go test ./deck/ -run TestEveryEmbeddedLayoutIsOverridable -v`
Expected: PASS. Ce test énumère les layouts embarqués ; il doit continuer à passer sans modification.

- [ ] **Step 6: Vérifications et commit**

```bash
go build -o /dev/null . && go test ./... && gofmt -l . | grep -v '^vendor/'
git add deck/layouts handlers/resources/demoit.css impact-framework/.demoit/layouts
git commit -m "$(cat <<'MSG'
feat(layouts): render the Tailwind chassis

Rewrites the seven embedded layouts and the two impact-framework
overrides onto the slide-header / slide-main / slide-prose chassis.
Every layout keeps carrying its own fallback class list, so a talk can
still override any of them.

Three classes are dropped rather than translated: center-left,
center-right and padding have no rule at all in beercss, so the eleven
places that used them never did anything. Two are translated to what
they actually did rather than what they read like: `max` inside a nav
resolved to flex:1, and `responsive` on an img resolved to a cropped
3rem square, which is why the header logos were thumbnails.

Height names now travel as written -- h-stage-xlarge, not h-stage-xl --
because split.html interpolates the frontmatter value verbatim and the
deck writes `xlarge`.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 5 : Le basculement

Objectif : retirer beercss, poser la scène, lier les nouvelles feuilles. **Première vérification visuelle complète.**

**Files:**
- Modify: `handlers/resources/index.tmpl.html` (réécriture complète)
- Modify: `handlers/resources/grid.tmpl.html`, `handlers/resources/speakernotes.tmpl.html`
- Create: `impact-framework/.demoit/tokens.css`, `components.css`, `tailwind.src.css`, `tailwind.css`
- Modify: `impact-framework/.demoit/talk.yml`, `impact-framework/.demoit/js/demoit.js`

**Interfaces:**
- Consumes: `Page.Stage`, `Page.Dark` (Tâche 2) ; `/demoit.css`, `/tailwind.css` (Tâche 1) ; le châssis des layouts (Tâche 4).
- Produces: `window.demoitTheme` — objet exposé par `demoit.js` avec `get()`, `set(theme)`, `toggle()`, `profile()`, `cycleProfile()` ; la classe `.dark` et l'attribut `data-display` sur `documentElement` ; l'émission `{theme, display}` sur `BroadcastChannel("demoit_nav")`.

- [ ] **Step 1: Écrire les trois fichiers CSS du talk**

`impact-framework/.demoit/tokens.css` — copie de `styles/tokens.css`, avec la palette du talk. **Uniquement des `@theme`.**

```css
/*
Design tokens of this talk. Copied from styles/tokens.css and free to diverge:
this is where the talk's palette lives.

Must contain @theme blocks and nothing else -- tailwind.src.css imports it with
`theme(reference)`, and Tailwind rejects a referenced file carrying any other
rule.
*/
@theme {
    --font-sans: "Poppins", ui-sans-serif, system-ui, sans-serif;
    --font-mono: "Roboto Mono", ui-monospace, monospace;

    --color-main: #0f15fd;
    --color-main-dark: #7c81ff;

    --color-fg: #3d4043;
    --color-fg-muted: #6b7075;
    --color-surface: #ffffff;
    --color-outline: #d5d7da;
    --color-void: #000000;

    --height-stage-small: 40%;
    --height-stage-medium: 65%;
    --height-stage-large: 85%;
    --height-stage-xlarge: 100%;

    --radius-card: 2rem;
}
```

`impact-framework/.demoit/components.css` — copie de `styles/components.css`, à l'identique.

`impact-framework/.demoit/tailwind.src.css` :

```css
/*
Build entry of this talk's stylesheet. Output: tailwind.css, committed and
served on /tailwind.css.

Utilities only: the preflight and the theme variables ship once, in the
engine's own /demoit.css.

The first import is not optional. Referencing only this talk's tokens leaves
Tailwind's own theme out of scope, and every utility backed by --spacing or
--color-* -- gap-4, text-white, px-8 -- is then silently not generated.
*/
@import "tailwindcss/theme.css" theme(reference);
@import "tailwindcss/utilities.css" layer(utilities);
@import "./tokens.css" theme(reference);
@import "./components.css";

/*
Scans the whole talk folder: demoit.md, demoit-en.html, and .demoit/layouts.
An explicit @source does reach into a dot-directory, unlike automatic source
detection.
*/
@source "../";
@source inline("col-span-{1,2,3,4,5,6,7,8,9,10,11,12}");
@source inline("h-stage-{small,medium,large,xlarge}");
```

- [ ] **Step 2: Construire la feuille du talk et la regarder**

Run:
```bash
hack/css.sh impact-framework
wc -l impact-framework/.demoit/tailwind.css
grep -c "box-sizing" impact-framework/.demoit/tailwind.css
grep -cE "^\s*:root" impact-framework/.demoit/tailwind.css
```
Expected: fichier non vide ; `box-sizing` à 0 (pas de preflight en double) ; `:root` à 0 (pas de variables de thème en double).

- [ ] **Step 3: Activer la scène et le thème dans `talk.yml`**

```yaml
layout: default
theme:
  stage: true
  dark: true
logos:
  - /images/zatsit_logo.svg
  - /images/logo.jpg
```

- [ ] **Step 4: Réécrire `handlers/resources/index.tmpl.html`**

```html
<!doctype html>
<html lang=en>
	<head>
		<meta charset="utf-8">
		<title>Demo {{ .CurrentStep }}/{{ .StepCount }}</title>
		<link rel="stylesheet" href="/demoit.css?hash={{ "demoit.css" | hash }}">
		<link rel="stylesheet" href="/tailwind.css?hash={{ "tailwind.css" | hash }}">
		<link rel="stylesheet" href="/style.css?hash={{ "style.css" | hash }}">

		<script>
			const CurrentStep = {{ .CurrentStep }};
			const StepCount = {{ .StepCount }};
			const NextURL = '{{ .NextURL }}';
			const PrevURL = '{{ .PrevURL }}';
			const ThemeEnabled = {{ .Dark }};
		</script>
		{{ if .Dark }}<script>
			// Runs before the first paint on purpose: setting the theme from
			// demoit.js would flash a white slide on every navigation, because
			// each slide is a full page load.
			(function () {
				const params = new URLSearchParams(location.search);
				let theme = params.get('theme');
				let display = params.get('display');
				try {
					theme = theme || localStorage.getItem('demoit_theme');
					display = display || localStorage.getItem('demoit_display');
				} catch (e) {
					// A file:// origin or a locked-down browser throws on
					// localStorage; the default theme is a fine answer.
				}
				if (theme === 'dark') {
					document.documentElement.classList.add('dark');
				}
				document.documentElement.dataset.display = display || 'screen';
			})();
		</script>{{ end }}
	</head>
	<body>
		<div id="app" class="{{ if .Stage }}stage{{ else }}no-stage{{ end }}">
			{{ .HTML }}
			<nav-arrows previous="{{ .PrevURL }}" next="{{ .NextURL }}"></nav-arrows>
			<footer class="slide-footer">
				<h5 id="progression" style="width: calc(1920px * {{ .CurrentStep }} / {{ .StepCount }})"></h5>
				<div class="px-4 pb-1 text-right text-fg-muted" style="font-weight: var(--weight-thin)">Slides framework adapted from <span class="font-bold">demoit</span></div>
			</footer>
		</div>
	</body>
	<script type="module" src="/js/demoit.js?hash={{ "js/demoit.js" | hash }}"></script>
	{{ if .DevMode }}<script src="/livereload.js"></script>{{ end }}
</html>
```

Changements et leur raison :
- beercss, `beer.min.js`, `material-dynamic-colors` retirés — les deux scripts n'étaient appelés nulle part.
- `<dialog id="maximized">` retiré — `#maximized` n'était référencé nulle part ; le bouton vert des fenêtres agit sur son propre shadow DOM.
- `class="light poppins-regular"` sur `<body>` retiré : la police vient du jeton `--font-sans` posé par `demoit.css`.
- `<div class="responsive max">` devient `.stage` ou `.no-stage` selon l'opt-in.
- La barre de progression est en `1920px` et non `100vw` : sur la scène, la largeur de référence est celle de la scène, pas de la fenêtre.
- Le `<div class="grid"><div class="s12">` autour du footer disparaît — deux niveaux de grille pour empiler deux blocs.
- L'ordre des feuilles compte : `demoit.css` (preflight + jetons), puis `tailwind.css` (utilitaires du talk), puis `style.css` (surcharges du talk, qui doivent gagner).

- [ ] **Step 5: Ajouter l'échelle de scène et le thème à `demoit.js`**

En tête de `impact-framework/.demoit/js/demoit.js`, avant `class BaseHTMLElement` :

```js
// The stage is 1920x1080 CSS pixels, scaled to fit the window. Dividing a
// length by a length in calc() is specified but not implemented in browsers,
// so the factor is computed here.
//
// Scaling with transform rather than by shrinking a root font size is
// deliberate: it keeps the tty iframe 1920 logical pixels wide, so xterm.js
// always measures the same grid and a demo does not reflow according to the
// room.
function fitStage() {
    const stage = document.querySelector('.stage');
    if (!stage) {
        return;
    }
    const scale = Math.min(window.innerWidth / 1920, window.innerHeight / 1080);
    document.documentElement.style.setProperty('--stage-scale', scale);
}

window.addEventListener('resize', fitStage);
fitStage();
```

Puis, après la définition du `channel` existant (`demoit.js:577`), la mécanique de thème :

```js
const DISPLAY_PROFILES = ['screen', 'tv', 'projector'];

// Theme and display profile live on documentElement, are remembered per
// browser, and are broadcast so that the speaker-notes window and /grid
// follow. They are read back before the first paint by the inline script in
// index.tmpl.html, which is what keeps a navigation from flashing white.
const demoitTheme = {
    get() {
        return document.documentElement.classList.contains('dark') ? 'dark' : 'light';
    },

    profile() {
        return document.documentElement.dataset.display || 'screen';
    },

    set(theme, display) {
        document.documentElement.classList.toggle('dark', theme === 'dark');
        document.documentElement.dataset.display = display;
        try {
            localStorage.setItem('demoit_theme', theme);
            localStorage.setItem('demoit_display', display);
        } catch (e) {
            // Nothing to remember is better than nothing to render.
        }
        document.dispatchEvent(new CustomEvent('demoit:theme', { detail: { theme, display } }));
    },

    toggle() {
        const next = this.get() === 'dark' ? 'light' : 'dark';
        this.set(next, this.profile());
        channel.postMessage({ theme: next, display: this.profile() });
    },

    cycleProfile() {
        const at = DISPLAY_PROFILES.indexOf(this.profile());
        const next = DISPLAY_PROFILES[(at + 1) % DISPLAY_PROFILES.length];
        this.set(this.get(), next);
        channel.postMessage({ theme: this.get(), display: next });
    },
};

window.demoitTheme = demoitTheme;

if (typeof ThemeEnabled !== 'undefined' && ThemeEnabled) {
    document.addEventListener('keydown', event => {
        // nav-arrows already binds the arrows, PageUp/PageDown and space;
        // t and d are free.
        if (event.key === 't') {
            demoitTheme.toggle();
        }
        if (event.key === 'd') {
            demoitTheme.cycleProfile();
        }
    });
}
```

Et dans le `channel.onmessage` existant (`demoit.js:597`), avant le test sur `destinationSlideId` :

```js
    if (e.data.hasOwnProperty("theme")) {
        demoitTheme.set(e.data.theme, e.data.display || demoitTheme.profile());
    }
```

- [ ] **Step 6: Reconstruire, lancer, et regarder les 18 slides**

Run:
```bash
hack/css.sh engine && hack/css.sh impact-framework
go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shootall.sh <scratchpad>/after-t5/png-1920 http://127.0.0.1:8899 17 18
```

Puis comparer chaque `after-t5/png-1920/slide-NN.png` à `before/png-1920/slide-NN.png`. À ce stade le deck porte encore ses classes beercss (Tâche 6), donc **on n'attend pas l'égalité** : on vérifie que le châssis est là et que rien ne plante — scène centrée et letterboxée, header avec titre et logos, footer avec la barre de progression, navigation au clavier, `t` et `d` qui commutent.

- [ ] **Step 7: Vérifier les trois autres routes**

Run:
```bash
curl -s -o /dev/null -w "grid %{http_code}\n" http://127.0.0.1:8899/grid
curl -s -o /dev/null -w "notes %{http_code}\n" http://127.0.0.1:8899/speakernotes
curl -s -o /dev/null -w "pdf %{http_code}\n" --max-time 300 http://127.0.0.1:8899/pdf
```
Expected: trois fois 200.

- [ ] **Step 8: Vérifications et commit**

```bash
go build -o /dev/null . && go test ./... && gofmt -l . | grep -v '^vendor/'
git add handlers/resources impact-framework/.demoit
git commit -m "$(cat <<'MSG'
feat: drop beercss and put the slides on a fixed stage

index.tmpl.html loses the beercss stylesheet, beer.min.js,
material-dynamic-colors and the maximized dialog. Each removal was
checked rather than assumed: neither script is called anywhere in either
talk's demoit.js, and #maximized is referenced nowhere -- a window's
green dot acts on its own shadow DOM.

In their place: the 1920x1080 stage, scaled to fit by transform, and a
theme read from localStorage by an inline script that runs before the
first paint. Both are opt-in through talk.yml, so sample/ keeps its
full-window rendering.

The theme sets a class on documentElement and broadcasts on the channel
the speaker-notes window already listens to, so the notes window and
/grid follow rather than staying white while the stage goes dark.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 6 : Les classes de `demoit.md`

Objectif : les 18 slides du deck Markdown portent des classes Tailwind, vérifiées une par une contre la référence.

**Files:**
- Modify: `impact-framework/demoit.md`

**Interfaces:**
- Consumes: le châssis (Tâche 4), les jetons et utilitaires (Tâche 1), `col-span-*` / `h-stage-*` (Tâches 1 et 3).

**La table de correspondance à appliquer** — c'est la référence de cette tâche, établie depuis les règles réelles de beercss 3.7.8 :

| beercss | Tailwind |
|---|---|
| `grid` | `grid grid-cols-12 gap-4` |
| `s1` … `s12` | `col-span-1` … `col-span-12` |
| `m6`, `l6` | supprimés — la scène n'a qu'une largeur ; `s12 m6 l6` devient `col-span-6` |
| `large-space` (div vide) | `h-12` |
| `medium-space` (div vide) | `h-8` |
| `small-space` (div vide) | `h-4` |
| `grid large-space` (sur la grille) | `grid grid-cols-12 gap-8` — sur une grille, `large-space` valait `gap:2rem`, pas une hauteur |
| `center-align` | `text-center justify-center` |
| `left-align` | `text-left justify-start` |
| `right-align` | `text-right justify-end` |
| `middle-align` | `flex items-center` |
| `top-align` | `flex items-start` |
| `middle` | `top-1/2 -translate-y-1/2` |
| `center` | `left-1/2 -translate-x-1/2` |
| `bottom` | `bottom-0` |
| `absolute` / `fixed` | `absolute` / `fixed` |
| `circle` | `rounded-full` |
| `round` | `rounded-card` |
| `extra` | `size-14` (3.5rem carré) |
| `transparent` | `bg-transparent shadow-none text-inherit` |
| `border` (sur `table`) | `border-collapse` + bordures explicites (voir slide 5) |
| `no-padding` / `no-margin` | `p-0` / `m-0` |
| `responsive` (sur bloc) | `w-full` |
| `responsive` (sur `img`) | `size-12 object-cover` |
| `max` (hors `nav`) | `max-w-full` |
| `main` | supprimé — `.main{height:100%}` du talk, remplacé par le flex du châssis |
| `large-height`, `xlarge-height` | `h-stage-large`, `h-stage-xlarge` |
| `large-text` | `text-base` |
| `xlarge-text` | `text-[1.6rem]` |
| `small` (sur `<a>`) | supprimé — sans effet sur un `<a>` |
| `center-left`, `center-right`, `padding` | supprimés — aucune règle |

- [ ] **Step 1: Slide 1 (lignes 1-7) — la couverture**

Frontmatter `class: responsive max center-align title` → `class: slide-main flex flex-col items-center justify-center text-center title`.
Le `<h1 class=" medium no-padding center-align top-align" style="color: white">` : `medium` sur un `h1` valait `font-size:2.8125rem`... non — `h1.small` et `h1.large` existent, `h1.medium` non. `medium` sur un `h1` tombait donc sur `.medium{block-size:2.5rem; inline-size:2.5rem; padding:0}` ou rien selon l'ordre : à vérifier sur la référence. Écrire `<h1 class="m-0 text-center" style="color: white">`, comparer à `before/png-1920/slide-00.png`, et n'ajuster la taille que si elle diffère.
`<h3 style="color: white">` : inchangé.

- [ ] **Step 2: Slide 2 (lignes 9-28)**

Les cinq `<div class="large-space"></div>` deviennent cinq `<div class="h-12"></div>`. Le dernier, ligne 28, aussi. Aucune autre classe.

Vérifier contre `slide-01.png` : c'est la slide où la règle de rythme vertical de beercss se voyait le plus (deux `##` séparés). Si les deux titres se collent, `slide-prose` n'est pas appliqué — le layout par défaut du talk est `default`, qui le porte.

- [ ] **Step 3: Slide 3 (lignes 30-79)**

`class: main responsive large-height middle-align` → `class: slide-main slide-prose flex items-center justify-center`.
Le bloc commenté (lignes 36-66) : **le laisser tel quel**, il documente l'état pré-migration Markdown. Il contient des `---` ? Non, vérifié — sinon il couperait une slide.
`<blockquote style="font-size: 2.0rem;">` : ajouter `class="border-l-4 border-main px-4 text-left"`, puisque beercss donnait à `blockquote` `padding:1rem; border-inline-start:.25rem solid var(--primary)` et que le talk le surchargeait avec `border-inline-start-color: var(--color-main)`.
`<div class="center-align middle-align">` (ligne 77) → `<div class="flex items-center justify-center text-center">`.

- [ ] **Step 4: Slide 4 (lignes 81-147) — la slide de contacts**

`class: main responsive xlarge-height center-align` → `class: slide-main h-stage-xlarge text-center`.

La grille (ligne 114) `<div class="grid large-space">` → `<div class="grid grid-cols-12 gap-8">`.
`<div class="m6 s6 l6">` ×2 → `<div class="col-span-6">`.
`<div class="round extra">` → `<div class="rounded-card">` — `extra` valait un carré de 3.5rem, ce qui n'a pas de sens autour d'une image de 15rem ; à confirmer sur `slide-03.png`.
`<div class="center-align middle">` → `<div class="text-center">` — `middle` positionnait en absolu à 50%, sans `absolute` c'était sans effet.
`<div class="s12"><h2 class="left-align">` → `<div class="col-span-12"><h2 class="text-left">`.

Les 12 lignes de contact, lignes 128-133 et 137-142, passent de deux div à une :
```html
  <div class="contact-row"><img src="/images/twitter.jpg"/><h5>@ldussart</h5></div>
```
`contact-row` porte `col-span-4`, le flex, le gap et l'arrondi de l'avatar. Les deux div vides `<div class="s1"></div><div class="s3"></div>` (lignes 134-135) disparaissent : elles ne servaient qu'à compléter une ligne de grille.

`<div class="medium-space"></div>` → `<div class="h-8"></div>`.
`<div class="s12">` (ligne 145) → `<div class="col-span-12">` — mais il est **hors** de la grille, donc `col-span-12` n'y fait rien : le remplacer par `<div>`.

- [ ] **Step 5: Slide 5 (lignes 149-196) — le tableau**

`class: main responsive large-height center-align middle-align` → `class: slide-main flex items-center justify-center`.

`<table class="border small-space xlarge-text">` : beercss donnait `table{inline-size:100%; border-spacing:0; font-size:.875rem}`, `table.border>tbody>tr:not(:last-child)>td, thead>tr>th{border-block-end:.0625rem solid var(--outline)}`, `table.small-space :is(th,td){padding:0}`, et le talk `.xlarge-text{font-size:1.6rem}`. D'où :

```html
<table class="w-full border-collapse text-[1.6rem] [&_td]:p-0 [&_th]:p-0 [&_tbody_tr:not(:last-child)_td]:border-b [&_tbody_tr:not(:last-child)_td]:border-outline">
```

- [ ] **Step 6: Slide 6 (lignes 198-215)**

`class: main responsive max xlarge-height` → `class: slide-main h-stage-xlarge`.
`<div class="grid ">` → `<div class="grid grid-cols-12 gap-4 h-full">`.
`<div class="s4">` → `<div class="col-span-4">`, `<div class="s8">` → `<div class="col-span-8">`.
`<a class="large-text small center large-space" ...>` : `large-text` valait `font-size:1rem`, `small` et `center` et `large-space` n'ont pas d'effet sur un `<a>` hors contexte flex — d'où `<a class="block text-base" ...>`. À confirmer sur `slide-05.png`.
`<img class="max large-height center" ...>` → `<img class="mx-auto max-w-full h-stage-large object-contain" ...>`.
Noter le `</h6>` orphelin ligne 207 : balise fermante sans ouvrante, laissée telle quelle (le HTML brut du deck n'est pas validé) — ou supprimée, c'est du bruit. **La supprimer.**

- [ ] **Step 7: Slide 7 (lignes 217-231)**

`class: main responsive max large-height` → `class: slide-main h-stage-large`.
`<div class="grid">` → `<div class="grid grid-cols-12 gap-4 h-full">`.
`<div class="s12 center-align">` → `<div class="col-span-12 text-center">`.
`<div class="s12 center-align responsive large-height ">` → `<div class="col-span-12 w-full h-stage-large text-center">`.

- [ ] **Step 8: Slide 8 (lignes 233-269)**

`class: main responsive max large-height center-align` → `class: slide-main h-stage-large text-center`.
Les `<div class="large-space">` → `<div class="h-12">`.
`<h3 class="left-align">` → `<h3 class="text-left">`.
`<div class="absolute center" >` → `<div class="absolute left-1/2 -translate-x-1/2">`.
`<h4 class="left-align">` ×3 → `<h4 class="text-left">`.

- [ ] **Step 9: Slides 9, 10, 11 (les `layout: split`)**

Aucune classe beercss : uniquement `height: xlarge` en frontmatter et des directives. **Rien à changer** — la Tâche 3 a fait le travail. Vérifier sur `slide-08.png` à `slide-10.png` que les panneaux occupent la scène et que le mermaid rend.

- [ ] **Step 10: Slide 12 (lignes 347-379) — la grille écrite à la main**

```
:::grid{class="grid grid-cols-12 gap-4 h-stage-xlarge"}
:::col{class="col-span-12"}
#### Pour pouvoir commencer à *mesurer*, il faut d'abord *observer* et *récolter* de la data.
:::

:::col{class="col-span-4 h-stage-xlarge"}
...
:::

:::col{class="col-span-8 h-stage-xlarge"}
::vscode{path=sources}
:::
:::
```
Dans le corps du premier `col` : `<p class="xlarge-text">` → `<p class="text-[1.6rem]">`, `<ul class="xlarge-text">` → `<ul class="list-disc pl-8 text-[1.6rem]">` — beercss faisait `:not(nav)>:is(ul,ol){all:revert}`, ce qui rendait les puces ; le preflight de Tailwind les enlève, donc il faut les redemander.
Le `<text>` (ligne ~356) n'est pas un élément HTML — le remplacer par `<div>`.

- [ ] **Step 11: Slides 13 à 16 (lignes 380-458)**

Slides 13 : rien (uniquement des directives).
Slides 14, 15, 16 : `class: responsive max xlarge-height` → `class: slide-main h-stage-xlarge`. Le grand bloc commenté de la slide 14 (lignes 399-429) reste tel quel.

- [ ] **Step 12: Slide 17 (lignes 460-479)**

`class: main responsive max large-height center-align` → `class: slide-main h-stage-large text-center`.
`large-space` → `h-12`. `absolute center` → `absolute left-1/2 -translate-x-1/2`. `h4 left-align` → `text-left`.
`<div class="absolute bottom center">` → `<div class="absolute bottom-0 left-1/2 -translate-x-1/2">`.

- [ ] **Step 13: Slide 18 (lignes 481-514)**

`class: main responsive xlarge-height center-align` → `class: slide-main h-stage-xlarge text-center`.
`<div class="responsive xlarge-height">` → `<div class="w-full h-stage-xlarge">`.
`<div class="absolute bottom">` → `<div class="absolute bottom-0 inset-x-0">`.
`<div class="grid  large-space">` → `<div class="grid grid-cols-12 gap-8">`.
Les 12 lignes de contact deviennent des `contact-row`, comme à la slide 4. Les deux div vides disparaissent.

- [ ] **Step 14: Reconstruire, recapturer, comparer les 18**

Run:
```bash
hack/css.sh impact-framework
go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shootall.sh <scratchpad>/after-t6/png-1920 http://127.0.0.1:8899 17 18
```
Comparer slide par slide à `before/png-1920/`. Critère : même intention, même hiérarchie, rien de perdu. Les différences attendues et acceptables sont le letterboxing hors 16:9, les tailles de logos (beercss les recadrait en carrés de 3rem) et les hauteurs de panneaux (recalibrées en parts de la zone principale).

- [ ] **Step 15: Vérifier qu'aucune classe beercss ne subsiste**

Run:
```bash
grep -nE 'class="[^"]*\b(s[0-9]+|m6|l6|large-space|medium-space|small-space|center-align|left-align|right-align|middle-align|top-align|responsive|circle|transparent|no-padding|no-margin|large-height|xlarge-height|large-text|xlarge-text|center-left|center-right)\b' impact-framework/demoit.md
```
Expected: aucune sortie **hors des blocs commentés**, qui sont conservés comme référence pré-migration.

- [ ] **Step 16: Commit**

```bash
git add impact-framework/demoit.md impact-framework/.demoit/tailwind.css
git commit -m "$(cat <<'MSG'
refactor(impact-framework): move the deck onto Tailwind classes

Rewrites the classes of all 18 slides, checked one by one against
screenshots of the pre-migration render.

Some substitutions are not what the class name suggested. `large-space`
on a grid was a 2rem gap, not a spacer height. `middle` without
`absolute` did nothing. `extra` around a 15rem image was a 3.5rem
square. `center-left`, `center-right` and `padding` have no rule in
beercss at all. And `s12 m6 l6` collapses to a single col-span-6,
because a fixed stage has one width.

Two things the preflight takes away and the slides have to ask for
again: list bullets, which beercss restored with `:not(nav)>:is(ul,ol)
{all:revert}`, and the vertical rhythm between block elements, now the
explicit slide-prose utility.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 7 : Le deck anglais

Objectif : `demoit-en.html` rend correctement, sans beercss.

**Files:**
- Modify: `impact-framework/demoit-en.html`

- [ ] **Step 1: Capturer la référence du deck anglais**

Run:
```bash
/tmp/demoit-before --port 8898 --shellport 9997 --locale en impact-framework &
bash <scratchpad>/shootall.sh <scratchpad>/before-en/png-1920 http://127.0.0.1:8898 17 18
```

Note : `/tmp/demoit-before` est le binaire construit avant la migration, conservé par la Tâche 0. S'il n'existe plus, le reconstruire depuis `d4af1bd` dans un second worktree.

- [ ] **Step 2: Appliquer la table de correspondance**

Le deck HTML n'a pas de layout pour absorber les classes à sa place : chaque slide porte son propre `<header><nav>` et son `<main>`. Appliquer la même table qu'à la Tâche 6, plus :
- `<header><nav>` → `<header class="slide-header">` sans le `<nav>`, comme dans `partials.html`.
- `<main class="responsive max">` → `<main class="slide-main slide-prose">`.
- `<div class="circle transparent s6"><img class="responsive">` → `<img class="size-12 rounded-full object-cover">`.
- `<h2 class="max center-left">` → `<h2 class="flex-1 m-0 text-left">`.

- [ ] **Step 3: Reconstruire, recapturer, comparer**

Run:
```bash
hack/css.sh impact-framework
/tmp/demoit-after --port 8899 --shellport 9998 --locale en impact-framework &
bash <scratchpad>/shootall.sh <scratchpad>/after-en/png-1920 http://127.0.0.1:8899 17 18
```
Comparer à `before-en/png-1920/`.

- [ ] **Step 4: Vérifier qu'aucune classe beercss ne subsiste**

Run: la même commande `grep` que la Tâche 6 Step 15, sur `impact-framework/demoit-en.html`.

- [ ] **Step 5: Commit**

```bash
git add impact-framework/demoit-en.html impact-framework/.demoit/tailwind.css
git commit -m "$(cat <<'MSG'
refactor(impact-framework): move the English deck onto Tailwind classes

demoit-en.html is served by --locale en and leans on beercss across 584
lines, with no layout to absorb its classes: every slide carries its own
header and main. Ported by hand against screenshots of the pre-migration
render, using the same mapping table as the Markdown deck.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 8 : Nettoyer `style.css` et vendorer Poppins

Objectif : ce qui est devenu un jeton disparaît de `style.css` ; Poppins vient du disque, plus du CDN.

**Files:**
- Modify: `impact-framework/.demoit/style.css`
- Create: `impact-framework/.demoit/fonts/poppins-{100,400,500,700}.woff2`
- Delete: `impact-framework/.demoit/poppins.css`

- [ ] **Step 1: Récupérer les quatre graisses**

Run:
```bash
mkdir -p impact-framework/.demoit/fonts
curl -sSfL -o /tmp/poppins.css "https://fonts.googleapis.com/css2?family=Poppins:wght@100;400;500;700&display=swap" -H "User-Agent: Mozilla/5.0"
grep -oE 'https://fonts.gstatic.com[^)]*\.woff2' /tmp/poppins.css | sort -u
```
Puis télécharger chaque URL dans `impact-framework/.demoit/fonts/`, en nommant les fichiers par graisse. L'`User-Agent` compte : sans lui Google sert du `.ttf` legacy.

- [ ] **Step 2: Déclarer les `@font-face` et retirer ce qui est devenu jeton**

Dans `impact-framework/.demoit/style.css`, remplacer le contenu par :

```css
/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
This talk's own layer. Loaded after /demoit.css and /tailwind.css so that it
wins. Everything that became a design token now lives in .demoit/tokens.css:
the palette, the fonts, the pane heights and the .poppins-* helpers, whose
seventeen classes were a font-weight each.

Poppins is served from .demoit/fonts through the /fonts/ route rather than
from Google's CDN: a talk must render with no network.
*/
@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 100;
    font-display: swap;
    src: url("/fonts/poppins-100.woff2") format("woff2");
}
@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 400;
    font-display: swap;
    src: url("/fonts/poppins-400.woff2") format("woff2");
}
@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 500;
    font-display: swap;
    src: url("/fonts/poppins-500.woff2") format("woff2");
}
@font-face {
    font-family: "Poppins";
    font-style: normal;
    font-weight: 700;
    font-display: swap;
    src: url("/fonts/poppins-700.woff2") format("woff2");
}

/* The cover's gradient plate. */
.title {
    background: linear-gradient(to bottom, var(--color-main) 10%, var(--color-surface) 90%);
}

/* An emphasis in this deck means "the accent colour", not italics. */
em {
    font-style: normal;
    color: var(--color-main);
}

a,
a:visited {
    color: var(--color-main);
    text-decoration: none;
}

/* mermaid renders asynchronously; hide the source until it has. */
.mermaid {
    visibility: hidden;
    block-size: 100%;
}

.mermaid[data-processed="true"] {
    visibility: visible;
}

.mermaid svg {
    min-inline-size: 100%;
    min-block-size: 100%;
}

/* The QR code overlay, positioned against the stage. */
#qr {
    position: absolute;
    inset-inline-end: 1rem;
    inset-block-end: 1rem;
    inline-size: 15rem;
    block-size: auto;
}

/* Speaker Notes window. */
#speaker-notes {
    display: table;
    position: fixed;
    inline-size: 80%;
    block-size: 80%;
    inset: 0 auto auto 0;
    z-index: 999;
    margin: 10%;
    padding: 2rem;
    border: 4px solid var(--color-outline);
    border-radius: 2rem;
    background-color: var(--color-surface);
    color: var(--color-fg);
}

#speaker-notes-contents {
    display: table-cell;
    vertical-align: middle;
    line-height: normal;
}

#speaker-notes-progress,
#current-slide-title {
    position: fixed;
    inset-block-start: 1rem;
    block-size: 3.5rem;
    padding: 0.5rem;
    font-size: 3rem;
    font-weight: bold;
    text-align: center;
}

#speaker-notes-progress {
    inset-inline-start: 1rem;
    background-color: var(--color-main);
    color: var(--color-surface);
}

#current-slide-title {
    inset-inline-start: 12rem;
    background-color: var(--color-fg);
    color: var(--color-surface);
}
```

Ce qui disparaît et pourquoi :
- `html, body, #top { font-size: 3vw; ... }` — règle morte : beercss est chargé après `style.css` et son propre `html{font-size:var(--size)}` gagne, donc `rem` vaut déjà 16 px. `#top` n'existe plus dans le template depuis longtemps.
- `#app { background-color: white }`, `body { display:flex; ... }`, `html { width/height }` — le châssis les porte.
- `.xsmall-height`, `.xsmall-width`, `.xlarge-height`, `.xlarge-text`, `.full-width`, `.main` — jetons et utilitaires.
- Les 17 `.poppins-*` — une graisse chacune, remplacées par `font-thin`, `font-normal`, `font-medium`, `font-bold`.
- `header { ... }` — `slide-header` le porte.
- `.main-slide`, `.gb-shape-1`, `.logos`, `img { height: 75% }`, `split-view { ... }`, `#progression` — soit morts (aucun usage dans le deck, vérifier au `grep` avant de supprimer), soit portés par le châssis. **`img { height: 75% }` en particulier : c'est une règle globale sur toutes les images du deck, à retirer avec attention et à vérifier sur les slides 4, 6, 7, 17 et 18.**
- `blockquote h5 { font-size: 2.1rem }`, `blockquote { border-inline-start-color }` — passés en classes sur les slides concernées (Tâche 6).

- [ ] **Step 3: Supprimer le fichier mort**

Run:
```bash
grep -rn "poppins.css" . --exclude-dir=vendor --exclude-dir=.git
git rm impact-framework/.demoit/poppins.css
```
Expected: le `grep` ne renvoie rien avant la suppression — le fichier n'est chargé par personne, et ses classes sont dupliquées dans `style.css`.

- [ ] **Step 4: Vérifier que la police vient bien du disque**

Run:
```bash
hack/css.sh impact-framework && go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
curl -s -o /dev/null -w "font %{http_code} %{size_download}\n" http://127.0.0.1:8899/fonts/poppins-400.woff2
grep -c "fonts.googleapis\|fonts.gstatic" handlers/resources/index.tmpl.html
```
Expected: la police répond 200 avec une taille non nulle ; le `grep` renvoie 0.

- [ ] **Step 5: Recapturer et comparer, puis commit**

```bash
bash <scratchpad>/shootall.sh <scratchpad>/after-t8/png-1920 http://127.0.0.1:8899 17 18
git add impact-framework/.demoit
git commit -m "$(cat <<'MSG'
refactor(impact-framework): thin out style.css and vendor Poppins

What became a design token leaves style.css: the palette, the pane
heights, and the seventeen .poppins-* classes that were one font-weight
The dead 3vw root font size goes too: beercss loaded after this sheet
and its own html rule won, so a rem has always been 16px here. Nothing
in the deck changes size.

Poppins now loads from .demoit/fonts through the /fonts/ route, which
existed and served nothing. A talk has to render with no network, and
the CDN link made that impossible.

Deletes .demoit/poppins.css: no template ever loaded it, and its classes
were already duplicated in style.css.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 9 : La chrome des composants et les runtimes

Objectif : les couleurs codées en dur des composants passent par des custom properties, et `source-code` et mermaid suivent le thème.

**Files:**
- Modify: `impact-framework/.demoit/js/demoit.js`

**Interfaces:**
- Consumes: `--dm-*` (Tâche 1), `window.demoitTheme` et l'événement `demoit:theme` (Tâche 5).

- [ ] **Step 1: Remplacer les couleurs en dur de `FakeWindow`**

Dans `static get styles()` de `FakeWindow` (`demoit.js:26-...`) :
- `background: #ddd` et `background-color: white` de `.main` → `background-color: var(--dm-window-bg, #fff)`
- `background: linear-gradient(to bottom, #edeaed 0%, #dddfdd 100%)` de `#bar` → `background: var(--dm-chrome-bg)`
- `border-bottom: 2px solid #cbcbcb` → `border-bottom: 2px solid var(--dm-chrome-border, #cbcbcb)`

Les custom properties traversent le shadow DOM ; les classes, non. C'est pour cette raison que le `class="large-height"` de `SourceCode.render()` (`demoit.js:266`) et toutes les classes beercss de `TitleBar.render()` (`demoit.js:549-559`) n'ont jamais rien fait.

- [ ] **Step 2: Remplacer les couleurs en dur de `SourceCode`**

Dans `static get styles()` de `SourceCode` (`demoit.js:177-244`) :
- `.chroma { color: #212121 }` → `color: var(--dm-code-fg, #212121)`
- `#tabs { background-color: rgb(243, 243, 243); border-bottom: 1.5px solid rgb(236, 236, 236) }` → `background-color: var(--dm-tabs-bg); border-bottom: 1.5px solid var(--dm-tab-bg)`
- `#tabs a { background: rgb(236, 236, 236); color: black }` → `background: var(--dm-tab-bg); color: var(--dm-tab-fg)`
- `#tabs a.selected { background: white }` → `background: var(--dm-window-bg)`
- `--default-color-selection: rgb(191, 214, 255)` → `var(--dm-code-selection)`

- [ ] **Step 3: Retirer les classes mortes**

`demoit.js:266` : `<div id="container" class="large-height">` → `<div id="container">`.
`demoit.js:546-562`, `TitleBar.render()` : remplacer le corps par du markup qui utilise le châssis, puisque les classes beercss n'y avaient aucun effet :

```js
    render() {
        return `
            <header class="slide-header">
                <h5 class="flex-1 m-0 text-left"><slot></slot></h5>
                <div class="flex items-center gap-4">
                    <img class="size-12 rounded-full object-cover" src="/images/zatsit_logo.svg">
                    <img class="size-12 rounded-full object-cover" src="/images/cloud_nord2.png">
                </div>
            </header>
            `;
    }
```

**Attention** : ces classes sont dans un shadow DOM, donc `slide-header` ne les atteindra pas davantage. Deux options : garder `TitleBar` avec des styles propres dans son `static get styles()` (cohérent avec les autres composants), ou lui faire renoncer au shadow DOM. **Choix : lui donner ses propres styles**, comme tous ses voisins, en reprenant les valeurs de `slide-header` :

```js
    static get styles() {
        return `
        header {
            display: flex;
            align-items: center;
            gap: 1rem;
            padding-inline: 0.3rem;
            color: var(--color-main);
            border-block-end: var(--rule, 1px) solid var(--color-main);
        }

        h5 {
            flex: 1;
            margin: 0;
            text-align: start;
        }

        img {
            block-size: 3rem;
            inline-size: 3rem;
            border-radius: 9999px;
            object-fit: cover;
        }`
    }

    render() {
        return `
            <header>
                <h5><slot></slot></h5>
                <div>
                    <img src="/images/zatsit_logo.svg">
                    <img src="/images/cloud_nord2.png">
                </div>
            </header>`;
    }
```
Noter aussi que `TitleBar.connectedCallback` appelle `render()` une seconde fois et ajoute le résultat au shadow root, en plus de ce que `BaseHTMLElement.connectedCallback` a déjà fait : le titre est donc rendu deux fois. Corriger en supprimant l'override de `connectedCallback`.

- [ ] **Step 4: Faire suivre le thème à `source-code`**

Dans `SourceCode`, remplacer la lecture du style et ajouter une réaction au thème :

```js
    // The chroma style is resolved from the theme rather than from the
    // attribute alone: /sourceCode ships its own stylesheet with the HTML
    // (html.Standalone(true) in handlers/code.go), so switching the theme is a
    // re-fetch and nothing has to be restyled here.
    currentStyle() {
        if (document.documentElement.classList.contains('dark')) {
            return this.getAttribute('code_style_dark') || 'github-dark';
        }

        return this.getAttribute('code_style') || 'vs';
    }
```

Dans `showCurrentTab`, remplacer `style=${this.code_style}` par `style=${this.currentStyle()}`, et mémoriser l'onglet courant :

```js
    async showCurrentTab(current) {
        this.current = current;
        const file = this.files[current];
        const startLines = this.startLines[current];
        const endLines = this.endLines[current];
        const url = `/sourceCode/${this.folder}/${file}?hash=${this.hash}&style=${this.currentStyle()}&startLine=${startLines}&endLine=${endLines}`;
        ...
    }
```

Et dans `connectedCallback`, après `this.showCurrentTab(0)` :

```js
        document.addEventListener('demoit:theme', () => this.showCurrentTab(this.current || 0));
```

- [ ] **Step 5: Faire suivre le thème à mermaid**

Remplacer `demoit.js:630-631` :

```js
import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@11.6.0/+esm';

function renderMermaid() {
    const dark = document.documentElement.classList.contains('dark');
    mermaid.initialize({ startOnLoad: false, theme: dark ? 'dark' : 'default' });

    // A re-render needs the processed marker cleared, and the original source
    // put back: mermaid replaces the element's content with the rendered SVG.
    document.querySelectorAll('.mermaid, pre.mermaid').forEach(node => {
        if (node.dataset.source === undefined) {
            node.dataset.source = node.textContent;
        } else {
            node.textContent = node.dataset.source;
        }
        delete node.dataset.processed;
    });

    mermaid.run();
}

renderMermaid();
document.addEventListener('demoit:theme', renderMermaid);
```

- [ ] **Step 6: Vérifier les deux thèmes sur les slides concernées**

Run:
```bash
hack/css.sh impact-framework && go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shoot1.sh <scratchpad>/after-t9/code-light.png "http://127.0.0.1:8899/9?theme=light" 20
bash <scratchpad>/shoot1.sh <scratchpad>/after-t9/code-dark.png "http://127.0.0.1:8899/9?theme=dark" 20
bash <scratchpad>/shoot1.sh <scratchpad>/after-t9/term-dark.png "http://127.0.0.1:8899/10?theme=dark" 20
```
Expected: la slide 9 porte le mermaid et le `::code`. En clair, code sur fond blanc et onglets gris ; en sombre, code clair sur fond sombre, onglets sombres, mermaid en thème sombre, et la chrome des fenêtres sombre. Le terminal reste sombre dans les deux cas — c'est voulu.

- [ ] **Step 7: Commit**

```bash
git add impact-framework/.demoit/js/demoit.js
git commit -m "$(cat <<'MSG'
feat(impact-framework): let the embedded runtimes follow the theme

Component styles live inside a shadow root, so the page's classes never
reach them -- which is why the class="large-height" on the code viewer's
container and every beercss class in TitleBar have never done anything.
Custom properties do cross that boundary, so the hard-coded window,
tab-strip and code colours become var(--dm-*) and follow the theme.

The code viewer re-fetches /sourceCode with a dark chroma style on a
theme change: the route already ships its own stylesheet with the HTML,
so there is nothing to restyle here. mermaid re-runs with its dark
theme, which needs the original source put back first, since a render
replaces the element's content with the SVG.

TitleBar also rendered its title twice: its connectedCallback appended a
second copy on top of the one BaseHTMLElement had already built.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 10 : Les fenêtres satellites et l'impression

Objectif : `/grid`, `/speakernotes` et `/pdf` se comportent correctement vis-à-vis du thème.

**Files:**
- Modify: `handlers/resources/grid.tmpl.html`, `handlers/resources/speakernotes.tmpl.html`, `handlers/pdf.go`

- [ ] **Step 1: Faire suivre le thème à la fenêtre de notes**

Dans `speakernotes.tmpl.html`, ajouter la feuille du moteur et le script pré-paint dans le `<head>`, à l'identique de `index.tmpl.html` (le même bloc, sans la garde `{{ if .Dark }}` : la page de notes n'a pas de `Page`), puis dans le `channel.onmessage` existant, avant le test sur `currentSlideId` :

```js
      if(e.data.hasOwnProperty("theme")) {
        document.documentElement.classList.toggle('dark', e.data.theme === 'dark');
        document.documentElement.dataset.display = e.data.display || 'screen';
        try {
          localStorage.setItem('demoit_theme', e.data.theme);
          localStorage.setItem('demoit_display', e.data.display || 'screen');
        } catch (err) {
          // Nothing to remember is better than nothing to render.
        }
      }
```

- [ ] **Step 2: Faire suivre le thème à `/grid`**

Dans `grid.tmpl.html`, ajouter le même `<link>` et le même script pré-paint, et remplacer les couleurs en dur du `<style>` (`border: 1px solid #BBB`) par `var(--color-outline)`. Les iframes sont sur la même origine et lisent le même `localStorage`, donc elles suivent sans travail supplémentaire.

- [ ] **Step 3: Forcer le clair à l'impression**

Dans `handlers/pdf.go`, ajouter `theme=light` à l'URL de chaque page rendue. Localiser la construction de l'URL et la faire passer par un paramètre :

```go
	// The theme is forced light for print: a PDF of dark slides is a PDF of
	// ink. The `?grid=true` parameter already sets the precedent for a
	// render-mode query parameter.
```
Le paramètre est lu par le script pré-paint de `index.tmpl.html`, écrit à la Tâche 5, qui donne priorité à l'URL sur `localStorage`.

- [ ] **Step 4: Vérifier les trois**

Run:
```bash
go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shoot1.sh <scratchpad>/after-t10/grid-dark.png "http://127.0.0.1:8899/grid?theme=dark" 25
bash <scratchpad>/shoot1.sh <scratchpad>/after-t10/notes.png "http://127.0.0.1:8899/speakernotes?theme=dark" 15
curl -s -o <scratchpad>/after-t10/deck.pdf --max-time 300 -w "pdf %{http_code} %{size_download}\n" http://127.0.0.1:8899/pdf
```
Expected: la grille et les notes en sombre ; le PDF en clair, quelle que soit la valeur mémorisée.

- [ ] **Step 5: Vérifications et commit**

```bash
go build -o /dev/null . && go test ./... && gofmt -l . | grep -v '^vendor/'
git add handlers/
git commit -m "$(cat <<'MSG'
feat(handlers): carry the theme to the satellite windows

The speaker-notes window and /grid read the theme from the same origin's
localStorage and listen on the channel the notes window already used, so
they no longer stay white while the stage goes dark -- a real problem in
a control room, where both windows are in front of the speaker.

/pdf forces the light theme regardless of what is remembered: a PDF of
dark slides is a PDF of ink.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 11 : Vérifier `sample/`

Objectif : `sample/` rend et navigue, sans scène ni thème.

**Files:**
- Modify (si nécessaire): `sample/.demoit/style.css`

- [ ] **Step 1: Capturer l'avant et l'après**

Run:
```bash
/tmp/demoit-before --port 8898 --shellport 9997 sample &
bash <scratchpad>/shootall.sh <scratchpad>/before-sample/png-1920 http://127.0.0.1:8898 5 15
/tmp/demoit-after --port 8899 --shellport 9998 sample &
bash <scratchpad>/shootall.sh <scratchpad>/after-sample/png-1920 http://127.0.0.1:8899 5 15
```
Note : `sample/demoit.html` a moins de slides que le deck d'impact-framework ; ajuster le dernier index en comptant les liens de `/grid`.

- [ ] **Step 2: Constater le delta attendu, corriger seulement ce qui est cassé**

`sample/` n'utilise que `.mermaid` et `.logos`, toutes deux définies dans son propre `style.css`. Il perd en revanche la feuille de base de beercss : tailles de `h1`…`h6`, `font-size` et `line-height` du `body`, `display:flex` de `header`/`footer`, et la règle de rythme vertical. Il reçoit à la place le preflight de Tailwind, qui remet les titres à la taille du texte courant et enlève les puces des listes.

Si le rendu est illisible, ajouter à `sample/.demoit/style.css` uniquement les règles de base qui manquent — c'est un fichier qui appartient au talk, donc le bon endroit :

```css
/*
beercss used to supply the base sheet this deck was written against: heading
sizes, the body font, and a `* + :is(h1..h6,p,ul,ol,...)` rule that gave every
slide its vertical rhythm. Tailwind's preflight makes none of those
assumptions, so the few this deck relies on are restated here rather than in
the engine -- sample/ deliberately stays off the new chassis.
*/
h1 { font-size: 3.5625rem; }
h2 { font-size: 2.8125rem; }
h3 { font-size: 2.25rem; }
h4 { font-size: 2rem; }
h5 { font-size: 1.75rem; }
h6 { font-size: 1.5rem; }

:not(nav) > :is(ul, ol) { all: revert; }

* + :is(blockquote, h1, h2, h3, h4, h5, h6, ol, p, pre, table, ul) {
    margin-block-start: 1rem;
}
```

- [ ] **Step 3: Vérifier que `sample/` ne reçoit ni scène ni commutateur**

Run:
```bash
curl -s http://127.0.0.1:8899/ | grep -E 'class="(stage|no-stage)"|demoit_theme'
```
Expected: `class="no-stage"` et **aucune** mention de `demoit_theme` — `sample/.demoit/talk.yml` ne déclare pas de bloc `theme:`. Si `talk.yml` n'existe pas pour `sample/`, c'est correct : les deux drapeaux valent faux par défaut (Tâche 2).

- [ ] **Step 4: Commit**

```bash
git add sample/
git commit -m "$(cat <<'MSG'
fix(sample): restate the base rules beercss used to supply

sample/ stays off the new chassis -- no stage, no theme switcher, since
its talk.yml declares no theme block. What the opt-in cannot protect it
from is the shared index template losing beercss altogether: heading
sizes, the body font and the sibling rule that gave every slide its
vertical rhythm all came from there.

Those few rules now live in the talk's own stylesheet, which is where a
talk's assumptions belong.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 12 : La documentation

Objectif : `CLAUDE.md` décrit le repo tel qu'il est.

**Files:**
- Modify: `CLAUDE.md`
- Modify: `impact-framework/demoit.html` (une ligne d'en-tête)

- [ ] **Step 1: Recenser ce qui est devenu faux**

Run: `grep -n "beercss\|s12 m6 l6\|style.css\|xlarge-height\|--color-main\|poppins" CLAUDE.md`

Les cinq passages à reprendre : la section « Commands » (ajouter `hack/css.sh`), le tableau des directives (`split{cols=}` et la grille 12 colonnes), la section `:::grid`/`:::col` (`impact-framework/demoit.md:349-380`, dont les numéros de ligne ont bougé), la section « Everything the browser loads comes from `<folder>/.demoit/` » (ajouter `tailwind.css`, `tokens.css`, `components.css`, `fonts/`, retirer `poppins.css`), et la dernière ligne des « Conventions » sur la palette.

- [ ] **Step 2: Ajouter les commandes**

Dans le bloc « Commands » :

```bash
hack/css.sh engine                      # rebuild handlers/resources/demoit.css (committed)
hack/css.sh impact-framework            # rebuild <talk>/.demoit/tailwind.css (committed)
hack/css.sh impact-framework --watch    # use alongside `demoit --dev impact-framework`
```

Et une phrase : les deux CSS sont des artefacts commités ; toute classe Tailwind inédite écrite dans une slide exige de relancer la commande du talk, sans quoi la classe n'existe pas dans la feuille.

- [ ] **Step 3: Corriger les sections sur les classes**

Remplacer les mentions de la grille beercss par la grille Tailwind à 12 colonnes, `s4` par `col-span-4`, `{height}-height` par `h-stage-{small,medium,large,xlarge}`, et ajouter le piège vérifié : les classes construites en Go n'existent que par `@source inline(...)` dans `styles/demoit.src.css`.

Documenter aussi le bloc `theme:` de `talk.yml` et ses deux drapeaux, la scène et ses conséquences (1920×1080, `rem` à 16 px inchangé, les `vw` qui mesurent désormais la scène, letterboxing hors 16:9), les profils d'affichage et les touches `t` et `d`.

- [ ] **Step 4: Marquer le deck de référence**

En tête de `impact-framework/demoit.html`, compléter le commentaire existant : le fichier est une référence **pré-Markdown et pré-Tailwind**, jamais servie, conservée à dessein.

- [ ] **Step 5: Commit**

```bash
git add CLAUDE.md impact-framework/demoit.html
git commit -m "$(cat <<'MSG'
docs: describe the Tailwind chassis

CLAUDE.md described beercss in five sections. Updates the grid and
height vocabulary, adds the hack/css.sh commands and the fact that both
stylesheets are committed artefacts, documents the talk.yml theme block,
the fixed stage and the display profiles.

Records the one trap the migration exposed: the classes deck/directive
assembles in Go reach the stylesheet only through `@source inline(...)`,
because no scanner will ever meet them in a source file.

Co-Authored-By: Claude Opus 5 (1M context) <noreply@anthropic.com>
MSG
)"
```

---

## Tâche 13 : La passe finale

Objectif : tout vérifier ensemble, thèmes × profils × routes.

- [ ] **Step 1: La suite automatique**

Run:
```bash
go build -o /dev/null .
go test ./...
gofmt -l . | grep -v '^vendor/'
golangci-lint run -c golangci.yml
```
Expected: build OK, tous les paquets `ok`, `gofmt` muet, lint muet.

- [ ] **Step 2: Aucune trace de beercss**

Run:
```bash
grep -rn "beercss\|beer.min\|material-dynamic" . --exclude-dir=vendor --exclude-dir=.git --exclude-dir=.tools
```
Expected: seules des occurrences dans `docs/superpowers/` (la spec et ce plan, qui racontent la migration) et dans les blocs commentés conservés comme référence.

- [ ] **Step 3: Les deux thèmes, les trois profils**

Run:
```bash
go build -o /tmp/demoit-after . && /tmp/demoit-after --port 8899 --shellport 9998 impact-framework &
bash <scratchpad>/shoot1.sh <scratchpad>/final/light-screen.png "http://127.0.0.1:8899/3?theme=light&display=screen" 20
bash <scratchpad>/shoot1.sh <scratchpad>/final/light-projector.png "http://127.0.0.1:8899/3?theme=light&display=projector" 20
bash <scratchpad>/shoot1.sh <scratchpad>/final/dark-screen.png "http://127.0.0.1:8899/3?theme=dark&display=screen" 20
bash <scratchpad>/shoot1.sh <scratchpad>/final/dark-tv.png "http://127.0.0.1:8899/3?theme=dark&display=tv" 20
```
Expected: six rendus cohérents. Le profil ne doit **jamais** déplacer un élément : seuls contraste, graisse et épaisseur de trait changent. Un décalage de position signale qu'un jeton de géométrie s'est glissé dans un profil.

- [ ] **Step 4: Le letterboxing**

Run:
```bash
bash <scratchpad>/shoot1.sh <scratchpad>/final/ratio-4-3.png "http://127.0.0.1:8899/3" 20
```
en modifiant `--window-size` de `shoot1.sh` à `1024,768`, puis à `3440,1440`.
Expected: la slide garde ses proportions et se centre, avec des bandes de `--color-void`. Aucun élément ne doit être coupé ni redistribué.

- [ ] **Step 5: Le deck complet, une dernière fois**

Run:
```bash
bash <scratchpad>/shootall.sh <scratchpad>/final/png-1920 http://127.0.0.1:8899 17 20
curl -s -o <scratchpad>/final/deck.pdf --max-time 300 -w "pdf %{http_code} %{size_download}\n" http://127.0.0.1:8899/pdf
```
Comparer les 18 slides à `before/png-1920/` et le PDF à `before/deck-before.pdf`. Critère : même intention, même hiérarchie, rien de perdu.

- [ ] **Step 6: Les routes de démo**

Une slide de chaque famille, à l'œil dans un navigateur réel : `::term` (slide 10), `::vscode` (slide 12, exige Docker), `::browser` (slide 6), `::code` (slide 9), `split{cols=}` (slide 13), `:::grid`/`:::col` (slide 12), mermaid (slide 9), `layout: cover` (slide 0), `layout: split` (slides 8-10). Plus `/last`, et la navigation clavier dans les deux sens.

- [ ] **Step 7: Commit final si des retouches ont été nécessaires**

---

## Self-Review

**Couverture de la spec :** chaque section a sa tâche. Architecture CSS → T1. Chaîne de build → T1. Classes invisibles au scanner → T1 Step 5, gardées par un test. Suppressions de `index.tmpl.html` → T5. Poppins et `poppins.css` → T8. La scène → T1 (CSS) et T5 (template et JS). Les profils → T1, vérifiés en T13. Le thème sombre → T1, T5, T9, T10. La table de correspondance → T4 (layouts), T6 (deck), T7 (deck anglais). Changements Go → T2 et T3. Le sort de `sample/` → T11. Vérification → T0 pour la référence, puis à chaque tâche, et T13 pour la passe finale. Documentation → T12.

**Deux écarts assumés par rapport à la spec, et leur raison :**

1. **Les hauteurs sont des parts, pas des longueurs.** La spec annonçait un recalibrage à l'œil de valeurs absolues. Des parts de la zone principale (`40%`, `65%`, `85%`, `100%`) suppriment le problème par construction plutôt que de le régler : un panneau ne peut plus être plus haut que la scène. Mettre la spec à jour.
2. **Les noms de tailles restent ceux du deck** (`h-stage-xlarge`, pas `h-stage-xl`), parce que `split.html` interpole `.Height` verbatim et que le frontmatter écrit `xlarge`. Traduire en Go aurait ajouté une table de correspondance pour rien.

**Cohérence des noms :** `stageHeightClass(height)` est défini en T3 (`render.go`) et utilisé par `toGrid` dans le même commit. `handlers.EngineCSS` est défini en T1 et routé en T1. `deck.Theme{Stage, Dark}` est défini en T2 et lu par `index.tmpl.html` en T5 sous les noms `.Stage` et `.Dark`. `window.demoitTheme` et l'événement `demoit:theme` sont définis en T5 et consommés en T9 et T10. Les jetons `--dm-*` sont déclarés en T1 et consommés en T9. Les noms de hauteur `small`/`medium`/`large`/`xlarge` sont fixés en T4 Step 2 et propagés aux jetons, à l'`@source inline`, aux tests de T3 et au test de T1 (T4 Step 5).
