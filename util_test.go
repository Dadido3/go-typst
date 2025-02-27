package typst

import (
	"bytes"
	"testing"
	"time"
)

func TestInjectValues(t *testing.T) {
	type args struct {
		values map[string]any
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{"empty", args{values: nil}, "", false},
		{"nil", args{values: map[string]any{"foo": nil}}, "#let foo = none\n", false},
		{"example", args{values: map[string]any{"foo": 1, "bar": 60 * time.Second}}, "#let bar = duration(seconds: 60)\n#let foo = 1\n", false},
		{"invalid identifier", args{values: map[string]any{"foo😀": 1}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			if err := InjectValues(output, tt.args.values); (err != nil) != tt.wantErr {
				t.Errorf("InjectValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("InjectValues() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
