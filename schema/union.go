package schema

import "bytes"

type UnionDefinition struct {
	Name []byte
	Types [][]byte
	Extentions []*UnionDefinition
	Directives []*Directive
}

func (u *UnionDefinition) GetFieldByName(name []byte) *FieldDefinition {
	for _, t := range u.Types {
		if bytes.Equal(t, name) {
			return &FieldDefinition{Name: name, Type: &FieldType{Name: name, Nullable: false}}
		}
	}

	return nil
}

func (u *UnionDefinition) HasType(name string) bool {
	for _, t := range u.Types {
		if string(t) == name {
			return true
		}
	}
	
	return false
}

type UnionDefinitions []*UnionDefinition

func (u UnionDefinitions) Has(name string) bool {
	for _, union := range u {
		if string(union.Name) == name {
			return true
		}
	}
	
	return false
}

func (u *UnionDefinition) TypeName() []byte {
	return u.Name
}