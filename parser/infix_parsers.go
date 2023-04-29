package parser

import "jnafolayan/sql-db/ast"

type infixParseFn func(*Parser, ast.Expression) (ast.Expression, error)

func parseInfixExpression(p *Parser, left ast.Expression) (ast.Expression, error) {
	infixExpr := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	op := p.getCurTokenPrecedence()
	p.nextToken()
	right, err := p.parseExpression(op - 1)
	if err != nil {
		return nil, err
	}

	infixExpr.Right = right

	return infixExpr, nil
}
