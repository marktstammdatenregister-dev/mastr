package spec

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"
)

type Reference struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type Field struct {
	Name       string     `yaml:"name"`
	Index      bool       `yaml:"index"`
	Xsd        string     `yaml:"xsd"`
	Sqlite     string     `yaml:"sqlite"`
	Psql       string     `yaml:"psql"`
	References *Reference `yaml:"references,omitempty"`
}

type Table struct {
	Root    string  `yaml:"root"`
	Element string  `yaml:"element"`
	Primary string  `yaml:"primary"`
	Fields  []Field `yaml:"fields"`
}

type ExportDescriptor struct {
	Prefix string
	Table  Table
}

func DecodeTable(fileName string) (*Table, error) {
	var table Table
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()
	d := yaml.NewDecoder(f)
	err = d.Decode(&table)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func DecodeExport(fileName string) ([]ExportDescriptor, error) {
	var tableDescriptorFiles []string
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()
	d := yaml.NewDecoder(f)
	err = d.Decode(&tableDescriptorFiles)
	if err != nil {
		return nil, err
	}

	dir := path.Dir(fileName)
	var export []ExportDescriptor
	for _, descriptorFileName := range tableDescriptorFiles {
		if !strings.HasSuffix(descriptorFileName, ".yaml") {
			return nil, fmt.Errorf("missing yaml suffix: %s", descriptorFileName)
		}
		table, err := DecodeTable(path.Join(dir, descriptorFileName))
		if err != nil {
			return export, err
		}
		export = append(export, ExportDescriptor{
			Prefix: strings.TrimSuffix(descriptorFileName, ".yaml"),
			Table:  *table,
		})
	}
	return export, nil
}
