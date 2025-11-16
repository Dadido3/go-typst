// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Dadido3/go-typst"
)

// Returns the TYPST_DOCKER_IMAGE environment variable.
// If that's not set, it will return an empty string, which makes the tests default to typst.DockerDefaultImage.
func typstDockerImage() string {
	return os.Getenv("TYPST_DOCKER_IMAGE")
}

func TestDocker_VersionString(t *testing.T) {
	caller := typst.Docker{
		Image: typstDockerImage(),
	}

	v, err := caller.VersionString()
	if err != nil {
		t.Fatalf("Failed to get typst version: %v.", err)
	}

	t.Logf("VersionString: %s", v)
}

func TestDocker_Fonts(t *testing.T) {
	caller := typst.Docker{
		Image: typstDockerImage(),
	}

	result, err := caller.Fonts(nil)
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) < 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want >= %d.", len(result), 4)
	}
}

func TestDocker_FontsWithOptions(t *testing.T) {
	caller := typst.Docker{
		Image: typstDockerImage(),
	}

	result, err := caller.Fonts(&typst.OptionsFonts{IgnoreSystemFonts: true})
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) != 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want %d.", len(result), 4)
	}
}

func TestDocker_FontsWithFontPaths(t *testing.T) {
	caller := typst.Docker{
		Image:   typstDockerImage(),
		Volumes: []string{"./test-files:/fonts"},
	}

	result, err := caller.Fonts(&typst.OptionsFonts{IgnoreSystemFonts: true, FontPaths: []string{"/fonts"}})
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) != 5 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want %d.", len(result), 5)
	}
}

// Test basic compile functionality.
func TestDocker_Compile(t *testing.T) {
	const inches = 1
	const ppi = 144

	typstCaller := typst.Docker{
		Image: typstDockerImage(),
	}

	r := bytes.NewBufferString(`#set page(width: ` + strconv.FormatInt(inches, 10) + `in, height: ` + strconv.FormatInt(inches, 10) + `in, margin: (x: 1mm, y: 1mm))
= Test

#lorem(5)`)

	opts := typst.OptionsCompile{
		Format: typst.OutputFormatPNG,
		PPI:    ppi,
	}

	var w bytes.Buffer
	if err := typstCaller.Compile(r, &w, &opts); err != nil {
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
func TestDocker_CompileWithWorkingDir(t *testing.T) {
	typstCaller := typst.Docker{
		Image:            typstDockerImage(),
		WorkingDirectory: filepath.Join(".", "test-files"),
		Volumes:          []string{".:/markup"},
	}

	r := bytes.NewBufferString(`#import "hello-world-template.typ": template
#show: doc => template()`)

	var w bytes.Buffer
	err := typstCaller.Compile(r, &w, &typst.OptionsCompile{Root: "/markup"})
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
	if w.Available() == 0 {
		t.Errorf("No output was written.")
	}
}
