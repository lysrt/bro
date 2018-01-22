# HTML Lexer

This package contains a simple HTML lexer. The design of the lexer is inspired by
[Writing An Interpreter In Go](https://interpreterbook.com/) and 
[Lexical Scanning In Go](https://talks.golang.org/2011/lex.slide).

## Goal

This is a learning exercice and not a full-featured HTML5 lexer.
Here is a list of what is implemented:

- [x] tokenize node
- [x] tokenize text
- [ ] tokenize comment
- [ ] tokenize CDATAcomment
- [ ] work on UTF-8 character
- [ ] replace HTML entities on the fly