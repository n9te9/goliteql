package introspection

import (
	"go/ast"

	"github.com/n9te9/goliteql/schema"
)

func generateDirectiveFuncDecl(directiveDefinition *schema.DirectiveDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("Directive" + string(directiveDefinition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{},
	}
}

func GenerateDirectiveFuncDecls(directiveDefinitions schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0, len(directiveDefinitions))

	for _, directiveDefinition := range directiveDefinitions {
		ret = append(ret, generateDirectiveFuncDecl(directiveDefinition))
	}

	return ret
}
