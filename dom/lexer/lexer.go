package lexer

import "fmt"

type lexFn func(l *Lexer) (Token, lexFn)

type Lexer struct {
	input        string
	position     int // current position in the input (current char)
	readPosition int // current reading position in the input (after current char)
	line         int
	linePosition int
	ch           byte // current char
	lex          lexFn
}

func New(input string) *Lexer {
	l := &Lexer{input: input, lex: lexNode}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	if l.ch == '\n' {
		l.line = 0
		l.linePosition = 0
	}
	l.position = l.readPosition
	l.readPosition++
	l.linePosition++
}

func (l *Lexer) NextToken() Token {
	tok, fn := l.lex(l)
	l.lex = fn
	return tok
}

func lexNode(l *Lexer) (tok Token, fn lexFn) {
	fn = lexNode

	l.skipWhitespace()

	switch l.ch {
	case '!':
		tok.Type = TokenBang
		tok.Literal = string(l.ch)
	case '=':
		tok.Type = TokenEqual
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = TokenSlash
		tok.Literal = string(l.ch)
	case '\'', '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case '<':
		tok.Type = TokenRBracket
		tok.Literal = string(l.ch)
	case '>':
		tok.Type = TokenLBracket
		tok.Literal = string(l.ch)
		fn = lexText
	case '0':
		tok.Literal, tok.Type = "", TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Type = TokenIdent
			tok.Literal = l.readIdentifier()
			return
		} else {
			tok.Type = TokenError
			tok.Literal = fmt.Sprintf("illegal character %q", l.ch)
		}
	}
	l.readChar()
	return
}

func lexText(l *Lexer) (tok Token, fn lexFn) {
	position := l.position
	if l.ch == 0 {
		tok.Type = TokenEOF
		return
	}
	whitespace := 0
	for {
		l.readChar()
		if isWhitespace(l.ch) {
			whitespace++
		}
		if l.ch == '<' || l.ch == 0 {
			tok = newToken(l)
			tok.Type = TokenText
			tok.Literal = l.input[position:l.position]
			fn = lexNode
			break
		}
		//TODO: detect HTML character
	}
	if l.position-position == whitespace+1 {
		tok, fn = fn(l)
	}
	return
}

func newToken(l *Lexer) Token {
	return Token{
		Position:     l.position,
		Line:         l.line,
		LinePosition: l.linePosition,
	}
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

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
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
