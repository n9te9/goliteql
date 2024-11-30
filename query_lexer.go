package ggparser

import (
	"errors"
	"unicode"
)

type QueryType string

type QueryToken struct {
	QueryType QueryType
	Value []byte
	Line int
	Column int
}

const (
	QueryTypeName QueryType = "NAME"
	QueryTypeInt QueryType = "INT"
	QueryTypeString QueryType = "STRING"
	QueryTypeQuery QueryType = "QUERY"
	QueryTypeMutate QueryType = "MUTATION"
	QueryTypeSubscription QueryType = "SUBSCRIPTION"
	QueryTypeEOF QueryType = "EOF"

	QueryTypeCurlyOpen    QueryType = "CURLY_OPEN"    // '{'
	QueryTypeCurlyClose   QueryType = "CURLY_CLOSE"   // '}'
	QueryTypeParenOpen    QueryType = "PAREN_OPEN"    // '('
	QueryTypeParenClose   QueryType = "PAREN_CLOSE"   // ')'
	QueryTypeColon        QueryType = "COLON"         // ':'
	QueryTypeAt           QueryType = "AT"           // '@'
	QueryTypeComma        QueryType = "COMMA"        // ','
	QueryTypeEqual        QueryType = "EQUAL"        // '='
	QueryTypeBracketOpen  QueryType = "BRACKET_OPEN" // '['
	QueryTypeBracketClose QueryType = "BRACKET_CLOSE" // ']'
)

var queryKeywords = map[string]QueryType{
	"query": QueryTypeQuery,
	"mutation": QueryTypeMutate,
	"subscription": QueryTypeSubscription,
}

func newFieldQueryToken(input []byte, cur, col, line int) (*QueryToken, int) {
	start := cur
	for cur < len(input) && unicode.IsLetter(rune(input[cur])) || unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	if tokenQueryType, ok := queryKeywords[string(input[start:cur])]; ok {
		return &QueryToken{QueryType: tokenQueryType, Value: input[start:cur], Column: col, Line: line}, cur
	}
	return &QueryToken{QueryType: QueryTypeName, Value: input[start:cur], Column: col, Line: line}, cur
}

func newIntQueryToken(input []byte, cur, col, line int) (*QueryToken, int) {
	start := cur
	for cur < len(input) && unicode.IsDigit(rune(input[cur])) {
		cur++
	}

	return &QueryToken{QueryType: QueryTypeInt, Value: input[start:cur], Column: col, Line: line}, cur
}

func newStringQueryToken(input []byte, cur, col, line int) (*QueryToken, int, error) {
	start := cur + 1
	cur++
	for cur < len(input) && input[cur] != '"' {
		cur++
	}

	if cur >= len(input) {
		return nil, -1, errors.New("unterminated string")
	}

	return &QueryToken{QueryType: QueryTypeString, Value: input[start:cur], Column: col, Line: line}, cur + 1, nil
}

func newEOFQueryToken(col, line int) *QueryToken {
	return &QueryToken{QueryType: QueryTypeEOF, Value: nil, Column: col, Line: line}
}

var queryPunctuators = map[byte]QueryType{
	'{': QueryTypeCurlyOpen,
	'}': QueryTypeCurlyClose,
	'(': QueryTypeParenOpen,
	')': QueryTypeParenClose,
	':': QueryTypeColon,
	'@': QueryTypeAt,
	',': QueryTypeComma,
	'=': QueryTypeEqual,
	'[': QueryTypeBracketOpen,
	']': QueryTypeBracketClose,
}

func lexQuery(input []byte) ([]*QueryToken, error) {
	tokens := make([]*QueryToken, 0)
	cur := 0
	col, line := 1, 1

	var token *QueryToken
	var err error
	for cur < len(input) {
		switch input[cur] {
		case ' ', '\t':
			col++
			cur++
			continue
		case '\n':
			line++
			col = 1
			cur++
			continue
		}

		switch input[cur] {
		case '{', '}', '(', ')', ':', '@', ',', '=':
			tokens = append(tokens, &QueryToken{QueryType: queryPunctuators[input[cur]], Value: []byte{input[cur]}, Column: col, Line: line})
			cur++
			col++
			continue
		}

		if unicode.IsLetter(rune(input[cur])) {
			token, cur = newFieldQueryToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if unicode.IsDigit(rune(input[cur])) {
			token, cur = newIntQueryToken(input, cur, col, line)
			tokens = append(tokens, token)
			col += len(token.Value)
			continue
		}

		if input[cur] == '"' {
			token, cur, err = newStringQueryToken(input, cur, col, line)
			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
			col += len(token.Value) + 2
			continue
		}
	}

	tokens = append(tokens, newEOFQueryToken(col, line))
	return tokens, nil
}
