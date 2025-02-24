// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"regexp"
	"strconv"
	"strings"
)

// ErrorDetails contains the details of a typst.Error.
type ErrorDetails struct {
	Message string // The parsed error message.
	Path    string // Path of the typst file where the error is located in. Zero value means that there is no further information.
	Line    int    // Line number of the error. Zero value means that there is no further information.
	Column  int    // Column of the error. Zero value means that there is no further information.
}

// Error represents a typst error.
// This can contain multiple sub-errors or sub-warnings.
type Error struct {
	Inner error

	Raw string // The raw output from stderr.

	// Raw output parsed into errors and warnings.
	Details []ErrorDetails
}

func (e *Error) Error() string {
	return e.Raw
}

func (e *Error) Unwrap() error {
	return e.Inner
}

var stderrRegex = regexp.MustCompile(`(?s)^(?<error>.+?)(?:(?:\n\s+┌─ (?<path>.+?):(?<line>\d+):(?<column>\d+)\n)|(?:$))`)

// ParseStderr will parse the given stderr output and return a typst.Error.
func ParseStderr(stderr string, inner error) error {
	err := Error{
		Inner: inner,
		Raw:   stderr,
	}

	// Get all "blocks" ending with double new lines.
	parts := strings.Split(stderr, "\n\n")
	parts = parts[:len(parts)-1]

	for _, part := range parts {
		if parsed := stderrRegex.FindStringSubmatch(part); parsed != nil {
			var details ErrorDetails

			if i := stderrRegex.SubexpIndex("error"); i > 0 && i < len(parsed) && parsed[i] != "" {
				details.Message = parsed[i]
			}
			if i := stderrRegex.SubexpIndex("path"); i > 0 && i < len(parsed) && parsed[i] != "" {
				details.Path = parsed[i]
			}
			if i := stderrRegex.SubexpIndex("line"); i > 0 && i < len(parsed) && parsed[i] != "" {
				if line, err := strconv.ParseInt(parsed[i], 10, 0); err == nil {
					details.Line = int(line)
				}
			}
			if i := stderrRegex.SubexpIndex("column"); i > 0 && i < len(parsed) && parsed[i] != "" {
				if column, err := strconv.ParseInt(parsed[i], 10, 0); err == nil {
					details.Column = int(column)
				}
			}

			err.Details = append(err.Details, details)
		}
	}

	return &err
}
