package token

import (
	"reflect"
	"strings"
)

type TokenType string

type TokenLocation struct {
	Line int
	Col  int
}

type Token struct {
	Type     TokenType
	Literal  string
	Location *TokenLocation
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "identifier"
	SELECT     TokenType = "SELECT"
	FROM       TokenType = "FROM"
	AS         TokenType = "AS"
	TABLE      TokenType = "TABLE"
	CREATE     TokenType = "CREATE"
	INSERT     TokenType = "INSERT"
	INTO       TokenType = "INTO"
	VALUES     TokenType = "VALUES"
	WHERE      TokenType = "WHERE"
	AND        TokenType = "AND"
	OR         TokenType = "OR"
	INT        TokenType = "INT"
	FLOAT      TokenType = "FLOAT"
	TEXT       TokenType = "TEXT"

	STRING TokenType = "STRING"

	// Symbols
	SEMICOLON TokenType = ";"
	ASTERISK  TokenType = "*"
	COMMA     TokenType = ","
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"

	PLUS  TokenType = "+"
	MINUS TokenType = "-"
	EQ    TokenType = "="
	N_EQ  TokenType = "!="
	GT    TokenType = ">"
	LT    TokenType = "<"
)

var keywords = map[string]TokenType{
	"SELECT": SELECT,
	"FROM":   FROM,
	"AS":     AS,
	"TABLE":  TABLE,
	"CREATE": CREATE,
	"INSERT": INSERT,
	"INTO":   INTO,
	"VALUES": VALUES,
	"WHERE":  WHERE,
	"AND":    AND,
	"OR":     OR,
	"INT":    INT,
	"FLOAT":  FLOAT,
	"TEXT":   TEXT,
}

func init() {
	// Add lowercase variants of keywords to allow case insensitive matching
	keys := reflect.ValueOf(keywords).MapKeys()
	for _, k := range keys {
		keywords[strings.ToLower(k.String())] = keywords[k.String()]
	}
}

func LookupIdentifier(str string) TokenType {
	if tt, ok := keywords[str]; ok {
		return tt
	}

	if tt, ok := keywords[strings.ToLower(str)]; ok {
		return tt
	}

	return IDENTIFIER
}

func IsKeyword(t *Token) bool {
	_, ok := keywords[t.Literal]
	return ok
}
