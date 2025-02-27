package main

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/Dadido3/go-typst"
)

// DataEntry contains data to be passed to Typst.
type DataEntry struct {
	Name string
	Size struct{ X, Y, Z float64 }

	Created time.Time
	Numbers []int
}

var TestData = []DataEntry{
	{Name: "Bell", Size: struct{ X, Y, Z float64 }{80, 40, 40}, Created: time.Date(2010, 12, 1, 12, 13, 14, 0, time.UTC), Numbers: []int{1, 2, 3}},
	{Name: "Scissor", Size: struct{ X, Y, Z float64 }{200, 30, 10}, Created: time.Date(2015, 5, 12, 23, 5, 10, 0, time.UTC), Numbers: []int{4, 5, 10, 15}},
	{Name: "Calculator", Size: struct{ X, Y, Z float64 }{150, 80, 8}, Created: time.Date(2016, 6, 10, 12, 15, 0, 0, time.UTC), Numbers: []int{16, 20, 30}},
	{Name: "Key", Size: struct{ X, Y, Z float64 }{25, 10, 2}, Created: time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC), Numbers: []int{100, 199, 205}},
}

func main() {
	var markup bytes.Buffer

	// Inject Go values as Typst markup.
	if err := typst.InjectValues(&markup, map[string]any{"data": TestData, "customText": "This data is coming from a Go application."}); err != nil {
		log.Panicf("Failed to inject values into Typst markup: %v.", err)
	}

	// Import the template and invoke the template function with the custom data.
	// Show is used to replace the current document with whatever content the template function in `template.typ` returns.
	markup.WriteString(`
#import "template.typ": template
#show: doc => template(data, customText)`)

	// Compile the prepared markup with Typst and write the result it into `output.pdf`.
	f, err := os.Create("output.pdf")
	if err != nil {
		log.Panicf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	typstCLI := typst.CLI{}
	if err := typstCLI.Compile(&markup, f, nil); err != nil {
		log.Panicf("Failed to compile document: %v.", err)
	}
}
