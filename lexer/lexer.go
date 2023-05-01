package lexer

import (
	"jnafolayan/sql-db/token"
)

type cursor struct {
	loc      *token.TokenLocation
	position int
	char     byte
}

type Lexer struct {
	source string
	cursor *cursor
}

func New(source string) *Lexer {
	l := &Lexer{
		source: source,
	}

	return l
}

func (l *Lexer) Cursor() *cursor {
	return l.cursor
}

func (l *Lexer) resetCursor() {
	if l.cursor == nil {
		l.cursor = &cursor{
			position: -1,
			char:     0,
			loc: &token.TokenLocation{
				Line: -1,
				Col:  -1,
			},
		}
	}

	l.cursor.position = -1
	l.cursor.char = 0
	l.cursor.loc.Line = 0
	l.cursor.loc.Col = 0
}

func (l *Lexer) readChar() {
	l.cursor.position++
	l.cursor.loc.Col++
	if l.cursor.position >= len(l.source) {
		l.cursor.char = 0
	} else {
		l.cursor.char = l.source[l.cursor.position]
	}
}

func (l *Lexer) peekChar() byte {
	p := l.cursor.position + 1
	if p >= len(l.source) {
		return 0
	} else {
		return l.source[p]
	}
}

func (l *Lexer) Tokenize() []*token.Token {
	var tokens []*token.Token

	l.resetCursor()
	l.readChar()

	l.cursor.loc.Line = 1

	for l.cursor.char != 0 {
		l.skipWhitespace()

		switch l.cursor.char {
		case ';':
			tokens = append(tokens, createToken(l.cursor, token.SEMICOLON))
		case '<':
			tokens = append(tokens, createToken(l.cursor, token.LT))
		case '>':
			tokens = append(tokens, createToken(l.cursor, token.GT))
		case '=':
			tokens = append(tokens, createToken(l.cursor, token.EQ))
		case '!':
			if l.peekChar() == '=' {
				tokens = append(tokens, createToken(l.cursor, token.N_EQ))
				l.readChar()
			}
		case '*':
			tokens = append(tokens, createToken(l.cursor, token.ASTERISK))
		case ',':
			tokens = append(tokens, createToken(l.cursor, token.COMMA))
		case '(':
			tokens = append(tokens, createToken(l.cursor, token.LPAREN))
		case ')':
			tokens = append(tokens, createToken(l.cursor, token.RPAREN))
		case '\'':
			t := &token.Token{
				Type: token.STRING,
				Location: &token.TokenLocation{
					Line: l.cursor.loc.Line,
					Col:  l.cursor.loc.Col,
				},
			}

			t.Literal = l.readString()
			tokens = append(tokens, t)

			// dont call readChar()
			continue
		default:
			if isLetter(l.cursor.char) {
				t := &token.Token{
					Location: &token.TokenLocation{
						Line: l.cursor.loc.Line,
						Col:  l.cursor.loc.Col,
					},
				}

				t.Literal = l.readIdentifier()
				t.Type = token.LookupIdentifier(t.Literal)
				tokens = append(tokens, t)

				// dont call readChar()
				continue
			} else if isDigit(l.cursor.char) {
				t := &token.Token{
					Location: &token.TokenLocation{
						Line: l.cursor.loc.Line,
						Col:  l.cursor.loc.Col,
					},
				}

				literal, tokenType := l.readNumber()
				t.Literal = literal
				t.Type = tokenType
				tokens = append(tokens, t)

				// dont call readChar()
				continue
			} else {
				tokens = append(tokens, createToken(l.cursor, token.ILLEGAL))
			}
		}

		l.readChar()
	}

	return tokens
}

func (l *Lexer) skipWhitespace() {
	for l.cursor.char == ' ' || l.cursor.char == '\t' || l.cursor.char == '\n' {
		if l.cursor.char == '\n' {
			l.cursor.loc.Line++
			l.cursor.loc.Col = -1
		}

		l.readChar()
	}
}

func createToken(cursor *cursor, tokenType token.TokenType) *token.Token {
	return &token.Token{
		Type:    tokenType,
		Literal: string(cursor.char),
		Location: &token.TokenLocation{
			Line: cursor.loc.Line,
			Col:  cursor.loc.Col,
		},
	}
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isAlphanum(char byte) bool {
	return isLetter(char) || isDigit(char) || char == '_'
}
