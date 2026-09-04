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

package directive

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// stackKey holds the containers open at the current point of the document.
var stackKey = parser.NewContextKey()

// openContainer is one entry of the stack of open containers.
type openContainer struct {
	node *Node
	// contentIndent is the indentation of the container's first non-blank
	// content line, stripped from every following line so that nesting can be
	// indented without turning into a Markdown code block.
	contentIndent int
	// contentStarted records that contentIndent has been measured.
	contentStarted bool
}

// blockParser parses `:::name{attrs}` containers and `::name{attrs}` leaves.
type blockParser struct{}

// NewParser returns the block parser of the directive extension.
func NewParser() parser.BlockParser {
	return &blockParser{}
}

// Trigger implements parser.BlockParser.
func (b *blockParser) Trigger() []byte {
	return []byte{':'}
}

// Open implements parser.BlockParser.
func (b *blockParser) Open(_ ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || line[pos] != ':' {
		return nil, parser.NoChildren
	}

	end := pos
	for ; end < len(line) && line[end] == ':'; end++ {
	}
	fenceLength := end - pos
	if fenceLength < 2 {
		return nil, parser.NoChildren
	}

	nameEnd := end
	for ; nameEnd < len(line) && isNameByte(line[nameEnd]); nameEnd++ {
	}
	name := string(line[end:nameEnd])
	if name == "" {
		return nil, parser.NoChildren
	}

	lineNumber, _ := reader.Position()
	reader.Advance(nameEnd)

	attrs, malformed := readAttributes(reader)

	node := NewNode(name, attrs, fenceLength, lineNumber+1)
	if malformed {
		addError(pc, node.Line, "the %s directive has an attribute block that does not parse: check the braces", name)
	}

	if !node.Container {
		return node, parser.NoChildren
	}

	push(pc, node)

	return node, parser.HasChildren
}

// Continue implements parser.BlockParser.
func (b *blockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	directiveNode, ok := node.(*Node)
	if !ok {
		return parser.Close
	}

	stack := containers(pc)
	depth := indexOf(stack, directiveNode)
	if depth < 0 {
		return parser.Close
	}

	line, segment := reader.PeekLine()
	width, pos := util.IndentWidth(line, reader.LineOffset())

	if !stack[depth].contentStarted && !util.IsBlank(line[pos:]) {
		stack[depth].contentStarted = true
		stack[depth].contentIndent = width
		pc.Set(stackKey, stack)
	}

	// Only the deepest open container consumes a closing fence: goldmark calls
	// Continue from the outermost container inwards, so an outer container must
	// leave the line to its children.
	if depth == len(stack)-1 && isClosingFence(line, pos, directiveNode.FenceLength) {
		directiveNode.closed = true
		pc.Set(stackKey, stack[:depth])
		reader.Advance(segment.Stop - segment.Start - trailingNewline(line) + segment.Padding)

		return parser.Close
	}

	// contentIndent is a column width, so a tab counts as four, but Advance
	// moves bytes. Clamp it to pos, the byte index of the line's first
	// non-space character: stripping stops at the content, never inside it.
	// Without the clamp an outdented continuation line loses characters and a
	// tab-indented directive is eaten into literal text, both silently.
	if indent := stack[depth].contentIndent; indent > 0 {
		if indent > pos {
			indent = pos
		}
		if limit := segment.Stop - segment.Start - 1; indent > limit {
			indent = limit
		}
		if indent > 0 {
			reader.Advance(indent)
		}
	}

	return parser.Continue | parser.HasChildren
}

// Close implements parser.BlockParser.
func (b *blockParser) Close(node ast.Node, _ text.Reader, pc parser.Context) {
	directiveNode, ok := node.(*Node)
	if !ok {
		return
	}

	if stack := containers(pc); indexOf(stack, directiveNode) >= 0 {
		pc.Set(stackKey, stack[:indexOf(stack, directiveNode)])
	}

	if directiveNode.Container && !directiveNode.closed {
		addError(pc, directiveNode.Line, "the %s directive is never closed: add a line of %d colons", directiveNode.Name, directiveNode.FenceLength)
	}
}

// CanInterruptParagraph implements parser.BlockParser.
func (b *blockParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine implements parser.BlockParser.
func (b *blockParser) CanAcceptIndentedLine() bool {
	return false
}

// readAttributes reads the `{key=value}` block a directive may carry, and
// reports whether the block was malformed — an unclosed brace, a key with no
// value, an unterminated quote — as opposed to simply absent. A malformed
// block must be reported: silently dropping it renders a component stripped of
// its attributes, which is a broken demo with no diagnostic.
//
// goldmark's own parser.ParseAttributes is deliberately not used. It accepts
// only [A-Za-z0-9_:.-] in an unquoted value and sends anything starting with a
// digit to its number parser, so a URL stops at the first slash, `11-20` stops
// after `11`, and `a.yml,b.yml` stops at the comma — and a value it cannot
// finish makes the whole block fail, stripping every attribute. The grammar
// here is the one the directives need and nothing more: an unquoted value runs
// to the next space or closing brace, and a value that must hold spaces or a
// brace is double-quoted.
func readAttributes(reader text.Reader) (map[string]string, bool) {
	attrs := map[string]string{}

	line, _ := reader.PeekLine()

	start := leadingSpaces(line)
	if start >= len(line) || line[start] != '{' {
		return attrs, false
	}

	end, ok := scanAttributes(line[start+1:], attrs)
	if !ok {
		// Discard whatever pairs were scanned before the failure. Keeping them
		// renders a component with partial, accidentally-plausible attributes
		// next to the error, which is worse than rendering none.
		return map[string]string{}, true
	}

	reader.Advance(start + 1 + end)

	return attrs, false
}

// scanAttributes reads `key=value` pairs up to the closing brace and returns
// the index just past that brace.
func scanAttributes(body []byte, attrs map[string]string) (int, bool) {
	for i := 0; i < len(body); {
		i += leadingSpaces(body[i:])
		if i >= len(body) {
			break
		}
		if body[i] == '}' {
			return i + 1, true
		}

		key, next, ok := scanAttributeKey(body, i)
		if !ok {
			return 0, false
		}

		value, next, ok := scanAttributeValue(body, next)
		if !ok {
			return 0, false
		}

		attrs[key] = value
		i = next
	}

	return 0, false
}

// scanAttributeKey reads a key starting at i and returns it with the index
// just past the `=` that must follow it.
func scanAttributeKey(body []byte, i int) (string, int, bool) {
	end := i
	for end < len(body) && isNameByte(body[end]) {
		end++
	}
	if end == i || end >= len(body) || body[end] != '=' {
		return "", 0, false
	}

	return string(body[i:end]), end + 1, true
}

// scanAttributeValue reads one value starting at i and returns it with the
// index just past it. A double-quoted value may hold spaces and braces; a bare
// one runs to the next space or closing brace.
func scanAttributeValue(body []byte, i int) (string, int, bool) {
	if i < len(body) && body[i] == '"' {
		end := bytes.IndexByte(body[i+1:], '"')
		if end < 0 {
			return "", 0, false
		}

		return string(body[i+1 : i+1+end]), i + end + 2, true
	}

	end := i
	for end < len(body) && isValueByte(body[end]) {
		end++
	}
	if end == i {
		return "", 0, false
	}

	return string(body[i:end]), end, true
}

// isValueByte reports whether c may appear in an unquoted attribute value. The
// newline is excluded so that a block whose brace never closes cannot swallow
// the end of the line into a value.
func isValueByte(c byte) bool {
	return c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '}'
}

// leadingSpaces returns the number of leading spaces and tabs of b.
func leadingSpaces(b []byte) int {
	i := 0
	for i < len(b) && (b[i] == ' ' || b[i] == '\t') {
		i++
	}

	return i
}

// isClosingFence reports whether the line is a run of at least length colons
// followed by nothing but spaces.
func isClosingFence(line []byte, pos, length int) bool {
	end := pos
	for ; end < len(line) && line[end] == ':'; end++ {
	}

	return end-pos >= length && util.IsBlank(line[end:])
}

// trailingNewline returns 1 when the line ends with a newline, 0 otherwise.
func trailingNewline(line []byte) int {
	if len(line) > 0 && line[len(line)-1] == '\n' {
		return 1
	}

	return 0
}

// isNameByte reports whether c may appear in a directive name.
func isNameByte(c byte) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9' ||
		c == '-' || c == '_'
}

// containers returns the stack of open containers held in pc.
func containers(pc parser.Context) []*openContainer {
	stack, _ := pc.Get(stackKey).([]*openContainer)

	return stack
}

// push adds a container to the stack held in pc.
func push(pc parser.Context, node *Node) {
	pc.Set(stackKey, append(containers(pc), &openContainer{node: node}))
}

// indexOf returns the depth of node in the stack, or -1 when it is absent.
func indexOf(stack []*openContainer, node *Node) int {
	for i, entry := range stack {
		if entry.node == node {
			return i
		}
	}

	return -1
}
