package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateApplyQueryResponseFuncDeclFromField(field *schema.FieldDefinition, indexes *schema.Indexes, typePrefix, operationPrefix string) ast.Decl {
	nestCount := getNestCount(field.Type, 0)

	stmts := generateApplyQueryResponseFuncStmts(field, indexes, 0, nestCount)
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
						Type: generateTypeExprFromFieldType(typePrefix, field.Type),
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
						Type: generateTypeExprFromFieldType("", generateFieldTypeForResponse(field.Type)),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
}

func generateApplyQueryResponseFuncStmts(field *schema.FieldDefinition, indexes *schema.Indexes, currentNestCount, nestCount int) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	var resultLh ast.Expr = ast.NewIdent("ret")
	if nestCount > 0 && currentNestCount > 0 {
		resultLh = ast.NewIdent(fmt.Sprintf("ret%d", currentNestCount))
	}

	if currentNestCount == nestCount {
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
						generateNestedArrayTypeForResponse(ast.NewIdent(fmt.Sprintf("%sResponse", field.Type.GetRootType().Name)), currentNestCount, nestCount),
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
				List: generateApplyQueryResponseFuncStmts(field, indexes, currentNestCount+1, nestCount),
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

func generateNestedArrayTypeForResponse(typeExpr ast.Expr, currentNestCount, nestCount int) ast.Expr {
	if currentNestCount >= nestCount {
		return &ast.StarExpr{
			X: typeExpr,
		}
	}

	return &ast.ArrayType{
		Elt: generateNestedArrayTypeForResponse(typeExpr, currentNestCount+1, nestCount),
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
				caseBody = append(caseBody, generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("nil"),
				}))

				var elmExpr ast.Expr = ast.NewIdent(retName)
				if !field.Type.Nullable {
					elmExpr = &ast.StarExpr{
						X: ast.NewIdent(retName),
					}
				}
				rh = generateObjectPointerExpr(string(field.Type.Name)+"Response", elmExpr)
			}
		} else {
			rh = &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					Sel: ast.NewIdent("NewNullable"),
					X:   ast.NewIdent("executor"),
				},
				Args: []ast.Expr{
					rh,
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

	return ret
}
