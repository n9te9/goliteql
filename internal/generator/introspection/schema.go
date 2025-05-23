package introspection

import "github.com/n9te9/goliteql/schema"

type FieldType struct {
	Name    []byte
	NonNull bool
	IsList  bool
	Child   *FieldType
}

func (f *FieldType) IsPrimitive() bool {
	name := string(f.Name)
	return name == "String" || name == "Int" || name == "Float" || name == "Boolean" || name == "ID"
}

func (f *FieldType) IsObject() bool {
	name := string(f.Name)
	return name != ""
}

func ExpandType(fieldType *schema.FieldType) *FieldType {
	return expandType(fieldType, false)
}

func expandType(fieldType *schema.FieldType, isExpandedNonNull bool) *FieldType {
	if !fieldType.Nullable && !isExpandedNonNull {
		return &FieldType{
			Name:    nil,
			NonNull: true,
			IsList:  false,
			Child:   expandType(fieldType, true),
		}
	}

	if fieldType.IsList {
		return &FieldType{
			Name:    nil,
			NonNull: fieldType.Nullable,
			IsList:  true,
			Child:   expandType(fieldType.ListType, false),
		}
	}

	return &FieldType{
		Name:    fieldType.Name,
		NonNull: fieldType.Nullable,
		IsList:  false,
		Child:   nil,
	}
}

func (f *FieldType) Unwrap() *FieldType {
	return f.Child
}
