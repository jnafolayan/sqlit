package parser

import (
	"errors"
	"fmt"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/lib"
	"jnafolayan/sql-db/token"
	"strconv"
)

type Parser struct {
	lexer             *lexer.Lexer
	it                *lib.Iterator[*token.Token]
	curToken          *token.Token
	OmitErrorLocation bool
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lexer:             l,
		OmitErrorLocation: false,
	}
}

func (p *Parser) nextToken() *token.Token {
	p.curToken = p.it.Next()
	return p.curToken
}

func (p *Parser) checkTokenIs(t *token.Token, tokenType token.TokenType) bool {
	if tokenType == token.EOF {
		return t == nil
	}

	return t != nil && t.Type == tokenType
}

func (p *Parser) checkCurToken(tokenType token.TokenType) bool {
	return p.checkTokenIs(p.curToken, tokenType)
}

func (p *Parser) checkPeekToken(tokenType token.TokenType) bool {
	return p.checkTokenIs(p.it.Peek(), tokenType)
}

func (p *Parser) expectPeekToken(tokenType token.TokenType) bool {
	if p.checkPeekToken(tokenType) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) Parse() (*ast.Program, error) {
	tokens := p.lexer.Tokenize()
	p.it = lib.NewIterator(tokens)

	program := &ast.Program{}

	for p.nextToken() != nil {
		stmt, err := p.parseStatement()
		if err != nil {
			loc := p.curToken.Location
			errorPrefix := ""
			if !p.OmitErrorLocation {
				errorPrefix = fmt.Sprintf("line %d, col %d: ", loc.Line, loc.Col)
			}
			return nil, fmt.Errorf("%s%w", errorPrefix, err)
		}

		program.Statements = append(program.Statements, stmt)
		if !p.checkPeekToken(token.EOF) && !p.checkCurToken(token.SEMICOLON) {
			return nil, expectedTokenError(token.SEMICOLON)
		}
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.SELECT:
		return p.parseSelectStatement()
	case token.CREATE:
		return p.parseCreateTableStatement()
	case token.INSERT:
		return p.parseInsertStatement()
	default:
		return nil, fmt.Errorf("invalid keyword %q", p.curToken.Literal)
	}
}

func (p *Parser) parseSelectStatement() (ast.Statement, error) {
	stmt := &ast.SelectStatement{}
	stmt.Columns = []*token.Token{}

	p.nextToken()
	for p.checkCurToken(token.IDENTIFIER) || p.checkCurToken(token.ASTERISK) {
		stmt.Columns = append(stmt.Columns, p.curToken)
		p.nextToken()
		if p.checkCurToken(token.COMMA) {
			p.nextToken()
		}
	}

	if len(stmt.Columns) == 0 {
		return nil, ErrEmptyColumnsList
	}

	if !p.checkCurToken(token.FROM) {
		return nil, expectedTokenError(token.FROM)
	}

	p.nextToken()
	if !p.checkCurToken(token.IDENTIFIER) {
		return nil, errors.New("expected table name")
	}

	stmt.Table = p.curToken

	if p.checkPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseCreateTableStatement() (ast.Statement, error) {
	stmt := &ast.CreateTableStatement{}

	if !p.expectPeekToken(token.TABLE) {
		return nil, expectedTokenError(token.TABLE)
	}

	if !p.expectPeekToken(token.IDENTIFIER) {
		return nil, errors.New("expected table name")
	}

	stmt.Table = p.curToken
	p.nextToken()

	if p.checkCurToken(token.LPAREN) {
		p.nextToken()
		for p.curToken != nil && !p.checkCurToken(token.RPAREN) {
			if !p.checkCurToken(token.IDENTIFIER) {
				return nil, errors.New("expected column name")
			}

			colName := p.curToken
			p.nextToken()
			if !token.IsKeyword(p.curToken) {
				return nil, errors.New("expected column type")
			}

			colType := p.curToken
			columnDef := &ast.ColumnDefinition{
				Name:     colName,
				DataType: colType,
			}
			stmt.Columns = append(stmt.Columns, columnDef)

			p.nextToken()
			if p.checkCurToken(token.COMMA) {
				p.nextToken()
			}
		}

		if !p.checkCurToken(token.RPAREN) {
			return nil, expectedTokenError(token.RPAREN)
		}
	}

	if p.checkPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseInsertStatement() (ast.Statement, error) {
	stmt := &ast.InsertStatement{}

	if !p.expectPeekToken(token.INTO) {
		return nil, expectedTokenError(token.INTO)
	}

	if !p.expectPeekToken(token.IDENTIFIER) {
		return nil, errors.New("expected table name")
	}

	stmt.Table = p.curToken

	if !p.expectPeekToken(token.LPAREN) {
		return nil, expectedTokenError(token.LPAREN)
	}

	p.nextToken()
	for p.curToken != nil && !p.checkCurToken(token.RPAREN) {
		stmt.Columns = append(stmt.Columns, p.curToken)
		p.nextToken()

		if p.checkCurToken(token.COMMA) {
			p.nextToken()
		}
	}

	if len(stmt.Columns) == 0 {
		return nil, ErrEmptyColumnsList
	}

	// TODO(jnafolayan): to check if curToken is RPAREN?

	if !p.expectPeekToken(token.VALUES) {
		return nil, expectedTokenError(token.VALUES)
	}

	if !p.expectPeekToken(token.LPAREN) {
		return nil, expectedTokenError(token.LPAREN)
	}

	p.nextToken()
	for p.curToken != nil && !p.checkCurToken(token.RPAREN) {
		expr := p.parseExpression()
		if expr != nil {
			stmt.Values = append(stmt.Values, expr)
		}
		p.nextToken()

		if p.checkCurToken(token.COMMA) {
			p.nextToken()
		}
	}

	return stmt, nil
}

func (p *Parser) parseExpression() ast.Expression {
	// TODO(jnafolayan): return complex expressions
	switch p.curToken.Type {
	case token.INT:
		v, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
		return &ast.IntegerLiteral{Token: p.curToken, Value: v}
	case token.STRING:
		return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	default:
		return nil
	}
}
