// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type CLI struct {
	ExecutablePath   string // The Typst executable path can be overridden here. Otherwise the default path will be used.
	WorkingDirectory string // The path where the Typst executable is run in. When left empty, the Typst executable will be run in the process's current directory.
}

// Ensure that CLI implements the Caller interface.
var _ Caller = CLI{}

// VersionString returns the version string as returned by Typst.
func (c CLI) VersionString() (string, error) {
	// Get path of executable.
	execPath := ExecutablePath
	if c.ExecutablePath != "" {
		execPath = c.ExecutablePath
	}
	if execPath == "" {
		return "", fmt.Errorf("not supported on this platform")
	}

	cmd := exec.Command(execPath, "--version")
	cmd.Dir = c.WorkingDirectory // This doesn't do anything, but we will do it anyways for consistency.

	var output, errBuffer bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return "", ParseStderr(errBuffer.String(), err)
		default:
			return "", err
		}
	}

	return output.String(), nil
}

// Compile takes a Typst document from input, and renders it into the output writer.
// The options parameter is optional.
func (c CLI) Compile(input io.Reader, output io.Writer, options *Options) error {
	args := []string{"c"}
	if options != nil {
		args = append(args, options.Args()...)
	}
	args = append(args, "--diagnostic-format", "human", "-", "-")

	// Get path of executable.
	execPath := ExecutablePath
	if c.ExecutablePath != "" {
		execPath = c.ExecutablePath
	}
	if execPath == "" {
		return fmt.Errorf("not supported on this platform")
	}

	cmd := exec.Command(execPath, args...)
	cmd.Dir = c.WorkingDirectory
	cmd.Stdin = input
	cmd.Stdout = output

	errBuffer := bytes.Buffer{}
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return ParseStderr(errBuffer.String(), err)
		default:
			return err
		}
	}

	return nil
}

// CompileWithVariables takes a Typst document from input, and renders it into the output writer.
// The options parameter is optional.
//
// Additionally this will inject the given map of variables into the global scope of the Typst document.
//
// Deprecated: You should use typst.InjectValues in combination with the normal Compile method instead.
func (c CLI) CompileWithVariables(input io.Reader, output io.Writer, options *Options, variables map[string]any) error {
	varBuffer := bytes.Buffer{}

	if err := InjectValues(&varBuffer, variables); err != nil {
		return fmt.Errorf("failed to inject values into Typst markup: %w", err)
	}

	reader := io.MultiReader(&varBuffer, input)

	return c.Compile(reader, output, options)
}
