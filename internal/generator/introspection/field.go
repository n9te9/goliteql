package introspection

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateFieldIsDeprecatedAssignStmt(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent("ret"),
				Sel: ast.NewIdent("IsDeprecated"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%t", fieldDefinition.IsDeprecated()),
			},
		},
	}
}

func generateFieldDeprecationReasonAssignStmt(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	deprecationReason := "nil"
	if fieldDefinition.IsDeprecated() {
		deprecationReason = fieldDefinition.DeprecatedReason()
	}

	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(fmt.Sprintf("field%s", fieldDefinition.Name)),
				Sel: ast.NewIdent("DeprecationReason"),
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
						Value: fmt.Sprintf(`"%s"`, deprecationReason),
					},
				},
			},
		},
	}
}

func generateFieldNameAssignStmt(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(fmt.Sprintf("field%s", fieldDefinition.Name)),
				Sel: ast.NewIdent("Name"),
			},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, string(fieldDefinition.Name)),
			},
		},
	}
}

func generateFieldDescriptionAssignStmt(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(fmt.Sprintf("field%s", fieldDefinition.Name)),
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
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "nil",
					},
				},
			},
		},
	}
}

func generateFieldArgsAssignStmt(fieldDefinition *schema.FieldDefinition) ast.Stmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.SelectorExpr{
				X:   ast.NewIdent(fmt.Sprintf("field%s", fieldDefinition.Name)),
				Sel: ast.NewIdent("Args"),
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
	}
}

func generateFieldTypeAssignStmt(fieldDefinition *schema.FieldDefinition) []ast.Stmt {
	return []ast.Stmt{
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent("t"),
				ast.NewIdent("err"),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("r"),
						Sel: ast.NewIdent(fmt.Sprintf("__schema__%s__type", fieldDefinition.Type.Name)),
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
			ast.NewIdent("ret"),
		}),
		&ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent(fmt.Sprintf("field%s", fieldDefinition.Name)),
					Sel: ast.NewIdent("Type"),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{ast.NewIdent("t")},
		},
	}
}

func GenerateFieldsCaseStmts(fieldDefinition schema.FieldDefinitions) []ast.Stmt {
	nameAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	descriptionAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	argsAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	typeAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	isDeprecatedAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	deprecationReasonAssignStmts := make([]ast.Stmt, 0, len(fieldDefinition))
	for _, field := range fieldDefinition {
		nameAssignStmts = append(nameAssignStmts, generateFieldNameAssignStmt(field))
		descriptionAssignStmts = append(descriptionAssignStmts, generateFieldDescriptionAssignStmt(field))
		argsAssignStmts = append(argsAssignStmts, generateFieldArgsAssignStmt(field))
		typeAssignStmts = append(typeAssignStmts, generateFieldTypeAssignStmt(field)...)
		isDeprecatedAssignStmts = append(isDeprecatedAssignStmts, generateFieldIsDeprecatedAssignStmt(field))
		deprecationReasonAssignStmts = append(deprecationReasonAssignStmts, generateFieldDeprecationReasonAssignStmt(field))
	}

	return []ast.Stmt{
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"name"`,
				},
			},
			Body: nameAssignStmts,
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"description"`,
				},
			},
			Body: descriptionAssignStmts,
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"args"`,
				},
			},
			Body: argsAssignStmts,
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"type"`,
				},
			},
			Body: typeAssignStmts,
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"isDeprecated"`,
				},
			},
			Body: isDeprecatedAssignStmts,
		},
		&ast.CaseClause{
			List: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: `"deprecationReason"`,
				},
			},
			Body: deprecationReasonAssignStmts,
		},
	}
}
