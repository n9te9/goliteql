package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateApplyResponseFuncDeclFromOperationDefinition(definition *schema.OperationDefinition, typePrefix string, indexes *schema.Indexes) []ast.Decl {
	ret := make([]ast.Decl, 0, len(definition.Fields))
	for _, field := range definition.Fields {
		ret = append(ret, generateApplyResponseFuncDeclFromFieldDefinition(field, typePrefix, indexes))
	}

	return ret
}

func generateApplyResponseFuncDeclFromFieldDefinition(fieldDefinition *schema.FieldDefinition, typePrefix string, indexes *schema.Indexes) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sQueryResponse", fieldDefinition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
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
						Names: []*ast.Ident{ast.NewIdent("resolverRet")},
						Type:  generateTypeExprFromFieldTypeForReturn(typePrefix, fieldDefinition.Type, indexes),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("node")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("Node"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("variables")},
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
				List: []*ast.Field{
					{
						Type: generateResponseNullable(fieldDefinition.Type),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateFieldTypeApplyBodyStmts(fieldDefinition.Type, indexes, typePrefix),
		},
	}
}

func generateResponseNullable(fieldType *schema.FieldType) ast.Expr {
	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateResponseNullable(fieldType.ListType),
		}
	}

	return &ast.SelectorExpr{
		X:   ast.NewIdent("executor"),
		Sel: ast.NewIdent("Nullable"),
	}
}

type Responsible interface {
	*schema.TypeDefinition |
		*schema.InterfaceDefinition |
		*schema.UnionDefinition |
		*schema.EnumDefinition |
		*schema.ScalarDefinition
}

func generateApplyResponseFuncDecl[T Responsible](definitions []T, indexes *schema.Indexes, typePrefix string) []ast.Decl {
	ret := make([]ast.Decl, 0, len(definitions))

	for _, definition := range definitions {
		switch def := any(definition).(type) {
		case *schema.TypeDefinition:
			if def.IsIntrospection() {
				continue
			}
			ret = append(ret, generateTypeApplyResponseFuncDecl(def, indexes, typePrefix))
		case *schema.InterfaceDefinition:
			ret = append(ret, generateInterfaceApplyResponseFuncDecl(def, indexes, typePrefix))
		case *schema.UnionDefinition:
			ret = append(ret, generateUnionApplyResponseFuncDecl(def, indexes, typePrefix))
		case *schema.EnumDefinition:
			ret = append(ret, generateEnumApplyResponseFuncDecl(def, indexes, typePrefix))
		case *schema.ScalarDefinition:
			ret = append(ret, generateScalarApplyResponseFuncDecl(def, indexes, typePrefix))
		}
	}

	return ret
}

func generateTypeApplyResponseFuncDecl(definition *schema.TypeDefinition, indexes *schema.Indexes, typePrefix string) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sResponse", definition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("resolverRet")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(typePrefix),
							Sel: ast.NewIdent(string(definition.Name)),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("node")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("Node"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("variables")},
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
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("Nullable"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateTypeApplyResponseFuncBody(definition, indexes),
		},
	}
}

func generateTypeApplyResponseFuncBody(definition *schema.TypeDefinition, indexes *schema.Indexes) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("ret"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CompositeLit{
					Type: ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
				},
			},
		},
		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("child"),
			Tok:   token.DEFINE,
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.IfStmt{
						Cond: &ast.CallExpr{
							Fun: &ast.SelectorExpr{X: ast.NewIdent("node"), Sel: ast.NewIdent("HasFragment")},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{X: ast.NewIdent("node"), Sel: ast.NewIdent("FragmentType")},
										},
										Op: token.NEQ,
										Y:  ast.NewIdent(fmt.Sprintf("%q", definition.Name)),
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													ast.NewIdent("nil"),
													&ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("fmt"),
															Sel: ast.NewIdent("Errorf"),
														},
														Args: []ast.Expr{
															ast.NewIdent(fmt.Sprintf(`"node type mismatch: expected %s"`, definition.Name)),
														},
													},
												},
											},
										},
									},
								},
								generateFragmentApplyResponseStmt(definition, indexes),
							},
						},
					},
					generateFieldApplyResponseStmts(definition, indexes),
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("executor"),
						Sel: ast.NewIdent("NewNullable"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ret"),
					},
				},
				ast.NewIdent("nil"),
			},
		},
	}
}

func generateFragmentApplyResponseStmt(definition *schema.TypeDefinition, indexes *schema.Indexes) ast.Stmt {
	return &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("fragmentChild"),
		Tok:   token.DEFINE,
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("child"),
			Sel: ast.NewIdent("Children"),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: &ast.SelectorExpr{
						X:   ast.NewIdent("fragmentChild"),
						Sel: ast.NewIdent("Name"),
					},
					Body: &ast.BlockStmt{
						List: generateCaseStmtsForTypeDefinition(definition, indexes, ast.NewIdent("fragmentChild")),
					},
				},
			},
		},
	}
}

func generateFieldApplyResponseStmts(definition *schema.TypeDefinition, indexes *schema.Indexes) ast.Stmt {
	return &ast.SwitchStmt{
		Tag: &ast.SelectorExpr{
			X:   ast.NewIdent("child"),
			Sel: ast.NewIdent("Name"),
		},
		Body: &ast.BlockStmt{
			List: generateCaseStmtsForTypeDefinition(definition, indexes, ast.NewIdent("child")),
		},
	}
}

func generateCaseStmtsForTypeDefinition(definition *schema.TypeDefinition, indexes *schema.Indexes, nestExpr ast.Expr) []ast.Stmt {
	ret := make([]ast.Stmt, 0, len(definition.Fields)+1)
	ret = append(ret, &ast.CaseClause{
		List: []ast.Expr{
			ast.NewIdent(`"__typename"`),
		},
		Body: []ast.Stmt{
			&ast.IfStmt{
				Cond: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.SelectorExpr{
							X:   nestExpr,
							Sel: ast.NewIdent("Directives"),
						},
						Sel: ast.NewIdent("ShouldInclude"),
					},
					Args: []ast.Expr{
						ast.NewIdent("variables"),
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("ret"),
									Sel: ast.NewIdent("GraphQLTypeName"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("executor"),
										Sel: ast.NewIdent("NewNullable"),
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: fmt.Sprintf(`"%s"`, definition.Name),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})

	for _, field := range definition.Fields {
		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%q", string(field.Name))),
			},
			Body: generateCaseBodyStmts(field, indexes, nestExpr),
		})
	}

	return ret
}

func generateCaseBodyStmts(field *schema.FieldDefinition, indexes *schema.Indexes, nestExpr ast.Expr) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)

	if field.Type.IsList {
		stmts = append(stmts, generateAssignMakeSliceForResponse(field))
		stmts = append(stmts, generateNestedArrayRangeStmts(field, field.Type, indexes, nestExpr, 0)...)
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   nestExpr,
						Sel: ast.NewIdent("Directives"),
					},
					Sel: ast.NewIdent("ShouldInclude"),
				},
				Args: []ast.Expr{
					ast.NewIdent("variables"),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent(toUpperCase(string(field.Name))),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", toUpperCase(string(field.Name)))),
								},
							},
						},
					},
				},
			},
		})
	} else {
		stmts = append(stmts, generateCaseRetAssignStmts(field, indexes, nestExpr)...)
	}

	return stmts
}

func generateNestedArrayRangeStmts(field *schema.FieldDefinition, fieldType *schema.FieldType, indexes *schema.Indexes, nestExpr ast.Expr, nestCount int) []ast.Stmt {
	if fieldType.IsList {
		var xExpr ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent("resolverRet"),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}
		if nestCount != 0 {
			xExpr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
		}

		if fieldType.Nullable {
			xExpr = &ast.StarExpr{
				X: xExpr,
			}
		}

		var sliceMakeStmt ast.Stmt = &ast.EmptyStmt{}
		var sliceAppendStmt ast.Stmt = &ast.EmptyStmt{}
		if nestCount > 0 {
			sliceMakeStmt = &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: generateResponseNullable(fieldType.ListType),
							},
							ast.NewIdent("0"),
							&ast.CallExpr{
								Fun: ast.NewIdent("len"),
								Args: []ast.Expr{
									xExpr,
								},
							},
						},
					},
				},
			}

			sliceAppendTarget := fmt.Sprintf("ret%d", nestCount)
			if nestCount == 1 {
				sliceAppendTarget = fmt.Sprintf("ret%s", toUpperCase(string(field.Name)))
			}

			sliceAppendStmt = &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(sliceAppendTarget),
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("append"),
						Args: []ast.Expr{
							ast.NewIdent(sliceAppendTarget),
							ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
						},
					},
				},
			}
		}

		return []ast.Stmt{
			sliceMakeStmt,
			&ast.RangeStmt{
				X:     xExpr,
				Tok:   token.DEFINE,
				Key:   ast.NewIdent("_"),
				Value: ast.NewIdent(fmt.Sprintf("v%d", nestCount)),
				Body: &ast.BlockStmt{
					List: generateNestedArrayRangeStmts(field, fieldType.ListType, indexes, nestExpr, nestCount+1),
				},
			},
			sliceAppendStmt,
		}
	}

	return generateCaseNestedRetAssignStmts(field, indexes, nestExpr, nestCount)
}

func generateCaseNestedRetAssignStmts(field *schema.FieldDefinition, indexes *schema.Indexes, nestExpr ast.Expr, nestCount int) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)
	_, isObject := indexes.TypeIndex[string(field.Type.GetRootType().Name)]
	_, isInterface := indexes.InterfaceIndex[string(field.Type.GetRootType().Name)]
	_, isUnion := indexes.UnionIndex[string(field.Type.GetRootType().Name)]
	if isObject || isInterface || isUnion {
		var argExpr ast.Expr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
		if field.Type.Nullable {
			argExpr = &ast.StarExpr{
				X: argExpr,
			}
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("ret%s", string(field.Type.GetRootType().Name))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", field.Type.GetRootType().Name)),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						argExpr,
						nestExpr,
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))

		appendTarget := fmt.Sprintf("ret%s", toUpperCase(string(field.Name)))
		if nestCount > 1 {
			appendTarget = fmt.Sprintf("ret%d", nestCount-1)
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(appendTarget),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent(appendTarget),
						ast.NewIdent(fmt.Sprintf("ret%s", string(field.Type.GetRootType().Name))),
					},
				},
			},
		})
	} else {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   nestExpr,
						Sel: ast.NewIdent("Directives"),
					},
					Sel: ast.NewIdent("ShouldInclude"),
				},
				Args: []ast.Expr{
					ast.NewIdent("variables"),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent(toUpperCase(string(field.Name))),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("resolverRet"),
										Sel: ast.NewIdent(toUpperCase(string(field.Name))),
									},
								},
							},
						},
					},
				},
			},
		})
	}

	return stmts
}

func generateCaseRetAssignStmts(field *schema.FieldDefinition, indexes *schema.Indexes, nestExpr ast.Expr) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)
	_, isObject := indexes.TypeIndex[string(field.Type.GetRootType().Name)]
	_, isInterface := indexes.InterfaceIndex[string(field.Type.GetRootType().Name)]
	_, isUnion := indexes.UnionIndex[string(field.Type.GetRootType().Name)]
	if isObject || isInterface || isUnion {
		var argExpr ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent("resolverRet"),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}
		if field.Type.Nullable {
			argExpr = &ast.StarExpr{
				X: argExpr,
			}
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("ret%s", toUpperCase(string(field.Name)))),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", field.Type.GetRootType().Name)),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						argExpr,
						nestExpr,
						ast.NewIdent("variables"),
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   nestExpr,
						Sel: ast.NewIdent("Directives"),
					},
					Sel: ast.NewIdent("ShouldInclude"),
				},
				Args: []ast.Expr{
					ast.NewIdent("variables"),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent(toUpperCase(string(field.Name))),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", toUpperCase(string(field.Name)))),
								},
							},
						},
					},
				},
			},
		})
	} else {
		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.SelectorExpr{
						X:   nestExpr,
						Sel: ast.NewIdent("Directives"),
					},
					Sel: ast.NewIdent("ShouldInclude"),
				},
				Args: []ast.Expr{
					ast.NewIdent("variables"),
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("ret"),
								Sel: ast.NewIdent(toUpperCase(string(field.Name))),
							},
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("resolverRet"),
										Sel: ast.NewIdent(toUpperCase(string(field.Name))),
									},
								},
							},
						},
					},
				},
			},
		})
	}

	return stmts
}

func generateAssignMakeSliceForResponse(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	var retExpr ast.Expr = ast.NewIdent(fmt.Sprintf("ret%s", toUpperCase(string(fieldDefinition.Name))))

	var argExpr ast.Expr = &ast.SelectorExpr{
		X:   ast.NewIdent("resolverRet"),
		Sel: ast.NewIdent(toUpperCase(string(fieldDefinition.Name))),
	}
	if fieldDefinition.Type.Nullable {
		argExpr = &ast.StarExpr{
			X: argExpr,
		}
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{retExpr},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					generateResponseNullable(fieldDefinition.Type),
					ast.NewIdent("0"),
					&ast.CallExpr{
						Fun: ast.NewIdent("len"),
						Args: []ast.Expr{
							argExpr,
						},
					},
				},
			},
		},
	}
}

func generateResponseTypeExpr(fieldType *schema.FieldType) ast.Expr {
	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateResponseTypeExpr(fieldType.ListType),
		}
	}

	return ast.NewIdent(fmt.Sprintf("%sResponse", fieldType.Name))
}

func generateFieldTypeApplyBodyStmts(fieldType *schema.FieldType, indexes *schema.Indexes, typePrefix string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if fieldType.IsList {
		ret = append(ret, generateFieldTypeRangeBodyStmts(fieldType, indexes, 0)...)
		ret = append(ret, &ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
				ast.NewIdent("nil"),
			},
		})
	} else {
		_, isObject := indexes.TypeIndex[string(fieldType.GetRootType().Name)]
		_, isInterface := indexes.InterfaceIndex[string(fieldType.GetRootType().Name)]
		_, isUnion := indexes.UnionIndex[string(fieldType.GetRootType().Name)]
		if isObject || isInterface || isUnion {
			var nilAssignStmt ast.Stmt = &ast.EmptyStmt{}

			var arg ast.Expr = ast.NewIdent("resolverRet")
			if isObject && fieldType.Nullable {
				nilAssignStmt = &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("resolverRet"),
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("executor"),
											Sel: ast.NewIdent("NewNullable"),
										},
										Args: []ast.Expr{
											ast.NewIdent("nil"),
										},
									},
									ast.NewIdent("nil"),
								},
							},
						},
					},
				}
				arg = &ast.StarExpr{
					X: arg,
				}
			}

			if isInterface || isUnion {
				nilAssignStmt = &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("resolverRet"),
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("executor"),
											Sel: ast.NewIdent("NewNullable"),
										},
										Args: []ast.Expr{
											ast.NewIdent("nil"),
										},
									},
									ast.NewIdent("nil"),
								},
							},
						},
					},
				}
			}

			ret = append(ret, nilAssignStmt)
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", fieldType.Name)),
						},
						Args: []ast.Expr{
							ast.NewIdent("ctx"),
							arg,
							ast.NewIdent("node"),
							ast.NewIdent("variables"),
						},
					},
				},
			})
			ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}))
			ret = append(ret, &ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("ret"),
					ast.NewIdent("nil"),
				},
			})
		}
	}

	return ret
}

func generateFieldTypeRangeBodyStmts(fieldType *schema.FieldType, indexes *schema.Indexes, nestCount int) []ast.Stmt {
	var xExpr ast.Expr = ast.NewIdent("resolverRet")
	var retExpr ast.Expr = ast.NewIdent("ret")
	if nestCount != 0 {
		xExpr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
		retExpr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
	}

	var nestValueExpr ast.Expr = ast.NewIdent(fmt.Sprintf("v%d", nestCount))
	var appendStmt ast.Stmt = &ast.EmptyStmt{}
	if nestCount > 0 {
		ret := fmt.Sprintf("ret%d", nestCount-1)
		appendArg := fmt.Sprintf("ret%d", nestCount)
		if nestCount == 1 {
			ret = "ret"
			appendArg = "ret1"
		}

		appendStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(ret),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent(ret),
						ast.NewIdent(appendArg),
					},
				},
			},
		}
	}

	if fieldType.Nullable && fieldType.IsList {
		if _, ok := indexes.TypeIndex[string(fieldType.GetRootType().Name)]; ok {
			xExpr = &ast.StarExpr{X: xExpr}
		}

		return []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{retExpr},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: generateResponseNullable(fieldType.ListType),
							},
							ast.NewIdent("0"),
							&ast.CallExpr{
								Fun: ast.NewIdent("len"),
								Args: []ast.Expr{
									xExpr,
								},
							},
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Value: nestValueExpr,
				X:     xExpr,
				Tok:   token.DEFINE,
				Body: &ast.BlockStmt{
					List: generateFieldTypeRangeBodyStmts(fieldType.ListType, indexes, nestCount+1),
				},
			},
			appendStmt,
		}
	}
	if fieldType.IsList {
		return []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{retExpr},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("make"),
						Args: []ast.Expr{
							&ast.ArrayType{
								Elt: generateResponseNullable(fieldType.ListType),
							},
							ast.NewIdent("0"),
							&ast.CallExpr{
								Fun: ast.NewIdent("len"),
								Args: []ast.Expr{
									xExpr,
								},
							},
						},
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Value: nestValueExpr,
				X:     xExpr,
				Tok:   token.DEFINE,
				Body: &ast.BlockStmt{
					List: generateFieldTypeRangeBodyStmts(fieldType.ListType, indexes, nestCount+1),
				},
			},
			appendStmt,
		}
	}

	var argExpr ast.Expr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
	if fieldType.Nullable {
		argExpr = &ast.StarExpr{
			X: argExpr,
		}
	}

	var appendCheckStmt ast.Stmt = appendStmt
	_, isUnion := indexes.UnionIndex[string(fieldType.GetRootType().Name)]
	if isUnion {
		appendCheckStmt = &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
				Op: token.NEQ,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					appendStmt,
				},
			},
		}
	}

	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", fieldType.Name)),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						argExpr,
						ast.NewIdent("node"),
						ast.NewIdent("variables"),
					},
				},
			},
		},
		generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
		appendCheckStmt,
	}
}

func generateInterfaceApplyResponseFuncDecl(definition *schema.InterfaceDefinition, indexes *schema.Indexes, typePrefix string) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sResponse", definition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("resolverRet")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(typePrefix),
							Sel: ast.NewIdent(string(definition.Name)),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("node")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("Node"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("variables")},
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
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("Nullable"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateInterfaceApplyResponseBody(definition, indexes, typePrefix),
		},
	}
}

func generateInterfaceApplyResponseBody[T *schema.InterfaceDefinition | *schema.UnionDefinition](definition T, indexes *schema.Indexes, typePrefix string) []ast.Stmt {
	switch definition := any(definition).(type) {
	case *schema.InterfaceDefinition:
		types := indexes.GetImplementedType(definition)
		return generateInterfaceApplySwitchStmtsForInterfaceDefinition(definition, types, indexes, typePrefix)
	case *schema.UnionDefinition:
		types := getUnionTypes(definition, indexes)
		return generateInterfaceApplySwitchStmtsForUnionDefinition(definition, types, indexes, typePrefix)
	}

	panic(fmt.Sprintf("unsupported definition type: %T", definition))
}

func getUnionTypes(definition *schema.UnionDefinition, indexes *schema.Indexes) schema.TypeDefinitions {
	types := make(schema.TypeDefinitions, 0, len(definition.Types))
	for _, typeName := range definition.Types {
		if typeDef, ok := indexes.TypeIndex[string(typeName)]; ok {
			types = append(types, typeDef)
		} else {
			panic(fmt.Sprintf("type %s not found in index", typeName))
		}
	}

	return types
}

func generateInterfaceApplySwitchStmtsForInterfaceDefinition(definition *schema.InterfaceDefinition, types schema.TypeDefinitions, indexes *schema.Indexes, typePrefix string) []ast.Stmt {
	typeCases := make([]ast.Stmt, 0, len(types))
	for _, typeDef := range types {
		typeCases = append(typeCases, &ast.CaseClause{
			List: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(typePrefix),
					Sel: ast.NewIdent(string(typeDef.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", typeDef.Name)),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("resolverRet"),
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
					},
				},
			},
		})

		typeCases = append(typeCases, &ast.CaseClause{
			List: []ast.Expr{
				&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(typePrefix),
						Sel: ast.NewIdent(string(typeDef.Name)),
					},
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", typeDef.Name)),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.StarExpr{
									X: ast.NewIdent("resolverRet"),
								},
								ast.NewIdent("node"),
								ast.NewIdent("variables"),
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("ret"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
					},
				},
			},
		})
	}

	return []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("ret")},
						Type:  ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
					},
				},
			},
		},
		&ast.TypeSwitchStmt{
			Assign: &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("resolverRet"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X:    ast.NewIdent("resolverRet"),
						Type: ast.NewIdent("type"),
					},
				},
			},
			Body: &ast.BlockStmt{
				List: typeCases,
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("ret"),
				Op: token.EQL,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("executor"),
						Sel: ast.NewIdent("NewNullable"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ret"),
					},
				},
				ast.NewIdent("nil"),
			},
		},
	}
}

func generateInterfaceApplySwitchStmtsForUnionDefinition(definition *schema.UnionDefinition, types schema.TypeDefinitions, indexes *schema.Indexes, typePrefix string) []ast.Stmt {
	typeCases := make([]ast.Stmt, 0, len(types))
	for _, typeDef := range types {
		typeCases = append(typeCases, &ast.CaseClause{
			List: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(typePrefix),
					Sel: ast.NewIdent(string(typeDef.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("node"),
								Sel: ast.NewIdent("FragmentType"),
							},
						},
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, typeDef.Name),
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("r"),
											Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", typeDef.Name)),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("resolverRet"),
											ast.NewIdent("child"),
											ast.NewIdent("variables"),
										},
									},
								},
							},
							generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("ret"),
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
								},
							},
						},
					},
				},
			},
		})

		typeCases = append(typeCases, &ast.CaseClause{
			List: []ast.Expr{
				&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(typePrefix),
						Sel: ast.NewIdent(string(typeDef.Name)),
					},
				},
			},
			Body: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("child"),
								Sel: ast.NewIdent("FragmentType"),
							},
						},
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, typeDef.Name),
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("r"),
											Sel: ast.NewIdent(fmt.Sprintf("apply%sResponse", typeDef.Name)),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											&ast.StarExpr{
												X: ast.NewIdent("resolverRet"),
											},
											ast.NewIdent("child"),
											ast.NewIdent("variables"),
										},
									},
								},
							},
							generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("nil")}),
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("ret"),
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									ast.NewIdent(fmt.Sprintf("ret%s", typeDef.Name)),
								},
							},
						},
					},
				},
			},
		})
	}

	return []ast.Stmt{
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("resolverRet"),
				Op: token.EQL,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("executor"),
									Sel: ast.NewIdent("NewNullable"),
								},
								Args: []ast.Expr{
									ast.NewIdent("nil"),
								},
							},
							ast.NewIdent("nil"),
						},
					},
				},
			},
		},

		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("ret")},
						Type:  ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
					},
				},
			},
		},

		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("child"),
			Tok:   token.DEFINE,
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("Children"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.TypeSwitchStmt{
						Assign: &ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("resolverRet"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent("resolverRet"),
									Type: ast.NewIdent("type"),
								},
							},
						},
						Body: &ast.BlockStmt{
							List: typeCases,
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("ret"),
				Op: token.EQL,
				Y: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "nil",
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("nil"),
							ast.NewIdent("nil"),
						},
					},
				},
			},
		},
		&ast.ReturnStmt{
			Results: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("executor"),
						Sel: ast.NewIdent("NewNullable"),
					},
					Args: []ast.Expr{
						ast.NewIdent("ret"),
					},
				},
				ast.NewIdent("nil"),
			},
		},
	}
}

func generateUnionApplyResponseFuncDecl(definition *schema.UnionDefinition, indexes *schema.Indexes, typePrefix string) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sResponse", definition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},

		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("resolverRet")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent(typePrefix),
							Sel: ast.NewIdent(string(definition.Name)),
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("node")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("Node"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("variables")},
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
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("Nullable"),
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateInterfaceApplyResponseBody(definition, indexes, typePrefix),
		},
	}
}

func generateScalarApplyResponseFuncDecl(definition *schema.ScalarDefinition, indexes *schema.Indexes, typePrefix string) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sResponse", definition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}
}

func generateEnumApplyResponseFuncDecl(definition *schema.EnumDefinition, indexes *schema.Indexes, typePrefix string) ast.Decl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%sResponse", definition.Name)),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{},
		},
	}
}
