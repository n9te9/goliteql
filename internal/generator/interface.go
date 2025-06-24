package generator

import (
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateInterfaceTypeDecls(interfaces []*schema.InterfaceDefinition) []ast.Decl {
	decls := make([]ast.Decl, 0, len(interfaces))

	for _, i := range interfaces {
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(string(i.Name)),
					Type: &ast.InterfaceType{
						Methods: &ast.FieldList{
							List: []*ast.Field{},
						},
					},
				},
			},
		})
	}

	return decls
}
