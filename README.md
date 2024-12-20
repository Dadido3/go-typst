# go-typst [![test](https://github.com/Dadido3/go-typst/actions/workflows/test.yml/badge.svg)](https://github.com/Dadido3/go-typst/actions/workflows/test.yml)

A library to generate documents and reports by utilizing the command line version of [typst].

This is basically a binding to typst-cli which exposes functions needed to compile documents into different formats like PDF, SVG or PNG. The goal is to make using typst as simple and "go like" as possible.

This module, along with typst itself, is a work in progress.
The API may change, and compatibility with different typst versions are also not set in stone.
There is no way to prevent this as long as typst has breaking changes.
To mitigate problems arising from this, most of the functionality is unit tested against different typst releases.
The supported and tested versions right now are:

- Typst 0.12.0

## Features

- PDF, SVG and PNG generation.
- All typst-cli parameters are [available as a struct](cli-options.go), which makes it easy to discover all available options.
- Encoder to convert go values into typst markup which can be injected into typst documents. This includes image.Image by using the [Image wrapper](image.go).
- Any stderr will be returned as go error value, including line number, column and file path of the error.
- Uses stdio; No temporary files will be created.
- Good unit test coverage.

## Installation

1. Use `go get github.com/Dadido3/go-typst` inside of your project to add this module to your project.
2. Install typst by following [the instructions in the typst repository].

## Runtime requirements

You need to have [typst] installed on any machine that you want to run your go project on.
You can install it by following [the instructions in the typst repository].

## Usage

Here we will create a simple PDF document by passing a reader with typst markup into `typstCLI.Compile` and then let it write the resulting PDF data into a file:

```go
func main() {
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

    if err := typstCLI.Compile(r, f, nil); err != nil {
        t.Fatalf("Failed to compile document: %v.", err)
    }
}
```

The resulting document will look like this:

![readme-1.svg](documentation/images/readme-1.svg)

[the instructions in the typst repository]: https://github.com/typst/typst?tab=readme-ov-file#installation
[typst]: https://typst.app/
