package main

import "text/scanner"

type Scanner struct {
	scanner.Scanner
}

func (s *Scanner) NextChar() rune {
	for s.Peek() == rune(' ') || s.Peek() == rune('\t') || s.Peek() == rune('\r') || s.Peek() == rune('\n') {
		s.Next()
	}
	return s.Peek()
}
