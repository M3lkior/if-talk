#!/usr/bin/env bash
# Builds the Tailwind stylesheets of demoit.
#
#   hack/css.sh engine                     -> handlers/resources/demoit.css
#   hack/css.sh <talk-folder>              -> <talk-folder>/.demoit/tailwind.css
#   hack/css.sh <talk-folder> --watch      -> rebuilds on change
#
# The CLI is a pinned standalone binary: no node_modules, and the network is
# only ever touched here, never while a talk is running.
#
# Both outputs are committed artefacts. Nothing in `go build` regenerates them,
# so a Tailwind class written in a slide does not exist until this runs.
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
    input="${target%/}/.demoit/tailwind.src.css"
    output="${target%/}/.demoit/tailwind.css"
    if [ ! -f "$input" ]; then
        echo "$input not found: is $target a talk folder?" >&2
        exit 1
    fi
fi

echo "building $output" >&2
# --optimize rather than --minify: the output is committed, so it stays
# readable in a diff while still being deduplicated.
exec "$CLI" --input "$input" --output "$output" --optimize $watch
