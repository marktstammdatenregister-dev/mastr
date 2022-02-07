package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"pvdb.de/mastr/internal/spec"
)

type tableMeta struct {
	Name   string `yaml:"name"`
	Schema string `yaml:"schema"`
}

type table struct {
	Table               tableMeta   `yaml:"table"`
	ObjectRelationships []objectRel `yaml:"object_relationships"`
	ArrayRelationships  []arrayRel  `yaml:"array_relationships"`
}

type objectRelFKC struct {
	ForeignKeyConstraintOn string `yaml:"foreign_key_constraint_on"`
}

type objectRel struct {
	Name  string       `yaml:"name"`
	Using objectRelFKC `yaml:"using"`
}

type arrayRelFKC struct {
	ForeignKeyConstraintOn struct {
		Column string    `yaml:"column"`
		Table  tableMeta `yaml:"table"`
	} `yaml:"foreign_key_constraint_on"`
}

type arrayRel struct {
	Name  string      `yaml:"name"`
	Using arrayRelFKC `yaml:"using"`
}

func ToHasura(schema string, td spec.Table) ([]byte, error) {
	var objectRels []objectRel
	for _, field := range td.Fields {
		if field.References == nil {
			continue
		}
		ref := field.References
		objectRels = append(objectRels, objectRel{
			Name: fmt.Sprintf("%s_%s", field.Name, ref.Table),
			Using: objectRelFKC{
				ForeignKeyConstraintOn: field.Name,
			},
		})
	}
	return yaml.Marshal(table{
		Table: tableMeta{
			Name:   td.Element,
			Schema: schema,
		},
		ObjectRelationships: objectRels,
		ArrayRelationships:  []arrayRel{},
	})
}
