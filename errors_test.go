// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Dadido3/go-typst"
	"github.com/google/go-cmp/cmp"
)

func TestErrors0(t *testing.T) {
	cli := typst.CLI{}

	r := bytes.NewBufferString(`This is a test!`)

	var w bytes.Buffer
	if err := cli.Compile(r, &w, nil); err != nil {
		t.Fatalf("Failed to compile document: %v", err)
	}
}

func TestErrors1(t *testing.T) {
	cli := typst.CLI{}

	r := bytes.NewBufferString(`This is a test!

#assert(1 < 1, message: "Test")`)

	var w bytes.Buffer
	if err := cli.Compile(r, &w, nil); err == nil {
		t.Fatalf("Expected error, but got nil")
	} else {
		var errTypst *typst.Error
		if errors.As(err, &errTypst) {
			if len(errTypst.Details) != 1 {
				t.Fatalf("Expected error doesn't contain the expected number of detail entries. Got %v, want %v", len(errTypst.Details), 1)
			}
			details := errTypst.Details[0]
			if details.Message != "error: assertion failed: Test" {
				t.Errorf("Expected error with error message %q, got %q", "error: assertion failed: Test", details.Message)
			}
			/*if details.Path != "" {
				t.Errorf("Expected error to point to path %q, got path %q", "", details.Path)
			}*/
			if details.Line != 3 {
				t.Errorf("Expected error to point at line %d, got line %d", 3, details.Line)
			}
			if details.Column != 1 {
				t.Errorf("Expected error to point at column %d, got column %d", 1, details.Column)
			}
		} else {
			t.Errorf("Expected error type %T, got %T: %v", errTypst, err, err)
		}
	}
}

func TestErrors2(t *testing.T) {
	cli := typst.CLI{}

	opts := typst.CLIOptions{
		Pages: "a",
	}

	r := bytes.NewBufferString(`This is a test!`)

	var w bytes.Buffer
	if err := cli.Compile(r, &w, &opts); err == nil {
		t.Fatalf("Expected error, but got nil")
	} else {
		var errTypst *typst.Error
		if errors.As(err, &errTypst) {
			if len(errTypst.Details) != 1 {
				t.Fatalf("Expected error doesn't contain the expected number of detail entries. Got %v, want %v", len(errTypst.Details), 1)
			}
			details := errTypst.Details[0]
			// Don't check the specific error message, as that may change over time.
			// The expected message should be similar to: error: invalid value 'a' for '--pages <PAGES>': not a valid page number.
			if details.Message == "" {
				t.Errorf("Expected error message, got %q", details.Message)
			}
		} else {
			t.Errorf("Expected error type %T, got %T: %v", errTypst, err, err)
		}
	}
}

func TestErrorParsing(t *testing.T) {
	var tests = map[string]struct {
		StdErr          string               // The original and raw stderr message.
		ExpectedDetails []typst.ErrorDetails // Expected parsed result.
	}{
		"Typst 0.13.0 HTML warning + error": {
			StdErr: "warning: html export is under active development and incomplete\n = hint: its behaviour may change at any time\n = hint: do not rely on this feature for production use cases\n = hint: see https://github.com/typst/typst/issues/5512 for more information\n\nerror: page configuration is not allowed inside of containers\n  ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>:1:1\n  │\n1 │ #set page(width: 100mm, height: auto, margin: 5mm)\n  │  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n\n",
			ExpectedDetails: []typst.ErrorDetails{
				{
					Message: "warning: html export is under active development and incomplete\n = hint: its behaviour may change at any time\n = hint: do not rely on this feature for production use cases\n = hint: see https://github.com/typst/typst/issues/5512 for more information",
				},
				{
					Message: "error: page configuration is not allowed inside of containers",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>",
					Line:    1,
					Column:  1,
				},
			},
		},
		"Typst 0.13.0 error with path": {
			StdErr: "error: expected expression\n   ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>:12:34\n   │\n12 │ - Test coverage of most features.#\n   │                                   ^\n\n",
			ExpectedDetails: []typst.ErrorDetails{
				{
					Message: "error: expected expression",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>",
					Line:    12,
					Column:  34,
				},
			},
		},
		"Typst 0.13.0 multiple errors with paths": {
			StdErr: "error: expected expression\n   ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>:11:53\n   │\n11 │ - Uses stdio; No temporary files need to be created.#\n   │                                                      ^\n\nerror: expected expression\n   ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>:12:34\n   │\n12 │ - Test coverage of most features.#\n   │                                   ^\n\n",
			ExpectedDetails: []typst.ErrorDetails{
				{
					Message: "error: expected expression",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>",
					Line:    11,
					Column:  53,
				},
				{
					Message: "error: expected expression",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>",
					Line:    12,
					Column:  34,
				},
			},
		},
		"Typst 0.13.0 stacked errors with paths": {
			StdErr: "error: expected expression\n  ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\test.typ:1:4\n  │\n1 │ hey#\n  │     ^\n\nhelp: error occurred while importing this module\n   ┌─ \\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>:14:9\n   │\n14 │ #include \"test.typ\"\n   │          ^^^^^^^^^^\n\n",
			ExpectedDetails: []typst.ErrorDetails{
				{
					Message: "error: expected expression",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\test.typ",
					Line:    1,
					Column:  4,
				},
				{
					Message: "help: error occurred while importing this module",
					Path:    "\\\\?\\C:\\Users\\David Vogel\\Desktop\\Synced\\Go\\Libraries\\go-typst\\<stdin>",
					Line:    14,
					Column:  9,
				},
			},
		},
		"Typst 0.13.0 error without path": {
			StdErr: "error: invalid value 'a' for '--pages <PAGES>': not a valid page number\n\nFor more information, try '--help'.\n",
			ExpectedDetails: []typst.ErrorDetails{
				{
					Message: "error: invalid value 'a' for '--pages <PAGES>': not a valid page number",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := typst.ParseStderr(tt.StdErr, nil)

			var typstError *typst.Error
			if errors.As(result, &typstError) {
				if !cmp.Equal(typstError.Details, tt.ExpectedDetails) {
					t.Errorf("Parsed details don't match expected details: %s", cmp.Diff(tt.ExpectedDetails, typstError.Details))
				}
			} else {
				t.Errorf("Parsed error is not of type %T", typstError)
			}

		})
	}
}
