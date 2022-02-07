package internal

import (
	"encoding/xml"
	"fmt"
	"pvdb.de/mastr/internal/spec"
)

const (
	startRoot = iota
	startItemOrEndRoot
	startFieldOrEndItem
	fieldValueOrEndField
	finished
)

type XMLReader struct {
	root    string
	element string
	state   int
	d       *xml.Decoder
}

func NewXMLReader(td *spec.Table, d *xml.Decoder) XMLReader {
	return XMLReader{
		root:    td.Root,
		element: td.Element,
		state:   startRoot,
		d:       d,
	}
}

func (s *XMLReader) Read() (map[string]string, error) {
	d := s.d
	root := s.root
	element := s.element

	item := make(map[string]string)
	var fieldName string
	var fieldValue []byte

	for {
		tok, err := d.Token()
		if err != nil {
			return item, err
		}
		switch s.state {
		case startRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != root {
					return item, fmt.Errorf("[%d] expected start of %s, got %s", s.state, root, name)
				}
				s.state = startItemOrEndRoot
			default: // ignore
			}
		case startItemOrEndRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != element {
					return item, fmt.Errorf("[%d] expected start of %s, got %s", s.state, element, name)
				}
				s.state = startFieldOrEndItem
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != root {
					return item, fmt.Errorf("[%d] expected start of %s, got %s", s.state, root, name)
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
					return item, fmt.Errorf("[%d] expected end of %s, got %s", s.state, element, name)
				}
				s.state = startItemOrEndRoot
				return item, nil
			default: // ignore
			}
		case fieldValueOrEndField:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				return item, fmt.Errorf("[%d] expected end of %s, got start of %s", s.state, fieldName, name)
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != fieldName {
					return item, fmt.Errorf("[%d] expected end of %s, got %s", s.state, fieldName, name)
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
				return item, fmt.Errorf("[%d] parsing finished, but got %v", s.state, tok)
			}
		}
	}
}
