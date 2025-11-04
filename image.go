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

// Image can be used to encode any image.Image into a Typst image.
//
// For this, just wrap any image.Image with this type before passing it to MarshalValue or a ValueEncoder:
//
//	typstImage := typst.Image{img}
//	typst.InjectValues(&r, map[string]any{"TestImage": typstImage})
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

// ImageRaw can be used to pass the raw data of any image to Typst.
// This will pass the raw byte values of a PNG, JPEG or any other image format that is supported by Typst.
//
// For this, just wrap any byte slice with this type before passing it to MarshalValue or a ValueEncoder:
//
//	typstImage := typst.ImageRaw(bufferPNG)
//	typst.InjectValues(&r, map[string]any{"TestImage": typstImage})
type ImageRaw []byte

func (i ImageRaw) MarshalTypstValue() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("image.decode(bytes((") // TODO: Pass bytes directly to image once Typst 0.12.0 is not supported anymore
	for _, b := range i {
		buf.WriteString(strconv.FormatUint(uint64(b), 10) + ",")
	}
	buf.WriteString(")))")

	return buf.Bytes(), nil
}
