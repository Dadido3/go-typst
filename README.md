# go-typst

A library to generate documents and reports by utilizing the command line version of [typst].

Features include:

- Encoder to convert go objects into typst objects which then can be injected into typst documents.
- Parsing of returned errors into go error objects.
- Uses stdio; No temporary files need to be created.
- Test coverage of most features.

## Installation

1. Use `go get github.com/Dadido3/go-typst` inside of your module to add this library to your project.
2. Install typst by following [the instructions in the typst repository].

## Runtime requirements

You need to have [typst] installed on any machine that you want to run your go project on.
You can install it by following [the instructions in the typst repository].

## Usage

ToDo

[the instructions in the typst repository]: https://github.com/typst/typst?tab=readme-ov-file#installation
[typst]: https://typst.app/
