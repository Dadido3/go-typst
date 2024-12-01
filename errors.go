package typst

import (
	"log"
	"regexp"
	"strconv"
)

var stderrRegex = regexp.MustCompile(`^error: (?<error>.+)\n  ┌─ (?<path>.+):(?<line>\d+):(?<column>\d+)\n`)

// Error represents a typst error.
type Error struct {
	Inner error

	Raw string // The raw error string as returned by the executable.

	Message string // Error message from typst.
	Path    string // Path of the typst file where the error is located in.
	Line    int    // Line number of the error.
	Column  int    // Column of the error.
}

// NewError returns a new error based on the stderr from the typst process.
func NewError(stderr string, inner error) *Error {
	err := Error{
		Raw:   stderr,
		Inner: inner,
	}

	parsed := stderrRegex.FindStringSubmatch(stderr)

	if i := stderrRegex.SubexpIndex("error"); i > 0 && i < len(parsed) {
		err.Message = parsed[i]
	}
	if i := stderrRegex.SubexpIndex("path"); i > 0 && i < len(parsed) {
		err.Path = parsed[i]
	}
	if i := stderrRegex.SubexpIndex("line"); i > 0 && i < len(parsed) {
		line, _ := strconv.ParseInt(parsed[i], 10, 0)
		err.Line = int(line)
	}
	if i := stderrRegex.SubexpIndex("column"); i > 0 && i < len(parsed) {
		column, _ := strconv.ParseInt(parsed[i], 10, 0)
		err.Column = int(column)
	}

	log.Printf("%#v", err)

	return &err
}

func (e *Error) Error() string {
	return e.Raw
}

func (e *Error) Unwrap() error {
	return e.Inner
}
