package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-parser/schema"
)


func generateModelField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldType := GraphQLType(f.Type.Name)
		var fieldTypeIdent *ast.Ident
		if fieldType.IsPrimitive() {
			fieldTypeIdent = golangType(f.Type, fieldType, "")
		} else {
			fieldTypeIdent = golangType(f.Type, fieldType, "")
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
			Tag: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}