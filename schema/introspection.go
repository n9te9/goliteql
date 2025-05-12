package schema

var (
	schemaIntrospectionFields = []*FieldDefinition{
		{
			Name:      []byte("description"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name:      []byte("queryType"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name:      []byte("mutationType"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name:      []byte("subscriptionType"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name:      []byte("types"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   true,
				Nullable: false,
				ListType: &FieldType{
					Name:     []byte("__Type"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name:      []byte("directives"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Directive"),
				IsList:   true,
				Nullable: false,
				ListType: &FieldType{
					Name:     []byte("__Directive"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
	}
	typeIntrospectionFields = []*FieldDefinition{
		{
			Name:      []byte("kind"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__TypeKind"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name:      []byte("name"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
		},
		{
			Name:      []byte("description"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
			Directives: nil,
			Default:    nil,
			Location:   nil,
		},
		{
			Name: []byte("fields"),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("includeDeprecated"),
					Type: &FieldType{
						Name:     []byte("Boolean"),
						IsList:   false,
						Nullable: true,
						ListType: nil,
					},
					Default: []byte("false"),
				},
			},
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: true,
				ListType: &FieldType{
					Name:     []byte("__Field"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name:      []byte("interfaces"),
			Arguments: nil,
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: true,
				ListType: &FieldType{
					Name:     []byte("__Type"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name:      []byte("possibleTypes"),
			Arguments: nil,
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: true,
				ListType: &FieldType{
					Name:     []byte("__Type"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name: []byte("enumValues"),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("includeDeprecated"),
					Type: &FieldType{
						Name:     []byte("Boolean"),
						IsList:   false,
						Nullable: true,
						ListType: nil,
					},
					Default: []byte("false"),
				},
			},
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: true,
				ListType: &FieldType{
					Name:     []byte("__EnumValue"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name: []byte("inputFields"),
			Arguments: []*ArgumentDefinition{
				{
					Name: []byte("includeDeprecated"),
					Type: &FieldType{
						Name:     []byte("Boolean"),
						IsList:   false,
						Nullable: true,
						ListType: nil,
					},
				},
			},
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: true,
				ListType: &FieldType{
					Name:     []byte("__InputValue"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name:      []byte("ofType"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
		},
	}
	fieldIntrospectionFields = []*FieldDefinition{
		{
			Name:      []byte("name"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
		},
		{
			Name:      []byte("description"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
		},
		{
			Name:      []byte("args"),
			Arguments: nil,
			Type: &FieldType{
				Name:     nil,
				IsList:   true,
				Nullable: false,
				ListType: &FieldType{
					Name:     []byte("__InputValue"),
					IsList:   false,
					Nullable: false,
					ListType: nil,
				},
			},
		},
		{
			Name:      []byte("type"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
		},
		{
			Name:      []byte("isDeprecated"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("Boolean"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
		},
		{
			Name:      []byte("deprecationReason"),
			Arguments: nil,
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
				ListType: nil,
			},
		},
	}
	inputValueIntrospectionFields = []*FieldDefinition{
		{
			Name: []byte("name"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("description"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
		{
			Name: []byte("type"),
			Type: &FieldType{
				Name:     []byte("__Type"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("defaultValue"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
		{
			Name: []byte("isDeprecated"),
			Type: &FieldType{
				Name:     []byte("Boolean"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("deprecationReason"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
	}

	enumValueIntrospectionFields = []*FieldDefinition{
		{
			Name: []byte("name"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("description"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
		{
			Name: []byte("isDeprecated"),
			Type: &FieldType{
				Name:     []byte("Boolean"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("deprecationReason"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
	}

	directiveIntrospectionFields = []*FieldDefinition{
		{
			Name: []byte("name"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: false,
			},
		},
		{
			Name: []byte("description"),
			Type: &FieldType{
				Name:     []byte("String"),
				IsList:   false,
				Nullable: true,
			},
		},
		{
			Name: []byte("locations"),
			Type: &FieldType{
				IsList:   true,
				Nullable: false,
				ListType: &FieldType{
					Name:     []byte("__DirectiveLocation"),
					IsList:   false,
					Nullable: false,
				},
			},
		},
		{
			Name: []byte("args"),
			Type: &FieldType{
				IsList:   true,
				Nullable: false,
				ListType: &FieldType{
					Name:     []byte("__InputValue"),
					IsList:   false,
					Nullable: false,
				},
			},
		},
		{
			Name: []byte("isRepeatable"),
			Type: &FieldType{
				Name:     []byte("Boolean"),
				IsList:   false,
				Nullable: false,
			},
		},
	}

	typeKindIntrospectionFields = []*EnumDefinition{
		{
			Name:       []byte("SCALAR"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("OBJECT"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INTERFACE"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("UNION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("ENUM"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INPUT_OBJECT"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("LIST"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("NON_NULL"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
	}
	directiveLocationIntrospectionFields = []*EnumDefinition{
		{
			Name:       []byte("QUERY"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("MUTATION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("SUBSCRIPTION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("FIELD"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("FRAGMENT_DEFINITION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("FRAGMENT_SPREAD"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INLINE_FRAGMENT"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("VARIABLE_DEFINITION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("SCHEMA"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("SCALAR"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("OBJECT"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("FIELD_DEFINITION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("ARGUMENT_DEFINITION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INTERFACE"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("UNION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("ENUM"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("ENUM_VALUE"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INPUT_OBJECT"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
		{
			Name:       []byte("INPUT_FIELD_DEFINITION"),
			Values:     nil,
			Directives: nil,
			Extentions: nil,
		},
	}
)

func withTypeIntrospection(schema *Schema) *Schema {
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__Schema"),
		Fields:     schemaIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__Type"),
		Fields:     typeIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__Field"),
		Fields:     fieldIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__InputValue"),
		Fields:     inputValueIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__EnumValue"),
		Fields:     enumValueIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__Directive"),
		Fields:     directiveIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Enums = append(schema.Enums, typeKindIntrospectionFields...)
	schema.Enums = append(schema.Enums, directiveLocationIntrospectionFields...)

	return schema
}
