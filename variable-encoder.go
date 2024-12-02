package typst

import (
	"encoding"
	"fmt"
	"io"
	"math"
	"reflect"
	"slices"
	"strconv"
	"unicode/utf8"
)

type VariableMarshaler interface {
	MarshalTypstVariable() ([]byte, error)
}

type VariableEncoder struct {
	indentLevel int

	writer io.Writer
}

func NewVariableEncoder(w io.Writer) *VariableEncoder {
	return &VariableEncoder{
		writer: w,
	}
}

func (e *VariableEncoder) Encode(v any) error {
	return e.marshal(reflect.ValueOf(v))
}

func (e *VariableEncoder) WriteString(s string) {
	e.WriteBytes([]byte(s))
}

func (e *VariableEncoder) WriteBytes(b []byte) {
	e.writer.Write(b)
}

func (e *VariableEncoder) WriteIndentationCharacters() {
	e.WriteBytes(slices.Repeat([]byte{' ', ' '}, e.indentLevel))
}

func (e *VariableEncoder) marshal(v reflect.Value) error {
	if !v.IsValid() {
		e.WriteString("none")
		return nil
		//return fmt.Errorf("invalid reflect.Value %v", v)
	}

	t := v.Type()

	if (t.Kind() == reflect.Pointer || t.Kind() == reflect.Interface) && v.IsNil() {
		e.WriteString("none")
		return nil
	}

	if t.Implements(reflect.TypeFor[VariableMarshaler]()) {
		if m, ok := v.Interface().(VariableMarshaler); ok {
			bytes, err := m.MarshalTypstVariable()
			e.WriteBytes(bytes)
			return err
		}
		e.WriteString("none")
		return nil
	}

	if t.Implements(reflect.TypeFor[encoding.TextMarshaler]()) {
		if m, ok := v.Interface().(encoding.TextMarshaler); ok {
			bytes, err := m.MarshalText()
			e.WriteBytes(bytes)
			return err
		}
		e.WriteString("none")
		return nil
	}

	// TODO: Handle images

	// TODO: Handle decimals

	// TODO: Handle Time

	// TODO: Handle durations

	switch t.Kind() {
	case reflect.Bool:
		e.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		e.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		e.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		switch {
		case math.IsNaN(f):
			e.WriteString("float.nan")
		case math.IsInf(f, 1):
			e.WriteString("float.inf")
		case math.IsInf(f, -1):
			e.WriteString("-float.inf")
		default:
			e.WriteString(strconv.FormatFloat(f, 'e', -1, 64))
		}
	case reflect.String:
		return e.encodeString(v)
	case reflect.Interface, reflect.Pointer:
		return e.marshal(v.Elem())
	case reflect.Map:
		return e.encodeMap(v)
	case reflect.Struct:
		return e.encodeStruct(v, t)
	case reflect.Slice:
		return e.encodeSlice(v)
	case reflect.Array:
		return e.encodeArray(v)
	default:
		return fmt.Errorf("unsupported type %q", t.String())
	}

	return nil
}

func (e *VariableEncoder) encodeString(v reflect.Value) error {

	src := v.String()

	dst := make([]byte, 0, len(src)+2)

	dst = append(dst, '"')

	for _, r := range src {
		switch r {
		case '\\', '"':
			dst = append(dst, '\\')
			dst = utf8.AppendRune(dst, r)
		case '\n':
			dst = append(dst, '\\', 'n')
		case '\r':
			dst = append(dst, '\\', 'r')
		case '\t':
			dst = append(dst, '\\', 't')
		}
		dst = utf8.AppendRune(dst, r)
	}

	dst = append(dst, '"')

	e.WriteBytes(dst)

	return nil
}

func (e *VariableEncoder) encodeStruct(v reflect.Value, t reflect.Type) error {
	if v.NumField() == 0 {
		e.WriteString("()")
		return nil
	}

	e.WriteString("(\n")

	e.indentLevel++

	for i := 0; i < t.NumField(); i++ {
		ft, fv := t.Field(i), v.Field(i)
		if ft.PkgPath == "" { // Ignore unexported fields.
			e.WriteIndentationCharacters()
			e.WriteString(ft.Name + ": ")
			if err := e.marshal(fv); err != nil {
				return fmt.Errorf("failed to encode value of struct field %q", ft.Name)
			}
			e.WriteString(",\n")
		}
	}

	e.indentLevel--

	e.WriteIndentationCharacters()
	e.WriteString(")")

	return nil
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
		e.WriteString("()")
		return nil
	}

	e.WriteString("(\n")

	e.indentLevel++

	mi := v.MapRange()
	for mi.Next() {
		mk, mv := mi.Key(), mi.Value()
		key, err := e.resolveKeyName(mk)
		if err != nil {
			return err
		}

		e.WriteIndentationCharacters()
		e.WriteString(key + ": ")
		if err := e.marshal(mv); err != nil {
			return fmt.Errorf("failed to encode map field %q", key)
		}

		e.WriteString(",\n")
	}

	e.indentLevel--

	e.WriteIndentationCharacters()
	e.WriteString(")")

	return nil
}

func (e *VariableEncoder) encodeSlice(v reflect.Value) error {
	e.WriteString("(")

	// TODO: Output byte slice as a base64 and use the typst based package to convert that into typst Bytes.

	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteString(", ")
		}
		if err := e.marshal(v.Index(i)); err != nil {
			return fmt.Errorf("failed to encode slice element %d of %d", i+1, n+1)
		}
	}

	e.WriteString(")")

	return nil
}

func (e *VariableEncoder) encodeArray(v reflect.Value) error {
	e.WriteString("(")

	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteString(", ")
		}
		if err := e.marshal(v.Index(i)); err != nil {
			return fmt.Errorf("failed to encode array element %d of %d", i+1, n+1)
		}
	}

	e.WriteString(")")

	return nil
}
