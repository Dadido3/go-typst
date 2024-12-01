//go:build windows

package typst

import "path/filepath"

// We assume the executable is in the current working directory.
var ExecutablePath = "." + string(filepath.Separator) + filepath.Join("typst.exe")
