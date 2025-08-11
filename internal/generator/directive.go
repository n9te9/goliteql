package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateDirectiveImport(modelPath string) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"context"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, modelPath),
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

func generateDirectiveDecls(typePrefix, modelPath string, directives schema.DirectiveDefinitions) []ast.Decl {
	ret := make([]ast.Decl, 0)

	ret = append(ret, generateDirectiveImport(modelPath))
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

func generateGetDirectiveStmt(directive *schema.DirectiveDefinition) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("dir"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("node"),
						Sel: ast.NewIdent("Directives"),
					},
					Sel: ast.NewIdent("FindByName"),
				},
				Args: []ast.Expr{
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf(`"%s"`, string(directive.Name)),
					},
				},
			},
		},
	}
}

func generateExtractDirectiveFuncArgFunc(typePrefix string, directive *schema.DirectiveDefinition, indexes *schema.Indexes) ast.Decl {
	stmts := make([]ast.Stmt, 0)
	stmts = append(stmts, generateExtractArgumentBody(typePrefix, directive.Arguments)...)
	stmts = append(stmts, generateGetDirectiveStmt(directive))
	stmts = append(stmts, generateExtractDirectiveArgumentBody(typePrefix, directive.Arguments, indexes))
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
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("Node"),
							},
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

func generateExtractDirectiveArgumentBody(typePrefix string, args schema.ArgumentDefinitions, indexes *schema.Indexes) ast.Stmt {
	return &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("arg"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("dir"),
			Sel: ast.NewIdent("Arguments"),
		},
		Tok: token.DEFINE,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				generateExtractDirectiveSwitchStmt(typePrefix, args, indexes),
			},
		},
	}
}

func generateExtractDirectiveSwitchStmt(typePrefix string, args schema.ArgumentDefinitions, indexes *schema.Indexes) ast.Stmt {
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
			List: generateExtractDirectiveCaseStmts(typePrefix, args, indexes),
		},
	}
}

func generateExtractDirectiveCaseStmts(typePrefix string, args schema.ArgumentDefinitions, indexes *schema.Indexes) []ast.Stmt {
	ret := make([]ast.Stmt, 0, len(args))

	for _, arg := range args {
		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(arg.Name)),
				},
			},
			Body: []ast.Stmt{
				generateExtractDirectiveMap(typePrefix, arg, indexes),
			},
		})
	}

	return ret
}

func generateExtractDirectiveMap(typePrefix string, arg *schema.ArgumentDefinition, indexes *schema.Indexes) ast.Stmt {
	var variableStmts, builtInStmts []ast.Stmt
	if d, ok := indexes.ScalarIndex[string(arg.Type.GetRootType().Name)]; ok {
		variableStmts = generateExtractDirectiveArgumentForVariable(typePrefix, arg, d)
		builtInStmts = generateExtractDirectiveArgumentForBuiltInVariable(typePrefix, arg, d)
	}

	if d, ok := indexes.EnumIndex[string(arg.Type.GetRootType().Name)]; ok {
		variableStmts = generateExtractDirectiveArgumentForVariable(typePrefix, arg, d)
		builtInStmts = generateExtractDirectiveArgumentForBuiltInVariable(typePrefix, arg, d)
	}

	return &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("arg"),
				Sel: ast.NewIdent("IsVariable"),
			},
		},
		Body: &ast.BlockStmt{
			List: variableStmts,
		},
		Else: &ast.BlockStmt{
			List: builtInStmts,
		},
	}
}

func generateExtractDirectiveArgumentForVariable[T *schema.ScalarDefinition | *schema.EnumDefinition | *schema.InputDefinition](typePrefix string, arg *schema.ArgumentDefinition, definition T) []ast.Stmt {
	switch d := any(definition).(type) {
	case *schema.ScalarDefinition:
		return generateExtractDirectiveScalarArgumentForVariable(typePrefix, arg, d)
	case *schema.EnumDefinition:
		return generateExtractDirectiveEnumArgumentForVariable(typePrefix, arg, d)
	case *schema.InputDefinition:

	}

	return nil
}

func generateExtractDirectiveScalarArgumentForVariable(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.ScalarDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
			ast.NewIdent("ok"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.IndexExpr{
				X: ast.NewIdent("variables"),
				Index: &ast.CallExpr{
					Fun: ast.NewIdent("string"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("arg"),
								Sel: ast.NewIdent("VariableAnnotation"),
							},
						},
					},
				},
			},
		},
	})

	ret = append(ret, &ast.IfStmt{
		Cond: &ast.UnaryExpr{
			Op: token.NOT,
			X:  ast.NewIdent("ok"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Errorf"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("`%s is required`", string(arg.Name)),
								},
							},
						},
					},
				},
			},
		},
	})

	if arg.Type.IsPrimitive() {
		ret = append(ret, generateExtractDirectivePremitiveScalarArgumentForVariable(typePrefix, arg, definition)...)
	} else {
		ret = append(ret, generateExtractDirectiveCustomScalarArgumentForVariable(arg)...)
	}

	return ret
}

func generateExtractDirectivePremitiveScalarArgumentForVariable(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.ScalarDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if arg.Type.IsBoolean() {
		var rhExpr ast.Expr = &ast.BinaryExpr{
			X:  ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
			Op: token.EQL,
			Y: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"true"`,
			},
		}

		if arg.Type.Nullable {
			rhExpr = pointerExpr(ast.NewIdent("bool"), rhExpr)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rhExpr,
			},
		})
	}

	if arg.Type.IsString() || arg.Type.IsID() {
		var rhExpr ast.Expr = &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SliceExpr{
					X: ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BasicLit{
						Kind:  token.INT,
						Value: "len(req" + toUpperCase(string(arg.Name)) + ") - 1",
					},
				},
			},
		}

		if arg.Type.Nullable {
			rhExpr = pointerExpr(ast.NewIdent("string"), rhExpr)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rhExpr,
			},
		})
	}

	if arg.Type.IsInt() {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.SliceExpr{
					X: ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BasicLit{
						Kind:  token.INT,
						Value: "len(arg.Value) - 1",
					},
				},
			},
		})

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sInt", arg.Name)),
				ast.NewIdent("_"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("strconv.Atoi"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
							},
						},
					},
				},
			},
		})

		var rh ast.Expr = ast.NewIdent(fmt.Sprintf("%sInt", arg.Name))
		if arg.Type.Nullable {
			rh = pointerExpr(ast.NewIdent("int"), rh)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rh,
			},
		})
	}

	if arg.Type.IsFloat() {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.SliceExpr{
					X: ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BasicLit{
						Kind:  token.INT,
						Value: "len(arg.Value) - 1",
					},
				},
			},
		})

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sFloat", arg.Name)),
				ast.NewIdent("_"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("strconv.ParseFloat"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
							},
						},
						ast.NewIdent("64"),
					},
				},
			},
		})

		var rh ast.Expr = ast.NewIdent(fmt.Sprintf("%sFloat", arg.Name))
		if arg.Type.Nullable {
			rh = pointerExpr(ast.NewIdent("float64"), rh)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rh,
			},
		})
	}

	return ret
}

func generateExtractDirectiveCustomScalarArgumentForVariable(arg *schema.ArgumentDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("json"),
						Sel: ast.NewIdent("Unmarshal"),
					},
					Args: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(string(arg.Name)),
						},
					},
				},
			},
		},
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `nil`,
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Errorf"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("`%s is required`", string(arg.Name)),
								},
							},
						},
					},
				},
			},
		},
	})

	return ret
}

func generateExtractDirectiveBuiltinCustomScalarArgument(arg *schema.ArgumentDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ast"),
			ast.NewIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent("parser"),
						},
						Sel: ast.NewIdent("ValueParser"),
					},
					Sel: ast.NewIdent("Parse"),
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
				},
			},
		},
	})
	ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent(string(arg.Name))}))
	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("jsonBytes"),
			ast.NewIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("ast"),
					Sel: ast.NewIdent("JSONMarshal"),
				},
			},
		},
	})
	ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent(string(arg.Name))}))
	ret = append(ret, &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("json"),
						Sel: ast.NewIdent("Unmarshal"),
					},
					Args: []ast.Expr{
						ast.NewIdent("jsonBytes"),
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(string(arg.Name)),
						},
					},
				},
			},
		},
	})

	return ret
}

func generateExtractDirectiveEnumArgumentForVariable(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.EnumDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
			ast.NewIdent("ok"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.IndexExpr{
				X: ast.NewIdent("variables"),
				Index: &ast.CallExpr{
					Fun: ast.NewIdent("string"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("arg"),
								Sel: ast.NewIdent("VariableAnnotation"),
							},
						},
					},
				},
			},
		},
	})

	ret = append(ret, &ast.IfStmt{
		Cond: &ast.UnaryExpr{
			Op: token.NOT,
			X:  ast.NewIdent("ok"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Errorf"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("`%s is required`", string(arg.Name)),
								},
							},
						},
					},
				},
			},
		},
	})

	caseStmts := make([]ast.Stmt, 0, len(definition.Values))
	for _, value := range definition.Values {
		var rhs ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(typePrefix),
			Sel: ast.NewIdent(string(value.Name)),
		}
		if arg.Type.Nullable {
			rhs = pointerExpr(&ast.SelectorExpr{
				X:   ast.NewIdent(typePrefix),
				Sel: ast.NewIdent(string(definition.Name)),
			}, rhs)
		}

		caseStmts = append(caseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(value.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						rhs,
					},
				},
			},
		})
	}

	ret = append(ret, &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("req%s", toUpperCase(string(arg.Name)))),
			},
		},
		Body: &ast.BlockStmt{
			List: caseStmts,
		},
	})

	return ret
}

func generateExtractDirectiveInputArgumentForVariable(arg *schema.ArgumentDefinition, definition *schema.EnumDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	return ret
}

func generateExtractDirectiveArgumentForBuiltInVariable[T *schema.ScalarDefinition | *schema.EnumDefinition | *schema.InputDefinition](typePrefix string, arg *schema.ArgumentDefinition, definition T) []ast.Stmt {
	switch d := any(definition).(type) {
	case *schema.ScalarDefinition:
		return generateExtractDirectiveScalarArgumentForBuiltInVariable(typePrefix, arg, d)
	case *schema.EnumDefinition:
		return generateExtractDirectiveEnumArgumentForBuiltInVariable(typePrefix, arg, d)
	case *schema.InputDefinition:

	}

	return nil
}

func generateExtractDirectiveScalarArgumentForBuiltInVariable(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.ScalarDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if arg.Type.IsPrimitive() {
		ret = append(ret, generateExtractDirectivePremitiveScalarArgumentForBuiltIn(typePrefix, arg, definition)...)
	} else {
		ret = append(ret, generateExtractDirectiveCustomScalarArgumentForBuiltIn(arg)...)
	}

	return ret
}

func generateExtractDirectivePremitiveScalarArgumentForBuiltIn(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.ScalarDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if arg.Type.IsBoolean() {
		var rhExpr ast.Expr = &ast.BinaryExpr{
			X: &ast.CallExpr{
				Fun: ast.NewIdent("string"),
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
				},
			},
			Op: token.EQL,
			Y: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"true"`,
			},
		}

		if arg.Type.Nullable {
			rhExpr = pointerExpr(ast.NewIdent("bool"), rhExpr)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rhExpr,
			},
		})
	}

	if arg.Type.IsString() || arg.Type.IsID() {
		var rhExpr ast.Expr = &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SliceExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("arg"),
									Sel: ast.NewIdent("Value"),
								},
							},
						},
						Op: token.SUB,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "1",
						},
					},
				},
			},
		}

		if arg.Type.Nullable {
			rhExpr = pointerExpr(ast.NewIdent("string"), rhExpr)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rhExpr,
			},
		})
	}

	if arg.Type.IsInt() {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.SliceExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("arg"),
									Sel: ast.NewIdent("Value"),
								},
							},
						},
						Op: token.SUB,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "1",
						},
					},
				},
			},
		})

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sInt", arg.Name)),
				ast.NewIdent("_"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("strconv.Atoi"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
							},
						},
					},
				},
			},
		})

		var rh ast.Expr = ast.NewIdent(fmt.Sprintf("%sInt", arg.Name))
		if arg.Type.Nullable {
			rh = pointerExpr(ast.NewIdent("int"), rh)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rh,
			},
		})
	}

	if arg.Type.IsFloat() {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.SliceExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
					Low: &ast.BasicLit{
						Kind:  token.INT,
						Value: "1",
					},
					High: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("arg"),
									Sel: ast.NewIdent("Value"),
								},
							},
						},
						Op: token.SUB,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "1",
						},
					},
				},
			},
		})

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sFloat", arg.Name)),
				ast.NewIdent("_"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("strconv.ParseFloat"),
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("%sStr", arg.Name)),
							},
						},
						ast.NewIdent("64"),
					},
				},
			},
		})

		var rh ast.Expr = ast.NewIdent(fmt.Sprintf("%sFloat", arg.Name))
		if arg.Type.Nullable {
			rh = pointerExpr(ast.NewIdent("float64"), rh)
		}

		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(string(arg.Name)),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rh,
			},
		})
	}

	return ret
}

func generateExtractDirectiveCustomScalarArgumentForBuiltIn(arg *schema.ArgumentDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ast"),
			ast.NewIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent("parser"),
						},
						Sel: ast.NewIdent("ValueParser"),
					},
					Sel: ast.NewIdent("Parse"),
				},
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("arg"),
						Sel: ast.NewIdent("Value"),
					},
				},
			},
		},
	})
	ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent(string(arg.Name))}))
	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("jsonBytes"),
			ast.NewIdent("err"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("ast"),
					Sel: ast.NewIdent("JSONBytes"),
				},
				Args: []ast.Expr{},
			},
		},
	})
	ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent(string(arg.Name))}))
	ret = append(ret, &ast.IfStmt{
		Init: &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("json"),
						Sel: ast.NewIdent("Unmarshal"),
					},
					Args: []ast.Expr{
						ast.NewIdent("jsonBytes"),
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(string(arg.Name)),
						},
					},
				},
			},
		},
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `nil`,
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
						ast.NewIdent("err"),
					},
				},
			},
		},
	})

	return ret
}

func generateExtractDirectiveEnumArgumentForBuiltInVariable(typePrefix string, arg *schema.ArgumentDefinition, definition *schema.EnumDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	caseStmts := make([]ast.Stmt, 0, len(definition.Values))
	for _, value := range definition.Values {
		var rhs ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(typePrefix),
			Sel: ast.NewIdent(string(value.Name)),
		}
		if arg.Type.Nullable {
			rhs = pointerExpr(&ast.SelectorExpr{
				X:   ast.NewIdent(typePrefix),
				Sel: ast.NewIdent(string(definition.Name)),
			}, rhs)
		}

		caseStmts = append(caseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(value.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(string(arg.Name)),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						rhs,
					},
				},
			},
		})
	}

	ret = append(ret, &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("arg"),
					Sel: ast.NewIdent("Value"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: caseStmts,
		},
	})

	return ret
}

func pointerExpr(typeExpr, valueExpr ast.Expr) ast.Expr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.IndexExpr{
			X: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: typeExpr,
				},
				Elts: []ast.Expr{
					valueExpr,
				},
			},
			Index: &ast.BasicLit{
				Kind:  token.INT,
				Value: "0",
			},
		},
	}
}

func generateExtractDirectiveInputArgumentForBuiltInVariable(arg *schema.ArgumentDefinition, definition *schema.InputDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	return ret
}
