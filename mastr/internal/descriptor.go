package internal

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Reference struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type FieldDescriptor struct {
	Name       string     `yaml:"name"`
	Mandatory  bool       `yaml:"mandatory"`
	Xsd        string     `yaml:"xsd"`
	Sqlite     string     `yaml:"sqlite"`
	Psql       string     `yaml:"psql"`
	References *Reference `yaml:"references,omitempty"`
}

type TableDescriptor struct {
	Root    string            `yaml:"root"`
	Element string            `yaml:"element"`
	Primary string            `yaml:"primary"`
	Fields  []FieldDescriptor `yaml:"fields"`
}

func DecodeDescriptor(descriptorFileName string) (*TableDescriptor, error) {
	var tableDescriptor TableDescriptor
	f, err := os.Open(descriptorFileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()
	d := yaml.NewDecoder(f)
	err = d.Decode(&tableDescriptor)
	if err != nil {
		return nil, err
	}
	return &tableDescriptor, nil
}
