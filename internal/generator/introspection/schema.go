package introspection

import (
	"github.com/n9te9/goliteql/schema"
)

type FieldType struct {
	Name            string
	NonNull         bool
	IsList          bool
	SchemaFieldType *schema.FieldType
	Child           *FieldType
}

func (f *FieldType) IsPrimitive() bool {
	name := string(f.Name)
	return name == "String" || name == "Int" || name == "Float" || name == "Boolean" || name == "ID"
}

func (f *FieldType) IsObject() bool {
	name := string(f.Name)
	return name != ""
}

func (f *FieldType) IsObjectType() bool {
	return f.IsObject() && !f.IsPrimitive()
}

func ExpandType(fieldType *schema.FieldType) *FieldType {
	return expandType(fieldType, false)
}

func expandType(fieldType *schema.FieldType, isExpandedNonNull bool) *FieldType {
	if !fieldType.Nullable && !isExpandedNonNull {
		return &FieldType{
			Name:    "",
			NonNull: true,
			IsList:  false,
			Child:   expandType(fieldType, true),
		}
	}

	if fieldType.IsList {
		return &FieldType{
			Name:    "",
			NonNull: false,
			IsList:  true,
			Child:   expandType(fieldType.ListType, false),
		}
	}

	return &FieldType{
		Name:            string(fieldType.Name),
		NonNull:         false,
		IsList:          false,
		SchemaFieldType: fieldType,
		Child:           nil,
	}
}

func (f *FieldType) Unwrap() *FieldType {
	return f.Child
}
