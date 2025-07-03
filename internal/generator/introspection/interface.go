package introspection

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateInterfaceFields(extentions []*schema.InterfaceDefinition) []ast.Stmt {
	ret := make([]ast.Stmt, 0)

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("fields"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: ast.NewIdent("make"),
				Args: []ast.Expr{
					&ast.ArrayType{
						Elt: ast.NewIdent("__Field"),
					},
					ast.NewIdent("0"),
					&ast.BasicLit{
						Kind:  token.INT,
						Value: string(len(extentions)),
					},
				},
			},
		},
	})

	ret = append(ret, &ast.AssignStmt{
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
	}, generateReturnErrorHandlingStmt([]ast.Expr{
		ast.NewIdent("ret"),
	}))

	for _, ext := range extentions {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sFields", ext.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", ext.Name)),
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
						ast.NewIdent(fmt.Sprintf("%sFields", ext.Name)),
					},
				},
			},
		})
	}

	ret = append(ret, &ast.AssignStmt{
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

	return ret
}

func generateInterfacePossibleTypeStmts(interfaceDefinition *schema.InterfaceDefinition, indexes *schema.Indexes) []ast.Stmt {
	implementedTypes := indexes.GetImplementedType(interfaceDefinition)

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
						Value: fmt.Sprintf("%d", len(implementedTypes)),
					},
				},
			},
		},
	})

	for _, implementedType := range implementedTypes {
		ret = append(ret, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(fmt.Sprintf("%sPossibleTypes", implementedType.Name)),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", implementedType.Name)),
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
				ast.NewIdent("possibleTypes"),
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: ast.NewIdent("append"),
					Args: []ast.Expr{
						ast.NewIdent("possibleTypes"),
						ast.NewIdent(fmt.Sprintf("%sPossibleTypes", implementedType.Name)),
					},
				},
			},
		})
	}

	ret = append(ret, &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("ret"),
				Sel: ast.NewIdent("PossibleTypes"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			ast.NewIdent("possibleTypes"),
		},
	})

	return ret
}

func GenerateInterfaceTypeCaseStmts(interfaceDefinition *schema.InterfaceDefinition, indexes *schema.Indexes) []ast.Stmt {
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
								ast.NewIdent(string(interfaceDefinition.Name)),
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
						ast.NewIdent("__TypeKind_INTERFACE"),
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
				},
				generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}),
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("fields"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("r"),
								Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__fields", interfaceDefinition.Name)),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("child"),
								ast.NewIdent("variables"),
								&ast.StarExpr{
									X: ast.NewIdent("includeDeprecated"),
								},
							},
						},
					},
				},
				generateReturnErrorHandlingStmt([]ast.Expr{
					ast.NewIdent("ret"),
				}),
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
								ast.NewIdent("fields"),
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
			Body: generateInterfaceFields(interfaceDefinition.Extentions),
		}, &ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"possibleTypes"`,
				},
			},
			Body: generateInterfacePossibleTypeStmts(interfaceDefinition, indexes),
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
