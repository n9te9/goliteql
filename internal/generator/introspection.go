package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/n9te9/goliteql/schema"
)

func generateIntrospectionModelAST(types []*schema.TypeDefinition) []ast.Decl {
	decls := make([]ast.Decl, 0, len(types))

	for _, t := range types {
		if t.IsIntrospection() {
			if t.PrimitiveTypeName != nil {
				decls = append(decls, &ast.GenDecl{
					Tok: token.TYPE,
					Specs: []ast.Spec{
						&ast.TypeSpec{
							Name: ast.NewIdent(string(t.Name)),
							Type: ast.NewIdent(string(t.PrimitiveTypeName)),
						},
					},
				})

				continue
			}

			decls = append(decls, &ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(string(t.Name)),
						Type: &ast.StructType{
							Fields: generateModelField(t.Fields),
						},
					},
				},
			})
		}
	}

	return decls
}

func generateModelFieldCaseAST(field *schema.FieldDefinition) ast.Stmt {
	return &ast.CaseClause{
		List: []ast.Expr{
			&ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, string(field.Name)),
			},
		},
		Body: []ast.Stmt{},
	}
}

func generateModelFieldCaseASTs(fields []*schema.FieldDefinition) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))
	for _, f := range fields {
		stmts = append(stmts, generateModelFieldCaseAST(f))
	}

	return stmts
}

func generateIntrospectionSchemaQueryAST(s *schema.Schema) ast.Decl {
	stmts := make([]ast.Stmt, 0)
	for _, t := range s.Types {
		if string(t.Name) == "__Schema" {
			stmts = append(stmts, generateModelFieldCaseASTs(t.Fields)...)
		}
	}

	body := make([]ast.Stmt, 0, len(s.Types))
	body = append(body, &ast.SwitchStmt{
		Tag: ast.NewIdent("string(node.Name)"),
		Body: &ast.BlockStmt{
			List: stmts,
		},
	})

	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("r"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("resolver"),
					},
				},
			},
		},
		Name: ast.NewIdent("__schema"),
		Type: &ast.FuncType{
			Params: generateOperationExecutorArgs(),
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
