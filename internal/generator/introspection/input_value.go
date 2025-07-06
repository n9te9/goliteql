package introspection

import (
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func GenerateInputValuesCaseStmts(fieldDefinition schema.FieldDefinitions) []ast.Stmt {
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
