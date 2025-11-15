package typst

import "io"

// TODO: Add an interface for the Typst caller and let CLI (and later docker and WASM) be implementations of that

// TODO: Add docker support to CLI, by calling docker run instead

// TODO: Add special type "Filename" (or similar) that implements a io.Reader/io.Writer that can be plugged into the input and output parameters of the Compile method

// Caller contains all functions that can be
type Caller interface {
	// VersionString returns the version string as returned by Typst.
	VersionString() (string, error)

	// Compile takes a Typst document from the supplied input reader, and renders it into the output writer.
	// The options parameter is optional, and can be nil.
	Compile(input io.Reader, output io.Writer, options *Options) error
}
