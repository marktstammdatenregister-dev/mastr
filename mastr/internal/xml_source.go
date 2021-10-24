package internal

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"pvdb.de/mastr/internal/spec"
)

const (
	startRoot = iota
	startItemOrEndRoot
	startFieldOrEndItem
	fieldValueOrEndField
	finished
)

// Implements CopyFromSource
type XMLSource struct {
	root    string
	element string
	fields  *Fields
	state   int
	d       *xml.Decoder
	values  []interface{}
	err     error
}

func NewXMLSource(td *spec.Table, d *xml.Decoder, fields *Fields) XMLSource {
	return XMLSource{
		root:    td.Root,
		element: td.Element,
		fields:  fields,
		state:   startRoot,
		d:       d,
		values:  nil,
		err:     nil,
	}
}

// Next() implements pgx.CopyFromSource.
func (s *XMLSource) Next() bool {
	values, err := s.nextValues()
	if errors.Is(err, io.EOF) {
		return false
	}
	if err != nil {
		s.err = err
		return false
	}
	s.values = values
	return true
}

// Values() implements pgx.CopyFromSource.
func (s *XMLSource) Values() ([]interface{}, error) {
	return s.values, s.err
}

// Err() implements pgx.CopyFromSource.
func (s *XMLSource) Err() error {
	return s.err
}

func (s *XMLSource) nextValues() ([]interface{}, error) {
	d := s.d
	root := s.root
	element := s.element
	fields := s.fields

	item := make(map[string]string)
	var fieldName string
	var fieldValue []byte

	for {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch s.state {
		case startRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != root {
					return nil, fmt.Errorf("[%d] expected start of %s, got %s", s.state, root, name)
				}
				s.state = startItemOrEndRoot
			default: // ignore
			}
		case startItemOrEndRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != element {
					return nil, fmt.Errorf("[%d] expected start of %s, got %s", s.state, element, name)
				}
				s.state = startFieldOrEndItem
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != root {
					return nil, fmt.Errorf("[%d] expected start of %s, got %s", s.state, root, name)
				}
				s.state = finished
			default: // ignore
			}
		case startFieldOrEndItem:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				fieldName = name
				s.state = fieldValueOrEndField
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != element {
					return nil, fmt.Errorf("[%d] expected end of %s, got %s", s.state, element, name)
				}
				s.state = startItemOrEndRoot
				return fields.Record(item)
			default: // ignore
			}
		case fieldValueOrEndField:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				return nil, fmt.Errorf("[%d] expected end of %s, got start of %s", s.state, fieldName, name)
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != fieldName {
					return nil, fmt.Errorf("[%d] expected end of %s, got %s", s.state, fieldName, name)
				}
				item[fieldName] = string(fieldValue)
				fieldValue = []byte{}
				s.state = startFieldOrEndItem
			case xml.CharData:
				fieldValue = append(fieldValue, []byte(xml.CharData(t))...)
			default: // ignore
			}
		case finished:
			switch tok.(type) {
			case xml.CharData: // ignore
			default:
				return nil, fmt.Errorf("[%d] parsing finished, but got %v", s.state, tok)
			}
		}
	}
}
