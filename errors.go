// Copyright (c) 2024 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"regexp"
	"strconv"
)

// Error represents a generic typst error.
type Error struct {
	Inner error

	Raw     string // The raw output from stderr.
	Message string // The parsed error message.
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return e.Raw
}

func (e *Error) Unwrap() error {
	return e.Inner
}

// ErrorWithPath represents a typst error that also contains information about its origin (filepath, line and column).
type ErrorWithPath struct {
	Inner error

	Raw     string // The raw error string as returned by the executable.
	Message string // Error message from typst.

	Path   string // Path of the typst file where the error is located in.
	Line   int    // Line number of the error.
	Column int    // Column of the error.
}

func (e *ErrorWithPath) Error() string {
	return e.Raw
}

func (e *ErrorWithPath) Unwrap() error {
	return e.Inner
}

var stderrRegex = regexp.MustCompile(`^error: (?<error>.+)\n`)
var stderrWithPathRegex = regexp.MustCompile(`^(?<path>.+):(?<line>\d+):(?<column>\d+): error: (?<error>.+)\n$`)

// ParseStderr will parse the given stderr output and return a suitable error object.
// Depending on the stderr message, this will return either a typst.Error or a typst.ErrorWithPath error.
func ParseStderr(stderr string, inner error) error {
	if parsed := stderrWithPathRegex.FindStringSubmatch(stderr); parsed != nil {
		err := ErrorWithPath{
			Raw:   stderr,
			Inner: inner,
		}

		if i := stderrWithPathRegex.SubexpIndex("error"); i > 0 && i < len(parsed) {
			err.Message = parsed[i]
		}
		if i := stderrWithPathRegex.SubexpIndex("path"); i > 0 && i < len(parsed) {
			err.Path = parsed[i]
		}
		if i := stderrWithPathRegex.SubexpIndex("line"); i > 0 && i < len(parsed) {
			line, _ := strconv.ParseInt(parsed[i], 10, 0)
			err.Line = int(line)
		}
		if i := stderrWithPathRegex.SubexpIndex("column"); i > 0 && i < len(parsed) {
			column, _ := strconv.ParseInt(parsed[i], 10, 0)
			err.Column = int(column)
		}

		return &err
	}

	if parsed := stderrRegex.FindStringSubmatch(stderr); parsed != nil {
		err := Error{
			Raw:   stderr,
			Inner: inner,
		}

		if i := stderrRegex.SubexpIndex("error"); i > 0 && i < len(parsed) {
			err.Message = parsed[i]
		}

		return &err
	}

	// Fall back to the raw error message.
	return &Error{
		Raw:   stderr,
		Inner: inner,
	}
}
