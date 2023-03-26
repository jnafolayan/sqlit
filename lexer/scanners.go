package lexer

func (l *Lexer) readIdentifier() string {
	p := l.cursor.position
	for isAlphanum(l.cursor.char) {
		l.readChar()
	}
	return l.source[p:l.cursor.position]
}

// TODO(jnafolayan): this currently does not support exponents
func (l *Lexer) readNumber() string {
	p := l.cursor.position
	for isDigit(l.cursor.char) {
		l.readChar()
	}

	if l.cursor.char == '.' {
		l.readChar()

		for isDigit(l.cursor.char) {
			l.readChar()
		}
	}

	return l.source[p:l.cursor.position]
}

func (l *Lexer) readString() string {
	l.readChar()

	p := l.cursor.position
	for l.cursor.char != '\'' && l.cursor.char != 0 {
		l.readChar()
	}

	if l.cursor.char == '\'' {
		l.readChar()
		return l.source[p : l.cursor.position-1]
	}

	return ""
}
