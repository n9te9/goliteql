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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: generateNestedArrayTypeForResponse(fieldDefinition.Type, indexes, true),
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
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
													ast.NewIdent("ret"),
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
				ast.NewIdent("ret"),
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
					Tag: &ast.CallExpr{
						Fun: ast.NewIdent("string"),
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("fragmentChild"),
								Sel: ast.NewIdent("Name"),
							},
						},
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
		Tag: &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("child"),
					Sel: ast.NewIdent("Name"),
				},
			},
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
		stmts = append(stmts, &ast.AssignStmt{
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
								Elt: generateResponseTypeExpr(fieldType.ListType),
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
						argExpr,
						nestExpr,
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("ret")}))

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
		stmts = append(stmts, &ast.AssignStmt{
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
						argExpr,
						nestExpr,
					},
				},
			},
		})
		stmts = append(stmts, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("ret")}))
		stmts = append(stmts, &ast.AssignStmt{
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
		})
	} else {
		stmts = append(stmts, &ast.AssignStmt{
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
					generateResponseTypeExpr(fieldDefinition.Type),
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
			var arg ast.Expr = ast.NewIdent("resolverRet")
			if isObject && fieldType.Nullable {
				arg = &ast.StarExpr{
					X: arg,
				}
			}

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
							arg,
							ast.NewIdent("node"),
						},
					},
				},
			})
			ret = append(ret, generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("ret")}))
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
		ret := fmt.Sprintf("ret%d", nestCount)
		if nestCount == 1 {
			ret = "ret"
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
						ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
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
								Elt: generateNestedArrayTypeForResponse(fieldType.ListType, indexes, false),
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
								Elt: generateNestedArrayTypeForResponse(fieldType.ListType, indexes, false),
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
						argExpr,
						ast.NewIdent("node"),
					},
				},
			},
		},
		generateReturnErrorHandlingStmt([]ast.Expr{ast.NewIdent("ret")}),
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
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
								ast.NewIdent("resolverRet"),
								ast.NewIdent("node"),
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
								&ast.StarExpr{
									X: ast.NewIdent("resolverRet"),
								},
								ast.NewIdent("node"),
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
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
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
											ast.NewIdent("resolverRet"),
											ast.NewIdent("child"),
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
											&ast.StarExpr{
												X: ast.NewIdent("resolverRet"),
											},
											ast.NewIdent("child"),
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
		&ast.ReturnStmt{
			Results: []ast.Expr{
				ast.NewIdent("ret"),
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", definition.Name)),
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

func generateApplyQueryResponseFuncDeclFromField(field *schema.FieldDefinition, indexes *schema.Indexes, typePrefix, operationPrefix string, isRoot bool) ast.Decl {
	nestCount := getNestCount(field.Type, 0)

	stmts := generateApplyQueryResponseFuncStmts(field, indexes, 0, nestCount, typePrefix, operationPrefix)
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("ret"),
			ast.NewIdent("nil"),
		},
	})

	if isRoot {
		operationPrefix = ""
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationPrefix, string(field.Name))),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("r")},
					Type:  &ast.StarExpr{X: ast.NewIdent("resolver")},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: stmts,
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("resolverRet"),
						},
						Type: generateTypeExprFromFieldTypeForReturn(typePrefix, field.Type, indexes),
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: generateNestedArrayTypeForResponse(field.Type, indexes, true),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
}

func generateApplyStmtForObjectQueryResponse(field *schema.FieldDefinition, typeDefinition *schema.TypeDefinition, indexes *schema.Indexes, currentNestCount, nestCount int, typePrefix, operationPrefix string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	var resultLh ast.Expr = ast.NewIdent("ret")
	if nestCount > 0 && currentNestCount > 0 {
		resultLh = ast.NewIdent(fmt.Sprintf("ret%d", currentNestCount))
	}

	if nestCount == 0 && field.Type.Nullable {
		ret = append(ret, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X:  ast.NewIdent("resolverRet"),
				Op: token.EQL,
				Y:  ast.NewIdent("nil"),
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
		})
	}

	rootType := field.Type.GetRootType()
	if currentNestCount == nestCount {
		if rootType.Nullable {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{resultLh},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("new"),
						Args: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)),
						},
					},
				},
			})
		} else {
			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{resultLh},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)),
					},
				},
			})
		}

		return append(ret, generateObjectFragmentStmtForQueryResponse(field, indexes, currentNestCount, nestCount, typePrefix, operationPrefix)...)
	}

	if currentNestCount < nestCount {
		var rangeXExpr ast.Expr = ast.NewIdent("resolverRet")
		if currentNestCount > 0 {
			rangeXExpr = ast.NewIdent(fmt.Sprintf("v%d", currentNestCount-1))
		}

		// e.g. ret = make([]*PostResponse, 0, len(v0))
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{resultLh},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						generateNestedArrayTypeForResponse(field.Type, indexes, true),
						ast.NewIdent("0"),
						&ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								rangeXExpr,
							},
						},
					},
				},
			},
		})

		// e.g. for _, v0 := range resolverRet {
		//   body
		//   append statement
		// }
		ret = append(ret, &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent(fmt.Sprintf("v%d", currentNestCount)),
			X:     rangeXExpr,
			Tok:   token.DEFINE,
			Body: &ast.BlockStmt{
				List: generateApplyQueryResponseFuncStmts(field, indexes, currentNestCount+1, nestCount, typePrefix, operationPrefix),
			},
		})

		// e.g. ret = append(ret, ret1)
		var lh ast.Expr = ast.NewIdent("ret")
		if currentNestCount > 0 {
			if currentNestCount > 1 {
				lh = ast.NewIdent(fmt.Sprintf("ret%d", currentNestCount-1))
			}

			ret = append(ret, &ast.AssignStmt{
				Lhs: []ast.Expr{lh},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("append"),
						Args: []ast.Expr{
							lh,
							resultLh,
						},
					},
				},
			})
		}
	}

	return ret
}

func generateApplyQueryResponseFuncStmts(field *schema.FieldDefinition, indexes *schema.Indexes, currentNestCount, nestCount int, typePrefix, operationPrefix string) []ast.Stmt {
	rootType := field.Type.GetRootType()

	_, isInterface := indexes.InterfaceIndex[string(rootType.Name)]
	_, isUnion := indexes.UnionIndex[string(rootType.Name)]
	if isInterface || isUnion {
		return generateApplyStmtForFragmentResponse(field, indexes, currentNestCount, nestCount, typePrefix, operationPrefix)
	}

	_, isObject := indexes.TypeIndex[string(rootType.Name)]
	if isObject {
		return generateApplyStmtForObjectQueryResponse(field, indexes.TypeIndex[string(rootType.Name)], indexes, currentNestCount, nestCount, typePrefix, operationPrefix)
	}

	return []ast.Stmt{}
}

func generateNestedArrayTypeForResponse(fieldType *schema.FieldType, indexes *schema.Indexes, isRoot bool) ast.Expr {
	_, isInterface := indexes.InterfaceIndex[string(fieldType.Name)]
	_, isUnion := indexes.UnionIndex[string(fieldType.Name)]
	if isInterface || isUnion {
		return ast.NewIdent(fmt.Sprintf("%sResponse", fieldType.Name))
	}

	if fieldType.IsList {
		return &ast.ArrayType{
			Elt: generateNestedArrayTypeForResponse(fieldType.ListType, indexes, false),
		}
	}

	return ast.NewIdent(fmt.Sprintf("%sResponse", fieldType.Name))
}

func generateObjectFragmentRangeStmtForQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationPrefix string) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)

	rootType := field.Type.GetRootType()

	stmts = append(stmts,
		&ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("node"),
						Sel: ast.NewIdent("FragmentType"),
					},
				},
				Op: token.NEQ,
				Y:  ast.NewIdent(fmt.Sprintf("%q", rootType.Name)),
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
									ast.NewIdent(fmt.Sprintf(`"node type mismatch: expected %s"`, rootType.Name)),
								},
							},
						},
					},
				},
			},
		})

	fragmentSwitch := generateApplyQueryResponseCaseStmtsForObjectFragment(indexes.TypeIndex[string(rootType.Name)], nestCount, string(field.Name), rootNestCount, indexes, operationPrefix)

	stmts = append(stmts, &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("child"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("node"),
			Sel: ast.NewIdent("Children"),
		},
		Tok: token.DEFINE,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.RangeStmt{
					Key:   ast.NewIdent("_"),
					Value: ast.NewIdent("fragmentChild"),
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("child"),
						Sel: ast.NewIdent("Children"),
					},
					Tok: token.DEFINE,
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.SwitchStmt{
								Tag: &ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("fragmentChild"),
											Sel: ast.NewIdent("Name"),
										},
									},
								},
								Body: &ast.BlockStmt{
									List: fragmentSwitch,
								},
							},
						},
					},
				},
			},
		},
	})

	return stmts
}

func generateObjectFragmentStmtForQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationPrefix string) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)

	stmts = append(stmts, &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("HasFragment"),
			},
		},
		Body: &ast.BlockStmt{
			List: generateObjectFragmentRangeStmtForQueryResponse(field, indexes, nestCount, rootNestCount, typePrefix, operationPrefix),
		},
	})

	stmts = append(stmts, generateApplySwitchStmtForQueryResponse(field, indexes, nestCount, rootNestCount, operationPrefix)...)

	return stmts
}

func generateApplySwitchStmtForQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, operationPrefix string) []ast.Stmt {
	rootType := field.Type.GetRootType()
	typeDefinition := indexes.TypeIndex[string(rootType.Name)]

	stmts := make([]ast.Stmt, 0)

	// e.g. switch string(child.Name) {
	//  case "$field":
	//    ret.$field = value
	// }
	switchStmt := &ast.SwitchStmt{
		Tag: &ast.CallExpr{
			Fun: ast.NewIdent("string"),
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("child"),
					Sel: ast.NewIdent("Name"),
				},
			},
		},
		Body: &ast.BlockStmt{
			List: generateApplyQueryResponseCaseStmts(typeDefinition, nestCount, string(field.Name), rootNestCount, indexes, operationPrefix),
		},
	}

	stmts = append(stmts, &ast.RangeStmt{
		Key:   ast.NewIdent("_"),
		Value: ast.NewIdent("child"),
		X: &ast.SelectorExpr{
			X:   ast.NewIdent("node"),
			Sel: ast.NewIdent("Children"),
		},
		Tok: token.DEFINE,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{switchStmt},
		},
	})

	// e.g. ret = append(ret, ret1)
	if nestCount > 0 && rootNestCount > 0 {
		var lh ast.Expr = ast.NewIdent("ret")
		if nestCount > 1 {
			lh = ast.NewIdent(fmt.Sprintf("ret%d", nestCount-1))
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				lh,
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						lh,
						ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
					},
				},
			},
		})
	}

	return stmts
}

func generateApplyStmtForFragmentResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationPrefix string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if rootNestCount == 0 {
		return []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{ast.NewIdent("ret")},
							Type:  ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)),
						},
					},
				},
			},
		}
	}

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("ret"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: generateApplyRetAssignStmtForFragmentNestResponse(field, indexes, nestCount+1, rootNestCount, typePrefix),
					},
					ast.NewIdent("0"),
					&ast.CallExpr{
						Fun: ast.NewIdent("len"),
						Args: []ast.Expr{
							ast.NewIdent("resolverRet"),
						},
					},
				},
			},
		},
	})
	ret = append(ret, generateApplyRangeStmtForFragmentResponse(field, indexes, nestCount, rootNestCount, typePrefix, operationPrefix)...)

	return ret
}

func generateApplyRetAssignStmtForFragmentNestResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix string) ast.Expr {
	if nestCount != rootNestCount {
		return &ast.ArrayType{
			Elt: generateApplyRetAssignStmtForFragmentNestResponse(field, indexes, nestCount+1, rootNestCount, typePrefix),
		}
	}

	return ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name))
}

func generateApplyRangeStmtForFragmentResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationName string) []ast.Stmt {
	rangeXExpr := ast.NewIdent("resolverRet")
	if nestCount > 0 {
		rangeXExpr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
	}

	if nestCount == rootNestCount {
		return generateApplySwitchStmtForFragmentQueryResponse(field, indexes, nestCount, rootNestCount, typePrefix, operationName)
	}

	var rangeRetAssignStmt, appendAssignStmt ast.Stmt = &ast.EmptyStmt{}, &ast.EmptyStmt{}
	if nestCount > 0 && nestCount < rootNestCount {
		rangeRetAssignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("make"),
					Args: []ast.Expr{
						&ast.ArrayType{
							Elt: generateApplyRetAssignStmtForFragmentNestResponse(field, indexes, nestCount+1, rootNestCount, typePrefix),
						},
						ast.NewIdent("0"),
						&ast.CallExpr{
							Fun: ast.NewIdent("len"),
							Args: []ast.Expr{
								rangeXExpr,
							},
						},
					},
				},
			},
		}

		var lh ast.Expr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
		if nestCount-1 == 0 {
			lh = ast.NewIdent("ret")
		} else {
			lh = ast.NewIdent(fmt.Sprintf("ret%d", nestCount-1))
		}

		appendAssignStmt = &ast.AssignStmt{
			Lhs: []ast.Expr{
				lh,
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						lh,
						ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
					},
				},
			},
		}
	}

	return []ast.Stmt{
		rangeRetAssignStmt,
		&ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent(fmt.Sprintf("v%d", nestCount)),
			X:     rangeXExpr,
			Tok:   token.DEFINE,
			Body: &ast.BlockStmt{
				List: generateApplyRangeStmtForFragmentResponse(field, indexes, nestCount+1, rootNestCount, typePrefix, operationName),
			},
		},
		appendAssignStmt,
	}
}

func getPossibleTypeDefinitions(field *schema.FieldDefinition, indexes *schema.Indexes) []*schema.TypeDefinition {
	fieldRootType := field.Type.GetRootType()
	ret := make(schema.TypeDefinitions, 0)
	if iface, ok := indexes.InterfaceIndex[string(fieldRootType.Name)]; ok {
		ret = append(ret, indexes.GetImplementedType(iface)...)
		return ret
	}

	if union, ok := indexes.UnionIndex[string(fieldRootType.Name)]; ok {
		for _, unionTypeName := range union.Types {
			if typeDef, ok := indexes.TypeIndex[string(unionTypeName)]; ok {
				ret = append(ret, typeDef)
			}
		}

		return ret
	}

	return ret
}

func generateApplyFragmentQueryResponseCaseStmtsForFragment(field *schema.FieldDefinition, typeDefinitions schema.TypeDefinitions, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationName string) []ast.Stmt {
	rangeSwitchStmts := make([]ast.Stmt, 0, len(typeDefinitions))
	for _, typeDefinition := range typeDefinitions {
		assignLh := ast.NewIdent("ret")
		if nestCount-1 > 0 {
			assignLh = ast.NewIdent(fmt.Sprintf("ret%s%d", typeDefinition.Name, rootNestCount))
		}

		switchCaseStmts := make([]ast.Stmt, 0)
		switchCaseStmts = append(switchCaseStmts, &ast.SwitchStmt{
			Tag: &ast.CallExpr{
				Fun: ast.NewIdent("string"),
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent("fragmentChild"),
						Sel: ast.NewIdent("Name"),
					},
				},
			},
			Body: &ast.BlockStmt{
				List: generateApplyQueryResponseCaseStmtsForFragment(typeDefinition, nestCount, string(field.Name), rootNestCount, indexes, operationName),
			},
		})

		rangeStmt := &ast.RangeStmt{
			Key:   ast.NewIdent("_"),
			Value: ast.NewIdent("fragmentChild"),
			X: &ast.SelectorExpr{
				X:   ast.NewIdent("child"),
				Sel: ast.NewIdent("Children"),
			},
			Tok: token.DEFINE,
			Body: &ast.BlockStmt{
				List: switchCaseStmts,
			},
		}

		rangeSwitchStmts = append(rangeSwitchStmts,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("ret%s%d", typeDefinition.Name, rootNestCount)),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", typeDefinition.Name)),
					},
				},
			})
		rangeSwitchStmts = append(rangeSwitchStmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("child"),
						Sel: ast.NewIdent("FragmentType"),
					},
				},
				Op: token.EQL,
				Y:  ast.NewIdent(fmt.Sprintf("%q", typeDefinition.Name)),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.TypeSwitchStmt{
						Assign: &ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(fmt.Sprintf("v%d", rootNestCount)),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X:    ast.NewIdent(fmt.Sprintf("v%d", rootNestCount-1)),
									Type: ast.NewIdent("type"),
								},
							},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent(typePrefix),
											Sel: ast.NewIdent(string(typeDefinition.Name)),
										},
									},
									Body: []ast.Stmt{
										rangeStmt,
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												assignLh,
											},
											Tok: token.ASSIGN,
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: ast.NewIdent("append"),
													Args: []ast.Expr{
														assignLh,
														ast.NewIdent(fmt.Sprintf("ret%s%d", typeDefinition.Name, rootNestCount)),
													},
												},
											},
										},
									},
								},
								&ast.CaseClause{
									List: []ast.Expr{
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(typePrefix),
												Sel: ast.NewIdent(string(typeDefinition.Name)),
											},
										},
									},
									Body: []ast.Stmt{
										rangeStmt,
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												assignLh,
											},
											Tok: token.ASSIGN,
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: ast.NewIdent("append"),
													Args: []ast.Expr{
														assignLh,
														ast.NewIdent(fmt.Sprintf("ret%s%d", typeDefinition.Name, rootNestCount)),
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.ExprStmt{
						X: &ast.BasicLit{
							Kind:  token.CONTINUE,
							Value: "continue",
						},
					},
				},
			},
		})
	}

	return rangeSwitchStmts
}

func generateApplySwitchStmtForFragmentQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix, operationName string) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)

	interfaceDefinition := indexes.InterfaceIndex[string(field.Type.GetRootType().Name)]
	typeDefinitions := getPossibleTypeDefinitions(field, indexes)

	stmt := &ast.IfStmt{
		Cond: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("node"),
				Sel: ast.NewIdent("HasFragment"),
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.RangeStmt{
					Key:   ast.NewIdent("_"),
					Value: ast.NewIdent("child"),
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("node"),
						Sel: ast.NewIdent("Children"),
					},
					Tok: token.DEFINE,
					Body: &ast.BlockStmt{
						List: generateApplyFragmentQueryResponseCaseStmtsForFragment(field, typeDefinitions, indexes, nestCount, rootNestCount, typePrefix, operationName),
					},
				},
			},
		},
	}

	if _, ok := indexes.InterfaceIndex[string(field.Type.GetRootType().Name)]; ok {
		stmt.Else = &ast.BlockStmt{
			List: generateInterfaceDefinitionApplyCaseStmts(interfaceDefinition, indexes, nestCount, rootNestCount, typePrefix),
		}
	}

	stmts = append(stmts, stmt)

	return stmts
}

func generateInterfaceDefinitionApplyCaseStmts(interfaceDefinition *schema.InterfaceDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix string) []ast.Stmt {
	caseStmts := make([]ast.Stmt, 0)
	assignLh := ast.NewIdent("ret")
	assignRh := ast.NewIdent("resolverRet")
	if nestCount > 0 {
		assignLh = ast.NewIdent(fmt.Sprintf("ret%d", rootNestCount))
		assignRh = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
	}

	for _, field := range interfaceDefinition.Fields {
		caseStmts = append(caseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(field.Name)),
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   assignLh,
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
									X:   ast.NewIdent(fmt.Sprintf("v%d", nestCount)),
									Sel: ast.NewIdent(toUpperCase(string(field.Name))),
								},
							},
						},
					},
				},
			},
		})
	}

	appendAssignLh := ast.NewIdent("ret")
	if nestCount-1 > 0 {
		appendAssignLh = ast.NewIdent(fmt.Sprintf("ret%d", nestCount-1))
	}

	typeCaseStmts := make([]ast.Stmt, 0)
	for _, typeDefinition := range indexes.GetImplementedType(interfaceDefinition) {
		rangeSwitchStmt := []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("ret%d", nestCount)),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CompositeLit{
						Type: ast.NewIdent(fmt.Sprintf("%sResponse", typeDefinition.Name)),
					},
				},
			},
			&ast.RangeStmt{
				Key:   ast.NewIdent("_"),
				Value: ast.NewIdent("child"),
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("node"),
					Sel: ast.NewIdent("Children"),
				},
				Tok: token.DEFINE,
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.SwitchStmt{
							Tag: &ast.CallExpr{
								Fun: ast.NewIdent("string"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("child"),
										Sel: ast.NewIdent("Name"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: caseStmts,
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					appendAssignLh,
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("append"),
						Args: []ast.Expr{
							appendAssignLh,
							assignLh,
						},
					},
				},
			},
		}

		typeCaseStmts = append(typeCaseStmts, &ast.CaseClause{
			List: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(typePrefix),
					Sel: ast.NewIdent(string(typeDefinition.Name)),
				},
			},
			Body: rangeSwitchStmt,
		}, &ast.CaseClause{
			List: []ast.Expr{
				&ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(typePrefix),
						Sel: ast.NewIdent(string(typeDefinition.Name)),
					},
				},
			},
			Body: rangeSwitchStmt,
		})
	}

	return []ast.Stmt{
		&ast.TypeSwitchStmt{
			Assign: &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(fmt.Sprintf("v%d", nestCount)),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X:    assignRh,
						Type: ast.NewIdent("type"),
					},
				},
			},
			Body: &ast.BlockStmt{
				List: typeCaseStmts,
			},
		},
	}
}

func generateApplyQueryResponseCaseStmts(typeDefinition *schema.TypeDefinition, nestCount int, fieldName string, rootNestCount int, indexes *schema.Indexes, operationName string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}
	var assignExpr ast.Expr = ast.NewIdent("ret")
	valueName := "resolverRet"
	if nestCount >= 0 && rootNestCount > 0 {
		valueName = fmt.Sprintf("v%d", nestCount-1)
		assignExpr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
	}
	for _, field := range typeDefinition.Fields {
		lhs := []ast.Expr{
			&ast.SelectorExpr{
				X:   assignExpr,
				Sel: ast.NewIdent(toUpperCase(string(field.Name))),
			},
		}

		caseBody := make([]ast.Stmt, 0)
		var rh ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(valueName),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}

		var arg ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(valueName),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}
		// if field.Type.Nullable {
		// 	arg = &ast.StarExpr{
		// 		X: arg,
		// 	}
		// }

		if field.Type.IsList && field.Type.Nullable {
			arg = &ast.StarExpr{
				X: arg,
			}
		}

		_, isScalar := indexes.ScalarIndex[string(field.Type.GetRootType().Name)]
		if field.Type.IsPrimitive() || isScalar {
			rh = generateNewNullableExpr(rh)
		} else {
			retName := fmt.Sprintf("ret%s", field.Name)
			fieldRetAssignStmt := &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(retName),
					ast.NewIdent("err"),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("r"),
							Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationName, string(field.Name))),
						},
						Args: []ast.Expr{
							arg,
							ast.NewIdent("child"),
						},
					},
				},
			}

			if _, ok := indexes.EnumIndex[string(field.Type.Name)]; ok {
				var arg ast.Expr = &ast.SelectorExpr{
					X:   ast.NewIdent(valueName),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				}
				if field.Type.Nullable {
					arg = &ast.StarExpr{
						X: arg,
					}
				}
				fieldRetAssignStmt = &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(retName),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("string"),
							Args: []ast.Expr{
								arg,
							},
						},
					},
				}
				caseBody = append(caseBody, fieldRetAssignStmt)
			} else {
				caseBody = append(caseBody, fieldRetAssignStmt)
				caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}))
			}

			// var elmExpr ast.Expr = ast.NewIdent(retName)
			// if field.Type.Nullable {
			// 	elmExpr = &ast.StarExpr{
			// 		X: ast.NewIdent(retName),
			// 	}
			// }

			rh = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("executor"),
					Sel: ast.NewIdent("NewNullable"),
				},
				Args: []ast.Expr{
					ast.NewIdent(retName),
				},
			}
		}
		caseBody = append(caseBody, &ast.AssignStmt{
			Lhs: lhs,
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				rh,
			},
		})

		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(field.Name)),
				},
			},
			Body: caseBody,
		})
	}

	ret = append(ret, &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"__typename"`,
			},
		},
		Body: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   assignExpr,
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
								Value: fmt.Sprintf(`"%s"`, string(typeDefinition.TypeName())),
							},
						},
					},
				},
			},
		},
	})

	return ret
}

func generateNewNullableExpr(expr ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   ast.NewIdent("executor"),
			Sel: ast.NewIdent("NewNullable"),
		},
		Args: []ast.Expr{expr},
	}
}

func generateApplyQueryResponseCaseStmtsForObjectFragment(typeDefinition *schema.TypeDefinition, nestCount int, fieldName string, rootNestCount int, indexes *schema.Indexes, operationName string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}

	var assignExpr ast.Expr = ast.NewIdent("ret")
	valueName := "resolverRet"
	if nestCount >= 0 && rootNestCount > 0 {
		valueName = fmt.Sprintf("v%d", nestCount-1)
		assignExpr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
	}

	for _, field := range typeDefinition.Fields {
		lhs := []ast.Expr{
			&ast.SelectorExpr{
				X:   assignExpr,
				Sel: ast.NewIdent(toUpperCase(string(field.Name))),
			},
		}

		caseBody := make([]ast.Stmt, 0)
		var rh ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(valueName),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}

		if !field.Type.Nullable {
			_, isScalar := indexes.ScalarIndex[string(field.Type.GetRootType().Name)]
			_, isEnum := indexes.EnumIndex[string(field.Type.GetRootType().Name)]
			if field.Type.GetRootType().IsPrimitive() || isScalar || isEnum {
				rh = generateNewNullableExpr(rh)
			} else {
				retName := fmt.Sprintf("ret%s", field.Name)
				var valueExpr ast.Expr = &ast.SelectorExpr{
					X:   ast.NewIdent(valueName),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				}

				fieldRetAssignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(retName),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationName, string(field.Name))),
							},
							Args: []ast.Expr{
								valueExpr,
								ast.NewIdent("child"),
							},
						},
					},
				}
				caseBody = append(caseBody, fieldRetAssignStmt)
				if rootNestCount > 0 || (rootNestCount == 0 && field.Type.Nullable) {
					caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
						ast.NewIdent("nil"),
					}))
				} else {
					caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
						ast.NewIdent("ret"),
					}))
				}

				if field.Type.IsList {
					if field.Type.GetRootType().Nullable {
						rh = &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(retName),
							},
						}
					} else {
						rh = &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(retName),
							},
						}
					}
				} else {
					rh = &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(retName),
						},
					}
				}
			}
			caseBody = append(caseBody, &ast.AssignStmt{
				Lhs: lhs,
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					rh,
				},
			})
		} else {
			_, isScalar := indexes.ScalarIndex[string(field.Type.GetRootType().Name)]
			_, isEnum := indexes.EnumIndex[string(field.Type.GetRootType().Name)]
			if !field.Type.GetRootType().IsPrimitive() && !isScalar && !isEnum {
				var firstArg ast.Expr = &ast.SelectorExpr{
					X:   ast.NewIdent(valueName),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				}

				if field.Type.IsList {
					firstArg = &ast.StarExpr{
						X: firstArg,
					}
				}

				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", field.Name)),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationName,
									string(field.Name))),
							},
							Args: []ast.Expr{
								firstArg,
								ast.NewIdent("child"),
							},
						},
					},
				})

				caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}))

				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   assignExpr,
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						generateNewNullableExpr(ast.NewIdent(fmt.Sprintf("ret%s", field.Name))),
					},
				})
			} else {
				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: lhs,
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						generateNewNullableExpr(rh),
					},
				})
			}
		}

		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(field.Name)),
				},
			},
			Body: caseBody,
		})
	}

	ret = append(ret, &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"__typename"`,
			},
		},
		Body: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   assignExpr,
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
								Value: fmt.Sprintf(`"%s"`, string(typeDefinition.TypeName())),
							},
						},
					},
				},
			},
		},
	})

	return ret
}

func generateApplyQueryResponseCaseStmtsForFragment(typeDefinition *schema.TypeDefinition, nestCount int, fieldName string, rootNestCount int, indexes *schema.Indexes, operationName string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}

	var assignExpr ast.Expr = ast.NewIdent("ret")
	valueName := "resolverRet"
	if nestCount >= 0 && rootNestCount > 0 {
		valueName = fmt.Sprintf("v%d", nestCount)
		assignExpr = ast.NewIdent(fmt.Sprintf("ret%s%d", typeDefinition.Name, nestCount))
	}

	for _, field := range typeDefinition.Fields {
		lhs := []ast.Expr{
			&ast.SelectorExpr{
				X:   assignExpr,
				Sel: ast.NewIdent(toUpperCase(string(field.Name))),
			},
		}

		caseBody := make([]ast.Stmt, 0)
		var rh ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent(valueName),
			Sel: ast.NewIdent(toUpperCase(string(field.Name))),
		}

		if !field.Type.Nullable {
			_, isScalar := indexes.ScalarIndex[string(field.Type.GetRootType().Name)]
			_, isEnum := indexes.EnumIndex[string(field.Type.GetRootType().Name)]
			if field.Type.GetRootType().IsPrimitive() || isScalar || isEnum {
				rh = generateNewNullableExpr(rh)
			} else {
				retName := fmt.Sprintf("ret%s", field.Name)
				var valueExpr ast.Expr = &ast.SelectorExpr{
					X:   ast.NewIdent(valueName),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				}

				fieldRetAssignStmt := &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(retName),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationName, string(field.Name))),
							},
							Args: []ast.Expr{
								valueExpr,
								ast.NewIdent("child"),
							},
						},
					},
				}
				caseBody = append(caseBody, fieldRetAssignStmt)
				if rootNestCount > 0 || (rootNestCount == 0 && field.Type.Nullable) {
					caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
						ast.NewIdent("nil"),
					}))
				} else {
					caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
						&ast.CompositeLit{
							Type: ast.NewIdent(string(typeDefinition.TypeName()) + "Response"),
						},
					}))
				}

				if field.Type.IsList {
					if field.Type.GetRootType().Nullable {
						rh = &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(retName),
							},
						}
					} else {
						rh = &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("executor"),
								Sel: ast.NewIdent("NewNullable"),
							},
							Args: []ast.Expr{
								ast.NewIdent(retName),
							},
						}
					}
				} else {
					rh = &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("executor"),
							Sel: ast.NewIdent("NewNullable"),
						},
						Args: []ast.Expr{
							ast.NewIdent(retName),
						},
					}
				}
			}
			caseBody = append(caseBody, &ast.AssignStmt{
				Lhs: lhs,
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					rh,
				},
			})
		} else {
			_, isScalar := indexes.ScalarIndex[string(field.Type.GetRootType().Name)]
			_, isEnum := indexes.EnumIndex[string(field.Type.GetRootType().Name)]
			if !field.Type.GetRootType().IsPrimitive() && !isScalar && !isEnum {
				var valueExpr ast.Expr = &ast.SelectorExpr{
					X:   ast.NewIdent(valueName),
					Sel: ast.NewIdent(toUpperCase(string(field.Name))),
				}
				if field.Type.IsList {
					valueExpr = &ast.StarExpr{
						X: valueExpr,
					}
				}

				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("ret%s", field.Name)),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", operationName,
									string(field.Name))),
							},
							Args: []ast.Expr{
								valueExpr,
								ast.NewIdent("child"),
							},
						},
					},
				})

				caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}))

				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   assignExpr,
							Sel: ast.NewIdent(toUpperCase(string(field.Name))),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						generateNewNullableExpr(ast.NewIdent(fmt.Sprintf("ret%s", field.Name))),
					},
				})
			} else {
				caseBody = append(caseBody, &ast.AssignStmt{
					Lhs: lhs,
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						generateNewNullableExpr(rh),
					},
				})
			}
		}

		ret = append(ret, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s"`, string(field.Name)),
				},
			},
			Body: caseBody,
		})
	}

	ret = append(ret, &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: `"__typename"`,
			},
		},
		Body: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					&ast.SelectorExpr{
						X:   assignExpr,
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
								Value: fmt.Sprintf(`"%s"`, string(typeDefinition.TypeName())),
							},
						},
					},
				},
			},
		},
	})

	return ret
}

func generateModelFieldForResponse(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldTypeExpr := generateExprForResponse(f.Type)

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeExpr,
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s,omitempty\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}
