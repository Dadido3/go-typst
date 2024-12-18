// Copyright (c) 2024 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"slices"
	"strconv"
	"time"
)

// VariableMarshaler can be implemented by types to support custom typst marshaling.
type VariableMarshaler interface {
	MarshalTypstVariable() ([]byte, error)
}

type VariableEncoder struct {
	indentLevel int

	writer io.Writer
}

// NewVariableEncoder returns a new encoder that writes into w.
func NewVariableEncoder(w io.Writer) *VariableEncoder {
	return &VariableEncoder{
		writer: w,
	}
}

func (e *VariableEncoder) Encode(v any) error {
	return e.marshal(reflect.ValueOf(v))
}

func (e *VariableEncoder) writeString(s string) error {
	return e.writeBytes([]byte(s))
}

func (e *VariableEncoder) writeStringLiteral(s []byte) error {
	dst := make([]byte, 0, len(s)+5)

	dst = append(dst, '"')

	for _, r := range s {
		switch r {
		case '\\', '"':
			dst = append(dst, '\\', r)
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		default:
			dst = append(dst, r)
		}
	}

	dst = append(dst, '"')

	return e.writeBytes(dst)
}

func (e *VariableEncoder) writeBytes(b []byte) error {
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("failed to write into writer: %w", err)
	}

	return nil
}

func (e *VariableEncoder) writeIndentationCharacters() error {
	return e.writeBytes(slices.Repeat([]byte{' ', ' '}, e.indentLevel))
}

func (e *VariableEncoder) marshal(v reflect.Value) error {
	if !v.IsValid() {
		return e.writeString("none")
		//return fmt.Errorf("invalid reflect.Value %v", v)
	}

	t := v.Type()

	switch i := v.Interface().(type) {
	case time.Time:
		if err := e.encodeTime(i); err != nil {
			return err
		}
		return nil
	case *time.Time:
		if i == nil {
			e.writeString("none")
			return nil
		}
		if err := e.encodeTime(*i); err != nil {
			return err
		}
		return nil
	case time.Duration:
		if err := e.encodeDuration(i); err != nil {
			return err
		}
		return nil
	case *time.Duration:
		if i == nil {
			e.writeString("none")
			return nil
		}
		if err := e.encodeDuration(*i); err != nil {
			return err
		}
		return nil
	}

	// TODO: Handle images, maybe create a wrapper type that does this

	if t.Implements(reflect.TypeFor[VariableMarshaler]()) {
		if m, ok := v.Interface().(VariableMarshaler); ok {
			bytes, err := m.MarshalTypstVariable()
			if err != nil {
				return fmt.Errorf("error calling MarshalTypstVariable for type %s: %w", t.String(), err)
			}
			return e.writeBytes(bytes)
		}
		return e.writeString("none")
	}

	if t.Implements(reflect.TypeFor[encoding.TextMarshaler]()) {
		if m, ok := v.Interface().(encoding.TextMarshaler); ok {
			b, err := m.MarshalText()
			if err != nil {
				return fmt.Errorf("error calling MarshalText for type %s: %w", t.String(), err)
			}
			return e.writeStringLiteral(b)
		}
		return e.writeString("none")
	}

	var err error
	switch t.Kind() {
	case reflect.Bool:
		err = e.writeString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		err = e.writeString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		err = e.writeString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		switch {
		case math.IsNaN(f):
			err = e.writeString("float.nan")
		case math.IsInf(f, 1):
			err = e.writeString("float.inf")
		case math.IsInf(f, -1):
			err = e.writeString("-float.inf")
		default:
			err = e.writeString(strconv.FormatFloat(f, 'e', -1, 64))
		}
	case reflect.String:
		return e.encodeString(v)
	case reflect.Interface, reflect.Pointer:
		if v.IsNil() {
			return e.writeString("none")
		}
		return e.marshal(v.Elem())
	case reflect.Map:
		return e.encodeMap(v)
	case reflect.Struct:
		return e.encodeStruct(v, t)
	case reflect.Slice:
		return e.encodeSlice(v, t)
	case reflect.Array:
		return e.encodeArray(v)
	default:
		return fmt.Errorf("unsupported type %q", t.String())
	}

	return err
}

func (e *VariableEncoder) encodeString(v reflect.Value) error {
	return e.writeStringLiteral([]byte(v.String()))
}

func (e *VariableEncoder) encodeStruct(v reflect.Value, t reflect.Type) error {
	if v.NumField() == 0 {
		return e.writeString("()")
	}

	if err := e.writeString("(\n"); err != nil {
		return err
	}

	e.indentLevel++

	for i := 0; i < t.NumField(); i++ {
		ft, fv := t.Field(i), v.Field(i)
		if ft.PkgPath == "" { // Ignore unexported fields.
			if err := e.writeIndentationCharacters(); err != nil {
				return err
			}
			// TODO: Allow name customization via struct tags
			if err := e.writeString(CleanIdentifier(ft.Name) + ": "); err != nil {
				return err
			}
			if err := e.marshal(fv); err != nil {
				return fmt.Errorf("failed to encode value of struct field %q: %w", ft.Name, err)
			}
			if err := e.writeString(",\n"); err != nil {
				return err
			}
		}
	}

	e.indentLevel--

	if err := e.writeIndentationCharacters(); err != nil {
		return err
	}

	return e.writeString(")")
}

func (e *VariableEncoder) resolveKeyName(v reflect.Value) (string, error) {
	// From encoding/json/encode.go.
	if v.Kind() == reflect.String {
		return v.String(), nil
	}
	if tm, ok := v.Interface().(encoding.TextMarshaler); ok {
		if v.Kind() == reflect.Pointer && v.IsNil() {
			return "", nil
		}
		buf, err := tm.MarshalText()
		return string(buf), err
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	}
	return "", fmt.Errorf("unsupported map key type %q", v.Type().String())
}

func (e *VariableEncoder) encodeMap(v reflect.Value) error {
	if v.Len() == 0 {
		return e.writeString("()")
	}

	if err := e.writeString("(\n"); err != nil {
		return err
	}

	e.indentLevel++

	// BUG: Map output needs to be sorted, otherwise this will cause the test to fail randomly

	mi := v.MapRange()
	for mi.Next() {
		mk, mv := mi.Key(), mi.Value()
		key, err := e.resolveKeyName(mk)
		if err != nil {
			return err
		}

		if err := e.writeIndentationCharacters(); err != nil {
			return err
		}
		if err := e.writeString(CleanIdentifier(key) + ": "); err != nil {
			return err
		}
		if err := e.marshal(mv); err != nil {
			return fmt.Errorf("failed to encode map field %q: %w", key, err)
		}

		if err := e.writeString(",\n"); err != nil {
			return err
		}
	}

	e.indentLevel--

	if err := e.writeIndentationCharacters(); err != nil {
		return err
	}

	return e.writeString(")")
}

func (e *VariableEncoder) EncodeByteSlice(bb []byte) error {
	if err := e.writeString("bytes(("); err != nil {
		return err
	}

	// TODO: Encode byte slice via base64 or similar and use a typst package to convert it into the corresponding bytes type

	for i, b := range bb {
		if i > 0 {
			if err := e.writeString(", "); err != nil {
				return err
			}
		}

		if err := e.writeString(strconv.FormatUint(uint64(b), 10)); err != nil {
			return err
		}
	}

	return e.writeString("))")
}

func (e *VariableEncoder) encodeSlice(v reflect.Value, t reflect.Type) error {

	// Special case for byte slices.
	if t.Elem().Kind() == reflect.Uint8 {
		return e.EncodeByteSlice(v.Bytes())
	}

	if err := e.writeString("("); err != nil {
		return err
	}

	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			if err := e.writeString(", "); err != nil {
				return err
			}
		}
		if err := e.marshal(v.Index(i)); err != nil {
			return fmt.Errorf("failed to encode slice element %d of %d: %w", i+1, n+1, err)
		}
	}

	return e.writeString(")")
}

func (e *VariableEncoder) encodeArray(v reflect.Value) error {
	if err := e.writeString("("); err != nil {
		return err
	}

	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			if err := e.writeString(", "); err != nil {
				return err
			}
		}
		if err := e.marshal(v.Index(i)); err != nil {
			return fmt.Errorf("failed to encode array element %d of %d: %w", i+1, n+1, err)
		}
	}

	return e.writeString(")")
}

func (e *VariableEncoder) encodeTime(t time.Time) error {
	return e.writeString(fmt.Sprintf("datetime(year: %d, month: %d, day: %d, hour: %d, minute: %d, second: %d)",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	))
}

func (e *VariableEncoder) encodeDuration(d time.Duration) error {
	return e.writeString(fmt.Sprintf("duration(seconds: %d)", int(math.Round(d.Seconds()))))
}
