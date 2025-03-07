// Copyright (c) 2024-2025 David Vogel
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package typst

import (
	"bytes"
	"cmp"
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"slices"
	"strconv"
	"time"
)

// MarshalValue takes any Go type and returns a Typst markup representation as a byte slice.
func MarshalValue(v any) ([]byte, error) {
	var buf bytes.Buffer

	enc := NewValueEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ValueMarshaler can be implemented by types to support custom Typst marshaling.
type ValueMarshaler interface {
	MarshalTypstValue() ([]byte, error)
}

type ValueEncoder struct {
	indentLevel int

	writer io.Writer
}

// NewValueEncoder returns a new encoder that writes into w.
func NewValueEncoder(w io.Writer) *ValueEncoder {
	return &ValueEncoder{
		writer: w,
	}
}

func (e *ValueEncoder) Encode(v any) error {
	return e.marshal(reflect.ValueOf(v))
}

func (e *ValueEncoder) writeString(s string) error {
	return e.writeBytes([]byte(s))
}

func (e *ValueEncoder) writeRune(r rune) error {
	return e.writeBytes([]byte{byte(r)})
}

func (e *ValueEncoder) writeStringLiteral(s []byte) error {
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

func (e *ValueEncoder) writeBytes(b []byte) error {
	if _, err := e.writer.Write(b); err != nil {
		return fmt.Errorf("failed to write into writer: %w", err)
	}

	return nil
}

func (e *ValueEncoder) writeIndentationCharacters() error {
	return e.writeBytes(slices.Repeat([]byte{' ', ' '}, e.indentLevel))
}

func (e *ValueEncoder) marshal(v reflect.Value) error {
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
			return e.writeString("none")
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
			return e.writeString("none")
		}
		if err := e.encodeDuration(*i); err != nil {
			return err
		}
		return nil
	}

	if t.Implements(reflect.TypeFor[ValueMarshaler]()) {
		if m, ok := v.Interface().(ValueMarshaler); ok {
			bytes, err := m.MarshalTypstValue()
			if err != nil {
				return fmt.Errorf("error calling MarshalTypstValue for type %s: %w", t.String(), err)
			}
			return e.writeBytes(bytes)
		}
		return e.writeString("none")
	}

	// TODO: Remove this in a future update, it's only here for compatibility reasons
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
		if v.Int() >= 0 {
			err = e.writeString(strconv.FormatInt(v.Int(), 10))
		} else {
			if err = e.writeRune('{'); err != nil {
				break
			}
			if err = e.writeString(strconv.FormatInt(v.Int(), 10)); err != nil {
				break
			}
			if err = e.writeRune('}'); err != nil {
				break
			}
		}
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
			err = e.writeString("{-float.inf}")
		case math.Signbit(f):
			if err = e.writeRune('{'); err != nil {
				break
			}
			if err = e.writeString(strconv.FormatFloat(f, 'e', -1, 64)); err != nil {
				break
			}
			if err = e.writeRune('}'); err != nil {
				break
			}
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

func (e *ValueEncoder) encodeString(v reflect.Value) error {
	return e.writeStringLiteral([]byte(v.String()))
}

func (e *ValueEncoder) encodeStruct(v reflect.Value, t reflect.Type) error {
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
			fieldName := ft.Name
			if name, ok := ft.Tag.Lookup("typst"); ok {
				// Omit fields that have their name set to "-".
				if name == "-" {
					continue
				}
				fieldName = name
			}

			if err := e.writeIndentationCharacters(); err != nil {
				return err
			}
			if err := e.writeStringLiteral([]byte(fieldName)); err != nil {
				return err
			}
			if err := e.writeString(": "); err != nil {
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

	return e.writeRune(')')
}

func (e *ValueEncoder) resolveKeyName(v reflect.Value) (string, error) {
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

func (e *ValueEncoder) encodeMap(v reflect.Value) error {
	if v.Len() == 0 {
		return e.writeString("()")
	}

	if err := e.writeString("(\n"); err != nil {
		return err
	}

	e.indentLevel++

	type pair struct {
		key   string
		value reflect.Value
	}

	// Get all key value pairs as reflect.Value.
	mi := v.MapRange()
	pairs := make([]pair, 0, v.Len())
	for mi.Next() {
		mk, mv := mi.Key(), mi.Value()
		key, err := e.resolveKeyName(mk)
		if err != nil {
			return err
		}
		pairs = append(pairs, pair{key, mv})
	}

	// Sort and then generate markup.
	slices.SortFunc(pairs, func(a, b pair) int { return cmp.Compare(a.key, b.key) })
	for _, pair := range pairs {
		key, value := pair.key, pair.value

		if err := e.writeIndentationCharacters(); err != nil {
			return err
		}
		if err := e.writeStringLiteral([]byte(key)); err != nil {
			return err
		}
		if err := e.writeString(": "); err != nil {
			return err
		}
		if err := e.marshal(value); err != nil {
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

	return e.writeRune(')')
}

func (e *ValueEncoder) EncodeByteSlice(bb []byte) error {
	if err := e.writeString("bytes(("); err != nil {
		return err
	}

	// TODO: Encode byte slice via base64 or similar and use a Typst package to convert it into the corresponding bytes type

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

	if len(bb) == 1 {
		if err := e.writeRune(','); err != nil {
			return err
		}
	}

	return e.writeString("))")
}

func (e *ValueEncoder) encodeSlice(v reflect.Value, t reflect.Type) error {

	// Special case for byte slices.
	if t.Elem().Kind() == reflect.Uint8 {
		return e.EncodeByteSlice(v.Bytes())
	}

	if err := e.writeRune('('); err != nil {
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

	if n == 1 {
		if err := e.writeRune(','); err != nil {
			return err
		}
	}

	return e.writeRune(')')
}

func (e *ValueEncoder) encodeArray(v reflect.Value) error {
	if err := e.writeRune('('); err != nil {
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

	if n == 1 {
		if err := e.writeRune(','); err != nil {
			return err
		}
	}

	return e.writeRune(')')
}

func (e *ValueEncoder) encodeTime(t time.Time) error {
	return e.writeString(fmt.Sprintf("datetime(year: %d, month: %d, day: %d, hour: %d, minute: %d, second: %d)",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
	))
}

func (e *ValueEncoder) encodeDuration(d time.Duration) error {
	return e.writeString(fmt.Sprintf("duration(seconds: %d)", int(math.Round(d.Seconds()))))
}
