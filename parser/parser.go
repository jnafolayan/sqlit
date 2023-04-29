package parser

import (
	"errors"
	"fmt"
	"jnafolayan/sql-db/ast"
	"jnafolayan/sql-db/lexer"
	"jnafolayan/sql-db/lib"
	"jnafolayan/sql-db/token"
)

type OperatorPrecedence int

const (
	LOWEST OperatorPrecedence = iota
	ASSIGN
	AND
	OR
	EQUALS
	LT_GT
	SUM
	PRODUCT
	PREFIX
	INDEX
	CALL
)

var precedences = map[token.TokenType]OperatorPrecedence{
	token.EQ:   EQUALS,
	token.N_EQ: EQUALS,
	token.LT:   LT_GT,
	token.GT:   LT_GT,
	token.AND:  AND,
	token.OR:   OR,
}

func getTokenPrecedence(tokenType token.TokenType) OperatorPrecedence {
	op, ok := precedences[tokenType]
	if !ok {
		return LOWEST
	}
	return op
}

type Parser struct {
	lexer             *lexer.Lexer
	it                *lib.Iterator[*token.Token]
	curToken          *token.Token
	peekToken         *token.Token
	OmitErrorLocation bool

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:             l,
		OmitErrorLocation: false,

		prefixParseFns: map[token.TokenType]prefixParseFn{},
		infixParseFns:  map[token.TokenType]infixParseFn{},
	}

	p.registerPrefixFn(token.INT, parseIntegerLiteral)
	p.registerPrefixFn(token.FLOAT, parseFloatLiteral)
	p.registerPrefixFn(token.STRING, parseStringLiteral)

	p.registerInfixFn(token.EQ, parseInfixExpression)
	p.registerInfixFn(token.N_EQ, parseInfixExpression)
	p.registerInfixFn(token.LT, parseInfixExpression)
	p.registerInfixFn(token.GT, parseInfixExpression)
	p.registerInfixFn(token.AND, parseInfixExpression)
	p.registerInfixFn(token.OR, parseInfixExpression)

	return p
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) getPeekTokenPrecedence() OperatorPrecedence {
	peek := p.it.Peek()
	if peek == nil {
		return LOWEST
	}

	return getTokenPrecedence(peek.Type)
}

func (p *Parser) getCurTokenPrecedence() OperatorPrecedence {
	return getTokenPrecedence(p.curToken.Type)
}

func (p *Parser) nextToken() *token.Token {
	p.curToken = p.it.Next()
	p.peekToken = p.it.Peek()
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
			errorPrefix := ""
			if !p.OmitErrorLocation {
				var loc *token.TokenLocation
				if p.curToken != nil {
					loc = p.curToken.Location
				} else {
					// we're at EOF
					lastToken := tokens[len(tokens)-1]
					loc = &token.TokenLocation{
						Line: lastToken.Location.Line,
						Col:  lastToken.Location.Col + len(lastToken.Literal) + 2, // 2 to leave a space between cur token and expected token location
					}
				}
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

	if p.expectPeekToken(token.WHERE) {
		// move to where
		p.nextToken()
		// move to next token
		p.nextToken()
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		stmt.Predicate = expr
	}

	if p.checkPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseCreateTableStatement() (ast.Statement, error) {
	stmt := &ast.CreateTableStatement{}

	if !p.expectPeekToken(token.TABLE) {
		p.nextToken()
		return nil, expectedTokenError(token.TABLE)
	}

	if !p.expectPeekToken(token.IDENTIFIER) {
		p.nextToken()
		p.nextToken()
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

	if len(stmt.Columns) == 0 {
		return nil, ErrEmptyColumnDefinitions
	}

	if p.checkPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseInsertStatement() (ast.Statement, error) {
	stmt := &ast.InsertStatement{}

	if !p.expectPeekToken(token.INTO) {
		p.nextToken()
		return nil, expectedTokenError(token.INTO)
	}

	if !p.expectPeekToken(token.IDENTIFIER) {
		p.nextToken()
		return nil, errors.New("expected table name")
	}

	stmt.Table = p.curToken

	if !p.expectPeekToken(token.LPAREN) {
		p.nextToken()
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
		p.nextToken()
		return nil, expectedTokenError(token.VALUES)
	}

	if !p.expectPeekToken(token.LPAREN) {
		p.nextToken()
		return nil, expectedTokenError(token.LPAREN)
	}

	p.nextToken()
	for p.curToken != nil && !p.checkCurToken(token.RPAREN) {
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}

		stmt.Values = append(stmt.Values, expr)
		p.nextToken()

		if p.checkCurToken(token.COMMA) {
			p.nextToken()
		}
	}

	if len(stmt.Values) != len(stmt.Columns) {
		return nil, errors.New("number of values must match number of columns")
	}

	if p.checkPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpression(precedence OperatorPrecedence) (ast.Expression, error) {
	tok := p.curToken
	prefixFn, ok := p.prefixParseFns[tok.Type]
	if !ok {
		return nil, fmt.Errorf("no prefix parse function for %s", tok.Type)
	}

	leftExpr, err := prefixFn(p)
	if err != nil {
		return nil, err
	}

	for !p.checkPeekToken(token.SEMICOLON) && precedence < p.getPeekTokenPrecedence() {
		infixFn, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return nil, fmt.Errorf("no infix parse function for %s", p.peekToken.Literal)
		}

		p.nextToken()
		leftExpr, err = infixFn(p, leftExpr)
		if err != nil {
			return nil, err
		}
	}

	return leftExpr, nil
}
