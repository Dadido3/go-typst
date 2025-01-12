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

A library to generate documents and reports by utilizing the command line version of Typst.
#footnote[https://typst.app/]

== Features

- Encoder to convert Go objects into Typst objects which then can be injected into Typst documents.
- Parsing of returned errors into Go error objects.
- Uses stdio; No temporary files need to be created.
- Test coverage of most features.`)

	typstCLI := typst.CLI{}

	f, err := os.Create("output.pdf")
	if err != nil {
		t.Fatalf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	if err := typstCLI.Compile(r, f, nil); err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}
