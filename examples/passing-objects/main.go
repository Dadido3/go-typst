package main

import (
	"log"
	"os"
	"time"

	"github.com/Dadido3/go-typst"
)

// DataEntry contains fake data to be passed to typst.
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
	typstCLI := typst.CLI{}

	r, err := os.Open("template.typ")
	if err != nil {
		log.Panicf("Failed to open template file for reading: %v.", err)
	}
	defer r.Close()

	f, err := os.Create("output.pdf")
	if err != nil {
		log.Panicf("Failed to create output file: %v.", err)
	}
	defer f.Close()

	if err := typstCLI.RenderWithVariables(r, f, nil, map[string]any{"Data": TestData}); err != nil {
		log.Panicf("Failed to render document: %v.", err)
	}
}
