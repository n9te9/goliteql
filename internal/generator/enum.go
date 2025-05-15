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
		genDecl := &ast.GenDecl{
			Tok: token.CONST,
		}

		specs := make([]ast.Spec, 0, len(e.Values))
		for _, v := range e.Values {
			specs = append(specs, &ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent(string(v.Name)),
				},
				Type: ast.NewIdent(string(e.Type.Name)),
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
