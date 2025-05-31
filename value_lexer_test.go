package goliteql_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/n9te9/goliteql"
)

func TestValueLexer_Lex(t *testing.T) {
	lexer := goliteql.NewValueLexer()

	tests := []struct {
		input    []byte
		expected []*goliteql.ValueLexerToken
	}{
		{
			input: []byte("123"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.INT, Value: []byte("123")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("-123"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.INT, Value: []byte("-123")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("123.56"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.FLOAT, Value: []byte("123.56")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("-123.56"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.FLOAT, Value: []byte("-123.56")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("true"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.BOOL, Value: []byte("true")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("false"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.BOOL, Value: []byte("false")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte("null"),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.NULL, Value: []byte("null")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
		{
			input: []byte(`"Hello, World!"`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.STRING, Value: []byte(`"Hello, World!"`)},
				{Type: goliteql.EOF, Value: nil},
			},
		}, {
			input: []byte(`"-123.56"`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.STRING, Value: []byte(`"-123.56"`)},
				{Type: goliteql.EOF, Value: nil},
			},
		}, {
			input: []byte(`"-123"`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.STRING, Value: []byte(`"-123"`)},
				{Type: goliteql.EOF, Value: nil},
			},
		}, {
			input: []byte(`{key: "value"}`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.STRING, Value: []byte(`"value"`)},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.EOF, Value: nil},
			},
		}, {
			input: []byte(`[{key: "value"}]`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.LBRACKET, Value: []byte("[")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.STRING, Value: []byte(`"value"`)},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.RBRACKET, Value: []byte("]")},
				{Type: goliteql.EOF, Value: nil},
			},
		}, {
			input: []byte(`[{string_key: "value", anotherKey: "anotherValue"}, {int_key: 1234, anotherIntKey: -1234}, {float_key: 1234.56, anotherKey2:-1234.56}, {bool_key: true, anotherBoolKey: false}, {null_key: null, anotherNullKey: null}]`),
			expected: []*goliteql.ValueLexerToken{
				{Type: goliteql.LBRACKET, Value: []byte("[")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`string_key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.STRING, Value: []byte(`"value"`)},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.IDENT, Value: []byte(`anotherKey`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.STRING, Value: []byte(`"anotherValue"`)},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`int_key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.INT, Value: []byte("1234")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.IDENT, Value: []byte(`anotherIntKey`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.INT, Value: []byte("-1234")},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`float_key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.FLOAT, Value: []byte("1234.56")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.IDENT, Value: []byte(`anotherKey2`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.FLOAT, Value: []byte("-1234.56")},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`bool_key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.BOOL, Value: []byte("true")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.IDENT, Value: []byte(`anotherBoolKey`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.BOOL, Value: []byte("false")},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.LBRACE, Value: []byte("{")},
				{Type: goliteql.IDENT, Value: []byte(`null_key`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.NULL, Value: []byte("null")},
				{Type: goliteql.COMMA, Value: []byte(",")},
				{Type: goliteql.IDENT, Value: []byte(`anotherNullKey`)},
				{Type: goliteql.COLON, Value: []byte(":")},
				{Type: goliteql.NULL, Value: []byte("null")},
				{Type: goliteql.RBRACE, Value: []byte("}")},
				{Type: goliteql.RBRACKET, Value: []byte("]")},
				{Type: goliteql.EOF, Value: nil},
			},
		},
	}

	for _, test := range tests {
		tokens, err := lexer.Lex(test.input)
		if err != nil {
			t.Errorf("Lexing error for input %s: %v", test.input, err)
			continue
		}

		if d := cmp.Diff(tokens, test.expected); d != "" {
			t.Errorf("Lexing mismatch for input %s:\n%s", test.input, d)
		}
	}
}
