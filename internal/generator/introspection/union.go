package introspection

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateUnionPossibleAssignStmts(unionDefinition *schema.UnionDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)
	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("possibleTypes"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: ast.NewIdent("__Type"),
					},
					ast.NewIdent("0"),
					&ast.BasicLit{
						Kind:  token.INT,
						Value: fmt.Sprintf("%d", len(unionDefinition.Types)),
					},
				},
			},
		},
	})

	for _, unionType := range unionDefinition.Types {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sPossibleType", unionType)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", unionType)),
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
				&ast.SelectorExpr{
					X:   ast.NewIdent("ret"),
					Sel: ast.NewIdent("PossibleTypes"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						Sel: ast.NewIdent("executor"),
						X:   ast.NewIdent("NewNullable"),
					},
					Args: []ast.Expr{
						ast.NewIdent(fmt.Sprintf("%sPossibleType", unionType)),
					},
				},
			},
		},
		)
	}

	return ret
}

func GenerateUnionTypeCaseStmts(unionDefinition *schema.UnionDefinition) []ast.Stmt {
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
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, string(unionDefinition.Name)),
								},
							},
						},
					},
				},
			},
		}, &ast.CaseClause{
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
						ast.NewIdent("__TypeKind_UNION"),
					},
				},
			},
		}, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fields"`,
				},
			},
			Body: []ast.Stmt{
				&ast.AssignStmt{
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
								&ast.CompositeLit{
									Type: &ast.ArrayType{
										Elt: ast.NewIdent("__Field"),
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
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"interfaces"`,
				},
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
										Elt: ast.NewIdent("__Type"),
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
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"possibleTypes"`,
				},
			},
			Body: generateUnionPossibleAssignStmts(unionDefinition),
		}, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"inputFields"`,
				},
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
										Elt: ast.NewIdent("__InputValue"),
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
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"enumValues"`,
				},
			},
			Body: []ast.Stmt{
				generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}),
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
										Elt: ast.NewIdent("__EnumValue"),
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
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"ofType"`,
				},
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
		},
	}
}
