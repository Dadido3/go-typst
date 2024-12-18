package typst_test

import (
	"bytes"
	"image"
	_ "image/png"
	"strconv"
	"testing"

	"github.com/Dadido3/go-typst"
)

// Test basic render functionality.
func TestCLI_Render(t *testing.T) {
	const inches = 1
	const ppi = 144

	cli := typst.CLI{}

	r := bytes.NewBufferString(`#set page(width: ` + strconv.FormatInt(inches, 10) + `in, height: ` + strconv.FormatInt(inches, 10) + `in, margin: (x: 1mm, y: 1mm))
= Test

#lorem(5)`)

	opts := typst.CLIOptions{
		Format: typst.OutputFormatPNG,
		PPI:    ppi,
	}

	var w bytes.Buffer
	if err := cli.Render(r, &w, &opts); err != nil {
		t.Fatalf("Failed to render document: %v.", err)
	}

	imgConf, imgType, err := image.DecodeConfig(&w)
	if err != nil {
		t.Fatalf("Failed to decode image: %v.", err)
	}
	if imgType != "png" {
		t.Fatalf("Resulting image is of type %q, expected %q.", imgType, "png")
	}
	if imgConf.Width != inches*ppi {
		t.Fatalf("Resulting image width is %d, expected %d.", imgConf.Width, inches*ppi)
	}
	if imgConf.Height != inches*ppi {
		t.Fatalf("Resulting image height is %d, expected %d.", imgConf.Height, inches*ppi)
	}
}
