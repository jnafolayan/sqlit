package parser

import (
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/token"
	"strconv"
)

type prefixParseFn func(*Parser) (ast.Expression, error)

func parseIdentifier(p *Parser) (ast.Expression, error) {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func parseIntegerLiteral(p *Parser) (ast.Expression, error) {
	v, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		return nil, expectedTokenError(token.INT)
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: v}, nil
}

func parseFloatLiteral(p *Parser) (ast.Expression, error) {
	v, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, expectedTokenError(token.FLOAT)
	}
	return &ast.FloatLiteral{Token: p.curToken, Value: v}, nil
}

func parseStringLiteral(p *Parser) (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}, nil
}
