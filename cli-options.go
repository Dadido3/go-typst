// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"os"
	"strconv"
	"time"
)

type OutputFormat string

const (
	OutputFormatAuto OutputFormat = ""

	OutputFormatPDF  OutputFormat = "pdf"
	OutputFormatPNG  OutputFormat = "png"
	OutputFormatSVG  OutputFormat = "svg"
	OutputFormatHTML OutputFormat = "html" // this format is only available since 0.13.0
)

type PDFStandard string

const (
	PDFStandard1_7  PDFStandard = "1.7"  // PDF 1.7
	PDFStandardA_2B PDFStandard = "a-2b" // PDF/A-2b
	PDFStandardA_3B PDFStandard = "a-3b" // PDF/A-3b (Available since Typst 0.13.0 and later)
)

type CLIOptions struct {
	Root              string            // Configures the project root (for absolute paths).
	Input             map[string]string // String key-value pairs visible through `sys.inputs`.
	FontPaths         []string          // Adds additional directories that are recursively searched for fonts.
	IgnoreSystemFonts bool              // Ensures system fonts won't be searched, unless explicitly included via FontPaths.
	CreationTime      time.Time         // The document's creation date. For more information, see https://reproducible-builds.org/specs/source-date-epoch/.
	PackagePath       string            // Custom path to local packages, defaults to system-dependent location.
	PackageCachePath  string            // Custom path to package cache, defaults to system-dependent location.
	Jobs              int               // Number of parallel jobs spawned during compilation, defaults to number of CPUs. Setting it to 1 disables parallelism.

	// Which pages to export. When unspecified, all document pages are exported.
	//
	// Pages to export are separated by commas, and can be either simple page numbers (e.g. '2,5' to export only pages 2 and 5) or page ranges (e.g. '2,3-6,8-' to export page 2, pages 3 to 6 (inclusive), page 8 and any pages after it).
	//
	// Page numbers are one-indexed and correspond to real page numbers in the document (therefore not being affected by the document's page counter).
	Pages string

	Format OutputFormat // The format of the output file, inferred from the extension by default.
	PPI    int          // The PPI (pixels per inch) to use for PNG export. Defaults to 144.

	// One (or multiple) PDF standards that Typst will enforce conformance with.
	//
	// See typst.PDFStandard for possible values.
	PDFStandards []PDFStandard

	Custom []string // Custom command line options go here.
}

// Args returns a list of CLI arguments that should be passed to the executable.
func (c *CLIOptions) Args() (result []string) {
	if c.Root != "" {
		result = append(result, "--root", c.Root)
	}

	for key, value := range c.Input {
		result = append(result, "--input", key+"="+value)
	}

	if len(c.FontPaths) > 0 {
		var paths string
		for i, path := range c.FontPaths {
			if i > 0 {
				paths += string(os.PathListSeparator)
			}
			paths += path
		}
		result = append(result, "--font-path", paths)
	}

	if c.IgnoreSystemFonts {
		result = append(result, "--ignore-system-fonts")
	}

	if !c.CreationTime.IsZero() {
		result = append(result, "--creation-timestamp", strconv.FormatInt(c.CreationTime.Unix(), 10))
	}

	if c.PackagePath != "" {
		result = append(result, "--package-path", c.PackagePath)
	}

	if c.PackageCachePath != "" {
		result = append(result, "--package-cache-path", c.PackageCachePath)
	}

	if c.Jobs > 0 {
		result = append(result, "-j", strconv.FormatInt(int64(c.Jobs), 10))
	}

	if c.Pages != "" {
		result = append(result, "--pages", c.Pages)
	}

	if c.Format != OutputFormatAuto {
		result = append(result, "-f", string(c.Format))
		if c.Format == OutputFormatHTML {
			// this is specific to version 0.13.0 where html
			// is a feature than need explicit activation
			// we must remove this when html becomes standard
			result = append(result, "--features", "html")
		}
	}

	if c.PPI > 0 {
		result = append(result, "--ppi", strconv.FormatInt(int64(c.PPI), 10))
	}

	if len(c.PDFStandards) > 0 {
		var standards string
		for i, standard := range c.PDFStandards {
			if i > 0 {
				standards += ","
			}
			standards += string(standard)
		}
		result = append(result, "--pdf-standard", standards)
	}

	result = append(result, c.Custom...)

	return
}
