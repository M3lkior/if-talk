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
	"fmt"

	"github.com/yuin/goldmark/parser"
)

// errorsKey holds the errors collected while parsing one document.
var errorsKey = parser.NewContextKey()

// Error is a problem found in a directive, with the line it was written on.
type Error struct {
	// Line is the 1-based line in the deck file.
	Line int
	// Message says what is wrong.
	Message string
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Message)
}

// Errors returns the directive errors collected while parsing with ctx.
//
// One context belongs to one document: errors accumulate in it, so reusing a
// context across two Convert calls reports the first document's problems again
// on the second.
func Errors(ctx parser.Context) []Error {
	collected, ok := ctx.Get(errorsKey).([]Error)
	if !ok {
		return nil
	}

	return collected
}

// addError records a directive error into ctx.
func addError(ctx parser.Context, line int, format string, args ...any) {
	collected, _ := ctx.Get(errorsKey).([]Error)
	ctx.Set(errorsKey, append(collected, Error{
		Line:    line,
		Message: fmt.Sprintf(format, args...),
	}))
}
