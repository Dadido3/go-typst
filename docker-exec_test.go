// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"image"
	"os/exec"
	"strconv"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestDockerExec(t *testing.T) {
	// Just to ensure that there is no container running.
	exec.Command("docker", "stop", "-t", "1", "typst-instance").Run() //nolint:errcheck
	exec.Command("docker", "rm", "typst-instance").Run()              //nolint:errcheck

	if err := exec.Command("docker", "run", "--name", "typst-instance", "-v", "./test-files:/test-files", "-id", "123marvin123/typst").Run(); err != nil {
		t.Fatalf("Failed to run Docker container: %v.", err)
	}
	t.Cleanup(func() {
		exec.Command("docker", "stop", "-t", "1", "typst-instance").Run() //nolint:errcheck
		exec.Command("docker", "rm", "typst-instance").Run()              //nolint:errcheck
	})

	tests := []struct {
		Name     string
		Function func(*testing.T)
	}{
		{"VersionString", dockerExec_VersionString},
		{"Fonts", dockerExec_Fonts},
		{"FontsWithOptions", dockerExec_FontsWithOptions},
		{"FontsWithFontPaths", dockerExec_FontsWithFontPaths},
		{"Compile", dockerExec_Compile},
		{"CompileWithWorkingDir", dockerExec_CompileWithWorkingDir},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			test.Function(t)
		})
	}
}

func dockerExec_VersionString(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
	}

	v, err := typstCaller.VersionString()
	if err != nil {
		t.Fatalf("Failed to get typst version: %v.", err)
	}

	t.Logf("VersionString: %s", v)
}

func dockerExec_Fonts(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
	}

	result, err := typstCaller.Fonts(nil)
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) < 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want >= %d.", len(result), 4)
	}
}

func dockerExec_FontsWithOptions(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
	}

	result, err := typstCaller.Fonts(&typst.OptionsFonts{IgnoreSystemFonts: true})
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) != 4 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want %d.", len(result), 4)
	}
}

func dockerExec_FontsWithFontPaths(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
	}

	result, err := typstCaller.Fonts(&typst.OptionsFonts{IgnoreSystemFonts: true, FontPaths: []string{"/test-files"}})
	if err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
	if len(result) != 5 {
		t.Errorf("Unexpected number of detected fonts. Got %d, want %d.", len(result), 5)
	}
}

// Test basic compile functionality.
func dockerExec_Compile(t *testing.T) {
	const inches = 1
	const ppi = 144

	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
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
func dockerExec_CompileWithWorkingDir(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "typst-instance",
	}

	r := bytes.NewBufferString(`#import "hello-world-template.typ": template
#show: doc => template()`)

	var w bytes.Buffer
	err := typstCaller.Compile(r, &w, &typst.OptionsCompile{Root: "/test-files"})
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
	if w.Available() == 0 {
		t.Errorf("No output was written.")
	}
}

func TestDockerExec_EmptyContainerName(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "",
	}

	_, err := typstCaller.VersionString()
	if err == nil {
		t.Errorf("Expected error, but got nil.")
	}
}

func TestDockerExec_NonRunningContainer(t *testing.T) {
	typstCaller := typst.DockerExec{
		ContainerName: "something-else",
	}

	_, err := typstCaller.VersionString()
	if err == nil {
		t.Errorf("Expected error, but got nil.")
	}
}
