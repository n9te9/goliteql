package generator

import (
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateEnumModelAST(enums []*schema.EnumDefinition) []ast.Decl {
	decls := make([]ast.Decl, 0, len(enums))

	for _, e := range enums {

		if e.IsIntrospection() {
			decls = append(decls, &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(string(e.Name)),
						Type: &ast.Ident{
							Name: string(e.Name),
						},
					},
				},
			})
		}
	}

	return decls
}
