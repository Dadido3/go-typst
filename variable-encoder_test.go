package typst

import (
	"bytes"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestVariableEncoder(t *testing.T) {

	tests := []struct {
		name    string
		params  any
		wantErr bool
		want    string
	}{
		{"nil", nil, false, "none"},
		{"bool false", false, false, "false"},
		{"bool true", true, false, "true"},
		{"int", int(-123), false, "-123"},
		{"int8", int8(-123), false, "-123"},
		{"int16", int16(-123), false, "-123"},
		{"int32", int32(-123), false, "-123"},
		{"int64", int64(-123), false, "-123"},
		{"uint", uint(123), false, "123"},
		{"uint8", uint8(123), false, "123"},
		{"uint16", uint16(123), false, "123"},
		{"uint32", uint32(123), false, "123"},
		{"uint64", uint64(123), false, "123"},
		{"float32", float32(1), false, "1e+00"},
		{"float64", float64(1), false, "1e+00"},
		{"float64 nan", float64(math.NaN()), false, "float.nan"},
		{"float64 +inf", float64(math.Inf(1)), false, "float.inf"},
		{"float64 -inf", float64(math.Inf(-1)), false, "-float.inf"},
		{"string", "Hey!", false, `"Hey!"`},
		{"struct", struct {
			Foo string
			Bar int
		}{"Hey!", 12345}, false, "(\n  Foo: \"Hey!\",\n  Bar: 12345,\n)"},
		{"struct empty", struct{}{}, false, "()"},
		{"struct empty pointer", (*struct{})(nil), false, "none"},
		{"map string string", map[string]string{"Foo": "Bar", "Foo2": "Electric Foogaloo"}, false, "(\n  Foo: \"Bar\",\n  Foo2: \"Electric Foogaloo\",\n)"},
		{"map string string empty", map[string]string{}, false, "()"},
		{"map string string nil", map[string]string(nil), false, "()"},
		{"string array", [5]string{"Foo", "Bar"}, false, `("Foo", "Bar", "", "", "")`},
		{"string slice", []string{"Foo", "Bar"}, false, `("Foo", "Bar")`},
		{"string slice empty", []string{}, false, `()`},
		{"string slice nil", []string(nil), false, `()`},
		{"string slice pointer", &[]string{"Foo", "Bar"}, false, `("Foo", "Bar")`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := bytes.Buffer{}
			vEnc := NewVariableEncoder(&result)

			err := vEnc.Encode(tt.params)
			switch {
			case err != nil && !tt.wantErr:
				t.Fatalf("Failed to encode typst variables: %v", err)
			case err == nil && tt.wantErr:
				t.Fatalf("Expected error, but got none")
			}

			if !cmp.Equal(result.String(), tt.want) {
				t.Errorf("Got unexpected result: %s", cmp.Diff(result.String(), tt.want))
			}
		})
	}
}
