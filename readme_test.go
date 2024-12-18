package typst_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestREADME1(t *testing.T) {
	r := bytes.NewBufferString(`#set page(width: 100mm, height: auto, margin: 5mm)
= go-typst

A library to generate documents and reports by utilizing the command line version of typst.
#footnote[https://typst.app/]

== Features

- Encoder to convert go objects into typst objects which then can be injected into typst documents.
- Parsing of returned errors into go error objects.
- Uses stdio; No temporary files need to be created.
- Test coverage of most features.`)

	typstCLI := typst.CLI{}

	f, err := os.Create("output.pdf")
	if err != nil {
		t.Fatalf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	if err := typstCLI.Render(r, f, nil); err != nil {
		t.Fatalf("Failed to render document: %v.", err)
	}
}
