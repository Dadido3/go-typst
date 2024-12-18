// Copyright (c) 2024 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"testing"
)

func TestCleanIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "_invalid_"},
		{"_", "_invalid_"},
		{"_-", "_-"},
		{"-foo-", "_foo-"},
		{"foo", "foo"},
		{"😊", "_invalid_"},
		{"foo😊", "foo_"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := CleanIdentifier(tt.input); got != tt.want {
				t.Errorf("IsIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"", false},
		{"_", false},
		{"_-", true},
		{"-foo", false},
		{"foo", true},
		{"😊", false},
		{"_😊", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsIdentifier(tt.input); got != tt.want {
				t.Errorf("IsIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
