package lexer

import (
	"jnafolayan/sql-db/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := "SELECT *, name, age FROM table24 44 20.45 'colors' '' WHERE AND OR;"
	expected := []struct {
		tokenType token.TokenType
		literal   string
	}{
		{token.SELECT, "SELECT"},
		{token.ASTERISK, "*"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "name"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "age"},
		{token.FROM, "FROM"},
		{token.IDENTIFIER, "table24"},
		{token.INT, "44"},
		{token.FLOAT, "20.45"},
		{token.STRING, "colors"},
		{token.STRING, ""},
		{token.WHERE, "WHERE"},
		{token.AND, "AND"},
		{token.OR, "OR"},
		{token.SEMICOLON, ";"},
	}

	l := New(input)
	tokens := l.Tokenize()

	for i, tt := range tokens {
		et := expected[i]

		if tt.Type != et.tokenType {
			t.Errorf("expected %s token type, got %s (%s)", et.tokenType, tt.Type, tt.Literal)
		}

		if tt.Literal != et.literal {
			t.Errorf("expected %s literal, got %s", et.literal, tt.Literal)
		}
	}
}
