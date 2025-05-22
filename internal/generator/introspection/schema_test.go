package introspection_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
				Name:    nil,
				NonNull: true,
				IsList:  false,
				Child: &introspection.FieldType{
					Name:    nil,
					NonNull: false,
					IsList:  true,
					Child: &introspection.FieldType{
						Name:    nil,
						NonNull: true,
						IsList:  false,
						Child: &introspection.FieldType{
							Name:    []byte("Post"),
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
				Name:    []byte("String"),
				NonNull: true,
				IsList:  false,
				Child:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := introspection.ExpandType(tt.input)

			if d := cmp.Diff(tt.expected, result); d != "" {
				t.Errorf("ExpandType() mismatch (-want +got):\n%s", d)
			}
		})
	}
}
