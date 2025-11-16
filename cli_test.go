// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"image"
	_ "image/png"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestCLI_VersionString(t *testing.T) {
	cli := typst.CLI{}

	v, err := cli.VersionString()
	if err != nil {
		t.Fatalf("Failed to get typst version: %v.", err)
	}

	t.Logf("VersionString: %s", v)
}

func TestCLI_Fonts(t *testing.T) {
	caller := typst.CLI{}

	result, err := caller.Fonts(nil)
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) < 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want >= %d.", len(result), 4)
	}
}

func TestCLI_FontsWithOptions(t *testing.T) {
	caller := typst.CLI{}

	result, err := caller.Fonts(&typst.OptionsFonts{IgnoreSystemFonts: true})
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) != 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want %d.", len(result), 4)
	}
}

// Test basic compile functionality.
func TestCLI_Compile(t *testing.T) {
	const inches = 1
	const ppi = 144

	cli := typst.CLI{}

	r := bytes.NewBufferString(`#set page(width: ` + strconv.FormatInt(inches, 10) + `in, height: ` + strconv.FormatInt(inches, 10) + `in, margin: (x: 1mm, y: 1mm))
= Test

#lorem(5)`)

	opts := typst.OptionsCompile{
		Format: typst.OutputFormatPNG,
		PPI:    ppi,
	}

	var w bytes.Buffer
	if err := cli.Compile(r, &w, &opts); err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}

	imgConf, imgType, err := image.DecodeConfig(&w)
	if err != nil {
		t.Fatalf("Failed to decode image: %v.", err)
	}
	if imgType != "png" {
		t.Fatalf("Resulting image is of type %q, expected %q.", imgType, "png")
	}
	if imgConf.Width != inches*ppi {
		t.Fatalf("Resulting image width is %d, expected %d.", imgConf.Width, inches*ppi)
	}
	if imgConf.Height != inches*ppi {
		t.Fatalf("Resulting image height is %d, expected %d.", imgConf.Height, inches*ppi)
	}
}

// Test basic compile functionality with a given working directory.
func TestCLI_CompileWithWorkingDir(t *testing.T) {
	cli := typst.CLI{
		WorkingDirectory: filepath.Join(".", "test-files"),
	}

	r := bytes.NewBufferString(`#import "hello-world-template.typ": template
#show: doc => template()`)

	var w bytes.Buffer
	err := cli.Compile(r, &w, nil)
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
	if w.Available() == 0 {
		t.Errorf("No output was written.")
	}
}
