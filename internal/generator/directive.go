package generator

import (
	"fmt"
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

func generateDirectiveDecls(typePrefix string, directives schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	ret = append(ret, generateDirectiveImport())
	ret = append(ret, generateDirectiveInterfaceDecl(typePrefix, directives))
	ret = append(ret, generateDirectiveImplementationDecl())
	ret = append(ret, generateVarDirectiveDecl())
	ret = append(ret, generateNewDirectiveFuncDecl())
	ret = append(ret, generateDirectiveFuncDecls(typePrefix, directives)...)

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

func generateDirectiveInterfaceDecl(typePrefix string, directives schema.DirectiveDefinitions) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Directive"),
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: generateInterfaceDirectiveMethods(typePrefix, directives),
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

func generateDirectiveFuncDecls(typePrefix string, directives schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, directive := range directives {
		if _, ok := builtinDirectiveNames[string(directive.Name)]; ok {
			continue
		}

		ret = append(ret, generateDirectiveFuncDecl(typePrefix, directive))
	}
	return ret
}

func generateDirectiveFuncDecl(typePrefix string, directive *schema.DirectiveDefinition) *ast.FuncDecl {
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
		Type: generateFuncType(typePrefix, directive),
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

func generateInterfaceDirectiveMethods(typePrefix string, directives schema.DirectiveDefinitions) []*ast.Field {
	methods := make([]*ast.Field, 0, len(directives))

	for _, directive := range directives {
		if _, ok := builtinDirectiveNames[string(directive.Name)]; ok {
			continue
		}

		methods = append(methods, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(toUpperCase(string(directive.Name))),
			},
			Type: generateFuncType(typePrefix, directive),
		})
	}

	return methods
}

func generateFuncType(typePrefix string, directive *schema.DirectiveDefinition) *ast.FuncType {
	argsFields := make([]*ast.Field, 0)
	argsFields = append(argsFields, &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("ctx"),
		},
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("context"),
			Sel: ast.NewIdent("Context"),
		},
	}, &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("ret"),
		},
		Type: ast.NewIdent("any"),
	})

	for _, arg := range directive.Arguments {
		argsFields = append(argsFields, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(string(arg.Name)),
			},
			Type: generateTypeExprFromFieldType(typePrefix, arg.Type),
		})
	}

	return &ast.FuncType{
		Params: &ast.FieldList{
			List: argsFields,
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

func generateExtractDirectiveFuncArgFuncs(directive *schema.DirectiveDefinition) []ast.Decl {
	ret := make([]ast.Decl, 0)

	ret = append(ret, generateExtractDirectiveFuncArgFunc(directive))
	return ret
}

func generateExtractDirectiveFuncArgFunc(directive *schema.DirectiveDefinition) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("extract%sArgs", toUpperCase(string(directive.Name)))),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("node"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("Node"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("variables"),
						},
						Type: &ast.MapType{
							Key:   ast.NewIdent("string"),
							Value: ast.NewIdent("json.RawMessage"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: generateArgumentResultTypeFieldSets(directive.Arguments),
			},
		},
	}
}

func generateArgumentResultTypeFieldSets(args schema.ArgumentDefinitions) []*ast.Field {
	fields := make([]*ast.Field, 0, len(args)+1)

	fields = append(fields, &ast.Field{
		Type: ast.NewIdent("error"),
	})

	return fields
}
