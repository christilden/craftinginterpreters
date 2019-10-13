package main

import (
	"fmt"
	"strconv"
)


type Parser struct {
	current   Token
	previous  Token
	hadError  bool
	panicMode bool
}

type Precedence uint8

const (
	PREC_NONE Precedence = iota
	PREC_ASSIGNMENT  // =
	PREC_OR          // or
	PREC_AND         // and
	PREC_EQUALITY    // == !=
	PREC_COMPARISON  // <  <= >=
	PREC_TERM        //+ -
	PREC_FACTOR      // * /
	PREC_UNARY       // ! - +
	PREC_CALL        // . () []
	PREC_PRIMARY
)

type (*ParseFn)() void

type ParseRule struct {
	prefix     ParseFn
	infix      ParseFn
	precedence Precedence
}

func currentChunk() *Chunk {
	return compilingChunk
}

func errorAt(token *Token, message *string) {
	if (parser.panicMode) { return }
	parser.panicMode = true

	fmt.Printf(stderr, "[line %d] Error", token.line)

	if (token.tokenType == TOKEN_EOF) {
		fmt.Printf(stderr, " at end")
	} else if (token.tokenType == TOKEN_ERROR) {
		// nothing.
	} else {
		fmt.Printf(stderr, " at '%d %d'", token.length, token.start)
	}

	fmt.Printf(stderr, ": %s\n", message)
	parser.hadError = true
}

func error(message *string) {
	errorAt(&parser.previous, message)
}

func errorAtCurrent(message *string) {
	errorAt(&parser.current, message)
}

var parser Parser
var compilingChunk *Chunk

func compile(source *string, chunk *Chunk) {
	initScanner(source)

	compilingChunk = chunk
	parser.hadError = false
	parser.panicMode = false

	advanceParser()
	expression()
	consume(TOKEN_EOF, "Expect end of expression.")
	endCompiler()
	return !parser.hadError
}

func advanceParser() {
	parser.previous = parser.current

	for {
		parser.current = scanToken()
		if (parser.current.tokenType != TOKEN_ERROR) { break }

		errorAtCurrent(parser.current.start)
	}
}

func consume(tokenType TokenType, message *string) {
	if (parser.current.tokenType == tokenType) {
		advance()
		return
	}

	errorAtCurrent(message)
}

func emitByte(byte uint8) {
	writeChunk(currentChunk(), byte, parser.previous.line)
}

func emitBytes(byte1 uint8, byte2 uint8) {
	emitByte(byte1)
	emitByte(byte2)
}

func emitReturn() {
	emitByte(OP_RETURN)
}

func makeConstant(value Value) {
	var constant int = addConstant(currentChunk(), value)
	UINT8_MAX := ^uint8(0)
	if (constant > UINT8_MAX) {
		error("Too many constants in one chunk.")
		return 0
	}

	return (uint8) constant
}

func emitConstant(value Value) {
	emitBytes(OP_CONSTANT, makeConstant(value))
}

func endCompiler() {
	emitReturn()
}

func grouping() {
	expression()
	consume(TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
}

func numberParser() {
	value, err := strconv.ParseFloat(parser.previous.start, 64)
	emitConstant(value)
}

func unary() {
	var operatorType TokenType = parser.previous.tokenType

	// compiles the operand.
	parsePrecedence(PREC_UNARY)

	// emits the operator instruction.
	switch (operatorType) {
		case TOKEN_MINUS: emitByte(OP_NEGATE)
		default:
			return // unreachable
	}
}

func binary() {
	// remembers the operator
	var operatorType TokenType = parser.previous.tokenType

	// compiles the right operand.
	var rule *ParseRule = getRule(operatorType)
	parsePrecedence((Precedence)(rule.precedence + 1))

	// emits the operator instruction.
	switch (operatorType) {
		case TOKEN_PLUS:         emitByte(OP_ADD)
		case TOKEN_MINUS:        emitByte(OP_SUBTRACT)
		case TOKEN_STAR:         emitByte(OP_MULTIPLY)
		case TOKEN_SLASH:        emitByte(OP_DIVIDE)
		default:
			return // unreachable
	}
}

ParseRule []rules = {
	{ grouping, NULL,    PREC_CALL },       // TOKEN_LEFT_PAREN
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_RIGHT_PAREN
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_LEFT_BRACE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_RIGHT_BRACE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_COMMA
	{ NULL,     NULL,    PREC_CALL },       // TOKEN_DOT
	{ unary,    binary,  PREC_TERM },       // TOKEN_MINUS
	{ NULL,     binary,  PREC_TERM },       // TOKEN_PLUS
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_SEMICOLON
	{ NULL,     binary,  PREC_FACTOR },     // TOKEN_SLASH
	{ NULL,     binary,  PREC_FACTOR },     // TOKEN_STAR
	{ NULL,     binary,  PREC_FACTOR },     // TOKEN_SLASH
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_BANG
	{ NULL,     NULL,    PREC_EQUALITY },   // TOKEN_BANG_EQUAL
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_EQUAL
	{ NULL,     NULL,    PREC_EQUALITY },   // TOKEN_EQUAL_EQUAL
	{ NULL,     NULL,    PREC_COMPARISON }, // TOKEN_GREATER
	{ NULL,     NULL,    PREC_COMPARISON }, // TOKEN_GREATER_EQUAL
	{ NULL,     NULL,    PREC_COMPARISON }, // TOKEN_LESS
	{ NULL,     NULL,    PREC_COMPARISON }, // TOKEN_LESS_EQUAL
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_IDENTIFIER
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_STRING
	{ number,   NULL,    PREC_NONE },       // TOKEN_NUMBER
	{ NULL,     NULL,    PREC_AND },        // TOKEN_AND
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_CLASS
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_ELSE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_FALSE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_FOR
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_FUN
	{ NULL,     NULL,    PREC_OR },         // TOKEN_OR
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_PRINT
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_RETURN
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_SUPER
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_THIS
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_TRUE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_VAR
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_WHILE
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_ERROR
	{ NULL,     NULL,    PREC_NONE },       // TOKEN_EOF
}

func parsePrecedence(Precedence precedence) {

}

func expression() {
	parsePrecedence(PREC_ASSIGNMENT)
}
