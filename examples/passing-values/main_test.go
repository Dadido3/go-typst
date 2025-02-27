package main

import (
	"testing"
)

// Run the example as a test.
func TestMain(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()

	main()
}
