package internal

import (
	"encoding/xml"
	"errors"
	"io"
	"marktstammdatenregister.dev/internal/spec"
)

type XMLSource struct {
	reader XMLReader
	fields *Fields
	skip   func(item map[string]string) bool
	values []interface{}
	err    error
}

func NewXMLSource(td *spec.Table, d *xml.Decoder, fields *Fields, skip func(map[string]string) bool) XMLSource {
	return XMLSource{
		reader: NewXMLReader(td, d),
		fields: fields,
		skip:   skip,
		values: nil,
		err:    nil,
	}
}

// Next() implements pgx.CopyFromSource.
func (s *XMLSource) Next() bool {
	for {
		values, err := s.reader.Read()
		if errors.Is(err, io.EOF) {
			return false
		}
		if err != nil {
			s.err = err
			return false
		}
		s.values, err = s.fields.Record(values)
		if err != nil {
			s.err = err
			return false
		}
		if s.skip(values) {
			continue
		}
		return true
	}
}

// Values() implements pgx.CopyFromSource.
func (s *XMLSource) Values() ([]interface{}, error) {
	return s.values, s.err
}

// Err() implements pgx.CopyFromSource.
func (s *XMLSource) Err() error {
	return s.err
}
