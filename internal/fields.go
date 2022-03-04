package internal

import (
	"fmt"
	"strconv"

	"marktstammdatenregister.dev/internal/spec"
)

var (
	unknownXsdType = "don't know how to handle XSD type '%s'"
	xsd2sqliteType = map[string]string{
		"nonNegativeInteger": "integer",
		"boolean":            "integer",
		"decimal":            "real",
		"date":               "text",
		"dateTime":           "text",
		"":                   "text",
	}
)

func Xsd2SqliteType(xsd string) (string, bool) {
	typ, ok := xsd2sqliteType[xsd]
	return typ, ok
}

type Fields struct {
	fields map[string]uint
	typ    map[string]string
}

func NewFields(fields []spec.Field) (*Fields, error) {
	f := make(map[string]uint)
	t := make(map[string]string)
	for i, field := range fields {
		typ, ok := Xsd2SqliteType(field.Xsd)
		if !ok {
			return nil, fmt.Errorf(unknownXsdType, field.Xsd)
		}
		f[field.Name] = uint(i)
		t[field.Name] = typ
	}
	return &Fields{fields: f, typ: t}, nil
}

func (f *Fields) Header() []string {
	n := len(f.fields)
	result := make([]string, n, n)
	for name, i := range f.fields {
		result[i] = name
	}
	return result
}

func (f *Fields) Record(item map[string]string) ([]interface{}, error) {
	n := len(f.fields)
	result := make([]interface{}, n, n)
	for name, value := range item {
		switch f.typ[name] {
		case "integer":
			if value == "" {
				result[f.fields[name]] = nil
				continue
			}
			v, err := strconv.Atoi(value)
			if err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "real":
			if value == "" {
				result[f.fields[name]] = nil
				continue
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "text", "":
			result[f.fields[name]] = value
		default:
			return nil, fmt.Errorf(unknownXsdType, f.typ[name])
		}
	}
	return result, nil
}
