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

// The default Docker image to use.
// This is the latest supported version of Typst.
const DockerDefaultImage = "ghcr.io/typst/typst:0.14.0"

// Docker allows you to invoke commands on a Typst Docker image.
type Docker struct {
	Image            string // The image to use, defaults to the latest supported offical Typst Docker image if left empty. See: typst.DockerDefaultImage.
	WorkingDirectory string // The working directory of Docker. When left empty, Docker will be run with the process's current working directory.

	// Additional bind-mounts or volumes that are passed via "--volume" flag to Docker.
	// For details, see: https://docs.docker.com/engine/storage/volumes/#syntax
	//
	// Example:
	//	typst.Docker{Volumes: []string{".:/markup"}} // This bind mounts the current working directory to "/markup" inside the container.
	//	typst.Docker{Volumes: []string{"/usr/share/fonts:/usr/share/fonts"}} // This makes all system fonts available to Typst running inside the container.
	Volumes []string
}

// Ensure that Docker implements the Caller interface.
var _ Caller = Docker{}

// VersionString returns the version string as returned by Typst.
func (d Docker) VersionString() (string, error) {
	image := DockerDefaultImage
	if d.Image != "" {
		image = d.Image
	}

	cmd := exec.Command("docker", "run", "-i", image, "--version")

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
func (d Docker) Fonts() ([]string, error) {
	image := DockerDefaultImage
	if d.Image != "" {
		image = d.Image
	}

	cmd := exec.Command("docker", "run", "-i", image, "fonts")

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
// The options parameter is optional.
func (d Docker) Compile(input io.Reader, output io.Writer, options *Options) error {
	image := DockerDefaultImage
	if d.Image != "" {
		image = d.Image
	}

	// Argument -i is needed for stdio to work.
	args := []string{"run", "-i"}

	// Add mounts.
	for _, volume := range d.Volumes {
		args = append(args, "-v", volume)
	}

	args = append(args, image)

	// From here on come Typst arguments.

	args = append(args, "c")
	if options != nil {
		args = append(args, options.Args()...)
	}
	args = append(args, "--diagnostic-format", "human", "-", "-") // TODO: Move these default arguments into Options

	cmd := exec.Command("docker", args...)
	cmd.Dir = d.WorkingDirectory
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
