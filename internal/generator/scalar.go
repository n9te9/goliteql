package generator

import (
	"fmt"
	"go/ast"
	"go/token"
)

func (g *Generator) generateScalar() {
	specs := make([]ast.Spec, 0)

	for _, scalar := range g.config.Scalars {
		specs = append(specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, scalar.Package),
			},
		})
	}

	g.scalarAST.Decls = append(g.scalarAST.Decls, &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	})

	for _, scalar := range g.config.Scalars {
		g.scalarAST.Decls = append(g.scalarAST.Decls, &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(scalar.Name),
					Type: &ast.Ident{
						Name: scalar.Type,
					},
				},
			},
		})

		g.scalarAST.Decls = append(g.scalarAST.Decls, generateScalarUnmarshalJSON(scalar))
		g.scalarAST.Decls = append(g.scalarAST.Decls, generateScalarMarshalJSON(scalar))
	}
}

func generateScalarUnmarshalJSON(scalar ScalarConfig) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("UnmarshalJSON"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: ast.NewIdent(scalar.Name)},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "data",
							},
						},
						Type: &ast.Ident{
							Name: "[]byte",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("panic"),
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"implement %s UnmarshalJSON"`, scalar.Name),
							},
						},
					},
				},
			},
		},
	}
}

func generateScalarMarshalJSON(scalar ScalarConfig) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("MarshalJSON"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("s")},
					Type:  &ast.StarExpr{X: ast.NewIdent(scalar.Name)},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "[]byte",
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("panic"),
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"implement %s MarshalJSON"`, scalar.Name),
							},
						},
					},
				},
			},
		},
	}
}
