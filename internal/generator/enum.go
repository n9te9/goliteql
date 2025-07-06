package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateEnumModelAST(enums []*schema.EnumDefinition) []ast.Decl {
	decls := make([]ast.Decl, 0, len(enums))

	for _, e := range enums {
		if string(e.Name) != "" {
			decls = append(decls, &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(string(e.Name)),
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			})
		}

		genDecl := &ast.GenDecl{
			Tok: token.CONST,
		}

		specs := make([]ast.Spec, 0, len(e.Values))
		for _, v := range e.Values {
			specs = append(specs, &ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent(string(v.Name)),
				},
				Type: ast.NewIdent(string(e.Name)),
				Values: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, string(v.Value)),
					},
				},
			})
		}

		genDecl.Specs = specs
		decls = append(decls, genDecl)
	}

	return decls
}
