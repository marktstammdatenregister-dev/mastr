package spec

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Reference struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type Field struct {
	Name       string     `yaml:"name"`
	Mandatory  bool       `yaml:"mandatory"`
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

func Decode(fileName string) (*Table, error) {
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
