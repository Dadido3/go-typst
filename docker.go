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
//
// This uses docker run to automatically pull and run a container.
// Therefore the container will start and stop automatically.
// To have more control over the lifetime of a Docker container see typst.DockerExec.
type Docker struct {
	Image            string // The image to use, defaults to the latest supported official Typst Docker image if left empty. See: typst.DockerDefaultImage.
	WorkingDirectory string // The working directory of Docker. When left empty, Docker will be run with the process's current working directory.

	// Additional bind-mounts or volumes that are passed via "--volume" flag to Docker.
	// For details, see: https://docs.docker.com/engine/storage/volumes/#syntax
	//
	// Example:
	//	typst.Docker{Volumes: []string{".:/markup"}} // This bind mounts the current working directory to "/markup" inside the container.
	//	typst.Docker{Volumes: []string{"/usr/share/fonts:/usr/share/fonts"}} // This makes all system fonts available to Typst running inside the container.
	Volumes []string

	// Custom "docker run" command line options go here.
	// For all available options, see: https://docs.docker.com/reference/cli/docker/container/run/
	//
	// Example:
	//	typst.Docker{Custom: []string{"--user", "1000"}} // Use a non-root user inside the docker container.
	Custom []string // Custom "docker run" command line options go here.
}

// Ensure that Docker implements the Caller interface.
var _ Caller = Docker{}

// args returns docker related arguments.
func (d Docker) args() []string {
	image := DockerDefaultImage
	if d.Image != "" {
		image = d.Image
	}

	// Argument -i is needed for stdio to work.
	args := []string{"run", "-i"}

	args = append(args, d.Custom...)

	// Add mounts.
	for _, volume := range d.Volumes {
		args = append(args, "-v", volume)
	}

	// Which docker image to use.
	args = append(args, image)

	return args
}

// VersionString returns the Typst version as a string.
func (d Docker) VersionString() (string, error) {
	args := append(d.args(), "--version")

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
func (d Docker) Fonts(options *OptionsFonts) ([]string, error) {
	args := d.args()

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
func (d Docker) Compile(input io.Reader, output io.Writer, options *OptionsCompile) error {
	args := d.args()

	if options == nil {
		options = new(OptionsCompile)
	}
	args = append(args, options.Args()...)

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
