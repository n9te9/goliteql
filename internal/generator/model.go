package generator

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/lkeix/gg-executor/schema"
)

func generateModelImport() *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"encoding/json"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fmt"`,
				},
			},
		},
	}
}

func generateSelectionSetInput(fields schema.FieldDefinitions) []ast.Decl {
	decls := make([]ast.Decl, 0, len(fields))

	generateTypeSpec := func(args schema.ArgumentDefinitions, operationName string) []ast.Spec {
		var ret []ast.Spec
		var list []*ast.Field

		for i, arg := range args {
			expr := generateExpr(arg.Type)
			list = append(list, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(toUpperCase(string(arg.Name))),
				},
				Tag: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("`json:\"arg%d\"`", i),
				},
				Type: expr,
			})
		}

		ret = append(ret, &ast.TypeSpec{
			Name: ast.NewIdent(toUpperCase(operationName) + "Args"),
			Type: &ast.StructType{
				Fields: &ast.FieldList{
					List: list,
				},
			},
		})

		return ret
	}

	for _, f := range fields {
		if len(f.Arguments) == 0 {
			continue
		}

		decl := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: generateTypeSpec(f.Arguments, string(f.Name)),
		}

		decls = append(decls, decl)
	}

	return decls
}

func generateModelField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldTypeExpr := generateExpr(f.Type)

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeExpr,
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateExpr(fieldType *schema.FieldType) ast.Expr {
	graphQLType := GraphQLType(fieldType.Name)
	if fieldType.Nullable {
		if graphQLType.IsPrimitive() {
			return &ast.StarExpr{
				X: &ast.Ident{
					Name: graphQLType.golangType(),
				},
			}
		} else {
			if fieldType.IsList {
				return &ast.StarExpr{
					X: &ast.ArrayType{
						Elt: generateExpr(fieldType.ListType),
					},
				}
			}

			return &ast.StarExpr{
				X: ast.NewIdent(graphQLType.golangType()),
			}
		}
	} else {
		if graphQLType.IsPrimitive() {
			return &ast.Ident{
				Name: graphQLType.golangType(),
			}
		} else {
			if fieldType.IsList {
				return &ast.ArrayType{
					Elt: generateExpr(fieldType.ListType),
				}
			}

			return ast.NewIdent(graphQLType.golangType())
		}
	}
}

func generateExprForMapper(fieldType *schema.FieldType) ast.Expr {
	graphQLType := GraphQLType(fieldType.Name)
	if graphQLType.IsPrimitive() {
		return &ast.StarExpr{
			X: &ast.Ident{
				Name: graphQLType.golangType(),
			},
		}
	} else {
		if fieldType.IsList {
			return &ast.ArrayType{
				Elt: generateExpr(fieldType.ListType),
			}
		}

		return &ast.StarExpr{
			X: ast.NewIdent(graphQLType.golangType()),
		}
	}
}

func generateModelMapperField(field schema.FieldDefinitions) *ast.FieldList {
	fields := make([]*ast.Field, 0, len(field))

	for _, f := range field {
		fieldTypeIdent := generateExprForMapper(f.Type)

		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{
				{
					Name: toUpperCase(string(f.Name)),
				},
			},
			Type: fieldTypeIdent,
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", string(f.Name)),
			},
		})
	}

	return &ast.FieldList{
		List: fields,
	}
}

func generateInputModelUnmarshalJSON(t *schema.InputDefinition) *ast.FuncDecl {
	var stmts []ast.Stmt
	stmts = append(stmts, generateUnmarshalJSONBody(t.Fields)...)
	stmts = append(stmts, generateMappingSchemaValidation(t)...)
	stmts = append(stmts, generateMapping(t.Fields)...)

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
			List: append(stmts, &ast.ReturnStmt{
				Results: []ast.Expr{
					ast.NewIdent("nil"),
				},
			}),
		},
	}
}

func generateUnmarshalJSONBody(fields schema.FieldDefinitions) []ast.Stmt {
	modelMapperType := generateModelMapperField(fields)

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
							Fields: modelMapperType,
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
				X:  ast.NewIdent("err"),
				Op: token.NEQ,
				Y:  ast.NewIdent("nil"),
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

func generateMapping(fields schema.FieldDefinitions) []ast.Stmt {
	stmts := make([]ast.Stmt, 0, len(fields))

	for _, f := range fields {
		var field ast.Expr
		field = &ast.SelectorExpr{
			X:   ast.NewIdent("mapper"),
			Sel: ast.NewIdent(toUpperCase(string(f.Name))),
		}
		if !f.Type.Nullable {
			field = &ast.StarExpr{
				X: field,
			}
		}

		if f.Type.Nullable && f.Type.IsList {
			field = &ast.UnaryExpr{
				Op: token.AND,
				X:  field,
			}
		}

		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.Ident{
						Name: "t",
					},
					Sel: ast.NewIdent(toUpperCase(string(f.Name))),
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				field,
			},
		})
	}

	return stmts
}

func generateMappingSchemaValidation[T *schema.InputDefinition | *schema.TypeDefinition](t T) []ast.Stmt {
	var schemaDefinition any = t
	var selectorX ast.Expr
	switch schemaDefinition.(type) {
	case *schema.InputDefinition:
		selectorX = ast.NewIdent("mapper")
	case *schema.TypeDefinition:
		selectorX = ast.NewIdent("t")
	}

	generateInputIfStmts := func(fields schema.FieldDefinitions) []ast.Stmt {
		stmts := make([]ast.Stmt, 0, len(fields))

		for _, f := range fields {
			if !f.Type.Nullable {
				stmts = append(stmts, &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.SelectorExpr{
							X:   selectorX,
							Sel: ast.NewIdent(toUpperCase(string(f.Name))),
						},
						Op: token.EQL,
						Y:  ast.NewIdent("nil"),
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
												Kind:  token.STRING,
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

	generateTypeIfStmts := func(fields schema.FieldDefinitions) []ast.Stmt {
		stmts := make([]ast.Stmt, 0, len(fields))

		for _, f := range fields {
			if !f.Type.Nullable {
				stmts = append(stmts, &ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.SelectorExpr{
							X:   selectorX,
							Sel: ast.NewIdent(toUpperCase(string(f.Name))),
						},
						Op: token.EQL,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("nil"),
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "fmt",
											},
											Sel: ast.NewIdent("Errorf"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
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

	var definition any = t
	switch d := definition.(type) {
	case *schema.InputDefinition:
		return generateInputIfStmts(d.Fields)
	case *schema.TypeDefinition:
		return generateTypeIfStmts(d.Fields)
	}

	return []ast.Stmt{}
}
