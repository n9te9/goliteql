package ggparser

type SchemaType string

const (
	SchemaTypeName SchemaType = "NAME"
	SchemaTypeInt SchemaType = "INT"
	SchemaTypeString SchemaType = "STRING"
	SchemaTypeSchema SchemaType = "QUERY"
	SchemaTypeMutate SchemaType = "MUTATION"
	SchemaTypeSubscription SchemaType = "SUBSCRIPTION"
	SchemaTypeEOF SchemaType = "EOF"

	SchemaTypeCurlyOpen    SchemaType = "CURLY_OPEN"    // '{'
	SchemaTypeCurlyClose   SchemaType = "CURLY_CLOSE"   // '}'
	SchemaTypeParenOpen    SchemaType = "PAREN_OPEN"    // '('
	SchemaTypeParenClose   SchemaType = "PAREN_CLOSE"   // ')'
	SchemaTypeColon        SchemaType = "COLON"         // ':'
	SchemaTypeAt           SchemaType = "AT"           // '@'
	SchemaTypeComma        SchemaType = "COMMA"        // ','
	SchemaTypeEqual        SchemaType = "EQUAL"        // '='
	SchemaTypeBracketOpen  SchemaType = "BRACKET_OPEN" // '['
	SchemaTypeBracketClose SchemaType = "BRACKET_CLOSE" // ']'
)

type SchemaToken struct {
	SchemaType SchemaType
	Value []byte
	Line int
	Column int
}

func lexSchema(input []byte) ([]*SchemaToken, error) {
	cur := 0
	col, line := 1, 1

	for cur < len(input) {
		switch input[cur] {
		case ' ', '\t':
			cur++
			col++
			continue
		case '\n':
			cur++
			line++
			col = 1
			continue
		}

	}
	return nil, nil
}