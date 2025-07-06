package introspection

import (
	"go/ast"
	"go/token"
)

func generateReturnErrorHandlingStmt(prefixReturnExpr []ast.Expr) ast.Stmt {
	prefixReturnExpr = append(prefixReturnExpr, ast.NewIdent("err"))
	return &ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X:  ast.NewIdent("err"),
			Op: token.NEQ,
			Y:  &ast.BasicLit{Kind: token.STRING, Value: "nil"},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: prefixReturnExpr,
				},
			},
		},
	}
}
