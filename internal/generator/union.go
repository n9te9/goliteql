package generator

import (
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateUnionTypeDecls(unions schema.UnionDefinitions) []ast.Decl {
	decls := make([]ast.Decl, 0, len(unions))

	for _, u := range unions {
		decls = append(decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(string(u.Name)),
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
