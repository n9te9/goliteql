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
	typeIntrospectionFields              = []*FieldDefinition{}
	fieldIntrospectionFields             = []*FieldDefinition{}
	inputValueIntrospectionFields        = []*FieldDefinition{}
	enumValueIntrospectionFields         = []*FieldDefinition{}
	directiveIntrospectionFields         = []*FieldDefinition{}
	typeKindIntrospectionFields          = []*FieldDefinition{}
	directiveLocationIntrospectionFields = []*FieldDefinition{}
	typeNameIntrospectionFields          = []*FieldDefinition{}
	inputTypeIntrospectionFields         = []*FieldDefinition{}
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
		Name:       []byte("__TypeKind"),
		Fields:     typeKindIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__DirectiveLocation"),
		Fields:     directiveLocationIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__TypeName"),
		Fields:     typeNameIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})
	schema.Types = append(schema.Types, &TypeDefinition{
		Name:       []byte("__InputType"),
		Fields:     inputTypeIntrospectionFields,
		Required:   make(map[*FieldDefinition]struct{}),
		tokens:     schema.tokens,
		Interfaces: nil,
		Directives: nil,
		Extentions: nil,
	})

	return schema
}
