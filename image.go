package typst

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strconv"
)

// Image can be used to encode any image.Image into a typst image.
//
// For this, just wrap any image.Image with this type before passing it to MarshalVariable or a VariableEncoder.
type Image struct{ image.Image }

func (i Image) MarshalTypstVariable() ([]byte, error) {
	var buffer bytes.Buffer

	if err := png.Encode(&buffer, i); err != nil {
		return nil, fmt.Errorf("failed to encode image as PNG: %w", err)
	}

	// TODO: Make image encoding more efficient: Use reader/writer, baseXX encoding

	var buf bytes.Buffer
	buf.WriteString("image.decode(bytes((")
	for _, b := range buffer.Bytes() {
		buf.WriteString(strconv.FormatUint(uint64(b), 10) + ",")
	}
	buf.WriteString(")))")

	return buf.Bytes(), nil
}
