package main

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type fieldDescriptor struct {
	Name       string `yaml:"name"`
	Mandatory  bool   `yaml:"mandatory"`
	Xsd        string `yaml:"xsd"`
	Sqlite     string `yaml:"sqlite"`
	References string `yaml:"references"`
}

type tableDescriptor struct {
	Root    string            `yaml:"root"`
	Element string            `yaml:"element"`
	Primary string            `yaml:"primary"`
	Fields  []fieldDescriptor `yaml:"fields"`
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
	d := yaml.NewDecoder(f)
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

	// NOTE(csv-header): It would be nicer to include the header to make the CSV files
	// self-contained. However, if we include the header here, we must skip it during the
	// import to SQLite. SQLite 3.32.0 added the --skip option to the .import command, but the
	// "ubuntu-latest" GitHub Actions runner uses Ubuntu 20.04, which ships with SQLite 3.31.
	// So overall it's just easier to not write the header here.
	//w.Write(fields.header())

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

	// Construct buffered reader and writer.
	const bufSize = 4096 * 1024
	br := bufio.NewReaderSize(os.Stdin, bufSize)
	bw := bufio.NewWriterSize(os.Stdout, bufSize)
	defer bw.Flush()

	err = convertXml(td, xml.NewDecoder(br), csv.NewWriter(bw))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
