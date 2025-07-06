package introspection

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateOperationFieldBodyStmts(fieldDefinitions schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fieldDefinitions))

	if fieldDefinitions.HasDeprecatedDirective() {
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("includeDeprecated"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent("extract__fieldsArgs"),
					},
					Args: []ast.Expr{
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		})
	}

	for _, fieldDefinition := range fieldDefinitions {
		stmts = append(stmts, generateOperationFieldAssignStmt(fieldDefinition)...)
	}

	stmts = append(stmts, &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("ret"),
				Sel: ast.NewIdent("Fields"),
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
					ast.NewIdent("fields"),
				},
			},
		},
	})

	return stmts
}

func generateOperationFieldAssignStmt(fieldDefinition *schema.FieldDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	if fieldDefinition.IsDeprecated() {
		ret = append(ret, &ast.IfStmt{
			Cond: &ast.StarExpr{
				X: ast.NewIdent("includeDeprecated"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent(fmt.Sprintf("%sfield", fieldDefinition.Name)),
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("r"),
									Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(fieldDefinition.Name))),
								},
								Args: []ast.Expr{
									ast.NewIdent("ctx"),
									ast.NewIdent("child"),
									ast.NewIdent("variables"),
								},
							},
						},
					},
					generateReturnErrorHandlingStmt([]ast.Expr{
						&ast.CompositeLit{
							Type: ast.NewIdent("__Type"),
							Elts: []ast.Expr{},
						},
					}), &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("fields"),
						},
						Tok: token.ASSIGN,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("append"),
								Args: []ast.Expr{
									ast.NewIdent("fields"),
									ast.NewIdent(fmt.Sprintf("%sfield", fieldDefinition.Name)),
								},
							},
						},
					},
				},
			},
		})
	} else {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sfield", fieldDefinition.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", string(fieldDefinition.Name))),
					},
					Args: []ast.Expr{
						ast.NewIdent("ctx"),
						ast.NewIdent("child"),
						ast.NewIdent("variables"),
					},
				},
			},
		}, generateReturnErrorHandlingStmt([]ast.Expr{
			ast.NewIdent("ret"),
		}), &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("fields"),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent("fields"),
						ast.NewIdent(fmt.Sprintf("%sfield", fieldDefinition.Name)),
					},
				},
			},
		})
	}

	return ret
}

func GenerateOperationCaseStmts(operationName string, operationDefinition *schema.OperationDefinition) []ast.Stmt {
	return []ast.Stmt{
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"name"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Name"),
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
								ast.NewIdent(fmt.Sprintf(`"%s"`, string(operationName))),
							},
						},
					},
				},
			},
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"kind"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Kind"),
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						ast.NewIdent("__TypeKind_OBJECT"),
					},
				},
			},
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"description"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Description"),
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
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fields"`,
				},
			},
			Body: generateOperationFieldBodyStmts(operationDefinition.Fields),
		},
		&ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"interfaces"`),
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("Interfaces"),
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
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.Ident{Name: "__Type"},
									},
								},
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"inputFields"`),
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("InputFields"),
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
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.Ident{Name: "__InputValue"},
									},
								},
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"enumValues"`),
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("EnumValues"),
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
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.Ident{Name: "__EnumValue"},
									},
								},
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"possibleTypes"`),
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("PossibleTypes"),
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
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: &ast.Ident{Name: "__Type"},
									},
									Elts: []ast.Expr{},
								},
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"ofType"`),
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.SelectorExpr{
							X:   ast.NewIdent("ret"),
							Sel: ast.NewIdent("OfType"),
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
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				ast.NewIdent(`"specifiedByURL"`),
			},
			Body: []ast.Stmt{},
		},
	}
}
