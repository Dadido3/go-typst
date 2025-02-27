// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

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
// For this, just wrap any image.Image with this type before passing it to MarshalValue or a ValueEncoder.
type Image struct{ image.Image }

func (i Image) MarshalTypstValue() ([]byte, error) {
	var buffer bytes.Buffer

	if err := png.Encode(&buffer, i); err != nil {
		return nil, fmt.Errorf("failed to encode image as PNG: %w", err)
	}

	// TODO: Make image encoding more efficient: Use reader/writer, baseXX encoding

	// TODO: Consider using raw pixel encoding instead of PNG

	var buf bytes.Buffer
	buf.WriteString("image.decode(bytes((") // TODO: Pass bytes directly to image once Typst 0.12.0 is not supported anymore
	for _, b := range buffer.Bytes() {
		buf.WriteString(strconv.FormatUint(uint64(b), 10) + ",")
	}
	buf.WriteString(")))")

	return buf.Bytes(), nil
}
