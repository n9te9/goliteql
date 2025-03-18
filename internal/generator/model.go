package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-parser/schema"
)

func generateModelImport() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: `"encoding/json"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: `"fmt"`,
				},
			},
		},
	}
}

func generateModelField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldType := GraphQLType(f.Type.Name)
		var fieldTypeIdent *ast.Ident
		if fieldType.IsPrimitive() {
			fieldTypeIdent = golangType(f.Type, fieldType, "")
		} else {
			fieldTypeIdent = golangType(f.Type, fieldType, "")
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
			Tag: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateTypeModelUnmarshalJSON(t *schema.TypeDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("UnmarshalJSON"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "t",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: string(t.Name),
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "data",
							},
						},
						Type: &ast.Ident{
							Name: "[]byte",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{

		},
	}
}

func generateModelMapperField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	starExpr := func(fieldType *schema.FieldType, graphQLType GraphQLType, modelPackagePath string) *ast.StarExpr {
		if fieldType.IsList {
			return &ast.StarExpr{
				X: &ast.Ident{
					Name: "[]" + golangType(fieldType.ListType, GraphQLType(fieldType.ListType.Name), modelPackagePath).Name,
				},
			}
		}

		if graphQLType.IsPrimitive() {
			if fieldType.Nullable {
				return &ast.StarExpr{
					X: &ast.Ident{
						Name: graphQLType.golangType(),
					},
				}
			}
		}

		return &ast.StarExpr{
			X: &ast.Ident{
				Name: graphQLType.golangType(),
			},
		}
	}

	for _, f := range field {
		fieldType := GraphQLType(f.Type.Name)
		var fieldTypeIdent *ast.StarExpr
		if fieldType.IsPrimitive() {
			fieldTypeIdent = starExpr(f.Type, fieldType, "")
		} else {
			fieldTypeIdent = starExpr(f.Type, fieldType, "")
		}

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
			Tag: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateInputModelUnmarshalJSON(t *schema.InputDefinition) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("UnmarshalJSON"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "t",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: string(t.Name),
						},
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "data",
							},
						},
						Type: &ast.Ident{
							Name: "[]byte",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: append(append([]ast.Stmt{}, generateUnmarshalJSONBody(t.Fields)...), generateMappingSchemaValidation(t)...),
		},
	}
}

func generateUnmarshalJSONBody(fields schema.FieldDefinitions) []ast.Stmt {
	return []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "mapper",
							},
						},
						Type: &ast.StructType{
							Fields: generateModelMapperField(fields),
							Incomplete: true,
						},
					},
				},
			},
		},
		&ast.IfStmt{
			Init: &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent("err"),
				},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "json",
							},
							Sel: &ast.Ident{
								Name: "Unmarshal",
							},
						},
						Args: []ast.Expr{
							ast.NewIdent("data"),
							ast.NewIdent("&mapper"),
						},
					},
				},
				Tok: token.DEFINE,
			},
			Cond: &ast.BinaryExpr{
				X: ast.NewIdent("err"),
				Op: token.NEQ,
				Y: ast.NewIdent("nil"),
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							ast.NewIdent("err"),
						},
					},
				},
			},
		},
	}
}

func generateMappingSchemaValidation(t *schema.InputDefinition) []ast.Stmt {
	generateIfStmts := func(fields schema.FieldDefinitions) []ast.Stmt {
		stmts := make([]ast.Stmt, 0, len(fields))

		for _, f := range fields {
			if !f.Type.Nullable {
				stmts = append(stmts, &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.SelectorExpr{
							X: ast.NewIdent("mapper"),
							Sel: ast.NewIdent(toUpperCase(string(f.Name))),
						},
						Op: token.EQL,
						Y: ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "fmt",
											},
											Sel: ast.NewIdent("Errorf"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind: token.STRING,
												Value: fmt.Sprintf("`%s is required`", string(f.Name)),
											},
										},
									},
								},
							},
						},
					},
				})
			}
		}

		return stmts
	}

	return generateIfStmts(t.Fields)
}