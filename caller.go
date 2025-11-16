// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import "io"

// TODO: Add WASM caller

// TODO: Add special type "Filename" (or similar) that implements a io.Reader/io.Writer that can be plugged into the input and output parameters of the Compile method to signal the use of input/output files instead of readers/writers

// Caller contains all Typst commands that are supported by this library.
type Caller interface {
	// VersionString returns the Typst version as a string.
	VersionString() (string, error)

	// Fonts returns all fonts that are available to Typst.
	Fonts() ([]string, error)

	// Compile takes a Typst document from the supplied input reader, and renders it into the output writer.
	// The options parameter is optional, and can be nil.
	Compile(input io.Reader, output io.Writer, options *OptionsCompile) error
}
