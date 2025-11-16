// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"os"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestOptions(t *testing.T) {
	o := typst.OptionsCompile{
		FontPaths: []string{"somepath/to/somewhere", "another/to/somewhere"},
	}
	args := o.Args()
	if len(args) != 2 {
		t.Errorf("wrong number of arguments, expected 2, got %d", len(args))
	}
	if args[0] != "--font-path" {
		t.Error("wrong font path option, expected --font-path, got", args[0])
	}
	if args[1] != "somepath/to/somewhere"+string(os.PathListSeparator)+"another/to/somewhere" {
		t.Error("wrong font path option, expected my two paths concatenated, got", args[1])
	}
}
