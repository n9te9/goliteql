package introspection_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/n9te9/goliteql/internal/generator/introspection"
	"github.com/n9te9/goliteql/schema"
)

func TestExpandType(t *testing.T) {
	tests := []struct {
		name     string
		input    *schema.FieldType
		expected *introspection.FieldType
	}{
		{
			name: "non-nullable scalar type",
			input: &schema.FieldType{
				Name:     nil,
				Nullable: false,
				IsList:   true,
				ListType: &schema.FieldType{
					Name:     []byte("Post"),
					Nullable: false,
					IsList:   false,
				},
			},
			expected: &introspection.FieldType{
				Name:    "",
				NonNull: true,
				IsList:  false,
				Child: &introspection.FieldType{
					Name:    "",
					NonNull: false,
					IsList:  true,
					Child: &introspection.FieldType{
						Name:    "",
						NonNull: true,
						IsList:  false,
						Child: &introspection.FieldType{
							Name:    "Post",
							NonNull: false,
							IsList:  false,
							Child:   nil,
						},
					},
				},
			},
		},
		{
			name: "nullable scalar type",
			input: &schema.FieldType{
				Name:     []byte("String"),
				Nullable: true,
				IsList:   false,
			},
			expected: &introspection.FieldType{
				Name:    "String",
				NonNull: false,
				IsList:  false,
				Child:   nil,
			},
		},
		{
			name: "non-null scalar type",
			input: &schema.FieldType{
				Name:     []byte("ID"),
				Nullable: false,
				IsList:   false,
			},
			expected: &introspection.FieldType{
				Name:    "",
				NonNull: true,
				IsList:  false,
				Child: &introspection.FieldType{
					Name:    "ID",
					NonNull: false,
					IsList:  false,
					Child:   nil,
					SchemaFieldType: &schema.FieldType{
						Name:     []byte("ID"),
						Nullable: false,
						IsList:   false,
					},
				},
				SchemaFieldType: &schema.FieldType{
					Name:     []byte("ID"),
					Nullable: false,
					IsList:   false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := introspection.ExpandType(tt.input)

			if d := cmp.Diff(tt.expected, result, cmpopts.IgnoreFields(introspection.FieldType{}, "SchemaFieldType")); d != "" {
				t.Errorf("ExpandType() mismatch (-want +got):\n%s", d)
			}
		})
	}
}
