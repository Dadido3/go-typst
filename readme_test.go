// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dadido3/go-typst"
)

func TestREADME1(t *testing.T) {
	input, output, options := new(bytes.Reader), new(bytes.Buffer), new(typst.OptionsCompile)
	// -----------------------
	typstCaller := typst.CLI{}

	err := typstCaller.Compile(input, output, options)
	// -----------------------
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}

func TestREADME3(t *testing.T) {
	input, output, options := new(bytes.Reader), new(bytes.Buffer), new(typst.OptionsCompile)
	// -----------------------
	typstCaller := typst.Docker{}

	err := typstCaller.Compile(input, output, options)
	// -----------------------
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}

func TestREADME4(t *testing.T) {
	// -----------------------
	typstCaller := typst.Docker{
		Volumes: []string{"./test-files:/markup"},
	}

	r := bytes.NewBufferString(`#include "hello-world.typ"`)

	var w bytes.Buffer
	err := typstCaller.Compile(r, &w, &typst.OptionsCompile{Root: "/markup"})
	// -----------------------
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}

func TestREADME5(t *testing.T) {
	// -----------------------
	typstCaller := typst.Docker{
		Volumes: []string{
			"./test-files:/markup",
			"/usr/share/fonts:/usr/share/fonts",
		},
	}
	// -----------------------

	if _, err := typstCaller.Fonts(nil); err != nil {
		t.Fatalf("Failed to get available fonts: %v.", err)
	}
}
func TestREADME6(t *testing.T) {
	input, output := new(bytes.Reader), new(bytes.Buffer)
	// -----------------------
	typstCaller := typst.Docker{
		Volumes: []string{"./test-files:/fonts"},
	}

	err := typstCaller.Compile(input, output, &typst.OptionsCompile{FontPaths: []string{"/fonts"}})
	// -----------------------
	if err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}

func TestREADME7(t *testing.T) {
	markup := bytes.NewBufferString(`#set page(width: 100mm, height: auto, margin: 5mm)
= go-typst

A library to generate documents and reports by utilizing the command line version of Typst.
#footnote[https://typst.app/]

== Features

- Encoder to convert Go objects into Typst objects which then can be injected into Typst documents.
- Parsing of returned errors into Go error objects.
- Uses stdio; No temporary files need to be created.
- Test coverage of most features.`)

	typstCaller := typst.CLI{}

	f, err := os.Create(filepath.Join(".", "documentation", "images", "readme-example-simple.svg"))
	if err != nil {
		t.Fatalf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	if err := typstCaller.Compile(markup, f, &typst.OptionsCompile{Format: typst.OutputFormatSVG}); err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}

func TestREADME8(t *testing.T) {
	customValues := map[string]any{
		"time":       time.Now(),
		"customText": "Hey there!",
		"struct": struct {
			Foo int
			Bar []string
		}{
			Foo: 123,
			Bar: []string{"this", "is", "a", "string", "slice"},
		},
	}

	// Inject Go values as Typst markup.
	var markup bytes.Buffer
	if err := typst.InjectValues(&markup, customValues); err != nil {
		t.Fatalf("Failed to inject values into Typst markup: %v.", err)
	}

	// Add some Typst markup using the previously injected values.
	markup.WriteString(`#set page(width: 100mm, height: auto, margin: 5mm)
#customText Today's date is #time.display("[year]-[month]-[day]") and the time is #time.display("[hour]:[minute]:[second]").

#struct`)

	f, err := os.Create(filepath.Join(".", "documentation", "images", "readme-example-injection.svg"))
	if err != nil {
		t.Fatalf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	typstCaller := typst.CLI{}
	if err := typstCaller.Compile(&markup, f, &typst.OptionsCompile{Format: typst.OutputFormatSVG}); err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}
