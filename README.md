# go-typst [![test](https://github.com/Dadido3/go-typst/actions/workflows/test.yml/badge.svg)](https://github.com/Dadido3/go-typst/actions/workflows/test.yml)

A library to generate documents and reports by utilizing the command line version of [typst].

## Features

- Encoder to convert go objects into typst objects which then can be injected into typst documents.
- Parsing of stderr into an go error object.
- Uses stdio; No temporary files need to be created.
- Good test coverage of features.

## Installation

1. Use `go get github.com/Dadido3/go-typst` inside of your project to add this module to your project.
2. Install typst by following [the instructions in the typst repository].

## Runtime requirements

You need to have [typst] installed on any machine that you want to run your go project on.
You can install it by following [the instructions in the typst repository].

## Usage

Here we will create a simple PDF document by passing a reader with typst markup into `typstCLI.Render` and then let it write the resulting PDF data into a file:

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

    if err := typstCLI.Render(r, f, nil); err != nil {
        t.Fatalf("Failed to render document: %v.", err)
    }
}
```

The resulting document will look like this:

![readme-1.svg](documentation/images/readme-1.svg)

[the instructions in the typst repository]: https://github.com/typst/typst?tab=readme-ov-file#installation
[typst]: https://typst.app/
