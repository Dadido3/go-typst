// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:build !(windows || unix)

package typst

// The path to the Typst executable.
// We leave that empty as we don't support this platform for now.
var ExecutablePath = ""
