package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateApplyQueryResponseFuncDeclFromField(field *schema.FieldDefinition, indexes *schema.Indexes, typePrefix, operationPrefix string) ast.Decl {
	nestCount := getNestCount(field.Type, 0)

	stmts := generateApplyQueryResponseFuncStmts(field, indexes, 0, nestCount, typePrefix)
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("ret"),
			ast.NewIdent("nil"),
		},
	})

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
						Type: generateNestedArrayTypeForResponse(ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)), 0, nestCount, field.Type.GetRootType().Nullable),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
}

func generateApplyQueryResponseFuncStmts(field *schema.FieldDefinition, indexes *schema.Indexes, currentNestCount, nestCount int, typePrefix string) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	var resultLh ast.Expr = ast.NewIdent("ret")
	if nestCount > 0 && currentNestCount > 0 {
		resultLh = ast.NewIdent(fmt.Sprintf("ret%d", currentNestCount))
	}

	rootType := field.Type.GetRootType()

	_, isInterface := indexes.InterfaceIndex[string(rootType.Name)]
	_, isUnion := indexes.UnionIndex[string(rootType.Name)]
	if isInterface || isUnion {
		return generateApplyStmtForFragmentResponse(field, indexes, currentNestCount, nestCount, typePrefix)
	}

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

		return append(ret, generateApplySwitchStmtForQueryResponse(field, indexes, currentNestCount, nestCount)...)
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
						generateNestedArrayTypeForResponse(ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)), currentNestCount, nestCount, field.Type.GetRootType().Nullable),
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
				List: generateApplyQueryResponseFuncStmts(field, indexes, currentNestCount+1, nestCount, typePrefix),
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

func generateNestedArrayTypeForResponse(typeExpr ast.Expr, currentNestCount, nestCount int, isNullable bool) ast.Expr {
	if currentNestCount >= nestCount {
		if isNullable {
			return &ast.SelectorExpr{
				X:   ast.NewIdent("executor"),
				Sel: ast.NewIdent("Nullable"),
			}
		}

		return typeExpr
	}

	return &ast.ArrayType{
		Elt: generateNestedArrayTypeForResponse(typeExpr, currentNestCount+1, nestCount, isNullable),
	}
}

func generateApplySwitchStmtForQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int) []ast.Stmt {
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
			List: generateApplyQueryResponseCaseStmts(typeDefinition, nestCount, string(field.Name), rootNestCount),
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

func generateApplyStmtForFragmentResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix string) []ast.Stmt {
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
	ret = append(ret, generateApplyRangeStmtForFragmentResponse(field, indexes, nestCount, rootNestCount, typePrefix)...)

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

func generateApplyRangeStmtForFragmentResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix string) []ast.Stmt {
	rangeXExpr := ast.NewIdent("resolverRet")
	if nestCount > 0 {
		rangeXExpr = ast.NewIdent(fmt.Sprintf("v%d", nestCount-1))
	}

	if nestCount == rootNestCount {
		return generateApplySwitchStmtForFragmentQueryResponse(field, indexes, nestCount, rootNestCount, typePrefix)
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
				List: generateApplyRangeStmtForFragmentResponse(field, indexes, nestCount+1, rootNestCount, typePrefix),
			},
		},
		appendAssignStmt,
	}
}

func generateApplySwitchStmtForFragmentQueryResponse(field *schema.FieldDefinition, indexes *schema.Indexes, nestCount, rootNestCount int, typePrefix string) []ast.Stmt {
	stmts := make([]ast.Stmt, 0)

	interfaceDefinition := indexes.InterfaceIndex[string(field.Type.GetRootType().Name)]
	typeDefinitions := indexes.GetImplementedType(interfaceDefinition)

	for _, typeDefinition := range typeDefinitions {
		rangeStmt := &ast.RangeStmt{
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
							List: generateApplyQueryResponseCaseStmtsForFragment(typeDefinition, nestCount, string(field.Name), rootNestCount),
						},
					},
				},
			},
		}

		stmts = append(stmts, &ast.IfStmt{
			Cond: &ast.BinaryExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("node"),
					Sel: ast.NewIdent("Type"),
				},
				Op: token.EQL,
				Y:  ast.NewIdent(fmt.Sprintf("%q", typeDefinition.Name)),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("ret%d", rootNestCount)),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CompositeLit{
								Type: ast.NewIdent(fmt.Sprintf("%sResponse", typeDefinition.Name)),
							},
						},
					},
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

func generateApplyQueryResponseCaseStmts(typeDefinition *schema.TypeDefinition, nestCount int, fieldName string, rootNestCount int) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}

	for _, field := range typeDefinition.Fields {
		var assignExpr ast.Expr = ast.NewIdent("ret")
		valueName := "resolverRet"
		if nestCount >= 0 && rootNestCount > 0 {
			valueName = fmt.Sprintf("v%d", nestCount-1)
			assignExpr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
		}
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
			if field.Type.IsBoolean() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsString() || field.Type.IsID() {
				rh = generateStringPointerExpr(rh)
			} else if field.Type.IsInt() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsFloat() {
				rh = generateFloatPointerExpr(rh)
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
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", fieldName, string(field.Name))),
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

				// var elmExpr ast.Expr = ast.NewIdent(retName)
				// if field.Type.Nullable {
				// 	elmExpr = &ast.StarExpr{
				// 		X: ast.NewIdent(retName),
				// 	}
				// }

				if field.Type.IsList {
					if field.Type.GetRootType().Nullable {
						rh = &ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(retName),
						}
					} else {
						rh = &ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(retName),
						}
					}
				} else {
					rh = &ast.UnaryExpr{
						Op: token.AND,
						X:  ast.NewIdent(retName),
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
			if !field.Type.IsPrimitive() {
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
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", fieldName,
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
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(fmt.Sprintf("ret%s", field.Name)),
						},
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

	return ret
}

func generateApplyQueryResponseCaseStmtsForFragment(typeDefinition *schema.TypeDefinition, nestCount int, fieldName string, rootNestCount int) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if typeDefinition == nil {
		return ret
	}

	for _, field := range typeDefinition.Fields {
		var assignExpr ast.Expr = ast.NewIdent("ret")
		valueName := "resolverRet"
		if nestCount >= 0 && rootNestCount > 0 {
			valueName = fmt.Sprintf("v%d", nestCount)
			assignExpr = ast.NewIdent(fmt.Sprintf("ret%d", nestCount))
		}
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
			if field.Type.IsBoolean() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsString() || field.Type.IsID() {
				rh = generateStringPointerExpr(rh)
			} else if field.Type.IsInt() {
				rh = generateIntPointerExpr(rh)
			} else if field.Type.IsFloat() {
				rh = generateFloatPointerExpr(rh)
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
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", fieldName, string(field.Name))),
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

				// var elmExpr ast.Expr = ast.NewIdent(retName)
				// if field.Type.Nullable {
				// 	elmExpr = &ast.StarExpr{
				// 		X: ast.NewIdent(retName),
				// 	}
				// }

				if field.Type.IsList {
					if field.Type.GetRootType().Nullable {
						rh = &ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(retName),
						}
					} else {
						rh = &ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(retName),
						}
					}
				} else {
					rh = &ast.UnaryExpr{
						Op: token.AND,
						X:  ast.NewIdent(retName),
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
			if !field.Type.IsPrimitive() {
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
								Sel: ast.NewIdent(fmt.Sprintf("apply%s%sQueryResponse", fieldName,
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
						&ast.UnaryExpr{
							Op: token.AND,
							X:  ast.NewIdent(fmt.Sprintf("ret%s", field.Name)),
						},
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

	return ret
}
