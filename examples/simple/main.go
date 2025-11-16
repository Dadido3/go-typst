package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Dadido3/go-typst"
)

func main() {
	// Convert a time.Time value into Typst markup.
	date, err := typst.MarshalValue(time.Now())
	if err != nil {
		log.Panicf("Failed to marshal date into Typst markup: %v", err)
	}

	// Write Typst markup into buffer.
	var markup bytes.Buffer
	fmt.Fprintf(&markup, `= Hello world

This document was created at #%s.display() using typst-go.`, date)

	// Compile the prepared markup with Typst and write the result it into `output.pdf`.
	f, err := os.Create("output.pdf")
	if err != nil {
		log.Panicf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	typstCaller := typst.CLI{}
	if err := typstCaller.Compile(&markup, f, nil); err != nil {
		log.Panic("failed to compile document: %w", err)
	}
}
