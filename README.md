# DemoIt

**DemoIt** is a tool that helps you create beautiful live-coding demonstrations.

## Why?

I'm doing lots of live-coding during conferences and I like tools like
[reveal.js](https://revealjs.com/) to create slides. What I wanted was a
tool that has some of the properties of reveal.js but with
capabilities to code and run commands in front of the audience.

Two things are really important to me:
 + For a given presentation, the slides should live in the same repository as the code.
 + The tool should allow context-switching-less live coding demos.
Attendees don't want to watch me switching between a browser, an IDE and a
terminal all the time.

This is how I came up with **DemoIt**.

## How?

**DemoIt** is a small command line tool written in Go. It serves
rich web content composed of text, images and smart web components.

Those web components make most of the *magic*.

 + One component displays multi-tab ttys that can be used to run any command.
 + Another is a web browser view that auto refreshes itself.
 + Another is a code viewer with highlighing and tabs that looks like a real IDE,

## Install

### Download binary from GitHub

```bash
curl -L -odemoit https://github.com/dgageot/demoit/releases/download/v1.0/demoit-`uname -s | tr '[:upper:]' '[:lower:]'`-`uname -m`
sudo install demoit /usr/local/bin/demoit
```

### Add shell font

To have a correct display in the web terminal, it's better to install the font `Inconsolata for Powerline` on your computer.
This font can be found [here](https://github.com/powerline/fonts/tree/master/Inconsolata).

## Get started from a Sample

```bash
# Create an empty directory
cd $HOME; mkdir my-demoit-presentation; cd $HOME/my-demoit-presentation
# Download a sample demo
curl -L https://github.com/dgageot/demoit/archive/master.tar.gz | tar xvf - --strip-components=2 demoit-master/sample
# Run demoit
demoit
```

Then, browse to http://localhost:8888

*Pro tip:* Run `demoit -dev` instead and enjoy live reload each time you change the content of the slides.

### How do I customize my presenttion then?

Basically, the idea is to:

 + Write content in `demoit.md` at the root of the project — see [Writing slides in Markdown](#writing-slides-in-markdown) below for the format. (You can still write `demoit.html` directly instead, with slides separated by `---`, if you'd rather author raw HTML.)
 + Add images, fonts and scripts in the `.demoit` folder at the root of the project.
 + Customize the style sheet in `.demoit/style.css`.

## Writing slides in Markdown

A presentation's slides live in `demoit.md` (or `demoit-<locale>.md` for a localized deck, tried before the non-localized name). If a talk still has a `demoit.html`, it keeps working exactly as before — Markdown is only used when the `.md` file exists.

A deck is a sequence of slides separated by a line of three or more dashes:

```markdown
---
layout: cover
---

# My talk

***

---
title: A second slide
source: https://example.com
---

Regular *Markdown* content: lists, `code`, images, tables, and so on.
```

Three things worth knowing about that separator:
 + Since `---` cuts a slide, it can't double as Markdown's horizontal-rule syntax anymore — write `***` for a rule instead.
 + For the same reason it can't underline a setext heading either. `Titre` on one line and `---` on the next used to be an `<h2>`; it is now a one-line slide followed by a new one. Write `## Titre`, or put the text in the slide's `title:` frontmatter key.
 + The splitter skips a `---` written inside a fenced code block, but **not** one written inside an HTML comment — a commented-out block containing a literal `---` still starts a new slide.

### Frontmatter

The block right after a separator (or at the very start of the file) is read as frontmatter when its first line looks like a lowercase `key: value` pair, and it runs until the next `---`. A slide needs no frontmatter at all if it has none of the keys below.

| Key | Meaning |
|---|---|
| `layout` | Which layout renders the slide (see below); defaults to the talk's own default layout. |
| `title` | The slide's title — rendered as Markdown, so `**bold**` and other inline formatting work. |
| `source` | Shown as a "Source: …" line under the slide. |
| `class` | CSS classes for the slide's main container. This **replaces** the layout's default classes rather than adding to them — give it every class you want, not just the extra ones. |
| `height` | Height of a `split` layout's columns (defaults to `xlarge`). |
| `speakernotes` | Notes shown only in the speaker window. Also Markdown, so a `-` bulleted list renders as a real list there. |

A body paragraph that happens to start the same way a frontmatter key does (a lowercase word, a colon, then more text) followed later by a `---` gets swallowed as frontmatter too — start such a paragraph with a capital letter if you meant it as regular text. Any frontmatter key demoit doesn't recognize is passed through untouched, for a custom layout to use however it likes.

### Layouts

Six layouts ship with demoit: `cover`, `default`, `content`, `quote`, `split`, `bare` — pick one with the `layout:` key. To change how one looks, or to add a layout of your own, drop a `.html` template at `.demoit/layouts/<name>.html` in your presentation; a same-named file there always wins over the built-in one, and a name demoit has never heard of works too as long as some slide's `layout:` points at it. A layout owns its default classes as well as its markup: write `<main class="{{ if .Class }}{{ .Class }}{{ else }}your defaults{{ end }}">` so a slide that names no `class:` key still gets something. (`partials` is not a layout — it only defines the shared header/source/notes fragments, and naming it in a `layout:` key is an error.)

### Directives

Directives put a live web component (a terminal, a browser, the code viewer, a split screen...) on a slide without writing raw HTML. There are two shapes:

 + `::name{key=value}` on its own line — a single, self-closing element.
 + `:::name{key=value}` ... `:::` — a block that can contain other Markdown, or other directives.

Attribute values are `key=value`, separated by spaces; wrap a value that contains a space in double quotes (`title="My Window"`). The braces are not optional and the directive takes the whole line: `::term path=sources` (braces forgotten) or `::term{path=a} and some text` (text after the block) is reported as an error on that slide rather than quietly rendering a component with nothing in it.

```markdown
:::split{cols=4,8 height=xlarge}
::term{path=sources}

::code{folder=sources files=main.go,go.mod lines=1-20,1-5}
:::
```

 + `::term{path=}`, `::browser{src=}` and `::vscode{path=}` open a terminal, an embedded browser, and a VS Code instance.
 + `::code{folder= files= lines= style=}` renders the tabbed, IDE-like code viewer. `files` and `lines` are comma-separated — `lines=11-20,4-9` gives one highlight range per file, in the same order as `files`. Those ranges **highlight** the given lines; the file is always shown in full, never truncated. `style` picks a chroma theme (`vs` if you omit it, which is also what the viewer defaults to on its own).
 + `:::split{...}` puts its content side by side. `height=` alone gives an auto-arranged split, the same one the `split` layout produces; add `cols=4,8` (one weight per pane, out of 12) for explicit column widths — each weight has to be a whole number from 1 to 12, and a `lines=` range has to be two numbers, or the slide reports the mistake instead of rendering it. Leave a **blank line** between panes when using `cols=` — the columns are counted one per paragraph/block, so two panes glued together without a blank line count as a single, wider pane and the widths won't line up.
 + `:::window{title=}` frames its content like a little macOS window; `:::speakernotes` marks content as speaker notes (the `speakernotes:` frontmatter key is the simpler way to do the same thing).

### `.demoit/talk.yml`

```yaml
layout: default
logos:
  - /images/logo.svg
  - /images/logo2.jpg
```

`layout` is the deck-wide default every slide falls back to when its own frontmatter names none. `logos` are shown in the header of `default`, `content`, `quote` and `split` (`bare` has no header, so it shows none), and separately again in `cover`'s own logo grid. `title` and `footer` also parse without error, but nothing in the shipped layouts currently reads either one — set them and nothing will change on screen.

## Contribute

### Build from sources

```bash

git clone https://github.com/dgageot/demoit.git
cd demoit
go install
```

```
export GOPATH=$HOME/Go
export GOROOT=/opt/homebrew/opt/go/libexec/
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin

```

*This requires Go 1.19 or later.*

