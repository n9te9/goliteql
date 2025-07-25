package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

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

	if fieldType.Nullable {
		return &ast.StarExpr{
			X: ast.NewIdent(fmt.Sprintf("%sResponse", fieldType.Name)),
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
								&ast.SelectorExpr{
									X:   ast.NewIdent(valueName),
									Sel: ast.NewIdent(toUpperCase(string(field.Name))),
								},
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
								&ast.SelectorExpr{
									X:   ast.NewIdent(valueName),
									Sel: ast.NewIdent(toUpperCase(string(field.Name))),
								},
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
