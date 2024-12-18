// Copyright (c) 2024 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestErrors0(t *testing.T) {
	cli := typst.CLI{}

	r := bytes.NewBufferString(`This is a test!`)

	var w bytes.Buffer
	if err := cli.Render(r, &w, nil); err != nil {
		t.Fatalf("Failed to render document: %v", err)
	}
}

func TestErrors1(t *testing.T) {
	cli := typst.CLI{}

	r := bytes.NewBufferString(`This is a test!

#assert(1 < 1, message: "Test")`)

	var w bytes.Buffer
	if err := cli.Render(r, &w, nil); err == nil {
		t.Fatalf("Expected error, but got nil")
	} else {
		var errWithPath *typst.ErrorWithPath
		if errors.As(err, &errWithPath) {
			if errWithPath.Message != "assertion failed: Test" {
				t.Errorf("Expected error with error message %q, got %q", "assertion failed: Test", errWithPath.Message)
			}
			/*if errWithPath.Path != "" {
				t.Errorf("Expected error to point to path %q, got path %q", "", errWithPath.Path)
			}*/
			if errWithPath.Line != 3 {
				t.Errorf("Expected error to point at line %d, got line %d", 3, errWithPath.Line)
			}
			if errWithPath.Column != 1 {
				t.Errorf("Expected error to point at column %d, got column %d", 1, errWithPath.Column)
			}
		} else {
			t.Errorf("Expected error type %T, got %T: %v", errWithPath, err, err)
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
	if err := cli.Render(r, &w, &opts); err == nil {
		t.Fatalf("Expected error, but got nil")
	} else {
		var errTypst *typst.Error
		if errors.As(err, &errTypst) {
			// Don't check the specific error message, as that may change over time.
			// The expected message should be similar to: invalid value 'a' for '--pages <PAGES>': not a valid page number.
			if errTypst.Message == "" {
				t.Errorf("Expected error message, got %q", errTypst.Message)
			}
		} else {
			t.Errorf("Expected error type %T, got %T: %v", errTypst, err, err)
		}
	}
}
