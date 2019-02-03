package main

import "unicode"


type TokenType uint8

const (
	// Single-character tokens.
	TOKEN_LEFT_PAREN TokenType = iota
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FUN
	TOKEN_FOR
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_ERROR
	TOKEN_EOF
)

type Token struct {
	tokenType TokenType
	source    *string
	start     int
	length    int
	line      int
}

type Scanner struct {
	source  *string
	length  int
	start   int
	current int
	line    int
}

func (s Scanner) TokenSource() string {
	return (*s.source)[s.start:s.current]
}

func (s Scanner) Start(advance ...int) rune {
	var pos int = 0
	if (len(advance) > 0) {
		pos = advance[0]
	}
	return (rune)((*s.source)[s.start + pos])
}

func (s Scanner) Current() rune {
	return (rune)((*s.source)[s.current])
}

func (s Scanner) Previous() rune {
	return (rune)((*s.source)[s.current - 1])
}

func (s Scanner) Next() rune {
	return (rune)((*s.source)[s.current + 1])
}

var scanner Scanner

func initScanner(source *string) {
	scanner.source = source
	scanner.length = len(*source)
	scanner.start = 0
	scanner.current = 0
	scanner.line = 1
}

func isAtEnd() bool {
	return scanner.current == (scanner.length)
}

func isAlpha(c rune) bool {
	return (c == '_' || unicode.IsLetter(c))
}

func isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

func makeToken(tokenType TokenType) Token {
	var token Token
	token.tokenType = tokenType
	token.source = scanner.source
	token.start = scanner.start
	token.length = (int)(scanner.current - scanner.start)
	token.line = scanner.line

	return token
}

func errorToken(message *string) Token {
	var token Token
	token.tokenType = TOKEN_ERROR
	token.source = message
	token.start = 0
	token.length = len(*message)
	token.line = scanner.line

	return token
}

func skipWhitespace() {
	for {
		c := peek()
		switch(c) {
			case ' ':
				advance()
			case '\r':
				advance()
			case '\t':
				advance()
			case '\n':
				scanner.line++
				advance()
			case '/':
				if (peekNext() == '/') {
					// comments go until the end of the line.
					for (peek() != '\n' && !isAtEnd()) { advance() }
				} else {
					return
				}
			default:
				return
		}
	}
}

func checkKeyword(start int, length int, rest string, tokenType TokenType) TokenType {
	var pre string = string(scanner.Start())
	if (start > 1) { pre += string(scanner.Start(1)) }
	if (scanner.current - scanner.start == start + length &&
			(pre + rest) == scanner.TokenSource()) {
		return tokenType
	}

	return TOKEN_IDENTIFIER
}

func identifierType() TokenType {
	switch (scanner.Start()) {
		case 'a': return checkKeyword(1, 2, "nd", TOKEN_AND)
		case 'c': return checkKeyword(1, 4, "lass", TOKEN_CLASS)
		case 'e': return checkKeyword(1, 3, "lse", TOKEN_ELSE)
		case 'f':
			if (scanner.current - scanner.start > 1) {
				switch (scanner.Start(1)) {
					case 'a': return checkKeyword(2, 3, "lse", TOKEN_FALSE)
					case 'o': return checkKeyword(2, 1, "r", TOKEN_FOR)
					case 'u': return checkKeyword(2, 1, "n", TOKEN_FUN)
				}
			}
		case 'i': return checkKeyword(1, 1, "f", TOKEN_IF)
		case 'n': return checkKeyword(1, 2, "il", TOKEN_NIL)
		case 'o': return checkKeyword(1, 1, "r", TOKEN_OR)
		case 'p': return checkKeyword(1, 4, "rint", TOKEN_PRINT)
		case 'r': return checkKeyword(1, 5, "eturn", TOKEN_RETURN)
		case 's': return checkKeyword(1, 4, "uper", TOKEN_SUPER)
		case 't':
			if (scanner.current - scanner.start > 1) {
				switch(scanner.Start(1)) {
					case 'h': return checkKeyword(2, 2, "is", TOKEN_THIS)
					case 'r': return checkKeyword(2,2, "ue", TOKEN_TRUE)
				}
			}
		case 'v': return checkKeyword(1, 2, "ar", TOKEN_VAR)
		case 'w': return checkKeyword(1, 4, "hile", TOKEN_WHILE)
	}

	return TOKEN_IDENTIFIER
}

func identifier() Token {
	for (isAlpha(peek()) || isDigit(peek())) { advance() }

	return makeToken(identifierType())
}

func str() Token {
	for (peek() != '"' && !isAtEnd()) {
		if (peek() == '\n') { scanner.line++ }
		advance()
	}

	if (isAtEnd()) {
		error := "Unterminated string."
		return errorToken(&error)
	}

	// the closing ".
	advance()
	return makeToken(TOKEN_STRING)
}

func number() Token {
	for (isDigit(peek())) { advance() }

	// looks for a fractional part.
	if (peek() == '.' && isDigit(peekNext())) {
		// consumes the "."
		advance()

		for (isDigit(peek())) { advance() }
	}

	return makeToken(TOKEN_NUMBER)
}

func advance() rune {
	if (isAtEnd()) { return rune(0) }
	scanner.current++;
	return scanner.Previous()
}

func peek() rune {
	if (isAtEnd()) { return rune(0) }
	return scanner.Current()
}

func peekNext() rune {
	if (isAtEnd()) { return rune(0) }
	return scanner.Next()
}

func match(expected rune) bool {
	if (isAtEnd()) { return false }
	if (scanner.Current() != expected) { return false }

	scanner.current++
	return true
}

func scanToken() Token {
	skipWhitespace()

	scanner.start = scanner.current

	if (isAtEnd()) {
		return makeToken(TOKEN_EOF)
	}

	var c rune = advance()

	if (isAlpha(c)) { return identifier() }
	if (isDigit(c)) { return number() }

	switch(c) {
		case '(': return makeToken(TOKEN_LEFT_PAREN)
		case ')': return makeToken(TOKEN_RIGHT_PAREN)
		case '{': return makeToken(TOKEN_LEFT_BRACE)
		case '}': return makeToken(TOKEN_RIGHT_BRACE)
		case ';': return makeToken(TOKEN_SEMICOLON)
		case ',': return makeToken(TOKEN_COMMA)
		case '.': return makeToken(TOKEN_DOT)
		case '-': return makeToken(TOKEN_MINUS)
		case '+': return makeToken(TOKEN_PLUS)
		case '/': return makeToken(TOKEN_SLASH)
		case '*': return makeToken(TOKEN_STAR)
		case '!':
			if (match('=')) {
				return makeToken(TOKEN_BANG_EQUAL)
			}
			return makeToken(TOKEN_BANG)
		case '=':
			if (match('=')) {
				return makeToken(TOKEN_EQUAL_EQUAL)
			}
			return makeToken(TOKEN_EQUAL)
		case '<':
			if (match('=')) {
				return makeToken(TOKEN_LESS_EQUAL)
			}
			return makeToken(TOKEN_LESS)
		case '>':
			if (match('=')) {
				return makeToken(TOKEN_GREATER_EQUAL)
			}
			return makeToken(TOKEN_GREATER)
		case '"': return str()
	}

	error := "Unexpected character."
	return errorToken(&error)
}
