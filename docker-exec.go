// Copyright (c) 2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

// Theoretically it's possible to use the Docker SDK directly:
// https://docs.docker.com/reference/api/engine/sdk/examples/
// But that dependency is unnecessarily huge, therefore we will just call the Docker executable.

// DockerExec allows you to invoke Typst commands in a running Docker container.
//
// This uses docker exec, and therefore needs you to set up a running container beforehand.
// For a less complex setup see typst.Docker.
type DockerExec struct {
	ContainerName string // The name of the running container you want to invoke Typst in.
	TypstPath     string // The path to the Typst executable inside of the container. Defaults to `typst` if left empty.

	// Custom "docker exec" command line options go here.
	// For all available options, see: https://docs.docker.com/reference/cli/docker/container/exec/
	//
	// Example:
	//	typst.DockerExec{Custom: []string{"--user", "1000"}} // Use a non-root user inside the docker container.
	Custom []string
}

// Ensure that DockerExec implements the Caller interface.
var _ Caller = DockerExec{}

// args returns docker related arguments.
func (d DockerExec) args() ([]string, error) {
	if d.ContainerName == "" {
		return nil, fmt.Errorf("the provided ContainerName field is empty")
	}

	typstPath := "typst"
	if d.TypstPath != "" {
		typstPath = d.TypstPath
	}

	// Argument -i is needed for stdio to work.
	args := []string{"exec", "-i"}

	args = append(args, d.Custom...)

	args = append(args, d.ContainerName, typstPath)

	return args, nil
}

// VersionString returns the Typst version as a string.
func (d DockerExec) VersionString() (string, error) {
	args, err := d.args()
	if err != nil {
		return "", err
	}
	args = append(args, "--version")

	cmd := exec.Command("docker", args...)

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

// Fonts returns all fonts that are available to Typst.
// The options parameter is optional, and can be nil.
func (d DockerExec) Fonts(options *OptionsFonts) ([]string, error) {
	args, err := d.args()
	if err != nil {
		return nil, err
	}

	if options == nil {
		options = new(OptionsFonts)
	}
	args = append(args, options.Args()...)

	cmd := exec.Command("docker", args...)

	var output, errBuffer bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return nil, ParseStderr(errBuffer.String(), err)
		default:
			return nil, err
		}
	}

	var result []string
	scanner := bufio.NewScanner(&output)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, nil
}

// Compile takes a Typst document from input, and renders it into the output writer.
// The options parameter is optional, and can be nil.
func (d DockerExec) Compile(input io.Reader, output io.Writer, options *OptionsCompile) error {
	args, err := d.args()
	if err != nil {
		return err
	}

	if options == nil {
		options = new(OptionsCompile)
	}
	args = append(args, options.Args()...)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = input
	cmd.Stdout = output

	errBuffer := bytes.Buffer{}
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			if err.ExitCode() >= 125 {
				// Most likely docker related error.
				// TODO: Find a better way to distinguish between Typst or Docker errors.
				return fmt.Errorf("exit code %d: %s", err.ExitCode(), errBuffer.String())
			} else {
				// Typst related error.
				return ParseStderr(errBuffer.String(), err)
			}
		default:
			return err
		}
	}

	return nil
}
