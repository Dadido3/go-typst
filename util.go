// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"fmt"
	"io"
	"maps"
	"slices"
)

// InjectValues will write the given key-value pairs as Typst markup into output.
// This can be used to inject Go values into typst documents.
//
// Every key in values needs to be a valid identifier, otherwise this function will return an error.
// Every value in values will be marshaled according to ValueEncoder into equivalent Typst markup.
//
// Passing {"foo": 1, "bar": 60 * time.Second} as values will produce the following output:
//
//	#let foo = 1
//	#let bar = duration(seconds: 60)
func InjectValues(output io.Writer, values map[string]any) error {
	enc := NewValueEncoder(output)

	// We will have to iterate over the sorted list of map keys.
	// Otherwise the output is not deterministic, and tests will fail randomly.
	for _, k := range slices.Sorted(maps.Keys(values)) {
		v := values[k]
		if !IsIdentifier(k) {
			return fmt.Errorf("%q is not a valid identifier", k)
		}
		if _, err := output.Write([]byte("#let " + CleanIdentifier(k) + " = ")); err != nil {
			return err
		}
		if err := enc.Encode(v); err != nil {
			return fmt.Errorf("failed to encode values with key %q: %w", k, err)
		}
		if _, err := output.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}
