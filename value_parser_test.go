package goliteql_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql"
)

func TestValueParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    goliteql.ValueParserExpr
		wantErr bool
	}{
		{
			name:    "simple integer",
			input:   []byte("123"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("123"), TokenType: goliteql.INT},
			wantErr: false,
		},
		{
			name:    "negative integer",
			input:   []byte("-123"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("-123"), TokenType: goliteql.INT},
			wantErr: false,
		},
		{
			name:    "simple float",
			input:   []byte("123.456"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("123.456"), TokenType: goliteql.FLOAT},
			wantErr: false,
		},
		{
			name:    "negative float",
			input:   []byte("-123.456"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("-123.456"), TokenType: goliteql.FLOAT},
			wantErr: false,
		},
		{
			name:    "boolean true",
			input:   []byte("true"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("true"), TokenType: goliteql.BOOL},
			wantErr: false,
		},
		{
			name:    "boolean false",
			input:   []byte("false"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("false"), TokenType: goliteql.BOOL},
			wantErr: false,
		},
		{
			name:    "null value",
			input:   []byte("null"),
			want:    &goliteql.ValueParserLiteral{Value: []byte("null"), TokenType: goliteql.NULL},
			wantErr: false,
		},
		{
			name:    "simple string",
			input:   []byte(`"hello"`),
			want:    &goliteql.ValueParserLiteral{Value: []byte(`"hello"`), TokenType: goliteql.STRING},
			wantErr: false,
		},
		{
			name:    "string with escaped characters",
			input:   []byte(`"hello\nworld"`),
			want:    &goliteql.ValueParserLiteral{Value: []byte(`"hello\nworld"`), TokenType: goliteql.STRING},
			wantErr: false,
		},
		{
			name:    "empty object",
			input:   []byte("{}"),
			want:    &goliteql.ValueParserObject{Fields: map[string]goliteql.ValueParserExpr{}},
			wantErr: false,
		},
		{
			name:  "object with one field",
			input: []byte(`{key: "value"}`),
			want: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key": &goliteql.ValueParserLiteral{Value: []byte(`"value"`), TokenType: goliteql.STRING, IsField: true},
				},
			},
			wantErr: false,
		},
		{
			name:  "object with multiple fields",
			input: []byte(`{key1: "value1", key2: 123}`),
			want: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`"value1"`), TokenType: goliteql.STRING, IsField: true},
					"key2": &goliteql.ValueParserLiteral{Value: []byte("123"), TokenType: goliteql.INT, IsField: true},
				},
			},
			wantErr: false,
		},
		{
			name:  "object with nested object",
			input: []byte(`{key1: "value1", key2: {nestedKey: "nestedValue"}}`),
			want: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`"value1"`), TokenType: goliteql.STRING, IsField: true},
					"key2": &goliteql.ValueParserObject{
						Fields: map[string]goliteql.ValueParserExpr{
							"nestedKey": &goliteql.ValueParserLiteral{Value: []byte(`"nestedValue"`), TokenType: goliteql.STRING, IsField: true},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "object with nested array object",
			input: []byte(`{key1: "value1", key2: [{key: "hogehoge"}, {anotherKey: "anotherValue"}]}`),
			want: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`"value1"`), TokenType: goliteql.STRING, IsField: true},
					"key2": &goliteql.ValueParserArray{
						Items: []goliteql.ValueParserExpr{
							&goliteql.ValueParserObject{
								Fields: map[string]goliteql.ValueParserExpr{
									"key": &goliteql.ValueParserLiteral{Value: []byte(`"hogehoge"`), TokenType: goliteql.STRING, IsField: true},
								},
							},
							&goliteql.ValueParserObject{
								Fields: map[string]goliteql.ValueParserExpr{
									"anotherKey": &goliteql.ValueParserLiteral{Value: []byte(`"anotherValue"`), TokenType: goliteql.STRING, IsField: true},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "array with integers",
			input: []byte(`[1, 2, 3]`),
			want: &goliteql.ValueParserArray{
				Items: []goliteql.ValueParserExpr{
					&goliteql.ValueParserLiteral{Value: []byte("1"), TokenType: goliteql.INT},
					&goliteql.ValueParserLiteral{Value: []byte("2"), TokenType: goliteql.INT},
					&goliteql.ValueParserLiteral{Value: []byte("3"), TokenType: goliteql.INT},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := goliteql.NewValueLexer()
			parser := goliteql.NewValueParser(lexer)

			got, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValueParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if d := cmp.Diff(got, tt.want); d != "" {
				t.Errorf("ValueParser.Parse() mismatch (-got +want):\n%s", d)
			}
		})
	}
}

func TestValueParserLiteral_JSONBytes(t *testing.T) {
	tests := []struct {
		name    string
		v       *goliteql.ValueParserLiteral
		want    []byte
		wantErr error
	}{
		{
			name: "integer literal to JSON bytes",
			v: &goliteql.ValueParserLiteral{
				Value:     []byte("123"),
				TokenType: goliteql.INT,
			},
			want:    []byte("123"),
			wantErr: nil,
		},
		{
			name: "float literal to JSON bytes",
			v: &goliteql.ValueParserLiteral{
				Value:     []byte("123.456"),
				TokenType: goliteql.FLOAT,
			},
			want:    []byte("123.456"),
			wantErr: nil,
		},
		{
			name: "boolean true literal to JSON bytes",
			v: &goliteql.ValueParserLiteral{
				Value:     []byte("true"),
				TokenType: goliteql.BOOL,
			},
			want:    []byte("true"),
			wantErr: nil,
		},
		{
			name: "boolean false literal to JSON bytes",
			v: &goliteql.ValueParserLiteral{
				Value:     []byte("false"),
				TokenType: goliteql.BOOL,
			},
			want:    []byte("false"),
			wantErr: nil,
		},
		{
			name: "null literal to JSON bytes",
			v: &goliteql.ValueParserLiteral{
				Value:     []byte("null"),
				TokenType: goliteql.NULL,
			},
			want:    []byte("null"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.JSONBytes()
			if err != nil && tt.wantErr == nil {
				t.Errorf("ValueParserLiteral.JSONBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("ValueParserLiteral.JSONBytes() expected error %v, got none", tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("ValueParserLiteral.JSONBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestValueParserObject_JSONBytes(t *testing.T) {
	tests := []struct {
		name    string
		v       *goliteql.ValueParserObject
		want    []byte
		wantErr error
	}{
		{
			name: "empty object to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{},
			},
			want:    []byte("{}"),
			wantErr: nil,
		},
		{
			name: "object with one field to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key": &goliteql.ValueParserLiteral{Value: []byte(`value`), TokenType: goliteql.STRING},
				},
			},
			want:    []byte(`{"key":"value"}`),
			wantErr: nil,
		},
		{
			name: "object with multiple fields to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`value1`), TokenType: goliteql.STRING},
					"key2": &goliteql.ValueParserLiteral{Value: []byte("123"), TokenType: goliteql.INT},
				},
			},
			want:    []byte(`{"key1":"value1","key2":123}`),
			wantErr: nil,
		},
		{
			name: "object with nested object to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`value1`), TokenType: goliteql.STRING},
					"key2": &goliteql.ValueParserObject{
						Fields: map[string]goliteql.ValueParserExpr{
							"nestedKey": &goliteql.ValueParserLiteral{Value: []byte(`nestedValue`), TokenType: goliteql.STRING},
						},
					},
				},
			},
			want: []byte(`{"key1":"value1","key2":{"nestedKey":"nestedValue"}}`),
		},
		{
			name: "object with nested array object to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"key1": &goliteql.ValueParserLiteral{Value: []byte(`value1`), TokenType: goliteql.STRING},
					"key2": &goliteql.ValueParserArray{
						Items: []goliteql.ValueParserExpr{
							&goliteql.ValueParserObject{
								Fields: map[string]goliteql.ValueParserExpr{
									"key": &goliteql.ValueParserLiteral{Value: []byte(`hogehoge`), TokenType: goliteql.STRING},
								},
							},
							&goliteql.ValueParserObject{
								Fields: map[string]goliteql.ValueParserExpr{
									"anotherKey": &goliteql.ValueParserLiteral{Value: []byte(`anotherValue`), TokenType: goliteql.STRING},
								},
							},
						},
					},
				},
			},
			want:    []byte(`{"key1":"value1","key2":[{"key":"hogehoge"},{"anotherKey":"anotherValue"}]}`),
			wantErr: nil,
		},
		{
			name: "object with mixed types to JSON bytes",
			v: &goliteql.ValueParserObject{
				Fields: map[string]goliteql.ValueParserExpr{
					"intKey":    &goliteql.ValueParserLiteral{Value: []byte("42"), TokenType: goliteql.INT},
					"floatKey":  &goliteql.ValueParserLiteral{Value: []byte("3.14"), TokenType: goliteql.FLOAT},
					"boolKey":   &goliteql.ValueParserLiteral{Value: []byte("true"), TokenType: goliteql.BOOL},
					"nullKey":   &goliteql.ValueParserLiteral{Value: []byte("null"), TokenType: goliteql.NULL},
					"stringKey": &goliteql.ValueParserLiteral{Value: []byte(`"hello"`), TokenType: goliteql.STRING},
					"nestedObj": &goliteql.ValueParserObject{
						Fields: map[string]goliteql.ValueParserExpr{
							"nestedKey": &goliteql.ValueParserLiteral{Value: []byte(`nestedValue`), TokenType: goliteql.STRING},
						},
					},
				},
			},
			want:    []byte(`{"intKey":42,"floatKey":3.14,"boolKey":true,"nullKey":null,"stringKey":""hello"","nestedObj":{"nestedKey":"nestedValue"}}`),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.v.JSONBytes()
			if err != nil && tt.wantErr == nil {
				t.Errorf("ValueParserObject.JSONBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("ValueParserObject.JSONBytes() expected error %v, got none", tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("ValueParserObject.JSONBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}
