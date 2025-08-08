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

func generateExtractDirectiveArgFuncs(typePrefix string, directives schema.DirectiveDefinitions, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0)

	for _, directive := range directives {
		if _, ok := builtinDirectiveNames[string(directive.Name)]; ok {
			continue
		}

		if len(directive.Arguments) == 0 {
			continue
		}

		ret = append(ret, generateExtractDirectiveFuncArgFunc(typePrefix, directive, indexes))
	}

	return ret
}

func generateExtractDirectiveFuncArgFunc(typePrefix string, directive *schema.DirectiveDefinition, indexes *schema.Indexes) ast.Decl {
	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, generateExtractArgumentBody(typePrefix, directive.Arguments)...)
	stmts = append(stmts, generateExtractDirectiveArgumentBody(directive.Arguments, indexes))
	stmts = append(stmts, generateDirectiveReturnStmt(directive.Arguments))

	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("extract%sDirectiveArgs", toUpperCase(string(directive.Name)))),
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
					},
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
							Key: ast.NewIdent("string"),
							Value: &ast.SelectorExpr{
								X:   ast.NewIdent("json"),
								Sel: ast.NewIdent("RawMessage"),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: generateDirectiveArgumentResultFieldSet(typePrefix, directive.Arguments),
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
	}
}

func generateDirectiveArgumentResultFieldSet(typePrefix string, arguemntDefinitions []*schema.ArgumentDefinition) []*ast.Field {
	ret := make([]*ast.Field, 0, len(arguemntDefinitions))

	for _, arg := range arguemntDefinitions {
		typeExpr := generateTypeExprFromFieldType(typePrefix, arg.Type)
		ret = append(ret, &ast.Field{
			Type: typeExpr,
		})
	}

	ret = append(ret, &ast.Field{
		Type: ast.NewIdent("error"),
	})

	return ret
}

func generateDirectiveReturnStmt(arguemntDefinitions []*schema.ArgumentDefinition) ast.Stmt {
	results := make([]ast.Expr, 0)
	for _, arg := range arguemntDefinitions {
		varName := string(arg.Name)
		results = append(results, ast.NewIdent(varName))
	}
	results = append(results, ast.NewIdent("nil"))

	return &ast.ReturnStmt{
		Results: results,
	}
}

func generateExtractArgumentBody(typePrefix string, arguemntDefinitions []*schema.ArgumentDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	for _, arg := range arguemntDefinitions {
		ret = append(ret, generateExtractArgumentVarDeclaration(typePrefix, arg))
	}

	return ret
}

func generateExtractArgumentVarDeclaration(typePrefix string, argumentDefinition *schema.ArgumentDefinition) ast.Stmt {
	varName := string(argumentDefinition.Name)
	typeExpr := generateTypeExprFromFieldType(typePrefix, argumentDefinition.Type)

	return &ast.DeclStmt{
		Decl: &ast.GenDecl{
			Tok: token.VAR,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names: []*ast.Ident{
						ast.NewIdent(varName),
					},
					Type: typeExpr,
				},
			},
		},
	}
}

func generateExtractDirectiveArgumentBody(args schema.ArgumentDefinitions, indexes *schema.Indexes) ast.Stmt {
	return &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("arg"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("node"),
			Sel: ast.NewIdent("Arguments"),
		},
		Tok: token.DEFINE,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				generateExtractDirectiveSwitchStmt(args, indexes),
			},
		},
	}
}

func generateExtractDirectiveSwitchStmt(args schema.ArgumentDefinitions, indexes *schema.Indexes) ast.Stmt {
	return &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("arg"),
					Sel: ast.NewIdent("Name"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateExtractDirectiveCaseStmts(args, indexes),
		},
	}
}

func generateExtractDirectiveCaseStmts(args schema.ArgumentDefinitions, indexes *schema.Indexes) []ast.Stmt {
	ret := make([]ast.Stmt, 0, len(args))

	for _, arg := range args {
		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(arg.Name)),
				},
			},
			Body: generateExtractDirectiveMap(arg, indexes),
		})
	}

	return ret
}

func generateExtractDirectiveMap(arg *schema.ArgumentDefinition, indexes *schema.Indexes) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	ret = append(ret, &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("arg"),
				Sel: ast.NewIdent("IsVariable"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
		Else: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	})

	return ret
}

func generateExtractDirectiveArgumentForVariable(arg *schema.ArgumentDefinition, indexes *schema.Indexes) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	return ret
}

func generateExtractDirectiveArgumentForBuildInVariable(arg *schema.ArgumentDefinition, indexes *schema.Indexes) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	return ret
}
