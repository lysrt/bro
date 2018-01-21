package lexer

type Lexer struct {
	input        string
	position     int  // current position in the input (current char)
	readPosition int  // current reading position in the input (after current char)
	ch           byte // current char
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '!':
		tok = newToken(Bang, l.ch)
	case '=':
		tok = newToken(Equal, l.ch)
	case '/':
		tok = newToken(Slash, l.ch)
	case '\'', '"':
		tok.Type = String
		tok.Literal = l.readString()
	case '<':
		tok = newToken(BAngleBracket, l.ch)
	case '>':
		tok = newToken(FAngleBracket, l.ch)
	case '0':
		tok.Literal, tok.Type = "", EOF
	default:
		if isLetter(l.ch) {
			tok.Type = Identifier
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(Illegal, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tt TokenType, ch byte) Token {
	return Token{Type: tt, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '-'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}
