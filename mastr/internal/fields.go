package internal

import (
	"fmt"
	"github.com/jackc/pgtype"
	"time"
)

var Location = time.UTC

type Fields struct {
	fields map[string]uint
	psqlty map[string]string
}

func NewFields(fieldDescriptors []FieldDescriptor) *Fields {
	f := make(map[string]uint)
	t := make(map[string]string)
	for i, fieldDescriptor := range fieldDescriptors {
		f[fieldDescriptor.Name] = uint(i)
		t[fieldDescriptor.Name] = fieldDescriptor.Psql
	}
	return &Fields{fields: f, psqlty: t}
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
		switch f.psqlty[name] {
		case "boolean":
			v := &pgtype.Bool{}
			if err := v.Set(value); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "date":
			ts, err := time.ParseInLocation("2006-01-02", value, Location)
			if err != nil {
				return result, err
			}
			v := &pgtype.Date{}
			if err := v.Set(ts); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "integer":
			v := &pgtype.Int4{}
			if err := v.Set(value); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "real":
			v := &pgtype.Float4{}
			if err := v.Set(value); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "text", "":
			v := &pgtype.Text{}
			if err := v.Set(value); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		case "timestamp":
			ts, err := time.ParseInLocation("2006-01-02T15:04:05.9999999", value, Location)
			if err != nil {
				return result, err
			}
			v := &pgtype.Timestamp{}
			if err := v.Set(ts); err != nil {
				return result, err
			}
			result[f.fields[name]] = v
		default:
			return nil, fmt.Errorf("unknown PostgreSQL type: %s", f.psqlty[name])
		}
	}
	return result, nil
}
