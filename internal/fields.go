package internal

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"marktstammdatenregister.dev/internal/spec"
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
			return nil, fmt.Errorf("unknown SQLite type: %s", f.sqlitety[name])
		}
	}
	return result, nil
}

func (f *Fields) Map(record []interface{}) (map[string]interface{}, error) {
	item := make(map[string]interface{})
	for field, i := range f.fields {
		if int(i) > len(record)-1 {
			return item, fmt.Errorf("record has %d fields, expected %d", len(record), len(f.fields))
		}
		item[field] = record[i]
	}
	return item, nil
}

func (f *Fields) ScanDest() []interface{} {
	n := len(f.fields)
	dest := make([]interface{}, n)
	for name, i := range f.fields {
		switch f.sqlitety[name] {
		case "integer":
			dest[i] = &sql.NullInt64{}
		case "real":
			dest[i] = &sql.NullFloat64{}
		case "text":
			dest[i] = &sql.NullString{}
		default:
			dest[i] = &sql.NullString{}
		}
	}
	return dest
}

func (f *Fields) GraphqlType(name string) (graphql.Output, error) {
	sqlitety, ok := f.sqlitety[name]
	if !ok {
		return nil, fmt.Errorf("field %s not found", name)
	}
	return graphqlType(sqlitety), nil
}

func graphqlType(sqlitety string) graphql.Output {
	switch sqlitety {
	case "integer":
		return graphql.Int
	case "real":
		return graphql.Float
	case "text":
		return graphql.String
	default:
		return graphql.String
	}
}
