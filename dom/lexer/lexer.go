package lexer

type scanFn func(l *Lexer) (Token, scanFn)

type Lexer struct {
	input        string
	position     int  // current position in the input (current char)
	readPosition int  // current reading position in the input (after current char)
	ch           byte // current char
	scan         scanFn
}

func New(input string) *Lexer {
	l := &Lexer{input: input, scan: scanNode}
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
	tok, fn := l.scan(l)
	l.scan = fn
	return tok
}

func scanNode(l *Lexer) (tok Token, fn scanFn) {
	fn = scanNode

	l.skipWhitespace()

	switch l.ch {
	case '!':
		tok = newToken(TokenBang, l.ch)
	case '=':
		tok = newToken(TokenEqual, l.ch)
	case '/':
		tok = newToken(TokenSlash, l.ch)
	case '\'', '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case '<':
		tok = newToken(TokenRBracket, l.ch)
	case '>':
		tok = newToken(TokenLBracket, l.ch)
		fn = scanText
	case '0':
		tok.Literal, tok.Type = "", TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Type = TokenIdent
			tok.Literal = l.readIdentifier()
			return
		} else {
			tok = newToken(TokenIllegal, l.ch)
		}
	}
	l.readChar()
	return
}

func scanText(l *Lexer) (tok Token, fn scanFn) {
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
			tok = Token{Type: TokenText, Literal: l.input[position:l.position]}
			fn = scanNode
			break
		}
		//TODO: detect HTML character
	}
	if l.position-position == whitespace+1 {
		tok, fn = fn(l)
	}
	return
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
