package typst_test

import (
	"os"
	"testing"

	"github.com/Dadido3/go-typst"
)

func TestCliOptions(t *testing.T) {
	o := typst.CLIOptions{
		FontPaths: []string{"somepath/to/somewhere", "another/to/somewhere"},
	}
	args := o.Args()
	if len(args) != 2 {
		t.Errorf("wrong number of arguments, expected 2, got %d", len(args))
	}
	if "--font-path" != args[0] {
		t.Error("wrong font path option, expected --font-path, got", args[0])
	}
	if "somepath/to/somewhere"+string(os.PathListSeparator)+"another/to/somewhere" != args[1] {
		t.Error("wrong font path option, expected my two paths concatenated, got", args[1])
	}
}
