// Copyright (c) 2024 David Vogel
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

// TODO: Add docker support to CLI

type CLI struct {
	ExecutablePath string // The typst executable path can be overridden here. Otherwise the default path will be used.
}

// TODO: Add method for querying the typst version resulting in a semver object

// VersionString returns the version string as returned by typst.
func (c CLI) VersionString() (string, error) {
	// Get path of executable.
	execPath := ExecutablePath
	if c.ExecutablePath != "" {
		execPath = c.ExecutablePath
	}

	cmd := exec.Command(execPath, "--version")

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

// Render takes a typst document from input, and renders it into the output writer.
// The options parameter is optional.
func (c CLI) Render(input io.Reader, output io.Writer, options *CLIOptions) error {
	args := []string{"c"}
	if options != nil {
		args = append(args, options.Args()...)
	}
	args = append(args, "--diagnostic-format", "short", "-", "-")

	// Get path of executable.
	execPath := ExecutablePath
	if c.ExecutablePath != "" {
		execPath = c.ExecutablePath
	}

	cmd := exec.Command(execPath, args...)
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

// Render takes a typst document from input, and renders it into the output writer.
// The options parameter is optional.
//
// Additionally this will inject the given map of variables into the global scope of the typst document.
func (c CLI) RenderWithVariables(input io.Reader, output io.Writer, options *CLIOptions, variables map[string]any) error {
	varBuffer := bytes.Buffer{}

	// TODO: Use io.pipe instead of a bytes.Buffer

	enc := NewVariableEncoder(&varBuffer)
	for k, v := range variables {
		varBuffer.WriteString("#let " + CleanIdentifier(k) + " = ")
		if err := enc.Encode(v); err != nil {
			return fmt.Errorf("failed to encode variables with key %q: %w", k, err)
		}
		varBuffer.WriteRune('\n')
	}

	reader := io.MultiReader(&varBuffer, input)

	return c.Render(reader, output, options)
}
