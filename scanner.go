package main

import (
	"text/scanner"
	"unicode"
)

// Scanner is a wrapper around text/Scanner, adding conveninent methods
type Scanner struct {
	scanner.Scanner
}

// NextChar places the buffer before the next non-whitespace rune, and returns it
func (s *Scanner) NextChar() rune {
	for unicode.IsSpace(s.Peek()) {
		s.Next()
	}
	return s.Peek()
}
