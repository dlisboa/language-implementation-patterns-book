package main

import (
	"fmt"
	"strings"
)

// page 31, Pattern 2:
// LL(1) Recursive-Descent Lexer

// Implementation
//
// Rules:
// * for each rule r in grammar G, create a method r() with the same name
// * if the rule has a subrule, call the subrule method from within r()
//
// Subrules:
// * alternatives like `rule := a | b | c` become either a switch or if-else
//   sequence for branches a, b or c
// * we test whether the branches a, b or c apply by looking at the next token
// * if a subrule is optional, remove the `default` or `else` part of the branch
//
// Tokens:
// * tokens have a type (a kind)
// * all references to token of type T becomes a call to match(T)
// * match(T) is a helper function/method that consumes a token if T is the
//   current lookahead token
// * if there's a mismatch and the current token is not T, match(T) returns an
//   error
// * we define the types of tokens with each type being descring by an integer
// * we define some meta token types like EOF and INVALID_TOKEN
//
// Lexer:
// * lexical rules are functions/methods in the lexer
// * it has a `Token()` function that produces the next token
// * this function looks at the next char and determines what it could be, then
//   calls the appropriate lexical rule function

// Grammar to be matched (ANTLR syntax):
//
// grammar NestedNameList;
// list     : '[' elements ']' ;       // match bracketed list
// elements : element (',' element)* ; // match comma-separated list
// element  : NAME | list ;            // element is name or nested list
// NAME     : ('a'..'z'|'A'..'Z')+ ;   // NAME is sequence of >=1 lette
//
// example strings that should be matched
// [a,b,c]
// [a,[b,c],d]

type Token struct {
	Type TokenType
	Text string
}

type TokenType int

// Token types
const (
	EOF TokenType = iota
	LBrack
	RBrack
	Name
	Comma
	Equals
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case LBrack:
		return "LBrack"
	case RBrack:
		return "RBrack"
	case Name:
		return "Name"
	case Comma:
		return "Comma"
	case Equals:
		return "Equals"
	default:
		return "Unknown"
	}
}

// Lexer goes through the input rune by rune and produces Tokens. Lexers are
// also called "scanners" or "tokenizers".
type Lexer struct {
	input   string // entire input
	pos     int    // current position index in the input
	current rune   // current rune
	stopped bool   // is the lexer stopped
}

// marks the end of input
var eof = rune(-1)

func NewLexer(input string) *Lexer {
	if input == "" {
		return &Lexer{current: eof}
	}
	// start at first rune
	// convert string to a rune slice so that indexing is rune-based and not byte-based
	current := []rune(input)[0]
	return &Lexer{input: input, current: current}
}

// isLetter is a helper function, only recognizes ASCII letters
func isLetter(r rune) bool {
	return r >= 'a' && r < 'z' || r >= 'A' && r <= 'Z'
}

func (lex *Lexer) Scan() bool {
	return !lex.stopped
}

// returns the Token at the current position
func (lex *Lexer) Next() (Token, error) {
	for lex.current != eof {
		switch lex.current {
		case ' ', '\t', '\n', '\r':
			lex.consume()
			continue
		case ',':
			lex.consume()
			return Token{Type: Comma, Text: ","}, nil
		case '[':
			lex.consume()
			return Token{Type: LBrack, Text: "["}, nil
		case ']':
			lex.consume()
			return Token{Type: RBrack, Text: "]"}, nil
		case '=':
			lex.consume()
			return Token{Type: Equals, Text: "="}, nil
		default:
			if isLetter(lex.current) {
				return lex.name()
			}
			lex.stopped = true
			return Token{}, fmt.Errorf("non-letter character: %c", lex.current)
		}
	}
	lex.stopped = true
	return Token{Type: EOF}, nil
}

// Lexical rule NAME. The string builder accumulates all the consecutive letters
// into a token.
func (lex *Lexer) name() (Token, error) {
	var s strings.Builder
	for isLetter(lex.current) {
		s.WriteRune(lex.current)
		lex.consume()
	}

	return Token{Type: Name, Text: s.String()}, nil
}

// Consume moves the current position forward by one and saves the next current
// rune.
func (lex *Lexer) consume() {
	lex.pos++

	if lex.pos >= len(lex.input) {
		// signals end of input
		lex.current = eof
	} else {
		// saves the next rune
		lex.current = []rune(lex.input)[lex.pos]
	}
}
