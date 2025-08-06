package generator

import (
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateDirectiveImport() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"context"`,
				},
			},
		},
	}
}

var builtinDirectiveNames = map[string]struct{}{
	"include":     {},
	"skip":        {},
	"deprecated":  {},
	"specifiedBy": {},
}

func generateDirectiveDecls(directives schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	ret = append(ret, generateDirectiveImport())
	ret = append(ret, generateDirectiveInterfaceDecl(directives))
	ret = append(ret, generateDirectiveImplementationDecl())
	ret = append(ret, generateVarDirectiveDecl())
	ret = append(ret, generateNewDirectiveFuncDecl())
	ret = append(ret, generateDirectiveFuncDecls(directives)...)

	return ret
}

func generateVarDirectiveDecl() ast.Decl {
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					ast.NewIdent("_"),
				},
				Type: ast.NewIdent("Directive"),
				Values: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.CompositeLit{
							Type: ast.NewIdent("directiveImpl"),
							Elts: []ast.Expr{},
						},
					},
				},
			},
		},
	}
}

func generateNewDirectiveFuncDecl() ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("NewDirective"),
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("directiveImpl"),
						},
					},
				},
			},
			Params: &ast.FieldList{
				List: []*ast.Field{},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: ast.NewIdent("directiveImpl"),
							},
						},
					},
				},
			},
		},
	}
}

func generateDirectiveInterfaceDecl(directives schema.DirectiveDefinitions) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Directive"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: generateInterfaceDirectiveMethods(directives),
					},
				},
			},
		},
	}
}

func generateDirectiveImplementationDecl() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("directiveImpl"),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{},
					},
				},
			},
		},
	}
}

func generateDirectiveFuncDecls(directives schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, directive := range directives {
		if _, ok := builtinDirectiveNames[string(directive.Name)]; ok {
			continue
		}

		ret = append(ret, generateDirectiveFuncDecl(directive))
	}
	return ret
}

func generateDirectiveFuncDecl(directive *schema.DirectiveDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("d")},
					Type: &ast.StarExpr{
						X: ast.NewIdent("directiveImpl"),
					},
				},
			},
		},
		Name: ast.NewIdent(toUpperCase(string(directive.Name))),
		Type: generateFuncType(directive),
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: ast.NewIdent("panic"),
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"not implemented"`,
							},
						},
					},
				},
			},
		},
	}
}

func generateInterfaceDirectiveMethods(directives schema.DirectiveDefinitions) []*ast.Field {
	methods := make([]*ast.Field, 0, len(directives))

	for _, directive := range directives {
		if _, ok := builtinDirectiveNames[string(directive.Name)]; ok {
			continue
		}

		methods = append(methods, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(toUpperCase(string(directive.Name))),
			},
			Type: generateFuncType(directive),
		})
	}

	return methods
}

func generateFuncType(directive *schema.DirectiveDefinition) *ast.FuncType {
	return &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("ctx"),
					},
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("context"),
						Sel: ast.NewIdent("Context"),
					},
				}, {
					Names: []*ast.Ident{
						ast.NewIdent("ret"),
					},
					Type: ast.NewIdent("any"),
				},
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				{
					Type: ast.NewIdent("any"),
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		},
	}
}
