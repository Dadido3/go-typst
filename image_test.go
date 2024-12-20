package typst_test

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"testing"

	"github.com/Dadido3/go-typst"
)

type testImage struct {
	Rect image.Rectangle
}

func (p *testImage) ColorModel() color.Model { return color.RGBAModel }

func (p *testImage) Bounds() image.Rectangle { return p.Rect }

func (p *testImage) At(x, y int) color.Color { return p.RGBAAt(x, y) }

func (p *testImage) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	return color.RGBA{uint8(x), uint8(y), uint8(x + y), 255}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *testImage) Opaque() bool {
	return true
}

func TestImage(t *testing.T) {
	img := &testImage{image.Rect(0, 0, 255, 255)}

	// Wrap image.
	typstImage := typst.Image{img}

	cli := typst.CLI{}

	r := bytes.NewBufferString(`= Image test

#TestImage`)

	if err := cli.CompileWithVariables(r, io.Discard, nil, map[string]any{"TestImage": typstImage}); err != nil {
		t.Fatalf("Failed to compile document: %v.", err)
	}
}
