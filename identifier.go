// Copyright (c) 2024 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"unicode/utf8"

	"github.com/smasher164/xid"
)

// CleanIdentifier will return the input cleaned up in a way so that it can safely be used as a typst identifier.
// This function will replace all illegal characters, which means collisions are possible in some cases.
//
// See https://github.com/typst/typst/blob/76c24ee6e35715cd14bb892d7b6b8d775c680bf7/crates/typst-syntax/src/lexer.rs#L932 for details.
func CleanIdentifier(input string) string {
	dst := make([]byte, 0, len(input))

	for i, r := range input {
		if i == 0 {
			// Handle first rune of input.
			switch {
			case xid.Start(r), r == '_':
				dst = utf8.AppendRune(dst, r)
			default:
				dst = append(dst, '_')
			}
		} else {
			// Handle all other runes of input.
			switch {
			case xid.Continue(r), r == '_', r == '-':
				dst = utf8.AppendRune(dst, r)
			default:
				dst = append(dst, '_')
			}
		}
	}

	// Don't allow empty identifiers.
	// We can't use a single placeholder ("_"), as it will cause errors when used in dictionaries.
	result := string(dst)
	if result == "_" || result == "" {
		return "_invalid_"
	}

	return string(dst)
}

// IsIdentifier will return whether input is a valid typst identifier.
//
// See https://github.com/typst/typst/blob/76c24ee6e35715cd14bb892d7b6b8d775c680bf7/crates/typst-syntax/src/lexer.rs#L932 for details.
func IsIdentifier(input string) bool {
	// Identifiers can't be empty.
	// We will also disallow a single underscore.
	if input == "" || input == "_" {
		return false
	}

	for i, r := range input {
		if i == 0 {
			// Handle first rune of input.
			switch {
			case xid.Start(r), r == '_':
			default:
				return false
			}
		} else {
			// Handle all other runes of input.
			switch {
			case xid.Continue(r), r == '_', r == '-':
			default:
				return false
			}
		}
	}

	return true
}
