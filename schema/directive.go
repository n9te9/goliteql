package schema

import (
	"bytes"
	"fmt"

	"github.com/n9te9/goliteql/query"
)

type Location struct {
	Name []byte
}

type Directive struct {
	Name      []byte
	Arguments []*DirectiveArgument
	Locations []*Location
}

type Directives []*Directive

func (d *Directives) Get(name []byte) *Directive {
	for _, directive := range *d {
		if bytes.Equal(directive.Name, name) {
			return directive
		}
	}
	return nil
}

type DirectiveArgument struct {
	Name  []byte
	Value []byte
}

type DirectiveDefinition struct {
	Name        []byte
	Description []byte
	Arguments   []*ArgumentDefinition
	Repeatable  bool
	Locations   []*Location
}

func (d *DirectiveDefinition) IsAllowedApplySchema() bool {
	for _, l := range d.Locations {
		if bytes.Equal(l.Name, []byte("SCHEMA")) ||
			bytes.Equal(l.Name, []byte("SCALAR")) ||
			bytes.Equal(l.Name, []byte("OBJECT")) ||
			bytes.Equal(l.Name, []byte("FIELD_DEFINITION")) ||
			bytes.Equal(l.Name, []byte("ARGUMENT_DEFINITION")) ||
			bytes.Equal(l.Name, []byte("INTERFACE")) ||
			bytes.Equal(l.Name, []byte("UNION")) ||
			bytes.Equal(l.Name, []byte("ENUM")) ||
			bytes.Equal(l.Name, []byte("ENUM_VALUE")) ||
			bytes.Equal(l.Name, []byte("INPUT_OBJECT")) ||
			bytes.Equal(l.Name, []byte("INPUT_FIELD_DEFINITION")) {
			return true
		}
	}

	return false
}

func (d *DirectiveDefinition) IsAllowedApplyField() bool {
	for _, l := range d.Locations {
		if bytes.Equal(l.Name, []byte("FIELD")) ||
			bytes.Equal(l.Name, []byte("FRAGMENT_DEFINITION")) ||
			bytes.Equal(l.Name, []byte("FRAGMENT_SPREAD")) ||
			bytes.Equal(l.Name, []byte("INLINE_FRAGMENT")) ||
			bytes.Equal(l.Name, []byte("VARIABLE_DEFINITION")) {
			return true
		}
	}

	return false
}

func (d *DirectiveDefinition) ValidateArguments(args []*query.DirectiveArgument) error {
	for _, def := range d.Arguments {
		required := !def.Type.Nullable
		found := false
		for _, arg := range args {
			if bytes.Equal(def.Name, arg.Name) {
				if err := def.ValidateValueType(arg.Value); err != nil {
					return fmt.Errorf("error validating argument %s: %w", def.Name, err)
				}

				found = true
			}
		}

		if required && !found {
			return fmt.Errorf("missing required argument: %s", def.Name)
		}
	}

	return nil
}

type DirectiveDefinitions []*DirectiveDefinition

func (d DirectiveDefinitions) IsAllowedApplySchema(fieldName []byte) bool {
	for _, directive := range d {
		if bytes.Equal(directive.Name, fieldName) && directive.IsAllowedApplySchema() {
			return true
		}
	}

	return false
}

func (d DirectiveDefinitions) IsAllowedApplyField(fieldName []byte) bool {
	for _, directive := range d {
		if bytes.Equal(directive.Name, fieldName) && directive.IsAllowedApplyField() {
			return true
		}
	}

	return false
}

func (d DirectiveDefinitions) Get(name []byte) *DirectiveDefinition {
	for _, directive := range d {
		if bytes.Equal(directive.Name, name) {
			return directive
		}
	}

	return nil
}

func NewBuildInDirectives() []*DirectiveDefinition {
	return []*DirectiveDefinition{
		{
			Name:        []byte("skip"),
			Description: []byte("Directs the executor to skip this field or fragment when the `if` argument is true."),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("if"),
					Type: &FieldType{Name: []byte("Boolean"), Nullable: false},
				},
			},
			Repeatable: false,
			Locations: []*Location{
				{
					Name: []byte("FIELD"),
				},
				{
					Name: []byte("FRAGMENT_SPREAD"),
				},
				{
					Name: []byte("INLINE_FRAGMENT"),
				},
			},
		},
		{
			Name:        []byte("include"),
			Description: []byte("Directs the executor to include this field or fragment only when the `if` argument is true."),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("if"),
					Type: &FieldType{Name: []byte("Boolean"), Nullable: false},
				},
			},
			Repeatable: false,
			Locations: []*Location{
				{
					Name: []byte("FIELD"),
				},
				{
					Name: []byte("FRAGMENT_SPREAD"),
				},
				{
					Name: []byte("INLINE_FRAGMENT"),
				},
			},
		},
		{
			Name:        []byte("deprecated"),
			Description: []byte("Marks an element of a GraphQL schema as no longer supported."),
			Arguments: []*ArgumentDefinition{
				{
					Name:    []byte("reason"),
					Type:    &FieldType{Name: []byte("String"), Nullable: true},
					Default: []byte("No longer supported"),
				},
			},
			Repeatable: false,
			Locations: []*Location{
				{
					Name: []byte("FIELD_DEFINITION"),
				},
				{
					Name: []byte("ENUM_VALUE"),
				},
			},
		},
		{
			Name:        []byte("specifiedBy"),
			Description: []byte("Exposes a URL that specifies the behaviour of this scalar."),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("url"),
					Type: &FieldType{Name: []byte("String"), Nullable: false},
				},
			},
			Repeatable: false,
			Locations: []*Location{
				{
					Name: []byte("SCALAR"),
				},
			},
		},
	}
}
