package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strconv"
)

type reference struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type fieldDescriptor struct {
	Name       string    `yaml:"name"`
	Mandatory  bool      `yaml:"mandatory"`
	Xsd        string    `yaml:"xsd"`
	Sqlite     string    `yaml:"sqlite"`
	References reference `yaml:"references"`
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

// Implements CopyFromSource
type xmlSource struct {
	root    string
	element string
	fields  *fields
	state   int
	d       *xml.Decoder
	values  []string
	err     error
}

func newXmlSource(td *tableDescriptor, d *xml.Decoder, fields *fields) xmlSource {
	return xmlSource{
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
func (s *xmlSource) Next() bool {
	values, err := s.nextValues()
	log.Printf("Next %v %v", values, err)
	if err == io.EOF {
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
func (s *xmlSource) Values() ([]interface{}, error) {
	values := make([]interface{}, len(s.values))
	//for i, v := range s.values {
	//	values[i] = v
	//}

	// TODO: Convert to the right type dynamically (with pgtype.Record?).
	values[0], _ = strconv.Atoi(s.values[0])
	values[1] = s.values[1]
	log.Printf("Values %v %v %v", s.values, values, s.err)
	return values, s.err
}

// Err() implements pgx.CopyFromSource.
func (s *xmlSource) Err() error {
	return s.err
}

func (s *xmlSource) nextValues() ([]string, error) {
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
				return fields.record(item), nil
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

func main() {
	descriptorFileName := flag.String("descriptor", "<undefined>", "file name of the table descriptor")
	databaseUrl := flag.String("database", "<undefined>", "PostgreSQL database URL")
	flag.Parse()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, *databaseUrl)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close(ctx)

	td, err := decodeDescriptor(*descriptorFileName)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Construct buffered reader and writer.
	const bufSize = 4096 * 1024
	br := bufio.NewReaderSize(os.Stdin, bufSize)
	bw := bufio.NewWriterSize(os.Stdout, bufSize)
	defer bw.Flush()

	//err = convertXml(td, xml.NewDecoder(br), csv.NewWriter(bw))
	//if err != nil {
	//	log.Fatalf("%v", err)
	//}

	fields := newFields(td.Fields)
	s := newXmlSource(td, xml.NewDecoder(br), fields)
	i, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{td.Element},
		fields.header(),
		&s,
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("wrote %d rows", i)
}
