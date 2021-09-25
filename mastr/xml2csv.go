package main

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type fieldDescriptor struct {
	Name       string `json:"name"`
	Mandatory  bool   `json:"mandatory"`
	Xsd        string `json:"xsd"`
	Sqlite     string `json:"sqlite"`
	References string `json:"references"`
}

type tableDescriptor struct {
	Root    string            `json:"root"`
	Element string            `json:"element"`
	Primary string            `json:"primary"`
	Fields  []fieldDescriptor `json:"fields"`
}

type fields struct {
	fields map[string]uint
}

const (
	startRoot = iota
	startItemOrEndRoot
	startFieldOrEndItem
	fieldValueOrEndField
	finished
)

func decodeDescriptor(descriptorFileName string) (*tableDescriptor, error) {
	var tableDescriptor tableDescriptor
	f, err := os.Open(descriptorFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	err = d.Decode(&tableDescriptor)
	if err != nil {
		return nil, err
	}
	return &tableDescriptor, nil
}

func newFields(fieldDescriptors []fieldDescriptor) *fields {
	f := make(map[string]uint)
	for i, fieldDescriptor := range fieldDescriptors {
		f[fieldDescriptor.Name] = uint(i)
	}
	return &fields{fields: f}
}

func (f *fields) header() []string {
	n := len(f.fields)
	result := make([]string, n, n)
	for name, i := range f.fields {
		result[i] = name
	}
	return result
}

func (f *fields) record(item map[string]string) []string {
	n := len(f.fields)
	result := make([]string, n, n)
	for name, value := range item {
		result[f.fields[name]] = value
	}
	return result
}

func convertXml(td *tableDescriptor, d *xml.Decoder, w *csv.Writer) error {
	root := td.Root
	element := td.Element
	fields := newFields(td.Fields)

	state := startRoot
	item := make(map[string]string)
	var fieldName string
	var fieldValue []byte

	w.Write(fields.header())

	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch state {
		case startRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != root {
					return fmt.Errorf("[%d] expected start of %s, got %s", state, root, name)
				}
				state = startItemOrEndRoot
			default: // ignore
			}
		case startItemOrEndRoot:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				if name != element {
					return fmt.Errorf("[%d] expected start of %s, got %s", state, element, name)
				}
				state = startFieldOrEndItem
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != root {
					return fmt.Errorf("[%d] expected start of %s, got %s", state, root, name)
				}
				state = finished
			default: // ignore
			}
		case startFieldOrEndItem:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				fieldName = name
				state = fieldValueOrEndField
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != element {
					return fmt.Errorf("[%d] expected end of %s, got %s", state, element, name)
				}
				w.Write(fields.record(item))
				item = make(map[string]string)
				state = startItemOrEndRoot
			default: // ignore
			}
		case fieldValueOrEndField:
			switch t := tok.(type) {
			case xml.StartElement:
				name := xml.StartElement(t).Name.Local
				return fmt.Errorf("[%d] expected end of %s, got start of %s", state, fieldName, name)
			case xml.EndElement:
				name := xml.EndElement(t).Name.Local
				if name != fieldName {
					return fmt.Errorf("[%d] expected end of %s, got %s", state, fieldName, name)
				}
				item[fieldName] = string(fieldValue)
				fieldValue = []byte{}
				state = startFieldOrEndItem
			case xml.CharData:
				fieldValue = append(fieldValue, []byte(xml.CharData(t))...)
			default: // ignore
			}
		case finished:
			switch tok.(type) {
			case xml.CharData: // ignore
			default:
				return fmt.Errorf("[%d] parsing finished, but got %v", state, tok)
			}
		}
	}
	w.Flush()
	return w.Error()
}

func main() {
	descriptorFileName := flag.String("descriptor", "<undefined>", "file name of the table descriptor")
	flag.Parse()

	td, err := decodeDescriptor(*descriptorFileName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = convertXml(td, xml.NewDecoder(os.Stdin), csv.NewWriter(os.Stdout))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
