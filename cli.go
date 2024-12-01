package typst

import (
	"bytes"
	"io"
	"os/exec"
)

type CLI struct {
	//ExecutablePath string
}

func (c CLI) Render(input io.Reader, output io.Writer) error {
	cmd := exec.Command(ExecutablePath, "c", "-", "-")
	cmd.Stdin = input
	cmd.Stdout = output

	errBuffer := bytes.Buffer{}
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		switch err := err.(type) {
		case *exec.ExitError:
			return NewError(errBuffer.String(), err)
		default:
			return err
		}
	}

	return nil
}

func (c CLI) RenderWithVariables(input io.Reader, output io.Writer, variables map[string]any) error {
	reader := io.MultiReader(nil, input)

	return c.Render(reader, output)
}
