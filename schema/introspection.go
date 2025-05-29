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
			Name: nil,
			Type: &FieldType{
				Name:     []byte("__TypeKind"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
			Values: []*EnumElement{
				{
					Name:  []byte("__TypeKind_SCALAR"),
					Value: []byte("SCALAR"),
				},
				{
					Name:  []byte("__TypeKind_OBJECT"),
					Value: []byte("OBJECT"),
				},
				{
					Name:  []byte("__TypeKind_INTERFACE"),
					Value: []byte("INTERFACE"),
				},
				{
					Name:  []byte("__TypeKind_UNION"),
					Value: []byte("UNION"),
				},
				{
					Name:  []byte("__TypeKind_ENUM"),
					Value: []byte("ENUM"),
				},
				{
					Name:  []byte("__TypeKind_INPUT_OBJECT"),
					Value: []byte("INPUT_OBJECT"),
				},
				{
					Name:  []byte("__TypeKind_LIST"),
					Value: []byte("LIST"),
				},
				{
					Name:  []byte("__TypeKind_NON_NULL"),
					Value: []byte("NON_NULL"),
				},
			},
			Directives: nil,
			Extentions: nil,
		},
	}
	directiveLocationIntrospectionFields = []*EnumDefinition{
		{
			Name: nil,
			Type: &FieldType{
				Name:     []byte("__DirectiveLocation"),
				IsList:   false,
				Nullable: false,
				ListType: nil,
			},
			Values: []*EnumElement{
				{
					Name:  []byte("__DirectiveLocation_QUERY"),
					Value: []byte("QUERY"),
				},
				{
					Name:  []byte("__DirectiveLocation_MUTATION"),
					Value: []byte("MUTATION"),
				},
				{
					Name:  []byte("__DirectiveLocation_SUBSCRIPTION"),
					Value: []byte("SUBSCRIPTION"),
				},
				{
					Name:  []byte("__DirectiveLocation_FIELD"),
					Value: []byte("FIELD"),
				},
				{
					Name:  []byte("__DirectiveLocation_FRAGMENT_DEFINITION"),
					Value: []byte("FRAGMENT_DEFINITION"),
				},
				{
					Name:  []byte("__DirectiveLocation_FRAGMENT_SPREAD"),
					Value: []byte("FRAGMENT_SPREAD"),
				},
				{
					Name:  []byte("__DirectiveLocation_INLINE_FRAGMENT"),
					Value: []byte("INLINE_FRAGMENT"),
				},
				{
					Name:  []byte("__DirectiveLocation_VARIABLE_DEFINITION"),
					Value: []byte("VARIABLE_DEFINITION"),
				},
				{
					Name:  []byte("__DirectiveLocation_SCHEMA"),
					Value: []byte("SCHEMA"),
				},
				{
					Name:  []byte("__DirectiveLocation_SCALAR"),
					Value: []byte("SCALAR"),
				},
				{
					Name:  []byte("__DirectiveLocation_OBJECT"),
					Value: []byte("OBJECT"),
				},
				{
					Name:  []byte("__DirectiveLocation_FIELD_DEFINITION"),
					Value: []byte("FIELD_DEFINITION"),
				},
				{
					Name:  []byte("__DirectiveLocation_ARGUMENT_DEFINITION"),
					Value: []byte("ARGUMENT_DEFINITION"),
				},
				{
					Name:  []byte("__DirectiveLocation_INTERFACE"),
					Value: []byte("INTERFACE"),
				},
				{
					Name:  []byte("__DirectiveLocation_UNION"),
					Value: []byte("UNION"),
				},
				{
					Name:  []byte("__DirectiveLocation_ENUM"),
					Value: []byte("ENUM"),
				},
				{
					Name:  []byte("__DirectiveLocation_ENUM_VALUE"),
					Value: []byte("ENUM_VALUE"),
				},
				{
					Name:  []byte("__DirectiveLocation_INPUT_OBJECT"),
					Value: []byte("INPUT_OBJECT"),
				},
				{
					Name:  []byte("__DirectiveLocation_INPUT_FIELD_DEFINITION"),
					Value: []byte("INPUT_FIELD_DEFINITION"),
				},
			},
			Directives: nil,
			Extentions: nil,
		},
	}

	fieldsIntrospectionOperationDefinition = &OperationDefinition{
		Name: []byte("__fields"),
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
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:              []byte("__DirectiveLocation"),
		Fields:            nil,
		Required:          make(map[*FieldDefinition]struct{}),
		tokens:            schema.tokens,
		PrimitiveTypeName: []byte("string"),
		Interfaces:        nil,
		Directives:        nil,
		Extentions:        nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:              []byte("__TypeKind"),
		Fields:            nil,
		Required:          make(map[*FieldDefinition]struct{}),
		tokens:            schema.tokens,
		PrimitiveTypeName: []byte("string"),
		Interfaces:        nil,
		Directives:        nil,
		Extentions:        nil,
	})

	schema.Enums = append(schema.Enums, typeKindIntrospectionFields...)
	schema.Enums = append(schema.Enums, directiveLocationIntrospectionFields...)

	return schema
}
