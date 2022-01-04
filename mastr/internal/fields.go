package internal

import (
	"fmt"
	"strconv"
	"time"

	"pvdb.de/mastr/internal/spec"
)

var Location = time.UTC

type Fields struct {
	fields   map[string]uint
	sqlitety map[string]string
}

func NewFields(fields []spec.Field) *Fields {
	f := make(map[string]uint)
	t := make(map[string]string)
	for i, field := range fields {
		f[field.Name] = uint(i)
		t[field.Name] = field.Sqlite
	}
	return &Fields{fields: f, sqlitety: t}
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
		switch f.sqlitety[name] {
		//case "boolean":
		//	v := &pgtype.Bool{}
		//	if err := v.Set(value); err != nil {
		//		return result, err
		//	}
		//	result[f.fields[name]] = v
		//case "date":
		//	ts, err := time.ParseInLocation("2006-01-02", value, Location)
		//	if err != nil {
		//		return result, err
		//	}
		//	v := &pgtype.Date{}
		//	if err := v.Set(ts); err != nil {
		//		return result, err
		//	}
		//	result[f.fields[name]] = v
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
		//case "timestamp":
		//	ts, err := time.ParseInLocation("2006-01-02T15:04:05.9999999", value, Location)
		//	if err != nil {
		//		return result, err
		//	}
		//	v := &pgtype.Timestamp{}
		//	if err := v.Set(ts); err != nil {
		//		return result, err
		//	}
		//	result[f.fields[name]] = v
		default:
			return nil, fmt.Errorf("unknown SQLite type: %s", f.sqlitety[name])
		}
	}
	return result, nil
}

func (f *Fields) Map(record []interface{}) (map[string]interface{}, error) {
	item := make(map[string]interface{})
	for field, i := range f.fields {
		if int(i) > len(record) - 1 {
			return item, fmt.Errorf("record has %d fields, expected %d", len(record), len(f.fields))
		}
		item[field] = record[i]
	}
	return item, nil
}

/*

When reading XML to SQLite: turn map[string]string into []interface{}
When serving SQLite as GraphQL: turn []interface{} into map[string]&graphql.Field ?

*/
