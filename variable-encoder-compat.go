// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import "io"

// This exists for compatibility reasons.

// Deprecated: Use NewValueEncoder instead, as this will be removed in a future version.
func NewVariableEncoder(w io.Writer) *ValueEncoder { return NewValueEncoder(w) }

// Deprecated: Use MarshalValue instead, as this will be removed in a future version.
func MarshalVariable(v any) ([]byte, error) { return MarshalValue(v) }

// Deprecated: Use ValueMarshaler interface instead, as this will be removed in a future version.
type VariableMarshaler interface {
	MarshalTypstVariable() ([]byte, error)
}
